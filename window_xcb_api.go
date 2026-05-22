//go:build linux

package syscore

import (
	"syscore/window/xcb/bindings"
	"syscore/window/xcb/loader"
)

/*
SYSCORE_Window_Xcb_Module holds a loaded libxcb.so.1 handle for X11/XCB windowing.

[Context]
Obtain via SYSCORE_Window_Xcb_ModuleLoad before binding commands. Linux only; other GOOS
return an error from load helpers.
*/
type SYSCORE_Window_Xcb_Module = loader.XcbModule

/*
SYSCORE_Window_Xcb_Commands holds bound XCB function pointers (1:1 with xcb_* exports).

[Context]
After SYSCORE_Window_Xcb_CommandsLoad, call through fields such as Connect, CreateWindow, and
Flush. Each field has hover documentation describing the corresponding xcb_* API (see
syscore/window/xcb/loader.XcbCommands and syscore/window/xcb/bindings.PFN_*).
*/
type SYSCORE_Window_Xcb_Commands = loader.XcbCommands

/*
SYSCORE_Window_Xcb_CommandManifest lists XCB entry points to bind selectively at runtime.

[Context]
Build with CommandManifestAdd* helpers, then pass to SYSCORE_Window_Xcb_CommandsLoadManifest.
*/
type SYSCORE_Window_Xcb_CommandManifest = loader.XcbCommandManifest

const (
	SYSCORE_Window_Xcb_WindowClassInputOutput = bindings.XCB_WINDOW_CLASS_INPUT_OUTPUT
	SYSCORE_Window_Xcb_PropModeReplace        = bindings.XCB_PROP_MODE_REPLACE
	SYSCORE_Window_Xcb_AtomNameWMName         = bindings.XCB_ATOM_WM_NAME
	SYSCORE_Window_Xcb_AtomNameNetWMName      = bindings.XCB_ATOM_NET_WM_NAME
	SYSCORE_Window_Xcb_AtomNameUTF8String     = bindings.XCB_ATOM_UTF8_STRING
)

type (
	// SYSCORE_Window_Xcb_ConnectionT is an opaque xcb_connection_t* (uintptr).
	SYSCORE_Window_Xcb_ConnectionT = bindings.XcbConnectionT
	// SYSCORE_Window_Xcb_WindowT is an xcb_window_t window ID.
	SYSCORE_Window_Xcb_WindowT = bindings.XcbWindowT
	// SYSCORE_Window_Xcb_VoidCookieT is an xcb_void_cookie_t request cookie.
	SYSCORE_Window_Xcb_VoidCookieT = bindings.XcbVoidCookieT
	// SYSCORE_Window_Xcb_ScreenIteratorT is the xcb_screen_iterator_t struct returned by xcb_setup_roots_iterator.
	SYSCORE_Window_Xcb_ScreenIteratorT = bindings.XcbScreenIteratorT
	// SYSCORE_Window_Xcb_PFN_connect is the C type for xcb_connect.
	SYSCORE_Window_Xcb_PFN_connect = bindings.PFN_xcb_connect
	// SYSCORE_Window_Xcb_PFN_get_setup is the C type for xcb_get_setup.
	SYSCORE_Window_Xcb_PFN_get_setup = bindings.PFN_xcb_get_setup
	// SYSCORE_Window_Xcb_PFN_setup_roots_iter is the C type for xcb_setup_roots_iterator.
	SYSCORE_Window_Xcb_PFN_setup_roots_iter = bindings.PFN_xcb_setup_roots_iterator
	// SYSCORE_Window_Xcb_PFN_generate_id is the C type for xcb_generate_id.
	SYSCORE_Window_Xcb_PFN_generate_id = bindings.PFN_xcb_generate_id
	// SYSCORE_Window_Xcb_PFN_create_window is the C type for xcb_create_window.
	SYSCORE_Window_Xcb_PFN_create_window = bindings.PFN_xcb_create_window
	// SYSCORE_Window_Xcb_PFN_map_window is the C type for xcb_map_window.
	SYSCORE_Window_Xcb_PFN_map_window = bindings.PFN_xcb_map_window
	// SYSCORE_Window_Xcb_PFN_flush is the C type for xcb_flush.
	SYSCORE_Window_Xcb_PFN_flush = bindings.PFN_xcb_flush
	// SYSCORE_Window_Xcb_AtomT is an xcb_atom_t property or type atom ID.
	SYSCORE_Window_Xcb_AtomT = bindings.XcbAtomT
	// SYSCORE_Window_Xcb_InternAtomCookieT is the cookie from xcb_intern_atom.
	SYSCORE_Window_Xcb_InternAtomCookieT = bindings.XcbInternAtomCookieT
	// SYSCORE_Window_Xcb_PFN_intern_atom is the C type for xcb_intern_atom.
	SYSCORE_Window_Xcb_PFN_intern_atom = bindings.PFN_xcb_intern_atom
	// SYSCORE_Window_Xcb_PFN_intern_atom_reply is the C type for xcb_intern_atom_reply.
	SYSCORE_Window_Xcb_PFN_intern_atom_reply = bindings.PFN_xcb_intern_atom_reply
	// SYSCORE_Window_Xcb_PFN_change_property is the C type for xcb_change_property.
	SYSCORE_Window_Xcb_PFN_change_property = bindings.PFN_xcb_change_property
	// SYSCORE_Window_Xcb_PFN_free is the C type for libc free(3) on XCB reply buffers.
	SYSCORE_Window_Xcb_PFN_free = bindings.PFN_c_free
)

