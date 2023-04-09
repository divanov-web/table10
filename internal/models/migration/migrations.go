package migration

import (
	"gorm.io/gorm"
	models2 "table10/internal/models"
	"table10/internal/models/seed"
	"table10/pkg/logging"
)

func RunMigrations(db *gorm.DB, logger *logging.Logger) error {
	logger.Info("Запуск базовых миграциий")

	// Здесь все модели, которые требуют миграции
	models := []interface{}{
		&models2.User{},
		&models2.Game{},
		&models2.UserGame{},
		&models2.Role{},
		&models2.Period{},
	}

	for _, model := range models {
		err := db.AutoMigrate(model)
		if err != nil {
			return err
		}
	}

	//миграция уровней доступа пользователей к играм
	err := seed.AddRole(db, logger)
	if err != nil {
		return err
	}
	//миграция предзаготовленных игр
	err = seed.AddGames(db, logger)
	if err != nil {
		return err
	}
	//обновление таблицы периодов
	err = seed.AddPeriods(db, logger)
	if err != nil {
		return err
	}

	logger.Info("Все базовые миграции выполнены успешно")
	return nil
}
