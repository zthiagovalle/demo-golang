package shared

import "github.com/go-playground/validator/v10"

type CustomValidator struct {
	Validator *validator.Validate
}

func (cv *CustomValidator) Validate(i any) error {
	return cv.Validator.Struct(i)
}

func NewCustomValidator() *CustomValidator {
	return &CustomValidator{Validator: validator.New()}
}
