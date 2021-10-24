package validator

import (
	"github.com/go-playground/validator/v10"
)

type Validator interface {
	ValidateStruct(s interface{}) error
}

type val struct {
	*validator.Validate
}

func (v *val) ValidateStruct(s interface{}) error {
	return v.Struct(s)
}

func NewValidator() Validator {
	return &val{validator.New()}
}
