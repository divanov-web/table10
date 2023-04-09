package gameAcceptPage

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
	"table10/internal/callbackdata"
	"table10/internal/constants"
	"table10/internal/models"
	"table10/internal/pages/base"
	"table10/internal/pages/interfaces"
	"table10/pkg/logging"
)

type page struct {
	base.AbstractPage
}

func NewPage(db *gorm.DB, logger *logging.Logger, ctx context.Context, user *models.User, callbackData *callbackdata.CallbackData) interfaces.Page {
	numericKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", constants.GamePageCode),
		),
	)

	return &page{
		AbstractPage: base.AbstractPage{
			Db:           db,
			Logger:       logger,
			Ctx:          ctx,
			User:         user,
			Name:         "Вступление в сервер игры",
			Description:  "Вы добавлены к серверу",
			Code:         constants.GameAcceptPageCode,
			KeyBoard:     &numericKeyboard,
			CallbackData: callbackData,
		},
	}
}

func (p *page) Generate() {
	gameCode, ok := p.CallbackData.Params["code"]
	if !ok {
		p.Logger.Errorf("Ошибка принятия сервера из-за ошибки передачи параметра")
	}
	p.Description = gameCode
}
