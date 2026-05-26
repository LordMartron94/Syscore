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
	{GoName: "ObjCClassNSApplication", Value: "NSApplication", Doc: "ObjCClassNSApplication is the NSApplication class name."},
	{GoName: "ObjCClassNSWindow", Value: "NSWindow", Doc: "ObjCClassNSWindow is the NSWindow class name."},
	{GoName: "ObjCClassNSString", Value: "NSString", Doc: "ObjCClassNSString is the NSString class name."},
	{GoName: "ObjCClassNSDate", Value: "NSDate", Doc: "ObjCClassNSDate is the NSDate class name."},
}

var maciosSelectorWhitelist = []maciosStringConstIR{
	{GoName: "ObjCSelSharedApplication", Value: "sharedApplication", Doc: "ObjCSelSharedApplication is +[NSApplication sharedApplication]."},
	{GoName: "ObjCSelAlloc", Value: "alloc", Doc: "ObjCSelAlloc is +alloc."},
	{GoName: "ObjCSelInitWithContentRectStyleMaskBackingDefer", Value: "initWithContentRect:styleMask:backing:defer:", Doc: "ObjCSelInitWithContentRectStyleMaskBackingDefer is NSWindow designated initializer."},
	{GoName: "ObjCSelMakeKeyAndOrderFront", Value: "makeKeyAndOrderFront:", Doc: "ObjCSelMakeKeyAndOrderFront shows and keys a window."},
	{GoName: "ObjCSelStringWithUTF8String", Value: "stringWithUTF8String:", Doc: "ObjCSelStringWithUTF8String is +[NSString stringWithUTF8String:]."},
	{GoName: "ObjCSelSetTitle", Value: "setTitle:", Doc: "ObjCSelSetTitle is -setTitle:."},
	{GoName: "ObjCSelClose", Value: "close", Doc: "ObjCSelClose is -close."},
	{GoName: "ObjCSelInit", Value: "init", Doc: "ObjCSelInit is -init."},
	{GoName: "ObjCSelSetDelegate", Value: "setDelegate:", Doc: "ObjCSelSetDelegate is -setDelegate:."},
	{GoName: "ObjCSelWindowShouldClose", Value: "windowShouldClose:", Doc: "ObjCSelWindowShouldClose is NSWindowDelegate -windowShouldClose:."},
	{GoName: "ObjCSelNextEventMatchingMask", Value: "nextEventMatchingMask:untilDate:inMode:dequeue:", Doc: "ObjCSelNextEventMatchingMask is the NSApplication event pump selector."},
	{GoName: "ObjCSelDistantPast", Value: "distantPast", Doc: "ObjCSelDistantPast is +[NSDate distantPast]."},
}

var maciosStringWhitelist = []maciosStringConstIR{
	{GoName: "NSRunLoopModeDefault", Value: "NSDefaultRunLoopMode", Doc: "NSRunLoopModeDefault is NSDefaultRunLoopMode."},
	{GoName: "ObjCTypeEncodingWindowShouldClose", Value: "c24@0:8@16", Doc: "ObjCTypeEncodingWindowShouldClose is the encoding for -windowShouldClose:."},
}

var maciosExtraConstants = []maciosEnumMemberIR{
	{
		GoName: "NSEventMaskAny",
		Value:  "0xFFFFFFFFFFFFFFFF",
		Doc:    "NSEventMaskAny matches all event types for nextEventMatchingMask.",
	},
}

var (
	maciosNativeAttrPattern = regexp.MustCompile(`\[Native\s*\(\s*"([^"]+)"\s*\)\]`)
	maciosEnumDeclPattern   = regexp.MustCompile(`public enum (\w+)`)
	maciosEnumMemberPattern = regexp.MustCompile(`^(\w+)\s*=\s*(.+),?\s*$`)
)

func maciosBindingsIRBuild(enumsPath string) (maciosBindingsIR, error) {
	content, err := os.ReadFile(enumsPath)
	if err != nil {
		return maciosBindingsIR{}, fmt.Errorf("read macios enums: %w", err)
	}

	ir := maciosBindingsIR{
		Classes:   append([]maciosStringConstIR(nil), maciosClassWhitelist...),
		Selectors: append([]maciosStringConstIR(nil), maciosSelectorWhitelist...),
		Strings:   append([]maciosStringConstIR(nil), maciosStringWhitelist...),
		Extras:    append([]maciosEnumMemberIR(nil), maciosExtraConstants...),
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
			value := strings.TrimSpace(matches[2])
			if idx := strings.Index(value, "//"); idx >= 0 {
				value = strings.TrimSpace(value[:idx])
			}
			value = strings.TrimSuffix(value, ",")
			current.Members = append(current.Members, maciosEnumMemberIR{
				Name:  matches[1],
				Value: value,
			})
		}
	}
	flush()
	return result, nil
}

func maciosEnumValueNormalize(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	raw = strings.TrimSuffix(raw, ",")
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
