package cabinet

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
			tgbotapi.NewInlineKeyboardButtonData("Счёт", "score"),
			tgbotapi.NewInlineKeyboardButtonData("Настройки", "config"),
		),
	)
	return &page{
		AbstractPage: base.AbstractPage{
			Name:        "Личный кабинет",
			Description: "Управление личным кабинетом",
			Command:     "cabinet",
			KeyBoard:    &numericKeyboard,
		},
	}
}
