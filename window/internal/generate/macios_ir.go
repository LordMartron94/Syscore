package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type maciosEnumMemberIR struct {
	Name   string
	GoName string
	Value  string
	Doc    string
}

type maciosEnumIR struct {
	NativeName string
	EnumName   string
	Members    []maciosEnumMemberIR
}

type maciosBindingsIR struct {
	Enums     []maciosEnumIR
	Classes   []maciosStringConstIR
	Selectors []maciosStringConstIR
	Strings   []maciosStringConstIR
	Extras    []maciosEnumMemberIR
}

type maciosStringConstIR struct {
	GoName string
	Value  string
	Doc    string
}

var maciosClassWhitelist = []maciosStringConstIR{
	{GoName: "ObjCClassNSApplication", Doc: "ObjCClassNSApplication is the NSApplication class name."},
	{GoName: "ObjCClassNSWindow", Doc: "ObjCClassNSWindow is the NSWindow class name."},
	{GoName: "ObjCClassNSString", Doc: "ObjCClassNSString is the NSString class name."},
	{GoName: "ObjCClassNSDate", Doc: "ObjCClassNSDate is the NSDate class name."},
}

var maciosSelectorWhitelist = []maciosStringConstIR{
	{GoName: "ObjCSelSharedApplication", Doc: "ObjCSelSharedApplication is +[NSApplication sharedApplication]."},
	{GoName: "ObjCSelAlloc", Doc: "ObjCSelAlloc is +alloc."},
	{GoName: "ObjCSelInitWithContentRectStyleMaskBackingDefer", Doc: "ObjCSelInitWithContentRectStyleMaskBackingDefer is NSWindow designated initializer."},
	{GoName: "ObjCSelMakeKeyAndOrderFront", Doc: "ObjCSelMakeKeyAndOrderFront shows and keys a window."},
	{GoName: "ObjCSelStringWithUTF8String", Doc: "ObjCSelStringWithUTF8String is +[NSString stringWithUTF8String:]."},
	{GoName: "ObjCSelSetTitle", Doc: "ObjCSelSetTitle is -setTitle:."},
	{GoName: "ObjCSelClose", Doc: "ObjCSelClose is -close."},
	{GoName: "ObjCSelInit", Doc: "ObjCSelInit is -init."},
	{GoName: "ObjCSelSetDelegate", Doc: "ObjCSelSetDelegate is -setDelegate:."},
	{GoName: "ObjCSelWindowShouldClose", Doc: "ObjCSelWindowShouldClose is NSWindowDelegate -windowShouldClose:."},
	{GoName: "ObjCSelNextEventMatchingMask", Doc: "ObjCSelNextEventMatchingMask is the NSApplication event pump selector."},
	{GoName: "ObjCSelDistantPast", Doc: "ObjCSelDistantPast is +[NSDate distantPast]."},
}

var maciosStringWhitelist = []maciosStringConstIR{
	{GoName: "NSRunLoopModeDefault", Doc: "NSRunLoopModeDefault is NSDefaultRunLoopMode."},
}

var maciosExtraConstants = []maciosEnumMemberIR{}

var (
	maciosNativeAttrPattern = regexp.MustCompile(`\[Native\s*\(\s*"([^"]+)"\s*\)\]`)
	maciosEnumDeclPattern   = regexp.MustCompile(`public enum (\w+)`)
	maciosEnumMemberPattern = regexp.MustCompile(`^(\w+)(?:\s*=\s*(.+))?,?\s*$`)
	maciosExportAttrPattern = regexp.MustCompile(`\[Export\s*\(\s*"([^"]+)"`)
	maciosFieldAttrPattern  = regexp.MustCompile(`\[Field\s*\(\s*"([^"]+)"`)
)

