package bindings

/*
WlDisplay is an opaque struct wl_display* (Wayland display connection).
*/
type WlDisplay uintptr

/*
WlProxy is an opaque struct wl_proxy* (base Wayland protocol object).
*/
type WlProxy uintptr

/*
WlRegistry is an opaque struct wl_registry* (global object registry).
*/
type WlRegistry uintptr

/*
WlCompositor is an opaque struct wl_compositor* (compositor global).
*/
type WlCompositor uintptr

/*
WlSurface is an opaque struct wl_surface* (client surface for Vulkan/shell integration).
*/
type WlSurface uintptr

/*
WlRegistryListener is the wl_registry_listener vtable passed to wl_proxy_add_listener.

Global and GlobalRemove are function pointers (use purego.NewCallback); unset entries must be 0.
*/
type WlRegistryListener struct {
	Global       uintptr
	GlobalRemove uintptr
}

/*
PFN_wl_display_connect connects to a Wayland compositor.

[Context]
Maps to wl_display_connect. Pass nil name to use WAYLAND_DISPLAY.

[Parameters]
name - Compositor socket name or nil.

[Returns]
Display handle, or 0 on failure.
*/
type PFN_wl_display_connect func(name *byte) WlDisplay

/*
PFN_wl_proxy_marshal_constructor marshals a request that constructs a new proxy object.

[Context]
Maps to wl_proxy_marshal_constructor. Used for wl_display_get_registry (opcode 1) and
wl_compositor_create_surface (opcode 0) when interface pointers are resolved from the library.

[Parameters]
proxy - Object implementing the interface (display or compositor as WlProxy).
opcode - Interface-specific request opcode.
interfacePtr - Address of wl_interface (for example wl_registry_interface).
args - Trailing request arguments expanded by purego (often nil or name for bind).

[Returns]
New proxy as WlProxy, or 0 on failure.
*/
type PFN_wl_proxy_marshal_constructor func(proxy WlProxy, opcode uint32, interfacePtr uintptr, args ...any) WlProxy

/*
PFN_wl_proxy_marshal_constructor_versioned marshals a constructor request with an interface version.

[Context]
Maps to wl_proxy_marshal_constructor_versioned. Used for wl_registry_bind when binding
wl_compositor from the registry global event.

[Parameters]
proxy - Registry proxy.
opcode - Request opcode (WL_REGISTRY_BIND_OPCODE).
interfacePtr - Target wl_interface pointer.
version - Interface version from the global advertisement.
args - Bind arguments (global name as uint32, then nil).

[Returns]
Bound global proxy.
*/
type PFN_wl_proxy_marshal_constructor_versioned func(proxy WlProxy, opcode uint32, interfacePtr uintptr, version uint32, args ...any) WlProxy

/*
PFN_wl_proxy_add_listener registers a listener vtable on a proxy.

[Context]
Maps to wl_proxy_add_listener. Used on the registry before wl_display_roundtrip.

[Parameters]
proxy - Target proxy (registry).
implementation - Pointer to wl_registry_listener in memory.
data - User data pointer passed to callbacks.

[Returns]
0 on success, negative on error.
*/
type PFN_wl_proxy_add_listener func(proxy WlProxy, implementation uintptr, data uintptr) int32

/*
PFN_wl_proxy_marshal marshals a request on an existing Wayland proxy.

[Context]
Maps to wl_proxy_marshal. Used for xdg_toplevel.set_title, xdg_surface.get_toplevel,
xdg_surface.ack_configure, and other non-constructor requests.

[Parameters]
proxy - Target proxy.
opcode - Interface-specific request opcode.
args - Trailing request arguments expanded by purego.

[Returns]
Return value depends on the request (often unused).
*/
type PFN_wl_proxy_marshal func(proxy WlProxy, opcode uint32, args ...any) WlProxy

/*
PFN_wl_display_roundtrip blocks until pending requests are processed and events dispatched.

[Context]
Maps to wl_display_roundtrip. Triggers registry global callbacks after add_listener.

[Parameters]
display - Display connection.

[Returns]
0 on success, negative on error.
*/
type PFN_wl_display_roundtrip func(display WlDisplay) int32
