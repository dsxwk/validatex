# validatex - Type-safe Go validation library

English | [简体中文](readme_zh.md)

`validatex` is a high-performance, zero-dependency Go validation library for API request validation, form validation, configuration validation, CLI input, and domain commands. It provides type-safe generic validators, scene-based rules, translated messages, nested error paths, wildcard rules, and framework-neutral validation results with optional generic HTTP error responses.

The core package uses only the Go standard library at runtime. HTTP framework integration belongs in the consuming application.

## Features

- Type-safe generic validators for strings, integers, floats, booleans, slices, and maps.
- Zero third-party runtime dependencies and no reflection in the core validation path.
- Nested struct, slice, and map validation with concrete error paths.
- Wildcard rules for repeated structures, including `*` and `**` path matching.
- Custom rules with typed predicates, translated messages, and named placeholders.
- Scene validation for create, update, and other application operations.
- Framework-neutral validation results with optional HTTP response helpers for Gin, Echo, Fiber, `net/http`, gRPC, and CLI tools.
- Built-in performance benchmarks with no external benchmark dependencies.

## Requirements

- Go `1.27.0` or newer.
- Module path: `github.com/dsxwk/validatex`.

## Installation

```bash
go get github.com/dsxwk/validatex@latest
```

```go
import "github.com/dsxwk/validatex"
```

## Project Structure

```text
validatex/
├── readme.md                # English documentation
├── readme_zh.md             # Simplified Chinese documentation
├── go.mod                   # Module definition and Go 1.27.0 requirement
├── validate.go              # Core interfaces, errors, and aggregator
├── translator.go            # Message lookup and wildcard translation
├── request.go               # Request validation and translator creation
├── response.go              # Framework-independent error responses
├── conditional.go           # When and Nullable helpers
├── field.go                 # Generic field base and primitive constructors
├── patterns.go              # Precompiled patterns and date layouts
├── string_rules.go          # String format and content rules
├── number_rules.go          # Integer and float rules
├── collection_rules.go      # Boolean, slice, and map rules
├── nested.go                # Nested struct and map validation
└── tests/
    ├── benchmark_test.go           # Self-contained performance benchmarks
    ├── internal_test.go            # Fuzz, allocation, key, and formatting tests
    ├── request_response_test.go    # Request and response tests
    ├── rules_test.go               # Built-in, conditional, and custom rules
    └── validate_test.go            # Core, nested, wildcard, and translation tests
```

## Quick Start

```go
type CreateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Age      int32  `json:"age"`
}

func (r *CreateUserRequest) Translates() validatex.Translates {
	return validatex.Translates{
		"username": "Username",
		"email":    "Email",
		"age":      "Age",
	}
}

func (r *CreateUserRequest) Messages() validatex.Messages {
	return validatex.Messages{
		"required": "{field} is required",
		"min":      "{field} must be at least {min}",
		"email":    "{field} is invalid",
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

## Framework-Independent Validation

The core package can be used from `net/http`, gRPC, CLI tools, jobs, or any other Go code:

```go
errs := validatex.ValidateRequest(&request)
if len(errs) > 0 {
	response := validatex.NewValidationErrorResponse(errs, "Validation failed")
	w.WriteHeader(response.Code)
	_ = json.NewEncoder(w).Encode(response)
	return
}
```

Reuse an existing translator when necessary:

```go
translator := validatex.NewRequestTranslator(&request)
errs := validatex.ValidateRequestWithTranslator(&request, translator)
```

## Built-In Rules

### Strings

```go
validatex.Str(translator, "username", value).
	Required().
	Min(3).
	Max(20).
	AlphaDash().
	Errors()
