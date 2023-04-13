package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	TelegramID   int    `gorm:"unique"`
	IsBot        bool   `gorm:"not null"`
	FirstName    string `gorm:"not null"`
	Username     string
	LanguageCode string     `gorm:"size:3"`
	LastPage     string     `gorm:"size:64"`
	Games        []UserGame `gorm:"foreignKey:UserID"`
}

type UserGame struct {
	gorm.Model
	UserID uint `gorm:"uniqueIndex:user_game_idx"`
	User   User
	GameID uint `gorm:"uniqueIndex:user_game_idx"`
	Game   Game
	RoleID uint
	Role   Role `gorm:"foreignKey:RoleID"`
	IsMain bool `gorm:"not null"`
}

func CreateUser(db *gorm.DB, user *User) error {
	return db.Create(user).Error
}
