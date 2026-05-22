package syscore

import "syscore/window/platform"

/*
SYSCORE_Window_DisplayAPI identifies a native windowing / display client API.

[Context]
Returned by detection helpers and used to choose which SYSCORE_Window_<Backend>_ loader to
initialize. SYSCORE_Window_DisplayAPIInvalid means no API was selected.
*/
type SYSCORE_Window_DisplayAPI = platform.DisplayAPI

const (
	// SYSCORE_Window_DisplayAPIInvalid is returned when detection fails or the API is unknown.
	SYSCORE_Window_DisplayAPIInvalid = platform.DisplayAPIInvalid
	// SYSCORE_Window_DisplayAPIWayland is the Wayland client API (libwayland-client).
	SYSCORE_Window_DisplayAPIWayland = platform.DisplayAPIWayland
	// SYSCORE_Window_DisplayAPIX11 is the X11 API via XCB (libxcb).
	SYSCORE_Window_DisplayAPIX11 = platform.DisplayAPIX11
	// SYSCORE_Window_DisplayAPIWin32 is the Win32 USER/GDI window API.
	SYSCORE_Window_DisplayAPIWin32 = platform.DisplayAPIWin32
	// SYSCORE_Window_DisplayAPIAppkit is the macOS AppKit / Objective-C runtime window API.
	SYSCORE_Window_DisplayAPIAppkit = platform.DisplayAPIAppkit
	// SYSCORE_Window_EnvWaylandDisplay is the Linux environment variable name "WAYLAND_DISPLAY".
	SYSCORE_Window_EnvWaylandDisplay = platform.EnvWaylandDisplay
	// SYSCORE_Window_EnvX11Display is the Linux environment variable name "DISPLAY".
	SYSCORE_Window_EnvX11Display = platform.EnvX11Display
)

/*
SYSCORE_Window_DisplayAPIName returns a stable lowercase name for a display API value.

[Context]
Maps enumeration values to "wayland", "x11", "win32", "appkit", or "invalid" for logging and
configuration. Does not inspect the environment.

[Parameters]
api - Value from SYSCORE_Window_DisplayAPIDetect or SYSCORE_Window_DisplayAPIListAvailable.

[Returns]
A short ASCII identifier.

[Side Effects]
None. Pure function.
*/
func SYSCORE_Window_DisplayAPIName(api SYSCORE_Window_DisplayAPI) string {
	return platform.DisplayAPIName(api)
}

/*
SYSCORE_Window_DisplayAPIDetect selects the recommended native display API for this process.

[Context]
Encodes the default session policy: on Linux prefer Wayland when WAYLAND_DISPLAY is set,
otherwise X11 when DISPLAY is set; on Windows use Win32; on macOS use AppKit. Does not load
libraries or verify that the compositor or X server is reachable.

[Returns]
The recommended API and nil error on success.

[Errors]
On Linux returns SYSCORE_Window_DisplayAPIInvalid and an error when neither WAYLAND_DISPLAY
nor DISPLAY is set. On unsupported GOOS returns Invalid and a platform error.

[Side Effects]
Reads environment variables only.

[Example]

	api, err := syscore.SYSCORE_Window_DisplayAPIDetect()
	if err != nil { ... }
	switch api {
	case syscore.SYSCORE_Window_DisplayAPIWayland:
		module, err := syscore.SYSCORE_Window_Wayland_ModuleLoad()
	}
*/
func SYSCORE_Window_DisplayAPIDetect() (SYSCORE_Window_DisplayAPI, error) {
	return platform.DisplayAPIDetect()
}

/*
SYSCORE_Window_DisplayAPIIsAvailable reports whether a display API can be used in the current environment.

[Context]
On Linux Wayland requires WAYLAND_DISPLAY; X11 requires DISPLAY. On Windows only Win32 is
available; on macOS only AppKit. Does not perform a connection handshake.

[Parameters]
api - API to test.

[Returns]
true when the API is supported on this GOOS and required environment (if any) is present.

[Side Effects]
May read environment variables on Linux.
*/
func SYSCORE_Window_DisplayAPIIsAvailable(api SYSCORE_Window_DisplayAPI) bool {
	return platform.DisplayAPIIsAvailable(api)
}

/*
SYSCORE_Window_DisplayAPIListAvailable returns every display API that is currently available.

[Context]
On Linux may return both Wayland and X11 when both environment variables are set. Order is
Wayland before X11 to match SYSCORE_Window_DisplayAPIDetect priority. On Windows and macOS
returns a single-element slice.

[Returns]
A possibly empty slice of APIs; no error is returned—use SYSCORE_Window_DisplayAPIDetect when
you need a single recommended choice or an error if none apply.

[Side Effects]
May read environment variables on Linux.
*/
func SYSCORE_Window_DisplayAPIListAvailable() []SYSCORE_Window_DisplayAPI {
	return platform.DisplayAPIListAvailable()
}
