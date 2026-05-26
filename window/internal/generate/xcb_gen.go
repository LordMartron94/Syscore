package main

import (
	"fmt"
	"os"

	"codegen"
	gocode "codegen/go"
)

const (
	xcbConstantsGenFile = "bindings_constants_gen.go"
	xcbTypesGenFile     = "bindings_types_gen.go"
)

func xcbBindingsGenerate(xprotoXML string, outputDir string) error {
	if _, err := os.Stat(xprotoXML); err != nil {
		return fmt.Errorf("xproto spec: %w", err)
	}

	ir, err := xprotoIRBuild(xprotoXML)
	if err != nil {
		return fmt.Errorf("xproto ir: %w", err)
	}

	if err := xcbConstantsGenerate(outputDir); err != nil {
		return err
	}
	if err := xcbTypesGenerate(outputDir, ir); err != nil {
		return err
	}
	return bindingFilesFormat(outputDir, xcbConstantsGenFile, xcbTypesGenFile)
}

func xcbConstantsGenerate(outputDir string) error {
	elements := bindingFilePreamble("bindings", "linux")

	stringConsts := []struct {
		name  string
		value string
		doc   string
	}{
		{"XCB_ATOM_WM_NAME", "WM_NAME", "XCB_ATOM_WM_NAME is the legacy WM_NAME property atom name."},
		{"XCB_ATOM_NET_WM_NAME", "_NET_WM_NAME", "XCB_ATOM_NET_WM_NAME is the EWMH _NET_WM_NAME property atom name."},
		{"XCB_ATOM_UTF8_STRING", "UTF8_STRING", "XCB_ATOM_UTF8_STRING is the UTF8_STRING type atom name."},
		{"XCB_ATOM_WM_NORMAL_HINTS", "WM_NORMAL_HINTS", "XCB_ATOM_WM_NORMAL_HINTS is the WM_NORMAL_HINTS property atom name."},
		{"XCB_ATOM_WM_SIZE_HINTS", "WM_SIZE_HINTS", "XCB_ATOM_WM_SIZE_HINTS is the WM_SIZE_HINTS type atom name for WM_NORMAL_HINTS values."},
		{"XCB_ATOM_WM_PROTOCOLS", "WM_PROTOCOLS", "XCB_ATOM_WM_PROTOCOLS is the WM_PROTOCOLS property atom name."},
		{"XCB_ATOM_WM_DELETE_WINDOW", "WM_DELETE_WINDOW", "XCB_ATOM_WM_DELETE_WINDOW is the WM_DELETE_WINDOW protocol atom name."},
	}
	specs := make([]gocode.ConstSpec, 0, 32)
	for _, item := range stringConsts {
		specs = append(specs, gocode.ConstSpecNew(item.name, nil, fmt.Sprintf("%q", item.value), item.doc))
	}

	numericConsts := []gocode.ConstSpec{
		gocode.ConstSpecNew("XCB_WINDOW_CLASS_INPUT_OUTPUT", gocode.TypeExprNamedPtr("uint16"), "1", "XCB_WINDOW_CLASS_INPUT_OUTPUT is the InputOutput window class."),
		gocode.ConstSpecNew("XCB_PROP_MODE_REPLACE", gocode.TypeExprNamedPtr("uint8"), "0", "XCB_PROP_MODE_REPLACE replaces the previous property value."),
		gocode.ConstSpecNew("XCB_EVENT_CLIENT_MESSAGE", gocode.TypeExprNamedPtr("uint8"), "33", "XCB_EVENT_CLIENT_MESSAGE is the response type for ClientMessage events."),
		gocode.ConstSpecNew("XCB_EVENT_DESTROY_NOTIFY", gocode.TypeExprNamedPtr("uint8"), "17", "XCB_EVENT_DESTROY_NOTIFY is the response type for DestroyNotify events."),
		gocode.ConstSpecNew("XcbClientMessageEventWindowOffset", nil, "4", "XcbClientMessageEventWindowOffset is the byte offset of window in xcb_client_message_event_t."),
		gocode.ConstSpecNew("XcbClientMessageEventData32Offset", nil, "12", "XcbClientMessageEventData32Offset is the byte offset of data32[0] in xcb_client_message_event_t."),
		gocode.ConstSpecNew("XcbDestroyNotifyEventWindowOffset", nil, "8", "XcbDestroyNotifyEventWindowOffset is the byte offset of window in xcb_destroy_notify_event_t."),
		gocode.ConstSpecNew("XCB_SIZE_HINTS_FLAG_PSIZE", gocode.TypeExprNamedPtr("int32"), "1 << 2", "XCB_SIZE_HINTS_FLAG_PSIZE sets width and height in xcb_size_hints_t."),
		gocode.ConstSpecNew("XCB_SIZE_HINTS_FLAG_PMIN_SIZE", gocode.TypeExprNamedPtr("int32"), "1 << 4", "XCB_SIZE_HINTS_FLAG_PMIN_SIZE sets min_width and min_height."),
		gocode.ConstSpecNew("XCB_SIZE_HINTS_FLAG_PMAX_SIZE", gocode.TypeExprNamedPtr("int32"), "1 << 5", "XCB_SIZE_HINTS_FLAG_PMAX_SIZE sets max_width and max_height."),
		gocode.ConstSpecNew("XcbInternAtomReplyAtomOffset", nil, "8", "XcbInternAtomReplyAtomOffset is the byte offset of atom in xcb_intern_atom_reply_t."),
		gocode.ConstSpecNew("XcbSizeHintsPropertyWordCount", gocode.TypeExprNamedPtr("uint32"), "18", "XcbSizeHintsPropertyWordCount is the xcb_change_property data length for XcbSizeHints."),
	}
	specs = append(specs, numericConsts...)
	elements = append(elements, gocode.FileElementFrom(gocode.DeclConstGroup(specs, "", true)))
	bindingBlankLine(&elements)
	return writeBindingFile(outputDir, xcbConstantsGenFile, elements)
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
