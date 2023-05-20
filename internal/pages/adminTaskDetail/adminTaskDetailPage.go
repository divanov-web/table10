package adminTaskDetailPage

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
	"table10/internal/services/task_straregy"
	"table10/pkg/logging"
	"table10/pkg/utils"
)

type page struct {
	base.AbstractPage
	task         *models.Task
	taskService  *services.TaskService
	taskStrategy task_straregy.TaskProgressionStrategy
}

func NewPage(db *gorm.DB, logger *logging.Logger, ctx context.Context, user *models.User, callbackData *callbackdata.CallbackData) interfaces.Page {

	return &page{
		AbstractPage: base.AbstractPage{
			Db:           db,
			Logger:       logger,
			Ctx:          ctx,
			User:         user,
			Name:         "Детальная задания юзера",
			Description:  "",
			Code:         pageCode.AdminTaskDetail,
			KeyBoard:     nil,
			CallbackData: callbackData,
		},
	}
}

func (p *page) Generate() {
	userTaskId, err := p.CallbackData.GetId()
	if err != nil {
		p.Logger.Errorf("%v", err)
		return
	}

	taskRepo := repository.NewTaskRepository(p.Db)
	userRepo := repository.NewUserRepository(p.Db)
	statusRepo := repository.NewStatusRepository(p.Db)
	taskService := services.NewTaskService(taskRepo, userRepo, statusRepo, p.Logger, p.Ctx)
	userTask, err1 := taskService.GetUserTaskById(userTaskId)
	if err1 != nil {
		p.Logger.Errorf("Ошибка при получении задания пользователя: %v", err1)
	}

	p.Description = fmt.Sprintf("*Выполняет:*[@%s](tg://user?id=%d)\n*%v*\nОписание:\n%v\n\n", userTask.User.Username, userTask.User.TelegramID, userTask.Task.GetName(), userTask.Task.GetShortDescription())

	callbackDataJSONAccept, err := utils.CreateCallbackDataJSON(map[string]string{"id": strconv.Itoa(int(userTask.Task.ID)), "action": "accept"})
	if err != nil {
		// Обработка ошибки
	}

	callbackDataJSONReject, err := utils.CreateCallbackDataJSON(map[string]string{"id": strconv.Itoa(int(userTask.Task.ID)), "action": "reject"})
	if err != nil {
		// Обработка ошибки
	}
	numericKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Подтвердить", pageCode.AdminTaskDetail+constants.ParamsSeparator+string(callbackDataJSONAccept)),
			tgbotapi.NewInlineKeyboardButtonData("Вернуть", pageCode.AdminTaskDetail+constants.ParamsSeparator+string(callbackDataJSONReject)),
			tgbotapi.NewInlineKeyboardButtonData("Назад", pageCode.AdminReview),
		),
	)
	p.KeyBoard = &numericKeyboard
}
