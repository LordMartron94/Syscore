package bindings

import "unsafe"

/*
XcbConnectionT is an opaque xcb_connection_t* returned by xcb_connect.
*/
type XcbConnectionT uintptr

/*
XcbWindowT is an X11 window ID (xcb_window_t).
*/
type XcbWindowT uint32

/*
XcbVoidCookieT is a void cookie (xcb_void_cookie_t) for unchecked requests.
*/
type XcbVoidCookieT uint32

/*
XcbScreenIteratorT is the xcb_screen_iterator_t struct from xcb_setup_roots_iterator.

Data points at the current xcb_screen_t; Rem is the number of screens remaining.
*/
type XcbScreenIteratorT struct {
	Data  uintptr
	Rem   int32
	Index int32
}

/*
PFN_xcb_connect opens a connection to the X window server.

[Context]
Maps to xcb_connect. Pass nil displayname to use DISPLAY from the environment.

[Parameters]
displayname - Connection name (DISPLAY) or nil.
screenp - Optional output screen number, or nil.

[Returns]
Connection handle, or 0 on failure.

[Reference]
https://xcb.freedesktop.org/manual/xcb-connect.html
*/
type PFN_xcb_connect func(displayname *byte, screenp *int32) XcbConnectionT

/*
PFN_xcb_get_setup returns the setup data from a connection.

[Context]
Maps to xcb_get_setup. The returned uintptr is *xcb_setup_t.

[Parameters]
c - Connection from xcb_connect.

[Returns]
Pointer to xcb_setup_t as uintptr.
*/
type PFN_xcb_get_setup func(c XcbConnectionT) uintptr

/*
PFN_xcb_setup_roots_iterator returns an iterator over screens in the setup.

[Context]
Maps to xcb_setup_roots_iterator. Use Data to read the root window ID from the first screen.

[Parameters]
setup - Pointer from xcb_get_setup.

[Returns]
Screen iterator struct by value.
*/
type PFN_xcb_setup_roots_iterator func(setup uintptr) XcbScreenIteratorT

/*
PFN_xcb_generate_id allocates a new X resource ID on a connection.

[Context]
Maps to xcb_generate_id. Required before xcb_create_window for the new window ID.

[Parameters]
c - Connection handle.

[Returns]
New resource ID.
*/
type PFN_xcb_generate_id func(c XcbConnectionT) uint32

/*
PFN_xcb_create_window creates an X window as a child of parent.

[Context]
Maps to xcb_create_window. This binding uses the unchecked variant (returns a cookie only).

[Parameters]
c - Connection handle.
depth - CopyFromParent when 0.
wid - Window ID from xcb_generate_id.
parent - Parent window (typically root).
x, y - Position.
width, height - Size in pixels.
borderWidth - Border width (often 0).
class - Window class (for example XCB_WINDOW_CLASS_INPUT_OUTPUT).
visual - CopyFromParent when 0.
valueMask - Attribute mask (0 for defaults).
valueList - Attribute values pointer or nil.

[Returns]
Request cookie.
*/
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

/*
PFN_xcb_map_window maps a window so it becomes visible when the server paints.

[Context]
Maps to xcb_map_window (unchecked).

[Parameters]
c - Connection handle.
window - Window ID to map.

[Returns]
Request cookie.
*/
type PFN_xcb_map_window func(c XcbConnectionT, window XcbWindowT) XcbVoidCookieT

/*
PFN_xcb_flush flushes the output buffer to the X server.

[Context]
Maps to xcb_flush. Call after a batch of requests so they reach the server.

[Parameters]
c - Connection handle.

[Returns]
Number of bytes flushed, or a negative value on error.
*/
type PFN_xcb_flush func(c XcbConnectionT) int32
