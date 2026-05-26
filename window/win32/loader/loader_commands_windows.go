//go:build windows

package loader

import (
	"fmt"

	"syscore/window/loadutil"
)

func Win32CommandsLoad(module Win32Module, commands *Win32Commands) error {
	user32Mappings := win32User32CommandMappings(commands)
	kernel32Mappings := win32Kernel32CommandMappings(commands)
	return win32CommandsBind(module, commands, user32Mappings, kernel32Mappings)
}

func win32CommandsBind(module Win32Module, commands *Win32Commands, user32Mappings, kernel32Mappings []loadutil.CommandMapping) error {
	if commands == nil {
		return fmt.Errorf("win32 loader: commands must not be nil")
	}
	if err := loadutil.LibraryCommandsBind(module.User32, user32Mappings); err != nil {
		return fmt.Errorf("win32 loader: user32: %w", err)
	}
	if err := loadutil.LibraryCommandsBind(module.Kernel32, kernel32Mappings); err != nil {
		return fmt.Errorf("win32 loader: kernel32: %w", err)
	}
	return nil
}
