/*
Package bindings provides Objective-C runtime PFN types and AppKit-related constants for macOS windows.

Window creation is performed via objc_msgSend on NSApplication and NSWindow after resolving
classes and selectors with objc_getClass and sel_registerName.
*/
package bindings
