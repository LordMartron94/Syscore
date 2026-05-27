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
