// Copyright (c) Microsoft Corporation. Licensed under the MIT License.
// Adapted from github.com/microsoft/go-winmd/cmd/gowinmd/internal/gowinmd for syscore window bindings.

package winmdgen

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"

	"github.com/microsoft/go-winmd/winmd"
)

// StringWriter adapts io.Writer to io.StringWriter for type emission.
type StringWriter struct {
	io.Writer
}

func (w StringWriter) WriteString(s string) (int, error) {
	return io.WriteString(w.Writer, s)
}

type TypeKind int

const (
	TypeKindEnum TypeKind = iota
	TypeKindStruct
	TypeKindNative
	TypeKindDelegate
)

type EnumMemberIR struct {
	Name  string
	Value string
}

type StructFieldIR struct {
	Name    string
	GoType  string
	ByValue bool
}

type TypeIR struct {
	Index          winmd.Index
	Arch           Arch
	Namespace      string
	GoName         string
	Kind           TypeKind
	UnderlyingType string
	EnumMembers    []EnumMemberIR
	StructFields   []StructFieldIR
}

type ParamIR struct {
	Name   string
	GoType string
}

type MethodIR struct {
	Index      winmd.Index
	Arch       Arch
	GoName     string
	Module     string
	Params     []ParamIR
	ReturnType string
	HasReturn  bool
}

type BindingsIR struct {
	Types   []TypeIR
	Methods []MethodIR
}

/*
ConstantIR represents a single constant exposed via the "Apis" class of a namespace.

[Context]
Win32 metadata models free-standing constants (for example CW_USEDEFAULT and WM_CLOSE) as static
constant fields on a per-namespace "Apis" TypeDef rather than as enum members. Emit these via
StreamApisConstants so the generator can produce real Go const declarations instead of struct fields.
*/
type ConstantIR struct {
	Namespace string
	Name      string
	GoType    string
	Value     string
}

func BuildBindingsIR(metadata *winmd.Metadata) (BindingsIR, error) {
	var ir BindingsIR
	if err := BuildBindingsIRStream(metadata, func(typeIR TypeIR) bool {
		ir.Types = append(ir.Types, typeIR)
		return true
	}); err != nil {
		return BindingsIR{}, err
	}
	if err := StreamMethodIRs(metadata, func(methodIR MethodIR) bool {
		ir.Methods = append(ir.Methods, methodIR)
		return true
	}); err != nil {
		return BindingsIR{}, err
	}
	return ir, nil
}

func BuildBindingsIRStream(metadata *winmd.Metadata, yield func(TypeIR) bool) error {
	context, err := NewContext(metadata)
	if err != nil {
		return err
	}

	indices := make([]winmd.Index, 0, metadata.Tables.TypeDef.Len())
	for idx := range metadata.Tables.TypeDef.Indices() {
		typeDef, err := metadata.Tables.TypeDef.At(idx)
		if err != nil {
			return err
		}
		if typeDef.Flags&winmd.TypeAttributes_NestedPublic != 0 {
			continue
		}
		indices = append(indices, idx)
	}
	sort.Slice(indices, func(i, j int) bool { return indices[i] < indices[j] })

	for _, idx := range indices {
		resolved, err := context.resolveTypeDef(idx)
		if err != nil {
			return fmt.Errorf("resolve typedef %d: %w", idx, err)
		}
		typeIR, err := context.typeIRFromResolved(resolved, ArchAll)
		if err != nil {
			continue
		}
		typeIR.Index = idx
		typeIR.Arch = context.TypeDefSupportedArch(idx)
		typeIR.Namespace = resolved.Namespace.String()
		if typeIR.GoName == "" || !win32GoNameValid(typeIR.GoName) {
			continue
		}
		if !yield(typeIR) {
			return nil
		}
	}
	return nil
}

