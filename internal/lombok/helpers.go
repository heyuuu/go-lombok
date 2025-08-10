package lombok

import (
	"iter"
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

// firstOf 返回 iter.Seq 的首个元素
func firstOf[T any](it iter.Seq[T]) (T, bool) {
	for item := range it {
		return item, true
	}
	var zero T
	return zero, false
}

func writeFileIfChanged(fileName string, content string) (changed bool, err error) {
	oldContent, err := os.ReadFile(fileName)
	if err == nil && string(oldContent) == content {
		return false, nil
	}

	err = os.WriteFile(fileName, []byte(content), 0644)
	return true, err
}

func deleteFileIfExists(fileName string) (exists bool, err error) {
	if _, err := os.Stat(fileName); err != nil {
		return false, nil
	}

	err = os.Remove(fileName)
	return true, err
}

// 判断是否为合法标识符名
// 规则: [a-zA-Z_][a-zA-Z0-9_]*
func isValidIdent(s string) bool {
	if s == "" {
		return false
	}

	for index, c := range []byte(s) {
		if ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z') || c == '_' {
			continue
		}
		if index > 0 && ('0' <= c && c <= '9') {
			continue
		}
		return false
	}

	return true
}
