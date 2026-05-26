//go:build darwin

package loader

import "syscore/window/appkit/bindings"

/*
AppkitCommands holds bound Objective-C runtime entry points for AppKit window setup.

[Context]
Populate with AppkitCommandsLoad or AppkitCommandsLoadManifest. NSWindow and NSApplication
interaction uses ObjcMsgSend with classes and selectors from bindings constants and SYSCORE_C_StringToCString.
*/
type AppkitCommands struct {
	/*
		ObjcGetClass is objc_getClass. Resolves "NSApplication", "NSWindow", etc.
	*/
	ObjcGetClass bindings.PFN_objc_getClass
	/*
		SelRegisterName is sel_registerName. Registers selector strings before objc_msgSend.
	*/
	SelRegisterName bindings.PFN_sel_registerName
	/*
		ObjcMsgSend is objc_msgSend. Dispatches Objective-C messages (alloc, init, makeKeyAndOrderFront:, ...).
	*/
	ObjcMsgSend bindings.PFN_objc_msgSend
	/*
		ObjcAllocateClassPair is objc_allocateClassPair. Creates a dynamic Objective-C subclass.
	*/
	ObjcAllocateClassPair bindings.PFN_objc_allocateClassPair
	/*
		ObjcRegisterClassPair is objc_registerClassPair. Registers a class pair after adding methods.
	*/
	ObjcRegisterClassPair bindings.PFN_objc_registerClassPair
	/*
		ClassAddMethod is class_addMethod. Adds an instance method to a dynamic class.
	*/
	ClassAddMethod bindings.PFN_class_addMethod
}
