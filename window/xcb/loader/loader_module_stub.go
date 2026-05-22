//go:build !linux

package loader

import (
	"fmt"

	"syscore/internal"
)

type XcbModule struct {
	Library internal.DynamicLibrary
}

func XcbModuleLoad() (XcbModule, error) {
	return XcbModule{}, fmt.Errorf("xcb loader: only supported on linux")
}

func XcbModuleLibraryName() string {
	return ""
}

func XcbModuleSymbolResolve(_ XcbModule, _ string) (uintptr, error) {
	return 0, fmt.Errorf("xcb loader: only supported on linux")
}
