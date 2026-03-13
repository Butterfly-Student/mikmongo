package validator

import (
	"github.com/go-playground/validator/v10"
)

// Validator wraps validator.Validate
type Validator struct {
	validate *validator.Validate
}

// New creates a new validator
func New() *Validator {
	return &Validator{
		validate: validator.New(),
	}
}

// Validate validates a struct
func (v *Validator) Validate(i interface{}) error {
	return v.validate.Struct(i)
}

// ValidateVar validates a single variable
func (v *Validator) ValidateVar(field interface{}, tag string) error {
	return v.validate.Var(field, tag)
}

// RegisterCustom registers a custom validation
func (v *Validator) RegisterCustom(tag string, fn validator.Func) error {
	return v.validate.RegisterValidation(tag, fn)
}
