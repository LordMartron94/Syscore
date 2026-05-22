package internal

import (
	"unicode/utf16"
	"unsafe"
)

func StringToCString(str string) ([]byte, *byte) {
	b := append([]byte(str), 0)
	return b, &b[0]
}

func CStringToString(str []byte) string {
	for i, b := range str {
		if b == 0 {
			return string(str[:i])
		}
	}
	return string(str[:])
}

func CStringPointerToString(ptr uintptr) string {
	if ptr == 0 {
		return ""
	}
	start := unsafe.Pointer(ptr)
	length := 0
	for {
		if *(*byte)(unsafe.Add(start, length)) == 0 {
			break
		}
		length++
	}
	if length == 0 {
		return ""
	}
	slice := unsafe.Slice((*byte)(start), length)
	return string(slice)
}

func StringToUTF16(str string) ([]uint16, *uint16) {
	encoded := utf16.Encode([]rune(str))
	encoded = append(encoded, 0)
	return encoded, &encoded[0]
}
