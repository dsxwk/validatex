package tests

import (
	"math"
	"regexp"
	"strings"
	"testing"

	"github.com/dsxwk/validatex"
)

// TestStringRules 验证字符串规则会拒绝无效输入
func TestStringRules(t *testing.T) {
	translator := validatex.NewTranslator()
	cases := []struct {
		name string
		errs validatex.Errors
	}{
		{"required", validatex.Str(translator, "value", "   ").Required().Errors()},
		{"prohibited", validatex.Str(translator, "value", "abc").Prohibited().Errors()},
		{"prohibited_if", validatex.Str(translator, "value", "abc").ProhibitedIf(true).Errors()},
		{"max", validatex.Str(translator, "value", "abcd").Max(3).Errors()},
		{"between", validatex.Str(translator, "value", "ab").Between(3, 5).Errors()},
		{"size", validatex.Str(translator, "value", "ab").Size(3).Errors()},
		{"alpha", validatex.Str(translator, "value", "abc1").Alpha().Errors()},
		{"alpha_dash", validatex.Str(translator, "value", "abc 1").AlphaDash().Errors()},
		{"alpha_num", validatex.Str(translator, "value", "abc-1").AlphaNum().Errors()},
		{"ascii", validatex.Str(translator, "value", "中文").ASCII().Errors()},
		{"lowercase", validatex.Str(translator, "value", "Abc").Lowercase().Errors()},
		{"uppercase", validatex.Str(translator, "value", "Abc").Uppercase().Errors()},
		{"starts_with", validatex.Str(translator, "value", "abc").StartsWith("x", "y").Errors()},
		{"doesnt_start_with", validatex.Str(translator, "value", "abc").DoesntStartWith("a").Errors()},
		{"ends_with", validatex.Str(translator, "value", "abc").EndsWith("x", "y").Errors()},
		{"doesnt_end_with", validatex.Str(translator, "value", "abc").DoesntEndWith("c").Errors()},
		{"contains", validatex.Str(translator, "value", "abc").Contains("x", "y").Errors()},
		{"doesnt_contain", validatex.Str(translator, "value", "abc").DoesntContain("b").Errors()},
		{"in", validatex.Str(translator, "value", "c").In("a", "b").Errors()},
		{"not_in", validatex.Str(translator, "value", "a").NotIn("a", "b").Errors()},
		{"regex", validatex.Str(translator, "value", "abc").Regex(`^\d+$`).Errors()},
		{"not_regex", validatex.Str(translator, "value", "123").NotRegex(`^\d+$`).Errors()},
		{"numeric", validatex.Str(translator, "value", "abc").Numeric().Errors()},
		{"integer", validatex.Str(translator, "value", "1.2").Integer().Errors()},
		{"decimal", validatex.Str(translator, "value", "1.2").Decimal(2, 3).Errors()},
		{"digits", validatex.Str(translator, "value", "123").Digits(2).Errors()},
		{"digits_between", validatex.Str(translator, "value", "123").DigitsBetween(1, 2).Errors()},
		{"url", validatex.Str(translator, "value", "not-url").URL().Errors()},
		{"ip", validatex.Str(translator, "value", "999.1.1.1").IP().Errors()},
		{"ipv4", validatex.Str(translator, "value", "::1").IPv4().Errors()},
		{"ipv6", validatex.Str(translator, "value", "127.0.0.1").IPv6().Errors()},
		{"mac_address", validatex.Str(translator, "value", "invalid").MACAddress().Errors()},
		{"uuid", validatex.Str(translator, "value", "invalid").UUID().Errors()},
		{"ulid", validatex.Str(translator, "value", "invalid").ULID().Errors()},
		{"hex_color", validatex.Str(translator, "value", "#12").HexColor().Errors()},
		{"timezone", validatex.Str(translator, "value", "Invalid/Zone").Timezone().Errors()},
		{"json", validatex.Str(translator, "value", "{invalid").JSON().Errors()},
		{"accepted", validatex.Str(translator, "value", "no").Accepted().Errors()},
		{"accepted_if", validatex.Str(translator, "value", "no").AcceptedIf(true).Errors()},
		{"declined", validatex.Str(translator, "value", "yes").Declined().Errors()},
		{"declined_if", validatex.Str(translator, "value", "yes").DeclinedIf(true).Errors()},
		{"boolean", validatex.Str(translator, "value", "maybe").Boolean().Errors()},
		{"same", validatex.Str(translator, "value", "a").Same("b").Errors()},
		{"different", validatex.Str(translator, "value", "a").Different("a").Errors()},
		{"confirmed", validatex.Str(translator, "value", "a").Confirmed("b").Errors()},
		{"date", validatex.Str(translator, "value", "invalid").Date().Errors()},
		{"date_format", validatex.Str(translator, "value", "2026-10-23").DateFormat("2006/01/02").Errors()},
		{"date_equals", validatex.Str(translator, "value", "2026-10-23").DateEquals("2026-10-24").Errors()},
		{"before", validatex.Str(translator, "value", "2026-10-24").Before("2026-10-23").Errors()},
		{"before_or_equal", validatex.Str(translator, "value", "2026-10-24").BeforeOrEqual("2026-10-23").Errors()},
		{"after", validatex.Str(translator, "value", "2026-10-22").After("2026-10-23").Errors()},
		{"after_or_equal", validatex.Str(translator, "value", "2026-10-22").AfterOrEqual("2026-10-23").Errors()},
		{"filled", validatex.Str(translator, "value", " ").Filled().Errors()},
		{"required_if", validatex.Str(translator, "value", "").RequiredIf(true).Errors()},
		{"required_unless", validatex.Str(translator, "value", "").RequiredUnless(false).Errors()},
		{"prohibited_unless", validatex.Str(translator, "value", "set").ProhibitedUnless(false).Errors()},
		{"not_regexp", validatex.Str(translator, "value", "123").NotRegexp(nil).Errors()},
	}

	for _, item := range cases {
		t.Run(item.name, func(t *testing.T) {
			if len(item.errs) != 1 {
				t.Fatalf("expected one error, got %v", item.errs)
			}
		})
	}
}

