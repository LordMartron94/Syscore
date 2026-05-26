package main

import (
	"fmt"
	"sort"

	gocode "codegen/go"
)

const (
	waylandConstantsGenFile = "bindings_constants_gen.go"
	waylandTypesGenFile     = "bindings_types_gen.go"
)

func waylandBindingsGenerate(waylandXML string, clientCoreHeader string, outputDir string, interfaceNames []string, listenerEvents map[string][]string, extraConstants []gocode.ConstSpec, emitClientPFNs bool) error {
	interfaces, err := waylandProtocolIRBuild(waylandXML, interfaceNames)
	if err != nil {
		return err
	}

	if err := waylandConstantsGenerate(interfaces, outputDir, extraConstants); err != nil {
		return err
	}
	if err := waylandTypesGenerate(interfaces, outputDir, listenerEvents, clientCoreHeader, emitClientPFNs); err != nil {
		return err
	}
	return bindingFilesFormat(outputDir, waylandConstantsGenFile, waylandTypesGenFile)
}

func waylandConstantsGenerate(interfaces map[string]waylandInterfaceIR, outputDir string, extraConstants []gocode.ConstSpec) error {
	names := make([]string, 0, len(interfaces))
	for name := range interfaces {
		names = append(names, name)
	}
	sort.Strings(names)

	elements := bindingFilePreamble("bindings", "linux")
	constSpecs := make([]gocode.ConstSpec, 0, 32)
	constSpecs = append(constSpecs, extraConstants...)

	for _, name := range names {
		iface := interfaces[name]
		symbolConst := waylandGoSymbolConst(name)
		constSpecs = append(constSpecs, gocode.ConstSpecNew(
			symbolConst,
			nil,
			fmt.Sprintf("%q", name+"_interface"),
			fmt.Sprintf("%s is the exported wl_interface symbol for %s.", symbolConst, name),
		))
		for _, request := range iface.Requests {
			key := waylandGoOpcodeConst(name, "request", request.Name)
			constSpecs = append(constSpecs, gocode.ConstSpecNew(
				key,
				gocode.TypeExprNamedPtr("uint32"),
				fmt.Sprintf("%d", request.Index),
				fmt.Sprintf("%s is %s.%s request opcode.", key, name, request.Name),
			))
		}
		for _, event := range iface.Events {
			key := waylandGoOpcodeConst(name, "event", event.Name)
			constSpecs = append(constSpecs, gocode.ConstSpecNew(
				key,
				gocode.TypeExprNamedPtr("uint32"),
				fmt.Sprintf("%d", event.Index),
				fmt.Sprintf("%s is %s.%s event opcode.", key, name, event.Name),
			))
		}
	}

	elements = append(elements, gocode.FileElementFrom(gocode.DeclConstGroup(constSpecs, "", true)))
	bindingBlankLine(&elements)
	return writeBindingFile(outputDir, waylandConstantsGenFile, elements)
}

func waylandTypesGenerate(interfaces map[string]waylandInterfaceIR, outputDir string, listenerEvents map[string][]string, clientCoreHeader string, emitClientPFNs bool) error {
	names := make([]string, 0, len(interfaces))
	for name := range interfaces {
		names = append(names, name)
	}
	sort.Strings(names)

	elements := bindingFilePreamble("bindings", "linux")
	elements = append(elements, gocode.FileElementFrom(
		gocode.DeclTypeDefined(
			"WlProxy",
			gocode.TypeExprNamed("uintptr"),
			false,
			"WlProxy is an opaque struct wl_proxy* (libwayland-client core).",
		),
	))
	bindingBlankLine(&elements)

	for _, name := range names {
		iface := interfaces[name]
		typeName := waylandGoOpaqueType(name)
		elements = append(elements, gocode.FileElementFrom(
			gocode.DeclTypeDefined(
				typeName,
				gocode.TypeExprNamed("uintptr"),
				false,
				fmt.Sprintf("%s is an opaque struct %s*.", typeName, name),
			),
		))
		bindingBlankLine(&elements)

		events := waylandListenerEventsSelect(iface.Events, listenerEvents[name])
		if len(events) == 0 {
			continue
		}
		listenerName := waylandGoListenerType(name)
		fields := make([]gocode.StructFieldDecl, 0, len(events))
		for _, event := range events {
			fieldName := waylandGoListenerFieldName(event.Name)
			fields = append(fields, gocode.StructFieldTypeDoc(
				fieldName,
				gocode.TypeExprNamed("uintptr"),
				fmt.Sprintf("%s callback for %s.%s (purego.NewCallback); unset must be 0.", fieldName, name, event.Name),
			))
		}
		elements = append(elements, gocode.FileElementFrom(
			gocode.DeclTypeStruct(
				listenerName,
				fields,
				fmt.Sprintf("%s is the %s_listener vtable for wl_proxy_add_listener.", listenerName, name),
			),
		))
		bindingBlankLine(&elements)
	}

	if emitClientPFNs {
		pfnSpecs, err := waylandPFNSpecsFromVendorHeader(clientCoreHeader, waylandLoaderCommands())
		if err != nil {
			return fmt.Errorf("wayland pfn: %w", err)
		}
		elements, err = bindingPFNElementsAppend(elements, pfnSpecs)
		if err != nil {
			return err
		}
	}

	return writeBindingFile(outputDir, waylandTypesGenFile, elements)
}

func waylandListenerEventsSelect(events []waylandMessageIR, whitelist []string) []waylandMessageIR {
	if len(whitelist) == 0 {
		return nil
	}
	allowed := make(map[string]struct{}, len(whitelist))
	for _, name := range whitelist {
		allowed[name] = struct{}{}
	}
	selected := make([]waylandMessageIR, 0, len(whitelist))
	for _, event := range events {
		if _, ok := allowed[event.Name]; ok {
			selected = append(selected, event)
		}
	}
	return selected
}
