package models

import (
	"gorm.io/gorm"
	"table10/pkg/logging"
	"table10/pkg/utils/formtating"
)

func seedGames(db *gorm.DB, logger *logging.Logger) error {
	logger.Info("Добавление записей в таблицу Game")

	// Создаем список игр для добавления
	gamesToAdd := []Game{
		{
			Name:         "Tashkent",
			LanguageCode: formtating.StrPtr("ru"),
		},
		{
			Name:         "Online Test",
			LanguageCode: formtating.StrPtr("ru"),
		},
	}

	for _, gameToAdd := range gamesToAdd {
		var game Game
		if err := db.Where("name = ? AND language_code = ?", gameToAdd.Name, gameToAdd.LanguageCode).First(&game).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				// Если запись не существует, создаем новую запись и сохраняем ее в таблице
				logger.Infof("Игра %s (%s) не существует, добавляем", gameToAdd.Name, *gameToAdd.LanguageCode)

				if err = db.Create(&gameToAdd).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}

	return nil
}
