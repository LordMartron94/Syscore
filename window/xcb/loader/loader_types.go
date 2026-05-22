package loader

import "syscore/window/xcb/bindings"

type XcbCommands struct {
	Connect            bindings.PFN_xcb_connect
	GetSetup           bindings.PFN_xcb_get_setup
	SetupRootsIterator bindings.PFN_xcb_setup_roots_iterator
	GenerateId         bindings.PFN_xcb_generate_id
	CreateWindow       bindings.PFN_xcb_create_window
	MapWindow          bindings.PFN_xcb_map_window
	Flush              bindings.PFN_xcb_flush
}
