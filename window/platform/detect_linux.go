//go:build linux

package platform

import (
	"fmt"
	"os"
)

func DisplayAPIDetect() (DisplayAPI, error) {
	if DisplayAPIIsAvailable(DisplayAPIWayland) {
		return DisplayAPIWayland, nil
	}
	if DisplayAPIIsAvailable(DisplayAPIX11) {
		return DisplayAPIX11, nil
	}
	return DisplayAPIInvalid, fmt.Errorf("window platform: no %s or %s environment found", EnvWaylandDisplay, EnvX11Display)
}

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
