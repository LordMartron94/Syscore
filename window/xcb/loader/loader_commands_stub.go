//go:build !linux

package loader

import "fmt"

/*
XcbCommands is the non-linux stub of the bound XCB function pointer table.

[Context]
Real XcbCommands fields are defined in loader_types.go under the linux build tag. Non-linux builds
need a placeholder so cross-platform consumers can compile and route around it at runtime.
*/
type XcbCommands struct{}

/*
XcbCommandManifest is the non-linux stub of the selective-load manifest type.
*/
type XcbCommandManifest struct{}

func XcbCommandsLoad(_ XcbModule, _ *XcbCommands) error {
	return fmt.Errorf("xcb loader: only supported on linux")
}
