//go:build !linux

package loader

import "fmt"

/*
WaylandCommands is the non-linux stub of the bound Wayland client function pointer table.

[Context]
Real WaylandCommands fields are defined in loader_types.go under the linux build tag. Non-linux
builds need a placeholder so cross-platform consumers can compile and route around it at runtime.
*/
type WaylandCommands struct{}

/*
WaylandCommandManifest is the non-linux stub of the selective-load manifest type.
*/
type WaylandCommandManifest struct{}

func WaylandCommandsLoad(_ WaylandModule, _ *WaylandCommands) error {
	return fmt.Errorf("wayland loader: only supported on linux")
}
