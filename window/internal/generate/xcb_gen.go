package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"codegen"
	gocode "codegen/go"
)

const (
	xcbConstantsGenFile = "bindings_constants_gen.go"
	xcbTypesGenFile     = "bindings_types_gen.go"
)

func xcbBindingsGenerate(xprotoXML string, xcbEventHeader string, xcbICCCMHeader string, xcbEWMHAtomlist string, icccmHTML string, outputDir string) error {
	if _, err := os.Stat(xprotoXML); err != nil {
		return fmt.Errorf("xproto spec: %w", err)
	}
	if _, err := os.Stat(xcbEventHeader); err != nil {
		return fmt.Errorf("xcb event header: %w", err)
	}
	if _, err := os.Stat(xcbICCCMHeader); err != nil {
		return fmt.Errorf("xcb icccm header: %w", err)
	}
	if _, err := os.Stat(xcbEWMHAtomlist); err != nil {
		return fmt.Errorf("xcb ewmh atomlist: %w", err)
	}
	if _, err := os.Stat(icccmHTML); err != nil {
		return fmt.Errorf("icccm html: %w", err)
	}

	ir, err := xprotoIRBuild(xprotoXML)
	if err != nil {
		return fmt.Errorf("xproto ir: %w", err)
	}

	eventResponseTypeMask, err := cHeaderDefineValueGet(xcbEventHeader, "XCB_EVENT_RESPONSE_TYPE_MASK")
	if err != nil {
		return err
	}
	sizeHintPSize, err := cHeaderEnumItemValueGet(xcbICCCMHeader, "XCB_ICCCM_SIZE_HINT_P_SIZE")
	if err != nil {
		return err
	}
	sizeHintPMinSize, err := cHeaderEnumItemValueGet(xcbICCCMHeader, "XCB_ICCCM_SIZE_HINT_P_MIN_SIZE")
	if err != nil {
		return err
	}
	sizeHintPMaxSize, err := cHeaderEnumItemValueGet(xcbICCCMHeader, "XCB_ICCCM_SIZE_HINT_P_MAX_SIZE")
	if err != nil {
		return err
	}
	ewmhAtomNames, err := xcbEWMHAtomNamesParse(xcbEWMHAtomlist)
	if err != nil {
		return err
	}
	wmDeleteWindow, err := xcbAtomNameFromConstName("XCB_ATOM_WM_DELETE_WINDOW")
	if err != nil {
		return err
	}
	if err := textTokenRequire(icccmHTML, wmDeleteWindow); err != nil {
		return err
	}

	if err := xcbConstantsGenerate(outputDir, ir, ewmhAtomNames, wmDeleteWindow, eventResponseTypeMask, sizeHintPSize, sizeHintPMinSize, sizeHintPMaxSize); err != nil {
		return err
	}
	if err := xcbTypesGenerate(outputDir, ir); err != nil {
		return err
	}
	return bindingFilesFormat(outputDir, xcbConstantsGenFile, xcbTypesGenFile)
}

