//go:build !windows

package loader

import "fmt"

type Win32Module struct{}

func Win32ModuleLoad() (Win32Module, error) {
	return Win32Module{}, fmt.Errorf("win32 loader: only supported on windows")
}

func Win32ModuleUser32LibraryName() string   { return "" }
func Win32ModuleKernel32LibraryName() string { return "" }

func Win32ModuleUser32SymbolResolve(_ Win32Module, _ string) (uintptr, error) {
	return 0, fmt.Errorf("win32 loader: only supported on windows")
}

func Win32ModuleKernel32SymbolResolve(_ Win32Module, _ string) (uintptr, error) {
	return 0, fmt.Errorf("win32 loader: only supported on windows")
}
