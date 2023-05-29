package httpServer

import "gopkg.in/go-playground/validator.v9"

// Validator ...
type Validator struct {
	validator *validator.Validate
}

// Validate validates structs
func (cv *Validator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}
