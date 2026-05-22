//go:build !windows

package loader

import "fmt"

func Win32CommandsLoad(_ Win32Module, _ *Win32Commands) error {
	return fmt.Errorf("win32 loader: only supported on windows")
}
