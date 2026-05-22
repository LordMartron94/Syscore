//go:build windows

package loader

import (
	"fmt"

	"syscore/internal"
	"syscore/window/loadutil"
)

const (
	win32User32LibraryName   = "user32.dll"
	win32Kernel32LibraryName = "kernel32.dll"
)

type Win32Module struct {
	User32   internal.DynamicLibrary
	Kernel32 internal.DynamicLibrary
}

func Win32ModuleLoad() (Win32Module, error) {
	user32, err := internal.DynamicLibraryLoad(win32User32LibraryName)
	if err != nil {
		return Win32Module{}, fmt.Errorf("win32 loader: user32: %w", err)
	}
	kernel32, err := internal.DynamicLibraryLoad(win32Kernel32LibraryName)
	if err != nil {
		return Win32Module{}, fmt.Errorf("win32 loader: kernel32: %w", err)
	}
	return Win32Module{User32: user32, Kernel32: kernel32}, nil
}

func Win32ModuleUser32LibraryName() string   { return win32User32LibraryName }
func Win32ModuleKernel32LibraryName() string { return win32Kernel32LibraryName }

func Win32ModuleUser32SymbolResolve(module Win32Module, symbolName string) (uintptr, error) {
	return loadutil.LibrarySymbolResolve(module.User32, symbolName)
}

func Win32ModuleKernel32SymbolResolve(module Win32Module, symbolName string) (uintptr, error) {
	return loadutil.LibrarySymbolResolve(module.Kernel32, symbolName)
}