/*
StreamApisConstants iterates constants exposed via per-namespace "Apis" classes and yields one
ConstantIR per integer-typed constant field. Non-integer constants (Guids, strings, struct values)
are skipped because they cannot be expressed as Go untyped const declarations.

[Context]
Used by the Win32 binding generator to emit Apis-class constants as Go `const` declarations
inside the type subpackage for each namespace they belong to.
*/
func StreamApisConstants(metadata *winmd.Metadata, yield func(ConstantIR) bool) error {
	context, err := NewContext(metadata)
	if err != nil {
		return err
	}

	indices := make([]winmd.Index, 0, metadata.Tables.TypeDef.Len())
	for idx := range metadata.Tables.TypeDef.Indices() {
		typeDef, err := metadata.Tables.TypeDef.At(idx)
		if err != nil {
			return err
		}
		if typeDef.Name.String() != "Apis" {
			continue
		}
		indices = append(indices, idx)
	}
	sort.Slice(indices, func(i, j int) bool { return indices[i] < indices[j] })

	for _, idx := range indices {
		typeDef, err := metadata.Tables.TypeDef.At(idx)
		if err != nil {
			return err
		}
		namespace := typeDef.Namespace.String()
		for fieldIdx := range typeDef.FieldList.All() {
			field, err := metadata.Tables.Field.At(fieldIdx)
			if err != nil {
				return err
			}
			constant, hasConstant := context.fieldConstant[fieldIdx]
			if !hasConstant {
				continue
			}
			value, err := constantValueHex(constant)
			if err != nil {
				continue
			}
			signature, err := metadata.FieldSignature(field.Signature)
			if err != nil {
				continue
			}
			goType, err := context.sigTypeGoString(&signature.Type, ArchAll)
			if err != nil {
				continue
			}
			constantIR := ConstantIR{
				Namespace: namespace,
				Name:      escapedUpper(field.Name.String()),
				GoType:    goType,
				Value:     value,
			}
			if !win32GoNameValid(constantIR.Name) {
				continue
			}
			if !yield(constantIR) {
				return nil
			}
		}
	}
	return nil
}

func StreamMethodIRs(metadata *winmd.Metadata, yield func(MethodIR) bool) error {
	context, err := NewContext(metadata)
	if err != nil {
		return err
	}

	indices := make([]winmd.Index, 0, context.Metadata.Tables.MethodDef.Len())
	for idx := range context.Metadata.Tables.MethodDef.Indices() {
		indices = append(indices, idx)
	}
	sort.Slice(indices, func(i, j int) bool { return indices[i] < indices[j] })

	for _, idx := range indices {
		method, err := context.Metadata.Tables.MethodDef.At(idx)
		if err != nil {
			return err
		}
		if context.MethodModuleName(idx) == "" {
			continue
		}
		methodIR, err := context.methodIRFromMethod(idx, method, ArchAll)
		if err != nil {
			return fmt.Errorf("method %s: %w", method.Name, err)
		}
		methodIR.Index = idx
		methodIR.Arch = context.MethodDefSupportedArch(idx)
		if !win32GoNameValid(methodIR.GoName) {
			continue
		}
		if !yield(methodIR) {
			return nil
		}
	}
	return nil
}

func (c *Context) sigTypeGoString(sig *winmd.SigType, arch Arch) (string, error) {
	var body strings.Builder
	if err := c.writeType(StringWriter{Writer: &body}, sig, arch); err != nil {
		return "", err
	}
	return strings.TrimSpace(body.String()), nil
}

func (c *Context) typeIRFromResolved(resolved *resolvedDef, arch Arch) (TypeIR, error) {
	if resolved.IsInterface() {
		return TypeIR{
			GoName:         resolved.GoName,
			Kind:           TypeKindDelegate,
			UnderlyingType: "uintptr",
		}, nil
	}

	switch resolved.def.Extends.Tag {
	case winmd.TypeDefOrRef_TypeRef:
		extendsRef, err := c.Metadata.Tables.TypeRef.At(resolved.def.Extends.Index)
		if err != nil {
			return TypeIR{}, err
		}
		if extendsRef.Namespace.String() == "System" {
			switch extendsRef.Name.String() {
			case "Enum":
				return c.enumIRFromResolved(resolved, arch)
			case "MulticastDelegate":
				return TypeIR{
					GoName:         resolved.GoName,
					Kind:           TypeKindDelegate,
					UnderlyingType: "uintptr",
				}, nil
			}
		}
		if resolved.Native {
			return c.nativeIRFromResolved(resolved, arch)
		}
		return c.structIRFromResolved(resolved, arch)
	case winmd.TypeDefOrRef_Null:
		return c.structIRFromResolved(resolved, arch)
	default:
		return TypeIR{
			GoName:         resolved.GoName,
			Kind:           TypeKindDelegate,
			UnderlyingType: "uintptr",
		}, nil
	}
}

