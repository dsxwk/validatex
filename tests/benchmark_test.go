package tests

import (
	"testing"

	"github.com/dsxwk/validatex"
)

type validatexBenchRequest struct {
	Username string
	Email    string
}

func (r *validatexBenchRequest) Validate(translator *validatex.Translator) validatex.Errors {
	return validatex.Check(
		validatex.Str(translator, "username", r.Username).Required().Min(3).Errors(),
		validatex.Str(translator, "email", r.Email).Required().Email().Errors(),
	)
}

func BenchmarkValid(b *testing.B) {
	translator := validatex.NewTranslator()
	validValidatex := &validatexBenchRequest{
		Username: "codex",
		Email:    "codex@example.com",
	}
	b.Run("validatex", func(b *testing.B) {
		b.ReportAllocs()
		for index := 0; index < b.N; index++ {
			validatex.ValidateRequestWithTranslator(validValidatex, "", translator)
		}
	})
	b.Run("validatex-request", func(b *testing.B) {
		b.ReportAllocs()
		for index := 0; index < b.N; index++ {
			validatex.ValidateRequest(validValidatex, "")
		}
	})
}

func BenchmarkInvalid(b *testing.B) {
	translator := validatex.NewTranslator()
	invalidValidatex := &validatexBenchRequest{
		Username: "x",
		Email:    "invalid",
	}
	b.Run("validatex", func(b *testing.B) {
		b.ReportAllocs()
		for index := 0; index < b.N; index++ {
			if len(validatex.ValidateRequestWithTranslator(invalidValidatex, "", translator)) == 0 {
				b.Fatal("validation passed")
			}
		}
	})
	b.Run("validatex-request", func(b *testing.B) {
		b.ReportAllocs()
		for index := 0; index < b.N; index++ {
			if len(validatex.ValidateRequest(invalidValidatex, "")) == 0 {
				b.Fatal("validation passed")
			}
		}
	})
}
