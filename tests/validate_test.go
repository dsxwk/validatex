package tests

import (
	"slices"
	"testing"

	"github.com/dsxwk/validatex"
)

type testAddress struct {
	Province string
}

// Validate 为嵌套测试实现validatex.Validatable
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

type pointerAddress struct{}

// Validate 为指针集成测试返回固定错误
func (*pointerAddress) Validate(*validatex.Translator) validatex.Errors {
	return validatex.Errors{{Path: "unexpected"}}
}

type customNested struct{}

// Validate 返回固定自定义消息
func (customNested) Validate(*validatex.Translator) validatex.Errors {
	return validatex.Errors{{Path: "code", Rule: "custom", Msg: "自定义消息"}}
}

type translatedNested struct{}

// Validate 为翻译测试返回仅包含规则名的错误
func (translatedNested) Validate(*validatex.Translator) validatex.Errors {
	return validatex.Errors{{Path: "code", Rule: "required"}}
}

type collidingMapKey string

// String 返回相同文本以测试路径唯一性
func (collidingMapKey) String() string {
	return "same"
}

type mapKeyValue string

// Validate 使用空路径校验当前值
func (value mapKeyValue) Validate(translator *validatex.Translator) validatex.Errors {
	return validatex.Str(translator, "", string(value)).Required().Errors()
}

type namedIntKey int

// Validate 为嵌套测试请求实现validatex.Validatable
func (u testUser) Validate(translator *validatex.Translator) validatex.Errors {
	return validatex.Check(
		validatex.Str(translator, "username", u.Username).Required().Min(3).Errors(),
		validatex.Str(translator, "email", u.Email).Required().Email().Errors(),
		validatex.Int(translator, "age", u.Age).Min(18).Errors(),
		validatex.Struct(translator, "address", u.Address),
		validatex.SliceStruct(translator, "items", u.Items),
	)
}

// TestCheckAndTranslate 验证错误聚合、翻译和显示名回退
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

	fallbackTranslator := &validatex.Translator{
		Translates: validatex.Translates{"user.email": "邮箱"},
	}
	errs = validatex.Str(fallbackTranslator, "user.email", "invalid").Email().Errors()
	if len(errs) != 1 || errs[0].Msg != "邮箱 is email" {
		t.Fatalf("unexpected fallback message: %v", errs)
	}
}

// TestStringMinUsesRuneCount 验证字符串长度按字符计数
func TestStringMinUsesRuneCount(t *testing.T) {
	translator := &validatex.Translator{Messages: validatex.Messages{"min": "{field}不能少于{min}"}}
	if errs := validatex.Str(translator, "name", "中文").Min(2).Errors(); len(errs) != 0 {
		t.Fatalf("expected two Chinese characters to pass, got %v", errs)
	}
	if errs := validatex.Str(translator, "name", "中").Min(2).Errors(); len(errs) != 1 {
		t.Fatalf("expected one Chinese character to fail, got %v", errs)
	}
}

// TestIntSupportsIntegerTypes 验证泛型整数构造器
func TestIntSupportsIntegerTypes(t *testing.T) {
	translator := &validatex.Translator{Messages: validatex.Messages{"min": "{field}不能少于{min}"}}
	if errs := validatex.Int(translator, "age", int32(16)).Min(18).Errors(); len(errs) != 1 {
		t.Fatalf("expected int32 validation to fail, got %v", errs)
	}
	if errs := validatex.Int(translator, "age", uint64(18)).Min(18).Errors(); len(errs) != 0 {
		t.Fatalf("expected uint64 validation to pass, got %v", errs)
	}
}

// TestCustomRules 验证命名和位置自定义参数
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

// TestNestedPaths 验证嵌套结构体和切片错误路径
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
	if indexes := errs[1].Indexes(); !slices.Equal(indexes, []int{1}) {
		t.Fatalf("unexpected slice path indexes: %v", indexes)
	}
	if errs[1].Msg != "省份必填" {
		t.Fatalf("unexpected nested message: %s", errs[1].Msg)
	}
}