// TestStringRulePasses 验证字符串规则会接受有效输入
func TestStringRulePasses(t *testing.T) {
	translator := validatex.NewTranslator()
	cases := []struct {
		name string
		errs validatex.Errors
	}{
		{"required", validatex.Str(translator, "value", "abc").Required().Errors()},
		{"max", validatex.Str(translator, "value", "中文").Max(2).Errors()},
		{"between", validatex.Str(translator, "value", "中文").Between(1, 2).Errors()},
		{"alpha", validatex.Str(translator, "value", "中文").Alpha().Errors()},
		{"alpha_dash", validatex.Str(translator, "value", "abc-1_2").AlphaDash().Errors()},
		{"url", validatex.Str(translator, "value", "https://example.com/path").URL().Errors()},
		{"ip", validatex.Str(translator, "value", "127.0.0.1").IP().Errors()},
		{"uuid", validatex.Str(translator, "value", "550e8400-e29b-41d4-a716-446655440000").UUID().Errors()},
		{"ulid", validatex.Str(translator, "value", "01ARZ3NDEKTSV4RRFFQ69G5FAV").ULID().Errors()},
		{"hex_color", validatex.Str(translator, "value", "#A1B2C3").HexColor().Errors()},
		{"timezone", validatex.Str(translator, "value", "Asia/Shanghai").Timezone().Errors()},
		{"json", validatex.Str(translator, "value", `{"name":"codex"}`).JSON().Errors()},
		{"decimal", validatex.Str(translator, "value", "12.345").Decimal(2, 3).Errors()},
		{"date", validatex.Str(translator, "value", "2026-10-23").Date().Errors()},
		{"after", validatex.Str(translator, "value", "2026-10-24").After("2026-10-23").Errors()},
	}

	for _, item := range cases {
		t.Run(item.name, func(t *testing.T) {
			if len(item.errs) != 0 {
				t.Fatalf("expected no errors, got %v", item.errs)
			}
		})
	}
}

