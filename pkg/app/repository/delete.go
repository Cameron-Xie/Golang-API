package repository

import (
	"github.com/jinzhu/gorm"
	"go-api/pkg/app/handler"
	"go-api/pkg/functional"
)

func Delete(db *gorm.DB) func(handler.ErrorHandler) functional.ComposeFunc {
	return func(errorHandler handler.ErrorHandler) functional.ComposeFunc {
		return func(model interface{}) (interface{}, error) {
			err := db.Delete(model).GetErrors()
			if len(err) != 0 {
				return nil, errorHandler(err)
			}

			return model, nil
		}
	}
}
