package syscore

import (
	"syscore/window/win32/bindings"
	"syscore/window/win32/loader"
)

/*
SYSCORE_Window_Win32_Module holds loaded user32.dll and kernel32.dll handles.

[Context]
Obtain via SYSCORE_Window_Win32_ModuleLoad on Windows before binding commands.
*/
type SYSCORE_Window_Win32_Module = loader.Win32Module

/*
SYSCORE_Window_Win32_Commands holds bound Win32 USER and KERNEL32 function pointers.

[Context]
After SYSCORE_Window_Win32_CommandsLoad, call through documented fields on loader.Win32Commands.
Use SYSCORE_C_StringToUTF16 for window and class names; see bindings constants for styles.
*/
type SYSCORE_Window_Win32_Commands = loader.Win32Commands

/*
SYSCORE_Window_Win32_CommandManifest lists Win32 entry points for selective binding.
*/
type SYSCORE_Window_Win32_CommandManifest = loader.Win32CommandManifest

const (
	SYSCORE_Window_Win32_StyleOverlappedWindow = bindings.WS_OVERLAPPEDWINDOW
	SYSCORE_Window_Win32_StyleExAppWindow      = bindings.WS_EX_APPWINDOW
	SYSCORE_Window_Win32_UseDefault            = bindings.CW_USEDEFAULT
	SYSCORE_Window_Win32_ClassStyleHRedraw     = bindings.CS_HREDRAW
	SYSCORE_Window_Win32_ClassStyleVRedraw     = bindings.CS_VREDRAW
	SYSCORE_Window_Win32_ColorWindow           = bindings.COLOR_WINDOW
	SYSCORE_Window_Win32_CursorArrow           = bindings.IDC_ARROW
	SYSCORE_Window_Win32_ShowWindow            = bindings.SW_SHOW
)

type (
	// SYSCORE_Window_Win32_HWND is a Win32 window handle.
	SYSCORE_Window_Win32_HWND = bindings.HWND
	// SYSCORE_Window_Win32_HINSTANCE is a Win32 module instance handle.
	SYSCORE_Window_Win32_HINSTANCE = bindings.HINSTANCE
	// SYSCORE_Window_Win32_HBRUSH is a Win32 brush handle.
	SYSCORE_Window_Win32_HBRUSH = bindings.HBRUSH
	// SYSCORE_Window_Win32_WNDCLASSEXW is the WNDCLASSEXW structure for RegisterClassExW.
	SYSCORE_Window_Win32_WNDCLASSEXW = bindings.WNDCLASSEXW
	// SYSCORE_Window_Win32_PFN_GetModuleHandleW is the C type for GetModuleHandleW.
	SYSCORE_Window_Win32_PFN_GetModuleHandleW = bindings.PFN_GetModuleHandleW
	// SYSCORE_Window_Win32_PFN_RegisterClassExW is the C type for RegisterClassExW.
	SYSCORE_Window_Win32_PFN_RegisterClassExW = bindings.PFN_RegisterClassExW
	// SYSCORE_Window_Win32_PFN_CreateWindowExW is the C type for CreateWindowExW.
	SYSCORE_Window_Win32_PFN_CreateWindowExW = bindings.PFN_CreateWindowExW
	// SYSCORE_Window_Win32_PFN_DefWindowProcW is the C type for DefWindowProcW.
	SYSCORE_Window_Win32_PFN_DefWindowProcW = bindings.PFN_DefWindowProcW
	// SYSCORE_Window_Win32_PFN_ShowWindow is the C type for ShowWindow.
	SYSCORE_Window_Win32_PFN_ShowWindow = bindings.PFN_ShowWindow
	// SYSCORE_Window_Win32_PFN_UpdateWindow is the C type for UpdateWindow.
	SYSCORE_Window_Win32_PFN_UpdateWindow = bindings.PFN_UpdateWindow
	// SYSCORE_Window_Win32_PFN_LoadCursorW is the C type for LoadCursorW.
	SYSCORE_Window_Win32_PFN_LoadCursorW = bindings.PFN_LoadCursorW
)

/*
SYSCORE_Window_Win32_ModuleLoad loads user32.dll and kernel32.dll.

[Context]
Entry point for native windows on Windows. Pair with SYSCORE_Window_DisplayAPIDetect (Win32).

[Returns]
Module with separate handles for each DLL.

[Errors]
Returns an error when not on windows or either LoadLibrary call fails.

[Side Effects]
Loads system DLLs into the process.
*/
func SYSCORE_Window_Win32_ModuleLoad() (SYSCORE_Window_Win32_Module, error) {
	return loader.Win32ModuleLoad()
}

/*
SYSCORE_Window_Win32_ModuleUser32LibraryName returns "user32.dll".
*/
func SYSCORE_Window_Win32_ModuleUser32LibraryName() string {
	return loader.Win32ModuleUser32LibraryName()
}

/*
SYSCORE_Window_Win32_ModuleKernel32LibraryName returns "kernel32.dll".
*/
func SYSCORE_Window_Win32_ModuleKernel32LibraryName() string {
	return loader.Win32ModuleKernel32LibraryName()
}

