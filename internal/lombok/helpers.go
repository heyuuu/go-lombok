package lombok

import (
	"os"
	"strings"
)

// pascalCase 字符串转大驼峰命名
func pascalCase(s string) string {
	parts := strings.Split(s, "_")

	var buf strings.Builder
	buf.Grow(len(s))
	for _, part := range parts {
		if part == "" {
			continue
		}
		buf.WriteString(strings.ToUpper(part[:1]))
		buf.WriteString(part[1:])
	}
	return buf.String()
}

// padRight 填充字符串到至少指定长度
func padRight(s string, size int, pad byte) string {
	if len(s) >= size {
		return s
	}
	return s + strings.Repeat(string(pad), size-len(s))
}

func writeFileIfChanged(fileName string, content string) (changed bool, err error) {
	existContent, err := os.ReadFile(fileName)
	if err == nil && string(existContent) == content {
		return false, nil
	}

	err = os.WriteFile(fileName, []byte(content), 0644)
	return true, err
}
