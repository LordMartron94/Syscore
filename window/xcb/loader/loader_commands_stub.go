//go:build !linux

package loader

import "fmt"

func XcbCommandsLoad(_ XcbModule, _ *XcbCommands) error {
	return fmt.Errorf("xcb loader: only supported on linux")
}
