package tests

import (
	"regexp"
	"testing"

	"github.com/dsxwk/validatex"
)

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
}

func TestPrecompiledRegexp(t *testing.T) {
	translator := validatex.NewTranslator()
	pattern := regexp.MustCompile(`^\d+$`)
	if errs := validatex.Str(translator, "value", "123").Regexp(pattern).Errors(); len(errs) != 0 {
		t.Fatalf("expected regexp to pass, got %v", errs)
	}
	if errs := validatex.Str(translator, "value", "abc").Regexp(pattern).Errors(); len(errs) != 1 {
		t.Fatalf("expected regexp to fail, got %v", errs)
	}
}
