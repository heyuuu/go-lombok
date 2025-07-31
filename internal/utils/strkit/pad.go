package strkit

import "strings"

func PadRight(str string, size int, pad byte) string {
	if len(str) >= size {
		return str
	}
	return str + strings.Repeat(string(pad), size-len(str))
}
