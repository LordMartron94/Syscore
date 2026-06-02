//go:build memforge_debug

package internal

import "github.com/ebitengine/purego"

func syscoreLibraryFunctionBindExecute(library DynamicLibrary, targetFn any, functionName string) {
	SyscoreFunctionBindingTargetValidate(targetFn, "LibraryFunctionBind")
	purego.RegisterLibFunc(targetFn, uintptr(library), functionName)
}

func syscoreFunctionBindAddressExecute(targetFn any, address uintptr) {
	SyscoreFunctionBindingTargetValidate(targetFn, "FunctionBindAddress")
	SyscoreFunctionBindingAddressValidate(address, "FunctionBindAddress")
	purego.RegisterFunc(targetFn, address)
}
