package services

import (
	"gopkg.in/go-playground/validator.v9"
	"go-api/pkg/app/handler"
	"go-api/pkg/functional"
)

func Validate(v *validator.Validate) func(handler.ErrorHandler) functional.ComposeFunc {
	return func(errorHandler handler.ErrorHandler) functional.ComposeFunc {
		return func(model interface{}) (interface{}, error) {
			if err := v.Struct(model); err != nil {
				return nil, errorHandler([]error{err})
			}

			return model, nil
		}
	}
}
