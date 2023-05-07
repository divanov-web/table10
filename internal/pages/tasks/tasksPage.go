// Package tasksPage предоставляет страницу со списком доступных заданий для пользователей.
package tasksPage

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
	"strconv"
	"strings"
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
	numericKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", pageCode.Main),
		),
	)
	return &page{
		AbstractPage: base.AbstractPage{
			Db:           db,
			Logger:       logger,
			Ctx:          ctx,
			User:         user,
			Name:         "Список заданий",
			Description:  "",
			Code:         pageCode.Tasks,
			KeyBoard:     &numericKeyboard,
			CallbackData: callbackData,
		},
	}
}

func (p *page) Generate() {
	game := p.User.Games[0].Game

	taskRepo := repository.NewTaskRepository(p.Db)
	userRepo := repository.NewUserRepository(p.Db)
	statusRepo := repository.NewStatusRepository(p.Db)
	taskService := services.NewTaskService(taskRepo, userRepo, statusRepo, p.Logger, p.Ctx)

	var filter *repository.TaskFilter

	filter = &repository.TaskFilter{
		Current:           true,
		User:              p.User,
		NotAssignedToUser: true,
	}

	tasks, err := taskService.GetTasks(&game, filter)
	if err != nil {
		p.Logger.Errorf("Ошибка при получении заданий: %v", err)
	}

	if len(tasks) == 0 {
		p.Description = fmt.Sprintf("На сегодняшний день нет доступных заданий, либо ты уже принял все задания\\.")
	} else {
		var taskDescriptions []string
		for _, task := range tasks {
			taskDescriptions = append(taskDescriptions, fmt.Sprintf("*%s* \\(%s\\)", task.Name, task.TaskType.GetName()))
		}

		taskList := strings.Join(taskDescriptions, "\n")
		p.Description = fmt.Sprintf("Список доступных заданий игры %v\\:\n%s", game.Name, taskList)
	}

	//Создание новых кнопок с заданиями
	taskButtons := make([][]tgbotapi.InlineKeyboardButton, 0)
	for _, task := range tasks {
		taskButton, err := createTaskButton(task)
		if err != nil {
			// Обработка ошибки
		}
		taskButtons = append(taskButtons, tgbotapi.NewInlineKeyboardRow(taskButton))
	}

	// Кнопка "Назад"
	backButtonRow := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Назад", pageCode.Main),
	)
	taskButtons = append(taskButtons, backButtonRow)

	p.KeyBoard = &tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: taskButtons,
	}
}

func createTaskButton(task models.Task) (tgbotapi.InlineKeyboardButton, error) {
	callbackDataJSON, err := utils.CreateCallbackDataJSON(map[string]string{"id": strconv.Itoa(int(task.ID))})
	if err != nil {
		return tgbotapi.InlineKeyboardButton{}, err
	}
	return tgbotapi.NewInlineKeyboardButtonData(task.GetClearedName(), pageCode.TaskDetail+constants.ParamsSeparator+string(callbackDataJSON)), nil
}
