package models

import "gorm.io/gorm"

type TaskType struct {
	gorm.Model
	Name        string `gorm:"not null"`
	Code        string `gorm:"unique;not null"`
	Description string
}
