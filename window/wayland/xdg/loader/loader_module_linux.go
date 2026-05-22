//go:build linux

package loader

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"

	"syscore/internal"
	"syscore/window/loadutil"
)

//go:embed embed/libsyscore_wayland_xdg.so
var xdgEmbeddedLibrary []byte

const xdgEmbeddedLibraryFileName = "libsyscore_wayland_xdg.so"

type XdgModule struct {
	Library internal.DynamicLibrary
}

func XdgModuleLoad() (XdgModule, error) {
	path, err := xdgEmbeddedLibraryMaterialize()
	if err != nil {
		return XdgModule{}, fmt.Errorf("xdg loader: %w", err)
	}
	library, err := internal.DynamicLibraryLoad(path)
	if err != nil {
		return XdgModule{}, fmt.Errorf("xdg loader: %w", err)
	}
	return XdgModule{Library: library}, nil
}

func XdgModuleEmbeddedLibraryFileName() string {
	return xdgEmbeddedLibraryFileName
}

func XdgModuleSymbolResolve(module XdgModule, symbolName string) (uintptr, error) {
	return loadutil.LibrarySymbolResolve(module.Library, symbolName)
}

func xdgEmbeddedLibraryMaterialize() (string, error) {
	path := filepath.Join(os.TempDir(), "syscore-"+xdgEmbeddedLibraryFileName)
	if err := os.WriteFile(path, xdgEmbeddedLibrary, 0o755); err != nil {
		return "", fmt.Errorf("write embedded library: %w", err)
	}
	return path, nil
}
