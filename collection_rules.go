package validatex

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

// mapKey 约束Map键为可比较类型
type mapKey interface {
	comparable
}

// formatMapKey 生成Map键的稳定文本
func formatMapKey[K mapKey](key K) string {
	var zero K
	if any(zero) == nil {
		return fmt.Sprintf("%T:%#v", key, key)
	}
	switch value := any(key).(type) {
	case string:
		return value
	case bool:
		return strconv.FormatBool(value)
	case int:
		return strconv.Itoa(value)
	case int8:
		return strconv.FormatInt(int64(value), 10)
	case int16:
		return strconv.FormatInt(int64(value), 10)
	case int32:
		return strconv.FormatInt(int64(value), 10)
	case int64:
		return strconv.FormatInt(value, 10)
	case uint:
		return strconv.FormatUint(uint64(value), 10)
	case uint8:
		return strconv.FormatUint(uint64(value), 10)
	case uint16:
		return strconv.FormatUint(uint64(value), 10)
	case uint32:
		return strconv.FormatUint(uint64(value), 10)
	case uint64:
		return strconv.FormatUint(value, 10)
	case uintptr:
		return strconv.FormatUint(uint64(value), 10)
	case float32:
		return strconv.FormatFloat(float64(value), 'g', -1, 32)
	case float64:
		return strconv.FormatFloat(value, 'g', -1, 64)
	case complex64:
		return strconv.FormatComplex(complex128(value), 'g', -1, 64)
	case complex128:
		return strconv.FormatComplex(value, 'g', -1, 128)
	default:
		return fmt.Sprintf("%#v", key)
	}
}

// displayMapKeyText 将空键转换为固定展示名
func displayMapKeyText(value string) string {
	if value == "" {
		return emptyMapKey
	}
	return value
}

// joinMapKeys 格式化Map键列表用于消息占位符
func joinMapKeys[K mapKey](keys []K) string {
	var builder strings.Builder
	for index, key := range keys {
		if index > 0 {
			builder.WriteByte(',')
		}
		builder.WriteString(displayMapKeyText(formatMapKey(key)))
	}
	return builder.String()
}

// BoolField 布尔字段校验器
type BoolField struct {
	baseField[bool]
}

// Bool 创建布尔字段校验器
func Bool(translator *Translator, path string, value bool) *BoolField {
	return &BoolField{baseField: newBaseField(translator, path, value)}
}

// Accepted 校验布尔值为真
func (f *BoolField) Accepted() *BoolField {
	f.check("accepted", !f.value)
	return f
}

// Declined 校验布尔值为假
func (f *BoolField) Declined() *BoolField {
	f.check("declined", f.value)
	return f
}

// AcceptedIf 条件成立时校验布尔值为真
func (f *BoolField) AcceptedIf(condition bool) *BoolField {
	f.check("accepted_if", condition && !f.value)
	return f
}

// DeclinedIf 条件成立时校验布尔值为假
func (f *BoolField) DeclinedIf(condition bool) *BoolField {
	f.check("declined_if", condition && f.value)
	return f
}

// Prohibited 校验布尔值为假
func (f *BoolField) Prohibited() *BoolField {
	f.check("prohibited", f.value)
	return f
}

// ProhibitedIf 条件成立时校验布尔值为假
func (f *BoolField) ProhibitedIf(condition bool) *BoolField {
	f.check("prohibited_if", condition && f.value)
	return f
}

// ProhibitedUnless 条件不成立时校验布尔值为假
func (f *BoolField) ProhibitedUnless(condition bool) *BoolField {
	f.check("prohibited_unless", !condition && f.value)
	return f
}

// True 校验布尔值为真
func (f *BoolField) True() *BoolField {
	f.check("true", !f.value)
	return f
}

// False 校验布尔值为假
func (f *BoolField) False() *BoolField {
	f.check("false", f.value)
	return f
}

// Same 校验布尔值与目标值相同
func (f *BoolField) Same(value bool) *BoolField {
	f.check("same", f.value != value, value)
	return f
}

