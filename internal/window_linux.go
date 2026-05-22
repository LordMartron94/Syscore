//go:build linux

package internal

import (
	"fmt"
	"os"
	"sync"
	"unsafe"

	"github.com/ebitengine/purego"
)

// --- XCB FFI ---

var (
	xcbLibInit sync.Once
	xcbLibErr  error

	xcb_connect              func(displayname *byte, screenp *int32) uintptr
	xcb_get_setup            func(c uintptr) uintptr
	xcb_setup_roots_iterator func(setup uintptr) xcb_screen_iterator_t
	xcb_generate_id          func(c uintptr) uint32
	xcb_create_window        func(c uintptr, depth uint8, wid uint32, parent uint32, x int16, y int16, width uint16, height uint16, border_width uint16, class uint16, visual uint32, value_mask uint32, value_list unsafe.Pointer) uint32
	xcb_map_window           func(c uintptr, window uint32) uint32
	xcb_flush                func(c uintptr) int32
)

type xcb_screen_iterator_t struct {
	data  uintptr
	rem   int32
	index int32
}

const xcbWindowClassInputOutput uint16 = 1

func xcbLibraryEnsureLoaded() error {
	xcbLibInit.Do(func() {
		lib, err := DynamicLibraryLoad("libxcb.so.1")
		if err != nil {
			xcbLibErr = err
			return
		}
		bindings := []struct {
			name string
			ptr  any
		}{
			{"xcb_connect", &xcb_connect},
			{"xcb_get_setup", &xcb_get_setup},
			{"xcb_setup_roots_iterator", &xcb_setup_roots_iterator},
			{"xcb_generate_id", &xcb_generate_id},
			{"xcb_create_window", &xcb_create_window},
			{"xcb_map_window", &xcb_map_window},
			{"xcb_flush", &xcb_flush},
		}
		for _, binding := range bindings {
			if err := LibraryFunctionBind(lib, binding.ptr, binding.name); err != nil {
				xcbLibErr = err
				return
			}
		}
	})
	return xcbLibErr
}

func windowCreateForX11() (uintptr, uintptr, error) {
	if err := xcbLibraryEnsureLoaded(); err != nil {
		return 0, 0, err
	}

	conn := xcb_connect(nil, nil)
	if conn == 0 {
		return 0, 0, fmt.Errorf("xcb_connect failed")
	}

	setup := xcb_get_setup(conn)
	iter := xcb_setup_roots_iterator(setup)
	if iter.data == 0 {
		return 0, 0, fmt.Errorf("xcb: no screens in setup")
	}

	rootWindowID := *(*uint32)(unsafe.Pointer(iter.data))
	windowID := xcb_generate_id(conn)

	xcb_create_window(
		conn,
		0,
		windowID,
		rootWindowID,
		0, 0,
		800, 600,
		0,
		xcbWindowClassInputOutput,
		0,
		0, nil,
	)

	xcb_map_window(conn, windowID)
	if xcb_flush(conn) < 0 {
		return 0, 0, fmt.Errorf("xcb_flush failed")
	}

	return conn, uintptr(windowID), nil
}

// --- Wayland FFI ---

var (
	waylandLibInit sync.Once
	waylandLibErr  error

	wl_display_connect                     func(name *byte) uintptr
	wl_proxy_marshal_constructor           func(proxy uintptr, opcode uint32, interfacePtr uintptr, args ...any) uintptr
	wl_proxy_marshal_constructor_versioned func(proxy uintptr, opcode uint32, interfacePtr uintptr, version uint32, args ...any) uintptr
	wl_proxy_add_listener                  func(proxy uintptr, implementation uintptr, data uintptr) int32
	wl_display_roundtrip                   func(display uintptr) int32
)

type wl_registry_listener struct {
	global        uintptr
	global_remove uintptr
}

type waylandWindowCreateState struct {
	compositor uintptr
}

