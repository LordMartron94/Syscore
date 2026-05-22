package loader

import "syscore/window/wayland/bindings"

type WaylandCommands struct {
	DisplayConnect                   bindings.PFN_wl_display_connect
	ProxyMarshalConstructor          bindings.PFN_wl_proxy_marshal_constructor
	ProxyMarshalConstructorVersioned bindings.PFN_wl_proxy_marshal_constructor_versioned
	ProxyAddListener                 bindings.PFN_wl_proxy_add_listener
	DisplayRoundtrip                 bindings.PFN_wl_display_roundtrip
}
