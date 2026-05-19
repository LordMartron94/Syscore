package internal

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
