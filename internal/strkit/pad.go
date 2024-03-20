package strkit

import "strings"

func PadLeft(str string, size int, pad byte) string {
	if len(str) >= size {
		return str
	}
	return strings.Repeat(string(pad), size-len(str)) + str
}
func PadRight(str string, size int, pad byte) string {
	if len(str) >= size {
		return str
	}
	return str + strings.Repeat(string(pad), size-len(str))
}
