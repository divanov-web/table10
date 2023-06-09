package models

import (
	"gorm.io/gorm"
	"table10/pkg/utils/formtating"
	"time"
)

const CloseDateOffset = 7 * 24 * time.Hour

type Task struct {
	gorm.Model
	IsActive         bool      `gorm:"not null:default:false"` //Задание активно
	GameID           uint      `gorm:"not null"`
	Game             Game      `gorm:"foreignKey:GameID"`
	TaskTypeID       uint      //Ссылка на тип задания
	TaskType         TaskType  `gorm:"foreignKey:TaskTypeID"`
	Name             string    `gorm:"not null"`
	StartDate        time.Time `gorm:"not null"` //дата начала задания
	EndDate          time.Time `gorm:"not null"` //Дата окончания принятия задания
	CloseDate        time.Time `gorm:"not null"` //Дата окончания возможности сдать задание
	ShortDescription *string   //Короткое описание задания. Пока отображается только оно в телеграме
	LongDescription  *string   `gorm:"type:text"` //Полное описание, пока не используется
	//Questions        *string    //Вопросы к заданию
	Url       *string    //ссылка на текст задания
	Points    int        `gorm:"type:int;not null;default:1"` //Очки, которые дают за задание (по-умолчанию 1)
	UserTasks []UserTask `gorm:"foreignKey:TaskID"`
}

func (t *Task) GetName() string {
	return t.Name
}

func (t *Task) GetShortDescription() string {
	if t.ShortDescription == nil || *t.ShortDescription == "" {
		return "-"
	}
	return *t.ShortDescription
}

func (t *Task) GetLongDescription() string {
	if t.LongDescription == nil || *t.LongDescription == "" {
		return "-"
	}
	return *t.LongDescription
}

/*func (t *Task) GetQuestions() string {
	if t.Questions == nil || *t.Questions == "" {
		return "-"
	}
	return *t.Questions
}*/

func (t *Task) GetClearedName() string {
	return formtating.UnescapeMarkdownV2(t.Name)
}
