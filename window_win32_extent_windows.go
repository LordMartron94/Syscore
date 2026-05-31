//go:build windows

package syscore

import (
	"fmt"

	foundation "syscore/window/win32/bindings/types/windows_win32_foundation"
)

/*
SYSCORE_Window_Win32_ClientExtentGet returns the client-area width and height of an HWND in pixels.

[Context]
Uses GetClientRect on the bound commands (via SYSCORE_Window_Win32_CommandsLoad or manifest).
*/
func SYSCORE_Window_Win32_ClientExtentGet(
	commands SYSCORE_Window_Win32_Commands,
	hwnd SYSCORE_Window_Win32_HWND,
) (widthPx uint32, heightPx uint32, err error) {
	var clientRect foundation.RECT
	if commands.GetClientRect(hwnd, &clientRect) == 0 {
		return 0, 0, fmt.Errorf("GetClientRect failed")
	}

	widthPx = uint32(clientRect.Right - clientRect.Left)
	heightPx = uint32(clientRect.Bottom - clientRect.Top)
	if widthPx == 0 || heightPx == 0 {
		return 0, 0, fmt.Errorf("GetClientRect reported zero extent")
	}

	return widthPx, heightPx, nil
}
