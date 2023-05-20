package models

import (
	"gorm.io/gorm"
	"table10/pkg/utils/formtating"
)

type User struct {
	gorm.Model
	TelegramID   int    `gorm:"unique"`
	IsBot        bool   `gorm:"not null"`
	FirstName    string `gorm:"not null"`
	Username     string
	LanguageCode string     `gorm:"size:3"`
	LastPage     string     `gorm:"size:64"`
	Games        []UserGame `gorm:"foreignKey:UserID"`
	Tasks        []UserTask `gorm:"foreignKey:UserID"`
}

func CreateUser(db *gorm.DB, user *User) error {
	return db.Create(user).Error
}

func (u *User) GetUserName() string {
	return formtating.EscapeMarkdownV2(u.Username)
}