func (c *Context) enumIRFromResolved(resolved *resolvedDef, arch Arch) (TypeIR, error) {
	var underlying string
	members := make([]EnumMemberIR, 0, resolved.def.FieldList.Len())

	for fieldIdx := range resolved.def.FieldList.All() {
		field, err := c.Metadata.Tables.Field.At(fieldIdx)
		if err != nil {
			return TypeIR{}, err
		}
		if field.Name.String() == "value__" {
			signature, err := c.Metadata.FieldSignature(field.Signature)
			if err != nil {
				return TypeIR{}, err
			}
			underlying, err = c.sigTypeGoString(&signature.Type, arch)
			if err != nil {
				return TypeIR{}, err
			}
			continue
		}
		constant, ok := c.fieldConstant[fieldIdx]
		if !ok {
			return TypeIR{}, fmt.Errorf("missing enum constant for %s", field.Name)
		}
		value, err := constantValueHex(constant)
		if err != nil {
			return TypeIR{}, err
		}
		members = append(members, EnumMemberIR{
			Name:  escapedUpper(field.Name.String()),
			Value: value,
		})
	}
	if underlying == "" {
		return TypeIR{}, fmt.Errorf("enum %s missing value__ field", resolved.GoName)
	}
	return TypeIR{
		GoName:         resolved.GoName,
		Kind:           TypeKindEnum,
		UnderlyingType: underlying,
		EnumMembers:    members,
	}, nil
}

func (c *Context) nativeIRFromResolved(resolved *resolvedDef, arch Arch) (TypeIR, error) {
	if resolved.def.FieldList.Start+1 != resolved.def.FieldList.End {
		return TypeIR{}, fmt.Errorf("native typedef %s expected one field", resolved.GoName)
	}
	field, err := c.Metadata.Tables.Field.At(resolved.def.FieldList.Start)
	if err != nil {
		return TypeIR{}, err
	}
	signature, err := c.Metadata.FieldSignature(field.Signature)
	if err != nil {
		return TypeIR{}, err
	}

	signatureType := signature.Type
	if resolved.NativePointer {
		inner, ok := signatureType.Value.(winmd.SigType)
		if !ok {
			return TypeIR{}, fmt.Errorf("native pointer typedef %s expected pointer signature", resolved.GoName)
		}
		signatureType = inner
	}

	underlying, err := c.sigTypeGoString(&signatureType, arch)
	if err != nil {
		return TypeIR{}, err
	}
	return TypeIR{
		GoName:         resolved.GoName,
		Kind:           TypeKindNative,
		UnderlyingType: underlying,
	}, nil
}

func (c *Context) structIRFromResolved(resolved *resolvedDef, arch Arch) (TypeIR, error) {
	usedFieldOffset := make(map[uint32]struct{})
	fields := make([]StructFieldIR, 0, resolved.def.FieldList.Len())

	for fieldIdx := range resolved.def.FieldList.All() {
		if offset, ok := c.fieldOffset[fieldIdx]; ok {
			if _, used := usedFieldOffset[offset]; used {
				continue
			}
			usedFieldOffset[offset] = struct{}{}
		}
		fieldIR, err := c.structFieldIRFromIndex(fieldIdx, arch)
		if err != nil {
			return TypeIR{}, err
		}
		for _, f := range fieldIR {
			fields = append(fields, f)
		}
	}

	return TypeIR{
		GoName:       resolved.GoName,
		Kind:         TypeKindStruct,
		StructFields: fields,
	}, nil
}

