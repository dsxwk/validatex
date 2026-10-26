# validatex - 类型安全的高性能 Go 校验库

[English](readme.md) | 简体中文

`validatex` 是一个高性能、无第三方依赖的 Go 数据校验库，适用于 API 请求校验、表单校验、配置校验、CLI 参数和领域命令校验。它提供类型安全的泛型校验器、场景规则、翻译消息、嵌套错误路径、通配符规则、与 Web 框架解耦的校验结果和可选通用 HTTP 错误结构。

核心包运行时仅使用 Go 标准库，HTTP 框架绑定由使用方自行集成。

## 功能特性

- 支持字符串、整数、浮点数、布尔值、切片和 Map 的类型安全泛型校验。
- 无第三方运行时依赖，核心校验路径不使用反射。
- 支持结构体、切片和 Map 嵌套校验，并生成完整错误路径。
- 支持 `*` 和 `**` 通配符规则，适合重复数据结构。
- 支持带类型判断函数、翻译消息和命名占位符的自定义规则。
- 支持创建、更新及其他业务场景校验。
- 提供与 Web 框架解耦的校验结果和可选 HTTP 响应辅助结构，可集成 Gin、Echo、Fiber、`net/http`、gRPC 和 CLI。
- 内置无第三方依赖的性能基准。

## 环境要求

- Go `1.27.0` 或更高版本。
- 模块路径为 `github.com/dsxwk/validatex`。
- 核心包运行时仅依赖 Go 标准库。

## 安装

```bash
go get github.com/dsxwk/validatex@latest
```

```go
import "github.com/dsxwk/validatex"
```

## 目录结构

```text
validatex/
├── readme.md                # 英文文档
├── readme_zh.md             # 简体中文文档
├── go.mod                   # 模块定义及 Go 1.27.0 版本要求
├── validate.go              # 核心接口,错误结构和聚合逻辑
├── translator.go            # 消息查找及通配符翻译
├── request.go               # 请求校验和翻译器创建
├── response.go              # 框架无关的错误响应
├── conditional.go           # When和Nullable条件校验
├── field.go                 # 泛型字段基础及基础类型构造器
├── patterns.go              # 预编译正则和日期布局
├── string_rules.go          # 字符串格式及内容规则
├── number_rules.go          # 整数及浮点数规则
├── collection_rules.go      # 布尔值,切片和Map规则
├── nested.go                # 嵌套结构体和Map校验
└── tests/
    ├── benchmark_test.go           # 无第三方依赖的性能基准
    ├── internal_test.go            # 模糊、分配、键及格式化测试
    ├── request_response_test.go    # 请求及响应测试
    ├── rules_test.go               # 内置、条件及自定义规则测试
    └── validate_test.go            # 核心、嵌套、通配符及翻译测试
```

## 快速开始

```go
type CreateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Age      int32  `json:"age"`
}

func (r *CreateUserRequest) Translates() validatex.Translates {
	return validatex.Translates{
		"username": "用户名",
		"email":    "邮箱",
		"age":      "年龄",
	}
}

func (r *CreateUserRequest) Messages() validatex.Messages {
	return validatex.Messages{
		"required": "{field}必填",
		"min":      "{field}不能小于{min}",
		"email":    "{field}格式不正确",
	}
}

func (r *CreateUserRequest) Validate(translator *validatex.Translator) validatex.Errors {
	return validatex.Check(
		validatex.Str(translator, "username", r.Username).Required().Min(3).Errors(),
		validatex.Str(translator, "email", r.Email).Required().Email().Errors(),
		validatex.Int(translator, "age", r.Age).Min(18).Errors(),
	)
}
```

## 通用校验

核心包不依赖 Gin，可以直接在 `net/http`、gRPC、命令行或其他 Go 代码中使用：

```go
errs := validatex.ValidateRequest(&request)
if len(errs) > 0 {
	response := validatex.NewValidationErrorResponse(errs, "参数校验失败")
	w.WriteHeader(response.Code)
	_ = json.NewEncoder(w).Encode(response)
	return
}
```

需要复用已有翻译器时：

```go
translator := validatex.NewRequestTranslator(&request)
errs := validatex.ValidateRequestWithTranslator(&request, translator)
```

## 内置规则

### 字符串

```go
validatex.Str(translator, "username", value).
	Required().
	Min(3).
	Max(20).
	AlphaDash().
	Errors()
```

