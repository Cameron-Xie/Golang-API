package repository

import (
	"github.com/jinzhu/gorm"
	"go-api/pkg/app/handler"
	"go-api/pkg/functional"
)

func Update(db *gorm.DB, model interface{}) func(handler.ErrorHandler) functional.ComposeFunc {
	return func(errorHandler handler.ErrorHandler) functional.ComposeFunc {
		return func(updates interface{}) (interface{}, error) {
			err := db.Model(model).Updates(updates).GetErrors()
			if len(err) != 0 {
				return nil, errorHandler(err)
			}

			return model, nil
		}
	}
}