```

| Rules | Description |
| --- | --- |
| `Required`, `RequiredIf`, `RequiredUnless` | Required and conditionally required |
| `Prohibited`, `ProhibitedIf`, `ProhibitedUnless` | Prohibited and conditionally prohibited |
| `Filled` | Must contain a non-whitespace value |
| `Min`, `Max`, `Between`, `Size` | Character-length constraints |
| `Email`, `URL` | Email and URL formats |
| `IP`, `IPv4`, `IPv6`, `MACAddress` | Network address formats |
| `UUID`, `ULID`, `HexColor`, `Timezone` | Common identifier and format rules |
| `JSON`, `Numeric`, `Integer`, `Decimal` | Structured and numeric string rules |
| `Alpha`, `AlphaDash`, `AlphaNum`, `ASCII` | Character-set rules |
| `Lowercase`, `Uppercase` | Case rules |
| `StartsWith`, `EndsWith`, `Contains` | Content presence rules |
| `DoesntStartWith`, `DoesntEndWith`, `DoesntContain` | Negative content rules |
| `In`, `NotIn` | Set membership |
| `Regex`, `NotRegex` | Runtime string patterns |
| `Regexp`, `NotRegexp` | Precompiled patterns for hot paths |
| `Digits`, `DigitsBetween` | Digit-count constraints |
| `Accepted`, `Declined`, `Boolean` | Boolean-like string values |
| `Same`, `Different`, `Confirmed` | Value comparison |
| `Date`, `DateFormat`, `DateEquals` | Date parsing and equality |
| `Before`, `BeforeOrEqual`, `After`, `AfterOrEqual` | Date comparison |

Format rules skip empty strings by default. Combine them with `Required` when the value must be present.

### Integers

```go
validatex.Int(translator, "age", value).Min(18).Max(60).Errors()
```

`Int` supports `int`, `int32`, `int64`, unsigned integer types, and named types based on those underlying types.

Available rules include `Required`, `Min`, `Max`, `Between`, `Size`, `Gt`, `Gte`, `Lt`, `Lte`, `In`, `NotIn`, `MultipleOf`, `Digits`, `DigitsBetween`, `Positive`, `Negative`, `RequiredIf`, `RequiredUnless`, `ProhibitedIf`, `ProhibitedUnless`, `Same`, and `Different`.

### Floats

```go
validatex.Float(translator, "price", value).
	Finite().
	Between(0.01, 9999.99).
	MultipleOf(0.01).
	Errors()
```

`Float` supports `float32` and `float64`, with required, range, comparison, multiple, sign, set, conditional, and custom rules.

### Booleans

```go
validatex.Bool(translator, "enabled", value).Accepted().Errors()
```

Available rules include `Accepted`, `Declined`, `AcceptedIf`, `DeclinedIf`, `Prohibited`, `True`, `False`, `Same`, `Different`, and custom `Rule`.

### Slices

```go
validatex.Slice(translator, "ids", values).
	Required().
	Distinct().
	Max(100).
	Errors()
```

Slice rules include `Required`, `Min`, `Max`, `Between`, `Size`, `Distinct`, `Contains`, `DoesntContain`, `Each`, and custom `Rule`.

`Each` reports an indexed error path:

```go
validatex.Slice(translator, "items", values).Each(
	"positive",
	"{field} must be greater than zero",
	func(value int) bool {
		return value > 0
	},
).Errors()
```

### Maps

```go
validatex.Map(translator, "data", values).
	HasKeys("name", "status").
	AllowedKeys("name", "status", "remark").
	Between(2, 3).
	Errors()
```

Map rules include `Required`, `Min`, `Max`, `Between`, `Size`, `HasKeys`, `AllowedKeys`, `Each`, and custom `Rule`.

Map keys support any Go `comparable` type, including integers, named scalar types, structs, and fixed-size arrays such as `[3]int`. Keys are formatted only when an error path is needed, so valid input does not pay that cost.

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

## Custom Rules

`Rule` receives a rule name, message, predicate, and optional placeholder arguments. The predicate returns `true` when validation passes.

```go
errs := validatex.Str(translator, "code", value).
	Rule("prefix", "{field} must start with {prefix}", func(value string) bool {
		return len(value) > 0 && value[0] == 'A'
	}, map[string]any{"prefix": "A"}).
	Errors()
