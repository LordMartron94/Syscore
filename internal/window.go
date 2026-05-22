package internal

const (
	waylandEnv = "WAYLAND_DISPLAY"
	x11Env     = "DISPLAY"
)

type WindowType uint8

const (
	WindowTypeWayland WindowType = iota
	WindowTypeX11
)

type Window struct {
	Type          WindowType
	DisplayHandle uintptr // Wayland: *wl_display       | X11: *xcb_connection_t
	WindowHandle  uintptr // Wayland: *wl_surface       | X11: xcb_window_t (cast to uintptr)
}
