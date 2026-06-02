//go:build memforge_debug

package internal

import (
	"fmt"
	"memcore"
	"reflect"
	"unsafe"
)

func SyscorePointerValidateMemcoreManaged(pointer unsafe.Pointer, context string) {
	if pointer == nil {
		return
	}
	address := uintptr(pointer)
	if memcore.MemcoreAddressBelongsToActiveRegion(address) {
		return
	}
	panic(fmt.Errorf(
		"syscore: %s received pointer outside memcore-managed region (0x%x)",
		context,
		address,
	))
}

func SyscorePointerAddressValidateMemcoreManaged(pointerAddress uintptr, context string) {
	if pointerAddress == 0 {
		return
	}
	if memcore.MemcoreAddressBelongsToActiveRegion(pointerAddress) {
		return
	}
	panic(fmt.Errorf(
		"syscore: %s received address outside memcore-managed region (0x%x)",
		context,
		pointerAddress,
	))
}

func SyscoreFunctionBindingTargetValidate(targetFn any, context string) {
	if targetFn == nil {
		panic(fmt.Errorf("syscore: %s target function must not be nil", context))
	}
	targetValue := reflect.ValueOf(targetFn)
	if targetValue.Kind() != reflect.Ptr || targetValue.IsNil() {
		panic(fmt.Errorf("syscore: %s target must be a non-nil pointer to a function", context))
	}
	if targetValue.Elem().Kind() != reflect.Func {
		panic(fmt.Errorf("syscore: %s target must point to a function", context))
	}
}

func SyscoreFunctionBindingAddressValidate(address uintptr, context string) {
	if address == 0 {
		panic(fmt.Errorf("syscore: %s binding address must not be zero", context))
	}
}
