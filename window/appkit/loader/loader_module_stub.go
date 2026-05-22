//go:build !darwin

package loader

import "fmt"

type AppkitModule struct{}

func AppkitModuleLoad() (AppkitModule, error) {
	return AppkitModule{}, fmt.Errorf("appkit loader: only supported on darwin")
}

func AppkitModuleObjcLibraryName() string   { return "" }
func AppkitModuleAppkitLibraryName() string { return "" }

func AppkitModuleObjcSymbolResolve(_ AppkitModule, _ string) (uintptr, error) {
	return 0, fmt.Errorf("appkit loader: only supported on darwin")
}

func AppkitModuleAppkitSymbolResolve(_ AppkitModule, _ string) (uintptr, error) {
	return 0, fmt.Errorf("appkit loader: only supported on darwin")
}
