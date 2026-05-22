//go:build !linux && !windows && !darwin

package internal

import (
	"fmt"
	"runtime"
)

func loadPlatformLibrary(name string) (uintptr, error) {
	return 0, fmt.Errorf("syscore: dynamic library load unsupported on GOOS %q", runtime.GOOS)
}

func platformLibrarySymbolResolve(_ DynamicLibrary, _ string) (uintptr, error) {
	return 0, fmt.Errorf("syscore: dynamic library symbol resolve unsupported on GOOS %q", runtime.GOOS)
}
