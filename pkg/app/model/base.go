package model

import "time"

type Base struct {
	ID        uint `gorm:"primary_key;auto_increment"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
