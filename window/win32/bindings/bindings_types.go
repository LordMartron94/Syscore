package bindings

import "unsafe"

type HWND uintptr
type HINSTANCE uintptr
type HICON uintptr
type HCURSOR uintptr
type HBRUSH uintptr
type ATOM uint16

type WNDCLASSEXW struct {
	CbSize        uint32
	Style         uint32
	LpfnWndProc   uintptr
	CbClsExtra    int32
	CbWndExtra    int32
	HInstance     HINSTANCE
	HIcon         HICON
	HCursor       HCURSOR
	HbrBackground HBRUSH
	LpszMenuName  *uint16
	LpszClassName *uint16
	HIconSm       HICON
}

type PFN_GetModuleHandleW func(moduleName *uint16) HINSTANCE

type PFN_RegisterClassExW func(class *WNDCLASSEXW) ATOM

type PFN_CreateWindowExW func(
	dwExStyle uint32,
	lpClassName *uint16,
	lpWindowName *uint16,
	dwStyle uint32,
	x int32,
	y int32,
	nWidth int32,
	nHeight int32,
	hWndParent HWND,
	hMenu uintptr,
	hInstance HINSTANCE,
	lpParam unsafe.Pointer,
) HWND

type PFN_DefWindowProcW func(hWnd HWND, msg uint32, wParam uintptr, lParam uintptr) uintptr

type PFN_ShowWindow func(hWnd HWND, nCmdShow int32) int32

type PFN_UpdateWindow func(hWnd HWND) int32

type PFN_LoadCursorW func(hInstance HINSTANCE, lpCursorName uintptr) HCURSOR
