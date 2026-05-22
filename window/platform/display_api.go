package platform

import "fmt"

type DisplayAPI uint8

const (
	DisplayAPIInvalid DisplayAPI = iota
	DisplayAPIWayland
	DisplayAPIX11
	DisplayAPIWin32
	DisplayAPIAppkit
)

const (
	EnvWaylandDisplay = "WAYLAND_DISPLAY"
	EnvX11Display     = "DISPLAY"
)

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
