package models

import (
	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	PeriodID         uint   `gorm:"not null"`
	Period           Period `gorm:"foreignKey:PeriodID"`
	TaskTypeID       uint
	TaskType         TaskType `gorm:"foreignKey:TaskTypeID"`
	Name             string   `gorm:"not null"`
	ShortDescription *string
	LongDescription  *string `gorm:"type:text"`
	Url              *string //ссылка на текст задания
	Points           int     `gorm:"type:int;not null;default:1"` //Очки, которые дают за задание (по-умолчанию 1)
}
