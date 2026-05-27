package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/microsoft/go-winmd/winmd"

	"codegen"
	gocode "codegen/go"

	"syscore/window/internal/generate/winmdgen"
)

const (
	win32ApisGenFile       = "apis_gen.go"
	win32DefaultFlushEvery = 75
)

/*
win32ArchEmission describes how an arch bucket maps to a generated file name and build tag.

[Context]
ArchAll content lives in the canonical "_gen.go" file with a "windows" build tag. Architecture-
specific TypeDefs and PFNs are split into their own files so the resulting bindings package
mirrors the spec's architectural variants behind Go build tags.
*/
type win32ArchEmission struct {
	Arch     winmdgen.Arch
	FileName string
	BuildTag string
}

func win32TypeFiles(arch winmdgen.Arch) win32ArchEmission {
	switch arch {
	case winmdgen.Arch386:
		return win32ArchEmission{arch, "types_386_gen.go", "windows && 386"}
	case winmdgen.ArchAMD64:
		return win32ArchEmission{arch, "types_amd64_gen.go", "windows && amd64"}
	case winmdgen.ArchARM64:
		return win32ArchEmission{arch, "types_arm64_gen.go", "windows && arm64"}
	default:
		return win32ArchEmission{winmdgen.ArchAll, "types_gen.go", "windows"}
	}
}

func win32PFNFiles(arch winmdgen.Arch) win32ArchEmission {
	switch arch {
	case winmdgen.Arch386:
		return win32ArchEmission{arch, "pfn_386_gen.go", "windows && 386"}
	case winmdgen.ArchAMD64:
		return win32ArchEmission{arch, "pfn_amd64_gen.go", "windows && amd64"}
	case winmdgen.ArchARM64:
		return win32ArchEmission{arch, "pfn_arm64_gen.go", "windows && arm64"}
	default:
		return win32ArchEmission{winmdgen.ArchAll, "pfn_gen.go", "windows"}
	}
}

/*
win32AllSpecificArchs is the bitmask of the three concrete Windows architectures the spec models.
*/
const win32AllSpecificArchs = winmdgen.Arch386 | winmdgen.ArchAMD64 | winmdgen.ArchARM64

/*
win32ArchsForWalk returns the architectures to query when resolving signatures.

[Context]
ArchAll callers must expand into every concrete architecture so arch-specific dependencies (for
example SLIST_HEADER, which only exists as 386 / amd64 / arm64 TypeDefs) are pulled into the
closure rather than being lost behind ArchAll resolution.
*/
func win32ArchsForWalk(supportedArch winmdgen.Arch) []winmdgen.Arch {
	if supportedArch == winmdgen.ArchAll || supportedArch == 0 {
		return []winmdgen.Arch{winmdgen.Arch386, winmdgen.ArchAMD64, winmdgen.ArchARM64}
	}
	return supportedArch.Unique()
}

/*
win32ExpandToSpecificArchs widens an Arch bitmask so ArchAll becomes the explicit union of all
three concrete archs. Used when propagating "where is this needed" tags through the closure.
*/
func win32ExpandToSpecificArchs(arch winmdgen.Arch) winmdgen.Arch {
	if arch == winmdgen.ArchAll || arch == 0 {
		return win32AllSpecificArchs
	}
	return arch
}

/*
win32EmissionArchs maps the bitmask of architectures a TypeDef or PFN is needed for into the list
of emission slots. A type needed on every architecture collapses into a single ArchAll slot; types
needed only on a subset go into one slot per architecture (the same definition is duplicated when
necessary so build tags route it correctly).
*/
func win32EmissionArchs(neededArchs winmdgen.Arch) []winmdgen.Arch {
	expanded := win32ExpandToSpecificArchs(neededArchs)
	if expanded == win32AllSpecificArchs {
		return []winmdgen.Arch{winmdgen.ArchAll}
	}
	return expanded.Unique()
}

