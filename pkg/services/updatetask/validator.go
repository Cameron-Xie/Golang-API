package updatetask

import (
	"errors"

	v "github.com/Cameron-Xie/Golang-API/pkg/validator"
)

type validator struct {
	v     v.Validator
	rules map[string]v.MapRule
}

func (o *validator) Validate(i map[string]interface{}) (map[string]interface{}, error) {
	if len(i) == 0 {
		return nil, errors.New("empty input")
	}

	return o.v.ValidateMap(i, o.rules)
}

func NewValidator() Validator {
	val, _ := v.New(nil)

	return &validator{
		v: val,
		rules: map[string]v.MapRule{
			"name": {
				Rule: "min=1,max=100",
				Msg:  "name must be at least 1 character and less than 100 characters",
			},
			"description": {
				Rule: "max=200",
				Msg:  "description must not be longer than 200 characters",
			},
		},
	}
}