// Different 校验布尔值与目标值不同
func (f *BoolField) Different(value bool) *BoolField {
	f.check("different", f.value == value, value)
	return f
}

// Rule 添加布尔自定义校验规则
func (f *BoolField) Rule(rule, message string, valid func(bool) bool, args ...any) *BoolField {
	f.custom(rule, message, valid, args...)
	return f
}

// SliceField 切片字段校验器
type SliceField[T comparable] struct {
	baseField[[]T]
}

// Slice 创建切片字段校验器
func Slice[T comparable](translator *Translator, path string, value []T) *SliceField[T] {
	return &SliceField[T]{baseField: newBaseField(translator, path, value)}
}

// Required 校验切片必填
func (f *SliceField[T]) Required() *SliceField[T] {
	f.check("required", len(f.value) == 0)
	return f
}

// RequiredIf 条件成立时校验切片必填
func (f *SliceField[T]) RequiredIf(condition bool) *SliceField[T] {
	f.check("required_if", condition && len(f.value) == 0)
	return f
}

// RequiredUnless 条件不成立时校验切片必填
func (f *SliceField[T]) RequiredUnless(condition bool) *SliceField[T] {
	f.check("required_unless", !condition && len(f.value) == 0)
	return f
}

// Prohibited 校验切片为空
func (f *SliceField[T]) Prohibited() *SliceField[T] {
	f.check("prohibited", len(f.value) != 0)
	return f
}

// ProhibitedIf 条件成立时校验切片为空
func (f *SliceField[T]) ProhibitedIf(condition bool) *SliceField[T] {
	f.check("prohibited_if", condition && len(f.value) != 0)
	return f
}

// ProhibitedUnless 条件不成立时校验切片为空
func (f *SliceField[T]) ProhibitedUnless(condition bool) *SliceField[T] {
	f.check("prohibited_unless", !condition && len(f.value) != 0)
	return f
}

// Min 校验切片最小元素数
func (f *SliceField[T]) Min(n int) *SliceField[T] {
	f.check("min", len(f.value) < n, n)
	return f
}

// Max 校验切片最大元素数
func (f *SliceField[T]) Max(n int) *SliceField[T] {
	f.check("max", len(f.value) > n, n)
	return f
}

// Between 校验切片元素数范围
func (f *SliceField[T]) Between(min, max int) *SliceField[T] {
	if len(f.value) < min || len(f.value) > max {
		f.check("between", true, arg("min", min), arg("max", max))
	}
	return f
}

// Size 校验切片元素数等于指定值
func (f *SliceField[T]) Size(n int) *SliceField[T] {
	f.check("size", len(f.value) != n, n)
	return f
}

// Distinct 校验切片元素不重复
func (f *SliceField[T]) Distinct() *SliceField[T] {
	if len(f.value) < 2 {
		return f
	}
	seen := make(map[T]struct{}, len(f.value))
	invalid := false
	for _, value := range f.value {
		if _, exists := seen[value]; exists {
			invalid = true
			break
		}
		seen[value] = struct{}{}
	}
	f.check("distinct", invalid)
	return f
}

// Contains 校验切片包含指定元素
func (f *SliceField[T]) Contains(value T) *SliceField[T] {
	f.check("contains", !slices.Contains(f.value, value), value)
	return f
}

// DoesntContain 校验切片不包含指定元素
func (f *SliceField[T]) DoesntContain(value T) *SliceField[T] {
	f.check("doesnt_contain", slices.Contains(f.value, value), value)
	return f
}

// Each 校验每个切片元素
func (f *SliceField[T]) Each(rule, message string, valid func(T) bool, args ...any) *SliceField[T] {
	for index, value := range f.value {
		if valid != nil && valid(value) {
			continue
		}
		path := f.path + "[" + strconv.Itoa(index) + "]"
		f.checkPath(path, rule, message, true, args...)
	}
	return f
}

// Rule 添加切片自定义校验规则
func (f *SliceField[T]) Rule(rule, message string, valid func([]T) bool, args ...any) *SliceField[T] {
	f.custom(rule, message, valid, args...)
	return f
}

