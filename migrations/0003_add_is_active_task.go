package migrations

import (
	"gorm.io/gorm"
	"table10/internal/models"
	"table10/pkg/logging"
)

func init() {
	migrationWithLogger = append(migrationWithLogger, &MigrationWithLogger{
		ID:          "0003_add_is_active_task",
		Description: "Add IsActive to Task model",
		Migrate: func(tx *gorm.DB, logger *logging.Logger) error {
			// Убеждаемся, что значение по умолчанию задано для всех существующих записей
			tx.Model(&models.Task{}).Update("is_active", true)

			// Изменение структуры таблицы
			return tx.Exec("ALTER TABLE tasks ADD COLUMN is_active boolean NOT NULL DEFAULT false").Error
		},
		Rollback: func(tx *gorm.DB, logger *logging.Logger) error {
			// удалить столбец
			return tx.Exec("ALTER TABLE tasks DROP COLUMN is_active").Error
		},
	})
}
