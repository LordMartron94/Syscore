//go:build windows

package syscore

import (
	"unsafe"

	kernel32 "syscore/window/win32/bindings/pfn/kernel32_dll"
	user32 "syscore/window/win32/bindings/pfn/user32_dll"
	foundation "syscore/window/win32/bindings/types/windows_win32_foundation"
	gdi "syscore/window/win32/bindings/types/windows_win32_graphics_gdi"
	wm "syscore/window/win32/bindings/types/windows_win32_ui_windowsandmessaging"
	"syscore/window/win32/loader"
)

/*
SYSCORE_Window_Win32_CursorArrow returns the *PWSTRElement form of the IDC_ARROW pseudo-resource.

[Context]
Win32 uses MAKEINTRESOURCE(32512) for IDC_ARROW: the integer is stored in the low bits of a
PWSTR-typed pointer. LoadCursorW accepts this as the lpCursorName argument. The returned pointer
must not be dereferenced - it is a tagged scalar, not a real string.

[Returns]
A *PWSTRElement value carrying 32512 in its address bits.
*/
func SYSCORE_Window_Win32_CursorArrow() *SYSCORE_Window_Win32_PWSTRElement {
	return (*SYSCORE_Window_Win32_PWSTRElement)(unsafe.Pointer(SYSCORE_Window_Win32_CursorArrowResource))
}

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
	SYSCORE_Window_Win32_StyleOverlappedWindow = wm.WS_OVERLAPPEDWINDOW
	// SYSCORE_Window_Win32_StyleOverlappedWindowFixed is WS_OVERLAPPED|WS_CAPTION|WS_SYSMENU|WS_MINIMIZEBOX (0x00CA0000); not a single winmd enum member.
	SYSCORE_Window_Win32_StyleOverlappedWindowFixed = wm.WS_OVERLAPPED | wm.WS_CAPTION | wm.WS_SYSMENU | wm.WS_MINIMIZEBOX
	SYSCORE_Window_Win32_StyleExAppWindow           = wm.WS_EX_APPWINDOW
	SYSCORE_Window_Win32_UseDefault                 = wm.CW_USEDEFAULT
	SYSCORE_Window_Win32_ClassStyleHRedraw          = wm.CS_HREDRAW
	SYSCORE_Window_Win32_ClassStyleVRedraw          = wm.CS_VREDRAW
	// SYSCORE_Window_Win32_ColorWindow exposes COLOR_WINDOW from the Graphics.Gdi SYS_COLOR_INDEX
	// enum. The system color index lives in graphics_gdi rather than ui_windowsandmessaging in
	// the win32 metadata.
	SYSCORE_Window_Win32_ColorWindow = gdi.COLOR_WINDOW
	// SYSCORE_Window_Win32_CursorArrowResource is the IDC_ARROW raw resource id; LoadCursorW
	// expects this value packed into a *PWSTRElement via MAKEINTRESOURCE semantics. Use
	// SYSCORE_Window_Win32_CursorArrow to obtain that pointer.
	SYSCORE_Window_Win32_CursorArrowResource uintptr = 32512
	SYSCORE_Window_Win32_ShowWindow                  = wm.SW_SHOW
	SYSCORE_Window_Win32_ShowWindowHide              = wm.SW_HIDE
	SYSCORE_Window_Win32_MessageClose                = wm.WM_CLOSE
	SYSCORE_Window_Win32_MessageDestroy              = wm.WM_DESTROY
	SYSCORE_Window_Win32_PeekMessageRemove           = wm.PM_REMOVE
)

