package bindings

type WlDisplay uintptr
type WlProxy uintptr
type WlRegistry uintptr
type WlCompositor uintptr
type WlSurface uintptr

type WlRegistryListener struct {
	Global       uintptr
	GlobalRemove uintptr
}

type PFN_wl_display_connect func(name *byte) WlDisplay

type PFN_wl_proxy_marshal_constructor func(proxy WlProxy, opcode uint32, interfacePtr uintptr, args ...any) WlProxy

type PFN_wl_proxy_marshal_constructor_versioned func(proxy WlProxy, opcode uint32, interfacePtr uintptr, version uint32, args ...any) WlProxy

type PFN_wl_proxy_add_listener func(proxy WlProxy, implementation uintptr, data uintptr) int32

type PFN_wl_display_roundtrip func(display WlDisplay) int32
