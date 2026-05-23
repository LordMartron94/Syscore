//go:build linux

package loader

import (
	"fmt"

	"syscore/window/loadutil"
)

func XcbCommandsLoadManifest(module XcbModule, manifest *XcbCommandManifest, commands *XcbCommands) error {
	if manifest == nil {
		return fmt.Errorf("xcb loader: manifest must not be nil")
	}
	if commands == nil {
		return fmt.Errorf("xcb loader: commands must not be nil")
	}

	mappings := make([]loadutil.CommandMapping, 0, len(manifest.fields))
	for _, field := range manifest.fields {
		switch field {
		case XcbCommandManifestFieldConnect:
			mappings = append(mappings, loadutil.CommandMapping{Target: &commands.Connect, Name: "xcb_connect"})
		case XcbCommandManifestFieldGetSetup:
			mappings = append(mappings, loadutil.CommandMapping{Target: &commands.GetSetup, Name: "xcb_get_setup"})
		case XcbCommandManifestFieldSetupRootsIterator:
			mappings = append(mappings, loadutil.CommandMapping{Target: &commands.SetupRootsIterator, Name: "xcb_setup_roots_iterator"})
		case XcbCommandManifestFieldGenerateId:
			mappings = append(mappings, loadutil.CommandMapping{Target: &commands.GenerateId, Name: "xcb_generate_id"})
		case XcbCommandManifestFieldCreateWindow:
			mappings = append(mappings, loadutil.CommandMapping{Target: &commands.CreateWindow, Name: "xcb_create_window"})
		case XcbCommandManifestFieldMapWindow:
			mappings = append(mappings, loadutil.CommandMapping{Target: &commands.MapWindow, Name: "xcb_map_window"})
		case XcbCommandManifestFieldFlush:
			mappings = append(mappings, loadutil.CommandMapping{Target: &commands.Flush, Name: "xcb_flush"})
		case XcbCommandManifestFieldInternAtom:
			mappings = append(mappings, loadutil.CommandMapping{Target: &commands.InternAtom, Name: "xcb_intern_atom"})
		case XcbCommandManifestFieldInternAtomReply:
			mappings = append(mappings, loadutil.CommandMapping{Target: &commands.InternAtomReply, Name: "xcb_intern_atom_reply"})
		case XcbCommandManifestFieldChangeProperty:
			mappings = append(mappings, loadutil.CommandMapping{Target: &commands.ChangeProperty, Name: "xcb_change_property"})
		case XcbCommandManifestFieldFree:
			mappings = append(mappings, loadutil.CommandMapping{Target: &commands.Free, Name: "free"})
		case XcbCommandManifestFieldDestroyWindow:
			mappings = append(mappings, loadutil.CommandMapping{Target: &commands.DestroyWindow, Name: "xcb_destroy_window"})
		case XcbCommandManifestFieldDisconnect:
			mappings = append(mappings, loadutil.CommandMapping{Target: &commands.Disconnect, Name: "xcb_disconnect"})
		default:
			return fmt.Errorf("xcb loader: unknown manifest field %d", field)
		}
	}

	return xcbCommandsBind(module, commands, mappings)
}