func maciosBindingsIRBuild(enumsPath string, appkitSourcePath string, foundationEnumsPath string) (maciosBindingsIR, error) {
	content, err := os.ReadFile(enumsPath)
	if err != nil {
		return maciosBindingsIR{}, fmt.Errorf("read macios enums: %w", err)
	}
	appkitSourceContent, err := os.ReadFile(appkitSourcePath)
	if err != nil {
		return maciosBindingsIR{}, fmt.Errorf("read appkit source: %w", err)
	}
	foundationEnumsContent, err := os.ReadFile(foundationEnumsPath)
	if err != nil {
		return maciosBindingsIR{}, fmt.Errorf("read foundation enums: %w", err)
	}

	ir := maciosBindingsIR{
		Extras: append([]maciosEnumMemberIR(nil), maciosExtraConstants...),
	}
	for _, item := range maciosClassWhitelist {
		resolved := item
		if resolved.Value == "" {
			resolved.Value = strings.TrimPrefix(item.GoName, "ObjCClass")
		}
		ir.Classes = append(ir.Classes, resolved)
	}
	exportSelectors := maciosExportSelectorSetParse(appkitSourceContent)
	for _, item := range maciosSelectorWhitelist {
		resolved := item
		selectorBase := maciosSelectorBaseFromGoName(item.GoName)
		resolvedValue, err := maciosSelectorResolve(exportSelectors, selectorBase)
		if err != nil {
			return maciosBindingsIR{}, err
		}
		resolved.Value = resolvedValue
		ir.Selectors = append(ir.Selectors, resolved)
	}
	foundationFieldByMember := maciosEnumFieldByMemberParse(foundationEnumsContent, "NSRunLoopMode")
	for _, item := range maciosStringWhitelist {
		resolved := item
		switch item.GoName {
		case "NSRunLoopModeDefault":
			value, ok := foundationFieldByMember["Default"]
			if !ok {
				return maciosBindingsIR{}, fmt.Errorf("Foundation.Enums NSRunLoopMode.Default [Field] not found")
			}
			resolved.Value = value
		default:
			if resolved.Value == "" {
				return maciosBindingsIR{}, fmt.Errorf("macios string constant %s is missing source mapping", item.GoName)
			}
		}
		ir.Strings = append(ir.Strings, resolved)
	}

	enums, err := maciosEnumsParse(content)
	if err != nil {
		return maciosBindingsIR{}, err
	}

	for _, enumIR := range enums {
		goPrefix := maciosEnumGoPrefix(enumIR)
		parsed := maciosEnumIR{
			NativeName: enumIR.NativeName,
			EnumName:   enumIR.EnumName,
		}
		for _, member := range enumIR.Members {
			value, err := maciosEnumValueNormalize(member.Value)
			if err != nil {
				continue
			}
			parsed.Members = append(parsed.Members, maciosEnumMemberIR{
				GoName: goPrefix + member.Name,
				Value:  value,
				Doc:    goPrefix + member.Name + " from dotnet/macios generated bindings.",
			})
		}
		if len(parsed.Members) > 0 {
			ir.Enums = append(ir.Enums, parsed)
		}
	}

	return ir, nil
}

func maciosEnumGoPrefix(enumIR maciosEnumIR) string {
	if enumIR.NativeName != "" {
		return enumIR.NativeName
	}
	return enumIR.EnumName
}

func maciosEnumsParse(content []byte) ([]maciosEnumIR, error) {
	lines := strings.Split(string(content), "\n")
	result := make([]maciosEnumIR, 0, 64)

	var pendingNative string
	var current *maciosEnumIR
	var enumNextImplicitValue int64

	flush := func() {
		if current == nil {
			return
		}
		result = append(result, *current)
		current = nil
	}

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if matches := maciosNativeAttrPattern.FindStringSubmatch(trimmed); len(matches) == 2 {
			flush()
			pendingNative = matches[1]
			continue
		}
		if matches := maciosEnumDeclPattern.FindStringSubmatch(trimmed); len(matches) == 2 {
			flush()
			current = &maciosEnumIR{
				NativeName: pendingNative,
				EnumName:   matches[1],
			}
			enumNextImplicitValue = 0
			pendingNative = ""
			continue
		}
		if current == nil {
			continue
		}
		if trimmed == "}" {
			flush()
			continue
		}
		if matches := maciosEnumMemberPattern.FindStringSubmatch(trimmed); len(matches) == 3 {
			value := ""
			if matches[2] != "" {
				value = strings.TrimSpace(matches[2])
				if idx := strings.Index(value, "//"); idx >= 0 {
					value = strings.TrimSpace(value[:idx])
				}
				value = strings.TrimSuffix(value, ",")
				if numericValue, ok := maciosEnumNumericValueTryParse(value); ok {
					enumNextImplicitValue = numericValue + 1
				}
			} else {
				value = strconv.FormatInt(enumNextImplicitValue, 10)
				enumNextImplicitValue++
			}
			current.Members = append(current.Members, maciosEnumMemberIR{
				Name:  matches[1],
				Value: value,
			})
		}
	}
	flush()
	return result, nil
}

