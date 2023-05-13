package seed

import (
	"gorm.io/gorm"
	"table10/internal/models"
	"table10/pkg/logging"
	"table10/pkg/utils/formtating"
)

func AddGames(db *gorm.DB, logger *logging.Logger) error {
	// Создаем список игр для добавления
	gamesToAdd := []models.Game{
		{
			Name:             "Tashkent",
			Code:             "tashkent",
			LanguageCode:     formtating.StrPtr("ru"),
			ShortDescription: formtating.StrPtr(formtating.EscapeMarkdownV2("Оффлайновые задания. Город Ташкент")),
			LongDescription: formtating.StrPtr(formtating.EscapeMarkdownV2("Эта версия игры Table 10 предназначена в основном для релокантов из России, " +
				"находящихся в Ташкенте, хотя местным жителям также может быть интересно присоединиться, " +
				"в игре нет ограничений по национальности. Многие из нас планируют временное пребывание в Ташкенте, " +
				"и цель игры заключается в том, чтобы провести это время с максимальной пользой: " +
				"посетить как можно больше достопримечательностей, исследовать другие города и районы, попробовать новые блюда и познакомиться с узбекской культурой.")),
		},
		{
			Name:             "Test",
			Code:             "test",
			LanguageCode:     formtating.StrPtr("en"),
			ShortDescription: formtating.StrPtr(formtating.EscapeMarkdownV2("Test server")),
			LongDescription:  formtating.StrPtr(formtating.EscapeMarkdownV2("Welcome to test server! This server is specifically designed for testing purposes. Here, you'll have the opportunity to evaluate and verify the functionality of my bot.")),
		},
	}

	for _, gameToAdd := range gamesToAdd {
		var game models.Game
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
			if err = db.Model(&game).Updates(&gameToAdd).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
