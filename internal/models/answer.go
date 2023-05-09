package models

import (
	"gorm.io/gorm"
)

type Answer struct {
	gorm.Model

	TaskID         uint   `gorm:"not null" json:"task_id"`
	UserID         uint   `gorm:"not null" json:"user_id"`
	Text           string `gorm:"size:255" json:"text"`
	ImagePath      string `gorm:"size:255" json:"image_path"`
	TelegramFileId string `gorm:"not null"` //ID файла в телеграмм

	Task Task `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	User User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
}
