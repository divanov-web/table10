package models

import (
	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	// Здесь все модели, которые требуют миграции
	err := db.AutoMigrate(&User{})
	if err != nil {
		return err
	}
	return nil
}
