package tasks

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"table10/internal/pages"
)

type page struct {
	pages.BasePage
}

func NewPage() pages.Page {
	numericKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Задание 1", "task1"),
			tgbotapi.NewInlineKeyboardButtonData("Задание 2", "task2"),
			tgbotapi.NewInlineKeyboardButtonData("Задание 3", "task3"),
		),
	)
	return &page{
		BasePage: pages.BasePage{
			Name:        "Задания",
			Description: "Список доступных заданий",
			Command:     "tasks",
			KeyBoard:    &numericKeyboard,
		},
	}
}
