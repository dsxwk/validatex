// Package validatex 提供泛型字段校验,场景校验,翻译消息和通用错误响应
package validatex

import (
	"fmt"
	"strings"
)

// Messages 规则翻译表
type Messages map[string]string

// Translates 字段翻译表
type Translates map[string]string

// Translator 校验翻译器
type Translator struct {
	Messages   Messages
	Translates Translates
}

// FieldError 单条字段错误
type FieldError struct {
	Path     string
	Rule     string
	Msg      string
	template string
	arg      any
	args     []any
	argCount int
}

// Errors 字段错误集合
type Errors []FieldError

// Error 返回所有错误的文本
func (e Errors) Error() string {
	items := make([]string, 0, len(e))
	for _, item := range e {
		items = append(items, fmt.Sprintf("[%s: %s]", item.Path, item.Msg))
	}
	return strings.Join(items, " ")
}

// ToMap 按字段路径聚合错误消息
func (e Errors) ToMap() map[string][]string {
	result := make(map[string][]string, len(e))
	for _, item := range e {
		result[item.Path] = append(result[item.Path], item.Msg)
	}
	return result
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

// Check 聚合多组字段错误
func Check(fields ...Errors) Errors {
	total := 0
	var only Errors
	for _, field := range fields {
		if len(field) == 0 {
			continue
		}
		total += len(field)
		only = field
	}
	if total == 0 {
		return nil
	}
	if total == len(only) {
		return only
	}
	all := make(Errors, 0, total)
	for _, field := range fields {
		all = append(all, field...)
	}
	return all
}

// arguments 返回错误参数
func (e FieldError) arguments() []any {
	if e.argCount == 0 {
		return nil
	}
	if e.argCount == 1 {
		return []any{e.arg}
	}
	return e.args
}
