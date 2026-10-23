package validatex

import (
	"strings"
	"unicode/utf8"
)

// integer 定义整数类型约束
type integer interface {
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
	template := message
	if message == "" {
		message = translate(f.translator, path, rule, args...)
	} else {
		message = formatMessage(f.translator, message, path, rule, args...)
	}
	fieldError := FieldError{
		Path:     path,
		Rule:     rule,
		Msg:      message,
		template: template,
		argCount: len(args),
	}
	if len(args) == 1 {
		fieldError.arg = args[0]
	} else if len(args) > 1 {
		fieldError.args = append([]any(nil), args...)
	}
	f.errs = append(f.errs, fieldError)
}

// custom 执行自定义校验规则
func (f *baseField[T]) custom(rule, message string, valid func(T) bool, args ...any) {
	f.checkMessage(rule, message, !valid(f.value), args...)
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
	f.check("email", f.value != "" && !emailPattern.MatchString(f.value))
	return f
}

// Rule 添加字符串自定义校验规则
func (f *StringField) Rule(rule, message string, valid func(string) bool, args ...any) *StringField {
	f.custom(rule, message, valid, args...)
	return f
}

// IntField 整数字段校验器
type IntField[T integer] struct {
	baseField[T]
}

// Int 创建整数字段校验器
func Int[T integer](translator *Translator, path string, value T) *IntField[T] {
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
