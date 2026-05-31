//go:build linux

package bindings

import "unsafe"

/*
XcbInternAtomReplyAtom reads the atom field from an xcb_intern_atom_reply_t pointer.

[Parameters]
reply - Pointer from xcb_intern_atom_reply, or 0.

[Returns]
Atom ID, or 0 when reply is nil.
*/
func XcbInternAtomReplyAtom(reply uintptr) XcbAtomT {
	if reply == 0 {
		return 0
	}
	return *(*XcbAtomT)(unsafe.Pointer(reply + XcbInternAtomReplyAtomOffset))
}

/*
XcbGetGeometryReplyWidth reads width from an xcb_get_geometry_reply_t pointer.

[Parameters]
reply - Pointer from xcb_get_geometry_reply, or 0.

[Returns]
Window width in pixels, or 0 when reply is nil.
*/
func XcbGetGeometryReplyWidth(reply uintptr) uint16 {
	if reply == 0 {
		return 0
	}
	return *(*uint16)(unsafe.Pointer(reply + XcbGetGeometryReplyWidthOffset))
}

/*
XcbGetGeometryReplyHeight reads height from an xcb_get_geometry_reply_t pointer.

[Parameters]
reply - Pointer from xcb_get_geometry_reply, or 0.

[Returns]
Window height in pixels, or 0 when reply is nil.
*/
func XcbGetGeometryReplyHeight(reply uintptr) uint16 {
	if reply == 0 {
		return 0
	}
	return *(*uint16)(unsafe.Pointer(reply + XcbGetGeometryReplyHeightOffset))
}
