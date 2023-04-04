package mainpage

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
			tgbotapi.NewInlineKeyboardButtonData("Задания", "tasks"),
			tgbotapi.NewInlineKeyboardButtonData("Личный кабинет", "cabinet"),
		),
	)
	return &page{
		BasePage: pages.BasePage{
			Name:        "Главное меню",
			Description: "Доступные пункты меню",
			Command:     "main",
			KeyBoard:    &numericKeyboard,
		},
	}
}
