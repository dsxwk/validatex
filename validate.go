// Package validatex 提供带本地化错误路径的泛型校验。
package validatex

import (
	"strconv"
	"strings"
)

// Messages 规则翻译表
type Messages map[string]string

// Translates 字段翻译表
type Translates map[string]string

// Translator 校验翻译器 校验期间Messages和Translates应视为只读
// 需要通配符预编译或并发只读快照时使用Compile
type Translator struct {
	Messages   Messages
	Translates Translates

	compiled           bool
	wildcardMessages   []compiledMessagePattern
	wildcardTranslates []compiledTranslatePattern
}

// FieldError 单条字段错误
type FieldError struct {
	Path     string
	Rule     string
	Msg      string
	template string
	arg      any
	args     []any
	source   messageSource
	hasArg   bool
}

type messageSource uint8

const (
	_ messageSource = iota // zero value marks already resolved messages
	sourceRule
	sourceTemplate
)

// Errors 字段错误集合
type Errors []FieldError

// Error 返回所有错误的文本
func (e Errors) Error() string {
	if len(e) == 0 {
		return ""
	}

	size := len(e)*5 - 1
	for _, item := range e {
		size += len(item.Path) + len(item.text())
	}

	var builder strings.Builder
	builder.Grow(size)
	for index, item := range e {
		if index > 0 {
			builder.WriteByte(' ')
		}
		builder.WriteByte('[')
		builder.WriteString(item.Path)
		builder.WriteString(": ")
		builder.WriteString(item.text())
		builder.WriteByte(']')
	}
	return builder.String()
}

// ToMap 按字段路径聚合错误消息
func (e Errors) ToMap() map[string][]string {
	result := make(map[string][]string, len(e))
	for _, item := range e {
		result[item.Path] = append(result[item.Path], item.text())
	}
	return result
}

// text 返回错误文本,缺失消息时回退到规则名
func (e FieldError) text() string {
	if e.Msg != "" {
		return e.Msg
	}
	return e.Rule
}

// Indexes 返回错误路径中的切片下标
func (e FieldError) Indexes() []int {
	var indexes []int
	for index := 0; index < len(e.Path); index++ {
		if e.Path[index] != '[' {
			continue
		}
		offset := strings.IndexByte(e.Path[index+1:], ']')
		if offset < 0 {
			break
		}
		closing := index + 1 + offset
		value, err := strconv.Atoi(e.Path[index+1 : closing])
		if err == nil && value >= 0 {
			indexes = append(indexes, value)
		}
		index = closing
	}
	return indexes
}

// Validatable 定义普通校验接口
type Validatable interface {
	Validate(translator *Translator) Errors
}

// SceneValidatable 定义场景校验接口
type SceneValidatable interface {
	ValidateWithScene(scene string, translator *Translator) Errors
}

// TranslatesProvider 定义字段翻译提供接口
type TranslatesProvider interface {
	Translates() Translates
}

// MessagesProvider 定义规则翻译提供接口
type MessagesProvider interface {
	Messages() Messages
}

// NewTranslator 创建校验翻译器
func NewTranslator() *Translator {
	return &Translator{
		Messages:   make(Messages),
		Translates: make(Translates),
	}
}

// Check 聚合多组字段错误 仅一组非空时直接返回该切片,调用方应视为只读
func Check(fields ...Errors) Errors {
	total := 0
	nonEmpty := 0
	var last Errors
	for _, field := range fields {
		if len(field) == 0 {
			continue
		}
		total += len(field)
		nonEmpty++
		last = field
	}
	if total == 0 {
		return nil
	}
	if nonEmpty == 1 {
		return last
	}
	all := make(Errors, 0, total)
	for _, field := range fields {
		all = append(all, field...)
	}
	return all
}
