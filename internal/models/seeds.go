package models

import (
	"gorm.io/gorm"
)

func seedGames(db *gorm.DB) error {
	// Проверяем, существует ли уже запись с заданными значениями в таблице Game
	var game Game
	languageCode := "ru"
	if err := db.Where("name = ? AND language_code = ?", "Tashkent", languageCode).First(&game).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Если запись не существует, создаем новую запись и сохраняем ее в таблице
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
