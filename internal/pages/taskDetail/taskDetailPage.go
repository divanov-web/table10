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
	task        *models.Task
	taskService *services.TaskService
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
	task, err := p.taskService.GetOneById(taskId, &repository.TaskFilter{User: p.User})
	if err != nil {
		p.Logger.Errorf("Ошибка при получении задания: %v", err)
	}
	p.task = task
	if len(task.Users) > 0 && action == "default" {
		action = "play"
	}
	switch action {
	case "accept":
		p.Accept()
	case "play":
		p.Play()
	case "to_review":
		p.ToReview()
	default:
		p.Detail()
	}
}

// Detail детальная страница задания
func (p *page) Detail() {
	task := p.task
	p.Description = fmt.Sprintf("*%v*\nОписание:\n%v\n\nТы можешь принять это задание или вернуться к списку заданий", task.GetName(), task.GetShortDescription())

	callbackDataJSON, err := utils.CreateCallbackDataJSON(map[string]string{"id": strconv.Itoa(int(task.ID)), "action": "accept"})
	if err != nil {
		// Обработка ошибки
	}
	numericKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Принять", pageCode.TaskDetail+constants.ParamsSeparator+string(callbackDataJSON)),
			tgbotapi.NewInlineKeyboardButtonData("Назад", pageCode.Tasks),
		),
	)
	p.KeyBoard = &numericKeyboard
}

// Play Детальная страница задания, если она уже в процессе игры
func (p *page) Play() {
	task := p.task
	p.Description = fmt.Sprintf("*%v*\nОписание:\n%v\n\nТы уже выполняешь это задание\\.\nМожно отправть его на проверку\\.", task.GetName(), task.GetShortDescription())

	callbackDataJSON, err := utils.CreateCallbackDataJSON(map[string]string{"id": strconv.Itoa(int(task.ID)), "action": "to_review"})
	if err != nil {
		// Обработка ошибки
	}
	numericKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("На проверку", pageCode.TaskDetail+constants.ParamsSeparator+string(callbackDataJSON)),
			tgbotapi.NewInlineKeyboardButtonData("Назад", pageCode.Tasks),
		),
	)
	p.KeyBoard = &numericKeyboard
}

// ToReview отправить задание на проверку
func (p *page) ToReview() {
	task := p.task
	p.Description = fmt.Sprintf("*%v*\nОписание:\n%v\n\nЗадание было отправлено на проверку\\.\nПодожди, пока наши модераторы проверят задание\\.", task.GetName(), task.GetShortDescription())

	numericKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			//tgbotapi.NewInlineKeyboardButtonData("На проверку", pageCode.TaskDetail+constants.ParamsSeparator+string(callbackDataJSON)),
			tgbotapi.NewInlineKeyboardButtonData("Назад", pageCode.Tasks),
		),
	)
	p.KeyBoard = &numericKeyboard
}

// Accept Станица принятия задания по id
func (p *page) Accept() {
	err := p.taskService.AddUserToTask(p.task, p.User)
	if err != nil {
		p.Logger.Errorf("Ошибка добавления пользователя в задания")
		p.Description = "Ошибка принятия задания"
	} else {
		p.Description = fmt.Sprintf("Вы успешно приняли задание %v", p.task.GetName())
		numericKeyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Список заданий", pageCode.Tasks),
			),
		)
		p.KeyBoard = &numericKeyboard
	}
}
