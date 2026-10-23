package validatex

import (
	"fmt"
	"strings"
)

// translate 获取翻译后的错误消息
func translate(translator *Translator, path, rule string, args ...any) string {
	if translator != nil && len(translator.Messages) > 0 {
		if message, ok := translator.Messages[path+"."+rule]; ok {
			return formatMessage(translator, message, path, rule, args...)
		}
		if message, ok := lookupWildcardMessage(translator.Messages, path, rule); ok {
			return formatMessage(translator, message, path, rule, args...)
		}
		if message, ok := translator.Messages[rule]; ok {
			return formatMessage(translator, message, path, rule, args...)
		}
	}
	return path + " is " + rule
}

// formatMessage 替换消息中的占位符
func formatMessage(translator *Translator, message, path, rule string, args ...any) string {
	message = strings.ReplaceAll(message, "{field}", displayField(translator, path))
	message = strings.ReplaceAll(message, "{path}", path)
	message = strings.ReplaceAll(message, "{rule}", rule)

	for index, arg := range args {
		value := fmt.Sprint(arg)
		message = strings.ReplaceAll(message, fmt.Sprintf("{%d}", index), value)
		if named, ok := arg.(map[string]any); ok {
			for name, item := range named {
				message = strings.ReplaceAll(message, "{"+name+"}", fmt.Sprint(item))
			}
		}
	}

	namedArgs := false
	if len(args) > 0 {
		_, namedArgs = args[0].(map[string]any)
	}
	if len(args) > 0 && !namedArgs {
		value := fmt.Sprint(args[0])
		for _, name := range []string{"value", "min", "max", "len", "size"} {
			message = strings.ReplaceAll(message, "{"+name+"}", value)
		}
	}
	if len(args) > 1 && !namedArgs {
		message = strings.ReplaceAll(message, "{max}", fmt.Sprint(args[1]))
	}
	return message
}

// displayField 获取字段显示名称
func displayField(translator *Translator, path string) string {
	if translator != nil {
		if field, ok := translator.Translates[path]; ok && field != "" {
			return field
		}
		if field, ok := lookupWildcardValue(translator.Translates, path); ok {
			return field
		}
		field := lastField(path)
		if translated, ok := translator.Translates[field]; ok && translated != "" {
			return translated
		}
	}
	return lastField(path)
}

// lastField 获取路径的最后一个字段
func lastField(path string) string {
	field := path
	if index := strings.LastIndex(field, "."); index >= 0 {
		field = field[index+1:]
	}
	if index := strings.Index(field, "["); index >= 0 {
		field = field[:index]
	}
	if field == "" {
		return path
	}
	return field
}

// lookupWildcardMessage 查找匹配路径的通配符消息
func lookupWildcardMessage(messages Messages, path, rule string) (string, bool) {
	suffix := "." + rule
	normalized := normalizePath(path)
	var (
		bestMessage string
		bestPattern string
		bestStars   = int(^uint(0) >> 1)
		bestLiteral int
	)
	for key, message := range messages {
		if !strings.HasSuffix(key, suffix) {
			continue
		}
		pattern := strings.TrimSuffix(key, suffix)
		if !strings.Contains(pattern, "*") || !matchWildcard(pattern, normalized) {
			continue
		}
		stars, literal := wildcardSpecificity(pattern)
		if stars < bestStars ||
			stars == bestStars && literal > bestLiteral ||
			stars == bestStars && literal == bestLiteral && pattern > bestPattern {
			bestMessage = message
			bestPattern = pattern
			bestStars = stars
			bestLiteral = literal
		}
	}
	return bestMessage, bestPattern != ""
}

// lookupWildcardValue 查找匹配路径的通配符值
func lookupWildcardValue(values map[string]string, path string) (string, bool) {
	normalized := normalizePath(path)
	var (
		bestValue   string
		bestPattern string
		bestStars   = int(^uint(0) >> 1)
		bestLiteral int
	)
	for pattern, value := range values {
		if value == "" || !strings.Contains(pattern, "*") || !matchWildcard(pattern, normalized) {
			continue
		}
		stars, literal := wildcardSpecificity(pattern)
		if stars < bestStars ||
			stars == bestStars && literal > bestLiteral ||
			stars == bestStars && literal == bestLiteral && pattern > bestPattern {
			bestValue = value
			bestPattern = pattern
			bestStars = stars
			bestLiteral = literal
		}
	}
	return bestValue, bestPattern != ""
}

// normalizePath 将索引路径转换为通配符路径
func normalizePath(path string) string {
	var builder strings.Builder
	for index := 0; index < len(path); {
		if path[index] == '[' {
			closeIndex := strings.IndexByte(path[index:], ']')
			if closeIndex >= 0 {
				builder.WriteString(".*")
				index += closeIndex + 1
				continue
			}
		}
		builder.WriteByte(path[index])
		index++
	}
	return builder.String()
}

// matchWildcard 判断通配符路径是否匹配实际路径
func matchWildcard(pattern, actual string) bool {
	return matchWildcardSegments(strings.Split(pattern, "."), strings.Split(actual, "."))
}

// matchWildcardSegments 按路径段匹配通配符
func matchWildcardSegments(pattern, actual []string) bool {
	if len(pattern) == 0 {
		return len(actual) == 0
	}
	if pattern[0] == "**" {
		if matchWildcardSegments(pattern[1:], actual) {
			return true
		}
		return len(actual) > 0 && matchWildcardSegments(pattern, actual[1:])
	}
	if len(actual) == 0 {
		return false
	}
	if pattern[0] != "*" && pattern[0] != actual[0] {
		return false
	}
	return matchWildcardSegments(pattern[1:], actual[1:])
}

// wildcardSpecificity 计算通配符和固定字符数量
func wildcardSpecificity(pattern string) (int, int) {
	stars := 0
	literal := 0
	for segment := range strings.SplitSeq(pattern, ".") {
		switch segment {
		case "**":
			stars += 2
		case "*":
			stars++
		default:
			literal += len(segment)
		}
	}
	return stars, literal
}
