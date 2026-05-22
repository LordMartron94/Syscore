//go:build !linux

package loader

import "fmt"

type XdgModule struct{}

func XdgModuleLoad() (XdgModule, error) {
	return XdgModule{}, fmt.Errorf("xdg loader: unsupported GOOS")
}

func XdgModuleEmbeddedLibraryFileName() string {
	return ""
}

func XdgModuleSymbolResolve(module XdgModule, symbolName string) (uintptr, error) {
	_ = module
	_ = symbolName
	return 0, fmt.Errorf("xdg loader: unsupported GOOS")
}
