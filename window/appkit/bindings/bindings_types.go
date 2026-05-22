package bindings

/*
ObjCObject is an opaque Objective-C id (instance pointer).
*/
type ObjCObject uintptr

/*
ObjCClass is an opaque Objective-C Class pointer.
*/
type ObjCClass uintptr

/*
ObjCSel is an opaque SEL method selector pointer.
*/
type ObjCSel uintptr

/*
NSRect is the Core Graphics / AppKit NSRect layout (origin x,y and size width,height as float64).

Passed by value to objc_msgSend for initWithContentRect:styleMask:backing:defer: on amd64.
*/
type NSRect struct {
	X      float64
	Y      float64
	Width  float64
	Height float64
}

/*
PFN_objc_getClass returns the Objective-C class for a given name.

[Context]
Maps to objc_getClass. Pass a NUL-terminated UTF-8 class name (for example "NSWindow").

[Parameters]
name - Class name as *byte.

[Returns]
Class pointer as ObjCClass.
*/
type PFN_objc_getClass func(name *byte) ObjCClass

/*
PFN_sel_registerName registers a selector name and returns its SEL.

[Context]
Maps to sel_registerName. Selector strings use Objective-C naming (for example "alloc").

[Parameters]
name - Selector name as *byte.

[Returns]
SEL handle.
*/
type PFN_sel_registerName func(name *byte) ObjCSel

/*
PFN_objc_msgSend sends an Objective-C message to receiver.

[Context]
Maps to objc_msgSend. Variadic args are expanded by purego for additional parameters after the
selector (structs, integers, object pointers). Floating-point returns may require objc_msgSend_fpret
for some selectors; NSWindow init paths used here work with objc_msgSend on amd64.

[Parameters]
receiver - Target object or class (as ObjCObject).
selector - Method SEL.
args - Method arguments per ABI.

[Returns]
Return value as ObjCObject (0 indicates failure for init methods).
*/
type PFN_objc_msgSend func(receiver ObjCObject, selector ObjCSel, args ...any) ObjCObject