/*
win32IRContextArch chooses the architecture to use when building a TypeDef or MethodDef IR.

[Context]
A type may be needed on architectures that are NOT part of its declared SupportedArchitecture set
(this happens when an ArchAll caller references a TypeDef whose metadata pins it to a single
architecture). We must still build a valid IR; we do so using one of the type's declared
architectures so its field signatures resolve cleanly, and emit the same definition into the file
that the consumer needs.
*/
func win32IRContextArch(supportedArch winmdgen.Arch, emissionArch winmdgen.Arch) winmdgen.Arch {
	if emissionArch == winmdgen.ArchAll {
		return winmdgen.ArchAll
	}
	if supportedArch == winmdgen.ArchAll {
		return emissionArch
	}
	if supportedArch&emissionArch == emissionArch {
		return emissionArch
	}
	candidates := supportedArch.Unique()
	if len(candidates) == 0 {
		return winmdgen.ArchAll
	}
	return candidates[0]
}

/*
win32IsArchSupported reports whether the supported architecture bitmask allows resolution under the
requested architecture (ArchAll always allowed).
*/
func win32IsArchSupported(supportedArch winmdgen.Arch, requested winmdgen.Arch) bool {
	if requested == winmdgen.ArchAll {
		return supportedArch == winmdgen.ArchAll
	}
	return supportedArch == winmdgen.ArchAll || supportedArch&requested == requested
}

func win32BindingsGenerate(winmdPath string, outputDir string, flushEvery int) error {
	if flushEvery <= 0 {
		flushEvery = win32DefaultFlushEvery
	}

	if err := win32BindingsClean(outputDir); err != nil {
		return fmt.Errorf("clean bindings dir: %w", err)
	}

	metadata, err := winmd.Open(winmdPath)
	if err != nil {
		return fmt.Errorf("open winmd: %w", err)
	}
	context, err := winmdgen.NewContext(metadata)
	if err != nil {
		return fmt.Errorf("new winmd context: %w", err)
	}

	seedMethods, err := win32CollectSeedMethods(context, metadata)
	if err != nil {
		return fmt.Errorf("collect seed methods: %w", err)
	}

	typeClosure, err := win32ComputeTypeClosure(context, seedMethods)
	if err != nil {
		return fmt.Errorf("compute type closure: %w", err)
	}

	registry := newWin32TypeRegistry()
	registry.Register(win32FoundationNamespace(), win32FoundationSlug(), win32FoundationGuidTypeIR().GoName)

	typeBuckets, err := win32BuildTypeBuckets(context, typeClosure, registry)
	if err != nil {
		return fmt.Errorf("build type buckets: %w", err)
	}

	methodBuckets, err := win32BuildMethodBuckets(context, seedMethods)
	if err != nil {
		return fmt.Errorf("build method buckets: %w", err)
	}

	apisBuckets, err := win32BuildApisBuckets(metadata, typeClosure)
	if err != nil {
		return fmt.Errorf("build apis buckets: %w", err)
	}

	if err := win32EmitTypeBuckets(outputDir, typeBuckets, registry, flushEvery); err != nil {
		return err
	}
	if err := win32EmitApisBuckets(outputDir, apisBuckets, flushEvery); err != nil {
		return err
	}
	if err := win32EmitMethodBuckets(outputDir, methodBuckets, registry, flushEvery); err != nil {
		return err
	}

	if err := bindingDirFormat(filepath.Join(outputDir, win32TypesSubdir)); err != nil {
		return err
	}
	if err := bindingDirFormat(filepath.Join(outputDir, win32PFNSubdir)); err != nil {
		return err
	}
	return bindingDirFormat(outputDir)
}

/*
win32CollectSeedMethods returns every MethodDef index whose ImplMap module is a seed DLL.
*/
func win32CollectSeedMethods(context *winmdgen.Context, metadata *winmd.Metadata) ([]winmd.Index, error) {
	indices := make([]winmd.Index, 0, metadata.Tables.MethodDef.Len())
	for idx := range metadata.Tables.MethodDef.Indices() {
		module := context.MethodModuleName(idx)
		if module == "" || !win32IsSeedDLL(module) {
			continue
		}
		indices = append(indices, idx)
	}
	sort.Slice(indices, func(i, j int) bool { return indices[i] < indices[j] })
	return indices, nil
}

