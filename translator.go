package validatex

import (
	"fmt"
	"maps"
	"slices"
	"strconv"
	"strings"
)

type compiledPattern struct {
	pattern  string
	segments []string
	stars    int
	literal  int
}

type compiledMessagePattern struct {
	compiledPattern
	rule    string
	message string
}

type compiledTranslatePattern struct {
	compiledPattern
	value string
}

// Compile 返回带有预编译通配符索引的只读翻译器快照
func (t *Translator) Compile() *Translator {
	if t == nil {
		return nil
	}
	if t.compiled {
		return t
	}

	messages := maps.Clone(t.Messages)
	translates := maps.Clone(t.Translates)
	compiled := &Translator{
		Messages:   messages,
		Translates: translates,
		compiled:   true,
	}
	compiled.wildcardMessages = compileMessagePatterns(messages)
	compiled.wildcardTranslates = compileTranslatePatterns(translates)
	return compiled
}

// compileMessagePatterns 编译通配符消息规则
func compileMessagePatterns(messages Messages) []compiledMessagePattern {
	patterns := make([]compiledMessagePattern, 0, len(messages))
	for key, message := range messages {
		separator := strings.LastIndexByte(key, '.')
		if separator <= 0 || separator == len(key)-1 {
			continue
		}
		pattern, ok := newCompiledPattern(key[:separator])
		if !ok {
			continue
		}
		patterns = append(patterns, compiledMessagePattern{
			compiledPattern: pattern,
			rule:            key[separator+1:],
			message:         message,
		})
	}
	return patterns
}

// compileTranslatePatterns 编译通配符字段名称规则
func compileTranslatePatterns(translates Translates) []compiledTranslatePattern {
	patterns := make([]compiledTranslatePattern, 0, len(translates))
	for key, value := range translates {
		if value == "" {
			continue
		}
		pattern, ok := newCompiledPattern(key)
		if !ok {
			continue
		}
		patterns = append(patterns, compiledTranslatePattern{
			compiledPattern: pattern,
			value:           value,
		})
	}
	return patterns
}

// newCompiledPattern 预解析通配符路径和优先级
func newCompiledPattern(pattern string) (compiledPattern, bool) {
	if !strings.Contains(pattern, "*") {
		return compiledPattern{}, false
	}
	segments := strings.Split(pattern, ".")
	stars := 0
	literal := 0
	for _, segment := range segments {
		switch segment {
		case "**":
			stars += 2
		case "*":
			stars++
		default:
			literal += len(segment)
		}
	}
	return compiledPattern{
		pattern:  pattern,
		segments: segments,
		stars:    stars,
		literal:  literal,
	}, true
}

// translateDisplay 按字段显示名获取翻译后的错误消息
func translateDisplay(translator *Translator, path, display, rule string, args ...any) string {
	if translator != nil && len(translator.Messages) > 0 {
		if message, ok := translator.Messages[path+"."+rule]; ok {
			return formatMessage(message, path, display, rule, args...)
		}
		if message, ok := translator.wildcardMessage(path, rule); ok {
			return formatMessage(message, path, display, rule, args...)
		}
		if message, ok := translator.Messages[rule]; ok {
			return formatMessage(message, path, display, rule, args...)
		}
	}
	if display == "" {
		display = path
	}
	return display + " is " + rule
}

// formatMessage 一次性替换消息中的占位符
func formatMessage(message, path, display, rule string, args ...any) string {
	if !strings.Contains(message, "{") {
		return message
	}

	var builder strings.Builder
	builder.Grow(len(message))
	for start := 0; start < len(message); {
		open := strings.IndexByte(message[start:], '{')
		if open < 0 {
			builder.WriteString(message[start:])
			break
		}
		open += start
		builder.WriteString(message[start:open])

		closeOffset := strings.IndexByte(message[open+1:], '}')
		if closeOffset < 0 {
			builder.WriteString(message[open:])
			break
		}
		closeIndex := open + 1 + closeOffset
		name := message[open+1 : closeIndex]
		if value, ok := placeholderValue(name, path, display, rule, args); ok {
			builder.WriteString(value)
		} else {
			builder.WriteString(message[open : closeIndex+1])
		}
		start = closeIndex + 1
	}
	return builder.String()
}

