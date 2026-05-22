//go:build windows

package loader

import (
	"fmt"

	"syscore/window/loadutil"
)

func Win32CommandsLoadManifest(module Win32Module, manifest *Win32CommandManifest, commands *Win32Commands) error {
	if manifest == nil {
		return fmt.Errorf("win32 loader: manifest must not be nil")
	}
	if commands == nil {
		return fmt.Errorf("win32 loader: commands must not be nil")
	}

	user32Mappings := make([]loadutil.CommandMapping, 0, len(manifest.fields))
	kernel32Mappings := make([]loadutil.CommandMapping, 0, len(manifest.fields))
	for _, field := range manifest.fields {
		switch field {
		case Win32CommandManifestFieldGetModuleHandleW:
			kernel32Mappings = append(kernel32Mappings, loadutil.CommandMapping{Target: &commands.GetModuleHandleW, Name: "GetModuleHandleW"})
		case Win32CommandManifestFieldRegisterClassExW:
			user32Mappings = append(user32Mappings, loadutil.CommandMapping{Target: &commands.RegisterClassExW, Name: "RegisterClassExW"})
		case Win32CommandManifestFieldCreateWindowExW:
			user32Mappings = append(user32Mappings, loadutil.CommandMapping{Target: &commands.CreateWindowExW, Name: "CreateWindowExW"})
		case Win32CommandManifestFieldDefWindowProcW:
			user32Mappings = append(user32Mappings, loadutil.CommandMapping{Target: &commands.DefWindowProcW, Name: "DefWindowProcW"})
		case Win32CommandManifestFieldShowWindow:
			user32Mappings = append(user32Mappings, loadutil.CommandMapping{Target: &commands.ShowWindow, Name: "ShowWindow"})
		case Win32CommandManifestFieldUpdateWindow:
			user32Mappings = append(user32Mappings, loadutil.CommandMapping{Target: &commands.UpdateWindow, Name: "UpdateWindow"})
		case Win32CommandManifestFieldLoadCursorW:
			user32Mappings = append(user32Mappings, loadutil.CommandMapping{Target: &commands.LoadCursorW, Name: "LoadCursorW"})
		default:
			return fmt.Errorf("win32 loader: unknown manifest field %d", field)
		}
	}

	return win32CommandsBind(module, commands, user32Mappings, kernel32Mappings)
}
