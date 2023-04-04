package models

import (
	"gorm.io/gorm"
	"table10/pkg/logging"
)

func seedGames(db *gorm.DB, logger *logging.Logger) error {
	logger.Info("Добавление записей в таблицу Game")
	// Проверяем, существует ли уже запись с заданными значениями в таблице Game
	var game Game
	languageCode := "ru"
	if err := db.Where("name = ? AND language_code = ?", "Tashkent", languageCode).First(&game).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Если запись не существует, создаем новую запись и сохраняем ее в таблице
			logger.Info("Игры Tashkent не сущестует, добавляем")
			newGame := Game{
				Name:         "Tashkent",
				LanguageCode: &languageCode,
			}

			if err = db.Create(&newGame).Error; err != nil {
				return err
			}
		} else {
			return err
		}
	}

	return nil
}