// placeholderValue 返回占位符对应的文本
func placeholderValue(name, path, display, rule string, args []any) (string, bool) {
	switch name {
	case "field":
		return display, true
	case "path":
		return path, true
	case "rule":
		return rule, true
	}

	for _, arg := range args {
		switch value := arg.(type) {
		case messageArg:
			if value.name == name {
				return formatArgument(value.value), true
			}
		case map[string]any:
			if item, ok := value[name]; ok {
				return formatArgument(item), true
			}
		}
	}

	if index, err := strconv.Atoi(name); err == nil && index >= 0 && index < len(args) {
		return formatArgument(args[index]), true
	}

	if len(args) == 0 || isNamedArgument(args[0]) {
		return "", false
	}
	if name == "max" && len(args) > 1 {
		return formatArgument(args[1]), true
	}
	switch name {
	case "value", "min", "max", "len", "size":
		return formatArgument(args[0]), true
	default:
		return "", false
	}
}

// formatArgument 将消息参数转换为文本
func formatArgument(value any) string {
	switch typed := value.(type) {
	case messageArg:
		return formatArgument(typed.value)
	case []string:
		return strings.Join(typed, ",")
	default:
		return fmt.Sprint(value)
	}
}

// isNamedArgument 判断参数是否使用命名占位符
func isNamedArgument(value any) bool {
	switch value.(type) {
	case messageArg, map[string]any:
		return true
	default:
		return false
	}
}

// displayField 获取字段显示名称
func displayField(translator *Translator, path string) string {
	return displayFieldWithFallback(translator, path, "")
}

