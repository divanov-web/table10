package main

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"table10/internal/config"
	"table10/pkg/logging"
)

func ConnectToDatabase(cfg *config.Config, logger *logging.Logger) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=UTC",
		cfg.Storage.Host,
		cfg.Storage.Port,
		cfg.Storage.Username,
		cfg.Storage.Password,
		cfg.Storage.Database,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
