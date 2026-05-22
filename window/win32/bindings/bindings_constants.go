package bindings

import "unsafe"

const (
	// WS_OVERLAPPEDWINDOW is a common overlapped window style with caption and resize borders.
	WS_OVERLAPPEDWINDOW uint32 = 0x00CF0000
	// WS_EX_APPWINDOW is an extended style that shows the window on the taskbar.
	WS_EX_APPWINDOW uint32 = 0x00040000
	// CW_USEDEFAULT lets the system choose position or size for CreateWindowExW.
	CW_USEDEFAULT int32 = -2147483648
)

const (
	// CS_HREDRAW redraws the window horizontally on resize.
	CS_HREDRAW uint32 = 0x0002
	// CS_VREDRAW redraws the window vertically on resize.
	CS_VREDRAW uint32 = 0x0001
	// COLOR_WINDOW is the system color index used for the default window background brush.
	COLOR_WINDOW int32 = 5
	// IDC_ARROW is the MAKEINTRESOURCE identifier for the standard arrow cursor.
	IDC_ARROW uintptr = 32512
)

const (
	// SW_SHOW displays a window in its current size and position.
	SW_SHOW int32 = 5
)

const (
	// WM_DESTROY is sent when a window is destroyed.
	WM_DESTROY uint32 = 0x0002
)

/*
WNDPROC is a uintptr window procedure pointer type alias for documentation clarity.
*/
type WNDPROC uintptr

/*
NullPointer is a nil unsafe.Pointer for CreateWindowExW lpParam when unused.
*/
var NullPointer = unsafe.Pointer(nil)
