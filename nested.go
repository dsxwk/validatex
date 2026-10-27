package validatex

import (
	"slices"
	"strconv"
	"strings"
)

// NestedOption 定义嵌套指针校验选项
type NestedOption uint8

const (
	nestedRequired NestedOption = 1 << iota
	nestedRequiredElements
)

// WithRequired 要求嵌套指针或容器必填
func WithRequired() NestedOption {
	return nestedRequired
}

// WithRequiredElements 要求切片或Map元素不能为nil
func WithRequiredElements() NestedOption {
	return nestedRequiredElements
}

// parseNestedOptions 解析嵌套指针校验选项
func parseNestedOptions(options []NestedOption) (required, requiredElements bool) {
	for _, option := range options {
		required = required || option&nestedRequired != 0
		requiredElements = requiredElements || option&nestedRequiredElements != 0
	}
	return required, requiredElements
}

// nestedRequiredErrors 创建嵌套字段必填错误
func nestedRequiredErrors[V any](translator *Translator, path string, value V) Errors {
	field := newBaseField(translator, path, value)
	field.check("required", true)
	return field.errs
}

// Struct 校验非nil的嵌套结构体值,指针请使用StructPointer
func Struct[T Validatable](translator *Translator, path string, value T) Errors {
	if any(value) == nil {
		return nil
	}
	return relocateErrors(translator, path, value.Validate(translator))
}

// StructPointer 校验嵌套结构体指针
func StructPointer[T any, P ValidatablePointer[T]](translator *Translator, path string, value P, options ...NestedOption) Errors {
	required, _ := parseNestedOptions(options)
	if value == nil {
		if required {
			return nestedRequiredErrors(translator, path, value)
		}
		return nil
	}
	return relocateErrors(translator, path, validatePointer(translator, value))
}

// SliceStruct 校验非nil结构体切片,指针元素请使用SliceStructPointer
func SliceStruct[T Validatable](translator *Translator, path string, values []T) Errors {
	var errs Errors
	for index, value := range values {
		if any(value) == nil {
			continue
		}
		childErrors := value.Validate(translator)
		if len(childErrors) == 0 {
			continue
		}
		if errs == nil {
			errs = make(Errors, 0, len(values))
		}
		prefix := path + "[" + strconv.Itoa(index) + "]"
		errs = append(errs, relocateErrors(translator, prefix, childErrors)...)
	}
	return errs
}

// SliceStructPointer 校验结构体指针切片
func SliceStructPointer[T any, P ValidatablePointer[T]](translator *Translator, path string, values []P, options ...NestedOption) Errors {
	required, requiredElements := parseNestedOptions(options)
	if required && len(values) == 0 {
		return nestedRequiredErrors(translator, path, values)
	}
	var errs Errors
	for index, value := range values {
		if value == nil {
			if requiredElements {
				itemPath := path + "[" + strconv.Itoa(index) + "]"
				errs = append(errs, nestedRequiredErrors(translator, itemPath, value)...)
			}
			continue
		}
		childErrors := validatePointer(translator, value)
		if len(childErrors) == 0 {
			continue
		}
		if errs == nil {
			errs = make(Errors, 0, len(values))
		}
		prefix := path + "[" + strconv.Itoa(index) + "]"
		errs = append(errs, relocateErrors(translator, prefix, childErrors)...)
	}
	return errs
}

// MapStruct 校验非nil Map中的结构体值,指针值请使用MapStructPointer
func MapStruct[K comparable, T Validatable](translator *Translator, path string, values map[K]T) Errors {
	var errs Errors
	for key, value := range values {
		if any(value) == nil {
			continue
		}
		childErrors := value.Validate(translator)
		if len(childErrors) == 0 {
			continue
		}
		if errs == nil {
			errs = make(Errors, 0, len(values))
		}
		keyText := formatMapKey(key)
		prefix := joinMapPath(path, keyText)
		errs = append(errs, relocateErrorsWithDisplay(translator, prefix, displayMapKeyText(keyText), childErrors)...)
	}
	sortErrorsByPath(errs)
	return errs
}