func xcbConstantsGenerate(outputDir string, ir *xprotoIR, ewmhAtomNames map[string]struct{}, wmDeleteWindow string, eventResponseTypeMaskValueExpr string, sizeHintPSizeValueExpr string, sizeHintPMinSizeValueExpr string, sizeHintPMaxSizeValueExpr string) error {
	elements := bindingFilePreamble("bindings", "linux")
	elements = bindingElementsWithImport(elements, "unsafe")

	wmName, err := xcbAtomNameFromXprotoRequire(ir, "WM_NAME")
	if err != nil {
		return err
	}
	wmNormalHints, err := xcbAtomNameFromXprotoRequire(ir, "WM_NORMAL_HINTS")
	if err != nil {
		return err
	}
	wmSizeHints, err := xcbAtomNameFromXprotoRequire(ir, "WM_SIZE_HINTS")
	if err != nil {
		return err
	}
	netWMName, err := xcbAtomNameFromEWMHRequire(ewmhAtomNames, "_NET_WM_NAME")
	if err != nil {
		return err
	}
	utf8String, err := xcbAtomNameFromEWMHRequire(ewmhAtomNames, "UTF8_STRING")
	if err != nil {
		return err
	}
	wmProtocols, err := xcbAtomNameFromEWMHRequire(ewmhAtomNames, "WM_PROTOCOLS")
	if err != nil {
		return err
	}

	stringConsts := []struct {
		name  string
		value string
		doc   string
	}{
		{"XCB_ATOM_WM_NAME", wmName, "XCB_ATOM_WM_NAME is the legacy WM_NAME property atom name."},
		{"XCB_ATOM_NET_WM_NAME", netWMName, "XCB_ATOM_NET_WM_NAME is the EWMH _NET_WM_NAME property atom name."},
		{"XCB_ATOM_UTF8_STRING", utf8String, "XCB_ATOM_UTF8_STRING is the UTF8_STRING type atom name."},
		{"XCB_ATOM_WM_NORMAL_HINTS", wmNormalHints, "XCB_ATOM_WM_NORMAL_HINTS is the WM_NORMAL_HINTS property atom name."},
		{"XCB_ATOM_WM_SIZE_HINTS", wmSizeHints, "XCB_ATOM_WM_SIZE_HINTS is the WM_SIZE_HINTS type atom name for WM_NORMAL_HINTS values."},
		{"XCB_ATOM_WM_PROTOCOLS", wmProtocols, "XCB_ATOM_WM_PROTOCOLS is the WM_PROTOCOLS property atom name."},
		{"XCB_ATOM_WM_DELETE_WINDOW", wmDeleteWindow, "XCB_ATOM_WM_DELETE_WINDOW is the WM_DELETE_WINDOW protocol atom name."},
	}
	specs := make([]gocode.ConstSpec, 0, 32)
	for _, item := range stringConsts {
		specs = append(specs, gocode.ConstSpecNew(item.name, nil, fmt.Sprintf("%q", item.value), item.doc))
	}

	windowClassInputOutput, err := xcbEnumValueRequireInt(ir, "WindowClass", "InputOutput")
	if err != nil {
		return err
	}
	propModeReplace, err := xcbEnumValueRequireInt(ir, "PropMode", "Replace")
	if err != nil {
		return err
	}
	eventClientMessage, err := xcbEventNumberRequire(ir, "ClientMessage")
	if err != nil {
		return err
	}
	eventDestroyNotify, err := xcbEventNumberRequire(ir, "DestroyNotify")
	if err != nil {
		return err
	}
	eventConfigureNotify, err := xcbEventNumberRequire(ir, "ConfigureNotify")
	if err != nil {
		return err
	}
	cwEventMask, err := xcbEnumValueRequireInt(ir, "CW", "EventMask")
	if err != nil {
		return err
	}
	eventMaskStructureNotify, err := xcbEnumValueRequireInt(ir, "EventMask", "StructureNotify")
	if err != nil {
		return err
	}

	numericConsts := []gocode.ConstSpec{
		gocode.ConstSpecNew(
			"XCB_WINDOW_CLASS_INPUT_OUTPUT",
			gocode.TypeExprNamedPtr("uint16"),
			fmt.Sprintf("%d", windowClassInputOutput),
			"XCB_WINDOW_CLASS_INPUT_OUTPUT is the InputOutput window class.",
		),
		gocode.ConstSpecNew(
			"XCB_PROP_MODE_REPLACE",
			gocode.TypeExprNamedPtr("uint8"),
			fmt.Sprintf("%d", propModeReplace),
			"XCB_PROP_MODE_REPLACE replaces the previous property value.",
		),
		gocode.ConstSpecNew(
			"XCB_EVENT_CLIENT_MESSAGE",
			gocode.TypeExprNamedPtr("uint8"),
			fmt.Sprintf("%d", eventClientMessage),
			"XCB_EVENT_CLIENT_MESSAGE is the response type for ClientMessage events.",
		),
		gocode.ConstSpecNew(
			"XCB_EVENT_DESTROY_NOTIFY",
			gocode.TypeExprNamedPtr("uint8"),
			fmt.Sprintf("%d", eventDestroyNotify),
			"XCB_EVENT_DESTROY_NOTIFY is the response type for DestroyNotify events.",
		),
		gocode.ConstSpecNew(
			"XCB_EVENT_CONFIGURE_NOTIFY",
			gocode.TypeExprNamedPtr("uint8"),
			fmt.Sprintf("%d", eventConfigureNotify),
			"XCB_EVENT_CONFIGURE_NOTIFY is the response type for ConfigureNotify events.",
		),
		gocode.ConstSpecNew(
			"XCB_CW_EVENT_MASK",
			gocode.TypeExprNamedPtr("uint32"),
			fmt.Sprintf("%d", cwEventMask),
			"XCB_CW_EVENT_MASK is the CreateWindow/ChangeWindowAttributes value mask bit for event_mask.",
		),
		gocode.ConstSpecNew(
			"XCB_EVENT_MASK_STRUCTURE_NOTIFY",
			gocode.TypeExprNamedPtr("uint32"),
			fmt.Sprintf("%d", eventMaskStructureNotify),
			"XCB_EVENT_MASK_STRUCTURE_NOTIFY selects ConfigureNotify on the window.",
		),
		gocode.ConstSpecNew("XCB_RESPONSE_TYPE_EVENT_CODE_MASK", gocode.TypeExprNamedPtr("uint8"), eventResponseTypeMaskValueExpr, "XCB_RESPONSE_TYPE_EVENT_CODE_MASK masks the core event code from response_type (X11 core protocol send_event flag is bit 7)."),
		gocode.ConstSpecNew("XCB_RESPONSE_TYPE_SENT_EVENT_FLAG", gocode.TypeExprNamedPtr("uint8"), "^XCB_RESPONSE_TYPE_EVENT_CODE_MASK", "XCB_RESPONSE_TYPE_SENT_EVENT_FLAG is set in response_type when the X server marks send_event=true (bit 7)."),
		gocode.ConstSpecNew("XcbClientMessageEventWindowOffset", nil, "unsafe.Offsetof(XcbClientMessageEvent{}.Window)", "XcbClientMessageEventWindowOffset is the byte offset of window in xcb_client_message_event_t."),
		gocode.ConstSpecNew("XcbClientMessageEventData32Offset", nil, "unsafe.Offsetof(XcbClientMessageEvent{}.Data)", "XcbClientMessageEventData32Offset is the byte offset of data32[0] in xcb_client_message_event_t."),
		gocode.ConstSpecNew("XcbDestroyNotifyEventWindowOffset", nil, "unsafe.Offsetof(XcbDestroyNotifyEvent{}.Window)", "XcbDestroyNotifyEventWindowOffset is the byte offset of window in xcb_destroy_notify_event_t."),
		gocode.ConstSpecNew("XCB_SIZE_HINTS_FLAG_PSIZE", gocode.TypeExprNamedPtr("int32"), sizeHintPSizeValueExpr, "XCB_SIZE_HINTS_FLAG_PSIZE sets width and height in xcb_size_hints_t."),
		gocode.ConstSpecNew("XCB_SIZE_HINTS_FLAG_PMIN_SIZE", gocode.TypeExprNamedPtr("int32"), sizeHintPMinSizeValueExpr, "XCB_SIZE_HINTS_FLAG_PMIN_SIZE sets min_width and min_height."),
		gocode.ConstSpecNew("XCB_SIZE_HINTS_FLAG_PMAX_SIZE", gocode.TypeExprNamedPtr("int32"), sizeHintPMaxSizeValueExpr, "XCB_SIZE_HINTS_FLAG_PMAX_SIZE sets max_width and max_height."),
		gocode.ConstSpecNew("XcbInternAtomReplyAtomOffset", nil, "unsafe.Offsetof(struct{ ResponseType uint8; Pad0 uint8; Sequence uint16; Length uint32; Atom XcbAtomT }{}.Atom)", "XcbInternAtomReplyAtomOffset is the byte offset of atom in xcb_intern_atom_reply_t."),
		gocode.ConstSpecNew("XcbGetGeometryReplyWidthOffset", nil, "unsafe.Offsetof(struct{ ResponseType uint8; Depth uint8; Sequence uint16; Length uint32; Root XcbWindowT; X int16; Y int16; Width uint16; Height uint16 }{}.Width)", "XcbGetGeometryReplyWidthOffset is the byte offset of width in xcb_get_geometry_reply_t."),
		gocode.ConstSpecNew("XcbGetGeometryReplyHeightOffset", nil, "unsafe.Offsetof(struct{ ResponseType uint8; Depth uint8; Sequence uint16; Length uint32; Root XcbWindowT; X int16; Y int16; Width uint16; Height uint16 }{}.Height)", "XcbGetGeometryReplyHeightOffset is the byte offset of height in xcb_get_geometry_reply_t."),
		gocode.ConstSpecNew("XcbSizeHintsPropertyWordCount", gocode.TypeExprNamedPtr("uint32"), "uint32(unsafe.Sizeof(XcbSizeHints{}) / unsafe.Sizeof(int32(0)))", "XcbSizeHintsPropertyWordCount is the xcb_change_property data length for XcbSizeHints."),
	}
	specs = append(specs, numericConsts...)
	elements = append(elements, gocode.FileElementFrom(gocode.DeclConstGroup(specs, "", true)))
	bindingBlankLine(&elements)
	return writeBindingFile(outputDir, xcbConstantsGenFile, elements)
}

