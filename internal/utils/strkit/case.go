package strkit

import "strings"

func UpperFirst(str string) string {
	switch len(str) {
	case 0:
		return ""
	case 1:
		return strings.ToUpper(str)
	default:
		return strings.ToUpper(str[0:1]) + str[1:]
	}
}

func LowerFirst(str string) string {
	switch len(str) {
	case 0:
		return ""
	case 1:
		return strings.ToLower(str)
	default:
		return strings.ToLower(str[0:1]) + str[1:]
	}
}

func CamelCase(str string) string {
	parts := strings.Split(str, "_")
	ucParts := make([]string, len(parts))
	for i, part := range parts {
		if i == 0 {
			ucParts[i] = part
		} else {
			ucParts[i] = UpperFirst(part)
		}
	}
	return strings.Join(ucParts, "")
}

func UpperCamelCase(str string) string {
	parts := strings.Split(str, "_")
	ucParts := make([]string, len(parts))
	for i, part := range parts {
		ucParts[i] = UpperFirst(part)
	}
	return strings.Join(ucParts, "")
}
