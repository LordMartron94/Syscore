package internal

import (
	"fmt"
	"reflect"
	"unsafe"
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
	syscoreLibraryFunctionBindExecute(library, targetFn, functionName)
	return nil
}

func DynamicLibrarySymbolResolve(library DynamicLibrary, symbolName string) (uintptr, error) {
	return platformLibrarySymbolResolve(library, symbolName)
}

func FunctionBindAddress(targetFn any, address uintptr) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("failed to bind function '%v': %v", targetFn, r)
		}
	}()
	syscoreFunctionBindAddressExecute(targetFn, address)
	return nil
}

func PlatformNewCallback(fn any) uintptr {
	return platformNewCallback(fn)
}

func FunctionPointerAddressSet(targetFnPointer any, address uintptr) error {
	SyscoreFunctionBindingTargetValidate(targetFnPointer, "FunctionPointerAddressSet")
	SyscoreFunctionBindingAddressValidate(address, "FunctionPointerAddressSet")

	if targetFnPointer == nil {
		return fmt.Errorf("failed to set function pointer address: nil target")
	}
	if address == 0 {
		return fmt.Errorf("failed to set function pointer address: zero address")
	}

	targetValue := reflect.ValueOf(targetFnPointer)
	if targetValue.Kind() != reflect.Ptr || targetValue.IsNil() {
		return fmt.Errorf("failed to set function pointer address: target must be a non-nil pointer to a function value")
	}
	if targetValue.Elem().Kind() != reflect.Func {
		return fmt.Errorf("failed to set function pointer address: target must point to a function value")
	}

	targetAddress := targetValue.Pointer()
	if targetAddress == 0 {
		return fmt.Errorf("failed to set function pointer address: invalid target pointer")
	}
	*(*uintptr)(unsafe.Pointer(targetAddress)) = address
	return nil
}
