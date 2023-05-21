package migrations

import (
	"gorm.io/gorm"
	"table10/internal/models"
	"table10/pkg/logging"
)

func init() {
	migrationWithLogger = append(migrationWithLogger, &MigrationWithLogger{
		ID:          "0002_add_status_rejected",
		Description: "Add status rejected",
		Migrate: func(tx *gorm.DB, logger *logging.Logger) error {
			logger.Infof("adding status with code rejected")
			itemsToAdd := []models.Status{
				{
					Name: "Отклонена",
					Code: "rejected",
				},
			}

			for _, itemToAdd := range itemsToAdd {
				var item models.Status
				if err := tx.Where("code = ?", itemToAdd.Code).First(&item).Error; err != nil {
					if err == gorm.ErrRecordNotFound {
						if err = tx.Create(&itemToAdd).Error; err != nil {
							return err
						}
					} else {
						return err
					}
				} else {
					if err = tx.Model(&item).Updates(&itemToAdd).Error; err != nil {
						return err
					}
				}
			}

			return nil
		},
		Rollback: func(tx *gorm.DB, logger *logging.Logger) error {
			return nil
		},
	})
}
