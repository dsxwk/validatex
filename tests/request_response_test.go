package tests

import (
	"net/http"
	"testing"

	"github.com/dsxwk/validatex"
)

type genericRequest struct {
	Name string
}

// Translates 实现validatex.TranslatesProvider
func (r *genericRequest) Translates() validatex.Translates {
	return validatex.Translates{"name": "名称"}
}

// Messages 实现validatex.MessagesProvider
func (r *genericRequest) Messages() validatex.Messages {
	return validatex.Messages{"required": "{field}必填"}
}

// Validate 实现validatex.Validatable
func (r *genericRequest) Validate(translator *validatex.Translator) validatex.Errors {
	return validatex.Check(
		validatex.Str(translator, "name", r.Name).Required().Errors(),
	)
}

type unexpectedValidationRequest struct{}

// Validate 在被意外调用时返回错误
func (*unexpectedValidationRequest) Validate(*validatex.Translator) validatex.Errors {
	return validatex.Errors{{Path: "unexpected"}}
}

type unexpectedSceneValidationRequest struct{}

// Validate 在被意外调用时返回错误
func (*unexpectedSceneValidationRequest) Validate(*validatex.Translator) validatex.Errors {
	return validatex.Errors{{Path: "unexpected"}}
}

// ValidateWithScene 在被意外调用时返回错误
func (*unexpectedSceneValidationRequest) ValidateWithScene(string, *validatex.Translator) validatex.Errors {
	return validatex.Errors{{Path: "unexpected"}}
}

type unexpectedTranslatorRequest struct{}

// Translates 在被意外调用时返回翻译配置
func (*unexpectedTranslatorRequest) Translates() validatex.Translates {
	return validatex.Translates{"unexpected": "translator"}
}

// Messages 在被意外调用时返回消息配置
func (*unexpectedTranslatorRequest) Messages() validatex.Messages {
	return validatex.Messages{"unexpected": "message"}
}

type valueValidationRequest struct{}

// Validate 使用值接收者实现validatex.Validatable
func (valueValidationRequest) Validate(*validatex.Translator) validatex.Errors {
	return validatex.Errors{{Path: "value"}}
}

type sceneOnlyRequest struct{}

// ValidateWithScene 实现validatex.SceneValidatable
func (*sceneOnlyRequest) ValidateWithScene(scene string, _ *validatex.Translator) validatex.Errors {
	return validatex.Errors{{Path: "scene", Msg: scene}}
}

// TestValidateRequestSkipsTypedNilPointer 验证类型化nil请求会被跳过
func TestValidateRequestSkipsTypedNilPointer(t *testing.T) {
	tests := []struct {
		name                   string
		validate               func() validatex.Errors
		validateWithTranslator func(*validatex.Translator) validatex.Errors
	}{
		{
			name: "validate",
			validate: func() validatex.Errors {
				var request *unexpectedValidationRequest
				return validatex.ValidateRequest(request)
			},
			validateWithTranslator: func(translator *validatex.Translator) validatex.Errors {
				var request *unexpectedValidationRequest
				return validatex.ValidateRequestWithTranslator(request, translator)
			},
		},
		{
			name: "validate with scene",
			validate: func() validatex.Errors {
				var request *unexpectedSceneValidationRequest
				return validatex.ValidateRequest(request)
			},
			validateWithTranslator: func(translator *validatex.Translator) validatex.Errors {
				var request *unexpectedSceneValidationRequest
				return validatex.ValidateRequestWithTranslator(request, translator)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if errs := test.validate(); errs != nil {
				t.Fatalf("unexpected validation errors for typed nil pointer: %v", errs)
			}
			if errs := test.validateWithTranslator(validatex.NewTranslator()); errs != nil {
				t.Fatalf("unexpected validation errors for typed nil pointer with translator: %v", errs)
			}
		})
	}
}

// TestValidateRequestSupportsValueReceiverThroughPointer 验证值接收者方法可通过指针调用
func TestValidateRequestSupportsValueReceiverThroughPointer(t *testing.T) {
	errs := validatex.ValidateRequest(&valueValidationRequest{})
	if len(errs) != 1 || errs[0].Path != "value" {
		t.Fatalf("unexpected value receiver validation errors: %v", errs)
	}
}

// TestValidateRequestSupportsSceneOnlyRequest 验证仅支持场景的请求
func TestValidateRequestSupportsSceneOnlyRequest(t *testing.T) {
	errs := validatex.ValidateSceneRequest(&sceneOnlyRequest{}, "create")
	if len(errs) != 1 || errs[0].Path != "scene" || errs[0].Msg != "create" {
		t.Fatalf("unexpected scene-only validation errors: %v", errs)
	}
	errs = validatex.ValidateSceneRequestWithTranslator(&sceneOnlyRequest{}, "update", validatex.NewTranslator())
	if len(errs) != 1 || errs[0].Msg != "update" {
		t.Fatalf("unexpected scene-only translator validation errors: %v", errs)
	}
}

// TestNewRequestTranslatorSkipsTypedNilRequest 验证nil请求不会提供翻译
func TestNewRequestTranslatorSkipsTypedNilRequest(t *testing.T) {
	var request *unexpectedTranslatorRequest
	translator := validatex.NewRequestTranslator(request)
	if len(translator.Translates) != 0 || len(translator.Messages) != 0 {
		t.Fatalf("typed nil request should not provide translations: %+v", translator)
	}
}

// TestNewRequestTranslatorLoadsRequestTranslations 验证请求翻译配置会被加载
func TestNewRequestTranslatorLoadsRequestTranslations(t *testing.T) {
	translator := validatex.NewRequestTranslator(&genericRequest{})
	if translator.Translates["name"] != "名称" || translator.Messages["required"] != "{field}必填" {
		t.Fatalf("unexpected request translations: %+v", translator)
	}
}

// TestValidateRequest 验证带请求翻译的请求校验
func TestValidateRequest(t *testing.T) {
	errs := validatex.ValidateRequest(&genericRequest{})
	if len(errs) != 1 || errs[0].Path != "name" || errs[0].Msg != "名称必填" {
		t.Fatalf("unexpected generic validation errors: %v", errs)
	}
}

// TestErrorResponse 验证校验错误和绑定错误响应
func TestErrorResponse(t *testing.T) {
	errs := validatex.Errors{
		{Path: "name", Rule: "required", Msg: "名称必填"},
	}
	validationResponse := validatex.NewValidationErrorResponse(errs, "参数校验失败")
	if validationResponse.Code != http.StatusUnprocessableEntity || validationResponse.Msg != "参数校验失败" {
		t.Fatalf("unexpected validation response: %+v", validationResponse)
	}
	data, ok := validationResponse.Data.(map[string][]string)
	if !ok || data["name"][0] != "名称必填" {
		t.Fatalf("unexpected validation response data: %+v", validationResponse.Data)
	}

	bindResponse := validatex.NewBindErrorResponse("invalid json", "参数格式错误")
	if bindResponse.Code != http.StatusBadRequest || bindResponse.Data != "invalid json" {
		t.Fatalf("unexpected bind response: %+v", bindResponse)
	}
}
