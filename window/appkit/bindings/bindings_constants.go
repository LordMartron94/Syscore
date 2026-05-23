package bindings

const (
	// NSWindowStyleMaskTitled includes a title bar.
	NSWindowStyleMaskTitled uint64 = 1 << 0
	// NSWindowStyleMaskClosable includes a close button.
	NSWindowStyleMaskClosable uint64 = 1 << 1
	// NSWindowStyleMaskMiniaturizable includes a minimize button.
	NSWindowStyleMaskMiniaturizable uint64 = 1 << 2
	// NSWindowStyleMaskResizable allows user resizing.
	NSWindowStyleMaskResizable uint64 = 1 << 3
	// NSBackingStoreBuffered selects buffered backing for the window.
	NSBackingStoreBuffered uint64 = 2
)

const (
	// ObjCClassNSApplication is the UTF-8 class name for NSApplication.
	ObjCClassNSApplication = "NSApplication"
	// ObjCClassNSWindow is the UTF-8 class name for NSWindow.
	ObjCClassNSWindow = "NSWindow"
	// ObjCClassNSString is the UTF-8 class name for NSString (Foundation).
	ObjCClassNSString = "NSString"
)

const (
	// ObjCSelSharedApplication is the selector for +[NSApplication sharedApplication].
	ObjCSelSharedApplication = "sharedApplication"
	// ObjCSelAlloc is the selector for +alloc / instance allocation.
	ObjCSelAlloc = "alloc"
	// ObjCSelInitWithContentRectStyleMaskBackingDefer is the NSWindow designated initializer selector.
	ObjCSelInitWithContentRectStyleMaskBackingDefer = "initWithContentRect:styleMask:backing:defer:"
	// ObjCSelMakeKeyAndOrderFront is the selector to show and key a window.
	ObjCSelMakeKeyAndOrderFront = "makeKeyAndOrderFront:"
	// ObjCSelStringWithUTF8String is the NSString class method +stringWithUTF8String:.
	ObjCSelStringWithUTF8String = "stringWithUTF8String:"
	// ObjCSelSetTitle is the NSWindow instance method -setTitle:.
	ObjCSelSetTitle = "setTitle:"
	// ObjCSelClose is the NSWindow instance method -close.
	ObjCSelClose = "close"
)
