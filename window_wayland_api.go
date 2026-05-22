package syscore

import (
	"syscore/window/wayland/bindings"
	"syscore/window/wayland/loader"
)

type SYSCORE_Window_Wayland_Module = loader.WaylandModule
type SYSCORE_Window_Wayland_Commands = loader.WaylandCommands
type SYSCORE_Window_Wayland_CommandManifest = loader.WaylandCommandManifest

type (
	SYSCORE_Window_Wayland_Display                           = bindings.WlDisplay
	SYSCORE_Window_Wayland_Proxy                             = bindings.WlProxy
	SYSCORE_Window_Wayland_Registry                          = bindings.WlRegistry
	SYSCORE_Window_Wayland_Compositor                        = bindings.WlCompositor
	SYSCORE_Window_Wayland_Surface                           = bindings.WlSurface
	SYSCORE_Window_Wayland_RegistryListener                  = bindings.WlRegistryListener
	SYSCORE_Window_Wayland_PFN_display_connect               = bindings.PFN_wl_display_connect
	SYSCORE_Window_Wayland_PFN_proxy_marshal_constructor     = bindings.PFN_wl_proxy_marshal_constructor
	SYSCORE_Window_Wayland_PFN_proxy_marshal_constructor_ver = bindings.PFN_wl_proxy_marshal_constructor_versioned
	SYSCORE_Window_Wayland_PFN_proxy_add_listener            = bindings.PFN_wl_proxy_add_listener
	SYSCORE_Window_Wayland_PFN_display_roundtrip             = bindings.PFN_wl_display_roundtrip
)

func SYSCORE_Window_Wayland_ModuleLoad() (SYSCORE_Window_Wayland_Module, error) {
	return loader.WaylandModuleLoad()
}

func SYSCORE_Window_Wayland_ModuleLibraryName() string {
	return loader.WaylandModuleLibraryName()
}

func SYSCORE_Window_Wayland_ModuleSymbolResolve(module SYSCORE_Window_Wayland_Module, symbolName string) (uintptr, error) {
	return loader.WaylandModuleSymbolResolve(module, symbolName)
}

func SYSCORE_Window_Wayland_CommandsLoad(module SYSCORE_Window_Wayland_Module, commands *SYSCORE_Window_Wayland_Commands) error {
	return loader.WaylandCommandsLoad(module, commands)
}

func SYSCORE_Window_Wayland_CommandManifestReset(manifest *SYSCORE_Window_Wayland_CommandManifest) {
	loader.WaylandCommandManifestReset(manifest)
}

func SYSCORE_Window_Wayland_CommandManifestAddDisplayConnect(manifest *SYSCORE_Window_Wayland_CommandManifest, target *SYSCORE_Window_Wayland_PFN_display_connect) error {
	return loader.WaylandCommandManifestAddDisplayConnect(manifest, target)
}

func SYSCORE_Window_Wayland_CommandManifestAddProxyMarshalConstructor(manifest *SYSCORE_Window_Wayland_CommandManifest, target *SYSCORE_Window_Wayland_PFN_proxy_marshal_constructor) error {
	return loader.WaylandCommandManifestAddProxyMarshalConstructor(manifest, target)
}

func SYSCORE_Window_Wayland_CommandManifestAddProxyMarshalConstructorVersioned(manifest *SYSCORE_Window_Wayland_CommandManifest, target *SYSCORE_Window_Wayland_PFN_proxy_marshal_constructor_ver) error {
	return loader.WaylandCommandManifestAddProxyMarshalConstructorVersioned(manifest, target)
}

func SYSCORE_Window_Wayland_CommandManifestAddProxyAddListener(manifest *SYSCORE_Window_Wayland_CommandManifest, target *SYSCORE_Window_Wayland_PFN_proxy_add_listener) error {
	return loader.WaylandCommandManifestAddProxyAddListener(manifest, target)
}

func SYSCORE_Window_Wayland_CommandManifestAddDisplayRoundtrip(manifest *SYSCORE_Window_Wayland_CommandManifest, target *SYSCORE_Window_Wayland_PFN_display_roundtrip) error {
	return loader.WaylandCommandManifestAddDisplayRoundtrip(manifest, target)
}

func SYSCORE_Window_Wayland_CommandsLoadManifest(module SYSCORE_Window_Wayland_Module, manifest *SYSCORE_Window_Wayland_CommandManifest, commands *SYSCORE_Window_Wayland_Commands) error {
	return loader.WaylandCommandsLoadManifest(module, manifest, commands)
}
