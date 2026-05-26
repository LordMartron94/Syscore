//go:build darwin

package loader

import (
	"fmt"

	"syscore/window/appkit/bindings"
)

type AppkitCommandManifestField uint8

const (
	AppkitCommandManifestFieldInvalid AppkitCommandManifestField = iota
	AppkitCommandManifestFieldObjcGetClass
	AppkitCommandManifestFieldSelRegisterName
	AppkitCommandManifestFieldObjcMsgSend
)

type AppkitCommandManifest struct {
	fields []AppkitCommandManifestField
}

func AppkitCommandManifestReset(manifest *AppkitCommandManifest) {
	if manifest == nil {
		return
	}
	manifest.fields = manifest.fields[:0]
}

func AppkitCommandManifestAddObjcGetClass(manifest *AppkitCommandManifest, target *bindings.PFN_objc_getClass) error {
	return appkitCommandManifestAdd(manifest, AppkitCommandManifestFieldObjcGetClass, target)
}

func AppkitCommandManifestAddSelRegisterName(manifest *AppkitCommandManifest, target *bindings.PFN_sel_registerName) error {
	return appkitCommandManifestAdd(manifest, AppkitCommandManifestFieldSelRegisterName, target)
}

func AppkitCommandManifestAddObjcMsgSend(manifest *AppkitCommandManifest, target *bindings.PFN_objc_msgSend) error {
	return appkitCommandManifestAdd(manifest, AppkitCommandManifestFieldObjcMsgSend, target)
}

func appkitCommandManifestAdd(manifest *AppkitCommandManifest, field AppkitCommandManifestField, target any) error {
	if manifest == nil {
		return fmt.Errorf("appkit loader: manifest must not be nil")
	}
	if target == nil {
		return fmt.Errorf("appkit loader: manifest target for field %d must not be nil", field)
	}
	manifest.fields = append(manifest.fields, field)
	return nil
}
