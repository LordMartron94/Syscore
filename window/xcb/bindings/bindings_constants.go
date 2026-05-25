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
	XCB_ATOM_WM_PROTOCOLS     = "WM_PROTOCOLS"
	XCB_ATOM_WM_DELETE_WINDOW = "WM_DELETE_WINDOW"
)

const (
	// XCB_EVENT_CLIENT_MESSAGE is the response type for ClientMessage events.
	XCB_EVENT_CLIENT_MESSAGE uint8 = 33
	// XCB_EVENT_DESTROY_NOTIFY is the response type for DestroyNotify events.
	XCB_EVENT_DESTROY_NOTIFY uint8 = 17
)

const (
	// XcbClientMessageEventWindowOffset is the byte offset of window in xcb_client_message_event_t.
	XcbClientMessageEventWindowOffset = 4
	// XcbClientMessageEventData32Offset is the byte offset of data32[0] in xcb_client_message_event_t.
	XcbClientMessageEventData32Offset = 12
	// XcbDestroyNotifyEventWindowOffset is the byte offset of window in xcb_destroy_notify_event_t.
	XcbDestroyNotifyEventWindowOffset = 8
)

const (
	// XCB_SIZE_HINTS_FLAG_PSIZE sets width and height in xcb_size_hints_t.
	XCB_SIZE_HINTS_FLAG_PSIZE int32 = 1 << 2
	// XCB_SIZE_HINTS_FLAG_PMIN_SIZE sets min_width and min_height in xcb_size_hints_t.
	XCB_SIZE_HINTS_FLAG_PMIN_SIZE int32 = 1 << 4
	// XCB_SIZE_HINTS_FLAG_PMAX_SIZE sets max_width and max_height in xcb_size_hints_t.
	XCB_SIZE_HINTS_FLAG_PMAX_SIZE int32 = 1 << 5
)

/*
XcbInternAtomReplyAtomOffset is the byte offset of the atom field in xcb_intern_atom_reply_t.
*/
const XcbInternAtomReplyAtomOffset = 8
