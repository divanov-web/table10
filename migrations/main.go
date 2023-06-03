package migrations

import (
	"context"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
	"table10/internal/config"
	"table10/internal/models"
	seed2 "table10/migrations/seed"
	"table10/pkg/logging"
	"time"
)

type MigrationWithLogger struct {
	ID          string
	Description string
	Migrate     func(*gorm.DB, *logging.Logger) error
	Rollback    func(*gorm.DB, *logging.Logger) error
}

var migrationWithLogger []*MigrationWithLogger

// CustomMigrationHistory модель для пользовательской таблицы миграций
type CustomMigrationHistory struct {
	ID          string `gorm:"primaryKey"`
	Description string `gorm:"type:text"`
	AppliedAt   *time.Time
}

var Migrations []*gormigrate.Migration

func Migrate(cfg *config.Config, db *gorm.DB, logger *logging.Logger) error {
	gormMigrations := make([]*gormigrate.Migration, len(migrationWithLogger))

	// Создание таблицы CustomMigrationHistory
	err := db.AutoMigrate()
	if err != nil {
		return err
	}

	m := gormigrate.New(db, gormigrate.DefaultOptions, gormMigrations)

	for i, migration := range migrationWithLogger {
		gormMigrations[i] = &gormigrate.Migration{
			ID: migration.ID,
			Migrate: func(tx *gorm.DB) error {
				return migration.Migrate(tx, logger)
			},
			Rollback: func(tx *gorm.DB) error {
				return migration.Rollback(tx, logger)
			},
		}
	}

	err = m.Migrate()
	if err != nil {
		return err
	}

	return runAlwaysMigrations(cfg, db, logger)
}

func runAlwaysMigrations(cfg *config.Config, db *gorm.DB, logger *logging.Logger) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	models := []interface{}{
		&models.User{},
		&models.Game{},
		&models.UserGame{},
		&models.Role{},
		&models.Period{},
		&models.TaskType{},
		//&models.Task{},
		&models.Status{},
		&models.UserTask{},
		&models.Answer{},
	}

	for _, model := range models {
		err := db.AutoMigrate(model)
		if err != nil {
			return err
		}
	}
	logger.Info("Все регулярные миграции выполнены успешно")

	//миграция уровней доступа пользователей к играм
	err := seed2.AddRole(db, logger)
	if err != nil {
		return err
	}
	//миграция статусов выполнения задач
	err = seed2.AddStatus(db, logger)
	if err != nil {
		return err
	}
	//миграция предзаготовленных игр
	err = seed2.AddGames(db, logger)
	if err != nil {
		return err
	}
	//обновление таблицы периодов
	err = seed2.AddPeriods(db, logger)
	if err != nil {
		return err
	}
	//типы игр
	err = seed2.AddTaskType(db, logger)
	if err != nil {
		return err
	}

	if cfg.IsProd != nil && *cfg.IsProd == false {
		err = seed2.AddTask(db, logger, ctx)
		if err != nil {
			return err
		}
	}

	return nil
}