// TestNumberRules 验证整数和浮点规则会拒绝无效输入
func TestNumberRules(t *testing.T) {
	translator := validatex.NewTranslator()
	invalid := []struct {
		name string
		errs validatex.Errors
	}{
		{"int_max", validatex.Int(translator, "value", 11).Max(10).Errors()},
		{"int_between", validatex.Int(translator, "value", 11).Between(1, 10).Errors()},
		{"int_size", validatex.Int(translator, "value", 11).Size(10).Errors()},
		{"int_gt", validatex.Int(translator, "value", 10).Gt(10).Errors()},
		{"int_gte", validatex.Int(translator, "value", 9).Gte(10).Errors()},
		{"int_lt", validatex.Int(translator, "value", 10).Lt(10).Errors()},
		{"int_lte", validatex.Int(translator, "value", 11).Lte(10).Errors()},
		{"int_in", validatex.Int(translator, "value", 3).In(1, 2).Errors()},
		{"int_not_in", validatex.Int(translator, "value", 1).NotIn(1, 2).Errors()},
		{"int_multiple", validatex.Int(translator, "value", 10).MultipleOf(3).Errors()},
		{"int_digits", validatex.Int(translator, "value", 123).Digits(2).Errors()},
		{"int_digits_between", validatex.Int(translator, "value", 123).DigitsBetween(1, 2).Errors()},
		{"int_positive", validatex.Int(translator, "value", 0).Positive().Errors()},
		{"int_negative", validatex.Int(translator, "value", 0).Negative().Errors()},
		{"int_required_if", validatex.Int(translator, "value", 0).RequiredIf(true).Errors()},
		{"int_required_unless", validatex.Int(translator, "value", 0).RequiredUnless(false).Errors()},
		{"int_prohibited", validatex.Int(translator, "value", 1).Prohibited().Errors()},
		{"int_prohibited_if", validatex.Int(translator, "value", 1).ProhibitedIf(true).Errors()},
		{"int_prohibited_unless", validatex.Int(translator, "value", 1).ProhibitedUnless(false).Errors()},
		{"int_same", validatex.Int(translator, "value", 1).Same(2).Errors()},
		{"int_different", validatex.Int(translator, "value", 1).Different(1).Errors()},
		{"float_min", validatex.Float(translator, "value", 1.5).Min(2).Errors()},
		{"float_max", validatex.Float(translator, "value", 3.5).Max(2).Errors()},
		{"float_between", validatex.Float(translator, "value", 3.5).Between(1, 2).Errors()},
		{"float_multiple", validatex.Float(translator, "value", 1.2).MultipleOf(0.5).Errors()},
		{"float_finite", validatex.Float(translator, "value", math.Inf(1)).Finite().Errors()},
		{"float_required_if", validatex.Float(translator, "value", 0.0).RequiredIf(true).Errors()},
		{"float_required_unless", validatex.Float(translator, "value", 0.0).RequiredUnless(false).Errors()},
		{"float_prohibited", validatex.Float(translator, "value", 1.0).Prohibited().Errors()},
		{"float_prohibited_if", validatex.Float(translator, "value", 1.0).ProhibitedIf(true).Errors()},
		{"float_prohibited_unless", validatex.Float(translator, "value", 1.0).ProhibitedUnless(false).Errors()},
		{"float_size", validatex.Float(translator, "value", 1.0).Size(2.0).Errors()},
		{"float_gt", validatex.Float(translator, "value", 1.0).Gt(2.0).Errors()},
		{"float_gte", validatex.Float(translator, "value", 1.0).Gte(2.0).Errors()},
		{"float_lt", validatex.Float(translator, "value", 2.0).Lt(1.0).Errors()},
		{"float_lte", validatex.Float(translator, "value", 2.0).Lte(1.0).Errors()},
		{"float_in", validatex.Float(translator, "value", 3.0).In(1.0, 2.0).Errors()},
		{"float_not_in", validatex.Float(translator, "value", 1.0).NotIn(1.0, 2.0).Errors()},
		{"float_positive", validatex.Float(translator, "value", 0.0).Positive().Errors()},
		{"float_negative", validatex.Float(translator, "value", 0.0).Negative().Errors()},
		{"float_same", validatex.Float(translator, "value", 1.0).Same(2.0).Errors()},
		{"float_different", validatex.Float(translator, "value", 1.0).Different(1.0).Errors()},
	}
	for _, item := range invalid {
		t.Run(item.name, func(t *testing.T) {
			if len(item.errs) != 1 {
				t.Fatalf("expected one error, got %v", item.errs)
			}
		})
	}

	if errs := validatex.Int(translator, "value", int32(10)).Between(1, 10).MultipleOf(2).Errors(); len(errs) != 0 {
		t.Fatalf("expected integer chain to pass, got %v", errs)
	}
	if errs := validatex.Float(translator, "value", 2.5).Between(1, 3).MultipleOf(0.5).Errors(); len(errs) != 0 {
		t.Fatalf("expected float chain to pass, got %v", errs)
	}
}

