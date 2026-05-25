//go:build linux

package syscore

import (
	"syscore/window/wayland/bindings"
	"syscore/window/wayland/loader"
)

/*
SYSCORE_Window_Wayland_Module holds a loaded libwayland-client.so.0 handle.

[Context]
Obtain via SYSCORE_Window_Wayland_ModuleLoad before binding commands. Linux only.
*/
type SYSCORE_Window_Wayland_Module = loader.WaylandModule

/*
SYSCORE_Window_Wayland_Commands holds bound Wayland client function pointers.

[Context]
After SYSCORE_Window_Wayland_CommandsLoad, invoke wl_* behavior through documented fields on
loader.WaylandCommands. Interface globals are resolved with SYSCORE_Window_Wayland_ModuleSymbolResolve.
*/
type SYSCORE_Window_Wayland_Commands = loader.WaylandCommands

/*
SYSCORE_Window_Wayland_CommandManifest lists Wayland entry points for selective binding.
*/
type SYSCORE_Window_Wayland_CommandManifest = loader.WaylandCommandManifest

const (
	SYSCORE_Window_Wayland_OpcodeDisplayGetRegistry      = bindings.WL_DISPLAY_GET_REGISTRY_OPCODE
	SYSCORE_Window_Wayland_OpcodeRegistryBind            = bindings.WL_REGISTRY_BIND_OPCODE
	SYSCORE_Window_Wayland_OpcodeCompositorCreateSurface = bindings.WL_COMPOSITOR_CREATE_SURFACE_OPCODE
	SYSCORE_Window_Wayland_OpcodeSurfaceCommit           = bindings.WL_SURFACE_COMMIT_OPCODE
	SYSCORE_Window_Wayland_SymbolRegistryInterface       = bindings.WL_REGISTRY_INTERFACE_SYMBOL
	SYSCORE_Window_Wayland_SymbolCompositorInterface     = bindings.WL_COMPOSITOR_INTERFACE_SYMBOL
	SYSCORE_Window_Wayland_SymbolSurfaceInterface        = bindings.WL_SURFACE_INTERFACE_SYMBOL
)

type (
	// SYSCORE_Window_Wayland_Display is an opaque struct wl_display*.
	SYSCORE_Window_Wayland_Display = bindings.WlDisplay
	// SYSCORE_Window_Wayland_Proxy is an opaque struct wl_proxy*.
	SYSCORE_Window_Wayland_Proxy = bindings.WlProxy
	// SYSCORE_Window_Wayland_Registry is an opaque struct wl_registry*.
	SYSCORE_Window_Wayland_Registry = bindings.WlRegistry
	// SYSCORE_Window_Wayland_Compositor is an opaque struct wl_compositor*.
	SYSCORE_Window_Wayland_Compositor = bindings.WlCompositor
	// SYSCORE_Window_Wayland_Surface is an opaque struct wl_surface*.
	SYSCORE_Window_Wayland_Surface = bindings.WlSurface
	// SYSCORE_Window_Wayland_RegistryListener is the wl_registry_listener vtable layout for callbacks.
	SYSCORE_Window_Wayland_RegistryListener = bindings.WlRegistryListener
	// SYSCORE_Window_Wayland_PFN_display_connect is the C type for wl_display_connect.
	SYSCORE_Window_Wayland_PFN_display_connect = bindings.PFN_wl_display_connect
	// SYSCORE_Window_Wayland_PFN_proxy_marshal_constructor is the C type for wl_proxy_marshal_constructor.
	SYSCORE_Window_Wayland_PFN_proxy_marshal_constructor = bindings.PFN_wl_proxy_marshal_constructor
	// SYSCORE_Window_Wayland_PFN_proxy_marshal_constructor_ver is the C type for wl_proxy_marshal_constructor_versioned.
	SYSCORE_Window_Wayland_PFN_proxy_marshal_constructor_ver = bindings.PFN_wl_proxy_marshal_constructor_versioned
	// SYSCORE_Window_Wayland_PFN_proxy_add_listener is the C type for wl_proxy_add_listener.
	SYSCORE_Window_Wayland_PFN_proxy_add_listener = bindings.PFN_wl_proxy_add_listener
	// SYSCORE_Window_Wayland_PFN_display_roundtrip is the C type for wl_display_roundtrip.
	SYSCORE_Window_Wayland_PFN_display_roundtrip = bindings.PFN_wl_display_roundtrip
	// SYSCORE_Window_Wayland_PFN_display_dispatch_pending is the C type for wl_display_dispatch_pending.
	SYSCORE_Window_Wayland_PFN_display_dispatch_pending = bindings.PFN_wl_display_dispatch_pending
	// SYSCORE_Window_Wayland_PFN_proxy_marshal is the C type for wl_proxy_marshal.
	SYSCORE_Window_Wayland_PFN_proxy_marshal = bindings.PFN_wl_proxy_marshal
	// SYSCORE_Window_Wayland_PFN_proxy_destroy is the C type for wl_proxy_destroy.
	SYSCORE_Window_Wayland_PFN_proxy_destroy = bindings.PFN_wl_proxy_destroy
	// SYSCORE_Window_Wayland_PFN_display_disconnect is the C type for wl_display_disconnect.
	SYSCORE_Window_Wayland_PFN_display_disconnect = bindings.PFN_wl_display_disconnect
)

