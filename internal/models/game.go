package models

import (
	"gorm.io/gorm"
)

type Game struct {
	gorm.Model
	Name             string  `gorm:"not null"`
	Code             string  `gorm:"not null"`
	ChatId           *int    //id чата в телеграмм, где будут объявления
	LanguageCode     *string `gorm:"size:2"`
	ShortDescription *string
	LongDescription  *string    `gorm:"type:text"`
	Users            []UserGame `gorm:"foreignKey:GameID"`
}

func (g *Game) GetName() string {
	return g.Name
}

func (g *Game) GetShortDescription() string {
	if g.ShortDescription == nil || *g.ShortDescription == "" {
		return "-"
	}
	return *g.ShortDescription
}

func (g *Game) GetLongDescription() string {
	if g.LongDescription == nil || *g.LongDescription == "" {
		return "-"
	}
	return *g.LongDescription
}