// MapStructPointer 校验Map中的结构体指针
func MapStructPointer[K comparable, T any, P ValidatablePointer[T]](translator *Translator, path string, values map[K]P, options ...NestedOption) Errors {
	required, requiredElements := parseNestedOptions(options)
	if required && len(values) == 0 {
		return nestedRequiredErrors(translator, path, values)
	}
	var errs Errors
	for key, value := range values {
		if value == nil {
			if requiredElements {
				keyText := formatMapKey(key)
				itemPath := joinMapPath(path, keyText)
				field := newBaseField(translator, itemPath, value)
				field.checkPathDisplay(itemPath, displayMapKeyText(keyText), "required", "")
				errs = append(errs, field.errs...)
			}
			continue
		}
		childErrors := validatePointer(translator, value)
		if len(childErrors) == 0 {
			continue
		}
		if errs == nil {
			errs = make(Errors, 0, len(values))
		}
		keyText := formatMapKey(key)
		prefix := joinMapPath(path, keyText)
		errs = append(errs, relocateErrorsWithDisplay(translator, prefix, displayMapKeyText(keyText), childErrors)...)
	}
	sortErrorsByPath(errs)
	return errs
}

// validatePointer 执行指针请求校验
func validatePointer[T any, P ValidatablePointer[T]](translator *Translator, value P) Errors {
	return value.Validate(translator)
}

// relocateErrors 重定位嵌套错误路径并重新翻译
func relocateErrors(translator *Translator, prefix string, values Errors) Errors {
	return relocateErrorsWithDisplay(translator, prefix, "", values)
}

// relocateErrorsWithDisplay 重定位嵌套错误路径,并为直接挂在容器键上的错误提供显示名
func relocateErrorsWithDisplay(translator *Translator, prefix, prefixDisplay string, values Errors) Errors {
	if len(values) == 0 {
		return nil
	}
	errs := make(Errors, 0, len(values))
	for _, value := range values {
		path := joinPath(prefix, value.Path)
		fallback := ""
		if value.Path == "" {
			fallback = prefixDisplay
		}
		display := displayFieldWithFallback(translator, path, fallback)
		message := relocatedMessage(translator, path, display, value)
		errs = append(errs, FieldError{
			Path:     path,
			Rule:     value.Rule,
			Msg:      message,
			template: value.template,
			source:   value.source,
			arg:      value.arg,
			args:     value.args,
			hasArg:   value.hasArg,
		})
	}
	return errs
}

// relocatedMessage 按新路径重新生成错误消息
func relocatedMessage(translator *Translator, path, display string, value FieldError) string {
	switch value.source {
	case sourceRule:
		if len(value.args) > 0 {
			return translateDisplay(translator, path, display, value.Rule, value.args...)
		}
		if value.hasArg {
			return translateDisplay(translator, path, display, value.Rule, value.arg)
		}
		return translateDisplay(translator, path, display, value.Rule)
	case sourceTemplate:
		if len(value.args) > 0 {
			return formatMessage(value.template, path, display, value.Rule, value.args...)
		}
		if value.hasArg {
			return formatMessage(value.template, path, display, value.Rule, value.arg)
		}
		return formatMessage(value.template, path, display, value.Rule)
	default:
		if value.Msg == "" && value.Rule != "" {
			return translateDisplay(translator, path, display, value.Rule)
		}
		return value.Msg
	}
}

// joinPath 拼接字段路径
func joinPath(prefix, path string) string {
	if prefix == "" {
		return path
	}
	if path == "" {
		return prefix
	}
	return prefix + "." + path
}

const (
	emptyMapKey      = "%empty"
	hexadecimalChars = "0123456789ABCDEF"
)

// joinMapPath 拼接Map路径,并转义可能破坏路径结构的key
func joinMapPath(prefix, key string) string {
	return joinPath(prefix, encodeMapKey(key))
}

// encodeMapKey 对Map key中的路径保留字符进行百分号编码
func encodeMapKey(key string) string {
	if key == "" {
		return emptyMapKey
	}
	if !strings.ContainsAny(key, ".[]%*") {
		return key
	}

	var builder strings.Builder
	builder.Grow(len(key))
	for index := 0; index < len(key); index++ {
		character := key[index]
		switch character {
		case '.', '[', ']', '%', '*':
			builder.WriteByte('%')
			builder.WriteByte(hexadecimalChars[character>>4])
			builder.WriteByte(hexadecimalChars[character&0x0F])
		default:
			builder.WriteByte(character)
		}
	}
	return builder.String()
}

// sortErrorsByPath 按路径稳定排序错误
func sortErrorsByPath(errs Errors) {
	if len(errs) < 2 {
		return
	}
	slices.SortStableFunc(errs, func(left, right FieldError) int {
		return strings.Compare(left.Path, right.Path)
	})
}