// TestCollectionRules 验证布尔、切片和Map规则
func TestCollectionRules(t *testing.T) {
	translator := validatex.NewTranslator()
	if errs := validatex.Slice(translator, "items", []int{1, 1}).Distinct().Errors(); len(errs) != 1 {
		t.Fatalf("expected distinct failure, got %v", errs)
	}
	if errs := validatex.Slice(translator, "items", []int{1}).Prohibited().Errors(); len(errs) != 1 {
		t.Fatalf("expected prohibited failure, got %v", errs)
	}
	if errs := validatex.Slice(translator, "items", []int{1, 2}).Contains(3).Errors(); len(errs) != 1 {
		t.Fatalf("expected contains failure, got %v", errs)
	}

	errs := validatex.Slice(translator, "items", []int{1, -1}).Each(
		"positive",
		"{field}必须大于零",
		func(value int) bool {
			return value > 0
		},
	).Errors()

	if len(errs) != 1 || errs[0].Path != "items[1]" {
		t.Fatalf("unexpected each errors: %v", errs)
	}

	data := map[string]any{"name": "codex", "extra": true}
	if errs := validatex.Map(translator, "data", data).HasKeys("name", "age").Errors(); len(errs) != 1 {
		t.Fatalf("expected has_keys failure, got %v", errs)
	}
	if errs := validatex.Map(translator, "data", data).AllowedKeys("name").Errors(); len(errs) != 1 {
		t.Fatalf("expected allowed_keys failure, got %v", errs)
	}
	if errs := validatex.Map(translator, "data", data).HasKeys("name").Min(2).Max(2).Errors(); len(errs) != 0 {
		t.Fatalf("expected map chain to pass, got %v", errs)
	}

	var emptySlice []int
	oneSlice := []int{1}
	assertOneErrorEach(t, map[string]validatex.Errors{
		"slice_required":          validatex.Slice(translator, "value", emptySlice).Required().Errors(),
		"slice_required_if":       validatex.Slice(translator, "value", emptySlice).RequiredIf(true).Errors(),
		"slice_required_unless":   validatex.Slice(translator, "value", emptySlice).RequiredUnless(false).Errors(),
		"slice_prohibited_if":     validatex.Slice(translator, "value", oneSlice).ProhibitedIf(true).Errors(),
		"slice_prohibited_unless": validatex.Slice(translator, "value", oneSlice).ProhibitedUnless(false).Errors(),
		"slice_min":               validatex.Slice(translator, "value", oneSlice).Min(2).Errors(),
		"slice_max":               validatex.Slice(translator, "value", oneSlice).Max(0).Errors(),
		"slice_between":           validatex.Slice(translator, "value", oneSlice).Between(2, 3).Errors(),
		"slice_size":              validatex.Slice(translator, "value", oneSlice).Size(2).Errors(),
		"slice_doesnt_contain":    validatex.Slice(translator, "value", oneSlice).DoesntContain(1).Errors(),
		"slice_rule": validatex.Slice(translator, "value", oneSlice).Rule(
			"custom",
			"",
			func([]int) bool { return false },
		).Errors(),
	})

	emptyMap := map[string]int{}
	oneMap := map[string]int{"value": 1}
	assertOneErrorEach(t, map[string]validatex.Errors{
		"map_required":          validatex.Map(translator, "value", emptyMap).Required().Errors(),
		"map_required_if":       validatex.Map(translator, "value", emptyMap).RequiredIf(true).Errors(),
		"map_required_unless":   validatex.Map(translator, "value", emptyMap).RequiredUnless(false).Errors(),
		"map_prohibited":        validatex.Map(translator, "value", oneMap).Prohibited().Errors(),
		"map_prohibited_if":     validatex.Map(translator, "value", oneMap).ProhibitedIf(true).Errors(),
		"map_prohibited_unless": validatex.Map(translator, "value", oneMap).ProhibitedUnless(false).Errors(),
		"map_between":           validatex.Map(translator, "value", oneMap).Between(2, 3).Errors(),
		"map_size":              validatex.Map(translator, "value", oneMap).Size(2).Errors(),
		"map_rule": validatex.Map(translator, "value", oneMap).Rule(
			"custom",
			"",
			func(map[string]int) bool { return false },
		).Errors(),
	})
}

