package bindings

const (
	// WL_DISPLAY_GET_REGISTRY_OPCODE is wl_display.get_registry (constructor opcode 1).
	WL_DISPLAY_GET_REGISTRY_OPCODE uint32 = 1
	// WL_REGISTRY_BIND_OPCODE is wl_registry.bind.
	WL_REGISTRY_BIND_OPCODE uint32 = 0
	// WL_COMPOSITOR_CREATE_SURFACE_OPCODE is wl_compositor.create_surface.
	WL_COMPOSITOR_CREATE_SURFACE_OPCODE uint32 = 0
	// WL_SURFACE_COMMIT_OPCODE is wl_surface.commit.
	WL_SURFACE_COMMIT_OPCODE uint32 = 6
)

const (
	// WL_REGISTRY_INTERFACE_SYMBOL is the exported wl_registry_interface global symbol name.
	WL_REGISTRY_INTERFACE_SYMBOL = "wl_registry_interface"
	// WL_COMPOSITOR_INTERFACE_SYMBOL is the exported wl_compositor_interface global symbol name.
	WL_COMPOSITOR_INTERFACE_SYMBOL = "wl_compositor_interface"
	// WL_SURFACE_INTERFACE_SYMBOL is the exported wl_surface_interface global symbol name.
	WL_SURFACE_INTERFACE_SYMBOL = "wl_surface_interface"
)
