package seed

import (
	"gorm.io/gorm"
	"table10/internal/models"
	"table10/pkg/logging"
)

func AddTaskType(db *gorm.DB, logger *logging.Logger) error {
	// Создаем список ролей
	itemsToAdd := []models.TaskType{
		{
			Name: "Свободные команды",
			Code: "common",
		},
		{
			Name: "Бади",
			Code: "buddy",
		},
		{
			Name: "Одиночное",
			Code: "solo",
		},
	}

	for _, itemToAdd := range itemsToAdd {
		var item models.TaskType
		if err := db.Where("code = ?", itemToAdd.Code).First(&item).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				// Если запись не существует, создаем новую запись и сохраняем ее в таблице
				logger.Infof("Тип задания %s не существует, добавляем", itemToAdd.Name)

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
