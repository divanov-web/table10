package taskAcceptPage

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
	"table10/internal/callbackdata"
	"table10/internal/constants/pageCode"
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
			tgbotapi.NewInlineKeyboardButtonData("Назад", pageCode.Tasks),
		),
	)

	return &page{
		AbstractPage: base.AbstractPage{
			Db:           db,
			Logger:       logger,
			Ctx:          ctx,
			User:         user,
			Name:         "Вступление в задание игры",
			Description:  "Вы взяли задание",
			Code:         pageCode.TaskAccept,
			KeyBoard:     &numericKeyboard,
			CallbackData: callbackData,
		},
	}
}

func (p *page) Generate() {
	taskId, err := p.CallbackData.GetTaskId()
	if err != nil {
		p.Logger.Errorf("%v", err)
		return
	}
	p.Description = fmt.Sprintf("Вы взяли задание id\\=%v", taskId)
}
