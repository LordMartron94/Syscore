//go:build darwin

package loader

import (
	"fmt"

	"syscore/internal"
	"syscore/window/loadutil"
)

const (
	appkitObjcLibraryName      = "/usr/lib/libobjc.A.dylib"
	appkitFrameworkLibraryName = "/System/Library/Frameworks/AppKit.framework/AppKit"
)

type AppkitModule struct {
	Objc   internal.DynamicLibrary
	Appkit internal.DynamicLibrary
}

func AppkitModuleLoad() (AppkitModule, error) {
	objc, err := internal.DynamicLibraryLoad(appkitObjcLibraryName)
	if err != nil {
		return AppkitModule{}, fmt.Errorf("appkit loader: objc: %w", err)
	}
	appkit, err := internal.DynamicLibraryLoad(appkitFrameworkLibraryName)
	if err != nil {
		return AppkitModule{}, fmt.Errorf("appkit loader: appkit: %w", err)
	}
	return AppkitModule{Objc: objc, Appkit: appkit}, nil
}

func AppkitModuleObjcLibraryName() string   { return appkitObjcLibraryName }
func AppkitModuleAppkitLibraryName() string { return appkitFrameworkLibraryName }

func AppkitModuleObjcSymbolResolve(module AppkitModule, symbolName string) (uintptr, error) {
	return loadutil.LibrarySymbolResolve(module.Objc, symbolName)
}

func AppkitModuleAppkitSymbolResolve(module AppkitModule, symbolName string) (uintptr, error) {
	return loadutil.LibrarySymbolResolve(module.Appkit, symbolName)
}
