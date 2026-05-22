//go:build linux

package platform

import (
	"fmt"
	"os"
)

/*
DisplayAPIDetect selects Wayland when WAYLAND_DISPLAY is set, otherwise X11 when DISPLAY is set.

See package syscore SYSCORE_Window_DisplayAPIDetect for the public facade contract.
*/
func DisplayAPIDetect() (DisplayAPI, error) {
	if DisplayAPIIsAvailable(DisplayAPIWayland) {
		return DisplayAPIWayland, nil
	}
	if DisplayAPIIsAvailable(DisplayAPIX11) {
		return DisplayAPIX11, nil
	}
	return DisplayAPIInvalid, fmt.Errorf("window platform: no %s or %s environment found", EnvWaylandDisplay, EnvX11Display)
}

/*
DisplayAPIIsAvailable reports whether api is usable on Linux (environment variables present).
*/
func DisplayAPIIsAvailable(api DisplayAPI) bool {
	switch api {
	case DisplayAPIWayland:
		_, ok := os.LookupEnv(EnvWaylandDisplay)
		return ok
	case DisplayAPIX11:
		_, ok := os.LookupEnv(EnvX11Display)
		return ok
	default:
		return false
	}
}

/*
DisplayAPIListAvailable returns Wayland and/or X11 when their environment variables are set (Wayland first).
*/
func DisplayAPIListAvailable() []DisplayAPI {
	available := make([]DisplayAPI, 0, 2)
	if DisplayAPIIsAvailable(DisplayAPIWayland) {
		available = append(available, DisplayAPIWayland)
	}
	if DisplayAPIIsAvailable(DisplayAPIX11) {
		available = append(available, DisplayAPIX11)
	}
	return available
}
