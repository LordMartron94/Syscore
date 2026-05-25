package bindings

/*
XdgWmBase is an opaque struct xdg_wm_base* (shell global from the registry).
*/
type XdgWmBase uintptr

/*
XdgSurface is an opaque struct xdg_surface* (shell surface role).
*/
type XdgSurface uintptr

/*
XdgToplevel is an opaque struct xdg_toplevel* (desktop window role).
*/
type XdgToplevel uintptr

/*
XdgSurfaceListener is the xdg_surface_listener vtable for wl_proxy_add_listener.

Configure is a purego callback invoked as (data, xdgSurface, serial).
*/
type XdgSurfaceListener struct {
	Configure uintptr
}

/*
XdgToplevelListener is the xdg_toplevel_listener vtable for wl_proxy_add_listener.

Configure is invoked as (data, xdgToplevel, width, height, states). Close is (data, xdgToplevel).
Unset entries must be 0.
*/
type XdgToplevelListener struct {
	Configure uintptr
	Close     uintptr
}
