package gamePage

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
	"table10/internal/models"
	"table10/internal/pages/base"
	gamePageInput "table10/internal/pages/gameInput"
	"table10/internal/pages/interfaces"
	"table10/pkg/logging"
)

const Command = "game"

type page struct {
	base.AbstractPage
}

func NewPage(db *gorm.DB, logger *logging.Logger, ctx context.Context, user *models.User) interfaces.Page {
	numericKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Поиск по коду", gamePageInput.Command),
		),
	)
	return &page{
		AbstractPage: base.AbstractPage{
			Name:        "Список серверов",
			Description: "Просмотр список серверов, в которых вы участвуете",
			Command:     Command,
			KeyBoard:    &numericKeyboard,
		},
	}
}
