package main

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"

	gocode "codegen/go"
)

var waylandVarargPFNSpecs = map[string]bindingPFNSpec{
	"wl_proxy_marshal": {
		Name: "PFN_wl_proxy_marshal",
		Func: gocode.TypeExprFunc(
			[]gocode.ParamType{
				{Name: "proxy", Type: gocode.TypeExprNamed("WlProxy")},
				{Name: "opcode", Type: gocode.TypeExprNamed("uint32")},
				{Name: "a0", Type: gocode.TypeExprNamed("uintptr")},
				{Name: "a1", Type: gocode.TypeExprNamed("uintptr")},
				{Name: "a2", Type: gocode.TypeExprNamed("uintptr")},
				{Name: "a3", Type: gocode.TypeExprNamed("uintptr")},
				{Name: "a4", Type: gocode.TypeExprNamed("uintptr")},
				{Name: "a5", Type: gocode.TypeExprNamed("uintptr")},
			},
			[]gocode.TypeExpr{gocode.TypeExprNamed("uintptr")},
		),
		Doc: "PFN_wl_proxy_marshal maps to wl_proxy_marshal (variadic args passed as uintptr slots).",
	},
	"wl_proxy_marshal_constructor": {
		Name: "PFN_wl_proxy_marshal_constructor",
		Func: gocode.TypeExprFunc(
			[]gocode.ParamType{
				{Name: "proxy", Type: gocode.TypeExprNamed("WlProxy")},
				{Name: "opcode", Type: gocode.TypeExprNamed("uint32")},
				{Name: "iface", Type: gocode.TypeExprNamed("uintptr")},
				{Name: "userdata", Type: gocode.TypeExprNamed("uintptr")},
			},
			[]gocode.TypeExpr{gocode.TypeExprNamed("WlProxy")},
		),
		Doc: "PFN_wl_proxy_marshal_constructor maps to wl_proxy_marshal_constructor (subset signature used by syscore).",
	},
	"wl_proxy_marshal_constructor_versioned": {
		Name: "PFN_wl_proxy_marshal_constructor_versioned",
		Func: gocode.TypeExprFunc(
			[]gocode.ParamType{
				{Name: "proxy", Type: gocode.TypeExprNamed("WlProxy")},
				{Name: "opcode", Type: gocode.TypeExprNamed("uint32")},
				{Name: "iface", Type: gocode.TypeExprNamed("uintptr")},
				{Name: "version", Type: gocode.TypeExprNamed("uint32")},
				{Name: "name", Type: gocode.TypeExprNamed("uint32")},
				{Name: "userdata", Type: gocode.TypeExprNamed("uintptr")},
			},
			[]gocode.TypeExpr{gocode.TypeExprNamed("WlProxy")},
		),
		Doc: "PFN_wl_proxy_marshal_constructor_versioned maps to wl_proxy_marshal_constructor_versioned (subset signature used by syscore).",
	},
}

func waylandPFNSpecsFromVendorHeader(headerPath string, commands []loaderCommandSpec) ([]bindingPFNSpec, error) {
	content, err := os.ReadFile(headerPath)
	if err != nil {
		return nil, fmt.Errorf("read vendor header: %w", err)
	}
	header := string(content)

	specs := make([]bindingPFNSpec, 0, len(commands))
	for _, command := range commands {
		if spec, ok := waylandVarargPFNSpecs[command.SymbolName]; ok {
			specs = append(specs, spec)
			continue
		}
		decl, err := waylandVendorFuncDeclExtract(header, command.SymbolName)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", command.SymbolName, err)
		}
		spec, err := waylandPFNSpecFromDecl(command.SymbolName, decl)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", command.SymbolName, err)
		}
		specs = append(specs, spec)
	}
	sort.Slice(specs, func(i, j int) bool { return specs[i].Name < specs[j].Name })
	return specs, nil
}

func waylandVendorFuncDeclExtract(header string, symbol string) (string, error) {
	pattern := regexp.MustCompile(`(?m)^[ \t]*` + regexp.QuoteMeta(symbol) + `\s*\(`)
	loc := pattern.FindStringIndex(header)
	if loc == nil {
		return "", fmt.Errorf("declaration not found in vendor wayland-client-core.h")
	}
	openParen := loc[1] - 1
	closeParen, ok := waylandMatchClosingParen(header, openParen)
	if !ok {
		return "", fmt.Errorf("unterminated declaration")
	}
	end := closeParen + 1
	if end < len(header) && header[end] == ';' {
		end++
	}
	start := 0
	if prev := strings.LastIndex(header[:loc[0]], ";"); prev >= 0 {
		start = prev + 1
	}
	return strings.TrimSpace(header[start:end]), nil
}

func waylandMatchClosingParen(header string, openParen int) (int, bool) {
	if openParen < 0 || openParen >= len(header) || header[openParen] != '(' {
		return 0, false
	}
	depth := 0
	for i := openParen; i < len(header); i++ {
		switch header[i] {
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 {
				return i, true
			}
		}
	}
	return 0, false
}

