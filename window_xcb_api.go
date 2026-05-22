package syscore

import (
	"syscore/window/xcb/bindings"
	"syscore/window/xcb/loader"
)

type SYSCORE_Window_Xcb_Module = loader.XcbModule
type SYSCORE_Window_Xcb_Commands = loader.XcbCommands
type SYSCORE_Window_Xcb_CommandManifest = loader.XcbCommandManifest

type (
	SYSCORE_Window_Xcb_ConnectionT          = bindings.XcbConnectionT
	SYSCORE_Window_Xcb_WindowT              = bindings.XcbWindowT
	SYSCORE_Window_Xcb_VoidCookieT          = bindings.XcbVoidCookieT
	SYSCORE_Window_Xcb_ScreenIteratorT      = bindings.XcbScreenIteratorT
	SYSCORE_Window_Xcb_PFN_connect          = bindings.PFN_xcb_connect
	SYSCORE_Window_Xcb_PFN_get_setup        = bindings.PFN_xcb_get_setup
	SYSCORE_Window_Xcb_PFN_setup_roots_iter = bindings.PFN_xcb_setup_roots_iterator
	SYSCORE_Window_Xcb_PFN_generate_id      = bindings.PFN_xcb_generate_id
	SYSCORE_Window_Xcb_PFN_create_window    = bindings.PFN_xcb_create_window
	SYSCORE_Window_Xcb_PFN_map_window       = bindings.PFN_xcb_map_window
	SYSCORE_Window_Xcb_PFN_flush            = bindings.PFN_xcb_flush
)

func SYSCORE_Window_Xcb_ModuleLoad() (SYSCORE_Window_Xcb_Module, error) {
	return loader.XcbModuleLoad()
}

func SYSCORE_Window_Xcb_ModuleLibraryName() string {
	return loader.XcbModuleLibraryName()
}

func SYSCORE_Window_Xcb_ModuleSymbolResolve(module SYSCORE_Window_Xcb_Module, symbolName string) (uintptr, error) {
	return loader.XcbModuleSymbolResolve(module, symbolName)
}

func SYSCORE_Window_Xcb_CommandsLoad(module SYSCORE_Window_Xcb_Module, commands *SYSCORE_Window_Xcb_Commands) error {
	return loader.XcbCommandsLoad(module, commands)
}

func SYSCORE_Window_Xcb_CommandManifestReset(manifest *SYSCORE_Window_Xcb_CommandManifest) {
	loader.XcbCommandManifestReset(manifest)
}

func SYSCORE_Window_Xcb_CommandManifestAddConnect(manifest *SYSCORE_Window_Xcb_CommandManifest, target *SYSCORE_Window_Xcb_PFN_connect) error {
	return loader.XcbCommandManifestAddConnect(manifest, target)
}

func SYSCORE_Window_Xcb_CommandManifestAddGetSetup(manifest *SYSCORE_Window_Xcb_CommandManifest, target *SYSCORE_Window_Xcb_PFN_get_setup) error {
	return loader.XcbCommandManifestAddGetSetup(manifest, target)
}

func SYSCORE_Window_Xcb_CommandManifestAddSetupRootsIterator(manifest *SYSCORE_Window_Xcb_CommandManifest, target *SYSCORE_Window_Xcb_PFN_setup_roots_iter) error {
	return loader.XcbCommandManifestAddSetupRootsIterator(manifest, target)
}

func SYSCORE_Window_Xcb_CommandManifestAddGenerateId(manifest *SYSCORE_Window_Xcb_CommandManifest, target *SYSCORE_Window_Xcb_PFN_generate_id) error {
	return loader.XcbCommandManifestAddGenerateId(manifest, target)
}

func SYSCORE_Window_Xcb_CommandManifestAddCreateWindow(manifest *SYSCORE_Window_Xcb_CommandManifest, target *SYSCORE_Window_Xcb_PFN_create_window) error {
	return loader.XcbCommandManifestAddCreateWindow(manifest, target)
}

func SYSCORE_Window_Xcb_CommandManifestAddMapWindow(manifest *SYSCORE_Window_Xcb_CommandManifest, target *SYSCORE_Window_Xcb_PFN_map_window) error {
	return loader.XcbCommandManifestAddMapWindow(manifest, target)
}

func SYSCORE_Window_Xcb_CommandManifestAddFlush(manifest *SYSCORE_Window_Xcb_CommandManifest, target *SYSCORE_Window_Xcb_PFN_flush) error {
	return loader.XcbCommandManifestAddFlush(manifest, target)
}

func SYSCORE_Window_Xcb_CommandsLoadManifest(module SYSCORE_Window_Xcb_Module, manifest *SYSCORE_Window_Xcb_CommandManifest, commands *SYSCORE_Window_Xcb_Commands) error {
	return loader.XcbCommandsLoadManifest(module, manifest, commands)
}
