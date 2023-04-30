package taskAvailablePage

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
	"table10/internal/services"
	"table10/pkg/logging"
)

type page struct {
	base.AbstractPage
}

func NewPage(db *gorm.DB, logger *logging.Logger, ctx context.Context, user *models.User) interfaces.Page {
	numericKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", pageCode.Tasks),
		),
	)

	return &page{
		AbstractPage: base.AbstractPage{
			Db:          db,
			Logger:      logger,
			Ctx:         ctx,
			User:        user,
			Name:        "Доступные задания",
			Description: "",
			Code:        pageCode.TasksAvailable,
			KeyBoard:    &numericKeyboard,
		},
	}
}

func (p *page) Generate() {
	periodRepo := repository.NewPeriodRepository(p.Db)
	periodService := services.NewPeriodService(periodRepo, p.Logger, p.Ctx)
	currentPeriod, err := periodService.ShowCurrent()
	if err != nil {
		p.Logger.Errorf("Текущий период не найден")
	}
	p.Description = fmt.Sprintf("Список доступных заданий на неделе %v", currentPeriod.WeekNumber)
}
