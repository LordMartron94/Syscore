//go:build linux

package syscore

import (
	"fmt"

	"syscore/window/xcb/bindings"
)

/*
SYSCORE_Window_Xcb_WindowExtentGet queries the pixel width and height of an XCB window.

[Context]
Uses xcb_get_geometry on the given window ID. The caller must have bound GetGeometry and
GetGeometryReply on commands (via SYSCORE_Window_Xcb_CommandsLoad or manifest).

[Side Effects]
Sends a synchronous X request and frees the reply buffer with commands.Free.
*/
func SYSCORE_Window_Xcb_WindowExtentGet(
	commands SYSCORE_Window_Xcb_Commands,
	conn SYSCORE_Window_Xcb_ConnectionT,
	window SYSCORE_Window_Xcb_WindowT,
) (widthPx uint32, heightPx uint32, err error) {
	cookie := commands.GetGeometry(conn, uint32(window))
	reply := commands.GetGeometryReply(conn, cookie, 0)
	if reply == 0 {
		return 0, 0, fmt.Errorf("xcb_get_geometry_reply returned nil")
	}
	defer commands.Free(reply)

	widthPx = uint32(bindings.XcbGetGeometryReplyWidth(reply))
	heightPx = uint32(bindings.XcbGetGeometryReplyHeight(reply))
	if widthPx == 0 || heightPx == 0 {
		return 0, 0, fmt.Errorf("xcb_get_geometry reported zero extent")
	}

	return widthPx, heightPx, nil
}