func (c *Context) structFieldIRFromIndex(fieldIndex winmd.Index, arch Arch) ([]StructFieldIR, error) {
	field, err := c.Metadata.Tables.Field.At(fieldIndex)
	if err != nil {
		return nil, err
	}
	signature, err := c.Metadata.FieldSignature(field.Signature)
	if err != nil {
		return nil, err
	}

	if signature.Type.Kind == winmd.ElementType_VALUETYPE {
		if coded, ok := signature.Type.Value.(winmd.CodedIndex[winmd.TypeDefOrRefOrSpec]); ok && coded.Tag == winmd.TypeDefOrRefOrSpec_TypeRef {
			ref, err := c.Metadata.Tables.TypeRef.At(coded.Index)
			if err != nil {
				return nil, err
			}
			if ref.ResolutionScope.Tag == winmd.ResolutionScope_TypeRef {
				def, err := c.resolveTypeRef(coded.Index, arch)
				if err != nil && err != errTypeDefNotDefinedInCurrentModule {
					return nil, err
				}
				if def != nil {
					nested, err := c.structIRFromResolved(def, arch)
					if err != nil {
						return nil, err
					}
					return nested.StructFields, nil
				}
			}
		}
	}

	goType, err := c.sigTypeGoString(&signature.Type, arch)
	if err != nil {
		goType = "uintptr"
	}
	return []StructFieldIR{{
		Name:   escapedUpper(field.Name.String()),
		GoType: goType,
	}}, nil
}

func (c *Context) methodIRsCollect(arch Arch) ([]MethodIR, error) {
	entries := make([]MethodIR, 0, c.Metadata.Tables.MethodDef.Len())
	for idx := range c.Metadata.Tables.MethodDef.Indices() {
		method, err := c.Metadata.Tables.MethodDef.At(idx)
		if err != nil {
			return nil, err
		}
		module := c.MethodModuleName(idx)
		if module == "" {
			continue
		}
		methodIR, err := c.methodIRFromMethod(idx, method, arch)
		if err != nil {
			return nil, fmt.Errorf("method %s: %w", method.Name, err)
		}
		if !win32GoNameValid(methodIR.GoName) {
			continue
		}
		entries = append(entries, methodIR)
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].GoName < entries[j].GoName })
	return entries, nil
}

func (c *Context) methodIRFromMethod(methodIndex winmd.Index, method winmd.MethodDef, arch Arch) (MethodIR, error) {
	sig, err := c.Metadata.MethodDefSignature(method.Signature)
	if err != nil {
		return MethodIR{}, err
	}

	params := make([]ParamIR, 0, len(sig.Param))
	for paramRowIndex := range method.ParamList.All() {
		param, err := c.Metadata.Tables.Param.At(paramRowIndex)
		if err != nil {
			return MethodIR{}, err
		}
		if param.Sequence == 0 {
			if param.Flags == 0 && param.Name.String() == "" {
				continue
			}
			return MethodIR{}, fmt.Errorf("unsupported param row sequence 0")
		}
		i := param.Sequence - 1
		if int(i) >= len(sig.Param) {
			return MethodIR{}, fmt.Errorf("param sequence out of range")
		}
		goType, err := c.sigTypeGoString(&sig.Param[i].Type, arch)
		if err != nil {
			goType = "uintptr"
		}
		params = append(params, ParamIR{
			Name:   escapeParam(param.Name.String()),
			GoType: goType,
		})
	}

	methodIR := MethodIR{
		GoName: escapedUpper(method.Name.String()),
		Module: c.MethodModuleName(methodIndex),
		Params: params,
	}
	if sig.RetType.Kind != winmd.SigRetTypeKind_Void {
		retType, err := c.sigTypeGoString(&sig.RetType.Type, arch)
		if err != nil {
			retType = "uintptr"
		}
		methodIR.ReturnType = retType
		methodIR.HasReturn = true
	}
	return methodIR, nil
}

/*
BuildTypeIR resolves a TypeDef and constructs an architecture-specific IR for it.

[Parameters]
typeDefIndex - Module-level TypeDef index. Nested types are not addressable here; resolve them via

	their parent's TypeIR.

arch         - Architecture context to apply when interpreting size-dependent signatures.
*/
func (c *Context) BuildTypeIR(typeDefIndex winmd.Index, arch Arch) (TypeIR, error) {
	resolved, err := c.resolveTypeDef(typeDefIndex)
	if err != nil {
		return TypeIR{}, err
	}
	typeIR, err := c.typeIRFromResolved(resolved, arch)
	if err != nil {
		return TypeIR{}, err
	}
	typeIR.Index = typeDefIndex
	typeIR.Arch = c.TypeDefSupportedArch(typeDefIndex)
	typeIR.Namespace = resolved.Namespace.String()
	return typeIR, nil
}

