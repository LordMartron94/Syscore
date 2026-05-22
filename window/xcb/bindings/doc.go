/*
Package bindings provides 1:1 XCB types and function pointer typedefs for libxcb.

Clients load symbols into syscore/window/xcb/loader.XcbCommands (or SYSCORE_Window_Xcb_Commands)
and invoke behavior through the bound fields, which use the PFN_* types defined here.
*/
package bindings
