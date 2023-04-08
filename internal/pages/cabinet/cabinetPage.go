package cabinetPage

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
	"table10/internal/models"
	"table10/internal/pages/base"
	gamePage "table10/internal/pages/game"
	"table10/internal/pages/interfaces"
	"table10/pkg/logging"
)

const Command = "cabinet"

type page struct {
	base.AbstractPage
}

func NewPage(db *gorm.DB, logger *logging.Logger, ctx context.Context, user *models.User) interfaces.Page {
	numericKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Счёт", "score"),
			tgbotapi.NewInlineKeyboardButtonData("Сервера игры", gamePage.Command),
		),
	)
	return &page{
		AbstractPage: base.AbstractPage{
			Name:        "Личный кабинет",
			Description: "Управление личным кабинетом",
			Command:     Command,
			KeyBoard:    &numericKeyboard,
		},
	}
}