/*
SYSCORE_Window_Win32_ModuleUser32SymbolResolve resolves an export from user32.dll.

[Parameters]
module - Loaded Win32 module.
symbolName - Export name (for example "CreateWindowExW").

[Returns]
Function pointer address suitable for WNDCLASSEXW.LpfnWndProc when using DefWindowProcW.
*/
func SYSCORE_Window_Win32_ModuleUser32SymbolResolve(module SYSCORE_Window_Win32_Module, symbolName string) (uintptr, error) {
	return loader.Win32ModuleUser32SymbolResolve(module, symbolName)
}

/*
SYSCORE_Window_Win32_ModuleKernel32SymbolResolve resolves an export from kernel32.dll.
*/
func SYSCORE_Window_Win32_ModuleKernel32SymbolResolve(module SYSCORE_Window_Win32_Module, symbolName string) (uintptr, error) {
	return loader.Win32ModuleKernel32SymbolResolve(module, symbolName)
}

/*
SYSCORE_Window_Win32_CommandsLoad binds all supported Win32 window entry points.

[Parameters]
module - Loaded module.
commands - Non-nil command holder.

[Errors]
Returns an error if commands is nil, not on windows, or any bind fails.
*/
func SYSCORE_Window_Win32_CommandsLoad(module SYSCORE_Window_Win32_Module, commands *SYSCORE_Window_Win32_Commands) error {
	return loader.Win32CommandsLoad(module, commands)
}

/*
SYSCORE_Window_Win32_CommandManifestReset clears manifest.
*/
func SYSCORE_Window_Win32_CommandManifestReset(manifest *SYSCORE_Window_Win32_CommandManifest) {
	loader.Win32CommandManifestReset(manifest)
}

/*
SYSCORE_Window_Win32_CommandManifestAddGetModuleHandleW registers GetModuleHandleW (kernel32) for selective loading.
*/
func SYSCORE_Window_Win32_CommandManifestAddGetModuleHandleW(manifest *SYSCORE_Window_Win32_CommandManifest, target *SYSCORE_Window_Win32_PFN_GetModuleHandleW) error {
	return loader.Win32CommandManifestAddGetModuleHandleW(manifest, target)
}

/*
SYSCORE_Window_Win32_CommandManifestAddRegisterClassExW registers RegisterClassExW for selective loading.
*/
func SYSCORE_Window_Win32_CommandManifestAddRegisterClassExW(manifest *SYSCORE_Window_Win32_CommandManifest, target *SYSCORE_Window_Win32_PFN_RegisterClassExW) error {
	return loader.Win32CommandManifestAddRegisterClassExW(manifest, target)
}

/*
SYSCORE_Window_Win32_CommandManifestAddCreateWindowExW registers CreateWindowExW for selective loading.
*/
func SYSCORE_Window_Win32_CommandManifestAddCreateWindowExW(manifest *SYSCORE_Window_Win32_CommandManifest, target *SYSCORE_Window_Win32_PFN_CreateWindowExW) error {
	return loader.Win32CommandManifestAddCreateWindowExW(manifest, target)
}

/*
SYSCORE_Window_Win32_CommandManifestAddDefWindowProcW registers DefWindowProcW for selective loading.
*/
func SYSCORE_Window_Win32_CommandManifestAddDefWindowProcW(manifest *SYSCORE_Window_Win32_CommandManifest, target *SYSCORE_Window_Win32_PFN_DefWindowProcW) error {
	return loader.Win32CommandManifestAddDefWindowProcW(manifest, target)
}

/*
SYSCORE_Window_Win32_CommandManifestAddShowWindow registers ShowWindow for selective loading.
*/
func SYSCORE_Window_Win32_CommandManifestAddShowWindow(manifest *SYSCORE_Window_Win32_CommandManifest, target *SYSCORE_Window_Win32_PFN_ShowWindow) error {
	return loader.Win32CommandManifestAddShowWindow(manifest, target)
}

/*
SYSCORE_Window_Win32_CommandManifestAddUpdateWindow registers UpdateWindow for selective loading.
*/
func SYSCORE_Window_Win32_CommandManifestAddUpdateWindow(manifest *SYSCORE_Window_Win32_CommandManifest, target *SYSCORE_Window_Win32_PFN_UpdateWindow) error {
	return loader.Win32CommandManifestAddUpdateWindow(manifest, target)
}

/*
SYSCORE_Window_Win32_CommandManifestAddLoadCursorW registers LoadCursorW for selective loading.
*/
func SYSCORE_Window_Win32_CommandManifestAddLoadCursorW(manifest *SYSCORE_Window_Win32_CommandManifest, target *SYSCORE_Window_Win32_PFN_LoadCursorW) error {
	return loader.Win32CommandManifestAddLoadCursorW(manifest, target)
}

/*
SYSCORE_Window_Win32_CommandsLoadManifest binds only manifest-listed Win32 entry points.

[Errors]
Returns an error if manifest or commands is nil or any bind fails.
*/
func SYSCORE_Window_Win32_CommandsLoadManifest(module SYSCORE_Window_Win32_Module, manifest *SYSCORE_Window_Win32_CommandManifest, commands *SYSCORE_Window_Win32_Commands) error {
	return loader.Win32CommandsLoadManifest(module, manifest, commands)
}
