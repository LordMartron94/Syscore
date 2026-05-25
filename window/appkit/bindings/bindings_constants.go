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
	// ObjCSelInit is the NSObject instance method -init.
	ObjCSelInit = "init"
	// ObjCSelSetDelegate is the NSWindow instance method -setDelegate:.
	ObjCSelSetDelegate = "setDelegate:"
	// ObjCSelWindowShouldClose is the NSWindowDelegate method -windowShouldClose:.
	ObjCSelWindowShouldClose = "windowShouldClose:"
	// ObjCSelNextEventMatchingMask is the NSApplication method -nextEventMatchingMask:untilDate:inMode:dequeue:.
	ObjCSelNextEventMatchingMask = "nextEventMatchingMask:untilDate:inMode:dequeue:"
)

const (
	// ObjCClassNSDate is the UTF-8 class name for NSDate.
	ObjCClassNSDate = "NSDate"
)

const (
	// ObjCSelDistantPast is the NSDate class method +distantPast.
	ObjCSelDistantPast = "distantPast"
)

const (
	// NSRunLoopModeDefault is the NSDefaultRunLoopMode constant string.
	NSRunLoopModeDefault = "NSDefaultRunLoopMode"
)

const (
	// NSEventMaskAny matches all event types for nextEventMatchingMask.
	NSEventMaskAny uint64 = 0xFFFFFFFFFFFFFFFF
)

const (
	// ObjCTypeEncodingWindowShouldClose is the type encoding for -windowShouldClose: returning BOOL.
	ObjCTypeEncodingWindowShouldClose = "c24@0:8@16"
)