/*
SYSCORE_Window_Xcb_ModuleLoad loads libxcb.so.1.

[Context]
First step for X11 windowing on Linux. Pair with SYSCORE_Window_DisplayAPIDetect when choosing
between Wayland and XCB.

[Returns]
Module handle on success.

[Errors]
Returns an error when not on linux or when dlopen fails.

[Side Effects]
Loads libxcb into the process.
*/
func SYSCORE_Window_Xcb_ModuleLoad() (SYSCORE_Window_Xcb_Module, error) {
	return loader.XcbModuleLoad()
}

/*
SYSCORE_Window_Xcb_ModuleLibraryName returns the soname used for XcbModuleLoad ("libxcb.so.1").
*/
func SYSCORE_Window_Xcb_ModuleLibraryName() string {
	return loader.XcbModuleLibraryName()
}

/*
SYSCORE_Window_Xcb_ModuleSymbolResolve resolves an exported symbol from the loaded XCB library.

[Parameters]
module - Loaded module from SYSCORE_Window_Xcb_ModuleLoad.
symbolName - Exported symbol (for example "xcb_connect").

[Returns]
Symbol address.

[Errors]
Returns an error if the symbol is missing.
*/
func SYSCORE_Window_Xcb_ModuleSymbolResolve(module SYSCORE_Window_Xcb_Module, symbolName string) (uintptr, error) {
	return loader.XcbModuleSymbolResolve(module, symbolName)
}

/*
SYSCORE_Window_Xcb_CommandsLoad binds every supported xcb_* entry point in the module.

[Parameters]
module - Loaded XCB module.
commands - Non-nil command holder to populate.

[Returns]
nil on success.

[Errors]
Returns an error if commands is nil, not on linux, or any symbol bind fails.
*/
func SYSCORE_Window_Xcb_CommandsLoad(module SYSCORE_Window_Xcb_Module, commands *SYSCORE_Window_Xcb_Commands) error {
	return loader.XcbCommandsLoad(module, commands)
}

/*
SYSCORE_Window_Xcb_CommandManifestReset clears all entries from manifest.

[Parameters]
manifest - Manifest to reset; nil is a no-op.
*/
func SYSCORE_Window_Xcb_CommandManifestReset(manifest *SYSCORE_Window_Xcb_CommandManifest) {
	loader.XcbCommandManifestReset(manifest)
}

/*
SYSCORE_Window_Xcb_CommandManifestAddConnect registers xcb_connect for selective loading.

[Parameters]
manifest - Manifest under construction.
target - Non-nil *SYSCORE_Window_Xcb_PFN_connect (typically &commands.Connect).

[Errors]
Returns an error if manifest or target is nil.
*/
func SYSCORE_Window_Xcb_CommandManifestAddConnect(manifest *SYSCORE_Window_Xcb_CommandManifest, target *SYSCORE_Window_Xcb_PFN_connect) error {
	return loader.XcbCommandManifestAddConnect(manifest, target)
}

/*
SYSCORE_Window_Xcb_CommandManifestAddGetSetup registers xcb_get_setup for selective loading.
*/
func SYSCORE_Window_Xcb_CommandManifestAddGetSetup(manifest *SYSCORE_Window_Xcb_CommandManifest, target *SYSCORE_Window_Xcb_PFN_get_setup) error {
	return loader.XcbCommandManifestAddGetSetup(manifest, target)
}

/*
SYSCORE_Window_Xcb_CommandManifestAddSetupRootsIterator registers xcb_setup_roots_iterator for selective loading.
*/
func SYSCORE_Window_Xcb_CommandManifestAddSetupRootsIterator(manifest *SYSCORE_Window_Xcb_CommandManifest, target *SYSCORE_Window_Xcb_PFN_setup_roots_iter) error {
	return loader.XcbCommandManifestAddSetupRootsIterator(manifest, target)
}