```

```go
errs := validatex.Int(translator, "count", value).
	Rule("even", "{field} must be even", func(value int32) bool {
		return value%2 == 0
	}).
	Errors()
```

Pass an empty message to use `Translator.Messages` instead:

```go
translator.Messages = validatex.Messages{
	"age.range": "{field} must be between {min} and {max}",
}
```

### Named Map Arguments

`map[string]any` values are available as named placeholders:

```go
errs := validatex.Int(translator, "age", value).
	Rule(
		"range",
		"{field} must be between {min} and {max}",
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

Supported placeholders include:

- `{field}`: translated field name.
- `{path}`: complete field path.
- `{rule}`: rule name.
- `{key}`: named map argument.
- `{0}`, `{1}`, etc.: positional arguments.

### Message Resolution

1. Exact field rule: `username.required`.
2. Wildcard field rule: `users.*.name.required`.
3. Global rule: `required`.
4. Fallback: translated field display plus `is rule`; if no display is available, the field path is used.

### Compiled Translation Tables

Wildcard-heavy translation tables can be compiled once at startup:

```go
translator := validatex.NewRequestTranslator(&request).Compile()
errs := validatex.ValidateRequestWithTranslator(&request, translator)
```

`Compile` returns a snapshot with pre-parsed wildcard patterns and copied translation maps. Treat the compiled translator and its maps as read-only; it is safe for concurrent validation reads as long as they are not modified. Changes to the original maps after compilation do not affect the snapshot.

`Check` may return one of its input slices directly when it is the only non-empty group. Treat returned `Errors` and the original validation result slices as read-only.

## Conditions and Nullable Values

Use `When` for conditional validation:

```go
errs := validatex.When(request.Status == "published", func() validatex.Errors {
	return validatex.Check(
		validatex.Str(translator, "content", request.Content).Required().Errors(),
	)
})
```

Use `Nullable` for pointers:

```go
errs := validatex.Nullable(request.Email, func(value string) validatex.Errors {
	return validatex.Check(
		validatex.Str(translator, "email", value).Required().Email().Errors(),
	)
})
```

Request validation entry points accept a request pointer. A `nil` pointer skips validation. If JSON missing fields must be distinguished from explicit `null`, record presence in the binding layer.

## Nested Validation

Nested values implement `Validate(translator *validatex.Translator) validatex.Errors`.

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

Error paths use concrete values, such as `address.province`, `items[0].province`, and `groups.admin.province`.

For optional nested pointers, use `StructPointer`, `SliceStructPointer`, and `MapStructPointer`. These functions use generic pointer constraints, so nil pointers are skipped without reflection:

```go
validatex.StructPointer(translator, "profile", request.Profile)
validatex.SliceStructPointer(translator, "items", request.Items)
validatex.MapStructPointer(translator, "groups", request.Groups)
```

Map keys keep their original form when safe. The reserved characters `.`, `[`, `]`, `%`, and `*` are percent-encoded in error paths (`%2E`, `%5B`, `%5D`, `%25`, `%2A`), and an empty key becomes `%empty`. Non-scalar keys use their Go representation so distinct types and values do not collapse into the same path; for example, `[3]int{1, 23, 4}` produces `data.%5B3%5Dint{1, 23, 4}`. Translation keys use the encoded path, while `{field}` still renders the readable key.

Use `Struct`, `SliceStruct`, and `MapStruct` for value elements, and `StructPointer`, `SliceStructPointer`, and `MapStructPointer` for pointer elements. The pointer variants explicitly skip nil.

### Wildcard Paths

Use wildcard message keys for repeated structures:

```go
translator.Messages = validatex.Messages{
	"users.*.name.required":      "{field} is required",
	"groups.*.province.required": "{field} is required",
	"users.**.email.required":    "{field} is required",
}

translator.Translates = validatex.Translates{
	"users.*.name":       "Username",
	"groups.*.province":  "Province",
	"users.**.email":     "Email",
}
```

- `*` matches exactly one path segment.
- `**` matches zero or more path segments.
- More specific patterns take precedence.

## Scene Validation

Implement `ValidateWithScene` to select rules by operation, then call `ValidateSceneRequest` or `ValidateSceneRequestWithTranslator`:

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

`ValidateRequest` and `ValidateRequestWithTranslator` are for requests implementing `Validatable`. Requests implementing `SceneValidatable` use `ValidateSceneRequest` or `ValidateSceneRequestWithTranslator`.

## Framework Integration

The package intentionally contains no Web framework dependency. A Gin integration can be implemented in the consuming application:

```go
func Bind[T any, P validatex.ValidatablePointer[T]](c *gin.Context, request P) bool {
	translator := validatex.NewRequestTranslator(request)
	if err := c.ShouldBindJSON(request); err != nil {
		response := validatex.NewBindErrorResponse(err.Error(), "Invalid request")
		c.JSON(response.Code, response)
		return false
	}

	errs := validatex.ValidateRequestWithTranslator(request, translator)
	if len(errs) > 0 {
		response := validatex.NewValidationErrorResponse(errs, "Validation failed")
		c.JSON(response.Code, response)
		return false
	}
	return true
}
```

The generic response structure can be encoded by any framework:

```json
{
  "code": 422,
  "msg": "Validation failed",
  "data": {
    "email": ["Email is invalid"],
    "address.province": ["Province is required"]
  }
}
```

## Performance

Benchmarks are in `tests/benchmark_test.go`.

```bash
go test -run '^$' -bench . -benchmem ./tests
```

Approximate median results on Windows with an Intel `i5-10400` on 2026-10-26:

| Scenario | Time | Memory | Allocations |
| --- | ---: | ---: | ---: |
| Valid hot path | 74 ns/op | 0 B/op | 0 |
| Valid request entry | 78 ns/op | 0 B/op | 0 |
| Invalid hot path | 448 ns/op | 480 B/op | 5 |
| Invalid request entry | 404 ns/op | 480 B/op | 5 |
| Valid parameterized rules | 27 ns/op | 0 B/op | 0 |
| Invalid parameterized rules | 258 ns/op | 224 B/op | 5 |
| Wildcard lookup, dynamic, 128 rules | 26.7 µs/op | 12.5 KB/op | 261 |
| Wildcard lookup, compiled, 128 rules | 6.1 µs/op | 256 B/op | 5 |
| Nested errors, 128 errors | 78 µs/op | 71.5 KB/op | 771 |

The benchmarks cover two string fields, parameterized string and integer rules, and wildcard-heavy translation lookup. They do not represent large nested structures, localization, database rules, or every input shape. Benchmarks introduce no third-party dependencies.

## Testing

Run from the module root:

```bash
go test ./...
go vet ./...
```

## FAQ

### Is validatex framework-independent?

Yes. The core module does not import Gin or another Web framework. Framework adapters belong in the consuming application.

### Does validatex use reflection?

The core validation path uses generics and does not use reflection.

### Does validatex support custom rules?

Yes. Use `Rule` with a typed predicate and optional positional or named message arguments.

### Does validatex support nested wildcard paths?

Yes. Use `*` for one path segment and `**` for multiple path segments.

### What Go version is required?

Go `1.27.0` or newer is required.

## Scope

The package covers validation rules that do not require external state. Rules such as `exists`, `unique`, `current_password`, file validation, MIME detection, and image dimensions depend on a database, authentication system, uploads, or application configuration and should be implemented by the consumer through custom `Rule` predicates.

## License

This project is licensed under the Mozilla Public License 2.0.

See [LICENSE](LICENSE) for the complete license text and [NOTICE](NOTICE) for project attribution information. Modified MPL-covered files must remain available under MPL-2.0, retain the original notices, and indicate that changes were made.
