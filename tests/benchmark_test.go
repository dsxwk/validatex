package tests

import (
	"strconv"
	"testing"

	"github.com/dsxwk/validatex"
)

type validatexBenchRequest struct {
	Username string
	Email    string
}

type validatexBenchItem struct {
	Value string
}

// Validate 为嵌套基准实现validatex.Validatable
func (item validatexBenchItem) Validate(fields validatex.Fields) validatex.Errors {
	return fields.Str("value", item.Value).Required().Errors()
}

type validatexBenchNestedRequest struct {
	Items  []validatexBenchItem
	Groups map[string]validatexBenchItem
}

// Validate 为嵌套基准实现 validatex.Validatable。
func (request *validatexBenchNestedRequest) Validate(fields validatex.Fields) validatex.Errors {
	return validatex.Check(
		fields.SliceStruct("items", request.Items),
		fields.MapStruct("groups", request.Groups),
	)
}

// Validate 为请求基准实现validatex.Validatable
func (r *validatexBenchRequest) Validate(fields validatex.Fields) validatex.Errors {
	return validatex.Check(
		fields.Str("username", r.Username).Required().Min(3).Errors(),
		fields.Str("email", r.Email).Required().Email().Errors(),
	)
}

// BenchmarkValid 测试有效请求的校验性能
func BenchmarkValid(b *testing.B) {
	translator := validatex.NewTranslator()
	validValidatex := &validatexBenchRequest{
		Username: "codex",
		Email:    "codex@example.com",
	}
	b.Run("validatex", func(b *testing.B) {
		b.ReportAllocs()
		for index := 0; index < b.N; index++ {
			validatex.ValidateRequestWithTranslator(validValidatex, translator)
		}
	})
	b.Run("validatex-request", func(b *testing.B) {
		b.ReportAllocs()
		for index := 0; index < b.N; index++ {
			validatex.ValidateRequest(validValidatex)
		}
	})
}

// BenchmarkInvalid 测试产生字段错误时的校验性能
func BenchmarkInvalid(b *testing.B) {
	translator := validatex.NewTranslator()
	invalidValidatex := &validatexBenchRequest{
		Username: "x",
		Email:    "invalid",
	}
	b.Run("validatex", func(b *testing.B) {
		b.ReportAllocs()
		for index := 0; index < b.N; index++ {
			if len(validatex.ValidateRequestWithTranslator(invalidValidatex, translator)) == 0 {
				b.Fatal("validation passed")
			}
		}
	})
	b.Run("validatex-request", func(b *testing.B) {
		b.ReportAllocs()
		for index := 0; index < b.N; index++ {
			if len(validatex.ValidateRequest(invalidValidatex)) == 0 {
				b.Fatal("validation passed")
			}
		}
	})
}

// BenchmarkParameterizedRules 测试参数化规则的有效和无效路径
func BenchmarkParameterizedRules(b *testing.B) {
	translator := validatex.NewTranslator()
	b.Run("valid", func(b *testing.B) {
		b.ReportAllocs()
		for index := 0; index < b.N; index++ {
			validatex.Str(translator, "name", "codex").
				Min(3).
				Max(10).
				Between(3, 10).
				Errors()
		}
	})
	b.Run("invalid", func(b *testing.B) {
		b.ReportAllocs()
		for index := 0; index < b.N; index++ {
			if len(validatex.Int(translator, "age", 16).Between(18, 60).Errors()) == 0 {
				b.Fatal("validation passed")
			}
		}
	})
}

// BenchmarkWildcardTranslation 对比动态和预编译的通配符查询
func BenchmarkWildcardTranslation(b *testing.B) {
	messages := validatex.Messages{
		"groups.*.field42.required": "{field} required",
	}
	translates := validatex.Translates{
		"groups.*.field42": "Group field",
	}
	for index := range 128 {
		field := "field" + strconv.Itoa(index)
		messages["groups.*."+field+".required"] = "{field} required"
		translates["groups.*."+field] = "Group field"
	}
	translator := &validatex.Translator{
		Messages:   messages,
		Translates: translates,
	}
	compiled := translator.Compile()

	for _, test := range []struct {
		name       string
		translator *validatex.Translator
	}{
		{name: "dynamic", translator: translator},
		{name: "compiled", translator: compiled},
	} {
		b.Run(test.name, func(b *testing.B) {
			b.ReportAllocs()
			for index := 0; index < b.N; index++ {
				if len(validatex.Str(test.translator, "groups.admin.field42", "").Required().Errors()) == 0 {
					b.Fatal("validation passed")
				}
			}
		})
	}
}

// BenchmarkNestedErrors 测试大量嵌套错误的校验性能
func BenchmarkNestedErrors(b *testing.B) {
	items := make([]validatexBenchItem, 64)
	groups := make(map[string]validatexBenchItem, 64)
	for index := range 64 {
		groups[strconv.Itoa(index)] = validatexBenchItem{}
	}
	request := &validatexBenchNestedRequest{
		Items:  items,
		Groups: groups,
	}
	translator := validatex.NewTranslator()

	b.ReportAllocs()
	for index := 0; index < b.N; index++ {
		if len(validatex.ValidateRequestWithTranslator(request, translator)) != 128 {
			b.Fatal("unexpected error count")
		}
	}
}
