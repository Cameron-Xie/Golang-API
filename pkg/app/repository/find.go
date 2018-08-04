package repository

import (
	"github.com/jinzhu/gorm"
	"go-api/pkg/app/handler"
)

func Find(db *gorm.DB) func(handler.ErrorHandler) func(interface{}) (interface{}, error) {
	return func(errorHandler handler.ErrorHandler) func(interface{}) (interface{}, error) {
		return func(project interface{}) (interface{}, error) {
			err := db.Find(project).GetErrors()
			if len(err) != 0 {
				return nil, errorHandler(err)
			}

			return project, nil
		}
	}
}