/*
win32ComputeTypeClosure performs a BFS over the type references reachable from seedMethods, tagging
each visited TypeDef with the bitmask of architectures it is needed on.

[Context]
This drives both bucketing and IR construction:

  - The returned bitmask determines which architecture file(s) a TypeDef gets emitted into. A
    TypeDef referenced from every arch collapses into the ArchAll bucket; one referenced only
    on a subset gets duplicated into each needed arch-specific file.
  - Walking a TypeDef uses its own SupportedArchitecture for resolution, but the discovered refs
    inherit the caller's needed-arch bitmask. This way an arm64 reference to an amd64-only
    TypeDef still propagates the arm64 need so the entire transitive closure remains available
    in arm64 build output.

Nested-public TypeDefs are inlined by the IR builder and never appear in the returned map.
*/
func win32ComputeTypeClosure(context *winmdgen.Context, seedMethods []winmd.Index) (map[winmd.Index]winmdgen.Arch, error) {
	closure := make(map[winmd.Index]winmdgen.Arch)
	type queueEntry struct {
		idx      winmd.Index
		newArchs winmdgen.Arch
	}
	queue := make([]queueEntry, 0, 64)

	enqueue := func(idx winmd.Index, archs winmdgen.Arch) {
		if archs == 0 {
			return
		}
		existing := closure[idx]
		newArchs := archs &^ existing
		if newArchs == 0 {
			return
		}
		nested, err := context.TypeDefIsNestedPublic(idx)
		if err == nil && nested {
			return
		}
		closure[idx] = existing | newArchs
		queue = append(queue, queueEntry{idx: idx, newArchs: newArchs})
	}

	for _, methodIndex := range seedMethods {
		supportedArch := context.MethodDefSupportedArch(methodIndex)
		neededArchs := win32ExpandToSpecificArchs(supportedArch)
		for _, walkArch := range win32ArchsForWalk(supportedArch) {
			refs, err := context.MethodReferencedTypeDefs(methodIndex, walkArch)
			if err != nil {
				return nil, err
			}
			for _, ref := range refs {
				enqueue(ref, neededArchs)
			}
		}
	}

	for len(queue) > 0 {
		entry := queue[0]
		queue = queue[1:]
		supportedArch := context.TypeDefSupportedArch(entry.idx)
		for _, walkArch := range win32ArchsForWalk(supportedArch) {
			refs, err := context.TypeReferencedTypeDefs(entry.idx, walkArch)
			if err != nil {
				return nil, err
			}
			for _, ref := range refs {
				enqueue(ref, entry.newArchs)
			}
		}
	}

	return closure, nil
}

/*
win32TypeBucketKey identifies a single emission slot in the types/<slug> package.
*/
type win32TypeBucketKey struct {
	Namespace string
	Slug      string
	Arch      winmdgen.Arch
}

type win32TypeBucket struct {
	Key   win32TypeBucketKey
	Types []winmdgen.TypeIR
}

/*
win32BuildTypeBuckets builds one TypeIR per (TypeDef, emission arch) tuple, where emission archs
come from the closure's per-TypeDef needed-arch bitmask. Each TypeIR is grouped by (namespace,
emission arch) and every unique (namespace, GoName) is registered in registry.
*/
func win32BuildTypeBuckets(context *winmdgen.Context, closure map[winmd.Index]winmdgen.Arch, registry *win32TypeRegistry) ([]win32TypeBucket, error) {
	buckets := map[win32TypeBucketKey]*win32TypeBucket{}

	addTypeIR := func(typeIR winmdgen.TypeIR, emissionArch winmdgen.Arch) {
		slug := win32NamespaceSlug(typeIR.Namespace)
		registry.Register(typeIR.Namespace, slug, typeIR.GoName)
		key := win32TypeBucketKey{Namespace: typeIR.Namespace, Slug: slug, Arch: emissionArch}
		bucket, ok := buckets[key]
		if !ok {
			bucket = &win32TypeBucket{Key: key}
			buckets[key] = bucket
		}
		bucket.Types = append(bucket.Types, typeIR)
	}

	indices := make([]winmd.Index, 0, len(closure))
	for idx := range closure {
		indices = append(indices, idx)
	}
	sort.Slice(indices, func(i, j int) bool { return indices[i] < indices[j] })

	for _, idx := range indices {
		supportedArch := context.TypeDefSupportedArch(idx)
		neededArchs := closure[idx]
		for _, emissionArch := range win32EmissionArchs(neededArchs) {
			irArch := win32IRContextArch(supportedArch, emissionArch)
			typeIR, err := context.BuildTypeIR(idx, irArch)
			if err != nil {
				fmt.Fprintf(os.Stderr, "skip typedef %d (%s): %v\n", idx, emissionArch, err)
				continue
			}
			if typeIR.GoName == "" || typeIR.GoName == "Apis" {
				continue
			}
			typeIR.Arch = emissionArch
			addTypeIR(typeIR, emissionArch)
		}
	}

	guidIR := win32FoundationGuidTypeIR()
	addTypeIR(guidIR, winmdgen.ArchAll)

	keys := make([]win32TypeBucketKey, 0, len(buckets))
	for key := range buckets {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		if keys[i].Slug != keys[j].Slug {
			return keys[i].Slug < keys[j].Slug
		}
		return keys[i].Arch < keys[j].Arch
	})

	ordered := make([]win32TypeBucket, 0, len(keys))
	for _, key := range keys {
		ordered = append(ordered, *buckets[key])
	}
	return ordered, nil
}

