package syscore

import "syscore/internal"

/*
SYSCORE_Pure_DynamicLibrary is an opaque handle to a dynamically loaded native library.

[Context]
Obtained from SYSCORE_Pure_LibraryLoad and passed to bind and symbol-resolve functions.
On Windows this is a module handle from LoadLibrary; on Linux and macOS it is a dlopen handle.
*/
type SYSCORE_Pure_DynamicLibrary = internal.DynamicLibrary

/*
SYSCORE_Pure_LibraryLoad loads a dynamic library into the current process.

[Context]
Wraps dlopen (Linux, macOS) or LoadLibrary (Windows) via purego. Used by window loaders and
consumers such as gpuarch that resolve ICD entry points.

[Parameters]
name - Library file name or path (for example "libxcb.so.1", "user32.dll").

[Returns]
A library handle on success.

[Errors]
Returns an error if the platform loader fails (library missing, wrong architecture, etc.).

[Side Effects]
Loads native code into the process; the library remains loaded until process exit on most platforms.

[Example]

	lib, err := syscore.SYSCORE_Pure_LibraryLoad("libvulkan.so.1")
*/
func SYSCORE_Pure_LibraryLoad(name string) (SYSCORE_Pure_DynamicLibrary, error) {
	return internal.DynamicLibraryLoad(name)
}

/*
SYSCORE_Pure_LibraryFunctionBind resolves an exported symbol by name and binds it to a Go function pointer.

[Context]
Uses the library's dynamic symbol table. targetFn must be a pointer to a func variable whose
signature matches the C ABI of the exported function.

[Parameters]
library - Handle from SYSCORE_Pure_LibraryLoad.
targetFn - Pointer to func variable to receive the binding (for example &myPFN).
functionName - Exported symbol name (for example "xcb_connect").

[Returns]
nil on success.

[Errors]
Returns an error if the symbol is missing or purego cannot register the function type.

[Side Effects]
Mutates the func variable at targetFn. Subsequent calls invoke native code.

[Example]

	var connect func(display *byte, screen *int32) uintptr
	err := syscore.SYSCORE_Pure_LibraryFunctionBind(lib, &connect, "xcb_connect")
*/
func SYSCORE_Pure_LibraryFunctionBind(library SYSCORE_Pure_DynamicLibrary, targetFn any, functionName string) error {
	return internal.LibraryFunctionBind(library, targetFn, functionName)
}

/*
SYSCORE_Pure_FunctionBindAddress binds a known native address to a Go function pointer.

[Context]
Use when the address was obtained from SYSCORE_Pure_LibrarySymbolResolve or from a Vulkan
GetInstanceProcAddr-style query rather than by exported name in the same library.

[Parameters]
targetFn - Pointer to func variable to receive the binding.
address - Native code pointer (non-zero).

[Returns]
nil on success.

[Errors]
Returns an error if targetFn is not a function pointer or the signature is unsupported by purego.

[Side Effects]
Mutates the func variable at targetFn.
*/
func SYSCORE_Pure_FunctionBindAddress(targetFn any, address uintptr) error {
	return internal.FunctionBindAddress(targetFn, address)
}

/*
SYSCORE_Pure_LibrarySymbolResolve returns the address of an exported symbol without creating a Go callable.

[Context]
Use for WNDPROC pointers, wl_*_interface globals, or passing the address into struct fields
before optional SYSCORE_Pure_FunctionBindAddress on a separate func variable.

[Parameters]
library - Handle from SYSCORE_Pure_LibraryLoad.
symbolName - Exported symbol name.

[Returns]
Symbol address on success.

[Errors]
Returns an error if the symbol is not found or the handle is invalid.

[Side Effects]
None.

[Example]

	addr, err := syscore.SYSCORE_Pure_LibrarySymbolResolve(lib, "wl_compositor_interface")
*/
func SYSCORE_Pure_LibrarySymbolResolve(library SYSCORE_Pure_DynamicLibrary, symbolName string) (uintptr, error) {
	return internal.DynamicLibrarySymbolResolve(library, symbolName)
}

/*
SYSCORE_Pure_NewCallback converts a Go function to a C-callable function pointer for native callbacks.

[Context]
Wraps purego.NewCallback (or the platform equivalent). Use for Wayland wl_registry_listener and
other FFI vtables. The function must use uintptr-sized arguments only (no string or slice parameters).

[Parameters]
fn - Go function matching the C callback ABI.

[Returns]
Trampoline address to store in listener vtables.

[Side Effects]
Callback slots are process-lifetime limited; memory is not released.

[Example]

	cb := syscore.SYSCORE_Pure_NewCallback(func(data uintptr, registry uintptr, name uint32, iface uintptr, version uint32) {
	    _ = syscore.SYSCORE_C_CStringPointerToString(iface)
	})
*/
func SYSCORE_Pure_NewCallback(fn any) uintptr {
	return internal.PlatformNewCallback(fn)
}

/*
SYSCORE_Pure_FunctionPointerAddressSet writes a native function address directly into a function-pointer slot.

[Context]
Use this for callback fields inside native struct mirrors where the native API will call back into Go through a
SYSCORE_Pure_NewCallback trampoline address. This avoids creating a Go callable wrapper, which is what
SYSCORE_Pure_FunctionBindAddress does.

[Parameters]
targetFnPointer - Pointer to a function-typed field or variable that stores a native callback address.
address - Native function address (for example return value from SYSCORE_Pure_NewCallback).

[Returns]
nil on success.

[Errors]
Returns an error if targetFnPointer is nil, not a pointer to a function value, or address is zero.

[Side Effects]
Mutates the function-pointer slot at targetFnPointer by writing address directly.

[Example]

	callbackAddress := syscore.SYSCORE_Pure_NewCallback(myCallback)
	err := syscore.SYSCORE_Pure_FunctionPointerAddressSet(&createInfo.PfnUserCallback, callbackAddress)
*/
func SYSCORE_Pure_FunctionPointerAddressSet(targetFnPointer any, address uintptr) error {
	return internal.FunctionPointerAddressSet(targetFnPointer, address)
}
