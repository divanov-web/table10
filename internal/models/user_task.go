package models

import "gorm.io/gorm"

type UserTask struct {
	gorm.Model
	UserID   uint `gorm:"uniqueIndex:user_task_idx"`
	User     User
	TaskID   uint `gorm:"uniqueIndex:user_task_idx"`
	Task     Task
	StatusID uint
	Status   Role `gorm:"foreignKey:StatusID"`
}