func xcbEnumValueRequireInt(ir *xprotoIR, enumName string, itemName string) (int64, error) {
	value, ok := xprotoEnumValueGet(ir, enumName, itemName)
	if !ok {
		return 0, fmt.Errorf("xproto: enum %s.%s not found in spec", enumName, itemName)
	}
	return value, nil
}

func xcbEventNumberRequire(ir *xprotoIR, eventName string) (uint32, error) {
	value, ok := xprotoEventNumberGet(ir, eventName)
	if !ok {
		return 0, fmt.Errorf("xproto: event %s number not found in spec", eventName)
	}
	return value, nil
}

func cHeaderDefineValueGet(headerPath string, defineName string) (string, error) {
	f, err := os.Open(headerPath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	prefix := "#define " + defineName
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !strings.HasPrefix(line, prefix) {
			continue
		}
		rest := strings.TrimSpace(strings.TrimPrefix(line, prefix))
		if rest == "" {
			return "", fmt.Errorf("header %s: %s has empty define value", headerPath, defineName)
		}
		return rest, nil
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return "", fmt.Errorf("header %s: missing #define %s", headerPath, defineName)
}

func cHeaderEnumItemValueGet(headerPath string, itemName string) (string, error) {
	f, err := os.Open(headerPath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !strings.Contains(line, itemName) {
			continue
		}
		equalsIndex := strings.Index(line, "=")
		if equalsIndex < 0 {
			continue
		}
		left := strings.TrimSpace(line[:equalsIndex])
		if left != itemName {
			continue
		}
		right := strings.TrimSpace(line[equalsIndex+1:])
		right = strings.TrimSuffix(right, ",")
		right = strings.TrimSpace(right)
		if right == "" {
			return "", fmt.Errorf("header %s: enum item %s has empty value", headerPath, itemName)
		}
		return right, nil
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return "", fmt.Errorf("header %s: missing enum item %s", headerPath, itemName)
}

func xcbEWMHAtomNamesParse(atomlistPath string) (map[string]struct{}, error) {
	content, err := os.ReadFile(atomlistPath)
	if err != nil {
		return nil, err
	}
	trimmed := strings.TrimSpace(string(content))
	trimmed = strings.TrimPrefix(trimmed, "DO(")
	trimmed = strings.TrimSuffix(trimmed, ")")
	parts := strings.Split(trimmed, ",")
	names := make(map[string]struct{}, len(parts))
	for _, part := range parts {
		name := strings.TrimSpace(part)
		if name == "" {
			continue
		}
		names[name] = struct{}{}
	}
	return names, nil
}

func textTokenRequire(path string, token string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if !strings.Contains(string(content), token) {
		return fmt.Errorf("%s: token %q not found", path, token)
	}
	return nil
}

func xcbAtomNameFromXprotoRequire(ir *xprotoIR, itemName string) (string, error) {
	if _, ok := xprotoEnumValueGet(ir, "Atom", itemName); !ok {
		return "", fmt.Errorf("xproto: Atom.%s not found in spec", itemName)
	}
	return itemName, nil
}

func xcbAtomNameFromEWMHRequire(atomNames map[string]struct{}, itemName string) (string, error) {
	if _, ok := atomNames[itemName]; !ok {
		return "", fmt.Errorf("ewmh atomlist: %s not found", itemName)
	}
	return itemName, nil
}

func xcbAtomNameFromConstName(constName string) (string, error) {
	const prefix = "XCB_ATOM_"
	if !strings.HasPrefix(constName, prefix) {
		return "", fmt.Errorf("atom const name %q missing %q prefix", constName, prefix)
	}
	return strings.TrimPrefix(constName, prefix), nil
}

func xcbTypesGenerate(outputDir string, ir *xprotoIR) error {
	elements := bindingFilePreamble("bindings", "linux")
	elements = bindingElementsWithImport(elements, "unsafe")

	opaqueTypes := []struct {
		name       string
		underlying string
	}{
		{"XcbConnectionT", "uintptr"},
		{"XcbWindowT", "uint32"},
		{"XcbVoidCookieT", "uint32"},
		{"XcbAtomT", "uint32"},
		{"XcbInternAtomCookieT", "uint32"},
		{"XcbGetGeometryCookieT", "uint32"},
	}
	for _, item := range opaqueTypes {
		elements = append(elements, gocode.FileElementFrom(
			gocode.DeclTypeDefined(item.name, gocode.TypeExprNamed(item.underlying), false, fmt.Sprintf("%s is an XCB wire/handle type.", item.name)),
		))
		bindingBlankLine(&elements)
	}

	elements = append(elements, gocode.FileElementFrom(
		gocode.DeclTypeStruct("XcbScreenIteratorT", []gocode.StructFieldDecl{
			gocode.StructFieldTypeDoc("Data", gocode.TypeExprNamed("uintptr"), "Data points at the current xcb_screen_t."),
			gocode.StructFieldTypeDoc("Rem", gocode.TypeExprNamed("int32"), "Rem is the number of screens remaining."),
			gocode.StructFieldTypeDoc("Index", gocode.TypeExprNamed("int32"), "Index is the screen index."),
		}, "XcbScreenIteratorT is the xcb_screen_iterator_t struct."),
	))
	bindingBlankLine(&elements)

	elements = append(elements, gocode.FileElementFrom(
		gocode.DeclTypeStruct("XcbClientMessageEvent", []gocode.StructFieldDecl{
			gocode.StructFieldTypeDoc("ResponseType", gocode.TypeExprNamed("uint8"), ""),
			gocode.StructFieldTypeDoc("Format", gocode.TypeExprNamed("uint8"), ""),
			gocode.StructFieldTypeDoc("Sequence", gocode.TypeExprNamed("uint16"), ""),
			gocode.StructFieldTypeDoc("Window", gocode.TypeExprNamed("XcbWindowT"), ""),
			gocode.StructFieldTypeDoc("Type", gocode.TypeExprNamed("XcbAtomT"), ""),
			gocode.StructFieldTypeDoc("Data", gocode.TypeExprArray("5", gocode.TypeExprNamed("uint32")), "Data is data32; index 0 is at byte offset 12."),
		}, "XcbClientMessageEvent is the xcb_client_message_event_t wire layout (32 bytes)."),
	))
	bindingBlankLine(&elements)

	elements = append(elements, gocode.FileElementFrom(
		gocode.DeclTypeStruct("XcbDestroyNotifyEvent", []gocode.StructFieldDecl{
			gocode.StructFieldTypeDoc("ResponseType", gocode.TypeExprNamed("uint8"), ""),
			gocode.StructFieldTypeDoc("Pad0", gocode.TypeExprNamed("uint8"), ""),
			gocode.StructFieldTypeDoc("Sequence", gocode.TypeExprNamed("uint16"), ""),
			gocode.StructFieldTypeDoc("Event", gocode.TypeExprNamed("XcbWindowT"), ""),
			gocode.StructFieldTypeDoc("Window", gocode.TypeExprNamed("XcbWindowT"), ""),
		}, "XcbDestroyNotifyEvent is the xcb_destroy_notify_event_t wire layout (12 bytes)."),
	))
	bindingBlankLine(&elements)

	elements = append(elements, gocode.FileElementFrom(
		gocode.DeclTypeStruct("XcbConfigureNotifyEvent", []gocode.StructFieldDecl{
			gocode.StructFieldTypeDoc("ResponseType", gocode.TypeExprNamed("uint8"), ""),
			gocode.StructFieldTypeDoc("Pad0", gocode.TypeExprNamed("uint8"), ""),
			gocode.StructFieldTypeDoc("Sequence", gocode.TypeExprNamed("uint16"), ""),
			gocode.StructFieldTypeDoc("Event", gocode.TypeExprNamed("XcbWindowT"), ""),
			gocode.StructFieldTypeDoc("Window", gocode.TypeExprNamed("XcbWindowT"), ""),
			gocode.StructFieldTypeDoc("AboveSibling", gocode.TypeExprNamed("XcbWindowT"), ""),
			gocode.StructFieldTypeDoc("X", gocode.TypeExprNamed("int16"), ""),
			gocode.StructFieldTypeDoc("Y", gocode.TypeExprNamed("int16"), ""),
			gocode.StructFieldTypeDoc("Width", gocode.TypeExprNamed("uint16"), ""),
			gocode.StructFieldTypeDoc("Height", gocode.TypeExprNamed("uint16"), ""),
			gocode.StructFieldTypeDoc("BorderWidth", gocode.TypeExprNamed("uint16"), ""),
			gocode.StructFieldTypeDoc("OverrideRedirect", gocode.TypeExprNamed("uint8"), ""),
			gocode.StructFieldTypeDoc("Pad1", gocode.TypeExprNamed("uint8"), ""),
		}, "XcbConfigureNotifyEvent is the xcb_configure_notify_event_t wire layout (32 bytes)."),
	))
	bindingBlankLine(&elements)

	sizeHintFields := []string{
		"Flags", "X", "Y", "Width", "Height", "MinWidth", "MinHeight", "MaxWidth", "MaxHeight",
		"WidthInc", "HeightInc", "MinAspectNum", "MinAspectDen", "MaxAspectNum", "MaxAspectDen",
		"BaseWidth", "BaseHeight", "WinGravity",
	}
	fields := make([]gocode.StructFieldDecl, 0, len(sizeHintFields))
	for _, field := range sizeHintFields {
		fields = append(fields, gocode.StructFieldTypeDoc(field, gocode.TypeExprNamed("int32"), ""))
	}
	elements = append(elements, gocode.FileElementFrom(
		gocode.DeclTypeStruct("XcbSizeHints", fields, "XcbSizeHints is the xcb_size_hints_t layout for WM_NORMAL_HINTS (ICCCM, 18x int32)."),
	))
	bindingBlankLine(&elements)
	elements = append(elements, xcbSizeofGuardElements()...)

	pfnSpecs, err := xcbPFNSpecsFromIR(ir, xcbLoaderCommands())
	if err != nil {
		return fmt.Errorf("xcb pfn: %w", err)
	}
	elements, err = bindingPFNElementsAppend(elements, pfnSpecs)
	if err != nil {
		return err
	}

	return writeBindingFile(outputDir, xcbTypesGenFile, elements)
}

func xcbSizeofGuardElements() []codegen.FileElement {
	guards := []struct {
		typ   string
		bytes int
	}{
		{"XcbClientMessageEvent", 32},
		{"XcbDestroyNotifyEvent", 12},
		{"XcbConfigureNotifyEvent", 32},
	}
	elements := make([]codegen.FileElement, 0, len(guards)+1)
	for _, guard := range guards {
		elements = append(elements, gocode.FileElementFrom(
			gocode.DeclVarArray(
				fmt.Sprintf("_xcbSizeofGuard%s", guard.typ),
				fmt.Sprintf("%d - unsafe.Sizeof(%s{})", guard.bytes, guard.typ),
				gocode.TypeExprNamed("byte"),
				fmt.Sprintf("compile-time guard: %s must be %d bytes.", guard.typ, guard.bytes),
			),
		))
	}
	bindingBlankLine(&elements)
	return elements
}