// TestBooleanRules 验证布尔专用规则
func TestBooleanRules(t *testing.T) {
	translator := validatex.NewTranslator()
	if errs := validatex.Bool(translator, "enabled", false).Accepted().Errors(); len(errs) != 1 {
		t.Fatalf("expected accepted failure, got %v", errs)
	}
	if errs := validatex.Bool(translator, "enabled", true).Accepted().Errors(); len(errs) != 0 {
		t.Fatalf("expected accepted pass, got %v", errs)
	}
	if errs := validatex.Bool(translator, "enabled", true).Declined().Errors(); len(errs) != 1 {
		t.Fatalf("expected declined failure, got %v", errs)
	}
	if errs := validatex.Bool(translator, "enabled", false).AcceptedIf(true).Errors(); len(errs) != 1 {
		t.Fatalf("expected accepted_if failure, got %v", errs)
	}
	assertOneErrorEach(t, map[string]validatex.Errors{
		"declined_if":       validatex.Bool(translator, "value", true).DeclinedIf(true).Errors(),
		"prohibited":        validatex.Bool(translator, "value", true).Prohibited().Errors(),
		"prohibited_if":     validatex.Bool(translator, "value", true).ProhibitedIf(true).Errors(),
		"prohibited_unless": validatex.Bool(translator, "value", true).ProhibitedUnless(false).Errors(),
		"true":              validatex.Bool(translator, "value", false).True().Errors(),
		"false":             validatex.Bool(translator, "value", true).False().Errors(),
		"same":              validatex.Bool(translator, "value", false).Same(true).Errors(),
		"different":         validatex.Bool(translator, "value", true).Different(true).Errors(),
		"rule": validatex.Bool(translator, "value", false).Rule(
			"custom",
			"",
			func(bool) bool { return false },
		).Errors(),
	})
}

// TestNumericRejectsSpecialFloatSyntax 验证数字规则拒绝非十进制语法
func TestNumericRejectsSpecialFloatSyntax(t *testing.T) {
	translator := validatex.NewTranslator()
	for _, value := range []string{"NaN", "Inf", "+Inf", "-Inf", "0x1p2"} {
		if errs := validatex.Str(translator, "value", value).Numeric().Errors(); len(errs) != 1 {
			t.Fatalf("expected %q to be rejected as numeric, got %v", value, errs)
		}
	}
	for _, value := range []string{"1e3", ".5", "-2.5", "+3"} {
		if errs := validatex.Str(translator, "value", value).Numeric().Errors(); len(errs) != 0 {
			t.Fatalf("expected %q to pass numeric validation, got %v", value, errs)
		}
	}
}