func waylandLibraryEnsureLoaded() error {
	waylandLibInit.Do(func() {
		lib, err := DynamicLibraryLoad("libwayland-client.so.0")
		if err != nil {
			waylandLibErr = err
			return
		}
		bindings := []struct {
			name string
			ptr  any
		}{
			{"wl_display_connect", &wl_display_connect},
			{"wl_proxy_marshal_constructor", &wl_proxy_marshal_constructor},
			{"wl_proxy_marshal_constructor_versioned", &wl_proxy_marshal_constructor_versioned},
			{"wl_proxy_add_listener", &wl_proxy_add_listener},
			{"wl_display_roundtrip", &wl_display_roundtrip},
		}
		for _, binding := range bindings {
			if err := LibraryFunctionBind(lib, binding.ptr, binding.name); err != nil {
				waylandLibErr = err
				return
			}
		}
	})
	return waylandLibErr
}

func waylandInterfaceResolve(library DynamicLibrary, symbolName string) (uintptr, error) {
	symbol, err := DynamicLibrarySymbolResolve(library, symbolName)
	if err != nil {
		return 0, fmt.Errorf("resolve %s: %w", symbolName, err)
	}
	if symbol == 0 {
		return 0, fmt.Errorf("resolve %s: symbol not found", symbolName)
	}
	return symbol, nil
}

func windowCreateForWayland() (uintptr, uintptr, error) {
	if err := waylandLibraryEnsureLoaded(); err != nil {
		return 0, 0, err
	}

	lib, err := DynamicLibraryLoad("libwayland-client.so.0")
	if err != nil {
		return 0, 0, err
	}

	wlRegistryInterface, err := waylandInterfaceResolve(lib, "wl_registry_interface")
	if err != nil {
		return 0, 0, err
	}
	wlCompositorInterface, err := waylandInterfaceResolve(lib, "wl_compositor_interface")
	if err != nil {
		return 0, 0, err
	}
	wlSurfaceInterface, err := waylandInterfaceResolve(lib, "wl_surface_interface")
	if err != nil {
		return 0, 0, err
	}

	display := wl_display_connect(nil)
	if display == 0 {
		return 0, 0, fmt.Errorf("wl_display_connect failed")
	}

	registry := wl_proxy_marshal_constructor(display, 1, wlRegistryInterface, nil)
	if registry == 0 {
		return 0, 0, fmt.Errorf("wl_display_get_registry failed")
	}

	state := waylandWindowCreateState{}
	globalCallback := purego.NewCallback(func(data uintptr, registryProxy uintptr, name uint32, interfaceName uintptr, version uint32) {
		if CStringPointerToString(interfaceName) != "wl_compositor" {
			return
		}
		createState := (*waylandWindowCreateState)(unsafe.Pointer(data))
		createState.compositor = wl_proxy_marshal_constructor_versioned(
			registryProxy,
			0,
			wlCompositorInterface,
			version,
			name,
			nil,
		)
	})

	listener := wl_registry_listener{
		global:        globalCallback,
		global_remove: 0,
	}
	listenerPtr := uintptr(unsafe.Pointer(&listener))

	if wl_proxy_add_listener(registry, listenerPtr, uintptr(unsafe.Pointer(&state))) < 0 {
		return 0, 0, fmt.Errorf("wl_proxy_add_listener failed")
	}

	if wl_display_roundtrip(display) < 0 {
		return 0, 0, fmt.Errorf("wl_display_roundtrip failed")
	}

	if state.compositor == 0 {
		return 0, 0, fmt.Errorf("wayland: wl_compositor global not advertised")
	}

	surface := wl_proxy_marshal_constructor(state.compositor, 0, wlSurfaceInterface, nil)
	if surface == 0 {
		return 0, 0, fmt.Errorf("wl_compositor_create_surface failed")
	}

	return display, surface, nil
}

func WindowCreate() (Window, error) {
	if _, waylandOk := os.LookupEnv(waylandEnv); waylandOk {
		display, surface, err := windowCreateForWayland()
		if err == nil {
			return Window{
				Type:          WindowTypeWayland,
				DisplayHandle: display,
				WindowHandle:  surface,
			}, nil
		}
	}

	if _, x11Ok := os.LookupEnv(x11Env); x11Ok {
		conn, windowID, err := windowCreateForX11()
		if err != nil {
			return Window{}, fmt.Errorf("x11 initialization failed: %w", err)
		}
		return Window{
			Type:          WindowTypeX11,
			DisplayHandle: conn,
			WindowHandle:  windowID,
		}, nil
	}

	return Window{}, fmt.Errorf("unsupported: no wayland or x11 environment found")
}
