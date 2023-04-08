package models

import (
	"gorm.io/gorm"
	"table10/pkg/logging"
)

func RunMigrations(db *gorm.DB, logger *logging.Logger) error {
	logger.Info("Запуск базовых миграциий")
	// Здесь все модели, которые требуют миграции
	models := []interface{}{
		&User{},
		&Game{},
		&Period{},
	}

	for _, model := range models {
		err := db.AutoMigrate(model)
		if err != nil {
			return err
		}
	}

	err := seedGames(db, logger)
	if err != nil {
		return err
	}
	//обновление таблицы периодов
	err = SeedPeriods(db, logger)
	if err != nil {
		return err
	}

	logger.Info("Все базовые миграции выполнены успешно")
	return nil
}
