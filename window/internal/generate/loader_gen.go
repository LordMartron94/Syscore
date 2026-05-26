package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"codegen"
	gocode "codegen/go"
)

const loaderCommandsGenFile = "loader_commands_gen.go"

type loaderCommandSpec struct {
	FieldName  string
	SymbolName string
}

type loaderGenSpec struct {
	OutputDir   string
	BuildTag    string
	PackageName string
	FuncName    string
	Imports     []string
	Commands    []loaderCommandSpec
}

func loaderCommandsGenerate(spec loaderGenSpec) error {
	commands := make([]loaderCommandSpec, len(spec.Commands))
	copy(commands, spec.Commands)
	sort.Slice(commands, func(i, j int) bool { return commands[i].FieldName < commands[j].FieldName })

	inits := make([]codegen.Expr, 0, len(commands))
	for _, command := range commands {
		inits = append(inits, gocode.ExprCompositeLit("loadutil.CommandMapping", []gocode.FieldInit{
			{
				Name:  "Target",
				Value: gocode.ExprAddressOf(gocode.ExprSelector(gocode.ExprIdent("commands"), command.FieldName)),
			},
			{Name: "Name", Value: gocode.ExprStringLit(command.SymbolName)},
		}))
	}

	elements := loaderGenHeaderElements(spec.BuildTag)
	elements = append(elements, gocode.FileElementFrom(gocode.DeclPackage(spec.PackageName)))
	bindingBlankLine(&elements)
	imports := spec.Imports
	if len(imports) > 0 {
		elements = append(elements, gocode.FileElementFrom(gocode.DeclImportBlock(imports...)))
	}
	bindingBlankLine(&elements)

	body := gocode.StmtBlock(
		gocode.StmtReturn(gocode.ExprSliceCompositeLit("loadutil.CommandMapping", inits)),
	)
	elements = append(elements, gocode.FileElementFrom(
		gocode.DeclFunc(spec.FuncName, []gocode.ParamType{
			{Name: "commands", Type: gocode.TypeExprPointer(gocode.TypeExprNamed(loaderCommandsStructName(spec.FuncName)))},
		}, []gocode.TypeExpr{gocode.TypeExprSlice(gocode.TypeExprNamed("loadutil.CommandMapping"))}, body),
	))

	file := gocode.DeclFile(elements...)
	content, err := gocode.GoFileRenderWithOptions(file, gocode.RenderOptionsEnumBindings())
	if err != nil {
		return fmt.Errorf("render loader: %w", err)
	}

	absPath := filepath.Join(spec.OutputDir, loaderCommandsGenFile)
	if err := os.WriteFile(absPath, []byte(content), 0o644); err != nil {
		return fmt.Errorf("write loader: %w", err)
	}
	return bindingFilesFormat(spec.OutputDir, loaderCommandsGenFile)
}

func loaderCommandsStructName(funcName string) string {
	switch funcName {
	case "waylandCommandMappingsAll":
		return "WaylandCommands"
	case "xcbCommandMappingsAll":
		return "XcbCommands"
	case "win32CommandMappingsAll":
		return "Win32Commands"
	case "appkitCommandMappingsAll":
		return "AppkitCommands"
	default:
		return "Commands"
	}
}

func waylandLoaderCommands() []loaderCommandSpec {
	return []loaderCommandSpec{
		{"DisplayConnect", "wl_display_connect"},
		{"ProxyMarshalConstructor", "wl_proxy_marshal_constructor"},
		{"ProxyMarshalConstructorVersioned", "wl_proxy_marshal_constructor_versioned"},
		{"ProxyAddListener", "wl_proxy_add_listener"},
		{"DisplayRoundtrip", "wl_display_roundtrip"},
		{"DisplayDispatchPending", "wl_display_dispatch_pending"},
		{"DisplayGetFd", "wl_display_get_fd"},
		{"DisplayPrepareRead", "wl_display_prepare_read"},
		{"DisplayReadEvents", "wl_display_read_events"},
		{"DisplayCancelRead", "wl_display_cancel_read"},
		{"ProxyMarshal", "wl_proxy_marshal"},
		{"ProxyDestroy", "wl_proxy_destroy"},
		{"DisplayDisconnect", "wl_display_disconnect"},
	}
}

