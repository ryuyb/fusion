package middleware

import "github.com/ryuyb/fusion/internal/infrastructure/provider/validator"

type StructValidator struct {
	Validator *validator.Validator
}

func (v *StructValidator) Validate(out any) error {
	return v.Validator.Validate(out)
}
