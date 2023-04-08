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
			Name:             "Tashkent",
			Code:             "tashkent",
			LanguageCode:     formtating.StrPtr("ru"),
			ShortDescription: formtating.StrPtr("Оффлайновые задания. Город Ташкент"),
			LongDescription:  formtating.StrPtr("Привет! это игра Table10(тестовый сервер). Город Ташкент"),
		},
		{
			Name:             "Online Test",
			Code:             "online_test",
			LanguageCode:     formtating.StrPtr("ru"),
			ShortDescription: formtating.StrPtr("Онлайновые задания для тестовой игры"),
			LongDescription:  formtating.StrPtr("Привет! это игра Table10(тестовый сервер) для онлайн игры"),
		},
		{
			Name:             "Москва",
			Code:             "moscow_test",
			LanguageCode:     formtating.StrPtr("ru"),
			ShortDescription: formtating.StrPtr("Оффлайновые задания. Город Москва"),
			LongDescription:  formtating.StrPtr("Привет! это игра Table10(тестовый сервер). Город Москва"),
		},
		{
			Name:             "Санкт-Петербург",
			Code:             "piter_test",
			LanguageCode:     formtating.StrPtr("ru"),
			ShortDescription: formtating.StrPtr(""),
			LongDescription:  formtating.StrPtr(""),
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
		} else {
			// Запись существует, обновляем ее
			logger.Infof("Игра %s (%s) существует, обновляем", gameToAdd.Name, *gameToAdd.LanguageCode)

			if err = db.Model(&game).Updates(&gameToAdd).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
