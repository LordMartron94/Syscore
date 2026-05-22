package syscore

import (
	"syscore/window/appkit/bindings"
	"syscore/window/appkit/loader"
)

type SYSCORE_Window_Appkit_Module = loader.AppkitModule
type SYSCORE_Window_Appkit_Commands = loader.AppkitCommands
type SYSCORE_Window_Appkit_CommandManifest = loader.AppkitCommandManifest

type (
	SYSCORE_Window_Appkit_Object               = bindings.ObjCObject
	SYSCORE_Window_Appkit_Class                = bindings.ObjCClass
	SYSCORE_Window_Appkit_Sel                  = bindings.ObjCSel
	SYSCORE_Window_Appkit_NSRect               = bindings.NSRect
	SYSCORE_Window_Appkit_PFN_objc_getClass    = bindings.PFN_objc_getClass
	SYSCORE_Window_Appkit_PFN_sel_registerName = bindings.PFN_sel_registerName
	SYSCORE_Window_Appkit_PFN_objc_msgSend     = bindings.PFN_objc_msgSend
)

func SYSCORE_Window_Appkit_ModuleLoad() (SYSCORE_Window_Appkit_Module, error) {
	return loader.AppkitModuleLoad()
}

func SYSCORE_Window_Appkit_ModuleObjcLibraryName() string {
	return loader.AppkitModuleObjcLibraryName()
}

func SYSCORE_Window_Appkit_ModuleAppkitLibraryName() string {
	return loader.AppkitModuleAppkitLibraryName()
}

func SYSCORE_Window_Appkit_ModuleObjcSymbolResolve(module SYSCORE_Window_Appkit_Module, symbolName string) (uintptr, error) {
	return loader.AppkitModuleObjcSymbolResolve(module, symbolName)
}

func SYSCORE_Window_Appkit_ModuleAppkitSymbolResolve(module SYSCORE_Window_Appkit_Module, symbolName string) (uintptr, error) {
	return loader.AppkitModuleAppkitSymbolResolve(module, symbolName)
}

func SYSCORE_Window_Appkit_CommandsLoad(module SYSCORE_Window_Appkit_Module, commands *SYSCORE_Window_Appkit_Commands) error {
	return loader.AppkitCommandsLoad(module, commands)
}

func SYSCORE_Window_Appkit_CommandManifestReset(manifest *SYSCORE_Window_Appkit_CommandManifest) {
	loader.AppkitCommandManifestReset(manifest)
}

func SYSCORE_Window_Appkit_CommandManifestAddObjcGetClass(manifest *SYSCORE_Window_Appkit_CommandManifest, target *SYSCORE_Window_Appkit_PFN_objc_getClass) error {
	return loader.AppkitCommandManifestAddObjcGetClass(manifest, target)
}

func SYSCORE_Window_Appkit_CommandManifestAddSelRegisterName(manifest *SYSCORE_Window_Appkit_CommandManifest, target *SYSCORE_Window_Appkit_PFN_sel_registerName) error {
	return loader.AppkitCommandManifestAddSelRegisterName(manifest, target)
}

func SYSCORE_Window_Appkit_CommandManifestAddObjcMsgSend(manifest *SYSCORE_Window_Appkit_CommandManifest, target *SYSCORE_Window_Appkit_PFN_objc_msgSend) error {
	return loader.AppkitCommandManifestAddObjcMsgSend(manifest, target)
}

func SYSCORE_Window_Appkit_CommandsLoadManifest(module SYSCORE_Window_Appkit_Module, manifest *SYSCORE_Window_Appkit_CommandManifest, commands *SYSCORE_Window_Appkit_Commands) error {
	return loader.AppkitCommandsLoadManifest(module, manifest, commands)
}
