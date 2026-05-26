/*
Package window provides platform-native windowing bindings and loaders.

Import the platform packages directly for full control:

	syscore/window/xcb/bindings
	syscore/window/xcb/loader

	syscore/window/wayland/bindings
	syscore/window/wayland/loader

	syscore/window/win32/bindings
	syscore/window/win32/loader

	syscore/window/appkit/bindings
	syscore/window/appkit/loader

The root syscore package exposes SYSCORE_Window_* facades over these loaders.

Display API selection (WAYLAND_DISPLAY / DISPLAY / GOOS) lives in syscore/window/platform.
*/
package window

//go:generate go run ./internal/generate -specRoot=./internal/spec -waylandBindings=./wayland/bindings -xdgBindings=./wayland/xdg/bindings -xcbBindings=./xcb/bindings -waylandLoader=./wayland/loader -xcbLoader=./xcb/loader -win32Bindings=./win32/bindings -win32Loader=./win32/loader -appkitBindings=./appkit/bindings -appkitLoader=./appkit/loader
