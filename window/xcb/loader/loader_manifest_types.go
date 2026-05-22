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
