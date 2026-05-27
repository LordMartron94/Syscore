//go:build linux

package loader

import (
	"fmt"

	"syscore/window/xcb/bindings"
)

type XcbCommandManifestField uint8

const (
	XcbCommandManifestFieldInvalid XcbCommandManifestField = iota
	XcbCommandManifestFieldConnect
	XcbCommandManifestFieldGetSetup
	XcbCommandManifestFieldSetupRootsIterator
	XcbCommandManifestFieldGenerateId
	XcbCommandManifestFieldCreateWindow
	XcbCommandManifestFieldMapWindow
	XcbCommandManifestFieldFlush
	XcbCommandManifestFieldInternAtom
	XcbCommandManifestFieldInternAtomReply
	XcbCommandManifestFieldChangeProperty
	XcbCommandManifestFieldFree
	XcbCommandManifestFieldDestroyWindow
	XcbCommandManifestFieldDisconnect
)

type XcbCommandManifest struct {
	fields []XcbCommandManifestField
}

func XcbCommandManifestReset(manifest *XcbCommandManifest) {
	if manifest == nil {
		return
	}
	manifest.fields = manifest.fields[:0]
}

func XcbCommandManifestAddConnect(manifest *XcbCommandManifest, target *bindings.PFN_xcb_connect) error {
	return xcbCommandManifestAdd(manifest, XcbCommandManifestFieldConnect, target)
}

func XcbCommandManifestAddGetSetup(manifest *XcbCommandManifest, target *bindings.PFN_xcb_get_setup) error {
	return xcbCommandManifestAdd(manifest, XcbCommandManifestFieldGetSetup, target)
}

func XcbCommandManifestAddSetupRootsIterator(manifest *XcbCommandManifest, target *bindings.PFN_xcb_setup_roots_iterator) error {
	return xcbCommandManifestAdd(manifest, XcbCommandManifestFieldSetupRootsIterator, target)
}

func XcbCommandManifestAddGenerateId(manifest *XcbCommandManifest, target *bindings.PFN_xcb_generate_id) error {
	return xcbCommandManifestAdd(manifest, XcbCommandManifestFieldGenerateId, target)
}

func XcbCommandManifestAddCreateWindow(manifest *XcbCommandManifest, target *bindings.PFN_xcb_create_window) error {
	return xcbCommandManifestAdd(manifest, XcbCommandManifestFieldCreateWindow, target)
}

func XcbCommandManifestAddMapWindow(manifest *XcbCommandManifest, target *bindings.PFN_xcb_map_window) error {
	return xcbCommandManifestAdd(manifest, XcbCommandManifestFieldMapWindow, target)
}

func XcbCommandManifestAddFlush(manifest *XcbCommandManifest, target *bindings.PFN_xcb_flush) error {
	return xcbCommandManifestAdd(manifest, XcbCommandManifestFieldFlush, target)
}

func XcbCommandManifestAddInternAtom(manifest *XcbCommandManifest, target *bindings.PFN_xcb_intern_atom) error {
	return xcbCommandManifestAdd(manifest, XcbCommandManifestFieldInternAtom, target)
}

func XcbCommandManifestAddInternAtomReply(manifest *XcbCommandManifest, target *bindings.PFN_xcb_intern_atom_reply) error {
	return xcbCommandManifestAdd(manifest, XcbCommandManifestFieldInternAtomReply, target)
}

func XcbCommandManifestAddChangeProperty(manifest *XcbCommandManifest, target *bindings.PFN_xcb_change_property) error {
	return xcbCommandManifestAdd(manifest, XcbCommandManifestFieldChangeProperty, target)
}

func XcbCommandManifestAddFree(manifest *XcbCommandManifest, target *bindings.PFN_c_free) error {
	return xcbCommandManifestAdd(manifest, XcbCommandManifestFieldFree, target)
}

func XcbCommandManifestAddDestroyWindow(manifest *XcbCommandManifest, target *bindings.PFN_xcb_destroy_window) error {
	return xcbCommandManifestAdd(manifest, XcbCommandManifestFieldDestroyWindow, target)
}

func XcbCommandManifestAddDisconnect(manifest *XcbCommandManifest, target *bindings.PFN_xcb_disconnect) error {
	return xcbCommandManifestAdd(manifest, XcbCommandManifestFieldDisconnect, target)
}

func xcbCommandManifestAdd(manifest *XcbCommandManifest, field XcbCommandManifestField, target any) error {
	if manifest == nil {
		return fmt.Errorf("xcb loader: manifest must not be nil")
	}
	if target == nil {
		return fmt.Errorf("xcb loader: manifest target for field %d must not be nil", field)
	}
	manifest.fields = append(manifest.fields, field)
	return nil
}
