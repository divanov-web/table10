package taskDetailPage

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
	"strconv"
	"table10/internal/callbackdata"
	"table10/internal/constants"
	"table10/internal/constants/pageCode"
	"table10/internal/models"
	"table10/internal/pages/base"
	"table10/internal/pages/interfaces"
	"table10/internal/repository"
	"table10/internal/services"
	"table10/pkg/logging"
	"table10/pkg/utils"
)

type page struct {
	base.AbstractPage
}

func NewPage(db *gorm.DB, logger *logging.Logger, ctx context.Context, user *models.User, callbackData *callbackdata.CallbackData) interfaces.Page {

	return &page{
		AbstractPage: base.AbstractPage{
			Db:           db,
			Logger:       logger,
			Ctx:          ctx,
			User:         user,
			Name:         "Детальная задания",
			Description:  "",
			Code:         pageCode.TaskDetail,
			KeyBoard:     nil,
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
	taskRepo := repository.NewTaskRepository(p.Db)
	taskService := services.NewTaskService(taskRepo, p.Logger, p.Ctx)
	task, err := taskService.GetOneById(taskId)
	if err != nil {
		p.Logger.Errorf("Ошибка при получении задания: %v", err)
	}

	p.Description = fmt.Sprintf("*%v*\nОписание:\n%v\n\nТы можешь принять это задание или вернуться к списку заданий", task.GetName(), task.GetShortDescription())

	callbackDataJSON, err := utils.CreateCallbackDataJSON(map[string]string{"id": strconv.Itoa(int(task.ID))})
	if err != nil {
		// Обработка ошибки
	}
	numericKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Принять", pageCode.TaskAccept+constants.ParamsSeparator+string(callbackDataJSON)),
			tgbotapi.NewInlineKeyboardButtonData("Назад", pageCode.TasksAvailable),
		),
	)
	p.KeyBoard = &numericKeyboard
}
