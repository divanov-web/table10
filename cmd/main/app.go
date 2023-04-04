package main

import (
	"table10/internal/config"
	"table10/internal/models"
	"table10/pkg/logging"
)

func main() {
	logger := logging.GetLogger()
	logger.Info("create telegram connection")

	cfg := config.GetConfig()

	// Подключение к БД
	db, err := ConnectToDatabase(cfg, logger)
	if err != nil {
		logger.Fatalf("Failed to connect to the database: %v", err)
	}
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	//обновление таблицы периодов
	err = models.SeedPeriods(db, logger)
	if err != nil {
		logger.Fatalf("Failed to run period update: %v", err)
	}

	// миграции базы данных
	err = models.RunMigrations(db, logger)
	if err != nil {
		logger.Fatalf("Failed to run database migrations: %v", err)
	}

	telegramStart(cfg, logger, db)

}
