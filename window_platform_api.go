package syscore

import "syscore/window/platform"

type SYSCORE_Window_DisplayAPI = platform.DisplayAPI

const (
	SYSCORE_Window_DisplayAPIInvalid = platform.DisplayAPIInvalid
	SYSCORE_Window_DisplayAPIWayland = platform.DisplayAPIWayland
	SYSCORE_Window_DisplayAPIX11     = platform.DisplayAPIX11
	SYSCORE_Window_DisplayAPIWin32   = platform.DisplayAPIWin32
	SYSCORE_Window_DisplayAPIAppkit  = platform.DisplayAPIAppkit
	SYSCORE_Window_EnvWaylandDisplay = platform.EnvWaylandDisplay
	SYSCORE_Window_EnvX11Display     = platform.EnvX11Display
)

/*
SYSCORE_Window_DisplayAPIName returns a stable lowercase name for a display API value.
*/
func SYSCORE_Window_DisplayAPIName(api SYSCORE_Window_DisplayAPI) string {
	return platform.DisplayAPIName(api)
}

/*
SYSCORE_Window_DisplayAPIDetect selects the recommended native display API for the current process.

On Linux, Wayland is preferred when WAYLAND_DISPLAY is set; otherwise X11 is used when DISPLAY is set.
On Windows the result is always Win32. On macOS the result is always AppKit.
*/
func SYSCORE_Window_DisplayAPIDetect() (SYSCORE_Window_DisplayAPI, error) {
	return platform.DisplayAPIDetect()
}

/*
SYSCORE_Window_DisplayAPIIsAvailable reports whether the given display API can be used in the
current environment (for example, Linux X11 requires DISPLAY).
*/
func SYSCORE_Window_DisplayAPIIsAvailable(api SYSCORE_Window_DisplayAPI) bool {
	return platform.DisplayAPIIsAvailable(api)
}

/*
SYSCORE_Window_DisplayAPIListAvailable returns every display API that is currently available.
Order matches detection priority on Linux (Wayland before X11 when both are set).
*/
func SYSCORE_Window_DisplayAPIListAvailable() []SYSCORE_Window_DisplayAPI {
	return platform.DisplayAPIListAvailable()
}
