package tests

import (
	"math"
	"testing"

	"github.com/dsxwk/validatex"
)

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
	}

	for _, item := range cases {
		t.Run(item.name, func(t *testing.T) {
			if len(item.errs) != 1 {
				t.Fatalf("expected one error, got %v", item.errs)
			}
		})
	}
}

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
		{"int_prohibited", validatex.Int(translator, "value", 1).Prohibited().Errors()},
		{"float_min", validatex.Float(translator, "value", 1.5).Min(2).Errors()},
		{"float_max", validatex.Float(translator, "value", 3.5).Max(2).Errors()},
		{"float_between", validatex.Float(translator, "value", 3.5).Between(1, 2).Errors()},
		{"float_multiple", validatex.Float(translator, "value", 1.2).MultipleOf(0.5).Errors()},
		{"float_finite", validatex.Float(translator, "value", math.Inf(1)).Finite().Errors()},
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
}

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
}
