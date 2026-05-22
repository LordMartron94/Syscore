package platform

import "fmt"

//go:generate stringer -type DisplayAPI

/*
DisplayAPI identifies a native windowing / display client API for the current platform.
*/
type DisplayAPI uint8

const (
	DisplayAPIInvalid DisplayAPI = iota
	DisplayAPIWayland
	DisplayAPIX11
	DisplayAPIWin32
	DisplayAPIAppkit
)

const (
	// EnvWaylandDisplay is the Linux environment variable used to locate a Wayland compositor.
	EnvWaylandDisplay = "WAYLAND_DISPLAY"
	// EnvX11Display is the Linux environment variable used to locate an X11 display server.
	EnvX11Display = "DISPLAY"
)

/*
DisplayAPIName returns a stable lowercase name for api.

[Context]
Maps enumeration values to "wayland", "x11", "win32", "appkit", or "invalid".

[Parameters]
api - Display API value.

[Returns]
Short ASCII identifier.

[Side Effects]
None.
*/
func DisplayAPIName(api DisplayAPI) string {
	switch api {
	case DisplayAPIWayland:
		return "wayland"
	case DisplayAPIX11:
		return "x11"
	case DisplayAPIWin32:
		return "win32"
	case DisplayAPIAppkit:
		return "appkit"
	default:
		return "invalid"
	}
}

func displayAPIUnsupportedError() error {
	return fmt.Errorf("window platform: display API detection unsupported on this GOOS")
}
