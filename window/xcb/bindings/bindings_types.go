package bindings

import "unsafe"

type XcbConnectionT uintptr
type XcbWindowT uint32
type XcbVoidCookieT uint32

type XcbScreenIteratorT struct {
	Data  uintptr
	Rem   int32
	Index int32
}

type PFN_xcb_connect func(displayname *byte, screenp *int32) XcbConnectionT

type PFN_xcb_get_setup func(c XcbConnectionT) uintptr

type PFN_xcb_setup_roots_iterator func(setup uintptr) XcbScreenIteratorT

type PFN_xcb_generate_id func(c XcbConnectionT) uint32

type PFN_xcb_create_window func(
	c XcbConnectionT,
	depth uint8,
	wid XcbWindowT,
	parent XcbWindowT,
	x int16,
	y int16,
	width uint16,
	height uint16,
	borderWidth uint16,
	class uint16,
	visual uint32,
	valueMask uint32,
	valueList unsafe.Pointer,
) XcbVoidCookieT

type PFN_xcb_map_window func(c XcbConnectionT, window XcbWindowT) XcbVoidCookieT

type PFN_xcb_flush func(c XcbConnectionT) int32
