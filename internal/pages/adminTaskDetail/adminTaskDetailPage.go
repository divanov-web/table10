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
	userTask    *models.UserTask
	taskService *services.TaskService
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
	action := p.CallbackData.GetAction()

	taskRepo := repository.NewTaskRepository(p.Db)
	userRepo := repository.NewUserRepository(p.Db)
	statusRepo := repository.NewStatusRepository(p.Db)
	p.taskService = services.NewTaskService(taskRepo, userRepo, statusRepo, p.Logger, p.Ctx)
	userTask, err1 := p.taskService.GetUserTaskById(userTaskId)
	p.userTask = userTask
	if err1 != nil {
		p.Logger.Errorf("Ошибка при получении задания пользователя: %v", err1)
	}

	p.Description = fmt.Sprintf("*Выполняет:*[@%s](tg://user?id=%d)\n*%v*\nОписание:\n%v\n\n", userTask.User.Username, userTask.User.TelegramID, userTask.Task.GetName(), userTask.Task.GetShortDescription())

	switch action {
	case "accept":
		p.Accept()
	case "reject":
		p.Reject()
	default:
		p.Detail()
	}
}

func (p *page) Detail() {
	callbackDataJSONAccept, err := utils.CreateCallbackDataJSON(map[string]string{"id": strconv.Itoa(int(p.userTask.ID)), "action": "accept"})
	if err != nil {
		// Обработка ошибки
	}

	callbackDataJSONReject, err := utils.CreateCallbackDataJSON(map[string]string{"id": strconv.Itoa(int(p.userTask.ID)), "action": "reject"})
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

	//list of User answers
	answerRepo := repository.NewAnswerRepository(p.Db)
	answerService := services.NewAnswerService(answerRepo, p.Logger, p.Ctx)
	answers, err := answerService.GetAnswers(&repository.AnswerFilter{UserTask: p.userTask})
	if err != nil {
		p.Description += fmt.Sprintf("\n\n*Ошибка получения списка ответов*")
		p.Logger.Errorf("Ошибка получения списка ответов. userTask id = %v", p.userTask.ID)
		return
	}
	var answerMessages []telegram.Message
	for _, answer := range answers {
		var text string
		if answer.Text != "" {
			text = fmt.Sprintf("[@%s](tg://user?id=%d)\\: ", answer.User.Username, answer.User.TelegramID) + answer.Text
		}
		message := telegram.Message{
			Text: text,
			Photo: telegram.Photo{
				FileId: answer.TelegramFileId,
				Url:    answer.ImagePath,
			},
		}
		answerMessages = append(answerMessages, message)
	}
	p.Messages = answerMessages
}

func (p *page) Accept() {

}

func (p *page) Reject() {

}
