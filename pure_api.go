package syscore

import "syscore/internal"

/*
SYSCORE_Pure_DynamicLibrary is the address to a dynamic library.
*/
type SYSCORE_Pure_DynamicLibrary = internal.DynamicLibrary

/*
SYSCORE_Pure_LibraryLoad loads a specific dynamic library into memory.

If the library cannot be loaded, it returns an error.
*/
func SYSCORE_Pure_LibraryLoad(name string) (SYSCORE_Pure_DynamicLibrary, error) {
	return internal.DynamicLibraryLoad(name)
}

/*
SYSCORE_Pure_LibraryFunctionBind binds a library function to a target function.

targetFn must be a function pointer.

Returns an error if the function cannot be found.
*/
func SYSCORE_Pure_LibraryFunctionBind(library SYSCORE_Pure_DynamicLibrary, targetFn any, functionName string) error {
	return internal.LibraryFunctionBind(library, targetFn, functionName)
}

/*
SYSCORE_Pure_FunctionBindAddress binds a specific C address to a Go function.

An error is produced if the type is not a function pointer or if the function returns more than 1 value.
*/
func SYSCORE_Pure_FunctionBindAddress(targetFn any, address uintptr) error {
	return internal.FunctionBindAddress(targetFn, address)
}

/*
SYSCORE_Pure_LibrarySymbolResolve returns the address of an exported symbol in a loaded library.
*/
func SYSCORE_Pure_LibrarySymbolResolve(library SYSCORE_Pure_DynamicLibrary, symbolName string) (uintptr, error) {
	return internal.DynamicLibrarySymbolResolve(library, symbolName)
}
