package main

import (
	"codegen"
	gocode "codegen/go"
)

const (
	appkitConstantsGenFile = "bindings_constants_gen.go"
	appkitTypesGenFile     = "bindings_types_gen.go"
)

func appkitBindingsGenerate(enumsPath string, outputDir string) error {
	ir, err := maciosBindingsIRBuild(enumsPath)
	if err != nil {
		return err
	}
	if err := appkitConstantsGenerate(outputDir, ir); err != nil {
		return err
	}
	if err := appkitTypesGenerate(outputDir); err != nil {
		return err
	}
	return bindingFilesFormat(outputDir, appkitConstantsGenFile, appkitTypesGenFile)
}

func appkitConstantsGenerate(outputDir string, ir maciosBindingsIR) error {
	elements := bindingFileStreamPreamble("bindings", "darwin")
	specs := make([]gocode.ConstSpec, 0, 64)

	for _, enumIR := range ir.Enums {
		enumSpecs := make([]gocode.ConstSpec, 0, len(enumIR.Members))
		for _, member := range enumIR.Members {
			enumSpecs = append(enumSpecs, gocode.ConstSpecNew(
				member.GoName,
				gocode.TypeExprNamedPtr("uint64"),
				member.Value,
				member.Doc,
			))
		}
		elements = append(elements, gocode.FileElementFrom(gocode.DeclConstGroup(enumSpecs, "", true)))
		bindingBlankLine(&elements)
	}
	for _, item := range ir.Classes {
		specs = append(specs, bindingConstSpecString(item.GoName, item.Value, item.Doc))
	}
	for _, item := range ir.Selectors {
		specs = append(specs, bindingConstSpecString(item.GoName, item.Value, item.Doc))
	}
	for _, item := range ir.Strings {
		specs = append(specs, bindingConstSpecString(item.GoName, item.Value, item.Doc))
	}
	for _, item := range ir.Extras {
		specs = append(specs, gocode.ConstSpecNew(
			item.GoName,
			gocode.TypeExprNamedPtr("uint64"),
			item.Value,
			item.Doc,
		))
	}
	if len(specs) > 0 {
		elements = append(elements, gocode.FileElementFrom(gocode.DeclConstGroup(specs, "", true)))
		bindingBlankLine(&elements)
	}
	return writeBindingFile(outputDir, appkitConstantsGenFile, elements)
}

func appkitTypesGenerate(outputDir string) error {
	elements := bindingFileStreamPreamble("bindings", "darwin")
	elements = bindingElementsWithImport(elements, "unsafe")
	elements = append(elements, gocode.FileElementFrom(
		gocode.DeclTypeDefined("ObjCObject", gocode.TypeExprNamed("uintptr"), false, "ObjCObject is an opaque Objective-C object handle."),
	))
	bindingBlankLine(&elements)
	elements = append(elements, gocode.FileElementFrom(
		gocode.DeclTypeDefined("ObjCClass", gocode.TypeExprNamed("uintptr"), false, "ObjCClass is an opaque Objective-C Class pointer."),
	))
	bindingBlankLine(&elements)
	elements = append(elements, gocode.FileElementFrom(
		gocode.DeclTypeDefined("ObjCSel", gocode.TypeExprNamed("uintptr"), false, "ObjCSel is an opaque SEL method selector pointer."),
	))
	bindingBlankLine(&elements)
	elements = append(elements, gocode.FileElementFrom(
		gocode.DeclTypeStruct("NSRect", []gocode.StructFieldDecl{
			gocode.StructFieldTypeDoc("X", gocode.TypeExprNamed("float64"), "X is the origin x coordinate."),
			gocode.StructFieldTypeDoc("Y", gocode.TypeExprNamed("float64"), "Y is the origin y coordinate."),
			gocode.StructFieldTypeDoc("Width", gocode.TypeExprNamed("float64"), "Width is the size width."),
			gocode.StructFieldTypeDoc("Height", gocode.TypeExprNamed("float64"), "Height is the size height."),
		}, "NSRect is the AppKit NSRect layout (flat CGFloat pairs) passed by value to objc_msgSend."),
	))
	bindingBlankLine(&elements)
	elements = append(elements, appkitSizeofGuardElements()...)

	var err error
	elements, err = bindingPFNElementsAppend(elements, appkitPFNSpecs())
	if err != nil {
		return err
	}
	return writeBindingFile(outputDir, appkitTypesGenFile, elements)
}

func appkitSizeofGuardElements() []codegen.FileElement {
	elements := make([]codegen.FileElement, 0, 2)
	elements = append(elements, gocode.FileElementFrom(
		gocode.DeclVarArray(
			"_appkitSizeofGuardNSRect",
			"32 - unsafe.Sizeof(NSRect{})",
			gocode.TypeExprNamed("byte"),
			"compile-time guard: NSRect must be 32 bytes (two CGFloat pairs).",
		),
	))
	bindingBlankLine(&elements)
	return elements
}
