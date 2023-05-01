package models

import (
	"gorm.io/gorm"
	"table10/pkg/utils/formtating"
)

type Game struct {
	gorm.Model
	Name             string `gorm:"not null"`
	Code             string `gorm:"not null"`
	ChatId           *int
	LanguageCode     *string `gorm:"size:2"`
	ShortDescription *string
	LongDescription  *string    `gorm:"type:text"`
	Users            []UserGame `gorm:"foreignKey:GameID"`
}

func (g *Game) GetName() string {
	return formtating.EscapeMarkdownV2(g.Name)
}

func (g *Game) GetShortDescription() string {
	if g.ShortDescription == nil || *g.ShortDescription == "" {
		return "-"
	}
	return formtating.EscapeMarkdownV2(*g.ShortDescription)
}

func (g *Game) GetLongDescription() string {
	if g.LongDescription == nil || *g.LongDescription == "" {
		return "-"
	}
	return formtating.EscapeMarkdownV2(*g.LongDescription)
}