// TestNestedValidationPreservesCustomMessage 验证固定消息在路径重定位后保持不变
func TestNestedValidationPreservesCustomMessage(t *testing.T) {
	translator := &validatex.Translator{
		Messages: validatex.Messages{"custom": "翻译消息"},
	}
	errs := validatex.Struct(translator, "parent", customNested{})
	if len(errs) != 1 || errs[0].Path != "parent.code" || errs[0].Msg != "自定义消息" {
		t.Fatalf("unexpected custom nested error: %+v", errs)
	}
}

// TestNestedValidationTranslatesManualRuleError 验证手写规则错误会被翻译
func TestNestedValidationTranslatesManualRuleError(t *testing.T) {
	translator := &validatex.Translator{
		Messages: validatex.Messages{"required": "{field} required"},
	}
	errs := validatex.Struct(translator, "parent", translatedNested{})
	if len(errs) != 1 || errs[0].Path != "parent.code" || errs[0].Msg != "code required" {
		t.Fatalf("unexpected translated nested error: %+v", errs)
	}
}

// TestFieldPathDoesNotDecodePercentEncoding 验证普通字段路径不会被解码
func TestFieldPathDoesNotDecodePercentEncoding(t *testing.T) {
	translator := &validatex.Translator{
		Messages: validatex.Messages{"required": "{field}"},
	}
	errs := validatex.Str(translator, "user%20name", "").Required().Errors()
	if len(errs) != 1 || errs[0].Msg != "user%20name" {
		t.Fatalf("unexpected field display: %+v", errs)
	}
}

// TestNestedValidationSkipsTypedNilPointer 验证指针专用嵌套API会跳过nil
func TestNestedValidationSkipsTypedNilPointer(t *testing.T) {
	var address *pointerAddress
	if errs := validatex.StructPointer(nil, "address", address); errs != nil {
		t.Fatalf("unexpected struct errors for typed nil pointer: %v", errs)
	}
	if errs := validatex.SliceStructPointer(nil, "items", []*pointerAddress{nil}); errs != nil {
		t.Fatalf("unexpected slice errors for typed nil pointer: %v", errs)
	}
	if errs := validatex.MapStructPointer(nil, "groups", map[string]*pointerAddress{"admin": nil}); errs != nil {
		t.Fatalf("unexpected map errors for typed nil pointer: %v", errs)
	}

	if errs := validatex.StructPointer(nil, "address", &pointerAddress{}); len(errs) != 1 || errs[0].Path != "address.unexpected" {
		t.Fatalf("expected non-nil pointer validation, got %v", errs)
	}
	if errs := validatex.MapStructPointer(nil, "groups", map[string]*pointerAddress{"a.b": {}}); len(errs) != 1 || errs[0].Path != "groups.a%2Eb.unexpected" {
		t.Fatalf("expected escaped pointer map path, got %v", errs)
	}
}

// TestMapStructErrorsAreDeterministic 验证Map结构体错误顺序稳定
func TestMapStructErrorsAreDeterministic(t *testing.T) {
	translator := &validatex.Translator{
		Messages: validatex.Messages{"required": "{field} required"},
	}
	addresses := map[string]testAddress{
		"beta":  {},
		"alpha": {},
	}

	for range 20 {
		errs := validatex.MapStruct(translator, "groups", addresses)
		if len(errs) != 2 || errs[0].Path != "groups.alpha.province" || errs[1].Path != "groups.beta.province" {
			t.Fatalf("unexpected map error order: %v", errs)
		}
	}
}

// TestMapStructEscapesPathSeparators 验证 Map 键不会破坏错误路径。
func TestMapStructEscapesPathSeparators(t *testing.T) {
	addresses := map[string]testAddress{
		"":       {},
		"*":      {},
		"a%b":    {},
		"a.b":    {},
		"a[b]":   {},
		"a%2Eb":  {},
		"a[0].b": {},
	}

	errs := validatex.MapStruct(nil, "groups", addresses)
	expected := map[string]bool{
		"groups.%2A.province":          true,
		"groups.%empty.province":       true,
		"groups.a%25b.province":        true,
		"groups.a%252Eb.province":      true,
		"groups.a%2Eb.province":        true,
		"groups.a%5B0%5D%2Eb.province": true,
		"groups.a%5Bb%5D.province":     true,
	}
	if len(errs) != len(expected) {
		t.Fatalf("expected %d errors, got %v", len(expected), errs)
	}
	seen := make(map[string]bool, len(errs))
	for _, err := range errs {
		if !expected[err.Path] {
			t.Fatalf("unexpected escaped path: %q", err.Path)
		}
		if seen[err.Path] {
			t.Fatalf("duplicate escaped path: %q", err.Path)
		}
		seen[err.Path] = true
	}
}

