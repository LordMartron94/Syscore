package main

import (
	"fmt"

	"codegen"
	gocode "codegen/go"
)

type bindingPFNSpec struct {
	Name    string
	Func    gocode.TypeExpr
	Doc     string
}

func bindingPFNElementsFromSpecs(specs []bindingPFNSpec) ([]codegen.FileElement, error) {
	elements := make([]codegen.FileElement, 0, len(specs)*2)
	for _, spec := range specs {
		elements = append(elements, gocode.FileElementFrom(
			gocode.DeclTypeDefined(spec.Name, spec.Func, false, spec.Doc),
		))
		bindingBlankLine(&elements)
	}
	return elements, nil
}

func appkitPFNSpecs() []bindingPFNSpec {
	return []bindingPFNSpec{
		{
			Name: "PFN_objc_getClass",
			Func: gocode.TypeExprFunc(
				[]gocode.ParamType{{Name: "name", Type: gocode.TypeExprNamed("uintptr")}},
				[]gocode.TypeExpr{gocode.TypeExprNamed("ObjCClass")},
			),
			Doc: "PFN_objc_getClass maps to objc_getClass.",
		},
		{
			Name: "PFN_sel_registerName",
			Func: gocode.TypeExprFunc(
				[]gocode.ParamType{{Name: "name", Type: gocode.TypeExprNamed("uintptr")}},
				[]gocode.TypeExpr{gocode.TypeExprNamed("ObjCSel")},
			),
			Doc: "PFN_sel_registerName maps to sel_registerName.",
		},
		{
			Name: "PFN_objc_msgSend",
			Func: gocode.TypeExprFunc(
				[]gocode.ParamType{
					{Name: "receiver", Type: gocode.TypeExprNamed("ObjCObject")},
					{Name: "selector", Type: gocode.TypeExprNamed("ObjCSel")},
					{Name: "args", Type: gocode.TypeExprSlice(gocode.TypeExprNamed("uintptr"))},
				},
				[]gocode.TypeExpr{gocode.TypeExprNamed("ObjCObject")},
			),
			Doc: "PFN_objc_msgSend maps to objc_msgSend (pass trailing args as uintptr slice via purego).",
		},
		{
			Name: "PFN_objc_allocateClassPair",
			Func: gocode.TypeExprFunc(
				[]gocode.ParamType{
					{Name: "superClass", Type: gocode.TypeExprNamed("ObjCClass")},
					{Name: "name", Type: gocode.TypeExprPointer(gocode.TypeExprNamed("byte"))},
					{Name: "extraBytes", Type: gocode.TypeExprNamed("uintptr")},
				},
				[]gocode.TypeExpr{gocode.TypeExprNamed("ObjCClass")},
			),
			Doc: "PFN_objc_allocateClassPair maps to objc_allocateClassPair.",
		},
		{
			Name: "PFN_objc_registerClassPair",
			Func: gocode.TypeExprFunc(
				[]gocode.ParamType{{Name: "cls", Type: gocode.TypeExprNamed("ObjCClass")}},
				nil,
			),
			Doc: "PFN_objc_registerClassPair maps to objc_registerClassPair.",
		},
		{
			Name: "PFN_class_addMethod",
			Func: gocode.TypeExprFunc(
				[]gocode.ParamType{
					{Name: "cls", Type: gocode.TypeExprNamed("ObjCClass")},
					{Name: "name", Type: gocode.TypeExprNamed("ObjCSel")},
					{Name: "imp", Type: gocode.TypeExprNamed("uintptr")},
					{Name: "types", Type: gocode.TypeExprPointer(gocode.TypeExprNamed("byte"))},
				},
				[]gocode.TypeExpr{gocode.TypeExprNamed("uint8")},
			),
			Doc: "PFN_class_addMethod maps to class_addMethod.",
		},
	}
}

func bindingPFNElementsAppend(elements []codegen.FileElement, specs []bindingPFNSpec) ([]codegen.FileElement, error) {
	pfnElements, err := bindingPFNElementsFromSpecs(specs)
	if err != nil {
		return nil, fmt.Errorf("pfn types: %w", err)
	}
	return append(elements, pfnElements...), nil
}
