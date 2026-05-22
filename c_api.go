package syscore

import "syscore/internal"

/*
SYSCORE_C_StringToCString converts a standard Go string to a null-terminated C string.

It returns both the actual array and the pointer to the first byte.
*/
func SYSCORE_C_StringToCString(str string) ([]byte, *byte) {
	return internal.StringToCString(str)
}

/*
SYSCORE_C_StringToCStringFirstByte is the same as SYSCORE_C_StringToCString but returns only the pointer to the first byte.
*/
func SYSCORE_C_StringToCStringFirstByte(str string) *byte {
	_, firstByte := internal.StringToCString(str)
	return firstByte
}

/*
SYSCORE_C_StringToCStringSlice is the same as SYSCORE_C_StringToCString but returns only the slice.
*/
func SYSCORE_C_StringToCStringSlice(str string) []byte {
	array, _ := internal.StringToCString(str)
	return array
}

/*
SYSCORE_C_CStringToString converts a null-terminated C string into a standard Go string.
*/
func SYSCORE_C_CStringToString(str []byte) string {
	return internal.CStringToString(str)
}

/*
SYSCORE_C_CStringPointerToString converts a NUL-terminated C string pointer into a Go string.
*/
func SYSCORE_C_CStringPointerToString(ptr uintptr) string {
	return internal.CStringPointerToString(ptr)
}

/*
SYSCORE_C_StringToUTF16 converts a Go string to a NUL-terminated UTF-16 string for Win32 APIs.

It returns both the code unit slice and a pointer to the first element.
*/
func SYSCORE_C_StringToUTF16(str string) ([]uint16, *uint16) {
	return internal.StringToUTF16(str)
}
