//go:build windows

package loader

import "syscore/window/win32/bindings"

/*
Win32Commands holds bound Win32 function pointers from user32.dll and kernel32.dll.

[Context]
Populate with Win32CommandsLoad or Win32CommandsLoadManifest. Pair with SYSCORE_C_StringToUTF16
for class and window title strings.
*/
type Win32Commands struct {
	/*
		GetModuleHandleW is GetModuleHandleW (kernel32). Returns HINSTANCE for the process or a DLL.
	*/
	GetModuleHandleW bindings.PFN_GetModuleHandleW
	/*
		RegisterClassExW is RegisterClassExW (user32). Registers a WNDCLASSEXW before CreateWindowExW.
	*/
	RegisterClassExW bindings.PFN_RegisterClassExW
	/*
		CreateWindowExW is CreateWindowExW (user32). Creates and returns an HWND.
	*/
	CreateWindowExW bindings.PFN_CreateWindowExW
	/*
		DefWindowProcW is DefWindowProcW (user32). Default window procedure; often used as LpfnWndProc address.
	*/
	DefWindowProcW bindings.PFN_DefWindowProcW
	/*
		ShowWindow is ShowWindow (user32). Shows the window (for example SW_SHOW).
	*/
	ShowWindow bindings.PFN_ShowWindow
	/*
		UpdateWindow is UpdateWindow (user32). Forces an immediate client-area update.
	*/
	UpdateWindow bindings.PFN_UpdateWindow
	/*
		LoadCursorW is LoadCursorW (user32). Loads a cursor (for example IDC_ARROW).
	*/
	LoadCursorW bindings.PFN_LoadCursorW
	/*
		DestroyWindow is DestroyWindow (user32). Destroys an HWND created with CreateWindowExW.
	*/
	DestroyWindow bindings.PFN_DestroyWindow
	/*
		PeekMessageW is PeekMessageW (user32). Non-blocking message queue peek.
	*/
	PeekMessageW bindings.PFN_PeekMessageW
	/*
		DispatchMessageW is DispatchMessageW (user32). Dispatches a message to the window procedure.
	*/
	DispatchMessageW bindings.PFN_DispatchMessageW
}