/*
SYSCORE_Window_Xcb_CommandManifestAddGenerateId registers xcb_generate_id for selective loading.
*/
func SYSCORE_Window_Xcb_CommandManifestAddGenerateId(manifest *SYSCORE_Window_Xcb_CommandManifest, target *SYSCORE_Window_Xcb_PFN_generate_id) error {
	return loader.XcbCommandManifestAddGenerateId(manifest, target)
}

/*
SYSCORE_Window_Xcb_CommandManifestAddCreateWindow registers xcb_create_window for selective loading.
*/
func SYSCORE_Window_Xcb_CommandManifestAddCreateWindow(manifest *SYSCORE_Window_Xcb_CommandManifest, target *SYSCORE_Window_Xcb_PFN_create_window) error {
	return loader.XcbCommandManifestAddCreateWindow(manifest, target)
}

/*
SYSCORE_Window_Xcb_CommandManifestAddMapWindow registers xcb_map_window for selective loading.
*/
func SYSCORE_Window_Xcb_CommandManifestAddMapWindow(manifest *SYSCORE_Window_Xcb_CommandManifest, target *SYSCORE_Window_Xcb_PFN_map_window) error {
	return loader.XcbCommandManifestAddMapWindow(manifest, target)
}

/*
SYSCORE_Window_Xcb_CommandManifestAddFlush registers xcb_flush for selective loading.
*/
func SYSCORE_Window_Xcb_CommandManifestAddFlush(manifest *SYSCORE_Window_Xcb_CommandManifest, target *SYSCORE_Window_Xcb_PFN_flush) error {
	return loader.XcbCommandManifestAddFlush(manifest, target)
}

/*
SYSCORE_Window_Xcb_CommandsLoadManifest binds only the entry points registered on manifest.

[Parameters]
module - Loaded XCB module.
manifest - Manifest built with CommandManifestAdd* calls.
commands - Command holder receiving bindings.

[Errors]
Returns an error if manifest or commands is nil, not on linux, or any listed symbol bind fails.
*/
func SYSCORE_Window_Xcb_CommandsLoadManifest(module SYSCORE_Window_Xcb_Module, manifest *SYSCORE_Window_Xcb_CommandManifest, commands *SYSCORE_Window_Xcb_Commands) error {
	return loader.XcbCommandsLoadManifest(module, manifest, commands)
}

/*
SYSCORE_Window_Xcb_InternAtomReplyAtom reads the atom field from an xcb_intern_atom_reply_t pointer.

[Parameters]
reply - Pointer returned by commands.InternAtomReply.

[Returns]
Atom ID, or 0 when reply is nil.
*/
func SYSCORE_Window_Xcb_InternAtomReplyAtom(reply uintptr) SYSCORE_Window_Xcb_AtomT {
	return bindings.XcbInternAtomReplyAtom(reply)
}

/*
SYSCORE_Window_Xcb_CommandManifestAddInternAtom registers xcb_intern_atom for selective loading.
*/
func SYSCORE_Window_Xcb_CommandManifestAddInternAtom(manifest *SYSCORE_Window_Xcb_CommandManifest, target *SYSCORE_Window_Xcb_PFN_intern_atom) error {
	return loader.XcbCommandManifestAddInternAtom(manifest, target)
}

/*
SYSCORE_Window_Xcb_CommandManifestAddInternAtomReply registers xcb_intern_atom_reply for selective loading.
*/
func SYSCORE_Window_Xcb_CommandManifestAddInternAtomReply(manifest *SYSCORE_Window_Xcb_CommandManifest, target *SYSCORE_Window_Xcb_PFN_intern_atom_reply) error {
	return loader.XcbCommandManifestAddInternAtomReply(manifest, target)
}

/*
SYSCORE_Window_Xcb_CommandManifestAddChangeProperty registers xcb_change_property for selective loading.
*/
func SYSCORE_Window_Xcb_CommandManifestAddChangeProperty(manifest *SYSCORE_Window_Xcb_CommandManifest, target *SYSCORE_Window_Xcb_PFN_change_property) error {
	return loader.XcbCommandManifestAddChangeProperty(manifest, target)
}

/*
SYSCORE_Window_Xcb_CommandManifestAddFree registers libc free for selective loading (XCB reply release).
*/
func SYSCORE_Window_Xcb_CommandManifestAddFree(manifest *SYSCORE_Window_Xcb_CommandManifest, target *SYSCORE_Window_Xcb_PFN_free) error {
	return loader.XcbCommandManifestAddFree(manifest, target)
}
