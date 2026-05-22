//go:build darwin

package internal

import (
	"github.com/ebitengine/purego"
)

func loadPlatformLibrary(name string) (uintptr, error) {
	return purego.Dlopen(name, purego.RTLD_NOW|purego.RTLD_GLOBAL)
}

func platformLibrarySymbolResolve(library DynamicLibrary, symbolName string) (uintptr, error) {
	return purego.Dlsym(uintptr(library), symbolName)
}
