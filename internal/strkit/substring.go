package strkit

import "strings"

func AfterByte(str string, c byte) (string, bool) {
	if pos := strings.IndexByte(str, c); pos >= 0 && pos+1 < len(str) {
		return str[pos+1:], true
	}
	return "", false
}

func LastAfterByte(str string, c byte) (string, bool) {
	if pos := strings.LastIndexByte(str, c); pos >= 0 {
		return str[pos+1:], true
	}
	return "", false
}

func BeforeByte(str string, c byte) (string, bool) {
	if pos := strings.IndexByte(str, c); pos >= 0 {
		return str[:pos], true
	}
	return "", false
}

func LastBeforeByte(str string, c byte) (string, bool) {
	if pos := strings.LastIndexByte(str, c); pos >= 0 {
		return str[:pos], true
	}
	return "", false
}

func MapByte(mapping func(b byte) byte, s string) string {
	if s == "" {
		return s
	}

	var b strings.Builder
	var changed = false
	for i, c := range []byte(s) {
		r := mapping(c)
		if !changed {
			if r == c {
				continue
			}
			changed = true
			b.WriteString(s[:i])
			b.WriteByte(r)
		} else {
			b.WriteByte(r)
		}
	}

	if changed {
		return b.String()
	} else {
		return s
	}
}