/*
win32MethodBucketKey identifies a single emission slot in the pfn/<slug> package.
*/
type win32MethodBucketKey struct {
	Module string
	Slug   string
	Arch   winmdgen.Arch
}

type win32MethodBucket struct {
	Key     win32MethodBucketKey
	Methods []winmdgen.MethodIR
}

func win32BuildMethodBuckets(context *winmdgen.Context, seedMethods []winmd.Index) ([]win32MethodBucket, error) {
	buckets := map[win32MethodBucketKey]*win32MethodBucket{}

	for _, idx := range seedMethods {
		module := context.MethodModuleName(idx)
		slug := win32ModuleSlug(module)
		supportedArch := context.MethodDefSupportedArch(idx)
		for _, emissionArch := range win32EmissionArchs(supportedArch) {
			irArch := win32IRContextArch(supportedArch, emissionArch)
			methodIR, err := context.BuildMethodIR(idx, irArch)
			if err != nil {
				fmt.Fprintf(os.Stderr, "skip method %d (%s): %v\n", idx, emissionArch, err)
				continue
			}
			methodIR.Arch = emissionArch
			key := win32MethodBucketKey{Module: module, Slug: slug, Arch: emissionArch}
			bucket, ok := buckets[key]
			if !ok {
				bucket = &win32MethodBucket{Key: key}
				buckets[key] = bucket
			}
			bucket.Methods = append(bucket.Methods, methodIR)
		}
	}

	keys := make([]win32MethodBucketKey, 0, len(buckets))
	for key := range buckets {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		if keys[i].Slug != keys[j].Slug {
			return keys[i].Slug < keys[j].Slug
		}
		return keys[i].Arch < keys[j].Arch
	})

	ordered := make([]win32MethodBucket, 0, len(keys))
	for _, key := range keys {
		ordered = append(ordered, *buckets[key])
	}
	return ordered, nil
}

/*
win32ApisBucket holds the integer Apis-class constants emitted for one source namespace.
*/
type win32ApisBucket struct {
	Namespace string
	Slug      string
	Constants []winmdgen.ConstantIR
}

/*
win32BuildApisBuckets returns Apis-class constants grouped by namespace, restricted to namespaces
that appear in the TypeDef closure (so we never emit a standalone constants file for a namespace
that has no types to host it).
*/
func win32BuildApisBuckets(metadata *winmd.Metadata, closure map[winmd.Index]winmdgen.Arch) ([]win32ApisBucket, error) {
	namespacesInClosure := map[string]bool{}
	for idx := range closure {
		def, err := metadata.Tables.TypeDef.At(idx)
		if err != nil {
			continue
		}
		namespacesInClosure[def.Namespace.String()] = true
	}

	bucketsByNamespace := map[string]*win32ApisBucket{}
	err := winmdgen.StreamApisConstants(metadata, func(constant winmdgen.ConstantIR) bool {
		if !namespacesInClosure[constant.Namespace] {
			return true
		}
		bucket, ok := bucketsByNamespace[constant.Namespace]
		if !ok {
			bucket = &win32ApisBucket{
				Namespace: constant.Namespace,
				Slug:      win32NamespaceSlug(constant.Namespace),
			}
			bucketsByNamespace[constant.Namespace] = bucket
		}
		bucket.Constants = append(bucket.Constants, constant)
		return true
	})
	if err != nil {
		return nil, err
	}

	keys := make([]string, 0, len(bucketsByNamespace))
	for namespace := range bucketsByNamespace {
		keys = append(keys, namespace)
	}
	sort.Strings(keys)

	ordered := make([]win32ApisBucket, 0, len(keys))
	for _, namespace := range keys {
		ordered = append(ordered, *bucketsByNamespace[namespace])
	}
	return ordered, nil
}