func xcbLoaderCommands() []loaderCommandSpec {
	return []loaderCommandSpec{
		{"Connect", "xcb_connect"},
		{"GetSetup", "xcb_get_setup"},
		{"SetupRootsIterator", "xcb_setup_roots_iterator"},
		{"GenerateId", "xcb_generate_id"},
		{"CreateWindow", "xcb_create_window"},
		{"MapWindow", "xcb_map_window"},
		{"Flush", "xcb_flush"},
		{"InternAtom", "xcb_intern_atom"},
		{"InternAtomReply", "xcb_intern_atom_reply"},
		{"ChangeProperty", "xcb_change_property"},
		{"Free", "free"},
		{"DestroyWindow", "xcb_destroy_window"},
		{"Disconnect", "xcb_disconnect"},
		{"PollForEvent", "xcb_poll_for_event"},
		{"GetFileDescriptor", "xcb_get_file_descriptor"},
		{"WaitForEvent", "xcb_wait_for_event"},
	}
}

func win32LoaderCommandsGenerate(spec loaderGenSpec) error {
	elements := loaderGenHeaderElements(spec.BuildTag)
	elements = append(elements, gocode.FileElementFrom(gocode.DeclPackage(spec.PackageName)))
	bindingBlankLine(&elements)
	imports := spec.Imports
	if len(imports) > 0 {
		elements = append(elements, gocode.FileElementFrom(gocode.DeclImportBlock(imports...)))
	}
	bindingBlankLine(&elements)

	user32Func, err := loaderCommandFuncElements("win32User32CommandMappings", "Win32Commands", []loaderCommandSpec{
		{"RegisterClassExW", "RegisterClassExW"},
		{"CreateWindowExW", "CreateWindowExW"},
		{"DefWindowProcW", "DefWindowProcW"},
		{"ShowWindow", "ShowWindow"},
		{"UpdateWindow", "UpdateWindow"},
		{"LoadCursorW", "LoadCursorW"},
		{"DestroyWindow", "DestroyWindow"},
		{"PeekMessageW", "PeekMessageW"},
		{"DispatchMessageW", "DispatchMessageW"},
	})
	if err != nil {
		return err
	}
	elements = append(elements, user32Func...)

	kernel32Func, err := loaderCommandFuncElements("win32Kernel32CommandMappings", "Win32Commands", []loaderCommandSpec{
		{"GetModuleHandleW", "GetModuleHandleW"},
	})
	if err != nil {
		return err
	}
	elements = append(elements, kernel32Func...)

	file := gocode.DeclFile(elements...)
	content, err := gocode.GoFileRenderWithOptions(file, gocode.RenderOptionsEnumBindings())
	if err != nil {
		return fmt.Errorf("render loader: %w", err)
	}
	absPath := filepath.Join(spec.OutputDir, loaderCommandsGenFile)
	if err := os.WriteFile(absPath, []byte(content), 0o644); err != nil {
		return fmt.Errorf("write loader: %w", err)
	}
	return bindingFilesFormat(spec.OutputDir, loaderCommandsGenFile)
}

func loaderCommandFuncElements(funcName string, commandsStruct string, commands []loaderCommandSpec) ([]codegen.FileElement, error) {
	sorted := make([]loaderCommandSpec, len(commands))
	copy(sorted, commands)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].FieldName < sorted[j].FieldName })

	inits := make([]codegen.Expr, 0, len(sorted))
	for _, command := range sorted {
		inits = append(inits, gocode.ExprCompositeLit("loadutil.CommandMapping", []gocode.FieldInit{
			{
				Name:  "Target",
				Value: gocode.ExprAddressOf(gocode.ExprSelector(gocode.ExprIdent("commands"), command.FieldName)),
			},
			{Name: "Name", Value: gocode.ExprStringLit(command.SymbolName)},
		}))
	}

	body := gocode.StmtBlock(
		gocode.StmtReturn(gocode.ExprSliceCompositeLit("loadutil.CommandMapping", inits)),
	)
	elements := make([]codegen.FileElement, 0, 2)
	elements = append(elements, gocode.FileElementFrom(
		gocode.DeclFunc(funcName, []gocode.ParamType{
			{Name: "commands", Type: gocode.TypeExprPointer(gocode.TypeExprNamed(commandsStruct))},
		}, []gocode.TypeExpr{gocode.TypeExprSlice(gocode.TypeExprNamed("loadutil.CommandMapping"))}, body),
	))
	elements = append(elements, gocode.GoBlankLine())
	return elements, nil
}

func appkitLoaderCommands() []loaderCommandSpec {
	return []loaderCommandSpec{
		{"ObjcGetClass", "objc_getClass"},
		{"SelRegisterName", "sel_registerName"},
		{"ObjcMsgSend", "objc_msgSend"},
		{"ObjcAllocateClassPair", "objc_allocateClassPair"},
		{"ObjcRegisterClassPair", "objc_registerClassPair"},
		{"ClassAddMethod", "class_addMethod"},
	}
}
