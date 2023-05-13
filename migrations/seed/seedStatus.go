package seed

import (
	"gorm.io/gorm"
	"table10/internal/models"
	"table10/pkg/logging"
)

func AddStatus(db *gorm.DB, logger *logging.Logger) error {
	itemsToAdd := []models.Status{
		{
			Name: "Новая",
			Code: "new",
		},
		{
			Name: "Поиск участников",
			Code: "search",
		},
		{
			Name: "Выполняется",
			Code: "in_progress",
		},
		{
			Name: "На проверке",
			Code: "under_review",
		},
		{
			Name: "Принята",
			Code: "accepted",
		},
	}

	for _, itemToAdd := range itemsToAdd {
		var item models.Status
		if err := db.Where("code = ?", itemToAdd.Code).First(&item).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				// Если запись не существует, создаем новую запись и сохраняем ее в таблице
				logger.Infof("Статус %s не существует, добавляем", itemToAdd.Name)

				if err = db.Create(&itemToAdd).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		} else {
			if err = db.Model(&item).Updates(&itemToAdd).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