// displayFieldWithFallback 获取字段显示名并保留指定回退值
func displayFieldWithFallback(translator *Translator, path, fallback string) string {
	if translator != nil {
		if field, ok := translator.Translates[path]; ok && field != "" {
			return field
		}
		if field, ok := translator.wildcardTranslate(path); ok {
			return field
		}
		field := lastField(path)
		if translated, ok := translator.Translates[field]; ok && translated != "" {
			return translated
		}
	}
	if fallback != "" {
		return fallback
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

// wildcardMessage 查找匹配当前路径和规则的通配符消息
func (t *Translator) wildcardMessage(path, rule string) (string, bool) {
	if t == nil {
		return "", false
	}
	if t.compiled {
		return lookupCompiledWildcardMessage(t.wildcardMessages, path, rule)
	}
	return lookupWildcardMessage(t.Messages, path, rule)
}

// wildcardTranslate 查找匹配当前路径的通配符字段名称
func (t *Translator) wildcardTranslate(path string) (string, bool) {
	if t == nil {
		return "", false
	}
	if t.compiled {
		return lookupCompiledWildcardTranslate(t.wildcardTranslates, path)
	}
	return lookupWildcardValue(t.Translates, path)
}

// lookupCompiledWildcardMessage 在预编译消息规则中查找最佳匹配
func lookupCompiledWildcardMessage(patterns []compiledMessagePattern, path, rule string) (string, bool) {
	var actual []string
	bestStars := int(^uint(0) >> 1)
	bestLiteral := 0
	bestPattern := ""
	bestMessage := ""
	for _, candidate := range patterns {
		if candidate.rule != rule {
			continue
		}
		if actual == nil {
			actual = strings.Split(normalizePath(path), ".")
		}
		if !matchWildcardSegments(candidate.segments, actual) {
			continue
		}
		if moreSpecificWildcard(candidate.pattern, candidate.stars, candidate.literal, bestPattern, bestStars, bestLiteral) {
			bestStars = candidate.stars
			bestLiteral = candidate.literal
			bestPattern = candidate.pattern
			bestMessage = candidate.message
		}
	}
	return bestMessage, bestPattern != ""
}

// lookupCompiledWildcardTranslate 在预编译字段规则中查找最佳匹配
func lookupCompiledWildcardTranslate(patterns []compiledTranslatePattern, path string) (string, bool) {
	var actual []string
	bestStars := int(^uint(0) >> 1)
	bestLiteral := 0
	bestPattern := ""
	bestValue := ""
	for _, candidate := range patterns {
		if actual == nil {
			actual = strings.Split(normalizePath(path), ".")
		}
		if !matchWildcardSegments(candidate.segments, actual) {
			continue
		}
		if moreSpecificWildcard(candidate.pattern, candidate.stars, candidate.literal, bestPattern, bestStars, bestLiteral) {
			bestStars = candidate.stars
			bestLiteral = candidate.literal
			bestPattern = candidate.pattern
			bestValue = candidate.value
		}
	}
	return bestValue, bestPattern != ""
}

// lookupWildcardMessage 查找匹配路径的通配符消息
func lookupWildcardMessage(messages Messages, path, rule string) (string, bool) {
	suffix := "." + rule
	var (
		bestMessage string
		bestPattern string
		bestStars   = int(^uint(0) >> 1)
		bestLiteral int
		actual      []string
	)
	for key, message := range messages {
		if !strings.Contains(key, "*") || !strings.HasSuffix(key, suffix) {
			continue
		}
		pattern := strings.TrimSuffix(key, suffix)
		if !strings.Contains(pattern, "*") {
			continue
		}
		if actual == nil {
			actual = strings.Split(normalizePath(path), ".")
		}
		if !matchWildcardSegments(strings.Split(pattern, "."), actual) {
			continue
		}
		stars, literal := wildcardSpecificity(pattern)
		if moreSpecificWildcard(pattern, stars, literal, bestPattern, bestStars, bestLiteral) {
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
	var (
		bestValue   string
		bestPattern string
		bestStars   = int(^uint(0) >> 1)
		bestLiteral int
		actual      []string
	)
	for pattern, value := range values {
		if value == "" || !strings.Contains(pattern, "*") {
			continue
		}
		if actual == nil {
			actual = strings.Split(normalizePath(path), ".")
		}
		if !matchWildcardSegments(strings.Split(pattern, "."), actual) {
			continue
		}
		stars, literal := wildcardSpecificity(pattern)
		if moreSpecificWildcard(pattern, stars, literal, bestPattern, bestStars, bestLiteral) {
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
	if !strings.Contains(path, "[") {
		return path
	}
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

// matchWildcardSegments 按路径段匹配通配符
func matchWildcardSegments(pattern, actual []string) bool {
	hasDeepWildcard := slices.Contains(pattern, "**")

	if !hasDeepWildcard {
		if len(pattern) != len(actual) {
			return false
		}
		for index, segment := range pattern {
			if segment != "*" && segment != actual[index] {
				return false
			}
		}
		return true
	}

	previous := make([]bool, len(actual)+1)
	current := make([]bool, len(actual)+1)
	previous[0] = true
	for _, segment := range pattern {
		clear(current)
		if segment == "**" {
			current[0] = previous[0]
			for index := 1; index <= len(actual); index++ {
				current[index] = previous[index] || current[index-1]
			}
		} else {
			for index := 1; index <= len(actual); index++ {
				current[index] = previous[index-1] &&
					(segment == "*" || segment == actual[index-1])
			}
		}
		previous, current = current, previous
	}
	return previous[len(actual)]
}

// moreSpecificWildcard 判断候选通配符是否优先于当前结果
func moreSpecificWildcard(pattern string, stars, literal int, bestPattern string, bestStars, bestLiteral int) bool {
	return stars < bestStars ||
		stars == bestStars && literal > bestLiteral ||
		stars == bestStars && literal == bestLiteral && pattern > bestPattern
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
