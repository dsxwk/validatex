package validatex

// Fields 绑定翻译器的字段构造器
type Fields struct {
	translator *Translator
}

// NewFields 创建绑定翻译器的字段构造器
func NewFields(translator *Translator) Fields {
	return Fields{translator: translator}
}

// Str 创建字符串校验字段
func (f Fields) Str(path, value string) *StringField {
	return Str(f.translator, path, value)
}

// Int 创建整数校验字段
func (f Fields) Int[T Integer](path string, value T) *IntField[T] {
	return Int(f.translator, path, value)
}

// Float 创建浮点数校验字段
func (f Fields) Float[T FloatValue](path string, value T) *FloatField[T] {
	return Float(f.translator, path, value)
}

// Bool 创建布尔校验字段
func (f Fields) Bool(path string, value bool) *BoolField {
	return Bool(f.translator, path, value)
}

// Slice 创建切片校验字段
func (f Fields) Slice[T comparable](path string, value []T) *SliceField[T] {
	return Slice(f.translator, path, value)
}

// Map 创建Map校验字段
func (f Fields) Map[K comparable, V any](path string, value map[K]V) *MapField[K, V] {
	return Map(f.translator, path, value)
}

// Struct 校验嵌套结构体值
func (f Fields) Struct[T Validatable](path string, value T) Errors {
	return Struct(f.translator, path, value)
}

// SliceStruct 校验结构体切片
func (f Fields) SliceStruct[T Validatable](path string, values []T) Errors {
	return SliceStruct(f.translator, path, values)
}

// MapStruct 校验结构体Map
func (f Fields) MapStruct[K comparable, T Validatable](path string, values map[K]T) Errors {
	return MapStruct(f.translator, path, values)
}

// StructPointer 校验嵌套结构体指针
func (f Fields) StructPointer[T any, P ValidatablePointer[T]](path string, value P, options ...NestedOption) Errors {
	return StructPointer(f.translator, path, value, options...)
}

// SliceStructPointer 校验结构体指针切片
func (f Fields) SliceStructPointer[T any, P ValidatablePointer[T]](path string, values []P, options ...NestedOption) Errors {
	return SliceStructPointer(f.translator, path, values, options...)
}

// MapStructPointer 校验Map中的结构体指针
func (f Fields) MapStructPointer[K comparable, T any, P ValidatablePointer[T]](path string, values map[K]P, options ...NestedOption) Errors {
	return MapStructPointer(f.translator, path, values, options...)
}
