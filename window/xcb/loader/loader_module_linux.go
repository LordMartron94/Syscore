//go:build linux

package loader

import (
	"fmt"

	"syscore/internal"
)

const xcbLibraryName = "libxcb.so.1"

type XcbModule struct {
	Library internal.DynamicLibrary
}

func XcbModuleLoad() (XcbModule, error) {
	library, err := internal.DynamicLibraryLoad(xcbLibraryName)
	if err != nil {
		return XcbModule{}, fmt.Errorf("xcb loader: %w", err)
	}
	return XcbModule{Library: library}, nil
}

func XcbModuleLibraryName() string {
	return xcbLibraryName
}

func XcbModuleSymbolResolve(module XcbModule, symbolName string) (uintptr, error) {
	return internal.DynamicLibrarySymbolResolve(module.Library, symbolName)
}
