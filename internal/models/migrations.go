package models

import (
	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	// Здесь все модели, которые требуют миграции
	models := []interface{}{
		&User{},
		&Game{},
	}

	for _, model := range models {
		err := db.AutoMigrate(model)
		if err != nil {
			return err
		}
	}

	err := seedGames(db)
	if err != nil {
		return err
	}

	return nil
}
