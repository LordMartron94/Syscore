//go:build !memforge_debug

package internal

import "unsafe"

func SyscorePointerValidateMemcoreManaged(pointer unsafe.Pointer, context string) {
}

func SyscorePointerAddressValidateMemcoreManaged(pointerAddress uintptr, context string) {
}

func SyscoreFunctionBindingTargetValidate(targetFn any, context string) {
}

func SyscoreFunctionBindingAddressValidate(address uintptr, context string) {
}
