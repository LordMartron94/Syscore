package main

import (
	"fmt"
	"sort"

	gocode "codegen/go"
)

func xcbPFNSpecsFromIR(ir *xprotoIR, commands []loaderCommandSpec) ([]bindingPFNSpec, error) {
	specs := make([]bindingPFNSpec, 0, len(commands))
	for _, command := range commands {
		spec, err := xcbPFNSpecFromSymbol(ir, command.SymbolName)
		if err != nil {
			return nil, fmt.Errorf("xcb pfn %s: %w", command.SymbolName, err)
		}
		specs = append(specs, spec)
	}
	sort.Slice(specs, func(i, j int) bool { return specs[i].Name < specs[j].Name })
	return specs, nil
}

func xcbPFNSpecFromSymbol(ir *xprotoIR, symbol string) (bindingPFNSpec, error) {
	switch symbol {
	case "free":
		return bindingPFNSpec{
			Name: "PFN_c_free",
			Func: gocode.TypeExprFunc(
				[]gocode.ParamType{{Name: "ptr", Type: gocode.TypeExprNamed("uintptr")}},
				nil,
			),
			Doc: "PFN_c_free maps to libc free(3) for XCB reply buffers.",
		}, nil
	}

	if symbol == "xcb_intern_atom_reply" {
		return bindingPFNSpec{
			Name: "PFN_xcb_intern_atom_reply",
			Func: gocode.TypeExprFunc(
				[]gocode.ParamType{
					{Name: "c", Type: gocode.TypeExprNamed("XcbConnectionT")},
					{Name: "cookie", Type: gocode.TypeExprNamed("XcbInternAtomCookieT")},
					{Name: "e", Type: gocode.TypeExprNamed("uintptr")},
				},
				[]gocode.TypeExpr{gocode.TypeExprNamed("uintptr")},
			),
			Doc: "PFN_xcb_intern_atom_reply maps to xcb_intern_atom_reply.",
		}, nil
	}

	if aux, ok := xcbAuxiliaryPFNSpec(symbol); ok {
		return aux, nil
	}

	requestName, ok := xcbRequestNameFromSymbol(symbol)
	if !ok {
		return bindingPFNSpec{}, fmt.Errorf("unknown symbol")
	}
	req, ok := ir.Requests[requestName]
	if !ok {
		return bindingPFNSpec{}, fmt.Errorf("request %q not in xproto", requestName)
	}
	return xcbRequestPFNSpec(symbol, req), nil
}

func xcbRequestPFNSpec(symbol string, req xprotoRequestIR) bindingPFNSpec {
	params := []gocode.ParamType{{Name: "c", Type: gocode.TypeExprNamed("XcbConnectionT")}}
	for _, field := range req.Fields {
		params = append(params, xcbFieldParam(field)...)
	}
	return bindingPFNSpec{
		Name: "PFN_" + symbol,
		Func: gocode.TypeExprFunc(params, []gocode.TypeExpr{gocode.TypeExprNamed(xcbRequestReturnType(req))}),
		Doc:  fmt.Sprintf("PFN_%s maps to %s.", symbol, symbol),
	}
}

func xcbFieldParam(field xprotoFieldIR) []gocode.ParamType {
	switch field.Kind {
	case "field":
		return []gocode.ParamType{{
			Name: xcbFieldGoName(field.Name),
			Type: gocode.TypeExprNamed(xcbWireTypeGo(field.Type)),
		}}
	case "list":
		return []gocode.ParamType{{
			Name: xcbFieldGoName(field.Name),
			Type: gocode.TypeExprNamed("uintptr"),
		}}
	case "switch":
		return []gocode.ParamType{{
			Name: xcbFieldGoName(field.Name),
			Type: gocode.TypeExprNamed("uintptr"),
		}}
	default:
		return nil
	}
}

func xcbRequestReturnType(req xprotoRequestIR) string {
	if req.HasReply && req.Name == "InternAtom" {
		return "XcbInternAtomCookieT"
	}
	return "XcbVoidCookieT"
}

func xcbWireTypeGo(wireType string) string {
	switch wireType {
	case "CARD8", "BYTE", "BOOL", "KEYCODE":
		return "uint8"
	case "CARD16", "INT16":
		return "uint16"
	case "CARD32", "INT32", "BOOL32", "VISUALID", "TIMESTAMP", "KEYSYM", "KEYCODE32":
		return "uint32"
	case "WINDOW":
		return "XcbWindowT"
	case "ATOM":
		return "XcbAtomT"
	default:
		return "uint32"
	}
}

