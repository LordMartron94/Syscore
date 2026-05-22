//go:build linux

package loader

import (
	"fmt"

	"syscore/window/loadutil"
)

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
		{Target: &commands.Free, Name: "xcb_free"},
	}
}

func xcbCommandsBind(module XcbModule, commands *XcbCommands, mappings []loadutil.CommandMapping) error {
	if commands == nil {
		return fmt.Errorf("xcb loader: commands must not be nil")
	}
	if err := loadutil.LibraryCommandsBind(module.Library, mappings); err != nil {
		return fmt.Errorf("xcb loader: %w", err)
	}
	return nil
}
