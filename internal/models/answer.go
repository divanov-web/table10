package models

import (
	"gorm.io/gorm"
)

type Answer struct {
	gorm.Model

	UserTaskID     uint   `gorm:"not null" json:"user_task_id"` //На какое задание ответ
	UserID         uint   `gorm:"not null" json:"user_id"`      //кто отвечает (пользователь или админ)
	Text           string `gorm:"size:255" json:"text"`
	ImagePath      string `gorm:"size:255" json:"image_path"`
	TelegramFileId string `gorm:"not null"` //ID файла в телеграмм

	UserTask UserTask `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	User     User     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
}