// TestEmailRuleRejectsMalformedAddresses 验证常见非法邮箱地址
func TestEmailRuleRejectsMalformedAddresses(t *testing.T) {
	translator := validatex.NewTranslator()
	for _, value := range []string{
		"a..b@example.com",
		".a@example.com",
		"a.@example.com",
		"a@-example.com",
		"a@example-.com",
		"a@example..com",
		"a@example",
	} {
		if errs := validatex.Str(translator, "email", value).Email().Errors(); len(errs) != 1 {
			t.Fatalf("expected %q to be rejected as email, got %v", value, errs)
		}
	}
	if errs := validatex.Str(translator, "email", "codex@example.com").Email().Errors(); len(errs) != 0 {
		t.Fatalf("expected valid email to pass, got %v", errs)
	}
}

// TestULIDRuleRejectsOverflowPrefix 验证ULID范围和字符集检查
func TestULIDRuleRejectsOverflowPrefix(t *testing.T) {
	translator := validatex.NewTranslator()
	overflow := "8" + strings.Repeat("Z", 25)
	if errs := validatex.Str(translator, "id", overflow).ULID().Errors(); len(errs) != 1 {
		t.Fatalf("expected overflow ULID to fail, got %v", errs)
	}
	if errs := validatex.Str(translator, "id", "01ARZ3NDEKTSV4RRFFQ69G5FAV").ULID().Errors(); len(errs) != 0 {
		t.Fatalf("expected valid ULID to pass, got %v", errs)
	}
}

// TestDecimalRejectsRepeatedSign 验证十进制符号语法检查
func TestDecimalRejectsRepeatedSign(t *testing.T) {
	translator := validatex.NewTranslator()
	if errs := validatex.Str(translator, "amount", "+-1").Decimal(0, 2).Errors(); len(errs) != 1 {
		t.Fatalf("expected repeated sign to fail, got %v", errs)
	}
	if errs := validatex.Str(translator, "amount", "-1.25").Decimal(0, 2).Errors(); len(errs) != 0 {
		t.Fatalf("expected signed decimal to pass, got %v", errs)
	}
}

// TestFloatMultipleOfHandlesLargeValues 验证大数倍数判断的十进制精度
func TestFloatMultipleOfHandlesLargeValues(t *testing.T) {
	translator := validatex.NewTranslator()
	if errs := validatex.Float(translator, "amount", 1e12).MultipleOf(0.1).Errors(); len(errs) != 0 {
		t.Fatalf("expected large decimal multiple to pass, got %v", errs)
	}
	if errs := validatex.Float(translator, "amount", 1.2).MultipleOf(0.5).Errors(); len(errs) != 1 {
		t.Fatalf("expected non-multiple to fail, got %v", errs)
	}
	if errs := validatex.Float(translator, "amount", 1e12).MultipleOf(0.3).Errors(); len(errs) != 1 {
		t.Fatalf("expected large non-multiple to fail, got %v", errs)
	}
	if errs := validatex.Float(translator, "amount", 1e16).MultipleOf(0.7).Errors(); len(errs) != 1 {
		t.Fatalf("expected large non-multiple to fail, got %v", errs)
	}
}

// TestNumericRequiredAndRegexpEmptyValues 验证零值和空可选值行为
func TestNumericRequiredAndRegexpEmptyValues(t *testing.T) {
	translator := validatex.NewTranslator()
	if errs := validatex.Int(translator, "count", 0).Required().Errors(); len(errs) != 1 {
		t.Fatalf("expected integer required failure, got %v", errs)
	}
	if errs := validatex.Float(translator, "amount", 0.0).Required().Errors(); len(errs) != 1 {
		t.Fatalf("expected float required failure, got %v", errs)
	}
	if errs := validatex.Str(translator, "value", "").Regexp(nil).Errors(); len(errs) != 0 {
		t.Fatalf("expected empty regexp value to pass, got %v", errs)
	}
	if errs := validatex.Str(translator, "value", "").NotRegexp(nil).Errors(); len(errs) != 0 {
		t.Fatalf("expected empty not-regexp value to pass, got %v", errs)
	}
}

