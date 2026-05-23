package loader

import (
	"fmt"

	"syscore/window/wayland/bindings"
)

type WaylandCommandManifestField uint8

const (
	WaylandCommandManifestFieldInvalid WaylandCommandManifestField = iota
	WaylandCommandManifestFieldDisplayConnect
	WaylandCommandManifestFieldProxyMarshalConstructor
	WaylandCommandManifestFieldProxyMarshalConstructorVersioned
	WaylandCommandManifestFieldProxyAddListener
	WaylandCommandManifestFieldDisplayRoundtrip
	WaylandCommandManifestFieldProxyMarshal
	WaylandCommandManifestFieldProxyDestroy
	WaylandCommandManifestFieldDisplayDisconnect
)

type WaylandCommandManifest struct {
	fields []WaylandCommandManifestField
}

func WaylandCommandManifestReset(manifest *WaylandCommandManifest) {
	if manifest == nil {
		return
	}
	manifest.fields = manifest.fields[:0]
}

func WaylandCommandManifestAddDisplayConnect(manifest *WaylandCommandManifest, target *bindings.PFN_wl_display_connect) error {
	return waylandCommandManifestAdd(manifest, WaylandCommandManifestFieldDisplayConnect, target)
}

func WaylandCommandManifestAddProxyMarshalConstructor(manifest *WaylandCommandManifest, target *bindings.PFN_wl_proxy_marshal_constructor) error {
	return waylandCommandManifestAdd(manifest, WaylandCommandManifestFieldProxyMarshalConstructor, target)
}

func WaylandCommandManifestAddProxyMarshalConstructorVersioned(manifest *WaylandCommandManifest, target *bindings.PFN_wl_proxy_marshal_constructor_versioned) error {
	return waylandCommandManifestAdd(manifest, WaylandCommandManifestFieldProxyMarshalConstructorVersioned, target)
}

func WaylandCommandManifestAddProxyAddListener(manifest *WaylandCommandManifest, target *bindings.PFN_wl_proxy_add_listener) error {
	return waylandCommandManifestAdd(manifest, WaylandCommandManifestFieldProxyAddListener, target)
}

func WaylandCommandManifestAddDisplayRoundtrip(manifest *WaylandCommandManifest, target *bindings.PFN_wl_display_roundtrip) error {
	return waylandCommandManifestAdd(manifest, WaylandCommandManifestFieldDisplayRoundtrip, target)
}

func WaylandCommandManifestAddProxyMarshal(manifest *WaylandCommandManifest, target *bindings.PFN_wl_proxy_marshal) error {
	return waylandCommandManifestAdd(manifest, WaylandCommandManifestFieldProxyMarshal, target)
}

func WaylandCommandManifestAddProxyDestroy(manifest *WaylandCommandManifest, target *bindings.PFN_wl_proxy_destroy) error {
	return waylandCommandManifestAdd(manifest, WaylandCommandManifestFieldProxyDestroy, target)
}

func WaylandCommandManifestAddDisplayDisconnect(manifest *WaylandCommandManifest, target *bindings.PFN_wl_display_disconnect) error {
	return waylandCommandManifestAdd(manifest, WaylandCommandManifestFieldDisplayDisconnect, target)
}

func waylandCommandManifestAdd(manifest *WaylandCommandManifest, field WaylandCommandManifestField, target any) error {
	if manifest == nil {
		return fmt.Errorf("wayland loader: manifest must not be nil")
	}
	if target == nil {
		return fmt.Errorf("wayland loader: manifest target for field %d must not be nil", field)
	}
	manifest.fields = append(manifest.fields, field)
	return nil
}
