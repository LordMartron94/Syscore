package loader

import "syscore/window/wayland/bindings"

/*
WaylandCommands holds bound Wayland client function pointers from libwayland-client.so.0.

[Context]
Populate with WaylandCommandsLoad or WaylandCommandsLoadManifest. Resolve wl_*_interface
globals with WaylandModuleSymbolResolve before marshal_constructor calls.
*/
type WaylandCommands struct {
	/*
		DisplayConnect is wl_display_connect. Connects to the compositor (nil uses WAYLAND_DISPLAY).
	*/
	DisplayConnect bindings.PFN_wl_display_connect
	/*
		ProxyMarshalConstructor is wl_proxy_marshal_constructor. Creates registry, surface, etc.
	*/
	ProxyMarshalConstructor bindings.PFN_wl_proxy_marshal_constructor
	/*
		ProxyMarshalConstructorVersioned is wl_proxy_marshal_constructor_versioned. Binds globals from registry.
	*/
	ProxyMarshalConstructorVersioned bindings.PFN_wl_proxy_marshal_constructor_versioned
	/*
		ProxyAddListener is wl_proxy_add_listener. Installs registry (and other) listeners.
	*/
	ProxyAddListener bindings.PFN_wl_proxy_add_listener
	/*
		DisplayRoundtrip is wl_display_roundtrip. Dispatches events after requests.
	*/
	DisplayRoundtrip bindings.PFN_wl_display_roundtrip
	/*
		DisplayDispatchPending is wl_display_dispatch_pending. Dispatches queued events without blocking.
	*/
	DisplayDispatchPending bindings.PFN_wl_display_dispatch_pending
	/*
		DisplayGetFd is wl_display_get_fd. Returns the compositor connection socket fd.
	*/
	DisplayGetFd bindings.PFN_wl_display_get_fd
	/*
		DisplayPrepareRead is wl_display_prepare_read. Reserves read access on the display fd.
	*/
	DisplayPrepareRead bindings.PFN_wl_display_prepare_read
	/*
		DisplayReadEvents is wl_display_read_events. Reads events from the fd into the queue.
	*/
	DisplayReadEvents bindings.PFN_wl_display_read_events
	/*
		DisplayCancelRead is wl_display_cancel_read. Cancels a prepared read when the fd is not readable.
	*/
	DisplayCancelRead bindings.PFN_wl_display_cancel_read
	/*
		ProxyMarshal is wl_proxy_marshal. Sends requests on existing proxies (xdg title, ack_configure, etc.).
	*/
	ProxyMarshal bindings.PFN_wl_proxy_marshal
	/*
		ProxyDestroy is wl_proxy_destroy. Destroys Wayland protocol objects.
	*/
	ProxyDestroy bindings.PFN_wl_proxy_destroy
	/*
		DisplayDisconnect is wl_display_disconnect. Closes the compositor connection.
	*/
	DisplayDisconnect bindings.PFN_wl_display_disconnect
}
