package mainPage

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
	"table10/internal/constants/pageCode"
	"table10/internal/models"
	"table10/internal/pages/base"
	"table10/internal/pages/interfaces"
	"table10/internal/repository"
	"table10/pkg/logging"
)

type page struct {
	base.AbstractPage
}

func NewPage(db *gorm.DB, logger *logging.Logger, ctx context.Context, user *models.User) interfaces.Page {
	numericKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Задания", pageCode.Tasks),
			tgbotapi.NewInlineKeyboardButtonData("Личный кабинет", pageCode.Cabinet),
			tgbotapi.NewInlineKeyboardButtonData("Об игре", pageCode.Welcome),
		),
	)

	var text string
	periodRepo := repository.NewPeriodRepository(db)
	currentPeriod, err := periodRepo.ShowCurrent(ctx)
	if err != nil {
		logger.Errorf("Текущая неделя не найдена")
	} else {
		text = fmt.Sprintf("Сейчас идёт %v неделя игры", currentPeriod.WeekNumber)
	}

	return &page{
		AbstractPage: base.AbstractPage{
			Name:        "Главное меню",
			Description: text,
			Code:        pageCode.Main,
			KeyBoard:    &numericKeyboard,
		},
	}
}
