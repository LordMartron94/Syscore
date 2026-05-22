package bindings

import "unsafe"

/*
HWND is a Win32 window handle.
*/
type HWND uintptr

/*
HINSTANCE is a Win32 module instance handle.
*/
type HINSTANCE uintptr

/*
HICON is a Win32 icon handle.
*/
type HICON uintptr

/*
HCURSOR is a Win32 cursor handle.
*/
type HCURSOR uintptr

/*
HBRUSH is a Win32 brush handle (used for window backgrounds).
*/
type HBRUSH uintptr

/*
ATOM is a registered class atom from RegisterClassExW.
*/
type ATOM uint16

/*
WNDCLASSEXW is the Win32 WNDCLASSEXW structure for RegisterClassExW.

LpfnWndProc is typically the address of DefWindowProcW from symbol resolve, not a Go callback,
unless you install a custom WNDPROC via purego.NewCallback.
*/
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

/*
PFN_GetModuleHandleW returns a module handle for the executable or a named DLL.

[Context]
Maps to GetModuleHandleW in kernel32. Pass nil moduleName for the executable.

[Parameters]
moduleName - Wide string module name or nil.

[Returns]
Module instance handle.
*/
type PFN_GetModuleHandleW func(moduleName *uint16) HINSTANCE

/*
PFN_RegisterClassExW registers a window class for subsequent CreateWindowExW calls.

[Context]
Maps to RegisterClassExW in user32.

[Parameters]
class - Populated WNDCLASSEXW including LpszClassName and LpfnWndProc.

[Returns]
Class atom, or 0 on failure.
*/
type PFN_RegisterClassExW func(class *WNDCLASSEXW) ATOM

/*
PFN_CreateWindowExW creates a top-level or child Win32 window.

[Context]
Maps to CreateWindowExW in user32.

[Parameters]
dwExStyle - Extended style (for example WS_EX_APPWINDOW).
lpClassName - Registered class name (UTF-16).
lpWindowName - Window title (UTF-16).
dwStyle - Window style (for example WS_OVERLAPPEDWINDOW).
x, y, nWidth, nHeight - Position and size (CW_USEDEFAULT allowed).
hWndParent - Parent HWND or 0.
hMenu - Menu handle or 0.
hInstance - Module from GetModuleHandleW.
lpParam - Creation param or NullPointer.

[Returns]
Window HWND, or 0 on failure.
*/
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

/*
PFN_DefWindowProcW is the default window procedure for unhandled messages.

[Context]
Maps to DefWindowProcW. Use its resolved address as WNDCLASSEXW.LpfnWndProc for a minimal window.
*/
type PFN_DefWindowProcW func(hWnd HWND, msg uint32, wParam uintptr, lParam uintptr) uintptr

/*
PFN_ShowWindow sets the show state of a window.

[Context]
Maps to ShowWindow. Use SW_SHOW to display after CreateWindowExW.

[Parameters]
hWnd - Target window.
nCmdShow - Show command (for example SW_SHOW).

[Returns]
Previous show state.
*/
type PFN_ShowWindow func(hWnd HWND, nCmdShow int32) int32

/*
PFN_UpdateWindow repaints the client area of a window immediately.

[Context]
Maps to UpdateWindow. Optional after ShowWindow.
*/
type PFN_UpdateWindow func(hWnd HWND) int32

/*
PFN_LoadCursorW loads a stock or custom cursor.

[Context]
Maps to LoadCursorW. Pass hInstance 0 and lpCursorName IDC_ARROW for the default arrow.
*/
type PFN_LoadCursorW func(hInstance HINSTANCE, lpCursorName uintptr) HCURSOR
