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