/*
win32EmitTypeBuckets writes one file per (namespace, arch) bucket under bindings/types/<slug>.
*/
func win32EmitTypeBuckets(outputDir string, buckets []win32TypeBucket, registry *win32TypeRegistry, flushEvery int) error {
	for _, bucket := range buckets {
		emission := win32TypeFiles(bucket.Key.Arch)
		packageDir := filepath.Join(outputDir, win32TypesSubdir, bucket.Key.Slug)
		if err := os.MkdirAll(packageDir, 0o755); err != nil {
			return err
		}
		packageName := win32GoPackageName(bucket.Key.Slug)

		stream, err := bindingFileStreamOpen(packageDir, emission.FileName, packageName, emission.BuildTag, flushEvery)
		if err != nil {
			return fmt.Errorf("open type stream %s/%s: %w", bucket.Key.Slug, emission.FileName, err)
		}

		exports := win32TypeExportsAssign(bucket.Types)
		seenEnumMembers := map[string]string{}
		imports := map[string]string{}
		typeBatches := make([][]codegen.FileElement, 0, len(exports))
		needsUnsafe := false
		for _, export := range exports {
			elements, batchNeedsUnsafe, err := win32TypeElementsGenerate(export, bucket.Key.Namespace, seenEnumMembers, registry, imports)
			if err != nil {
				fmt.Fprintf(os.Stderr, "skip type %s in %s: %v\n", export.ExportName, bucket.Key.Slug, err)
				continue
			}
			if batchNeedsUnsafe {
				needsUnsafe = true
			}
			typeBatches = append(typeBatches, elements)
		}
		if needsUnsafe {
			imports["unsafe"] = ""
		}
		if len(imports) > 0 {
			if err := win32BindingStreamWriteImportBlock(stream, imports); err != nil {
				_ = bindingFileStreamClose(stream)
				return fmt.Errorf("write type imports %s: %w", bucket.Key.Slug, err)
			}
		}
		for _, batch := range typeBatches {
			if err := bindingFileStreamWriteElements(stream, batch...); err != nil {
				fmt.Fprintf(os.Stderr, "write type batch in %s: %v\n", bucket.Key.Slug, err)
			}
		}
		if err := bindingFileStreamClose(stream); err != nil {
			return err
		}
	}
	return nil
}

/*
win32EmitApisBuckets writes one apis_gen.go per namespace with integer Apis-class constants.
*/
func win32EmitApisBuckets(outputDir string, buckets []win32ApisBucket, flushEvery int) error {
	for _, bucket := range buckets {
		packageDir := filepath.Join(outputDir, win32TypesSubdir, bucket.Slug)
		if err := os.MkdirAll(packageDir, 0o755); err != nil {
			return err
		}
		packageName := win32GoPackageName(bucket.Slug)

		stream, err := bindingFileStreamOpen(packageDir, win32ApisGenFile, packageName, "windows", flushEvery)
		if err != nil {
			return fmt.Errorf("open apis stream %s: %w", bucket.Slug, err)
		}

		exports := win32ConstantExportsAssign(bucket.Constants)
		specs := make([]gocode.ConstSpec, 0, len(exports))
		for _, export := range exports {
			specs = append(specs, gocode.ConstSpecNew(
				export.ExportName,
				nil,
				export.ConstantIR.Value,
				export.ExportName+" from Windows.Win32 metadata (originally typed "+export.ConstantIR.GoType+"). Emitted untyped so callers convert to the receiving Go type at the use site.",
			))
		}
		if len(specs) == 0 {
			_ = bindingFileStreamClose(stream)
			continue
		}
		elements := []codegen.FileElement{
			gocode.FileElementFrom(gocode.DeclConstGroup(specs, "", true)),
		}
		bindingBlankLine(&elements)
		if err := bindingFileStreamWriteElements(stream, elements...); err != nil {
			_ = bindingFileStreamClose(stream)
			return fmt.Errorf("write apis %s: %w", bucket.Slug, err)
		}
		if err := bindingFileStreamClose(stream); err != nil {
			return err
		}
	}
	return nil
}

