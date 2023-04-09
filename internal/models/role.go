package models

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	Name        string `gorm:"not null"`
	Code        string `gorm:"unique;not null"`
	Description string
}
