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
XCB_ATOM_WM_NORMAL_HINTS is the WM_NORMAL_HINTS property atom name (intern before use).
*/
const XCB_ATOM_WM_NORMAL_HINTS = "WM_NORMAL_HINTS"

/*
XCB_ATOM_WM_SIZE_HINTS is the WM_SIZE_HINTS type atom name for WM_NORMAL_HINTS values.
*/
const XCB_ATOM_WM_SIZE_HINTS = "WM_SIZE_HINTS"

const (
	// XCB_SIZE_HINTS_FLAG_PSIZE sets width and height in XSizeHints.
	XCB_SIZE_HINTS_FLAG_PSIZE int64 = 1 << 2
	// XCB_SIZE_HINTS_FLAG_PMIN_SIZE sets min_width and min_height in XSizeHints.
	XCB_SIZE_HINTS_FLAG_PMIN_SIZE int64 = 1 << 4
	// XCB_SIZE_HINTS_FLAG_PMAX_SIZE sets max_width and max_height in XSizeHints.
	XCB_SIZE_HINTS_FLAG_PMAX_SIZE int64 = 1 << 5
)

/*
XcbInternAtomReplyAtomOffset is the byte offset of the atom field in xcb_intern_atom_reply_t.
*/
const XcbInternAtomReplyAtomOffset = 8
