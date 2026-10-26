package tests

import (
	"testing"

	"github.com/dsxwk/validatex"
)

var allocationErrors validatex.Errors

type allocationMapValue string

// Validate 为分配测试实现validatex.Validatable
func (value allocationMapValue) Validate(translator *validatex.Translator) validatex.Errors {
	return validatex.Str(translator, "", string(value)).Required().Errors()
}

// TestValidRulesDoNotAllocate 验证有效输入不会产生分配
func TestValidRulesDoNotAllocate(t *testing.T) {
	translator := validatex.NewTranslator()
	values := map[string]string{"name": "codex"}
	arrayValues := map[[3]int]string{{1, 23, 4}: "codex"}
	arrayStructValues := map[[3]int]allocationMapValue{{1, 23, 4}: "codex"}
	sliceStructValues := []allocationMapValue{"one", "two"}

	tests := []struct {
		name   string
		errors func() validatex.Errors
	}{
		{
			name: "string bounds",
			errors: func() validatex.Errors {
				return validatex.Str(translator, "name", "codex").Min(3).Max(10).Between(3, 10).Errors()
			},
		},
		{
			name: "named string argument",
			errors: func() validatex.Errors {
				return validatex.Str(translator, "name", "codex").StartsWith("code", "id-").Errors()
			},
		},
		{
			name: "integer bounds",
			errors: func() validatex.Errors {
				return validatex.Int(translator, "age", 18).Min(18).Max(60).Between(18, 60).Errors()
			},
		},
		{
			name: "map keys",
			errors: func() validatex.Errors {
				return validatex.Map(translator, "values", values).HasKeys("name").AllowedKeys("name").Errors()
			},
		},
		{
			name: "comparable map keys",
			errors: func() validatex.Errors {
				return validatex.Map(translator, "values", arrayValues).
					HasKeys([3]int{1, 23, 4}).
					AllowedKeys([3]int{1, 23, 4}).
					Each("required", "", func(_ [3]int, value string) bool {
						return value != ""
					}).
					Errors()
			},
		},
		{
			name: "valid map struct keys",
			errors: func() validatex.Errors {
				return validatex.MapStruct(translator, "values", arrayStructValues)
			},
		},
		{
			name: "valid slice structs",
			errors: func() validatex.Errors {
				return validatex.SliceStruct(translator, "values", sliceStructValues)
			},
		},
		{
			name: "empty optional fields",
			errors: func() validatex.Errors {
				return validatex.Str(translator, "optional", "").
					Regex(`^\d+$`).
					URL().
					JSON().
					Timezone().
					UUID().
					Date().
					DateFormat("2006-01-02").
					Errors()
			},
		},
		{
			name: "case insensitive values",
			errors: func() validatex.Errors {
				return validatex.Str(translator, "accepted", " TRUE ").Accepted().Errors()
			},
		},
		{
			name: "lowercase",
			errors: func() validatex.Errors {
				return validatex.Str(translator, "value", "codex-123").Lowercase().Errors()
			},
		},
		{
			name: "uppercase",
			errors: func() validatex.Errors {
				return validatex.Str(translator, "value", "CODEX-123").Uppercase().Errors()
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			allocs := testing.AllocsPerRun(1000, func() {
				allocationErrors = test.errors()
			})
			if allocs != 0 {
				t.Fatalf("expected no allocations, got %.2f", allocs)
			}
			if len(allocationErrors) != 0 {
				t.Fatalf("expected valid input, got %v", allocationErrors)
			}
		})
	}
}