/*
win32EmitMethodBuckets writes one file per (DLL, arch) bucket under bindings/pfn/<slug>.
*/
func win32EmitMethodBuckets(outputDir string, buckets []win32MethodBucket, registry *win32TypeRegistry, flushEvery int) error {
	for _, bucket := range buckets {
		emission := win32PFNFiles(bucket.Key.Arch)
		packageDir := filepath.Join(outputDir, win32PFNSubdir, bucket.Key.Slug)
		if err := os.MkdirAll(packageDir, 0o755); err != nil {
			return err
		}
		packageName := win32GoPackageName(bucket.Key.Slug)

		stream, err := bindingFileStreamOpen(packageDir, emission.FileName, packageName, emission.BuildTag, flushEvery)
		if err != nil {
			return fmt.Errorf("open pfn stream %s/%s: %w", bucket.Key.Slug, emission.FileName, err)
		}

		exports := win32MethodExportsAssign(bucket.Methods)
		imports := map[string]string{}
		batches := make([][]codegen.FileElement, 0, len(exports))
		needsUnsafe := false
		for _, export := range exports {
			elements, batchNeedsUnsafe, err := win32PFNTypeElementsGenerate(export, registry, imports)
			if err != nil {
				fmt.Fprintf(os.Stderr, "skip pfn %s in %s: %v\n", export.PFNTypeName, bucket.Key.Slug, err)
				continue
			}
			if batchNeedsUnsafe {
				needsUnsafe = true
			}
			batches = append(batches, elements)
		}

		if needsUnsafe {
			imports["unsafe"] = ""
		}
		if len(imports) > 0 {
			if err := win32BindingStreamWriteImportBlock(stream, imports); err != nil {
				_ = bindingFileStreamClose(stream)
				return fmt.Errorf("write pfn imports %s: %w", bucket.Key.Slug, err)
			}
		}
		for _, batch := range batches {
			if err := bindingFileStreamWriteElements(stream, batch...); err != nil {
				fmt.Fprintf(os.Stderr, "write pfn batch in %s: %v\n", bucket.Key.Slug, err)
			}
		}
		if err := bindingFileStreamClose(stream); err != nil {
			return err
		}
	}
	return nil
}

func win32TypeElementsGenerate(
	export win32TypeExport,
	preferNamespace string,
	seenEnumMembers map[string]string,
	registry *win32TypeRegistry,
	imports map[string]string,
) ([]codegen.FileElement, bool, error) {
	typeIR := export.TypeIR
	exportedName := export.ExportName
	needsUnsafe := false
	elements := make([]codegen.FileElement, 0, 8)
	doc := exportedName + " is generated from Windows.Win32 metadata."

	switch typeIR.Kind {
	case winmdgen.TypeKindEnum:
		enumUnderlying := win32RewriteGoType(typeIR.UnderlyingType, registry, preferNamespace, imports)
		underlying, err := gocode.TypeExprFromGoTypeString(enumUnderlying)
		if err != nil {
			return nil, false, err
		}
		elements = append(elements, gocode.FileElementFrom(
			gocode.DeclTypeDefined(exportedName, underlying, false, doc),
		))
		bindingBlankLine(&elements)

		specs := make([]gocode.ConstSpec, 0, len(typeIR.EnumMembers))
		enumType := gocode.TypeExprNamedPtr(exportedName)
		for _, member := range typeIR.EnumMembers {
			memberName := win32EnumMemberExportName(exportedName, member.Name, seenEnumMembers)
			specs = append(specs, gocode.ConstSpecNew(
				memberName,
				enumType,
				member.Value,
				memberName+" from Windows.Win32 metadata.",
			))
		}
		elements = append(elements, gocode.FileElementFrom(gocode.DeclConstGroup(specs, "", true)))
		bindingBlankLine(&elements)

	case winmdgen.TypeKindStruct:
		fieldNameUses := map[string]int{}
		fields := make([]gocode.StructFieldDecl, 0, len(typeIR.StructFields))
		for _, field := range typeIR.StructFields {
			fieldType := win32RewriteGoType(field.GoType, registry, preferNamespace, imports)
			typ, err := gocode.TypeExprFromGoTypeString(fieldType)
			if err != nil {
				return nil, false, err
			}
			if strings.Contains(fieldType, "unsafe.Pointer") {
				needsUnsafe = true
			}
			name := field.Name
			if count := fieldNameUses[field.Name]; count > 0 {
				name = fmt.Sprintf("%s_%d", field.Name, count)
			}
			fieldNameUses[field.Name]++
			fields = append(fields, gocode.StructFieldTypeDoc(name, typ, ""))
		}
		elements = append(elements, gocode.FileElementFrom(
			gocode.DeclTypeStruct(exportedName, fields, doc),
		))
		bindingBlankLine(&elements)

	case winmdgen.TypeKindNative, winmdgen.TypeKindDelegate:
		underlyingType := win32RewriteGoType(typeIR.UnderlyingType, registry, preferNamespace, imports)
		underlying, err := gocode.TypeExprFromGoTypeString(underlyingType)
		if err != nil {
			return nil, false, err
		}
		if strings.Contains(underlyingType, "unsafe.Pointer") {
			needsUnsafe = true
		}
		elements = append(elements, gocode.FileElementFrom(
			gocode.DeclTypeDefined(exportedName, underlying, false, doc),
		))
		bindingBlankLine(&elements)
	}

	return elements, needsUnsafe, nil
}

