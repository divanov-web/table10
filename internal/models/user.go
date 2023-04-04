package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	TelegramID   int    `gorm:"unique"`
	IsBot        bool   `gorm:"not null"`
	FirstName    string `gorm:"not null"`
	Username     string
	LanguageCode string `gorm:"size:3"`
	LastPage     string `gorm:"size:64"`
}

func CreateUser(db *gorm.DB, user *User) error {
	return db.Create(user).Error
}
