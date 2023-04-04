package cabinet

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
			tgbotapi.NewInlineKeyboardButtonData("Счёт", "score"),
			tgbotapi.NewInlineKeyboardButtonData("Настройки", "config"),
		),
	)
	return &page{
		BasePage: pages.BasePage{
			Name:        "Личный кабинет",
			Description: "Управление личным кабинетом",
			Command:     "cabinet",
			KeyBoard:    &numericKeyboard,
		},
	}
}
