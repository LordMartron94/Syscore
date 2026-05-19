//go:build windows

package internal

import "golang.org/x/sys/windows"

func loadPlatformLibrary(name string) (DynamicLibrary, error) {
	handle, err := windows.LoadLibrary(name)
	if err != nil {
		return 0, err
	}
	return DynamicLibrary(handle), nil
}