type (
	// SYSCORE_Window_Win32_HWND is a Win32 window handle.
	SYSCORE_Window_Win32_HWND = foundation.HWND
	// SYSCORE_Window_Win32_HINSTANCE is a Win32 module instance handle.
	SYSCORE_Window_Win32_HINSTANCE = foundation.HINSTANCE
	// SYSCORE_Window_Win32_HMODULE is a Win32 loaded module handle (returned by GetModuleHandleW).
	SYSCORE_Window_Win32_HMODULE = foundation.HMODULE
	// SYSCORE_Window_Win32_HBRUSH is a Win32 brush handle.
	SYSCORE_Window_Win32_HBRUSH = gdi.HBRUSH
	// SYSCORE_Window_Win32_HCURSOR is a Win32 cursor handle.
	SYSCORE_Window_Win32_HCURSOR = wm.HCURSOR
	// SYSCORE_Window_Win32_WPARAM is the wParam scalar passed to a window procedure.
	SYSCORE_Window_Win32_WPARAM = foundation.WPARAM
	// SYSCORE_Window_Win32_LPARAM is the lParam scalar passed to a window procedure.
	SYSCORE_Window_Win32_LPARAM = foundation.LPARAM
	// SYSCORE_Window_Win32_LRESULT is the LRESULT returned by a window procedure.
	SYSCORE_Window_Win32_LRESULT = foundation.LRESULT
	// SYSCORE_Window_Win32_WNDPROC is the window-procedure function pointer type used in WNDCLASSEXW.LpfnWndProc.
	SYSCORE_Window_Win32_WNDPROC = wm.WNDPROC
	// SYSCORE_Window_Win32_PWSTRElement is the UTF-16 code unit type used to back PWSTR/LPCWSTR pointers.
	SYSCORE_Window_Win32_PWSTRElement = foundation.PWSTRElement
	// SYSCORE_Window_Win32_WindowStyle is the WINDOW_STYLE enum (WS_*) accepted by CreateWindowExW.
	SYSCORE_Window_Win32_WindowStyle = wm.WINDOW_STYLE
	// SYSCORE_Window_Win32_WindowExStyle is the WINDOW_EX_STYLE enum (WS_EX_*) accepted by CreateWindowExW.
	SYSCORE_Window_Win32_WindowExStyle = wm.WINDOW_EX_STYLE
	// SYSCORE_Window_Win32_ShowWindowCmd is the SHOW_WINDOW_CMD enum (SW_*) accepted by ShowWindow.
	SYSCORE_Window_Win32_ShowWindowCmd = wm.SHOW_WINDOW_CMD
	// SYSCORE_Window_Win32_PeekMessageRemoveType is the PEEK_MESSAGE_REMOVE_TYPE enum (PM_*) accepted by PeekMessageW.
	SYSCORE_Window_Win32_PeekMessageRemoveType = wm.PEEK_MESSAGE_REMOVE_TYPE
	// SYSCORE_Window_Win32_WNDCLASSEXW is the WNDCLASSEXW structure for RegisterClassExW.
	SYSCORE_Window_Win32_WNDCLASSEXW = wm.WNDCLASSEXW
	// SYSCORE_Window_Win32_PFN_GetModuleHandleW is the C type for GetModuleHandleW.
	SYSCORE_Window_Win32_PFN_GetModuleHandleW = kernel32.PFN_GetModuleHandleW
	// SYSCORE_Window_Win32_PFN_RegisterClassExW is the C type for RegisterClassExW.
	SYSCORE_Window_Win32_PFN_RegisterClassExW = user32.PFN_RegisterClassExW
	// SYSCORE_Window_Win32_PFN_CreateWindowExW is the C type for CreateWindowExW.
	SYSCORE_Window_Win32_PFN_CreateWindowExW = user32.PFN_CreateWindowExW
	// SYSCORE_Window_Win32_PFN_DefWindowProcW is the C type for DefWindowProcW.
	SYSCORE_Window_Win32_PFN_DefWindowProcW = user32.PFN_DefWindowProcW
	// SYSCORE_Window_Win32_PFN_ShowWindow is the C type for ShowWindow.
	SYSCORE_Window_Win32_PFN_ShowWindow = user32.PFN_ShowWindow
	// SYSCORE_Window_Win32_PFN_UpdateWindow is the C type for UpdateWindow.
	SYSCORE_Window_Win32_PFN_UpdateWindow = user32.PFN_UpdateWindow
	// SYSCORE_Window_Win32_PFN_LoadCursorW is the C type for LoadCursorW.
	SYSCORE_Window_Win32_PFN_LoadCursorW = user32.PFN_LoadCursorW
	// SYSCORE_Window_Win32_PFN_DestroyWindow is the C type for DestroyWindow.
	SYSCORE_Window_Win32_PFN_DestroyWindow = user32.PFN_DestroyWindow
	// SYSCORE_Window_Win32_PFN_PeekMessageW is the C type for PeekMessageW.
	SYSCORE_Window_Win32_PFN_PeekMessageW = user32.PFN_PeekMessageW
	// SYSCORE_Window_Win32_PFN_DispatchMessageW is the C type for DispatchMessageW.
	SYSCORE_Window_Win32_PFN_DispatchMessageW = user32.PFN_DispatchMessageW
	// SYSCORE_Window_Win32_MSG is the Win32 MSG structure.
	SYSCORE_Window_Win32_MSG = wm.MSG
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
SYSCORE_Window_Win32_CommandManifestAddDestroyWindow registers DestroyWindow for selective loading.
*/
func SYSCORE_Window_Win32_CommandManifestAddDestroyWindow(manifest *SYSCORE_Window_Win32_CommandManifest, target *SYSCORE_Window_Win32_PFN_DestroyWindow) error {
	return loader.Win32CommandManifestAddDestroyWindow(manifest, target)
}

/*
SYSCORE_Window_Win32_CommandsLoadManifest binds only manifest-listed Win32 entry points.

[Errors]
Returns an error if manifest or commands is nil or any bind fails.
*/
func SYSCORE_Window_Win32_CommandsLoadManifest(module SYSCORE_Window_Win32_Module, manifest *SYSCORE_Window_Win32_CommandManifest, commands *SYSCORE_Window_Win32_Commands) error {
	return loader.Win32CommandsLoadManifest(module, manifest, commands)
}
