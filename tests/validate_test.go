package tests

import (
	"testing"

	"github.com/dsxwk/validatex"
)

type testAddress struct {
	Province string
}

func (a testAddress) Validate(translator *validatex.Translator) validatex.Errors {
	return validatex.Check(
		validatex.Str(translator, "province", a.Province).Required().Errors(),
	)
}

type testUser struct {
	Username string
	Email    string
	Age      int64
	Address  testAddress
	Items    []testAddress
}

func (u testUser) Validate(translator *validatex.Translator) validatex.Errors {
	return validatex.Check(
		validatex.Str(translator, "username", u.Username).Required().Min(3).Errors(),
		validatex.Str(translator, "email", u.Email).Required().Email().Errors(),
		validatex.Int(translator, "age", u.Age).Min(18).Errors(),
		validatex.Struct(translator, "address", u.Address),
		validatex.SliceStruct(translator, "items", u.Items),
	)
}

func TestCheckAndTranslate(t *testing.T) {
	translator := &validatex.Translator{
		Messages: validatex.Messages{
			"required": "{field}必填",
			"min":      "{field}不能少于{min}",
			"email":    "{field}格式不正确",
		},
		Translates: validatex.Translates{
			"username": "用户名",
			"email":    "邮箱",
			"age":      "年龄",
		},
	}

	errs := validatex.Check(
		validatex.Str(translator, "username", "").Required().Min(3).Errors(),
		validatex.Str(translator, "email", "invalid").Required().Email().Errors(),
		validatex.Int(translator, "age", 16).Min(18).Errors(),
	)
	if len(errs) != 3 {
		t.Fatalf("expected 3 errors, got %d", len(errs))
	}
	if errs[0].Msg != "用户名必填" {
		t.Fatalf("unexpected required message: %s", errs[0].Msg)
	}
	if errs[1].Msg != "邮箱格式不正确" {
		t.Fatalf("unexpected email message: %s", errs[1].Msg)
	}
	if errs[2].Msg != "年龄不能少于18" {
		t.Fatalf("unexpected min message: %s", errs[2].Msg)
	}
}

func TestStringMinUsesRuneCount(t *testing.T) {
	translator := &validatex.Translator{Messages: validatex.Messages{"min": "{field}不能少于{min}"}}
	if errs := validatex.Str(translator, "name", "中文").Min(2).Errors(); len(errs) != 0 {
		t.Fatalf("expected two Chinese characters to pass, got %v", errs)
	}
	if errs := validatex.Str(translator, "name", "中").Min(2).Errors(); len(errs) != 1 {
		t.Fatalf("expected one Chinese character to fail, got %v", errs)
	}
}

func TestIntSupportsIntegerTypes(t *testing.T) {
	translator := &validatex.Translator{Messages: validatex.Messages{"min": "{field}不能少于{min}"}}
	if errs := validatex.Int(translator, "age", int32(16)).Min(18).Errors(); len(errs) != 1 {
		t.Fatalf("expected int32 validation to fail, got %v", errs)
	}
	if errs := validatex.Int(translator, "age", uint64(18)).Min(18).Errors(); len(errs) != 0 {
		t.Fatalf("expected uint64 validation to pass, got %v", errs)
	}
}

func TestCustomRules(t *testing.T) {
	translator := &validatex.Translator{
		Translates: validatex.Translates{
			"code":  "编码",
			"count": "数量",
		},
	}
	prefixRule := func(value string) bool {
		return len(value) > 0 && value[0] == 'A'
	}
	if errs := validatex.Str(translator, "code", "A123").Rule("prefix", "{field}必须以{prefix}开头", prefixRule, map[string]any{"prefix": "A"}).Errors(); len(errs) != 0 {
		t.Fatalf("expected custom string rule to pass, got %v", errs)
	}

	errs := validatex.Check(
		validatex.Str(translator, "code", "123").Rule("prefix", "{field}必须以{prefix}开头", prefixRule, map[string]any{"prefix": "A"}).Errors(),
		validatex.Int(translator, "count", int32(3)).Rule("even", "{field}必须是偶数", func(value int32) bool {
			return value%2 == 0
		}).Errors(),
	)
	if len(errs) != 2 {
		t.Fatalf("expected 2 custom rule errors, got %d", len(errs))
	}
	if errs[0].Rule != "prefix" || errs[0].Msg != "编码必须以A开头" {
		t.Fatalf("unexpected custom string error: %+v", errs[0])
	}
	if errs[1].Rule != "even" || errs[1].Msg != "数量必须是偶数" {
		t.Fatalf("unexpected custom integer error: %+v", errs[1])
	}
}

func TestNestedPaths(t *testing.T) {
	translator := &validatex.Translator{
		Messages:   validatex.Messages{"required": "{field}必填"},
		Translates: validatex.Translates{"province": "省份"},
	}
	user := testUser{
		Username: "codex",
		Email:    "codex@example.com",
		Age:      18,
		Items:    []testAddress{{Province: "浙江"}, {}},
	}

	errs := user.Validate(translator)
	if len(errs) != 2 {
		t.Fatalf("expected 2 errors, got %d", len(errs))
	}
	if errs[0].Path != "address.province" {
		t.Fatalf("unexpected struct path: %s", errs[0].Path)
	}
	if errs[1].Path != "items[1].province" {
		t.Fatalf("unexpected slice path: %s", errs[1].Path)
	}
	if errs[1].Msg != "省份必填" {
		t.Fatalf("unexpected nested message: %s", errs[1].Msg)
	}
}