/*
SYSCORE_Window_Wayland_ModuleLoad loads libwayland-client.so.0.

[Context]
Use when SYSCORE_Window_DisplayAPIDetect returns Wayland. Requires WAYLAND_DISPLAY in the
environment for a successful connection later via DisplayConnect.

[Returns]
Module handle on success.

[Errors]
Returns an error when not on linux or when dlopen fails.

[Side Effects]
Loads libwayland-client into the process.
*/
func SYSCORE_Window_Wayland_ModuleLoad() (SYSCORE_Window_Wayland_Module, error) {
	return loader.WaylandModuleLoad()
}

/*
SYSCORE_Window_Wayland_ModuleLibraryName returns the soname used for WaylandModuleLoad.
*/
func SYSCORE_Window_Wayland_ModuleLibraryName() string {
	return loader.WaylandModuleLibraryName()
}

/*
SYSCORE_Window_Wayland_ModuleSymbolResolve resolves exported globals such as wl_registry_interface.

[Parameters]
module - Loaded Wayland module.
symbolName - Symbol name (for example "wl_compositor_interface").

[Returns]
Address of the exported wl_interface or other global.

[Errors]
Returns an error if the symbol is missing.
*/
func SYSCORE_Window_Wayland_ModuleSymbolResolve(module SYSCORE_Window_Wayland_Module, symbolName string) (uintptr, error) {
	return loader.WaylandModuleSymbolResolve(module, symbolName)
}

/*
SYSCORE_Window_Wayland_CommandsLoad binds all supported Wayland client entry points.

[Parameters]
module - Loaded module.
commands - Non-nil command holder.

[Errors]
Returns an error if commands is nil, not on linux, or any bind fails.
*/
func SYSCORE_Window_Wayland_CommandsLoad(module SYSCORE_Window_Wayland_Module, commands *SYSCORE_Window_Wayland_Commands) error {
	return loader.WaylandCommandsLoad(module, commands)
}

/*
SYSCORE_Window_Wayland_CommandManifestReset clears manifest.
*/
func SYSCORE_Window_Wayland_CommandManifestReset(manifest *SYSCORE_Window_Wayland_CommandManifest) {
	loader.WaylandCommandManifestReset(manifest)
}

/*
SYSCORE_Window_Wayland_CommandManifestAddDisplayConnect registers wl_display_connect for selective loading.
*/
func SYSCORE_Window_Wayland_CommandManifestAddDisplayConnect(manifest *SYSCORE_Window_Wayland_CommandManifest, target *SYSCORE_Window_Wayland_PFN_display_connect) error {
	return loader.WaylandCommandManifestAddDisplayConnect(manifest, target)
}

