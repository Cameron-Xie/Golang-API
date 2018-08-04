package repository

import (
	"go-api/pkg/app/handler"
	"github.com/jinzhu/gorm"
	"go-api/pkg/functional"
)

func Create(db *gorm.DB) func(handler.ErrorHandler) functional.ComposeFunc {
	return func(errorHandler handler.ErrorHandler) functional.ComposeFunc {
		return func(model interface{}) (interface{}, error) {
			err := db.Create(model).GetErrors()
			if len(err) != 0 {
				return nil, errorHandler(err)
			}

			return model, nil
		}
	}
}
