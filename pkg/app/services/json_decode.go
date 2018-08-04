package services

import (
	"io"
	"encoding/json"
	"go-api/pkg/app/handler"
	"go-api/pkg/functional"
)

func Decode(reader io.Reader) func(handler.ErrorHandler) functional.ComposeFunc {
	return func(errorHandler handler.ErrorHandler) functional.ComposeFunc {
		return func(model interface{}) (interface{}, error) {
			if err := json.NewDecoder(reader).Decode(model); err != nil {
				return nil, errorHandler([]error{err})
			}

			return model, nil
		}
	}
}
