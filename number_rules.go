package validatex

import (
	"fmt"
	"math"
	"math/big"
	"reflect"
	"slices"
	"strconv"
)

// FloatValue 定义内置浮点类型约束
type FloatValue interface {
	~float32 | ~float64
}

// Required 校验整数必填
func (f *IntField[T]) Required() *IntField[T] {
	var zero T
	f.check("required", f.value == zero)
	return f
}

// RequiredIf 条件成立时校验整数必填
func (f *IntField[T]) RequiredIf(condition bool) *IntField[T] {
	var zero T
	f.check("required_if", condition && f.value == zero)
	return f
}

// RequiredUnless 条件不成立时校验整数必填
func (f *IntField[T]) RequiredUnless(condition bool) *IntField[T] {
	var zero T
	f.check("required_unless", !condition && f.value == zero)
	return f
}

// Prohibited 校验整数等于零
func (f *IntField[T]) Prohibited() *IntField[T] {
	var zero T
	f.check("prohibited", f.value != zero)
	return f
}

// ProhibitedIf 条件成立时校验整数等于零
func (f *IntField[T]) ProhibitedIf(condition bool) *IntField[T] {
	var zero T
	f.check("prohibited_if", condition && f.value != zero)
	return f
}

// ProhibitedUnless 条件不成立时校验整数等于零
func (f *IntField[T]) ProhibitedUnless(condition bool) *IntField[T] {
	var zero T
	f.check("prohibited_unless", !condition && f.value != zero)
	return f
}

// Max 校验整数最大值
func (f *IntField[T]) Max(n T) *IntField[T] {
	f.check("max", f.value > n, n)
	return f
}

// Between 校验整数范围
func (f *IntField[T]) Between(min, max T) *IntField[T] {
	if f.value < min || f.value > max {
		f.check("between", true, arg("min", min), arg("max", max))
	}
	return f
}

// Size 校验整数等于指定值
func (f *IntField[T]) Size(n T) *IntField[T] {
	f.check("size", f.value != n, n)
	return f
}

// Gt 校验整数大于目标值
func (f *IntField[T]) Gt(n T) *IntField[T] {
	f.check("gt", f.value <= n, n)
	return f
}

// Gte 校验整数大于等于目标值
func (f *IntField[T]) Gte(n T) *IntField[T] {
	f.check("gte", f.value < n, n)
	return f
}

// Lt 校验整数小于目标值
func (f *IntField[T]) Lt(n T) *IntField[T] {
	f.check("lt", f.value >= n, n)
	return f
}

// Lte 校验整数小于等于目标值
func (f *IntField[T]) Lte(n T) *IntField[T] {
	f.check("lte", f.value > n, n)
	return f
}

// In 校验整数属于指定集合
func (f *IntField[T]) In(values ...T) *IntField[T] {
	if !slices.Contains(values, f.value) {
		f.check("in", true, arg("values", fmt.Sprint(values)))
	}
	return f
}

// NotIn 校验整数不属于指定集合
func (f *IntField[T]) NotIn(values ...T) *IntField[T] {
	if slices.Contains(values, f.value) {
		f.check("not_in", true, arg("values", fmt.Sprint(values)))
	}
	return f
}

// MultipleOf 校验整数为指定值的倍数
func (f *IntField[T]) MultipleOf(n T) *IntField[T] {
	f.check("multiple_of", n == 0 || f.value%n != 0, n)
	return f
}

// Digits 校验整数位数
func (f *IntField[T]) Digits(n int) *IntField[T] {
	if integerDigits(f.value) != n {
		f.check("digits", true, arg("digits", n))
	}
	return f
}

// DigitsBetween 校验整数位数范围
func (f *IntField[T]) DigitsBetween(min, max int) *IntField[T] {
	digits := integerDigits(f.value)
	if digits < min || digits > max {
		f.check("digits_between", true, arg("min", min), arg("max", max))
	}
	return f
}

// Positive 校验整数大于零
func (f *IntField[T]) Positive() *IntField[T] {
	var zero T
	f.check("positive", f.value <= zero)
	return f
}

// Negative 校验整数小于零
func (f *IntField[T]) Negative() *IntField[T] {
	var zero T
	f.check("negative", f.value >= zero)
	return f
}

// Same 校验整数与目标值相同
func (f *IntField[T]) Same(value T) *IntField[T] {
	f.check("same", f.value != value, value)
	return f
}

// Different 校验整数与目标值不同
func (f *IntField[T]) Different(value T) *IntField[T] {
	f.check("different", f.value == value, value)
	return f
}

// FloatField 浮点字段校验器
type FloatField[T FloatValue] struct {
	baseField[T]
}

// Float 创建浮点字段校验器
func Float[T FloatValue](translator *Translator, path string, value T) *FloatField[T] {
	return &FloatField[T]{baseField: newBaseField(translator, path, value)}
}

// Required 校验浮点数必填
func (f *FloatField[T]) Required() *FloatField[T] {
	var zero T
	f.check("required", f.value == zero)
	return f
}

// RequiredIf 条件成立时校验浮点数必填
func (f *FloatField[T]) RequiredIf(condition bool) *FloatField[T] {
	var zero T
	f.check("required_if", condition && f.value == zero)
	return f
}

// RequiredUnless 条件不成立时校验浮点数必填
func (f *FloatField[T]) RequiredUnless(condition bool) *FloatField[T] {
	var zero T
	f.check("required_unless", !condition && f.value == zero)
	return f
}

// Prohibited 校验浮点数等于零
func (f *FloatField[T]) Prohibited() *FloatField[T] {
	var zero T
	f.check("prohibited", f.value != zero)
	return f
}

