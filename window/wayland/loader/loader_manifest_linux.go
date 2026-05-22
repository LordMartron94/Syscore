//go:build linux

package loader

import (
	"fmt"

	"syscore/window/loadutil"
)

func WaylandCommandsLoadManifest(module WaylandModule, manifest *WaylandCommandManifest, commands *WaylandCommands) error {
	if manifest == nil {
		return fmt.Errorf("wayland loader: manifest must not be nil")
	}
	if commands == nil {
		return fmt.Errorf("wayland loader: commands must not be nil")
	}

	mappings := make([]loadutil.CommandMapping, 0, len(manifest.fields))
	for _, field := range manifest.fields {
		switch field {
		case WaylandCommandManifestFieldDisplayConnect:
			mappings = append(mappings, loadutil.CommandMapping{Target: &commands.DisplayConnect, Name: "wl_display_connect"})
		case WaylandCommandManifestFieldProxyMarshalConstructor:
			mappings = append(mappings, loadutil.CommandMapping{Target: &commands.ProxyMarshalConstructor, Name: "wl_proxy_marshal_constructor"})
		case WaylandCommandManifestFieldProxyMarshalConstructorVersioned:
			mappings = append(mappings, loadutil.CommandMapping{Target: &commands.ProxyMarshalConstructorVersioned, Name: "wl_proxy_marshal_constructor_versioned"})
		case WaylandCommandManifestFieldProxyAddListener:
			mappings = append(mappings, loadutil.CommandMapping{Target: &commands.ProxyAddListener, Name: "wl_proxy_add_listener"})
		case WaylandCommandManifestFieldDisplayRoundtrip:
			mappings = append(mappings, loadutil.CommandMapping{Target: &commands.DisplayRoundtrip, Name: "wl_display_roundtrip"})
		case WaylandCommandManifestFieldProxyMarshal:
			mappings = append(mappings, loadutil.CommandMapping{Target: &commands.ProxyMarshal, Name: "wl_proxy_marshal"})
		default:
			return fmt.Errorf("wayland loader: unknown manifest field %d", field)
		}
	}

	return waylandCommandsBind(module, commands, mappings)
}