func maciosEnumNumericValueTryParse(raw string) (int64, bool) {
	normalized := strings.TrimSpace(raw)
	normalized = strings.TrimPrefix(normalized, "(")
	normalized = strings.TrimSuffix(normalized, ")")
	normalized = strings.TrimSpace(normalized)
	if strings.HasPrefix(normalized, "0x") || strings.HasPrefix(normalized, "0X") {
		value, err := strconv.ParseInt(normalized[2:], 16, 64)
		return value, err == nil
	}
	value, err := strconv.ParseInt(normalized, 10, 64)
	return value, err == nil
}

func maciosEnumValueNormalize(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	raw = strings.TrimSuffix(raw, ",")
	raw = strings.ReplaceAll(raw, "(ulong)", "")
	raw = strings.ReplaceAll(raw, "unchecked(", "")
	raw = strings.ReplaceAll(raw, ")", "")
	raw = strings.TrimSpace(raw)
	if strings.Contains(raw, "UInt64.MaxValue") {
		return "^uint64(0)", nil
	}
	if strings.Contains(raw, "<<") {
		parts := strings.Split(raw, "<<")
		if len(parts) != 2 {
			return "", fmt.Errorf("unsupported shift expression: %s", raw)
		}
		base, err := maciosEnumValueNormalize(strings.TrimSpace(parts[0]))
		if err != nil {
			return "", err
		}
		shift := strings.TrimSpace(parts[1])
		shift = strings.Trim(shift, "()")
		if _, err := strconv.Atoi(shift); err != nil {
			return "", fmt.Errorf("unsupported shift operand: %s", shift)
		}
		return base + " << " + shift, nil
	}
	if strings.HasPrefix(raw, "0x") || strings.HasPrefix(raw, "0X") {
		return raw, nil
	}
	if _, err := strconv.ParseUint(raw, 10, 64); err == nil {
		return raw, nil
	}
	if strings.HasPrefix(raw, "1UL") || strings.HasPrefix(raw, "1u") {
		return "1", nil
	}
	return "", fmt.Errorf("unsupported enum value: %s", raw)
}

func maciosExportSelectorSetParse(content []byte) map[string]struct{} {
	lines := strings.Split(string(content), "\n")
	selectorSet := make(map[string]struct{})
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		matches := maciosExportAttrPattern.FindStringSubmatch(trimmed)
		if len(matches) != 2 {
			continue
		}
		selectorSet[matches[1]] = struct{}{}
	}
	return selectorSet
}

func maciosSelectorBaseFromGoName(goName string) string {
	base := strings.TrimPrefix(goName, "ObjCSel")
	if base == "" {
		return ""
	}
	return strings.ToLower(base[:1]) + base[1:]
}

func maciosSelectorResolve(selectorSet map[string]struct{}, selectorBase string) (string, error) {
	if selectorBase == "" {
		return "", fmt.Errorf("empty selector base")
	}
	if _, ok := selectorSet[selectorBase]; ok {
		return selectorBase, nil
	}
	candidate := ""
	bestColonCount := -1
	for selector := range selectorSet {
		if !strings.HasPrefix(selector, selectorBase) {
			continue
		}
		if len(selector) <= len(selectorBase) {
			continue
		}
		next := selector[len(selectorBase)]
		if next != ':' {
			continue
		}
		if selector == selectorBase+":" {
			return selector, nil
		}
		colonCount := strings.Count(selector, ":")
		if colonCount > bestColonCount {
			candidate = selector
			bestColonCount = colonCount
			continue
		}
		if colonCount == bestColonCount && selector != candidate {
			return "", fmt.Errorf("ambiguous selector base %q: %q and %q", selectorBase, candidate, selector)
		}
	}
	if candidate != "" {
		return candidate, nil
	}
	return selectorBase, nil
}

func maciosEnumFieldByMemberParse(content []byte, enumName string) map[string]string {
	lines := strings.Split(string(content), "\n")
	result := make(map[string]string)
	inEnum := false
	pendingField := ""
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if !inEnum {
			if strings.Contains(trimmed, "public enum "+enumName) {
				inEnum = true
			}
			continue
		}
		if strings.HasPrefix(trimmed, "}") {
			return result
		}
		if matches := maciosFieldAttrPattern.FindStringSubmatch(trimmed); len(matches) == 2 {
			pendingField = matches[1]
			continue
		}
		matches := maciosEnumMemberPattern.FindStringSubmatch(trimmed)
		if len(matches) != 3 {
			continue
		}
		memberName := matches[1]
		if pendingField != "" {
			result[memberName] = pendingField
			pendingField = ""
		}
	}
	return result
}
