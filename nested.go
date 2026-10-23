package validatex

import "fmt"

// Struct 校验嵌套结构体
func Struct[T Validatable](translator *Translator, path string, value T) Errors {
	if any(value) == nil {
		return nil
	}
	return relocateErrors(translator, path, value.Validate(translator))
}

// SliceStruct 校验结构体切片
func SliceStruct[T Validatable](translator *Translator, path string, values []T) Errors {
	var errs Errors
	for index, value := range values {
		if any(value) == nil {
			continue
		}
		prefix := fmt.Sprintf("%s[%d]", path, index)
		errs = append(errs, relocateErrors(translator, prefix, value.Validate(translator))...)
	}
	return errs
}

// MapStruct 校验Map中的结构体
func MapStruct[K comparable, T Validatable](translator *Translator, path string, values map[K]T) Errors {
	var errs Errors
	for key, value := range values {
		if any(value) == nil {
			continue
		}
		prefix := joinPath(path, fmt.Sprint(key))
		errs = append(errs, relocateErrors(translator, prefix, value.Validate(translator))...)
	}
	return errs
}

// relocateErrors 重定位嵌套错误路径并重新翻译
func relocateErrors(translator *Translator, prefix string, values Errors) Errors {
	errs := make(Errors, 0, len(values))
	for _, value := range values {
		path := joinPath(prefix, value.Path)
		message := value.Msg
		if value.template != "" {
			message = formatMessage(translator, value.template, path, value.Rule, value.arguments()...)
		} else if translator != nil {
			message = translate(translator, path, value.Rule, value.arguments()...)
		}
		errs = append(errs, FieldError{
			Path:     path,
			Rule:     value.Rule,
			Msg:      message,
			template: value.template,
			arg:      value.arg,
			args:     append([]any(nil), value.args...),
			argCount: value.argCount,
		})
	}
	return errs
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
