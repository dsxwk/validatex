package tests

import (
	"testing"

	"github.com/dsxwk/validatex"
)

type wildcardOption struct {
	Label string
}

func (o wildcardOption) Validate(translator *validatex.Translator) validatex.Errors {
	return validatex.Check(
		validatex.Str(translator, "label", o.Label).Required().Errors(),
	)
}

type wildcardUser struct {
	Name    string
	Options map[string]wildcardOption
}

func (u wildcardUser) Validate(translator *validatex.Translator) validatex.Errors {
	return validatex.Check(
		validatex.Str(translator, "name", u.Name).Required().Errors(),
		validatex.MapStruct(translator, "options", u.Options),
	)
}

type wildcardRequest struct {
	Users  []wildcardUser
	Groups map[string]wildcardOption
}

func (r wildcardRequest) Validate(translator *validatex.Translator) validatex.Errors {
	return validatex.Check(
		validatex.SliceStruct(translator, "users", r.Users),
		validatex.MapStruct(translator, "groups", r.Groups),
	)
}

func TestWildcardPaths(t *testing.T) {
	translator := &validatex.Translator{
		Messages: validatex.Messages{
			"users.*.name.required":            "{field}不能为空",
			"users.*.options.*.label.required": "{field}不能为空",
			"groups.**.label.required":         "群组标签不能为空",
			"data.*.required":                  "{field}必填",
		},
		Translates: validatex.Translates{
			"users.*.name":            "用户名",
			"users.*.options.*.label": "选项标签",
			"data.*":                  "映射值",
		},
	}
	request := wildcardRequest{
		Users: []wildcardUser{{
			Options: map[string]wildcardOption{
				"theme": {},
			},
		}},
		Groups: map[string]wildcardOption{
			"admin": {},
		},
	}

	errs := request.Validate(translator)
	if len(errs) != 3 {
		t.Fatalf("expected 3 wildcard errors, got %v", errs)
	}
	if errs[0].Path != "users[0].name" || errs[0].Msg != "用户名不能为空" {
		t.Fatalf("unexpected wildcard struct error: %+v", errs[0])
	}
	if errs[1].Path != "users[0].options.theme.label" || errs[1].Msg != "选项标签不能为空" {
		t.Fatalf("unexpected wildcard map error: %+v", errs[1])
	}
	if errs[2].Path != "groups.admin.label" || errs[2].Msg != "群组标签不能为空" {
		t.Fatalf("unexpected deep wildcard error: %+v", errs[2])
	}

	data := map[string]string{"name": ""}
	errs = validatex.Map(translator, "data", data).Each(
		"required",
		"",
		func(_ string, value string) bool {
			return value != ""
		},
	).Errors()

	if len(errs) != 1 || errs[0].Path != "data.name" || errs[0].Msg != "映射值必填" {
		t.Fatalf("unexpected wildcard map each error: %v", errs)
	}
}
