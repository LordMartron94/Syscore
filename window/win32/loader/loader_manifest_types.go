package loader

import (
	"fmt"

	"syscore/window/win32/bindings"
)

type Win32CommandManifestField uint8

const (
	Win32CommandManifestFieldInvalid Win32CommandManifestField = iota
	Win32CommandManifestFieldGetModuleHandleW
	Win32CommandManifestFieldRegisterClassExW
	Win32CommandManifestFieldCreateWindowExW
	Win32CommandManifestFieldDefWindowProcW
	Win32CommandManifestFieldShowWindow
	Win32CommandManifestFieldUpdateWindow
	Win32CommandManifestFieldLoadCursorW
)

type Win32CommandManifest struct {
	fields []Win32CommandManifestField
}

func Win32CommandManifestReset(manifest *Win32CommandManifest) {
	if manifest == nil {
		return
	}
	manifest.fields = manifest.fields[:0]
}

func Win32CommandManifestAddGetModuleHandleW(manifest *Win32CommandManifest, target *bindings.PFN_GetModuleHandleW) error {
	return win32CommandManifestAdd(manifest, Win32CommandManifestFieldGetModuleHandleW, target)
}

func Win32CommandManifestAddRegisterClassExW(manifest *Win32CommandManifest, target *bindings.PFN_RegisterClassExW) error {
	return win32CommandManifestAdd(manifest, Win32CommandManifestFieldRegisterClassExW, target)
}

func Win32CommandManifestAddCreateWindowExW(manifest *Win32CommandManifest, target *bindings.PFN_CreateWindowExW) error {
	return win32CommandManifestAdd(manifest, Win32CommandManifestFieldCreateWindowExW, target)
}

func Win32CommandManifestAddDefWindowProcW(manifest *Win32CommandManifest, target *bindings.PFN_DefWindowProcW) error {
	return win32CommandManifestAdd(manifest, Win32CommandManifestFieldDefWindowProcW, target)
}

func Win32CommandManifestAddShowWindow(manifest *Win32CommandManifest, target *bindings.PFN_ShowWindow) error {
	return win32CommandManifestAdd(manifest, Win32CommandManifestFieldShowWindow, target)
}

func Win32CommandManifestAddUpdateWindow(manifest *Win32CommandManifest, target *bindings.PFN_UpdateWindow) error {
	return win32CommandManifestAdd(manifest, Win32CommandManifestFieldUpdateWindow, target)
}

func Win32CommandManifestAddLoadCursorW(manifest *Win32CommandManifest, target *bindings.PFN_LoadCursorW) error {
	return win32CommandManifestAdd(manifest, Win32CommandManifestFieldLoadCursorW, target)
}

func win32CommandManifestAdd(manifest *Win32CommandManifest, field Win32CommandManifestField, target any) error {
	if manifest == nil {
		return fmt.Errorf("win32 loader: manifest must not be nil")
	}
	if target == nil {
		return fmt.Errorf("win32 loader: manifest target for field %d must not be nil", field)
	}
	manifest.fields = append(manifest.fields, field)
	return nil
}
