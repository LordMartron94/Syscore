//go:build windows

package loader

import (
	"fmt"

	"syscore/window/loadutil"
)

func Win32CommandsLoad(module Win32Module, commands *Win32Commands) error {
	user32Mappings, kernel32Mappings := win32CommandMappingsAll(commands)
	return win32CommandsBind(module, commands, user32Mappings, kernel32Mappings)
}

func win32CommandMappingsAll(commands *Win32Commands) (user32 []loadutil.CommandMapping, kernel32 []loadutil.CommandMapping) {
	user32 = []loadutil.CommandMapping{
		{Target: &commands.RegisterClassExW, Name: "RegisterClassExW"},
		{Target: &commands.CreateWindowExW, Name: "CreateWindowExW"},
		{Target: &commands.DefWindowProcW, Name: "DefWindowProcW"},
		{Target: &commands.ShowWindow, Name: "ShowWindow"},
		{Target: &commands.UpdateWindow, Name: "UpdateWindow"},
		{Target: &commands.LoadCursorW, Name: "LoadCursorW"},
		{Target: &commands.DestroyWindow, Name: "DestroyWindow"},
		{Target: &commands.PeekMessageW, Name: "PeekMessageW"},
		{Target: &commands.DispatchMessageW, Name: "DispatchMessageW"},
	}
	kernel32 = []loadutil.CommandMapping{
		{Target: &commands.GetModuleHandleW, Name: "GetModuleHandleW"},
	}
	return user32, kernel32
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
