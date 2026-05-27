package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"syscore/window/internal/generate/winmdgen"
)

const (
	win32TypesSubdir    = "types"
	win32PFNSubdir      = "pfn"
	win32BindingsModule = "syscore/window/win32/bindings"
)

/*
win32SeedDLLs lists Win32 DLL modules whose function pointers seed the binding closure.

[Context]
Only methods exported from these DLLs are emitted as PFN typedefs, and only types reachable from
their parameter or return signatures (transitively) are emitted under types/. This keeps the
generated surface focused on the windowing subset of the Win32 metadata while remaining fully
driven by the spec.
*/
var win32SeedDLLs = map[string]bool{
	"user32.dll":   true,
	"kernel32.dll": true,
}

/*
win32IsSeedDLL reports whether a Win32 DLL module is in the windowing seed set.

[Parameters]
module - Lowercase module name as reported by Context.MethodModuleName.
*/
func win32IsSeedDLL(module string) bool {
	return win32SeedDLLs[strings.ToLower(module)]
}

/*
win32GoTypePrimitives lists Go primitive identifiers that win32RewriteGoType must pass through
unchanged. Anything else encountered in a type expression is looked up in the type registry.
*/
var win32GoTypePrimitives = map[string]bool{
	"bool": true, "byte": true, "rune": true,
	"int": true, "int8": true, "uint8": true, "int16": true, "uint16": true,
	"int32": true, "uint32": true, "int64": true, "uint64": true,
	"float32": true, "float64": true,
	"uintptr": true, "any": true, "string": true, "error": true,
}

type win32TypeExport struct {
	TypeIR     winmdgen.TypeIR
	ExportName string
}

type win32MethodExport struct {
	MethodIR    winmdgen.MethodIR
	ExportName  string
	PFNTypeName string
}

type win32ConstantExport struct {
	ConstantIR winmdgen.ConstantIR
	ExportName string
}

/*
win32TypeRegistry tracks the canonical (package, ExportName) for every emitted TypeDef Go name.

[Context]
A registry entry is recorded once per (namespace, goName) so cross-arch variants share the same
qualified reference. Re-registering with the same (namespace, goName) is idempotent.
*/
type win32TypeRegistry struct {
	byNamespace map[string]map[string]win32TypeBinding
	byGoName    map[string][]win32TypeBinding
}

type win32TypeBinding struct {
	Namespace  string
	Slug       string
	Package    string
	ExportName string
}

func newWin32TypeRegistry() *win32TypeRegistry {
	return &win32TypeRegistry{
		byNamespace: make(map[string]map[string]win32TypeBinding),
		byGoName:    make(map[string][]win32TypeBinding),
	}
}

/*
Register records a TypeDef binding in the registry. Re-registering the same (namespace, goName)
is a no-op (later TypeIRs from other architectures share the canonical export name).
*/
func (registry *win32TypeRegistry) Register(namespace, slug, goName string) string {
	if existing, ok := registry.byNamespace[namespace][goName]; ok {
		return existing.ExportName
	}
	if registry.byNamespace[namespace] == nil {
		registry.byNamespace[namespace] = map[string]win32TypeBinding{}
	}
	packageName := win32GoPackageName(slug)
	binding := win32TypeBinding{
		Namespace:  namespace,
		Slug:       slug,
		Package:    packageName,
		ExportName: goName,
	}
	registry.byNamespace[namespace][goName] = binding
	registry.byGoName[goName] = append(registry.byGoName[goName], binding)
	return binding.ExportName
}

/*
Resolve returns the qualified form of identifier together with the import path that must be added
when the caller emits a file in preferNamespace.

[Returns]
qualified  - Qualified identifier (for example "foundation.HWND") or the bare identifier when it

	belongs to preferNamespace.

importPath - Import path needed in the current file, or empty when no import is required.
ok         - True when the identifier is a primitive or a registered TypeDef. False means the

	identifier is unknown; callers should leave it as-is so the missing dependency
	surfaces as a Go compile error rather than being silently masked.
*/
func (registry *win32TypeRegistry) Resolve(identifier string, preferNamespace string) (qualified string, importPath string, ok bool) {
	if win32GoTypePrimitives[identifier] {
		return identifier, "", true
	}
	if preferNamespace != "" {
		if binding, ok := registry.byNamespace[preferNamespace][identifier]; ok {
			return binding.ExportName, "", true
		}
	}
	bindings, exists := registry.byGoName[identifier]
	if !exists || len(bindings) == 0 {
		return "", "", false
	}
	binding := bindings[0]
	return binding.Package + "." + binding.ExportName, win32TypesImportPath(binding.Slug), true
}

