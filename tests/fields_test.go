package tests

import (
	"testing"

	"github.com/dsxwk/validatex"
)

// TestFieldsBindTranslator 验证字段构造器绑定翻译器
func TestFieldsBindTranslator(t *testing.T) {
	translator := &validatex.Translator{
		Messages: validatex.Messages{
			"required": "{field}必填",
			"gt":       "{field}必须大于{0}",
		},
		Translates: validatex.Translates{
			"Name":  "名称",
			"Age":   "年龄",
			"IDs":   "ID列表",
			"IDs.*": "ID",
		},
	}
	fields := validatex.NewFields(translator)
	errs := validatex.Check(
		fields.Str("Name", "").Required().Errors(),
		fields.Int("Age", 0).Gt(0).Errors(),
		fields.Slice("IDs", []int64{0}).Each("gt", "", func(value int64) bool {
			return value > 0
		}, int64(0)).Errors(),
	)
	if len(errs) != 3 {
		t.Fatalf("unexpected field errors: %v", errs)
	}
	if errs[0].Msg != "名称必填" || errs[1].Msg != "年龄必须大于0" || errs[2].Path != "IDs[0]" {
		t.Fatalf("unexpected translated fields: %v", errs)
	}
}

// TestFieldsNestedPointerOptions 验证嵌套指针字段选项
func TestFieldsNestedPointerOptions(t *testing.T) {
	var address *pointerAddress
	fields := validatex.NewFields(validatex.NewTranslator())

	errs := fields.StructPointer("Address", address, validatex.WithRequired())
	if len(errs) != 1 || errs[0].Path != "Address" || errs[0].Rule != "required" {
		t.Fatalf("unexpected required pointer field errors: %v", errs)
	}

	errs = fields.SliceStructPointer(
		"Items",
		[]*pointerAddress{nil},
		validatex.WithRequiredElements(),
	)
	if len(errs) != 1 || errs[0].Path != "Items[0]" || errs[0].Rule != "required" {
		t.Fatalf("unexpected required pointer element errors: %v", errs)
	}

	errs = fields.MapStructPointer(
		"Groups",
		map[string]*pointerAddress{"admin": nil},
		validatex.WithRequiredElements(),
	)
	if len(errs) != 1 || errs[0].Path != "Groups.admin" || errs[0].Rule != "required" {
		t.Fatalf("unexpected required pointer map errors: %v", errs)
	}
}
