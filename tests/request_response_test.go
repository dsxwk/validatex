package tests

import (
	"net/http"
	"testing"

	"github.com/dsxwk/validatex"
)

type genericRequest struct {
	Name string
}

func (r *genericRequest) Translates() validatex.Translates {
	return validatex.Translates{"name": "名称"}
}

func (r *genericRequest) Messages() validatex.Messages {
	return validatex.Messages{"required": "{field}必填"}
}

func (r *genericRequest) Validate(translator *validatex.Translator) validatex.Errors {
	return validatex.Check(
		validatex.Str(translator, "name", r.Name).Required().Errors(),
	)
}

type unexpectedValidationRequest struct{}

func (*unexpectedValidationRequest) Validate(*validatex.Translator) validatex.Errors {
	return validatex.Errors{{Path: "unexpected"}}
}

type unexpectedSceneValidationRequest struct{}

func (*unexpectedSceneValidationRequest) Validate(*validatex.Translator) validatex.Errors {
	return validatex.Errors{{Path: "unexpected"}}
}

func (*unexpectedSceneValidationRequest) ValidateWithScene(string, *validatex.Translator) validatex.Errors {
	return validatex.Errors{{Path: "unexpected"}}
}

type valueValidationRequest struct{}

func (valueValidationRequest) Validate(*validatex.Translator) validatex.Errors {
	return validatex.Errors{{Path: "value"}}
}

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
				return validatex.ValidateRequest(request, "create")
			},
			validateWithTranslator: func(translator *validatex.Translator) validatex.Errors {
				var request *unexpectedValidationRequest
				return validatex.ValidateRequestWithTranslator(request, "create", translator)
			},
		},
		{
			name: "validate with scene",
			validate: func() validatex.Errors {
				var request *unexpectedSceneValidationRequest
				return validatex.ValidateRequest(request, "create")
			},
			validateWithTranslator: func(translator *validatex.Translator) validatex.Errors {
				var request *unexpectedSceneValidationRequest
				return validatex.ValidateRequestWithTranslator(request, "create", translator)
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

func TestValidateRequestSupportsValueReceiverThroughPointer(t *testing.T) {
	errs := validatex.ValidateRequest(&valueValidationRequest{}, "create")
	if len(errs) != 1 || errs[0].Path != "value" {
		t.Fatalf("unexpected value receiver validation errors: %v", errs)
	}
}

func TestValidateRequest(t *testing.T) {
	errs := validatex.ValidateRequest(&genericRequest{}, "create")
	if len(errs) != 1 || errs[0].Path != "name" || errs[0].Msg != "名称必填" {
		t.Fatalf("unexpected generic validation errors: %v", errs)
	}
}

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
