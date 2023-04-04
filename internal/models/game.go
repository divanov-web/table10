package models

import "gorm.io/gorm"

type Game struct {
	gorm.Model
	Name         string `gorm:"not null"`
	ChatId       *int
	LanguageCode *string `gorm:"size:2"`
}
