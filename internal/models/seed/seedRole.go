package seed

import (
	"gorm.io/gorm"
	"table10/internal/models"
	"table10/pkg/logging"
)

func AddRole(db *gorm.DB, logger *logging.Logger) error {
	// Создаем список ролей
	rolesToAdd := []models.Role{
		{
			Name: "Администратор",
			Code: "admin",
		},
		{
			Name: "Модератор",
			Code: "moderator",
		},
		{
			Name: "Пользователь",
			Code: "user",
		},
	}

	for _, itemToAdd := range rolesToAdd {
		var role models.Role
		if err := db.Where("code = ?", itemToAdd.Code).First(&role).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				// Если запись не существует, создаем новую запись и сохраняем ее в таблице
				logger.Infof("Роль %s не существует, добавляем", itemToAdd.Name)

				if err = db.Create(&itemToAdd).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		} else {
			if err = db.Model(&role).Updates(&itemToAdd).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
