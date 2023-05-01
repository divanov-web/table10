package taskAvailablePage

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
	"strconv"
	"strings"
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

	taskRepo := repository.NewTaskRepository(p.Db)
	taskService := services.NewTaskService(taskRepo, p.Logger, p.Ctx)
	tasks, err := taskService.GetTasks(currentPeriod)
	if err != nil {
		p.Logger.Errorf("Ошибка при получении заданий: %v", err)
	}

	var taskDescriptions []string
	for _, task := range tasks {
		taskDescriptions = append(taskDescriptions, fmt.Sprintf("*%s* \\(%s\\)", task.GetName(), task.TaskType.GetName()))
	}

	taskList := strings.Join(taskDescriptions, "\n")
	p.Description = fmt.Sprintf("Список доступных заданий на неделе %v\\:\n%s", currentPeriod.WeekNumber, taskList)

	//Создание новых кнопок с заданиями
	taskButtons := make([][]tgbotapi.InlineKeyboardButton, 0)
	for _, task := range tasks {
		taskButton, err := createTaskButton(task)
		if err != nil {
			// Обработка ошибки
		}
		taskButtons = append(taskButtons, tgbotapi.NewInlineKeyboardRow(taskButton))
	}

	// Добавьте кнопку "Назад"
	backButtonRow := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Назад", pageCode.Tasks),
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
	return tgbotapi.NewInlineKeyboardButtonData(task.Name, pageCode.TaskDetail+constants.ParamsSeparator+string(callbackDataJSON)), nil
}