func win32BindingsClean(outputDir string) error {
	entries, err := os.ReadDir(outputDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasPrefix(name, "bindings_") && strings.HasSuffix(name, "_gen.go") {
			if err := os.Remove(filepath.Join(outputDir, name)); err != nil {
				return fmt.Errorf("remove %s: %w", name, err)
			}
		}
	}
	if err := os.RemoveAll(filepath.Join(outputDir, win32TypesSubdir)); err != nil {
		return err
	}
	return os.RemoveAll(filepath.Join(outputDir, win32PFNSubdir))
}

func win32GoPackageName(slug string) string {
	name := strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_' {
			return r
		}
		return '_'
	}, slug)
	name = strings.Trim(name, "_")
	if name == "" {
		return "unknown"
	}
	if name[0] >= '0' && name[0] <= '9' {
		return "pkg_" + name
	}
	return name
}

func win32TypesImportPath(slug string) string {
	return win32BindingsModule + "/" + win32TypesSubdir + "/" + slug
}

func win32PFNImportPath(slug string) string {
	return win32BindingsModule + "/" + win32PFNSubdir + "/" + slug
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
	var builder strings.Builder
	for _, r := range strings.ToLower(value) {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			builder.WriteRune(r)
		case r == '.', r == '-', r == '_':
			builder.WriteRune('_')
		default:
			if r >= 128 {
				continue
			}
		}
	}
	slug := strings.Trim(builder.String(), "_")
	if len(slug) > 80 {
		slug = slug[:80]
	}
	if slug == "" {
		return "unknown"
	}
	return slug
}

/*
win32EnumMemberExportName picks a unique Go identifier for an enum member, suffixing the enum
name when the member name collides with a member already chosen for a different enum.
*/
func win32EnumMemberExportName(enumName string, memberName string, seen map[string]string) string {
	if ownerEnum, ok := seen[memberName]; ok && ownerEnum != enumName {
		return enumName + "_" + memberName
	}
	seen[memberName] = enumName
	return memberName
}

var win32GoTypeIdentifierPattern = regexp.MustCompile(`[A-Za-z_][A-Za-z0-9_]*`)

/*
win32RewriteGoType rewrites the type identifiers inside goType to refer to their canonical Go
exports. Identifiers from preferNamespace stay unqualified; identifiers from other namespaces are
qualified with their package selector and the import path is collected via imports.

[Context]
Unknown identifiers (those not in the registry) are returned as-is. This is intentional: the
binding generator only emits TypeDefs that participate in the windowing closure, so any reference
to a non-emitted type signals an incomplete closure and should surface as a Go compile error.
*/
func win32RewriteGoType(goType string, registry *win32TypeRegistry, preferNamespace string, imports map[string]string) string {
	return win32GoTypeIdentifierPattern.ReplaceAllStringFunc(goType, func(identifier string) string {
		if win32GoTypePrimitives[identifier] {
			return identifier
		}
		qualified, importPath, ok := registry.Resolve(identifier, preferNamespace)
		if !ok {
			return identifier
		}
		if importPath != "" {
			imports[importPath] = ""
		}
		return qualified
	})
}

