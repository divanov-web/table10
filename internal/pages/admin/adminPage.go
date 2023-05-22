package adminPage

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
	"table10/internal/constants/pageCode"
	"table10/internal/models"
	"table10/internal/pages/base"
	"table10/internal/pages/interfaces"
	"table10/pkg/logging"
)

type page struct {
	base.AbstractPage
}

func NewPage(db *gorm.DB, logger *logging.Logger, ctx context.Context, user *models.User) interfaces.Page {
	numericKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Меню", pageCode.Main),
		),
	)

	return &page{
		AbstractPage: base.AbstractPage{
			Db:          db,
			Logger:      logger,
			Ctx:         ctx,
			User:        user,
			Name:        "Админка",
			Description: "У вас нет доступа к этой странице",
			Code:        pageCode.Admin,
			KeyBoard:    &numericKeyboard,
		},
	}
}

func (p *page) Generate() {
	canModerate := p.CanModerate()
	if !canModerate {
		p.Description = "У вас нет доступа к этой странице"
		p.Logger.Errorf("Getting access to admin interface. user_id=%v", p.User.ID)
		return
	}

	p.Description = "Административный интерфейс"
	numericKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Задания на проверку", pageCode.AdminReview),
			tgbotapi.NewInlineKeyboardButtonData("Меню", pageCode.Main),
		),
	)
	p.KeyBoard = &numericKeyboard

}