// TestMapEachEscapesKeysAndDisplaysDecodedField 验证Map键转义和显示名
func TestMapEachEscapesKeysAndDisplaysDecodedField(t *testing.T) {
	translator := &validatex.Translator{
		Messages: validatex.Messages{
			"data.*.required":     "{field} required",
			"data.a%2Eb.required": "{field} exact",
		},
	}
	values := map[string]string{
		"":    "",
		"a.b": "",
	}

	errs := validatex.Map(translator, "data", values).Each(
		"required",
		"",
		func(_ string, value string) bool {
			return value != ""
		},
	).Errors()
	if len(errs) != 2 {
		t.Fatalf("expected two map errors, got %v", errs)
	}
	if errs[0].Path != "data.%empty" || errs[0].Msg != "%empty required" {
		t.Fatalf("unexpected empty key error: %+v", errs[0])
	}
	if errs[1].Path != "data.a%2Eb" || errs[1].Msg != "a.b exact" {
		t.Fatalf("unexpected escaped key error: %+v", errs[1])
	}
}

// TestMapEachKeepsNamedKeyPathsUnique 验证命名键不会产生重复路径
func TestMapEachKeepsNamedKeyPathsUnique(t *testing.T) {
	values := map[collidingMapKey]string{
		"first":  "",
		"second": "",
	}
	errs := validatex.Map(validatex.NewTranslator(), "data", values).Each(
		"required",
		"",
		func(_ collidingMapKey, value string) bool {
			return value != ""
		},
	).Errors()
	if len(errs) != 2 || errs[0].Path == errs[1].Path {
		t.Fatalf("expected unique named map key paths, got %+v", errs)
	}
}

// TestMapSupportsNumericAndArrayKeys 验证通用可比较Map键
func TestMapSupportsNumericAndArrayKeys(t *testing.T) {
	translator := validatex.NewTranslator()
	arrayValues := map[[3]int]string{
		{1, 23, 4}: "",
	}
	errs := validatex.Map(translator, "data", arrayValues).Each(
		"required",
		"",
		func(_ [3]int, value string) bool {
			return value != ""
		},
	).Errors()
	if len(errs) != 1 || errs[0].Path != "data.%5B3%5Dint{1, 23, 4}" {
		t.Fatalf("unexpected array key path: %v", errs)
	}

	numberValues := map[int]testAddress{
		23: {},
		1:  {},
	}
	errs = validatex.MapStruct(translator, "items", numberValues)
	if len(errs) != 2 || errs[0].Path != "items.1.province" || errs[1].Path != "items.23.province" {
		t.Fatalf("unexpected numeric key paths: %v", errs)
	}

	translator = &validatex.Translator{
		Messages: validatex.Messages{"required": "{field} required"},
	}
	nestedValues := map[[3]int]mapKeyValue{
		{1, 23, 4}: "",
	}
	errs = validatex.MapStruct(translator, "data", nestedValues)
	if len(errs) != 1 || errs[0].Msg != "[3]int{1, 23, 4} required" {
		t.Fatalf("unexpected direct map value error: %v", errs)
	}

	mixedValues := map[any]string{
		nil:            "",
		"1":            "",
		1:              "",
		namedIntKey(1): "",
	}
	errs = validatex.Map(validatex.NewTranslator(), "data", mixedValues).Each(
		"required",
		"",
		func(_ any, value string) bool {
			return value != ""
		},
	).Errors()
	if len(errs) != 4 {
		t.Fatalf("expected four interface key errors, got %v", errs)
	}
	seenPaths := make(map[string]bool, len(errs))
	for _, validationError := range errs {
		if seenPaths[validationError.Path] {
			t.Fatalf("duplicate interface map key path: %v", errs)
		}
		seenPaths[validationError.Path] = true
	}
	if len(seenPaths) != 4 {
		t.Fatalf("interface map keys should produce unique paths: %v", errs)
	}
}

