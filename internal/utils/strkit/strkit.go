package strkit

import "strings"

func Cut(s string, sep string) (before, after string, found bool) {
	return strings.Cut(s, sep)
}

func LastCut(s string, sep string) (before, after string, found bool) {
	if i := strings.LastIndex(s, sep); i >= 0 {
		return s[:i], s[i+len(sep):], true
	}
	return "", s, false
}

func ContainsByte(s string, c byte) bool {
	return strings.IndexByte(s, c) >= 0
}

func JoinFunc[T any](elems []T, sep string, fn func(T) string) string {
	switch len(elems) {
	case 0:
		return ""
	case 1:
		return fn(elems[0])
	}

	var elemStrings = make([]string, len(elems))
	for i, elem := range elems {
		elemStrings[i] = fn(elem)
	}
	return strings.Join(elemStrings, sep)
}
