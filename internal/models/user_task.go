package models

import "gorm.io/gorm"

type UserTask struct {
	gorm.Model
	UserID   uint `gorm:"index:idx_user_task,unique"`
	User     User
	TaskID   uint `gorm:"index:idx_user_task,unique"`
	Task     Task
	StatusID uint
	Status   Status   `gorm:"foreignKey:StatusID"`
	Answers  []Answer `gorm:"foreignKey:UserTaskID" json:"-"`
}