// ProhibitedIf 条件成立时校验浮点数等于零
func (f *FloatField[T]) ProhibitedIf(condition bool) *FloatField[T] {
	var zero T
	f.check("prohibited_if", condition && f.value != zero)
	return f
}

// ProhibitedUnless 条件不成立时校验浮点数等于零
func (f *FloatField[T]) ProhibitedUnless(condition bool) *FloatField[T] {
	var zero T
	f.check("prohibited_unless", !condition && f.value != zero)
	return f
}

// Finite 校验浮点数为有限值
func (f *FloatField[T]) Finite() *FloatField[T] {
	f.check("finite", !isFinite(f.value))
	return f
}

// Min 校验浮点数最小值
func (f *FloatField[T]) Min(n T) *FloatField[T] {
	f.check("min", !isFinite(f.value) || !isFinite(n) || f.value < n, n)
	return f
}

// Max 校验浮点数最大值
func (f *FloatField[T]) Max(n T) *FloatField[T] {
	f.check("max", !isFinite(f.value) || !isFinite(n) || f.value > n, n)
	return f
}

// Between 校验浮点数范围
func (f *FloatField[T]) Between(min, max T) *FloatField[T] {
	invalid := !isFinite(f.value) || !isFinite(min) || !isFinite(max) || f.value < min || f.value > max
	if invalid {
		f.check("between", true, arg("min", min), arg("max", max))
	}
	return f
}

// Size 校验浮点数等于指定值
func (f *FloatField[T]) Size(n T) *FloatField[T] {
	f.check("size", !isFinite(f.value) || !isFinite(n) || f.value != n, n)
	return f
}

// Gt 校验浮点数大于目标值
func (f *FloatField[T]) Gt(n T) *FloatField[T] {
	f.check("gt", !isFinite(f.value) || !isFinite(n) || f.value <= n, n)
	return f
}

// Gte 校验浮点数大于等于目标值
func (f *FloatField[T]) Gte(n T) *FloatField[T] {
	f.check("gte", !isFinite(f.value) || !isFinite(n) || f.value < n, n)
	return f
}

// Lt 校验浮点数小于目标值
func (f *FloatField[T]) Lt(n T) *FloatField[T] {
	f.check("lt", !isFinite(f.value) || !isFinite(n) || f.value >= n, n)
	return f
}

// Lte 校验浮点数小于等于目标值
func (f *FloatField[T]) Lte(n T) *FloatField[T] {
	f.check("lte", !isFinite(f.value) || !isFinite(n) || f.value > n, n)
	return f
}

// In 校验浮点数属于指定集合
func (f *FloatField[T]) In(values ...T) *FloatField[T] {
	if !isFinite(f.value) || !slices.Contains(values, f.value) {
		f.check("in", true, arg("values", fmt.Sprint(values)))
	}
	return f
}

// NotIn 校验浮点数不属于指定集合
func (f *FloatField[T]) NotIn(values ...T) *FloatField[T] {
	if isFinite(f.value) && slices.Contains(values, f.value) {
		f.check("not_in", true, arg("values", fmt.Sprint(values)))
	}
	return f
}

// MultipleOf 校验浮点数为指定值的倍数
func (f *FloatField[T]) MultipleOf(n T) *FloatField[T] {
	f.check("multiple_of", !isFloatMultipleOf(f.value, n), n)
	return f
}

// Positive 校验浮点数大于零
func (f *FloatField[T]) Positive() *FloatField[T] {
	var zero T
	f.check("positive", !isFinite(f.value) || f.value <= zero)
	return f
}

// Negative 校验浮点数小于零
func (f *FloatField[T]) Negative() *FloatField[T] {
	var zero T
	f.check("negative", !isFinite(f.value) || f.value >= zero)
	return f
}

// Same 校验浮点数与目标值相同
func (f *FloatField[T]) Same(value T) *FloatField[T] {
	f.check("same", !isFinite(f.value) || !isFinite(value) || f.value != value, value)
	return f
}

// Different 校验浮点数与目标值不同
func (f *FloatField[T]) Different(value T) *FloatField[T] {
	f.check("different", isFinite(f.value) && isFinite(value) && f.value == value, value)
	return f
}

// Rule 添加浮点数自定义校验规则
func (f *FloatField[T]) Rule(rule, message string, valid func(T) bool, args ...any) *FloatField[T] {
	f.custom(rule, message, valid, args...)
	return f
}

// isFinite 判断浮点数是否为有限值
func isFinite[T FloatValue](value T) bool {
	result := float64(value)
	return !math.IsNaN(result) && !math.IsInf(result, 0)
}

// isFloatMultipleOf 按十进制表示判断浮点数是否为指定值的倍数
func isFloatMultipleOf[T FloatValue](value, divisor T) bool {
	if !isFinite(value) || !isFinite(divisor) {
		return false
	}
	bitSize := reflect.TypeOf(value).Bits()
	valueRat, valueOK := new(big.Rat).SetString(strconv.FormatFloat(float64(value), 'g', -1, bitSize))
	divisorRat, divisorOK := new(big.Rat).SetString(strconv.FormatFloat(float64(divisor), 'g', -1, bitSize))
	if !valueOK || !divisorOK || divisorRat.Sign() == 0 {
		return false
	}
	return new(big.Rat).Quo(valueRat, divisorRat).IsInt()
}

// integerDigits 获取整数的十进制位数
func integerDigits[T Integer](value T) int {
	digits := 0
	var zero T
	for {
		value /= 10
		digits++
		if value == zero {
			return digits
		}
	}
}
