package storetask

import (
	v "github.com/Cameron-Xie/Golang-API/pkg/validator"
)

type validator struct {
	v v.Validator
}

func (o *validator) Validate(i *Task) error {
	return o.v.ValidateStruct(i)
}

func NewValidator() Validator {
	val, _ := v.New(nil)

	return &validator{v: val}
}
