package bindings

const (
	// XDG_WM_BASE_GLOBAL_NAME is the registry global name for xdg_wm_base.
	XDG_WM_BASE_GLOBAL_NAME = "xdg_wm_base"

	// XDG_WM_BASE_INTERFACE_SYMBOL is the exported wl_interface symbol in libsyscore_wayland_xdg.so.
	XDG_WM_BASE_INTERFACE_SYMBOL = "xdg_wm_base_interface"
	// XDG_SURFACE_INTERFACE_SYMBOL is the exported wl_interface symbol for xdg_surface.
	XDG_SURFACE_INTERFACE_SYMBOL = "xdg_surface_interface"
	// XDG_TOPLEVEL_INTERFACE_SYMBOL is the exported wl_interface symbol for xdg_toplevel.
	XDG_TOPLEVEL_INTERFACE_SYMBOL = "xdg_toplevel_interface"

	// XDG_WM_BASE_GET_XDG_SURFACE_OPCODE is xdg_wm_base.get_xdg_surface.
	XDG_WM_BASE_GET_XDG_SURFACE_OPCODE uint32 = 2

	// XDG_SURFACE_GET_TOPLEVEL_OPCODE is xdg_surface.get_toplevel.
	XDG_SURFACE_GET_TOPLEVEL_OPCODE uint32 = 1
	// XDG_SURFACE_ACK_CONFIGURE_OPCODE is xdg_surface.ack_configure.
	XDG_SURFACE_ACK_CONFIGURE_OPCODE uint32 = 4
	// XDG_SURFACE_CONFIGURE_EVENT_OPCODE is the configure event on xdg_surface.
	XDG_SURFACE_CONFIGURE_EVENT_OPCODE uint32 = 0

	// XDG_TOPLEVEL_SET_TITLE_OPCODE is xdg_toplevel.set_title.
	XDG_TOPLEVEL_SET_TITLE_OPCODE uint32 = 2
	// XDG_TOPLEVEL_SET_MIN_SIZE_OPCODE is xdg_toplevel.set_min_size.
	XDG_TOPLEVEL_SET_MIN_SIZE_OPCODE uint32 = 8
)