| 规则 | 说明 |
| --- | --- |
| `Required` | 去除首尾空白后不能为空 |
| `RequiredIf` | 条件成立时必填 |
| `RequiredUnless` | 条件不成立时必填 |
| `Prohibited` / `ProhibitedIf` / `ProhibitedUnless` | 禁止有值及条件禁止 |
| `Filled` | 去除首尾空白后必须有值 |
| `Min` / `Max` | 最小或最大字符数 |
| `Between` | 字符数范围，例如 `Between(3, 20)` |
| `Size` | 字符数等于指定值 |
| `Email` | 邮箱格式 |
| `URL` | URL 格式 |
| `IP` / `IPv4` / `IPv6` | IP 地址格式 |
| `MACAddress` | MAC 地址格式 |
| `UUID` / `ULID` | UUID 或 ULID 格式 |
| `HexColor` | 十六进制颜色格式 |
| `Timezone` | 合法时区名称 |
| `JSON` | 合法 JSON 字符串 |
| `Alpha` | 仅包含字母 |
| `AlphaDash` | 仅包含字母、数字、短横线和下划线 |
| `AlphaNum` | 仅包含字母和数字 |
| `ASCII` | 仅包含 ASCII 字符 |
| `Lowercase` / `Uppercase` | 小写或大写 |
| `StartsWith` / `EndsWith` | 以任一目标开头或结尾 |
| `DoesntStartWith` / `DoesntEndWith` | 不以任何目标开头或结尾 |
| `Contains` / `DoesntContain` | 包含或不包含任一目标 |
| `In` / `NotIn` | 属于或不属于指定集合 |
| `Regex` / `NotRegex` | 使用字符串正则匹配 |
| `Regexp` / `NotRegexp` | 使用预编译正则匹配，适合热点路径 |
| `Numeric` | 数字字符串 |
| `Integer` | 整数字符串 |
| `Decimal` | 小数位数范围，例如 `Decimal(2, 4)` |
| `Digits` / `DigitsBetween` | 数字字符串长度或长度范围 |
| `Accepted` / `Declined` | 接受值或拒绝值 |
| `AcceptedIf` / `DeclinedIf` | 条件成立时校验接受值或拒绝值 |
| `Boolean` | 布尔字符串 |
| `Same` / `Different` / `Confirmed` | 与目标值相同、不同或一致 |
| `Date` | 常用日期格式 |
| `DateFormat` | 按 Go 布局校验日期 |
| `DateEquals` | 与目标日期是同一天 |
| `Before` / `BeforeOrEqual` | 早于或不晚于目标日期 |
| `After` / `AfterOrEqual` | 晚于或不早于目标日期 |

字符串格式规则对空字符串默认跳过，需要必填时组合 `Required`。

### 整数

```go
validatex.Int(translator, "age", value).Min(18).Max(60).Errors()
```

`Int` 支持 `int`、`int32`、`int64`、无符号整数以及基于这些类型定义的自定义整数类型。

| 规则 | 说明 |
| --- | --- |
| `Min` / `Max` | 最小值或最大值 |
| `Between` | 数值范围 |
| `Size` | 等于指定值 |
| `Gt` / `Gte` / `Lt` / `Lte` | 大小比较 |
| `In` / `NotIn` | 属于或不属于指定集合 |
| `MultipleOf` | 为指定值的倍数 |
| `Digits` / `DigitsBetween` | 十进制位数或位数范围 |
| `Positive` / `Negative` | 大于零或小于零 |
| `Required` / `RequiredIf` / `RequiredUnless` | 必填和条件必填，零值视为空 |
| `Prohibited` / `ProhibitedIf` / `ProhibitedUnless` | 禁止有值及条件禁止 |
| `Same` / `Different` | 与目标值相同或不同 |

### 浮点数

```go
validatex.Float(translator, "price", value).
	Finite().
	Between(0.01, 9999.99).
	MultipleOf(0.01).
	Errors()
```

`Float` 支持 `float32`、`float64`，提供 `Required`、`Finite`、`Min`、`Max`、`Between`、`Size`、`Gt`、`Gte`、`Lt`、`Lte`、`In`、`NotIn`、`MultipleOf`、`Positive`、`Negative`、`Same`、`Different`、条件必填和条件禁止规则。

### 布尔值

```go
validatex.Bool(translator, "enabled", value).Accepted().Errors()
```

提供 `Accepted`、`Declined`、`AcceptedIf`、`DeclinedIf`、`Prohibited`、`ProhibitedIf`、`ProhibitedUnless`、`True`、`False`、`Same`、`Different` 和自定义 `Rule`。

