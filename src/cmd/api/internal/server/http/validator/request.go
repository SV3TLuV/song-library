package validator

import "github.com/go-playground/validator"

type requestValidator struct {
	validator *validator.Validate
}

func NewRequestValidator() *requestValidator {
	v := validator.New()
	_ = v.RegisterValidation("uuid", validateUUID)
	return &requestValidator{
		validator: v,
	}
}

func (v *requestValidator) Validate(i any) error {
	return v.validator.Struct(i)
}
