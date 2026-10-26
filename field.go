package validatex

import (
	"strings"
	"unicode/utf8"
)

// Integer 定义内置整数类型约束
type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// baseField 通用字段校验器
type baseField[T any] struct {
	path       string
	value      T
	translator *Translator
	errs       Errors
}

// newBaseField 创建通用字段校验器
func newBaseField[T any](translator *Translator, path string, value T) baseField[T] {
	return baseField[T]{translator: translator, path: path, value: value}
}

// check 按条件添加字段错误
func (f *baseField[T]) check(rule string, invalid bool, args ...any) {
	f.checkMessage(rule, "", invalid, args...)
}

// checkMessage 按条件和消息添加字段错误
func (f *baseField[T]) checkMessage(rule, message string, invalid bool, args ...any) {
	f.checkPath(f.path, rule, message, invalid, args...)
}

// checkPath 按指定路径添加字段错误
func (f *baseField[T]) checkPath(path, rule, message string, invalid bool, args ...any) {
	if !invalid {
		return
	}
	f.addError(path, displayField(f.translator, path), rule, message, args...)
}

// checkPathDisplay 按指定路径和字段显示名添加字段错误
func (f *baseField[T]) checkPathDisplay(path, display, rule, message string, args ...any) {
	f.addError(path, display, rule, message, args...)
}

// addError 添加已经确定无效的字段错误
func (f *baseField[T]) addError(path, display, rule, message string, args ...any) {
	template := message
	source := sourceTemplate
	if message == "" {
		source = sourceRule
		message = translateDisplay(f.translator, path, display, rule, args...)
	} else {
		message = formatMessage(message, path, display, rule, args...)
	}
	fieldError := FieldError{
		Path:     path,
		Rule:     rule,
		Msg:      message,
		template: template,
		source:   source,
	}
	if len(args) == 1 {
		fieldError.arg = args[0]
		fieldError.hasArg = true
	} else if len(args) > 1 {
		fieldError.args = append([]any(nil), args...)
	}
	f.errs = append(f.errs, fieldError)
}

// messageArg 定义带名称的消息占位符参数
type messageArg struct {
	name  string
	value any
}

// arg 创建带名称的消息占位符参数
func arg(name string, value any) messageArg {
	return messageArg{name: name, value: value}
}

// custom 执行自定义校验规则
func (f *baseField[T]) custom(rule, message string, valid func(T) bool, args ...any) {
	f.checkMessage(rule, message, valid == nil || !valid(f.value), args...)
}

// Errors 返回字段错误
func (f *baseField[T]) Errors() Errors {
	return f.errs
}

// StringField 字符串字段校验器
type StringField struct {
	baseField[string]
}

// Str 创建字符串字段校验器
func Str(translator *Translator, path, value string) *StringField {
	return &StringField{baseField: newBaseField(translator, path, value)}
}

// Required 校验字符串必填
func (f *StringField) Required() *StringField {
	f.check("required", strings.TrimSpace(f.value) == "")
	return f
}

// Min 校验字符串最小字符数
func (f *StringField) Min(n int) *StringField {
	f.check("min", f.value != "" && n > 0 && utf8.RuneCountInString(f.value) < n, n)
	return f
}

// Email 校验邮箱格式
func (f *StringField) Email() *StringField {
	f.check("email", f.value != "" && !isEmail(f.value))
	return f
}

// Rule 添加字符串自定义校验规则
func (f *StringField) Rule(rule, message string, valid func(string) bool, args ...any) *StringField {
	f.custom(rule, message, valid, args...)
	return f
}

// IntField 整数字段校验器
type IntField[T Integer] struct {
	baseField[T]
}

// Int 创建整数字段校验器
func Int[T Integer](translator *Translator, path string, value T) *IntField[T] {
	return &IntField[T]{baseField: newBaseField(translator, path, value)}
}

// Min 校验整数最小值
func (f *IntField[T]) Min(n T) *IntField[T] {
	f.check("min", f.value < n, n)
	return f
}

// Rule 添加整数自定义校验规则
func (f *IntField[T]) Rule(rule, message string, valid func(T) bool, args ...any) *IntField[T] {
	f.custom(rule, message, valid, args...)
	return f
}
