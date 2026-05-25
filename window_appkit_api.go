//go:build darwin

package syscore

import (
	"syscore/window/appkit/bindings"
	"syscore/window/appkit/loader"
)

/*
SYSCORE_Window_Appkit_Module holds loaded libobjc.A.dylib and AppKit framework handles.

[Context]
Obtain via SYSCORE_Window_Appkit_ModuleLoad on macOS. Window creation uses objc_msgSend through
bound commands after resolving classes and selectors with SYSCORE_C_StringToCString.
*/
type SYSCORE_Window_Appkit_Module = loader.AppkitModule

/*
SYSCORE_Window_Appkit_Commands holds bound Objective-C runtime entry points used for AppKit windows.

[Context]
Call through documented fields on loader.AppkitCommands (ObjcGetClass, SelRegisterName, ObjcMsgSend).
NSWindow setup uses bindings constants for class/selector names and ObjcMsgSend for dispatch.
*/
type SYSCORE_Window_Appkit_Commands = loader.AppkitCommands

/*
SYSCORE_Window_Appkit_CommandManifest lists Objective-C runtime entry points for selective binding.
*/
type SYSCORE_Window_Appkit_CommandManifest = loader.AppkitCommandManifest

const (
	SYSCORE_Window_Appkit_WindowStyleTitled             = bindings.NSWindowStyleMaskTitled
	SYSCORE_Window_Appkit_WindowStyleClosable           = bindings.NSWindowStyleMaskClosable
	SYSCORE_Window_Appkit_WindowStyleMiniaturizable     = bindings.NSWindowStyleMaskMiniaturizable
	SYSCORE_Window_Appkit_WindowStyleResizable          = bindings.NSWindowStyleMaskResizable
	SYSCORE_Window_Appkit_BackingStoreBuffered          = bindings.NSBackingStoreBuffered
	SYSCORE_Window_Appkit_ClassNSApplication            = bindings.ObjCClassNSApplication
	SYSCORE_Window_Appkit_ClassNSWindow                 = bindings.ObjCClassNSWindow
	SYSCORE_Window_Appkit_ClassNSString                 = bindings.ObjCClassNSString
	SYSCORE_Window_Appkit_SelSharedApplication          = bindings.ObjCSelSharedApplication
	SYSCORE_Window_Appkit_SelAlloc                      = bindings.ObjCSelAlloc
	SYSCORE_Window_Appkit_SelInitWindow                 = bindings.ObjCSelInitWithContentRectStyleMaskBackingDefer
	SYSCORE_Window_Appkit_SelMakeKeyAndOrderFront       = bindings.ObjCSelMakeKeyAndOrderFront
	SYSCORE_Window_Appkit_SelStringWithUTF8String       = bindings.ObjCSelStringWithUTF8String
	SYSCORE_Window_Appkit_SelSetTitle                   = bindings.ObjCSelSetTitle
	SYSCORE_Window_Appkit_SelClose                      = bindings.ObjCSelClose
	SYSCORE_Window_Appkit_SelInit                       = bindings.ObjCSelInit
	SYSCORE_Window_Appkit_SelSetDelegate                = bindings.ObjCSelSetDelegate
	SYSCORE_Window_Appkit_SelWindowShouldClose          = bindings.ObjCSelWindowShouldClose
	SYSCORE_Window_Appkit_SelNextEventMatchingMask      = bindings.ObjCSelNextEventMatchingMask
	SYSCORE_Window_Appkit_ClassNSDate                   = bindings.ObjCClassNSDate
	SYSCORE_Window_Appkit_SelDistantPast                = bindings.ObjCSelDistantPast
	SYSCORE_Window_Appkit_RunLoopModeDefault            = bindings.NSRunLoopModeDefault
	SYSCORE_Window_Appkit_EventMaskAny                  = bindings.NSEventMaskAny
	SYSCORE_Window_Appkit_TypeEncodingWindowShouldClose = bindings.ObjCTypeEncodingWindowShouldClose
)

type (
	// SYSCORE_Window_Appkit_Object is an opaque Objective-C object (id).
	SYSCORE_Window_Appkit_Object = bindings.ObjCObject
	// SYSCORE_Window_Appkit_Class is an opaque Objective-C Class pointer.
	SYSCORE_Window_Appkit_Class = bindings.ObjCClass
	// SYSCORE_Window_Appkit_Sel is an opaque SEL selector pointer.
	SYSCORE_Window_Appkit_Sel = bindings.ObjCSel
	// SYSCORE_Window_Appkit_NSRect is the NSRect layout (origin and size as float64 values).
	SYSCORE_Window_Appkit_NSRect = bindings.NSRect
	// SYSCORE_Window_Appkit_PFN_objc_getClass is the C type for objc_getClass.
	SYSCORE_Window_Appkit_PFN_objc_getClass = bindings.PFN_objc_getClass
	// SYSCORE_Window_Appkit_PFN_sel_registerName is the C type for sel_registerName.
	SYSCORE_Window_Appkit_PFN_sel_registerName = bindings.PFN_sel_registerName
	// SYSCORE_Window_Appkit_PFN_objc_msgSend is the C type for objc_msgSend.
	SYSCORE_Window_Appkit_PFN_objc_msgSend = bindings.PFN_objc_msgSend
	// SYSCORE_Window_Appkit_PFN_objc_allocateClassPair is the C type for objc_allocateClassPair.
	SYSCORE_Window_Appkit_PFN_objc_allocateClassPair = bindings.PFN_objc_allocateClassPair
	// SYSCORE_Window_Appkit_PFN_objc_registerClassPair is the C type for objc_registerClassPair.
	SYSCORE_Window_Appkit_PFN_objc_registerClassPair = bindings.PFN_objc_registerClassPair
	// SYSCORE_Window_Appkit_PFN_class_addMethod is the C type for class_addMethod.
	SYSCORE_Window_Appkit_PFN_class_addMethod = bindings.PFN_class_addMethod
)