// TestMapKeyFormatting 验证可比较Map键生成稳定路径
func TestMapKeyFormatting(t *testing.T) {
	tests := []struct {
		name string
		got  string
		want string
	}{
		{name: "string", got: mapKeyPath("name"), want: "data.name"},
		{name: "bool", got: mapKeyPath(true), want: "data.true"},
		{name: "int", got: mapKeyPath(int64(-23)), want: "data.-23"},
		{name: "float", got: mapKeyPath(1.5), want: "data.1%2E5"},
		{name: "complex", got: mapKeyPath(complex(1, 2)), want: "data.(1+2i)"},
		{name: "array", got: mapKeyPath([3]int{1, 23, 4}), want: "data.%5B3%5Dint{1, 23, 4}"},
		{name: "interface", got: mapKeyPath[any]("name"), want: `data.string:"name"`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.got != test.want {
				t.Fatalf("map key path = %q, want %q", test.got, test.want)
			}
		})
	}
}

// mapKeyPath 返回单个Map键生成的校验路径
func mapKeyPath[K comparable](key K) string {
	values := map[K]string{key: ""}
	errors := validatex.Map(nil, "data", values).Each(
		"required",
		"",
		func(_ K, value string) bool {
			return value != ""
		},
	).Errors()
	if len(errors) != 1 {
		return ""
	}
	return errors[0].Path
}

// TestMessageFormatting 验证占位符展开和回退行为
func TestMessageFormatting(t *testing.T) {
	rangeTranslator := &validatex.Translator{
		Messages: validatex.Messages{
			"between": "{field}|{min}|{max}|{0}|{1}|{missing}",
		},
	}
	if errors := validatex.Int(rangeTranslator, "name", 5).Between(1, 2).Errors(); len(errors) != 1 ||
		errors[0].Msg != "name|1|2|1|2|{missing}" {
		t.Fatalf("unexpected named message formatting: %v", errors)
	}

	if errors := validatex.Int(nil, "value", 5).Rule(
		"range",
		"{0}/{1}/{max}",
		func(int) bool { return false },
		10,
		20,
	).Errors(); len(errors) != 1 || errors[0].Msg != "10/20/20" {
		t.Fatalf("unexpected positional message formatting: %v", errors)
	}

	if errors := validatex.Int(nil, "value", 5).Rule(
		"named",
		"{field}:{value}",
		func(int) bool { return false },
		map[string]any{"value": 3},
	).Errors(); len(errors) != 1 || errors[0].Msg != "value:3" {
		t.Fatalf("unexpected map message formatting: %v", errors)
	}

	translator := &validatex.Translator{
		Messages:   validatex.Messages{"required": "{field}"},
		Translates: validatex.Translates{"value": "{path}"},
	}
	if errors := validatex.Str(translator, "value", "").Required().Errors(); len(errors) != 1 ||
		errors[0].Msg != "{path}" {
		t.Fatalf("inserted placeholders should not expand again: %v", errors)
	}
}

// FuzzNumericString 模糊测试公开的数字字符串校验
func FuzzNumericString(f *testing.F) {
	for _, value := range []string{"", "1", "-1.5", "1e3", "NaN", "0x1p2"} {
		f.Add(value)
	}
	translator := validatex.NewTranslator()
	f.Fuzz(func(t *testing.T, value string) {
		_ = validatex.Str(translator, "value", value).Numeric().Errors()
	})
}

// FuzzEmail 模糊测试公开的邮箱校验
func FuzzEmail(f *testing.F) {
	for _, value := range []string{"", "a@example.com", "a..b@example.com", "@example.com"} {
		f.Add(value)
	}
	translator := validatex.NewTranslator()
	f.Fuzz(func(t *testing.T, value string) {
		_ = validatex.Str(translator, "value", value).Email().Errors()
	})
}

// FuzzWildcardPath 使用任意路径模糊测试通配符翻译
func FuzzWildcardPath(f *testing.F) {
	f.Add("users.**.email", "users[0].profile.email")
	f.Add("**.*.**", "a.b.c")
	f.Fuzz(func(t *testing.T, pattern, actual string) {
		translator := &validatex.Translator{
			Messages: validatex.Messages{pattern + ".required": "invalid"},
		}
		_ = validatex.Str(translator, actual, "").Required().Errors()
	})
}
