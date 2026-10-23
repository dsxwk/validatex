package validatex

// When 条件成立时执行校验
func When(condition bool, validate func() Errors) Errors {
	if !condition || validate == nil {
		return nil
	}
	return validate()
}

// Nullable 值不为空时执行校验
func Nullable[T any](value *T, validate func(T) Errors) Errors {
	if value == nil || validate == nil {
		return nil
	}
	return validate(*value)
}
