//go:build !linux

package loader

import (
	"fmt"

	"syscore/internal"
)

type WaylandModule struct {
	Library internal.DynamicLibrary
}

func WaylandModuleLoad() (WaylandModule, error) {
	return WaylandModule{}, fmt.Errorf("wayland loader: only supported on linux")
}

func WaylandModuleLibraryName() string {
	return ""
}

func WaylandModuleSymbolResolve(_ WaylandModule, _ string) (uintptr, error) {
	return 0, fmt.Errorf("wayland loader: only supported on linux")
}
