package migrations

import (
	"gorm.io/gorm"
	"table10/internal/models"
	"table10/pkg/logging"
)

func init() {
	migrationWithLogger = append(migrationWithLogger, &MigrationWithLogger{
		ID:          "0001_create_tasks",
		Description: "Create tasks table",
		Migrate: func(tx *gorm.DB, logger *logging.Logger) error {
			logger.Infof("Разовая миграция выполнена")
			return tx.AutoMigrate(&models.Task{})
		},
		Rollback: func(tx *gorm.DB, logger *logging.Logger) error {
			return tx.Migrator().DropTable("tasks")
		},
	})
}
