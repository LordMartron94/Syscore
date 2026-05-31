//go:build windows

package loader

import (
	"fmt"

	kernel32 "syscore/window/win32/bindings/pfn/kernel32_dll"
	user32 "syscore/window/win32/bindings/pfn/user32_dll"
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
	Win32CommandManifestFieldDestroyWindow
	Win32CommandManifestFieldGetClientRect
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

func Win32CommandManifestAddGetModuleHandleW(manifest *Win32CommandManifest, target *kernel32.PFN_GetModuleHandleW) error {
	return win32CommandManifestAdd(manifest, Win32CommandManifestFieldGetModuleHandleW, target)
}

func Win32CommandManifestAddRegisterClassExW(manifest *Win32CommandManifest, target *user32.PFN_RegisterClassExW) error {
	return win32CommandManifestAdd(manifest, Win32CommandManifestFieldRegisterClassExW, target)
}

func Win32CommandManifestAddCreateWindowExW(manifest *Win32CommandManifest, target *user32.PFN_CreateWindowExW) error {
	return win32CommandManifestAdd(manifest, Win32CommandManifestFieldCreateWindowExW, target)
}

func Win32CommandManifestAddDefWindowProcW(manifest *Win32CommandManifest, target *user32.PFN_DefWindowProcW) error {
	return win32CommandManifestAdd(manifest, Win32CommandManifestFieldDefWindowProcW, target)
}

func Win32CommandManifestAddShowWindow(manifest *Win32CommandManifest, target *user32.PFN_ShowWindow) error {
	return win32CommandManifestAdd(manifest, Win32CommandManifestFieldShowWindow, target)
}

func Win32CommandManifestAddUpdateWindow(manifest *Win32CommandManifest, target *user32.PFN_UpdateWindow) error {
	return win32CommandManifestAdd(manifest, Win32CommandManifestFieldUpdateWindow, target)
}

func Win32CommandManifestAddLoadCursorW(manifest *Win32CommandManifest, target *user32.PFN_LoadCursorW) error {
	return win32CommandManifestAdd(manifest, Win32CommandManifestFieldLoadCursorW, target)
}

func Win32CommandManifestAddDestroyWindow(manifest *Win32CommandManifest, target *user32.PFN_DestroyWindow) error {
	return win32CommandManifestAdd(manifest, Win32CommandManifestFieldDestroyWindow, target)
}

func Win32CommandManifestAddGetClientRect(manifest *Win32CommandManifest, target *user32.PFN_GetClientRect) error {
	return win32CommandManifestAdd(manifest, Win32CommandManifestFieldGetClientRect, target)
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
