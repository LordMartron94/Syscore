//go:build !memforge_debug

package internal

import "github.com/ebitengine/purego"

func syscoreLibraryFunctionBindExecute(library DynamicLibrary, targetFn any, functionName string) {
	purego.RegisterLibFunc(targetFn, uintptr(library), functionName)
}

func syscoreFunctionBindAddressExecute(targetFn any, address uintptr) {
	purego.RegisterFunc(targetFn, address)
}
