package syscore

import (
	"syscore/window/win32/bindings"
	"syscore/window/win32/loader"
)

type SYSCORE_Window_Win32_Module = loader.Win32Module
type SYSCORE_Window_Win32_Commands = loader.Win32Commands
type SYSCORE_Window_Win32_CommandManifest = loader.Win32CommandManifest

type (
	SYSCORE_Window_Win32_HWND                 = bindings.HWND
	SYSCORE_Window_Win32_HINSTANCE            = bindings.HINSTANCE
	SYSCORE_Window_Win32_WNDCLASSEXW          = bindings.WNDCLASSEXW
	SYSCORE_Window_Win32_PFN_GetModuleHandleW = bindings.PFN_GetModuleHandleW
	SYSCORE_Window_Win32_PFN_RegisterClassExW = bindings.PFN_RegisterClassExW
	SYSCORE_Window_Win32_PFN_CreateWindowExW  = bindings.PFN_CreateWindowExW
	SYSCORE_Window_Win32_PFN_DefWindowProcW   = bindings.PFN_DefWindowProcW
	SYSCORE_Window_Win32_PFN_ShowWindow       = bindings.PFN_ShowWindow
	SYSCORE_Window_Win32_PFN_UpdateWindow     = bindings.PFN_UpdateWindow
	SYSCORE_Window_Win32_PFN_LoadCursorW      = bindings.PFN_LoadCursorW
)

func SYSCORE_Window_Win32_ModuleLoad() (SYSCORE_Window_Win32_Module, error) {
	return loader.Win32ModuleLoad()
}

func SYSCORE_Window_Win32_ModuleUser32LibraryName() string {
	return loader.Win32ModuleUser32LibraryName()
}

func SYSCORE_Window_Win32_ModuleKernel32LibraryName() string {
	return loader.Win32ModuleKernel32LibraryName()
}

func SYSCORE_Window_Win32_ModuleUser32SymbolResolve(module SYSCORE_Window_Win32_Module, symbolName string) (uintptr, error) {
	return loader.Win32ModuleUser32SymbolResolve(module, symbolName)
}

func SYSCORE_Window_Win32_ModuleKernel32SymbolResolve(module SYSCORE_Window_Win32_Module, symbolName string) (uintptr, error) {
	return loader.Win32ModuleKernel32SymbolResolve(module, symbolName)
}

func SYSCORE_Window_Win32_CommandsLoad(module SYSCORE_Window_Win32_Module, commands *SYSCORE_Window_Win32_Commands) error {
	return loader.Win32CommandsLoad(module, commands)
}

func SYSCORE_Window_Win32_CommandManifestReset(manifest *SYSCORE_Window_Win32_CommandManifest) {
	loader.Win32CommandManifestReset(manifest)
}

func SYSCORE_Window_Win32_CommandManifestAddGetModuleHandleW(manifest *SYSCORE_Window_Win32_CommandManifest, target *SYSCORE_Window_Win32_PFN_GetModuleHandleW) error {
	return loader.Win32CommandManifestAddGetModuleHandleW(manifest, target)
}

func SYSCORE_Window_Win32_CommandManifestAddRegisterClassExW(manifest *SYSCORE_Window_Win32_CommandManifest, target *SYSCORE_Window_Win32_PFN_RegisterClassExW) error {
	return loader.Win32CommandManifestAddRegisterClassExW(manifest, target)
}

func SYSCORE_Window_Win32_CommandManifestAddCreateWindowExW(manifest *SYSCORE_Window_Win32_CommandManifest, target *SYSCORE_Window_Win32_PFN_CreateWindowExW) error {
	return loader.Win32CommandManifestAddCreateWindowExW(manifest, target)
}

func SYSCORE_Window_Win32_CommandManifestAddDefWindowProcW(manifest *SYSCORE_Window_Win32_CommandManifest, target *SYSCORE_Window_Win32_PFN_DefWindowProcW) error {
	return loader.Win32CommandManifestAddDefWindowProcW(manifest, target)
}

func SYSCORE_Window_Win32_CommandManifestAddShowWindow(manifest *SYSCORE_Window_Win32_CommandManifest, target *SYSCORE_Window_Win32_PFN_ShowWindow) error {
	return loader.Win32CommandManifestAddShowWindow(manifest, target)
}

func SYSCORE_Window_Win32_CommandManifestAddUpdateWindow(manifest *SYSCORE_Window_Win32_CommandManifest, target *SYSCORE_Window_Win32_PFN_UpdateWindow) error {
	return loader.Win32CommandManifestAddUpdateWindow(manifest, target)
}

func SYSCORE_Window_Win32_CommandManifestAddLoadCursorW(manifest *SYSCORE_Window_Win32_CommandManifest, target *SYSCORE_Window_Win32_PFN_LoadCursorW) error {
	return loader.Win32CommandManifestAddLoadCursorW(manifest, target)
}

func SYSCORE_Window_Win32_CommandsLoadManifest(module SYSCORE_Window_Win32_Module, manifest *SYSCORE_Window_Win32_CommandManifest, commands *SYSCORE_Window_Win32_Commands) error {
	return loader.Win32CommandsLoadManifest(module, manifest, commands)
}