// TestMapEachErrorsAreDeterministic 验证Map遍历错误顺序稳定
func TestMapEachErrorsAreDeterministic(t *testing.T) {
	translator := validatex.NewTranslator()
	values := map[string]string{
		"beta":  "",
		"alpha": "",
	}

	for range 20 {
		errs := validatex.Map(translator, "data", values).Each(
			"required",
			"",
			func(_ string, value string) bool {
				return value != ""
			},
		).Errors()
		if len(errs) != 2 || errs[0].Path != "data.alpha" || errs[1].Path != "data.beta" {
			t.Fatalf("unexpected map each error order: %v", errs)
		}
	}
}

// TestConditionalValidation 验证条件辅助函数和预编译正则
func TestConditionalValidation(t *testing.T) {
	translator := validatex.NewTranslator()
	if errs := validatex.When(false, func() validatex.Errors {
		return validatex.Str(translator, "name", "").Required().Errors()
	}); len(errs) != 0 {
		t.Fatalf("expected false condition to skip validation, got %v", errs)
	}
	if errs := validatex.When(true, func() validatex.Errors {
		return validatex.Str(translator, "name", "").Required().Errors()
	}); len(errs) != 1 {
		t.Fatalf("expected true condition to validate, got %v", errs)
	}

	var optional *string
	if errs := validatex.Nullable(optional, func(value string) validatex.Errors {
		return validatex.Str(translator, "name", value).Required().Errors()
	}); len(errs) != 0 {
		t.Fatalf("expected nil value to skip validation, got %v", errs)
	}
	value := ""
	if errs := validatex.Nullable(&value, func(value string) validatex.Errors {
		return validatex.Str(translator, "name", value).Required().Errors()
	}); len(errs) != 1 {
		t.Fatalf("expected non-nil value to validate, got %v", errs)
	}

	pattern := regexp.MustCompile(`^\d+$`)
	if errs := validatex.Str(translator, "value", "123").Regexp(pattern).Errors(); len(errs) != 0 {
		t.Fatalf("expected regexp to pass, got %v", errs)
	}
	if errs := validatex.Str(translator, "value", "abc").Regexp(pattern).Errors(); len(errs) != 1 {
		t.Fatalf("expected regexp to fail, got %v", errs)
	}
}

// TestNilRulePredicatesFailClosed 验证nil回调不会panic且按失败处理
func TestNilRulePredicatesFailClosed(t *testing.T) {
	translator := validatex.NewTranslator()
	if errs := validatex.Str(translator, "value", "set").Rule("custom", "", nil).Errors(); len(errs) != 1 {
		t.Fatalf("expected nil string predicate to fail, got %v", errs)
	}
	if errs := validatex.Slice(translator, "value", []int{1}).Each("custom", "", nil).Errors(); len(errs) != 1 {
		t.Fatalf("expected nil slice predicate to fail, got %v", errs)
	}
	if errs := validatex.Map(translator, "value", map[string]int{"one": 1}).Each("custom", "", nil).Errors(); len(errs) != 1 {
		t.Fatalf("expected nil map predicate to fail, got %v", errs)
	}
}

// assertOneErrorEach 验证每个命名用例都只产生一个错误
func assertOneErrorEach(t *testing.T, tests map[string]validatex.Errors) {
	t.Helper()
	for name, errs := range tests {
		t.Run(name, func(t *testing.T) {
			if len(errs) != 1 {
				t.Fatalf("expected one error, got %v", errs)
			}
		})
	}
}
