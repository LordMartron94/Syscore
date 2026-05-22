//go:build !windows

package loader

import "fmt"

func Win32CommandsLoadManifest(_ Win32Module, _ *Win32CommandManifest, _ *Win32Commands) error {
	return fmt.Errorf("win32 loader: only supported on windows")
}
