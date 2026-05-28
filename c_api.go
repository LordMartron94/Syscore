package syscore

import (
	"memarch"
	"memcore"
	"syscore/internal"
)

/*
SYSCORE_C_StringToCString converts a Go string to a NUL-terminated UTF-8 C string.

[Context]
Many FFI calls expect char* pointing at bytes that remain valid for the duration of the
call. The returned slice owns the storage; keep it alive while the C API uses the pointer.

[Parameters]
str - Source string; an extra NUL byte is appended.

[Returns]
A byte slice including the terminator and a pointer to its first element (suitable for
*byte parameters in purego bindings).

[Side Effects]
Allocates a new slice. Pure aside from allocation.

[Example]

	name, namePtr := syscore.SYSCORE_C_StringToCString("vkCreateInstance")
	_ = name
	commands.CreateInstance(namePtr, ...)
*/
func SYSCORE_C_StringToCString(str string) ([]byte, *byte) {
	return internal.StringToCString(str)
}

/*
SYSCORE_C_StringToCStringManual allocates a NUL-terminated UTF-8 C string in manual memory.

[Context]
Use this for FFI pointers that must not rely on Go heap slice lifetime. The returned pointer
is backed by memory allocated via the provided memarch allocation function.

[Parameters]
allocFn - Allocation function used to reserve manual memory.
str - Source string encoded as UTF-8 with an appended NUL terminator.

[Returns]
The manual memory mark and a pointer to the first byte of the C string.

[Side Effects]
Allocates memory through allocFn.
*/
func SYSCORE_C_StringToCStringManual(allocFn memarch.AllocationFn, str string) (memcore.MarkRaw, *byte) {
	return memarch.MemArchCStringCreate(allocFn, str)
}

/*
SYSCORE_C_StringToCStringFirstByteManual allocates a NUL-terminated UTF-8 C string in manual
memory and returns only the first-byte pointer.

[Context]
Use this for APIs that only require *byte while still keeping storage off the Go heap.

[Parameters]
allocFn - Allocation function used to reserve manual memory.
str - Source string encoded as UTF-8 with an appended NUL terminator.

[Returns]
Pointer to the first byte of the manually allocated C string.

[Side Effects]
Allocates memory through allocFn.
*/
func SYSCORE_C_StringToCStringFirstByteManual(allocFn memarch.AllocationFn, str string) *byte {
	_, firstByte := memarch.MemArchCStringCreate(allocFn, str)
	return firstByte
}

/*
SYSCORE_C_StringToCStringFirstByte is the same as SYSCORE_C_StringToCString but returns only the pointer.

[Context]
Use when the backing slice is kept alive elsewhere or only the pointer is passed into C.

[Parameters]
str - Source string; an extra NUL byte is appended in a temporary slice that remains
reachable only if you also retain the slice from SYSCORE_C_StringToCString.

[Returns]
Pointer to the first byte of the NUL-terminated encoding.

[Side Effects]
Allocates a slice internally; the pointer is invalid after that slice is garbage-collected
unless you retain the full slice from SYSCORE_C_StringToCString.

[Example]

	ptr := syscore.SYSCORE_C_StringToCStringFirstByte("xcb_connect")
	// Prefer SYSCORE_C_StringToCString when the slice must outlive the call.
*/
func SYSCORE_C_StringToCStringFirstByte(str string) *byte {
	_, firstByte := internal.StringToCString(str)
	return firstByte
}

/*
SYSCORE_C_StringToCStringSlice is the same as SYSCORE_C_StringToCString but returns only the slice.

[Context]
Use when you need to retain the backing storage in a []byte without using the pointer form.

[Parameters]
str - Source string; an extra NUL byte is appended.

[Returns]
NUL-terminated UTF-8 bytes.

[Side Effects]
Allocates a new slice.
*/
func SYSCORE_C_StringToCStringSlice(str string) []byte {
	array, _ := internal.StringToCString(str)
	return array
}

