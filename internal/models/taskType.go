package models

import (
	"gorm.io/gorm"
	"table10/pkg/utils/formtating"
)

type TaskType struct {
	gorm.Model
	Name        string `gorm:"not null"`
	Code        string `gorm:"unique;not null"`
	Description string
}

func (t *TaskType) GetName() string {
	return formtating.EscapeMarkdownV2(t.Name)
}