/*
SYSCORE_Window_Appkit_ModuleLoad loads libobjc.A.dylib and the AppKit framework.

[Context]
Required on macOS for native windows. objc symbols bind from the Objc handle; Appkit handle is
available for additional symbol resolution if needed.

[Returns]
Module with Objc and Appkit library handles.

[Errors]
Returns an error when not on darwin or dlopen fails for either library.

[Side Effects]
Loads Objective-C runtime and AppKit into the process.
*/
func SYSCORE_Window_Appkit_ModuleLoad() (SYSCORE_Window_Appkit_Module, error) {
	return loader.AppkitModuleLoad()
}

/*
SYSCORE_Window_Appkit_ModuleObjcLibraryName returns the path used to load the Objective-C runtime.
*/
func SYSCORE_Window_Appkit_ModuleObjcLibraryName() string {
	return loader.AppkitModuleObjcLibraryName()
}

/*
SYSCORE_Window_Appkit_ModuleAppkitLibraryName returns the path used to load the AppKit framework.
*/
func SYSCORE_Window_Appkit_ModuleAppkitLibraryName() string {
	return loader.AppkitModuleAppkitLibraryName()
}

/*
SYSCORE_Window_Appkit_ModuleObjcSymbolResolve resolves a symbol from libobjc.A.dylib.
*/
func SYSCORE_Window_Appkit_ModuleObjcSymbolResolve(module SYSCORE_Window_Appkit_Module, symbolName string) (uintptr, error) {
	return loader.AppkitModuleObjcSymbolResolve(module, symbolName)
}

/*
SYSCORE_Window_Appkit_ModuleAppkitSymbolResolve resolves a symbol from the AppKit framework image.
*/
func SYSCORE_Window_Appkit_ModuleAppkitSymbolResolve(module SYSCORE_Window_Appkit_Module, symbolName string) (uintptr, error) {
	return loader.AppkitModuleAppkitSymbolResolve(module, symbolName)
}

/*
SYSCORE_Window_Appkit_CommandsLoad binds objc_getClass, sel_registerName, and objc_msgSend.

[Parameters]
module - Loaded Appkit module.
commands - Non-nil command holder.

[Errors]
Returns an error if commands is nil, not on darwin, or any bind fails.
*/
func SYSCORE_Window_Appkit_CommandsLoad(module SYSCORE_Window_Appkit_Module, commands *SYSCORE_Window_Appkit_Commands) error {
	return loader.AppkitCommandsLoad(module, commands)
}

/*
SYSCORE_Window_Appkit_CommandManifestReset clears manifest.
*/
func SYSCORE_Window_Appkit_CommandManifestReset(manifest *SYSCORE_Window_Appkit_CommandManifest) {
	loader.AppkitCommandManifestReset(manifest)
}

/*
SYSCORE_Window_Appkit_CommandManifestAddObjcGetClass registers objc_getClass for selective loading.
*/
func SYSCORE_Window_Appkit_CommandManifestAddObjcGetClass(manifest *SYSCORE_Window_Appkit_CommandManifest, target *SYSCORE_Window_Appkit_PFN_objc_getClass) error {
	return loader.AppkitCommandManifestAddObjcGetClass(manifest, target)
}

/*
SYSCORE_Window_Appkit_CommandManifestAddSelRegisterName registers sel_registerName for selective loading.
*/
func SYSCORE_Window_Appkit_CommandManifestAddSelRegisterName(manifest *SYSCORE_Window_Appkit_CommandManifest, target *SYSCORE_Window_Appkit_PFN_sel_registerName) error {
	return loader.AppkitCommandManifestAddSelRegisterName(manifest, target)
}

/*
SYSCORE_Window_Appkit_CommandManifestAddObjcMsgSend registers objc_msgSend for selective loading.
*/
func SYSCORE_Window_Appkit_CommandManifestAddObjcMsgSend(manifest *SYSCORE_Window_Appkit_CommandManifest, target *SYSCORE_Window_Appkit_PFN_objc_msgSend) error {
	return loader.AppkitCommandManifestAddObjcMsgSend(manifest, target)
}

/*
SYSCORE_Window_Appkit_CommandsLoadManifest binds only manifest-listed Objective-C runtime entry points.

[Errors]
Returns an error if manifest or commands is nil or any bind fails.
*/
func SYSCORE_Window_Appkit_CommandsLoadManifest(module SYSCORE_Window_Appkit_Module, manifest *SYSCORE_Window_Appkit_CommandManifest, commands *SYSCORE_Window_Appkit_Commands) error {
	return loader.AppkitCommandsLoadManifest(module, manifest, commands)
}
