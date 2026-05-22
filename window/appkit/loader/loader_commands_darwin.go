//go:build darwin

package loader

import (
	"fmt"

	"syscore/window/loadutil"
)

func AppkitCommandsLoad(module AppkitModule, commands *AppkitCommands) error {
	return appkitCommandsBind(module, commands, appkitCommandMappingsAll(commands))
}

func appkitCommandMappingsAll(commands *AppkitCommands) []loadutil.CommandMapping {
	return []loadutil.CommandMapping{
		{Target: &commands.ObjcGetClass, Name: "objc_getClass"},
		{Target: &commands.SelRegisterName, Name: "sel_registerName"},
		{Target: &commands.ObjcMsgSend, Name: "objc_msgSend"},
	}
}

func appkitCommandsBind(module AppkitModule, commands *AppkitCommands, mappings []loadutil.CommandMapping) error {
	if commands == nil {
		return fmt.Errorf("appkit loader: commands must not be nil")
	}
	if err := loadutil.LibraryCommandsBind(module.Objc, mappings); err != nil {
		return fmt.Errorf("appkit loader: objc: %w", err)
	}
	return nil
}
