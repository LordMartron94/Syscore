//go:build linux

package loader

import (
	"fmt"

	"syscore/window/loadutil"
)

func WaylandCommandsLoad(module WaylandModule, commands *WaylandCommands) error {
	return waylandCommandsBind(module, commands, waylandCommandMappingsAll(commands))
}

func waylandCommandMappingsAll(commands *WaylandCommands) []loadutil.CommandMapping {
	return []loadutil.CommandMapping{
		{Target: &commands.DisplayConnect, Name: "wl_display_connect"},
		{Target: &commands.ProxyMarshalConstructor, Name: "wl_proxy_marshal_constructor"},
		{Target: &commands.ProxyMarshalConstructorVersioned, Name: "wl_proxy_marshal_constructor_versioned"},
		{Target: &commands.ProxyAddListener, Name: "wl_proxy_add_listener"},
		{Target: &commands.DisplayRoundtrip, Name: "wl_display_roundtrip"},
		{Target: &commands.ProxyMarshal, Name: "wl_proxy_marshal"},
	}
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
