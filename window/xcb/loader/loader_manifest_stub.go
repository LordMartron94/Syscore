//go:build !linux

package loader

import "fmt"

func XcbCommandsLoadManifest(_ XcbModule, _ *XcbCommandManifest, _ *XcbCommands) error {
	return fmt.Errorf("xcb loader: only supported on linux")
}
