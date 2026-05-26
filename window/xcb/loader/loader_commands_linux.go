//go:build linux

package loader

import (
	"fmt"

	"syscore/internal"
	"syscore/window/loadutil"
)

const xcbLibcLibraryName = "libc.so.6"

func XcbCommandsLoad(module XcbModule, commands *XcbCommands) error {
	return xcbCommandsBind(module, commands, xcbCommandMappingsAll(commands))
}

func xcbCommandsBind(module XcbModule, commands *XcbCommands, mappings []loadutil.CommandMapping) error {
	if commands == nil {
		return fmt.Errorf("xcb loader: commands must not be nil")
	}

	xcbMappings := make([]loadutil.CommandMapping, 0, len(mappings))
	var freeMapping loadutil.CommandMapping
	hasFree := false
	for _, mapping := range mappings {
		if mapping.Name == "free" {
			freeMapping = mapping
			hasFree = true
			continue
		}
		xcbMappings = append(xcbMappings, mapping)
	}

	if err := loadutil.LibraryCommandsBind(module.Library, xcbMappings); err != nil {
		return fmt.Errorf("xcb loader: %w", err)
	}
	if !hasFree {
		return nil
	}

	libc, err := internal.DynamicLibraryLoad(xcbLibcLibraryName)
	if err != nil {
		return fmt.Errorf("xcb loader: libc: %w", err)
	}
	if err := loadutil.LibraryCommandsBind(libc, []loadutil.CommandMapping{freeMapping}); err != nil {
		return fmt.Errorf("xcb loader: libc free: %w", err)
	}
	return nil
}
