package loader

import "syscore/window/win32/bindings"

type Win32Commands struct {
	GetModuleHandleW bindings.PFN_GetModuleHandleW
	RegisterClassExW bindings.PFN_RegisterClassExW
	CreateWindowExW  bindings.PFN_CreateWindowExW
	DefWindowProcW   bindings.PFN_DefWindowProcW
	ShowWindow       bindings.PFN_ShowWindow
	UpdateWindow     bindings.PFN_UpdateWindow
	LoadCursorW      bindings.PFN_LoadCursorW
}