/*
win32TypeExportsAssign deduplicates TypeIR variants that share a Go name within a single bucket.

[Context]
A bucket is one (namespace, arch) emission slot. Identical signatures collapse to a single export;
distinct signatures receive numeric suffixes so they can coexist in the same generated file.
*/
func win32TypeExportsAssign(typeIRs []winmdgen.TypeIR) []win32TypeExport {
	byName := make(map[string][]winmdgen.TypeIR)
	for _, typeIR := range typeIRs {
		byName[typeIR.GoName] = append(byName[typeIR.GoName], typeIR)
	}

	names := make([]string, 0, len(byName))
	for name := range byName {
		names = append(names, name)
	}
	sort.Strings(names)

	exports := make([]win32TypeExport, 0, len(typeIRs))
	for _, goName := range names {
		variants := byName[goName]
		seenSignatures := map[string]bool{}
		variantIndex := 0
		for _, typeIR := range variants {
			signature := win32TypeIRSignatureKey(typeIR)
			if seenSignatures[signature] {
				continue
			}
			seenSignatures[signature] = true

			exportName := goName
			if variantIndex > 0 {
				exportName = fmt.Sprintf("%s_%d", goName, variantIndex)
			}
			variantIndex++

			exports = append(exports, win32TypeExport{
				TypeIR:     typeIR,
				ExportName: exportName,
			})
		}
	}
	return exports
}

/*
win32MethodExportsAssign deduplicates MethodIR variants for a single bucket.
*/
func win32MethodExportsAssign(methodIRs []winmdgen.MethodIR) []win32MethodExport {
	byName := make(map[string][]winmdgen.MethodIR)
	for _, methodIR := range methodIRs {
		byName[methodIR.GoName] = append(byName[methodIR.GoName], methodIR)
	}

	names := make([]string, 0, len(byName))
	for name := range byName {
		names = append(names, name)
	}
	sort.Strings(names)

	exports := make([]win32MethodExport, 0, len(methodIRs))
	for _, goName := range names {
		variants := byName[goName]
		seenSignatures := map[string]bool{}
		variantIndex := 0
		for _, methodIR := range variants {
			signature := win32MethodIRSignatureKey(methodIR)
			if seenSignatures[signature] {
				continue
			}
			seenSignatures[signature] = true

			exportName := goName
			if variantIndex > 0 {
				exportName = fmt.Sprintf("%s_%d", goName, variantIndex)
			}
			variantIndex++

			exports = append(exports, win32MethodExport{
				MethodIR:    methodIR,
				ExportName:  exportName,
				PFNTypeName: "PFN_" + exportName,
			})
		}
	}
	return exports
}

/*
win32ConstantExportsAssign deduplicates Apis-class constants within a single namespace.
*/
func win32ConstantExportsAssign(constants []winmdgen.ConstantIR) []win32ConstantExport {
	byName := make(map[string][]winmdgen.ConstantIR)
	for _, constant := range constants {
		byName[constant.Name] = append(byName[constant.Name], constant)
	}

	names := make([]string, 0, len(byName))
	for name := range byName {
		names = append(names, name)
	}
	sort.Strings(names)

	exports := make([]win32ConstantExport, 0, len(constants))
	for _, goName := range names {
		variants := byName[goName]
		seenSignatures := map[string]bool{}
		variantIndex := 0
		for _, constant := range variants {
			signature := constant.GoType + "=" + constant.Value
			if seenSignatures[signature] {
				continue
			}
			seenSignatures[signature] = true

			exportName := goName
			if variantIndex > 0 {
				exportName = fmt.Sprintf("%s_%d", goName, variantIndex)
			}
			variantIndex++

			exports = append(exports, win32ConstantExport{
				ConstantIR: constant,
				ExportName: exportName,
			})
		}
	}
	return exports
}

/*
win32FoundationGuidTypeIR returns the synthetic Guid TypeIR injected into the Foundation package.

[Context]
The win32 metadata models GUID via System.Guid, which is not part of the winmd TypeDef table.
Many windowing-related structs reference Guid (for example through extension subsystems), so the
generator injects this shim so cross-package references resolve.
*/
func win32FoundationGuidTypeIR() winmdgen.TypeIR {
	return winmdgen.TypeIR{
		Namespace: "Windows.Win32.Foundation",
		GoName:    "Guid",
		Kind:      winmdgen.TypeKindStruct,
		Arch:      winmdgen.ArchAll,
		StructFields: []winmdgen.StructFieldIR{
			{Name: "Data1", GoType: "uint32"},
			{Name: "Data2", GoType: "uint16"},
			{Name: "Data3", GoType: "uint16"},
			{Name: "Data4", GoType: "[8]uint8"},
		},
	}
}

func win32FoundationNamespace() string {
	return "Windows.Win32.Foundation"
}

func win32FoundationSlug() string {
	return win32NamespaceSlug(win32FoundationNamespace())
}