// MapField Map字段校验器
type MapField[K comparable, V any] struct {
	baseField[map[K]V]
}

// Map 创建Map字段校验器
func Map[K comparable, V any](translator *Translator, path string, value map[K]V) *MapField[K, V] {
	return &MapField[K, V]{baseField: newBaseField(translator, path, value)}
}

// Required 校验Map必填
func (f *MapField[K, V]) Required() *MapField[K, V] {
	f.check("required", len(f.value) == 0)
	return f
}

// RequiredIf 条件成立时校验Map必填
func (f *MapField[K, V]) RequiredIf(condition bool) *MapField[K, V] {
	f.check("required_if", condition && len(f.value) == 0)
	return f
}

// RequiredUnless 条件不成立时校验Map必填
func (f *MapField[K, V]) RequiredUnless(condition bool) *MapField[K, V] {
	f.check("required_unless", !condition && len(f.value) == 0)
	return f
}

// Prohibited 校验Map为空
func (f *MapField[K, V]) Prohibited() *MapField[K, V] {
	f.check("prohibited", len(f.value) != 0)
	return f
}

// ProhibitedIf 条件成立时校验Map为空
func (f *MapField[K, V]) ProhibitedIf(condition bool) *MapField[K, V] {
	f.check("prohibited_if", condition && len(f.value) != 0)
	return f
}

// ProhibitedUnless 条件不成立时校验Map为空
func (f *MapField[K, V]) ProhibitedUnless(condition bool) *MapField[K, V] {
	f.check("prohibited_unless", !condition && len(f.value) != 0)
	return f
}

// Min 校验Map最小键值数
func (f *MapField[K, V]) Min(n int) *MapField[K, V] {
	f.check("min", len(f.value) < n, n)
	return f
}

// Max 校验Map最大键值数
func (f *MapField[K, V]) Max(n int) *MapField[K, V] {
	f.check("max", len(f.value) > n, n)
	return f
}

// Between 校验Map键值数范围
func (f *MapField[K, V]) Between(min, max int) *MapField[K, V] {
	if len(f.value) < min || len(f.value) > max {
		f.check("between", true, arg("min", min), arg("max", max))
	}
	return f
}

// Size 校验Map键值数等于指定值
func (f *MapField[K, V]) Size(n int) *MapField[K, V] {
	f.check("size", len(f.value) != n, n)
	return f
}

// HasKeys 校验Map包含全部指定键
func (f *MapField[K, V]) HasKeys(keys ...K) *MapField[K, V] {
	invalid := false
	for _, key := range keys {
		if _, exists := f.value[key]; !exists {
			invalid = true
			break
		}
	}
	if invalid {
		f.check("has_keys", true, arg("keys", joinMapKeys(keys)))
	}
	return f
}

// AllowedKeys 校验Map仅包含指定键
func (f *MapField[K, V]) AllowedKeys(keys ...K) *MapField[K, V] {
	invalid := false
	for key := range f.value {
		allowed := slices.Contains(keys, key)
		if !allowed {
			invalid = true
			break
		}
	}
	if invalid {
		f.check("allowed_keys", true, arg("keys", joinMapKeys(keys)))
	}
	return f
}

// Each 校验每个Map键值
func (f *MapField[K, V]) Each(rule, message string, valid func(K, V) bool, args ...any) *MapField[K, V] {
	start := len(f.errs)
	for key, value := range f.value {
		if valid != nil && valid(key, value) {
			continue
		}
		keyText := formatMapKey(key)
		encodedKey := encodeMapKey(keyText)
		path := joinPath(f.path, encodedKey)
		display := displayFieldWithFallback(f.translator, path, displayMapKeyText(keyText))
		f.checkPathDisplay(path, display, rule, message, args...)
	}
	sortErrorsByPath(f.errs[start:])
	return f
}

// Rule 添加Map自定义校验规则
func (f *MapField[K, V]) Rule(rule, message string, valid func(map[K]V) bool, args ...any) *MapField[K, V] {
	f.custom(rule, message, valid, args...)
	return f
}