/*
BuildMethodIR constructs an architecture-specific IR for a MethodDef row.
*/
func (c *Context) BuildMethodIR(methodIndex winmd.Index, arch Arch) (MethodIR, error) {
	method, err := c.Metadata.Tables.MethodDef.At(methodIndex)
	if err != nil {
		return MethodIR{}, err
	}
	methodIR, err := c.methodIRFromMethod(methodIndex, method, arch)
	if err != nil {
		return MethodIR{}, err
	}
	methodIR.Index = methodIndex
	methodIR.Arch = c.MethodDefSupportedArch(methodIndex)
	return methodIR, nil
}

/*
TypeDefNamespace returns the source namespace string for a TypeDef row.
*/
func (c *Context) TypeDefNamespace(typeDefIndex winmd.Index) (string, error) {
	def, err := c.Metadata.Tables.TypeDef.At(typeDefIndex)
	if err != nil {
		return "", err
	}
	return def.Namespace.String(), nil
}

/*
TypeDefGoName returns the Go-exportable name for a TypeDef row.
*/
func (c *Context) TypeDefGoName(typeDefIndex winmd.Index) (string, error) {
	resolved, err := c.resolveTypeDef(typeDefIndex)
	if err != nil {
		return "", err
	}
	return resolved.GoName, nil
}

/*
TypeDefIsNestedPublic reports whether a TypeDef is a nested public type (i.e. defined inside
another TypeDef and not addressable at module scope).
*/
func (c *Context) TypeDefIsNestedPublic(typeDefIndex winmd.Index) (bool, error) {
	def, err := c.Metadata.Tables.TypeDef.At(typeDefIndex)
	if err != nil {
		return false, err
	}
	return def.Flags&winmd.TypeAttributes_NestedPublic != 0, nil
}

/*
MethodReferencedTypeDefs returns every module-level TypeDef index referenced from a method's
parameter or return signature for the given architecture context.

[Context]
Used to seed the type closure walker. Pointers to System types (for example System.Guid) are
not part of the winmd TypeDef table and are silently skipped here; callers must handle
synthetic stand-ins (such as the Foundation Guid shim) separately.
*/
func (c *Context) MethodReferencedTypeDefs(methodIndex winmd.Index, arch Arch) ([]winmd.Index, error) {
	method, err := c.Metadata.Tables.MethodDef.At(methodIndex)
	if err != nil {
		return nil, err
	}
	signature, err := c.Metadata.MethodDefSignature(method.Signature)
	if err != nil {
		return nil, err
	}
	visited := map[winmd.Index]struct{}{}
	visit := func(idx winmd.Index) {
		visited[idx] = struct{}{}
	}
	for paramIdx := range signature.Param {
		if err := c.walkSigTypeRefs(&signature.Param[paramIdx].Type, arch, visit); err != nil {
			return nil, err
		}
	}
	if signature.RetType.Kind != winmd.SigRetTypeKind_Void {
		if err := c.walkSigTypeRefs(&signature.RetType.Type, arch, visit); err != nil {
			return nil, err
		}
	}
	return sortedIndexKeys(visited), nil
}

/*
TypeReferencedTypeDefs returns every module-level TypeDef index referenced from a TypeDef's
field signatures (or underlying type for enums and native typedefs) for the given architecture.

[Context]
Anonymous nested struct types referenced via a parent's field are inlined: their TypeDef index is
not returned (it has no module-level identity) but the types referenced by their own fields are
walked so they appear in the closure.
*/
func (c *Context) TypeReferencedTypeDefs(typeDefIndex winmd.Index, arch Arch) ([]winmd.Index, error) {
	def, err := c.Metadata.Tables.TypeDef.At(typeDefIndex)
	if err != nil {
		return nil, err
	}
	visited := map[winmd.Index]struct{}{}
	visit := func(idx winmd.Index) {
		visited[idx] = struct{}{}
	}
	for fieldIdx := range def.FieldList.All() {
		if err := c.walkFieldRefs(fieldIdx, arch, visit); err != nil {
			return nil, err
		}
	}
	return sortedIndexKeys(visited), nil
}

func (c *Context) walkFieldRefs(fieldIndex winmd.Index, arch Arch, visit func(winmd.Index)) error {
	field, err := c.Metadata.Tables.Field.At(fieldIndex)
	if err != nil {
		return err
	}
	if field.Name.String() == "value__" {
		return nil
	}
	signature, err := c.Metadata.FieldSignature(field.Signature)
	if err != nil {
		return nil
	}
	return c.walkSigTypeRefs(&signature.Type, arch, visit)
}

