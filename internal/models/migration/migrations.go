package migration

import (
	"context"
	"gorm.io/gorm"
	"table10/internal/config"
	models2 "table10/internal/models"
	"table10/internal/models/seed"
	"table10/pkg/logging"
	"time"
)

func RunMigrations(cfg *config.Config, db *gorm.DB, logger *logging.Logger) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	logger.Info("Запуск базовых миграциий")

	// Здесь все модели, которые требуют миграции
	models := []interface{}{
		&models2.User{},
		&models2.Game{},
		&models2.UserGame{},
		&models2.Role{},
		&models2.Period{},
		&models2.TaskType{},
		&models2.Task{},
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
	//типы игр
	err = seed.AddTaskType(db, logger)
	if err != nil {
		return err
	}

	if cfg.IsProd != nil && *cfg.IsProd {
		err = seed.AddTask(db, logger, ctx)
		if err != nil {
			return err
		}
	}

	logger.Info("Все базовые миграции выполнены успешно")
	return nil
}
