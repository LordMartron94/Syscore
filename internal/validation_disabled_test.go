//go:build !memforge_debug

package internal

import "testing"

func TestSyscoreCStringBytePointerLengthAllowsHeapPointerWhenDebugDisabled(t *testing.T) {
	heapCString, firstByte := StringToCString("heap")
	_ = heapCString

	length := CStringBytePointerLength(firstByte)
	if length != 4 {
		t.Fatalf("expected length 4, got %d", length)
	}
}
