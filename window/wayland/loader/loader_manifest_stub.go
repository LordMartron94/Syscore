//go:build !linux

package loader

import "fmt"

func WaylandCommandsLoadManifest(_ WaylandModule, _ *WaylandCommandManifest, _ *WaylandCommands) error {
	return fmt.Errorf("wayland loader: only supported on linux")
}