// TestErrorsError 验证错误字符串和Map渲染
func TestErrorsError(t *testing.T) {
	errs := validatex.Errors{
		{Path: "name", Msg: "required"},
		{Path: "email", Msg: "invalid"},
	}
	if message := errs.Error(); message != "[name: required] [email: invalid]" {
		t.Fatalf("unexpected error string: %s", message)
	}
	if message := (validatex.Errors{}).Error(); message != "" {
		t.Fatalf("expected empty error string, got %q", message)
	}
	ruleOnly := validatex.Errors{{Path: "age", Rule: "required"}}
	if message := ruleOnly.Error(); message != "[age: required]" {
		t.Fatalf("expected rule fallback in error string, got %q", message)
	}
	if message := ruleOnly.ToMap()["age"][0]; message != "required" {
		t.Fatalf("expected rule fallback in error map, got %q", message)
	}
}

// TestFieldErrorIndexes 验证错误路径中的切片下标提取
func TestFieldErrorIndexes(t *testing.T) {
	tests := []struct {
		name string
		path string
		want []int
	}{
		{name: "no index", path: "name"},
		{name: "single index", path: "users[2].email", want: []int{2}},
		{name: "nested indexes", path: "users[2].items[4].email", want: []int{2, 4}},
		{name: "adjacent indexes", path: "matrix[7][9]", want: []int{7, 9}},
		{name: "invalid index", path: "users[first].email"},
		{name: "unclosed index", path: "users[2.email"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := (validatex.FieldError{Path: test.path}).Indexes()
			if !slices.Equal(got, test.want) {
				t.Fatalf("Indexes() = %v, want %v", got, test.want)
			}
		})
	}
}

type wildcardOption struct {
	Label string
}

// Validate 为通配符测试实现 validatex.Validatable。
func (option wildcardOption) Validate(translator *validatex.Translator) validatex.Errors {
	return validatex.Check(
		validatex.Str(translator, "label", option.Label).Required().Errors(),
	)
}

type wildcardUser struct {
	Name    string
	Options map[string]wildcardOption
}

// Validate 为通配符测试实现validatex.Validatable
func (user wildcardUser) Validate(translator *validatex.Translator) validatex.Errors {
	return validatex.Check(
		validatex.Str(translator, "name", user.Name).Required().Errors(),
		validatex.MapStruct(translator, "options", user.Options),
	)
}

type wildcardRequest struct {
	Users  []wildcardUser
	Groups map[string]wildcardOption
}

// Validate 为通配符测试实现validatex.Validatable
func (request wildcardRequest) Validate(translator *validatex.Translator) validatex.Errors {
	return validatex.Check(
		validatex.SliceStruct(translator, "users", request.Users),
		validatex.MapStruct(translator, "groups", request.Groups),
	)
}

// TestWildcardPaths 验证动态和预编译通配符翻译
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

	compiled := translator.Compile()
	translator.Messages["users.*.name.required"] = "changed"
	translator.Translates["users.*.name"] = "changed"
	errs = request.Validate(compiled)
	if len(errs) != 3 ||
		errs[0].Path != "users[0].name" || errs[0].Msg != "用户名不能为空" ||
		errs[1].Path != "users[0].options.theme.label" || errs[1].Msg != "选项标签不能为空" ||
		errs[2].Path != "groups.admin.label" || errs[2].Msg != "群组标签不能为空" {
		t.Fatalf("unexpected compiled wildcard errors: %v", errs)
	}

	errs = validatex.Map(compiled, "data", data).Each(
		"required",
		"",
		func(_ string, value string) bool {
			return value != ""
		},
	).Errors()
	if len(errs) != 1 || errs[0].Path != "data.name" || errs[0].Msg != "映射值必填" {
		t.Fatalf("unexpected compiled wildcard map error: %v", errs)
	}
}
