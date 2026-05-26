package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/microsoft/go-winmd/winmd"

	"codegen"
	gocode "codegen/go"

	"syscore/window/internal/generate/winmdgen"
)

const (
	win32TypesFilePrefix   = "bindings_types_"
	win32PFNFilePrefix     = "bindings_pfn_"
	win32GeneratedSuffix   = "_gen.go"
	win32DefaultFlushEvery = 75
)

var win32LegacyBindingFiles = []string{
	"bindings_types_gen.go",
	"bindings_pfn_gen.go",
}

func win32BindingsGenerate(winmdPath string, outputDir string, flushEvery int) error {
	if flushEvery <= 0 {
		flushEvery = win32DefaultFlushEvery
	}

	for _, name := range win32LegacyBindingFiles {
		legacyPath := filepath.Join(outputDir, name)
		if err := os.Remove(legacyPath); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("remove %s: %w", legacyPath, err)
		}
	}

	metadata, err := winmd.Open(winmdPath)
	if err != nil {
		return fmt.Errorf("open winmd: %w", err)
	}

	typeStreams := make(map[string]*bindingFileStream)
	pfnStreams := make(map[string]*bindingFileStream)

	if err := winmdgen.BuildBindingsIRStream(metadata, func(typeIR winmdgen.TypeIR) bool {
		slug := win32NamespaceSlug(typeIR.Namespace)
		stream, err := win32TypesStreamEnsure(typeStreams, outputDir, slug, flushEvery)
		if err != nil {
			fmt.Fprintf(os.Stderr, "types stream %s: %v\n", slug, err)
			return true
		}
		typeElements, _, err := win32TypeElementsGenerate(typeIR)
		if err != nil {
			fmt.Fprintf(os.Stderr, "skip type %s: %v\n", typeIR.GoName, err)
			return true
		}
		if err := bindingFileStreamWriteElements(stream, typeElements...); err != nil {
			fmt.Fprintf(os.Stderr, "write type %s: %v\n", typeIR.GoName, err)
		}
		return true
	}); err != nil {
		return fmt.Errorf("winmd types: %w", err)
	}

	if err := winmdgen.StreamMethodIRs(metadata, func(methodIR winmdgen.MethodIR) bool {
		slug := win32ModuleSlug(methodIR.Module)
		stream, err := win32PFNStreamEnsure(pfnStreams, outputDir, slug, flushEvery)
		if err != nil {
			fmt.Fprintf(os.Stderr, "pfn stream %s: %v\n", slug, err)
			return true
		}
		pfnElements, _, err := win32PFNTypeElementsGenerate(methodIR)
		if err != nil {
			fmt.Fprintf(os.Stderr, "skip pfn %s: %v\n", methodIR.GoName, err)
			return true
		}
		if err := bindingFileStreamWriteElements(stream, pfnElements...); err != nil {
			fmt.Fprintf(os.Stderr, "write pfn %s: %v\n", methodIR.GoName, err)
		}
		return true
	}); err != nil {
		return fmt.Errorf("winmd methods: %w", err)
	}

	for _, stream := range typeStreams {
		if err := bindingFileStreamClose(stream); err != nil {
			return err
		}
	}
	for _, stream := range pfnStreams {
		if err := bindingFileStreamClose(stream); err != nil {
			return err
		}
	}
	if len(typeStreams) == 0 && len(pfnStreams) == 0 {
		return nil
	}
	return bindingDirFormat(outputDir)
}

func win32TypesStreamEnsure(streams map[string]*bindingFileStream, outputDir, slug string, flushEvery int) (*bindingFileStream, error) {
	if stream, ok := streams[slug]; ok {
		return stream, nil
	}
	stream, err := win32BindingStreamOpen(outputDir, win32TypesFileName(slug), flushEvery)
	if err != nil {
		return nil, err
	}
	streams[slug] = stream
	return stream, nil
}

func win32PFNStreamEnsure(streams map[string]*bindingFileStream, outputDir, slug string, flushEvery int) (*bindingFileStream, error) {
	if stream, ok := streams[slug]; ok {
		return stream, nil
	}
	stream, err := win32BindingStreamOpen(outputDir, win32PFNFileName(slug), flushEvery)
	if err != nil {
		return nil, err
	}
	streams[slug] = stream
	return stream, nil
}

func win32BindingStreamOpen(outputDir, fileName string, flushEvery int) (*bindingFileStream, error) {
	elements := bindingFileStreamPreamble("bindings", "windows")
	return bindingFileStreamCreate(outputDir, fileName, elements, flushEvery)
}

func win32TypesFileName(slug string) string {
	return win32TypesFilePrefix + slug + win32GeneratedSuffix
}

func win32PFNFileName(slug string) string {
	return win32PFNFilePrefix + slug + win32GeneratedSuffix
}

