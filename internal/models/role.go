package models

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	Name        string `gorm:"not null"`
	Code        string `gorm:"unique;not null"`
	Description string
}

func (g *Role) CanModerate() bool {
	return g.Code == "moderator" || g.Code == "admin"
}
