//go:build memforge_debug

package internal

import (
	"memcore"
	"strings"
	"testing"
	"unsafe"
)

func TestSyscoreCStringBytePointerLengthAllowsHeapPointer(t *testing.T) {
	heapCString, firstByte := StringToCString("heap")
	_ = heapCString

	length := CStringBytePointerLength(firstByte)
	if length != 4 {
		t.Fatalf("expected length 4, got %d", length)
	}
}

func TestSyscoreCStringBytePointerToStringAllowsMemcoreRegion(t *testing.T) {
	buffer := []byte{'o', 'k', 0}
	regionID := memcoreRegionRegisterForTest(buffer)
	defer memcoreRegionUnregisterForTest(regionID)

	ptr := (*byte)(unsafe.Pointer(&buffer[0]))
	value := CStringBytePointerToString(ptr)
	if value != "ok" {
		t.Fatalf("expected \"ok\", got %q", value)
	}
}

func TestFunctionBindAddressValidatesAtBindingExecution(t *testing.T) {
	var target func()
	err := FunctionBindAddress(&target, 0)
	if err == nil {
		t.Fatal("expected error for zero binding address under memforge_debug")
	}
	if !strings.Contains(err.Error(), "failed to bind function") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func memcoreRegionRegisterForTest(buffer []byte) uint32 {
	return memcore.MemcoreRegionRegister(uintptr(unsafe.Pointer(&buffer[0])), uint64(len(buffer)))
}

func memcoreRegionUnregisterForTest(regionID uint32) {
	memcore.MemcoreRegionUnregister(regionID)
}
