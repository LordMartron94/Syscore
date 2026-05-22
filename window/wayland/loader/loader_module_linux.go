//go:build linux

package loader

import (
	"fmt"

	"syscore/internal"
	"syscore/window/loadutil"
)

const waylandLibraryName = "libwayland-client.so.0"

type WaylandModule struct {
	Library internal.DynamicLibrary
}

func WaylandModuleLoad() (WaylandModule, error) {
	library, err := internal.DynamicLibraryLoad(waylandLibraryName)
	if err != nil {
		return WaylandModule{}, fmt.Errorf("wayland loader: %w", err)
	}
	return WaylandModule{Library: library}, nil
}

func WaylandModuleLibraryName() string {
	return waylandLibraryName
}

func WaylandModuleSymbolResolve(module WaylandModule, symbolName string) (uintptr, error) {
	return loadutil.LibrarySymbolResolve(module.Library, symbolName)
}
