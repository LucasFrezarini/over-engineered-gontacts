package validator

import (
	v "github.com/go-playground/validator/v10"
)

// CustomValidator defines a custom validator to use to validate incoming requests
// in the Server's Echo instance
type CustomValidator struct {
	Validator *v.Validate
}

// Validate uses the go-playground/validator to validate the request body passed as parameter
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.Validator.Struct(i)
}

// NewCustomValidator returns a ready to use CustomValidator to integrate with Echo
func NewCustomValidator() *CustomValidator {
	return &CustomValidator{Validator: v.New()}
}
