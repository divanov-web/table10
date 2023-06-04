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

	canModerate := user.Games[0].Role.CanModerate()

	rows := [][]tgbotapi.InlineKeyboardButton{
		{
			tgbotapi.NewInlineKeyboardButtonData("Мои задания", pageCode.TasksAccepted),
			tgbotapi.NewInlineKeyboardButtonData("Доступные задания", pageCode.Tasks),
		},
		{
			//tgbotapi.NewInlineKeyboardButtonData("Личный кабинет", pageCode.Cabinet),
			tgbotapi.NewInlineKeyboardButtonData("Об игре", pageCode.Welcome),
		},
	}

	if canModerate {
		extraRow := []tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData("Администрирование", pageCode.Admin),
		}
		rows = append(rows, extraRow)
	}

	numericKeyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

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