### 切片

```go
validatex.Slice(translator, "ids", values).
	Required().
	Distinct().
	Max(100).
	Errors()
```

提供 `Required`、`RequiredIf`、`RequiredUnless`、`Prohibited`、`ProhibitedIf`、`ProhibitedUnless`、`Min`、`Max`、`Between`、`Size`、`Distinct`、`Contains`、`DoesntContain`、`Each` 和自定义 `Rule`。

`Each` 会为每个失败元素生成独立错误路径：

```go
validatex.Slice(translator, "items", values).Each(
	"positive",
	"{field}必须大于零",
	func(value int) bool {
		return value > 0
	},
).Errors()
```

错误路径格式为 `items[1]`。

### Map

```go
validatex.Map(translator, "data", values).
	HasKeys("name", "status").
	AllowedKeys("name", "status", "remark").
	Between(2, 3).
	Errors()
```

提供 `Required`、`RequiredIf`、`RequiredUnless`、`Prohibited`、`ProhibitedIf`、`ProhibitedUnless`、`Min`、`Max`、`Between`、`Size`、`HasKeys`、`AllowedKeys`、`Each` 和自定义 `Rule`。

Map key 支持任意 Go `comparable` 类型，包括整数、命名标量类型、结构体和 `[3]int` 这类定长数组。只有产生错误路径时才会格式化 key，校验通过的热路径不承担这部分开销。

```go
values := map[[3]int]string{
	{1, 23, 4}: "",
}

errs := validatex.Map(translator, "data", values).Each(
	"required",
	"",
	func(_ [3]int, value string) bool {
		return value != ""
	},
).Errors()
```

`Each` 的错误路径使用Map键名：

```go
validatex.Map(translator, "labels", values).Each(
	"required",
	"",
	func(_ string, value string) bool {
		return value != ""
	},
).Errors()
```

错误路径格式为 `labels.zh`。

## 自定义规则

`Rule` 接收规则名、消息和判断函数。判断函数返回 `true` 表示校验通过。

```go
errs := validatex.Str(translator, "code", value).
	Rule("prefix", "{field}必须以{prefix}开头", func(value string) bool {
		return len(value) > 0 && value[0] == 'A'
	}, map[string]any{"prefix": "A"}).
	Errors()
```

```go
errs := validatex.Int(translator, "count", value).
	Rule("even", "{field}必须是偶数", func(value int32) bool {
		return value%2 == 0
	}).
	Errors()
```

`Rule` 的参数顺序为 `规则名`、`消息`、`判断函数`、`占位符参数`。消息传空字符串时使用 `Translator.Messages` 中的翻译配置。

### Map参数

`map[string]any` 用作消息命名占位符，键名与 `{键名}` 对应：

```go
errs := validatex.Int(translator, "age", value).
	Rule(
		"range",
		"{field}必须在{min}到{max}之间",
		func(value int32) bool {
			return value >= 18 && value <= 60
		},
		map[string]any{
			"min": 18,
			"max": 60,
		},
	).
	Errors()
```

当字段配置了 `Translates` 时，`{field}` 会替换成对应字段名称：

```go
translator.Translates = validatex.Translates{
	"age": "年龄",
}
```

最终消息为：

```text
年龄必须在18到60之间
```

`map[string]any` 也可以与 `Translator.Messages` 配合使用：

```go
translator.Messages = validatex.Messages{
	"age.range": "{field}必须在{min}到{max}之间",
}

errs := validatex.Int(translator, "age", value).
	Rule(
		"range",
		"",
		func(value int32) bool {
			return value >= 18 && value <= 60
		},
		map[string]any{
			"min": 18,
			"max": 60,
		},
	).
	Errors()
```

消息中的内置占位符包括：

- `{field}`：翻译后的字段名称，未配置时使用路径最后一段。
- `{path}`：完整字段路径，例如 `items[0].age`。
- `{rule}`：当前规则名。
- `{key}`：`map[string]any` 中对应的命名参数。
- `{0}`、`{1}`：按位置传入的参数。

配置翻译时，规则消息按以下顺序读取：

1. `字段路径.规则名`，例如 `username.required`。
2. 通配符路径规则，例如 `users.*.name.required`。
3. `规则名`，例如 `required`。
4. 未配置时优先使用字段显示名生成 `字段显示名 is 规则名`；没有显示名时回退到字段路径。

