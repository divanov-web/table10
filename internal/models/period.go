package models

import (
	"gorm.io/gorm"
	"time"
)

const CloseDateOffset = 7 * 24 * time.Hour

type Period struct {
	gorm.Model
	GameID     uint      `gorm:"not null"`
	Game       Game      `gorm:"foreignKey:GameID"`
	WeekNumber int       `gorm:"type:int2;not null"`
	StartDate  time.Time `gorm:"not null"`
	EndDate    time.Time `gorm:"not null"`
}

// GetStartDate возвращает дату начала периода.
func (p *Period) GetStartDate() time.Time {
	return p.StartDate
}

// GetEndDate возвращает дату окончания периода.
func (p *Period) GetEndDate() time.Time {
	return p.EndDate
}

// GetCloseDate возвращает последнюю дату сдачи задания, которая равна EndDate + 7 дней.
func (p *Period) GetCloseDate() time.Time {
	return p.EndDate.Add(CloseDateOffset)
}
