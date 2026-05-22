package bindings

const (
	NSWindowStyleMaskTitled         uint64 = 1 << 0
	NSWindowStyleMaskClosable       uint64 = 1 << 1
	NSWindowStyleMaskMiniaturizable uint64 = 1 << 2
	NSWindowStyleMaskResizable      uint64 = 1 << 3
	NSBackingStoreBuffered          uint64 = 2
)

const (
	ObjCClassNSApplication = "NSApplication"
	ObjCClassNSWindow      = "NSWindow"
)

const (
	ObjCSelSharedApplication                        = "sharedApplication"
	ObjCSelAlloc                                    = "alloc"
	ObjCSelInitWithContentRectStyleMaskBackingDefer = "initWithContentRect:styleMask:backing:defer:"
	ObjCSelMakeKeyAndOrderFront                     = "makeKeyAndOrderFront:"
)