消息支持 `{field}`、`{path}`、`{rule}`、`{0}`、`{1}` 等占位符，也支持 `map[string]any` 中的命名占位符。

### 编译翻译表

通配符规则较多时，可以在启动阶段编译一次：

```go
translator := validatex.NewRequestTranslator(&request).Compile()
errs := validatex.ValidateRequestWithTranslator(&request, translator)
```

`Compile` 返回包含预解析通配符索引和翻译表副本的快照。编译后的翻译器及其 Map 应视为只读；只要不修改，就可以安全地并发读取。编译后修改原始 Map 不会影响该快照。

`Check` 在只有一组非空错误时会直接返回该切片。返回的 `Errors` 以及原始校验结果切片都应视为只读。

## 条件与可空

`When` 用于条件校验：

```go
errs := validatex.When(request.Status == "published", func() validatex.Errors {
	return validatex.Check(
		validatex.Str(translator, "content", request.Content).Required().Errors(),
	)
})
```

`Nullable` 用于可空指针：

```go
errs := validatex.Nullable(request.Email, func(value string) validatex.Errors {
	return validatex.Check(
		validatex.Str(translator, "email", value).Required().Email().Errors(),
	)
})
```

请求校验入口接收请求指针。`nil` 指针会跳过校验；非 `nil` 指针会使用解引用后的值执行校验。如果 JSON 字段缺失和显式 `null` 需要区分，应在绑定层使用额外字段记录存在状态。

## 嵌套校验

嵌套结构体和切片需要实现 `Validate(translator *validatex.Translator) validatex.Errors`。

```go
type Address struct {
	Province string `json:"province"`
}

func (a Address) Validate(translator *validatex.Translator) validatex.Errors {
	return validatex.Check(
		validatex.Str(translator, "province", a.Province).Required().Errors(),
	)
}

func (r *CreateUserRequest) Validate(translator *validatex.Translator) validatex.Errors {
	return validatex.Check(
		validatex.Struct(translator, "address", r.Address),
		validatex.SliceStruct(translator, "items", r.Items),
		validatex.MapStruct(translator, "groups", r.Groups),
	)
}
```

嵌套错误路径会自动拼接为 `address.province`、`items[0].province`、`groups.admin.province`。

可空嵌套指针使用 `StructPointer`、`SliceStructPointer` 和 `MapStructPointer`。这些函数通过泛型指针约束识别 nil，不依赖反射：

```go
validatex.StructPointer(translator, "profile", request.Profile)
validatex.SliceStructPointer(translator, "items", request.Items)
validatex.MapStructPointer(translator, "groups", request.Groups)
```

Map key 在安全时保持原样；路径保留字符 `.`、`[`、`]`、`%`、`*` 会进行百分号编码（`%2E`、`%5B`、`%5D`、`%25`、`%2A`），空 key 编码为 `%empty`。非标量 key 使用 Go 类型表示，避免不同类型或值生成相同路径；例如 `[3]int{1, 23, 4}` 会生成 `data.%5B3%5Dint{1, 23, 4}`。翻译配置使用编码后的路径，`{field}` 仍展示可读的 key。

值元素使用 `Struct`、`SliceStruct`、`MapStruct`；指针元素使用 `StructPointer`、`SliceStructPointer`、`MapStructPointer`。指针版本会显式跳过 nil。

### 通配符路径

结构体切片、结构体Map和Map逐项校验支持通配符翻译：

```go
translator.Messages = validatex.Messages{
	"users.*.name.required":       "{field}不能为空",
	"groups.*.province.required":  "{field}不能为空",
	"users.**.email.required":     "{field}不能为空",
}

translator.Translates = validatex.Translates{
	"users.*.name":              "用户名",
	"groups.*.province":         "省份",
	"users.**.email":            "邮箱",
}
```

- `*` 匹配一个路径段，例如 `users.*.name` 匹配 `users[0].name`。
- `**` 匹配零个或多个路径段，例如 `users.**.email` 匹配 `users[0].profile.email`。
- 实际错误路径保留具体下标和Map键，例如 `users[0].name`、`groups.admin.province`。
- 多个通配符规则同时匹配时，优先使用固定路径段更多、通配符更少的规则。

## 场景校验

实现 `ValidateWithScene` 后，使用 `ValidateSceneRequest` 或 `ValidateSceneRequestWithTranslator` 调用对应场景。

