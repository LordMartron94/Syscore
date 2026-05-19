package internal

import (
	"fmt"

	"github.com/ebitengine/purego"
)

type DynamicLibrary uintptr

func DynamicLibraryLoad(name string) (DynamicLibrary, error) {
	handle, err := loadPlatformLibrary(name)
	if err != nil {
		return 0, fmt.Errorf("failed to load %s: %w", name, err)
	}

	return DynamicLibrary(handle), nil
}

func LibraryFunctionBind(library DynamicLibrary, targetFn any, functionName string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("failed to bind function '%s': %v", functionName, r)
		}
	}()

	purego.RegisterLibFunc(targetFn, uintptr(library), functionName)
	return nil
}

func FunctionBindAddress(targetFn any, address uintptr) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("failed to bind function '%v': %v", targetFn, r)
		}
	}()

	purego.RegisterFunc(targetFn, address)
	return nil
}
