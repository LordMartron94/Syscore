//go:build linux

package loader

import (
	"fmt"

	"syscore/window/loadutil"
)

func WaylandCommandsLoad(module WaylandModule, commands *WaylandCommands) error {
	return waylandCommandsBind(module, commands, waylandCommandMappingsAll(commands))
}

func waylandCommandsBind(module WaylandModule, commands *WaylandCommands, mappings []loadutil.CommandMapping) error {
	if commands == nil {
		return fmt.Errorf("wayland loader: commands must not be nil")
	}
	if err := loadutil.LibraryCommandsBind(module.Library, mappings); err != nil {
		return fmt.Errorf("wayland loader: %w", err)
	}
	return nil
}
