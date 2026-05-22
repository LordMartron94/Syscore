package bindings

/*
XCB_WINDOW_CLASS_INPUT_OUTPUT is the InputOutput window class (copy from parent visual/depth when combined with zero depth/visual).
*/
const XCB_WINDOW_CLASS_INPUT_OUTPUT uint16 = 1

/*
XCB_PROP_MODE_REPLACE replaces the previous property value (XCB_PROP_MODE_REPLACE).
*/
const XCB_PROP_MODE_REPLACE uint8 = 0

/*
XCB_ATOM_WM_NAME is the legacy WM_NAME property atom name (intern before use).
*/
const XCB_ATOM_WM_NAME = "WM_NAME"

/*
XCB_ATOM_NET_WM_NAME is the EWMH _NET_WM_NAME property atom name (intern before use).
*/
const XCB_ATOM_NET_WM_NAME = "_NET_WM_NAME"

/*
XCB_ATOM_UTF8_STRING is the UTF8_STRING type atom name (intern before use).
*/
const XCB_ATOM_UTF8_STRING = "UTF8_STRING"

/*
XcbInternAtomReplyAtomOffset is the byte offset of the atom field in xcb_intern_atom_reply_t.
*/
const XcbInternAtomReplyAtomOffset = 8