func waylandPFNSpecFromDecl(symbol string, decl string) (bindingPFNSpec, error) {
	open := strings.Index(decl, symbol+"(")
	if open < 0 {
		return bindingPFNSpec{}, fmt.Errorf("invalid declaration")
	}
	close := strings.LastIndex(decl, ")")
	if close < 0 || close <= open {
		return bindingPFNSpec{}, fmt.Errorf("invalid declaration")
	}
	returnType := waylandDeclReturnType(decl[:open], symbol)

	paramsRaw := strings.TrimSpace(decl[open+len(symbol)+1 : close])
	paramsRaw = strings.ReplaceAll(paramsRaw, "...", "")
	params := waylandCParamsParse(paramsRaw)

	goParams := make([]gocode.ParamType, 0, len(params))
	for i, param := range params {
		name := waylandGoParamName(param.name, i)
		goParams = append(goParams, gocode.ParamType{
			Name: name,
			Type: gocode.TypeExprNamed(waylandCTypeGo(param.typ)),
		})
	}

	goReturns := waylandCReturnGo(returnType)
	return bindingPFNSpec{
		Name: "PFN_" + symbol,
		Func: gocode.TypeExprFunc(goParams, goReturns),
		Doc:  fmt.Sprintf("PFN_%s maps to %s.", symbol, symbol),
	}, nil
}

type waylandCParam struct {
	typ  string
	name string
}

func waylandCParamsParse(paramsRaw string) []waylandCParam {
	paramsRaw = strings.TrimSpace(paramsRaw)
	if paramsRaw == "" || paramsRaw == "void" {
		return nil
	}
	parts := waylandCSplitParams(paramsRaw)
	params := make([]waylandCParam, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" || part == "void" {
			continue
		}
		if part == "..." {
			continue
		}
		name, typ := waylandCParamNameType(part)
		params = append(params, waylandCParam{typ: typ, name: name})
	}
	return params
}

func waylandCSplitParams(paramsRaw string) []string {
	var parts []string
	depth := 0
	start := 0
	for i, r := range paramsRaw {
		switch r {
		case '(':
			depth++
		case ')':
			depth--
		case ',':
			if depth == 0 {
				parts = append(parts, paramsRaw[start:i])
				start = i + 1
			}
		}
	}
	parts = append(parts, paramsRaw[start:])
	return parts
}

func waylandDeclReturnType(prefix string, symbol string) string {
	lines := strings.Split(prefix, "\n")
	filtered := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.Contains(line, symbol) {
			continue
		}
		filtered = append(filtered, line)
	}
	return strings.Join(filtered, " ")
}

func waylandCParamNameType(part string) (string, string) {
	part = strings.TrimSpace(part)
	if part == "" {
		return "", "uintptr"
	}
	if strings.Contains(part, "**") {
		return "listener", "void **"
	}
	if idx := strings.LastIndex(part, "*"); idx >= 0 {
		after := strings.TrimSpace(part[idx+1:])
		if after != "" && !strings.Contains(after, " ") {
			typ := strings.TrimSpace(part[:idx+1])
			name := strings.TrimPrefix(after, "*")
			return name, waylandCTypeNormalize(typ)
		}
	}
	i := strings.LastIndex(part, " ")
	if i < 0 {
		return "", waylandCTypeNormalize(part)
	}
	name := strings.TrimSpace(part[i+1:])
	name = strings.TrimPrefix(name, "*")
	return name, waylandCTypeNormalize(part[:i])
}

func waylandCTypeNormalize(typ string) string {
	typ = strings.TrimSpace(typ)
	typ = strings.ReplaceAll(typ, "const ", "")
	typ = strings.ReplaceAll(typ, "struct ", "")
	typ = strings.ReplaceAll(typ, "  ", " ")
	return strings.TrimSpace(typ)
}

func waylandCTypeGo(typ string) string {
	typ = waylandCTypeNormalize(typ)
	switch typ {
	case "wl_display *", "wl_display*", "struct wl_display *":
		return "WlDisplay"
	case "wl_proxy *", "wl_proxy*", "struct wl_proxy *":
		return "WlProxy"
	case "uint32_t", "uint32":
		return "uint32"
	case "int", "int32_t":
		return "int32"
	case "char *", "char*", "const char *", "const char*":
		return "uintptr"
	case "wl_interface *", "wl_interface*":
		return "uintptr"
	case "void **", "void**":
		return "uintptr"
	case "void *", "void*":
		return "uintptr"
	default:
		if strings.HasSuffix(typ, "*") {
			return "uintptr"
		}
		return "uintptr"
	}
}

func waylandGoParamName(name string, index int) string {
	name = strings.TrimPrefix(strings.TrimSpace(name), "*")
	if name == "" || name == "interface" {
		switch index {
		case 0:
			return "proxy"
		case 1:
			return "opcode"
		case 2:
			return "iface"
		case 3:
			return "userdata"
		default:
			return fmt.Sprintf("arg%d", index)
		}
	}
	if name == "implementation" {
		return "listener"
	}
	return name
}

func waylandCReturnGo(returnType string) []gocode.TypeExpr {
	returnType = waylandCTypeNormalize(returnType)
	if returnType == "" || returnType == "void" {
		return nil
	}
	switch returnType {
	case "wl_display *", "wl_display*", "struct wl_display *":
		return []gocode.TypeExpr{gocode.TypeExprNamed("WlDisplay")}
	case "wl_proxy *", "wl_proxy*", "struct wl_proxy *":
		return []gocode.TypeExpr{gocode.TypeExprNamed("WlProxy")}
	case "int", "int32_t":
		return []gocode.TypeExpr{gocode.TypeExprNamed("int32")}
	case "uint32_t":
		return []gocode.TypeExpr{gocode.TypeExprNamed("uint32")}
	default:
		if strings.HasSuffix(returnType, "*") {
			return []gocode.TypeExpr{gocode.TypeExprNamed("uintptr")}
		}
		return []gocode.TypeExpr{gocode.TypeExprNamed("uintptr")}
	}
}
