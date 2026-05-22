package bindings

const (
	WL_DISPLAY_GET_REGISTRY_OPCODE      uint32 = 1
	WL_REGISTRY_BIND_OPCODE             uint32 = 0
	WL_COMPOSITOR_CREATE_SURFACE_OPCODE uint32 = 0
)

const (
	WL_REGISTRY_INTERFACE_SYMBOL   = "wl_registry_interface"
	WL_COMPOSITOR_INTERFACE_SYMBOL = "wl_compositor_interface"
	WL_SURFACE_INTERFACE_SYMBOL    = "wl_surface_interface"
)
