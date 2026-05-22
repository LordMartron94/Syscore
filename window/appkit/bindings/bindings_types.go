package bindings

type ObjCObject uintptr
type ObjCClass uintptr
type ObjCSel uintptr

type NSRect struct {
	X      float64
	Y      float64
	Width  float64
	Height float64
}

type PFN_objc_getClass func(name *byte) ObjCClass

type PFN_sel_registerName func(name *byte) ObjCSel

type PFN_objc_msgSend func(receiver ObjCObject, selector ObjCSel, args ...any) ObjCObject
