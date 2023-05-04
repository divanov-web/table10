package models

import "gorm.io/gorm"

type UserGame struct {
	gorm.Model
	UserID     uint `gorm:"uniqueIndex:user_game_idx"`
	User       User
	GameID     uint `gorm:"uniqueIndex:user_game_idx"`
	Game       Game
	RoleID     uint
	Role       Role `gorm:"foreignKey:RoleID"`
	IsMain     bool `gorm:"not null"`
	IsBlocked  bool `gorm:"not null:default:false"`
	IsAccepted bool `gorm:"not null:default:true"`
}
