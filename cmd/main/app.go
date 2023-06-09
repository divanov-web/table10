package main

import (
	"table10/internal/config"
	"table10/migrations"
	"table10/pkg/logging"
)

func main() {
	logger := logging.GetLogger()
	logger.Info("create telegram connection")

	cfg := config.GetConfig()

	// DB connection
	db, err := ConnectToDatabase(cfg, logger)
	if err != nil {
		logger.Fatalf("Failed to connect to the database: %v", err)
	}
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	// GORMigrate
	if err := migrations.Migrate(cfg, db, logger); err != nil {
		logger.Fatalf("Failed to apply migrations: %v", err)
	}

	telegramStart(cfg, logger, db)

}
