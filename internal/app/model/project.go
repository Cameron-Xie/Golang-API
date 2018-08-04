package model

import (
	"go-api/pkg/app/model"
	_ "github.com/jinzhu/gorm"
)

type Project struct {
	model.Base
	Name        string `gorm:"type:varchar(100);unique_index" json:"name" validate:"required,min=3,max=120"`
	Description string `gorm:"type:varchar(100)" json:"description"`
}
