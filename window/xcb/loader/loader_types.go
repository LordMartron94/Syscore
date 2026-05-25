package loader

import "syscore/window/xcb/bindings"

/*
XcbCommands holds bound XCB function pointers resolved from libxcb.so.1.

[Context]
Populate with XcbCommandsLoad or XcbCommandsLoadManifest, then call through the fields.
Each field name matches the xcb_* symbol bound at load time.
*/
type XcbCommands struct {
	/*
		Connect is xcb_connect. Opens an XCB connection; use nil displayname for DISPLAY.
	*/
	Connect bindings.PFN_xcb_connect
	/*
		GetSetup is xcb_get_setup. Returns setup pointer for screen/root discovery.
	*/
	GetSetup bindings.PFN_xcb_get_setup
	/*
		SetupRootsIterator is xcb_setup_roots_iterator. Iterates screens from setup.
	*/
	SetupRootsIterator bindings.PFN_xcb_setup_roots_iterator
	/*
		GenerateId is xcb_generate_id. Allocates a new window/resource ID.
	*/
	GenerateId bindings.PFN_xcb_generate_id
	/*
		CreateWindow is xcb_create_window. Creates a child window (unchecked cookie).
	*/
	CreateWindow bindings.PFN_xcb_create_window
	/*
		MapWindow is xcb_map_window. Maps a window for display (unchecked cookie).
	*/
	MapWindow bindings.PFN_xcb_map_window
	/*
		Flush is xcb_flush. Sends buffered requests to the X server.
	*/
	Flush bindings.PFN_xcb_flush
	/*
		InternAtom is xcb_intern_atom. Resolves property and type atom names.
	*/
	InternAtom bindings.PFN_xcb_intern_atom
	/*
		InternAtomReply is xcb_intern_atom_reply. Blocks for an intern_atom reply.
	*/
	InternAtomReply bindings.PFN_xcb_intern_atom_reply
	/*
		ChangeProperty is xcb_change_property. Sets window properties such as WM_NAME.
	*/
	ChangeProperty bindings.PFN_xcb_change_property
	/*
		Free is libc free(3). Releases malloc'd reply buffers from InternAtomReply.
	*/
	Free bindings.PFN_c_free
	/*
		DestroyWindow is xcb_destroy_window. Destroys a window (unchecked cookie).
	*/
	DestroyWindow bindings.PFN_xcb_destroy_window
	/*
		Disconnect is xcb_disconnect. Closes the XCB connection.
	*/
	Disconnect bindings.PFN_xcb_disconnect
	/*
		PollForEvent is xcb_poll_for_event. Returns the next queued event or nil.
	*/
	PollForEvent bindings.PFN_xcb_poll_for_event
}