/*
SYSCORE_C_CStringToString converts a NUL-terminated C string stored in a Go byte slice to a Go string.

[Context]
Decodes FFI output or inbound char* data that has already been copied into Go memory.

[Parameters]
str - Bytes up to and including the first 0 byte; content after the first NUL is ignored.

[Returns]
A Go string without the terminator.

[Side Effects]
None. Does not mutate str.
*/
func SYSCORE_C_CStringToString(str []byte) string {
	return internal.CStringToString(str)
}

/*
SYSCORE_C_CStringPointerToString reads a NUL-terminated C string from a raw address.

[Context]
Used for Wayland registry global callbacks and similar FFI surfaces that pass const char*
as uintptr. The memory must be readable in the current process.

[Parameters]
ptr - Address of the first byte; 0 returns "".

[Returns]
Decoded string up to the first NUL.

[Side Effects]
Reads foreign memory at ptr until NUL. No writes.

[Edge Cases]
Invalid or unmapped ptr can panic; callers must ensure the pointer came from a live C callback.
*/
func SYSCORE_C_CStringPointerToString(ptr uintptr) string {
	return internal.CStringPointerToString(ptr)
}

/*
SYSCORE_C_CStringBytePointerToString reads a NUL-terminated C string from a *byte pointer.

[Context]
Use when purego bindings expose char* as *byte (for example PFN names passed to
vkGetInstanceProcAddr). Prefer this over SYSCORE_C_CStringPointerToString when the FFI
surface already types the argument as *byte.

[Parameters]
ptr - Pointer to the first byte; nil returns "".

[Returns]
Decoded string up to the first NUL.

[Side Effects]
Reads foreign memory at ptr until NUL. No writes.

[Edge Cases]
Invalid or unmapped ptr can panic; callers must ensure the pointer came from valid FFI memory.
*/
func SYSCORE_C_CStringBytePointerToString(ptr *byte) string {
	return internal.CStringBytePointerToString(ptr)
}

/*
SYSCORE_C_CStringBytePointerLength returns the payload byte length of a NUL-terminated C string.

[Parameters]
ptr - Pointer to the first byte; nil returns 0.

[Returns]
Length in bytes excluding the terminating NUL.
*/
func SYSCORE_C_CStringBytePointerLength(ptr *byte) uint64 {
	return internal.CStringBytePointerLength(ptr)
}

/*
SYSCORE_C_CStringBytePointerToBytes returns a byte-slice view over a C string pointer.

[Context]
This is a zero-copy view into foreign/manual memory. The returned slice is valid only while the
pointed memory remains valid and unchanged.

[Parameters]
ptr - Pointer to a NUL-terminated C string.
includeTerminator - When true, includes the trailing NUL byte in the returned slice.

[Returns]
A byte slice view of the C string bytes; nil for nil pointer (or empty payload when terminator is excluded).
*/
func SYSCORE_C_CStringBytePointerToBytes(ptr *byte, includeTerminator bool) []byte {
	return internal.CStringBytePointerToBytes(ptr, includeTerminator)
}

/*
SYSCORE_C_StringToUTF16 converts a Go string to a NUL-terminated UTF-16 string for Win32 LPCWSTR.

[Context]
Win32 window APIs require wide-character strings. The returned slice owns the storage.

[Parameters]
str - Source Unicode string; a UTF-16 NUL code unit is appended.

[Returns]
UTF-16 code units including the terminator and a pointer to the first element.

[Side Effects]
Allocates a new slice.

[Example]

	class, classPtr := syscore.SYSCORE_C_StringToUTF16("MyWindowClass")
	commands.RegisterClassExW(&wndClass) // wndClass.LpszClassName = classPtr
*/
func SYSCORE_C_StringToUTF16(str string) ([]uint16, *uint16) {
	return internal.StringToUTF16(str)
}
