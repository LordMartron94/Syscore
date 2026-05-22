package bindings

import "unsafe"

const (
	WS_OVERLAPPEDWINDOW uint32 = 0x00CF0000
	WS_EX_APPWINDOW     uint32 = 0x00040000
	CW_USEDEFAULT       int32  = -2147483648 // 0x80000000 as signed int32
)

const (
	CS_HREDRAW   uint32  = 0x0002
	CS_VREDRAW   uint32  = 0x0001
	COLOR_WINDOW int32   = 5
	IDC_ARROW    uintptr = 32512
)

const (
	SW_SHOW int32 = 5
)

const (
	WM_DESTROY uint32 = 0x0002
)

type WNDPROC uintptr

var NullPointer = unsafe.Pointer(nil)
