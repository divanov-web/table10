// Package adminReviewPage Page with tasks under review
package adminReviewPage

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
	StatusCode "table10/internal/constants/statusCode"
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
			tgbotapi.NewInlineKeyboardButtonData("Назад", pageCode.Admin),
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
			Code:         pageCode.AdminReview,
			KeyBoard:     &numericKeyboard,
			CallbackData: callbackData,
		},
	}
}

func (p *page) Generate() {
	gameId := p.User.Games[0].GameID

	taskRepo := repository.NewTaskRepository(p.Db)
	userRepo := repository.NewUserRepository(p.Db)
	statusRepo := repository.NewStatusRepository(p.Db)
	taskService := services.NewTaskService(taskRepo, userRepo, statusRepo, p.Logger, p.Ctx)

	filter := &repository.UserTaskFilter{
		GameId:     gameId,
		StatusCode: StatusCode.UnderReview,
	}

	userTasks, err := taskService.GetUserTasks(filter)
	if err != nil {
		p.Logger.Errorf("Ошибка при получении заданий: %v", err)
	}

	if len(userTasks) == 0 {
		p.Description = fmt.Sprintf("Нет заданий на проверку\\.")
	} else {
		var taskDescriptions []string
		for _, userTask := range userTasks {
			taskDescriptions = append(taskDescriptions, fmt.Sprintf("*%s* \\(%s\\)", userTask.Task.GetName(), userTask.User.GetUserName()))
		}

		taskList := strings.Join(taskDescriptions, "\n")

		p.Description = fmt.Sprintf("Список заданий на проверку \\:\n%s", taskList)
	}

	//Создание новых кнопок с заданиями
	taskButtons := make([][]tgbotapi.InlineKeyboardButton, 0)
	for _, userTask := range userTasks {
		taskButton, err := createTaskButton(userTask)
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

func createTaskButton(userTask models.UserTask) (tgbotapi.InlineKeyboardButton, error) {
	callbackDataJSON, err := utils.CreateCallbackDataJSON(map[string]string{"id": strconv.Itoa(int(userTask.ID))})
	if err != nil {
		return tgbotapi.InlineKeyboardButton{}, err
	}
	return tgbotapi.NewInlineKeyboardButtonData(userTask.Task.GetClearedName(), pageCode.AdminUserTaskDetail+constants.ParamsSeparator+string(callbackDataJSON)), nil
}
