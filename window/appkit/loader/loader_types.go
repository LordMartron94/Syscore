package loader

import "syscore/window/appkit/bindings"

type AppkitCommands struct {
	ObjcGetClass    bindings.PFN_objc_getClass
	SelRegisterName bindings.PFN_sel_registerName
	ObjcMsgSend     bindings.PFN_objc_msgSend
}
