//go:build !linux

package loader

import "fmt"

func WaylandCommandsLoad(_ WaylandModule, _ *WaylandCommands) error {
	return fmt.Errorf("wayland loader: only supported on linux")
}