func (c *Context) walkSigTypeRefs(p *winmd.SigType, arch Arch, visit func(winmd.Index)) error {
	if p == nil {
		return nil
	}
	if p.Kind == winmd.ElementType_PTR {
		if t, ok := p.Value.(winmd.SigType); ok && t.Kind == winmd.ElementType_VOID {
			return nil
		}
	}
	switch v := p.Value.(type) {
	case winmd.CodedIndex[winmd.TypeDefOrRefOrSpec]:
		switch v.Tag {
		case winmd.TypeDefOrRefOrSpec_TypeDef:
			visit(v.Index)
		case winmd.TypeDefOrRefOrSpec_TypeRef:
			ref, err := c.Metadata.Tables.TypeRef.At(v.Index)
			if err != nil {
				return err
			}
			if ref.ResolutionScope.Tag == winmd.ResolutionScope_TypeRef {
				def, err := c.resolveTypeRef(v.Index, arch)
				if err != nil && !errors.Is(err, errTypeDefNotDefinedInCurrentModule) {
					return err
				}
				if def != nil {
					for fieldIdx := range def.def.FieldList.All() {
						if err := c.walkFieldRefs(fieldIdx, arch, visit); err != nil {
							return err
						}
					}
				}
			} else {
				def, err := c.resolveTypeRef(v.Index, arch)
				if err != nil && !errors.Is(err, errTypeDefNotDefinedInCurrentModule) {
					return err
				}
				if def != nil {
					visit(def.Index)
				}
			}
		}
	case winmd.SigType:
		return c.walkSigTypeRefs(&v, arch, visit)
	case winmd.SigArray:
		return c.walkSigTypeRefs(&v.Type, arch, visit)
	}
	return nil
}

func sortedIndexKeys(set map[winmd.Index]struct{}) []winmd.Index {
	out := make([]winmd.Index, 0, len(set))
	for idx := range set {
		out = append(out, idx)
	}
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}

func win32GoNameValid(name string) bool {
	for _, r := range name {
		switch {
		case r >= 'A' && r <= 'Z':
		case r >= 'a' && r <= 'z':
		case r >= '0' && r <= '9':
		case r == '_':
		default:
			return false
		}
	}
	return name != ""
}

func constantValueHex(constant winmd.Constant) (string, error) {
	switch constant.Type {
	case winmd.ElementType_I1:
		v := int8(constant.Value[0])
		if v < 0 {
			return "-0x" + strconv.FormatInt(int64(-v), 16), nil
		}
		return "0x" + strconv.FormatUint(uint64(v), 16), nil
	case winmd.ElementType_I2:
		v := int16(binary.LittleEndian.Uint16(constant.Value))
		if v < 0 {
			return "-0x" + strconv.FormatInt(int64(-v), 16), nil
		}
		return "0x" + strconv.FormatUint(uint64(v), 16), nil
	case winmd.ElementType_I4:
		v := int32(binary.LittleEndian.Uint32(constant.Value))
		if v < 0 {
			return strconv.FormatInt(int64(v), 10), nil
		}
		return "0x" + strconv.FormatUint(uint64(v), 16), nil
	case winmd.ElementType_I8:
		v := int64(binary.LittleEndian.Uint64(constant.Value))
		if v < 0 {
			return "-0x" + strconv.FormatInt(-v, 16), nil
		}
		return "0x" + strconv.FormatUint(uint64(v), 16), nil
	case winmd.ElementType_U1:
		return "0x" + strconv.FormatUint(uint64(constant.Value[0]), 16), nil
	case winmd.ElementType_U2:
		return "0x" + strconv.FormatUint(uint64(binary.LittleEndian.Uint16(constant.Value)), 16), nil
	case winmd.ElementType_U4:
		return "0x" + strconv.FormatUint(uint64(binary.LittleEndian.Uint32(constant.Value)), 16), nil
	case winmd.ElementType_U8:
		return "0x" + strconv.FormatUint(binary.LittleEndian.Uint64(constant.Value), 16), nil
	default:
		return "", fmt.Errorf("unsupported constant type %v", constant.Type)
	}
}
