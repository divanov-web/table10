package models

import "gorm.io/gorm"

type Status struct {
	gorm.Model
	Name        string `gorm:"not null"`
	Code        string `gorm:"unique;not null"`
	Description string
}

func (s *Status) GetName() string {
	return s.Name
}
