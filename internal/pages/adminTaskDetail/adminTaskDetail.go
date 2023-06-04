// Package adminTaskDetailPage Task detail page for administration
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
	"table10/internal/structs/telegram"
	"table10/pkg/logging"
	"table10/pkg/utils"
)

type page struct {
	base.AbstractPage
	task        *models.Task
	taskService *services.TaskService
}

func NewPage(db *gorm.DB, logger *logging.Logger, ctx context.Context, user *models.User, callbackData *callbackdata.CallbackData) interfaces.Page {
	numericKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", pageCode.AdminTasks),
		),
	)
	return &page{
		AbstractPage: base.AbstractPage{
			Db:           db,
			Logger:       logger,
			Ctx:          ctx,
			User:         user,
			Name:         "Задание",
			Description:  "",
			Code:         pageCode.AdminTaskDetail,
			KeyBoard:     &numericKeyboard,
			CallbackData: callbackData,
		},
	}
}

func (p *page) Generate() {
	taskId, err := p.CallbackData.GetId()
	if err != nil {
		p.Logger.Errorf("%v", err)
		return
	}
	action := p.CallbackData.GetAction()

	taskRepo := repository.NewTaskRepository(p.Db)
	userRepo := repository.NewUserRepository(p.Db)
	statusRepo := repository.NewStatusRepository(p.Db)
	p.taskService = services.NewTaskService(taskRepo, userRepo, statusRepo, p.Logger, p.Ctx)
	task, _, err1 := p.taskService.GetOneById(taskId, &repository.TaskFilter{})
	if err1 != nil {
		p.Logger.Errorf("Ошибка при получении задания: %v", err1)
	}
	p.task = task

	switch action {
	case "activate":
		p.Activate()
	default:
		p.NotActive()
	}
}

func (p *page) NotActive() {
	task := p.task
	p.Description = fmt.Sprintf("*%v*\n_Не активно_\n%v", task.GetName(), task.GetShortDescription())

	callbackDataJSON, err := utils.CreateCallbackDataJSON(map[string]string{"id": strconv.Itoa(int(task.ID)), "action": "activate"})
	if err != nil {
		// Обработка ошибки
	}
	numericKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Активировать", pageCode.AdminTaskDetail+constants.ParamsSeparator+string(callbackDataJSON)),
			tgbotapi.NewInlineKeyboardButtonData("Назад", pageCode.AdminTasks),
		),
	)
	p.KeyBoard = &numericKeyboard
}

func (p *page) Activate() {
	task := p.task

	err := p.taskService.ChangeActive(task, true)
	if err != nil {
		p.Logger.Errorf("Ошибка активации задания: %v", err)
		p.Description = fmt.Sprintf("Ошибка активации задания: %v", err)
		return
	}

	p.Description = fmt.Sprintf("*%v*\n_Активировано_\n", task.GetName())

	numericKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", pageCode.AdminTasks),
		),
	)
	p.KeyBoard = &numericKeyboard

	//Additional message to channel about new task
	var answerMessages []telegram.Message
	message := telegram.Message{
		Text:   fmt.Sprintf("Опубликовано новое задание\n\n*%v*\n%v", task.GetName(), task.GetShortDescription()),
		ChatId: task.Game.ChatId,
	}
	p.Messages = append(answerMessages, message)
}