func win32PFNTypeElementsGenerate(
	export win32MethodExport,
	registry *win32TypeRegistry,
	imports map[string]string,
) ([]codegen.FileElement, bool, error) {
	methodIR := export.MethodIR
	needsUnsafe := false
	params := make([]gocode.ParamType, 0, len(methodIR.Params))
	for _, param := range methodIR.Params {
		goType := win32RewriteGoType(param.GoType, registry, "", imports)
		typ, err := gocode.TypeExprFromGoTypeString(goType)
		if err != nil {
			return nil, false, err
		}
		if strings.Contains(goType, "unsafe.Pointer") {
			needsUnsafe = true
		}
		params = append(params, gocode.ParamType{Name: param.Name, Type: typ})
	}

	returns := []gocode.TypeExpr(nil)
	if methodIR.HasReturn {
		returnType := win32RewriteGoType(methodIR.ReturnType, registry, "", imports)
		ret, err := gocode.TypeExprFromGoTypeString(returnType)
		if err != nil {
			return nil, false, err
		}
		if strings.Contains(returnType, "unsafe.Pointer") {
			needsUnsafe = true
		}
		returns = []gocode.TypeExpr{ret}
	}

	funcType := gocode.TypeExprFunc(params, returns)
	elements := make([]codegen.FileElement, 0, 2)
	elements = append(elements, gocode.FileElementFrom(
		gocode.DeclTypeDefined(
			export.PFNTypeName,
			funcType,
			false,
			export.PFNTypeName+" maps to "+methodIR.GoName+" in "+methodIR.Module+".",
		),
	))
	bindingBlankLine(&elements)
	return elements, needsUnsafe, nil
}

func win32MethodIRSignatureKey(methodIR winmdgen.MethodIR) string {
	var signature strings.Builder
	signature.WriteString(fmt.Sprintf("arch=%d;", methodIR.Arch))
	for _, param := range methodIR.Params {
		signature.WriteString(param.GoType)
		signature.WriteByte(';')
	}
	if methodIR.HasReturn {
		signature.WriteString(methodIR.ReturnType)
	}
	return signature.String()
}

func win32TypeIRSignatureKey(typeIR winmdgen.TypeIR) string {
	var signature strings.Builder
	signature.WriteString(fmt.Sprintf("kind=%d;arch=%d;", typeIR.Kind, typeIR.Arch))
	signature.WriteString(typeIR.UnderlyingType)
	signature.WriteByte(';')
	for _, field := range typeIR.StructFields {
		signature.WriteString(field.Name)
		signature.WriteByte('=')
		signature.WriteString(field.GoType)
		signature.WriteByte(';')
	}
	for _, member := range typeIR.EnumMembers {
		signature.WriteString(member.Name)
		signature.WriteByte('=')
		signature.WriteString(member.Value)
		signature.WriteByte(';')
	}
	return signature.String()
}