/*
SYSCORE_Window_Wayland_CommandManifestAddProxyMarshalConstructor registers wl_proxy_marshal_constructor for selective loading.
*/
func SYSCORE_Window_Wayland_CommandManifestAddProxyMarshalConstructor(manifest *SYSCORE_Window_Wayland_CommandManifest, target *SYSCORE_Window_Wayland_PFN_proxy_marshal_constructor) error {
	return loader.WaylandCommandManifestAddProxyMarshalConstructor(manifest, target)
}

/*
SYSCORE_Window_Wayland_CommandManifestAddProxyMarshalConstructorVersioned registers wl_proxy_marshal_constructor_versioned for selective loading.
*/
func SYSCORE_Window_Wayland_CommandManifestAddProxyMarshalConstructorVersioned(manifest *SYSCORE_Window_Wayland_CommandManifest, target *SYSCORE_Window_Wayland_PFN_proxy_marshal_constructor_ver) error {
	return loader.WaylandCommandManifestAddProxyMarshalConstructorVersioned(manifest, target)
}

/*
SYSCORE_Window_Wayland_CommandManifestAddProxyAddListener registers wl_proxy_add_listener for selective loading.
*/
func SYSCORE_Window_Wayland_CommandManifestAddProxyAddListener(manifest *SYSCORE_Window_Wayland_CommandManifest, target *SYSCORE_Window_Wayland_PFN_proxy_add_listener) error {
	return loader.WaylandCommandManifestAddProxyAddListener(manifest, target)
}

/*
SYSCORE_Window_Wayland_CommandManifestAddDisplayRoundtrip registers wl_display_roundtrip for selective loading.
*/
func SYSCORE_Window_Wayland_CommandManifestAddDisplayRoundtrip(manifest *SYSCORE_Window_Wayland_CommandManifest, target *SYSCORE_Window_Wayland_PFN_display_roundtrip) error {
	return loader.WaylandCommandManifestAddDisplayRoundtrip(manifest, target)
}

/*
SYSCORE_Window_Wayland_CommandsLoadManifest binds only manifest-listed Wayland entry points.

[Errors]
Returns an error if manifest or commands is nil or any bind fails.
*/
func SYSCORE_Window_Wayland_CommandsLoadManifest(module SYSCORE_Window_Wayland_Module, manifest *SYSCORE_Window_Wayland_CommandManifest, commands *SYSCORE_Window_Wayland_Commands) error {
	return loader.WaylandCommandsLoadManifest(module, manifest, commands)
}

/*
SYSCORE_Window_Wayland_CommandManifestAddProxyMarshal registers wl_proxy_marshal for selective loading.
*/
func SYSCORE_Window_Wayland_CommandManifestAddProxyMarshal(manifest *SYSCORE_Window_Wayland_CommandManifest, target *SYSCORE_Window_Wayland_PFN_proxy_marshal) error {
	return loader.WaylandCommandManifestAddProxyMarshal(manifest, target)
}

/*
SYSCORE_Window_Wayland_CommandManifestAddProxyDestroy registers wl_proxy_destroy for selective loading.
*/
func SYSCORE_Window_Wayland_CommandManifestAddProxyDestroy(manifest *SYSCORE_Window_Wayland_CommandManifest, target *SYSCORE_Window_Wayland_PFN_proxy_destroy) error {
	return loader.WaylandCommandManifestAddProxyDestroy(manifest, target)
}

/*
SYSCORE_Window_Wayland_CommandManifestAddDisplayDisconnect registers wl_display_disconnect for selective loading.
*/
func SYSCORE_Window_Wayland_CommandManifestAddDisplayDisconnect(manifest *SYSCORE_Window_Wayland_CommandManifest, target *SYSCORE_Window_Wayland_PFN_display_disconnect) error {
	return loader.WaylandCommandManifestAddDisplayDisconnect(manifest, target)
}
