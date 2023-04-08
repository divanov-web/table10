package tasks

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
	"table10/internal/pages/base"
	"table10/internal/pages/interfaces"
	"table10/pkg/logging"
)

type page struct {
	base.AbstractPage
}

func NewPage(db *gorm.DB, logger *logging.Logger) interfaces.Page {
	numericKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Задание 1", "task1"),
			tgbotapi.NewInlineKeyboardButtonData("Задание 2", "task2"),
			tgbotapi.NewInlineKeyboardButtonData("Задание 3", "task3"),
		),
	)
	return &page{
		AbstractPage: base.AbstractPage{
			Name:        "Задания",
			Description: "Список доступных заданий",
			Command:     "tasks",
			KeyBoard:    &numericKeyboard,
		},
	}
}
