package models

import (
	"gorm.io/gorm"
	"table10/pkg/utils/formtating"
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

func (t *Task) GetName() string {
	return formtating.EscapeMarkdownV2(t.Name)
}

func (t *Task) GetShortDescription() string {
	if t.ShortDescription == nil || *t.ShortDescription == "" {
		return "-"
	}
	return formtating.EscapeMarkdownV2(*t.ShortDescription)
}

func (t *Task) GetLongDescription() string {
	if t.LongDescription == nil || *t.LongDescription == "" {
		return "-"
	}
	return formtating.EscapeMarkdownV2(*t.LongDescription)
}
