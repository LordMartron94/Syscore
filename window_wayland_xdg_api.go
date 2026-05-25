//go:build linux

package syscore

import (
	"syscore/window/wayland/xdg/bindings"
	"syscore/window/wayland/xdg/loader"
)

/*
SYSCORE_Window_Wayland_Xdg_Module holds the embedded xdg-shell protocol library.

[Context]
Load via SYSCORE_Window_Wayland_Xdg_ModuleLoad, then resolve wl_interface globals with
SYSCORE_Window_Wayland_Xdg_ModuleSymbolResolve before xdg_wm_base bind and surface setup.
Linux only.
*/
type SYSCORE_Window_Wayland_Xdg_Module = loader.XdgModule

const (
	SYSCORE_Window_Wayland_Xdg_GlobalNameWmBase            = bindings.XDG_WM_BASE_GLOBAL_NAME
	SYSCORE_Window_Wayland_Xdg_SymbolWmBaseInterface       = bindings.XDG_WM_BASE_INTERFACE_SYMBOL
	SYSCORE_Window_Wayland_Xdg_SymbolSurfaceInterface      = bindings.XDG_SURFACE_INTERFACE_SYMBOL
	SYSCORE_Window_Wayland_Xdg_SymbolToplevelInterface     = bindings.XDG_TOPLEVEL_INTERFACE_SYMBOL
	SYSCORE_Window_Wayland_Xdg_OpcodeWmBaseGetXdgSurface   = bindings.XDG_WM_BASE_GET_XDG_SURFACE_OPCODE
	SYSCORE_Window_Wayland_Xdg_OpcodeSurfaceGetToplevel    = bindings.XDG_SURFACE_GET_TOPLEVEL_OPCODE
	SYSCORE_Window_Wayland_Xdg_OpcodeSurfaceAckConfigure   = bindings.XDG_SURFACE_ACK_CONFIGURE_OPCODE
	SYSCORE_Window_Wayland_Xdg_OpcodeSurfaceConfigureEvent = bindings.XDG_SURFACE_CONFIGURE_EVENT_OPCODE
	SYSCORE_Window_Wayland_Xdg_OpcodeToplevelSetTitle      = bindings.XDG_TOPLEVEL_SET_TITLE_OPCODE
	SYSCORE_Window_Wayland_Xdg_OpcodeToplevelSetMinSize    = bindings.XDG_TOPLEVEL_SET_MIN_SIZE_OPCODE
	SYSCORE_Window_Wayland_Xdg_OpcodeToplevelSetMaxSize    = bindings.XDG_TOPLEVEL_SET_MAX_SIZE_OPCODE
)

type (
	// SYSCORE_Window_Wayland_Xdg_WmBase is an opaque xdg_wm_base*.
	SYSCORE_Window_Wayland_Xdg_WmBase = bindings.XdgWmBase
	// SYSCORE_Window_Wayland_Xdg_Surface is an opaque xdg_surface*.
	SYSCORE_Window_Wayland_Xdg_Surface = bindings.XdgSurface
	// SYSCORE_Window_Wayland_Xdg_Toplevel is an opaque xdg_toplevel*.
	SYSCORE_Window_Wayland_Xdg_Toplevel = bindings.XdgToplevel
	// SYSCORE_Window_Wayland_Xdg_SurfaceListener is the xdg_surface_listener layout for ProxyAddListener.
	SYSCORE_Window_Wayland_Xdg_SurfaceListener = bindings.XdgSurfaceListener
)

/*
SYSCORE_Window_Wayland_Xdg_ModuleLoad materializes and loads the embedded xdg-shell library.

[Returns]
Module handle on success.

[Errors]
Returns an error when not on linux or dlopen fails.

[Side Effects]
Writes the embedded .so to the process temp directory, then loads it.
*/
func SYSCORE_Window_Wayland_Xdg_ModuleLoad() (SYSCORE_Window_Wayland_Xdg_Module, error) {
	return loader.XdgModuleLoad()
}

/*
SYSCORE_Window_Wayland_Xdg_ModuleEmbeddedLibraryFileName returns the embedded soname file name.
*/
func SYSCORE_Window_Wayland_Xdg_ModuleEmbeddedLibraryFileName() string {
	return loader.XdgModuleEmbeddedLibraryFileName()
}

/*
SYSCORE_Window_Wayland_Xdg_ModuleSymbolResolve resolves wl_interface globals from the xdg module.

[Parameters]
module - Loaded xdg module.
symbolName - Symbol (for example SYSCORE_Window_Wayland_Xdg_SymbolWmBaseInterface).

[Returns]
Address of the wl_interface global.
*/
func SYSCORE_Window_Wayland_Xdg_ModuleSymbolResolve(module SYSCORE_Window_Wayland_Xdg_Module, symbolName string) (uintptr, error) {
	return loader.XdgModuleSymbolResolve(module, symbolName)
}