```go
func (r *UpdateUserRequest) ValidateWithScene(scene string, translator *validatex.Translator) validatex.Errors {
	switch scene {
	case "create":
		return validatex.Check(
			validatex.Str(translator, "username", r.Username).Required().Errors(),
		)
	case "update":
		return validatex.Check(
			validatex.Int(translator, "id", r.ID).Min(1).Errors(),
		)
	default:
		return nil
	}
}
```

`ValidateRequest` 和 `ValidateRequestWithTranslator` 用于实现 `Validatable` 的请求；实现 `SceneValidatable` 的请求使用 `ValidateSceneRequest` 或 `ValidateSceneRequestWithTranslator`。

## 框架集成

框架绑定只需要组合 `NewRequestTranslator`、`ValidateRequestWithTranslator` 和通用响应：

```go
func Bind[T any, P validatex.ValidatablePointer[T]](c *gin.Context, request P) bool {
	translator := validatex.NewRequestTranslator(request)
	if err := c.ShouldBindJSON(request); err != nil {
		response := validatex.NewBindErrorResponse(err.Error(), "参数格式错误")
		c.JSON(response.Code, response)
		return false
	}

	errs := validatex.ValidateRequestWithTranslator(request, translator)
	if len(errs) > 0 {
		response := validatex.NewValidationErrorResponse(errs, "参数校验失败")
		c.JSON(response.Code, response)
		return false
	}
	return true
}
```

控制器调用：

```go
func (ctl *UserController) Create(c *gin.Context) {
	var request CreateUserRequest
	if !Bind(c, &request, "create") {
		return
	}

	_ = request
}
```

框架适配层可以自行决定：

- JSON 格式错误返回 HTTP `400`。
- 字段校验失败返回 HTTP `422`。
- 校验错误会按字段路径聚合到 `data`。

```json
{
  "code": 422,
  "msg": "参数校验失败",
  "data": {
    "email": ["邮箱格式不正确"],
    "address.province": ["省份必填"]
  }
}
```

## 性能基准

基准代码位于模块内的 `tests/benchmark_test.go`，覆盖相同用户名、邮箱字段的有效和无效输入：

```bash
go test -run '^$' -bench . -benchmem ./tests
```

2026-10-26 在 Windows、`i5-10400` 环境的三次基准中位数近似结果：

| 场景 | 耗时 | 内存 | 分配 |
| --- | ---: | ---: | ---: |
| 有效输入热路径 | 74 ns/op | 0 B/op | 0 |
| 有效输入完整入口 | 78 ns/op | 0 B/op | 0 |
| 无效输入热路径 | 448 ns/op | 480 B/op | 5 |
| 无效输入完整入口 | 404 ns/op | 480 B/op | 5 |
| 有效参数化规则 | 27 ns/op | 0 B/op | 0 |
| 无效参数化规则 | 258 ns/op | 224 B/op | 5 |
| 通配符查询，动态，128 条规则 | 26.7 µs/op | 12.5 KB/op | 261 |
| 通配符查询，编译后，128 条规则 | 6.1 µs/op | 256 B/op | 5 |
| 嵌套错误，128 条错误 | 78 µs/op | 71.5 KB/op | 771 |

基准覆盖两个字符串字段、参数化字符串和整数规则，以及通配符规则较多的翻译查询。热路径使用预建翻译器，完整入口在没有翻译提供者时直接传入 `nil` 翻译器。结果不代表大型嵌套结构、国际化、数据库规则或所有输入场景。基准不引入第三方依赖。

## 测试

测试代码位于模块内的 `tests` 目录：

```bash
go test ./...
```

## 常见问题

### validatex 是否依赖 Web 框架？

不依赖。核心模块不导入 Gin 或其他 Web 框架，框架适配由业务项目完成。

### validatex 是否使用反射？

核心校验路径使用泛型，不使用反射。

### 是否支持自定义规则？

支持。可使用 `Rule` 传入类型安全的判断函数，并配置位置参数或命名占位符。

### 是否支持嵌套通配符路径？

支持。`*` 匹配一个路径段，`**` 匹配多个路径段。

### 最低 Go 版本是多少？

需要 Go `1.27.0` 或更高版本。

## 许可证

本项目采用 Mozilla Public License 2.0。

完整许可证见 [LICENSE](LICENSE)，项目来源和修改说明见 [NOTICE](NOTICE)。修改后的 MPL 覆盖文件必须继续使用 MPL-2.0，保留原始声明，并注明进行了修改。
