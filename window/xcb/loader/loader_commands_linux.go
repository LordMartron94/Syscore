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

func xcbCommandMappingsAll(commands *XcbCommands) []loadutil.CommandMapping {
	return []loadutil.CommandMapping{
		{Target: &commands.Connect, Name: "xcb_connect"},
		{Target: &commands.GetSetup, Name: "xcb_get_setup"},
		{Target: &commands.SetupRootsIterator, Name: "xcb_setup_roots_iterator"},
		{Target: &commands.GenerateId, Name: "xcb_generate_id"},
		{Target: &commands.CreateWindow, Name: "xcb_create_window"},
		{Target: &commands.MapWindow, Name: "xcb_map_window"},
		{Target: &commands.Flush, Name: "xcb_flush"},
		{Target: &commands.InternAtom, Name: "xcb_intern_atom"},
		{Target: &commands.InternAtomReply, Name: "xcb_intern_atom_reply"},
		{Target: &commands.ChangeProperty, Name: "xcb_change_property"},
		{Target: &commands.Free, Name: "free"},
	}
}

func xcbCommandsBind(module XcbModule, commands *XcbCommands, mappings []loadutil.CommandMapping) error {
	if commands == nil {
		return fmt.Errorf("xcb loader: commands must not be nil")
	}

	xcbMappings, libcMappings := xcbCommandMappingsPartition(commands, mappings)
	if err := loadutil.LibraryCommandsBind(module.Library, xcbMappings); err != nil {
		return fmt.Errorf("xcb loader: %w", err)
	}
	if len(libcMappings) == 0 {
		return nil
	}

	libc, err := internal.DynamicLibraryLoad(xcbLibcLibraryName)
	if err != nil {
		return fmt.Errorf("xcb loader: libc for reply free: %w", err)
	}
	if err := loadutil.LibraryCommandsBind(libc, libcMappings); err != nil {
		return fmt.Errorf("xcb loader: %w", err)
	}
	return nil
}

func xcbCommandMappingsPartition(commands *XcbCommands, mappings []loadutil.CommandMapping) (xcbMappings []loadutil.CommandMapping, libcMappings []loadutil.CommandMapping) {
	for _, mapping := range mappings {
		if mapping.Target == &commands.Free {
			libcMappings = append(libcMappings, loadutil.CommandMapping{
				Target: mapping.Target,
				Name:   "free",
			})
			continue
		}
		xcbMappings = append(xcbMappings, mapping)
	}
	return xcbMappings, libcMappings
}