// xcbAuxiliaryPFNSpec covers libxcb helpers not modeled as xproto requests.
func xcbAuxiliaryPFNSpec(symbol string) (bindingPFNSpec, bool) {
	switch symbol {
	case "xcb_connect":
		return bindingPFNSpec{
			Name: "PFN_xcb_connect",
			Func: gocode.TypeExprFunc(
				[]gocode.ParamType{
					{Name: "display", Type: gocode.TypeExprNamed("uintptr")},
					{Name: "screen", Type: gocode.TypeExprNamed("uintptr")},
				},
				[]gocode.TypeExpr{gocode.TypeExprNamed("XcbConnectionT")},
			),
			Doc: "PFN_xcb_connect maps to xcb_connect.",
		}, true
	case "xcb_get_setup":
		return bindingPFNSpec{
			Name: "PFN_xcb_get_setup",
			Func: gocode.TypeExprFunc(
				[]gocode.ParamType{{Name: "c", Type: gocode.TypeExprNamed("XcbConnectionT")}},
				[]gocode.TypeExpr{gocode.TypeExprNamed("uintptr")},
			),
			Doc: "PFN_xcb_get_setup maps to xcb_get_setup.",
		}, true
	case "xcb_setup_roots_iterator":
		return bindingPFNSpec{
			Name: "PFN_xcb_setup_roots_iterator",
			Func: gocode.TypeExprFunc(
				[]gocode.ParamType{{Name: "setup", Type: gocode.TypeExprNamed("uintptr")}},
				[]gocode.TypeExpr{gocode.TypeExprNamed("XcbScreenIteratorT")},
			),
			Doc: "PFN_xcb_setup_roots_iterator maps to xcb_setup_roots_iterator (libxcb iterator over xproto Setup.roots).",
		}, true
	case "xcb_generate_id":
		return bindingPFNSpec{
			Name: "PFN_xcb_generate_id",
			Func: gocode.TypeExprFunc(
				[]gocode.ParamType{{Name: "c", Type: gocode.TypeExprNamed("XcbConnectionT")}},
				[]gocode.TypeExpr{gocode.TypeExprNamed("XcbWindowT")},
			),
			Doc: "PFN_xcb_generate_id maps to xcb_generate_id.",
		}, true
	case "xcb_flush":
		return bindingPFNSpec{
			Name: "PFN_xcb_flush",
			Func: gocode.TypeExprFunc(
				[]gocode.ParamType{{Name: "c", Type: gocode.TypeExprNamed("XcbConnectionT")}},
				[]gocode.TypeExpr{gocode.TypeExprNamed("int32")},
			),
			Doc: "PFN_xcb_flush maps to xcb_flush.",
		}, true
	case "xcb_disconnect":
		return bindingPFNSpec{
			Name: "PFN_xcb_disconnect",
			Func: gocode.TypeExprFunc(
				[]gocode.ParamType{{Name: "c", Type: gocode.TypeExprNamed("XcbConnectionT")}},
				nil,
			),
			Doc: "PFN_xcb_disconnect maps to xcb_disconnect.",
		}, true
	case "xcb_poll_for_event":
		return bindingPFNSpec{
			Name: "PFN_xcb_poll_for_event",
			Func: gocode.TypeExprFunc(
				[]gocode.ParamType{{Name: "c", Type: gocode.TypeExprNamed("XcbConnectionT")}},
				[]gocode.TypeExpr{gocode.TypeExprNamed("uintptr")},
			),
			Doc: "PFN_xcb_poll_for_event maps to xcb_poll_for_event.",
		}, true
	case "xcb_get_file_descriptor":
		return bindingPFNSpec{
			Name: "PFN_xcb_get_file_descriptor",
			Func: gocode.TypeExprFunc(
				[]gocode.ParamType{{Name: "c", Type: gocode.TypeExprNamed("XcbConnectionT")}},
				[]gocode.TypeExpr{gocode.TypeExprNamed("int32")},
			),
			Doc: "PFN_xcb_get_file_descriptor maps to xcb_get_file_descriptor.",
		}, true
	case "xcb_wait_for_event":
		return bindingPFNSpec{
			Name: "PFN_xcb_wait_for_event",
			Func: gocode.TypeExprFunc(
				[]gocode.ParamType{{Name: "c", Type: gocode.TypeExprNamed("XcbConnectionT")}},
				[]gocode.TypeExpr{gocode.TypeExprNamed("uintptr")},
			),
			Doc: "PFN_xcb_wait_for_event maps to xcb_wait_for_event.",
		}, true
	default:
		return bindingPFNSpec{}, false
	}
}