func win32NamespaceSlug(namespace string) string {
	return win32SlugFromString(namespace)
}

func win32ModuleSlug(module string) string {
	slug := win32SlugFromString(module)
	if slug == "" {
		return "unknown"
	}
	return slug
}

func win32SlugFromString(value string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(value) {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
		case r == '.', r == '-', r == '_':
			b.WriteRune('_')
		default:
			if unicode.IsLetter(r) || unicode.IsDigit(r) {
				b.WriteRune(unicode.ToLower(r))
			}
		}
	}
	slug := strings.Trim(b.String(), "_")
	if len(slug) > 80 {
		slug = slug[:80]
	}
	if slug == "" {
		return "unknown"
	}
	return slug
}

func win32TypeElementsGenerate(typeIR winmdgen.TypeIR) ([]codegen.FileElement, bool, error) {
	needsUnsafe := false
	elements := make([]codegen.FileElement, 0, 8)
	doc := typeIR.GoName + " is generated from Windows.Win32 metadata."

	switch typeIR.Kind {
	case winmdgen.TypeKindEnum:
		underlying, err := gocode.TypeExprFromGoTypeString(typeIR.UnderlyingType)
		if err != nil {
			return nil, false, err
		}
		elements = append(elements, gocode.FileElementFrom(
			gocode.DeclTypeDefined(typeIR.GoName, underlying, false, doc),
		))
		bindingBlankLine(&elements)

		specs := make([]gocode.ConstSpec, 0, len(typeIR.EnumMembers))
		enumType := gocode.TypeExprNamedPtr(typeIR.GoName)
		for _, member := range typeIR.EnumMembers {
			specs = append(specs, gocode.ConstSpecNew(
				member.Name,
				enumType,
				member.Value,
				member.Name+" from Windows.Win32 metadata.",
			))
		}
		elements = append(elements, gocode.FileElementFrom(gocode.DeclConstGroup(specs, "", true)))
		bindingBlankLine(&elements)

	case winmdgen.TypeKindStruct:
		fields := make([]gocode.StructFieldDecl, 0, len(typeIR.StructFields))
		for _, field := range typeIR.StructFields {
			typ, err := gocode.TypeExprFromGoTypeString(field.GoType)
			if err != nil {
				return nil, false, err
			}
			if strings.Contains(field.GoType, "unsafe.Pointer") {
				needsUnsafe = true
			}
			fields = append(fields, gocode.StructFieldTypeDoc(field.Name, typ, ""))
		}
		elements = append(elements, gocode.FileElementFrom(
			gocode.DeclTypeStruct(typeIR.GoName, fields, doc),
		))
		bindingBlankLine(&elements)

	case winmdgen.TypeKindNative, winmdgen.TypeKindDelegate:
		underlying, err := gocode.TypeExprFromGoTypeString(typeIR.UnderlyingType)
		if err != nil {
			return nil, false, err
		}
		if strings.Contains(typeIR.UnderlyingType, "unsafe.Pointer") {
			needsUnsafe = true
		}
		elements = append(elements, gocode.FileElementFrom(
			gocode.DeclTypeDefined(typeIR.GoName, underlying, false, doc),
		))
		bindingBlankLine(&elements)
	}

	return elements, needsUnsafe, nil
}

func win32PFNTypeElementsGenerate(methodIR winmdgen.MethodIR) ([]codegen.FileElement, bool, error) {
	needsUnsafe := false
	params := make([]gocode.ParamType, 0, len(methodIR.Params))
	for _, param := range methodIR.Params {
		typ, err := gocode.TypeExprFromGoTypeString(param.GoType)
		if err != nil {
			return nil, false, err
		}
		if strings.Contains(param.GoType, "unsafe.Pointer") {
			needsUnsafe = true
		}
		params = append(params, gocode.ParamType{Name: param.Name, Type: typ})
	}

	returns := []gocode.TypeExpr(nil)
	if methodIR.HasReturn {
		ret, err := gocode.TypeExprFromGoTypeString(methodIR.ReturnType)
		if err != nil {
			return nil, false, err
		}
		if strings.Contains(methodIR.ReturnType, "unsafe.Pointer") {
			needsUnsafe = true
		}
		returns = []gocode.TypeExpr{ret}
	}

	funcType := gocode.TypeExprFunc(params, returns)
	elements := make([]codegen.FileElement, 0, 2)
	elements = append(elements, gocode.FileElementFrom(
		gocode.DeclTypeDefined(
			"PFN_"+methodIR.GoName,
			funcType,
			false,
			"PFN_"+methodIR.GoName+" maps to "+methodIR.GoName+" in "+methodIR.Module+".",
		),
	))
	bindingBlankLine(&elements)
	return elements, needsUnsafe, nil
}
