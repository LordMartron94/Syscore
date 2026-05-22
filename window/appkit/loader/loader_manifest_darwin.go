//go:build darwin

package loader

import (
	"fmt"

	"syscore/window/loadutil"
)

func AppkitCommandsLoadManifest(module AppkitModule, manifest *AppkitCommandManifest, commands *AppkitCommands) error {
	if manifest == nil {
		return fmt.Errorf("appkit loader: manifest must not be nil")
	}
	if commands == nil {
		return fmt.Errorf("appkit loader: commands must not be nil")
	}

	mappings := make([]loadutil.CommandMapping, 0, len(manifest.fields))
	for _, field := range manifest.fields {
		switch field {
		case AppkitCommandManifestFieldObjcGetClass:
			mappings = append(mappings, loadutil.CommandMapping{Target: &commands.ObjcGetClass, Name: "objc_getClass"})
		case AppkitCommandManifestFieldSelRegisterName:
			mappings = append(mappings, loadutil.CommandMapping{Target: &commands.SelRegisterName, Name: "sel_registerName"})
		case AppkitCommandManifestFieldObjcMsgSend:
			mappings = append(mappings, loadutil.CommandMapping{Target: &commands.ObjcMsgSend, Name: "objc_msgSend"})
		default:
			return fmt.Errorf("appkit loader: unknown manifest field %d", field)
		}
	}

	return appkitCommandsBind(module, commands, mappings)
}
