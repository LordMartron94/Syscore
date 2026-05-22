//go:build windows

package internal

import "golang.org/x/sys/windows"

func loadPlatformLibrary(name string) (uintptr, error) {
	handle, err := windows.LoadLibrary(name)
	if err != nil {
		return 0, err
	}
	return uintptr(handle), nil
}

func platformLibrarySymbolResolve(library DynamicLibrary, symbolName string) (uintptr, error) {
	proc, err := windows.GetProcAddress(windows.Handle(library), symbolName)
	if err != nil {
		return 0, err
	}
	return proc, nil
}
