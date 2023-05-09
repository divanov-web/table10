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
	task, taskStrategy, err1 := p.taskService.GetOneById(taskId, &repository.TaskFilter{User: p.User})
	if err1 != nil {
		p.Logger.Errorf("Ошибка при получении задания: %v", err1)
	}
	p.task = task
	p.taskStrategy = taskStrategy
	if len(task.UserTasks) > 0 && action == "default" {
		action = task.UserTasks[0].Status.Code
	}
	switch action {
	case "accept":
		p.Accept()
	case "in_progress":
		p.InProgress()
	case "to_review":
		p.ToReview()
	case "under_review":
		p.UnderReview()
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

// Accept Станица принятия задания по id
func (p *page) Accept() {
	err := p.taskService.AddUserToTask(p.task, p.User, p.taskStrategy)
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

// InProgress Детальная страница задания, если она уже в процессе игры
func (p *page) InProgress() {
	task := p.task
	p.Description = fmt.Sprintf("*%v*\nОписание:\n%v\n\nТы уже выполняешь это задание\\.\nМожно отправть его на проверку\\.", task.GetName(), task.GetShortDescription())

	callbackDataJSON, err := utils.CreateCallbackDataJSON(map[string]string{"id": strconv.Itoa(int(task.ID)), "action": "to_review"})
	if err != nil {
		// Обработка ошибки
	}
	numericKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("На проверку", pageCode.TaskDetail+constants.ParamsSeparator+string(callbackDataJSON)),
			tgbotapi.NewInlineKeyboardButtonData("Назад", pageCode.TasksAccepted),
		),
	)
	p.KeyBoard = &numericKeyboard
}

// ToReview отправить ответы и файлы перед подтверждением отправки задания на проверку
func (p *page) ToReview() {
	p.Description = fmt.Sprintf("Напиши ответы на вопросы из задания и\\/или прикрепи фото \\(прикрепи его именно как фото, а не файл\\)\\. Сообщения можно отправлять несколько раз\\. Нажмите на кнопку Подтвердить только после отправки сообщений с ответами\\.")
	userText := p.GetUserText()
	userPhoto := p.GetUserPhoto()
	answerRepo := repository.NewAnswerRepository(p.Db)
	answerService := services.NewAnswerService(answerRepo, p.Logger, p.Ctx)
	if userText != "" || userPhoto != nil {
		err := answerService.AddAnswer(userText, userPhoto, p.User, p.task)
		if err != nil {
			p.Logger.Errorf("Ошибка добавления ответа пользователя к заданию: %v", err)
			p.Description = "Ошибка добавления ответа"
			return
		}
	}
	callbackDataJSON, err := utils.CreateCallbackDataJSON(map[string]string{"id": strconv.Itoa(int(p.task.ID)), "action": "under_review"})
	if err != nil {
		// Обработка ошибки
	}
	numericKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Подтвердить", pageCode.TaskDetail+constants.ParamsSeparator+string(callbackDataJSON)),
			tgbotapi.NewInlineKeyboardButtonData("Мои задания", pageCode.TasksAccepted),
		),
	)
	p.KeyBoard = &numericKeyboard
}

// UnderReview отправить задание на проверку
func (p *page) UnderReview() {
	answerRepo := repository.NewAnswerRepository(p.Db)
	answerService := services.NewAnswerService(answerRepo, p.Logger, p.Ctx)
	answers, err := answerService.GetAnswers(&repository.AnswerFilter{UserTask: &p.task.UserTasks[0]})
	if err != nil {
		p.Description = fmt.Sprintf("Ошибка отправки задания на рповерку")
		p.Logger.Errorf("Ошибка отправки задания на проверку. userTask id = %v", p.task.UserTasks[0])
		return
	}
	//Если у пользователя нет ответов на задание
	if len(answers) == 0 {
		p.Description = fmt.Sprintf("Ошибка\\: ты не отправил ответы на задание\\.\nПерейди по кнопке \\'Отправка ответов\\' и отправь ответы на задание\\.")
		callbackDataJSON, err := utils.CreateCallbackDataJSON(map[string]string{"id": strconv.Itoa(int(p.task.ID)), "action": "to_review"})
		if err != nil {
			// Обработка ошибки
		}
		numericKeyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Отправка ответов", pageCode.TaskDetail+constants.ParamsSeparator+string(callbackDataJSON)),
			),
		)
		p.KeyBoard = &numericKeyboard
		return
	}

	p.Description = fmt.Sprintf("Задание было отправлено на проверку\\.\nПодожди, пока наши модераторы проверят задание\\.")
	if p.task.UserTasks[0].Status.Code == "in_progress" {
		err := p.taskService.ChangeStatus(p.task, "under_review")
		if err != nil {
			p.Description = fmt.Sprintf("Ошибка отправки задания на проверку")
		}
	} else {
		p.Description += fmt.Sprintf("\n\n*%v*\nОписание:\n%v", p.task.GetName(), p.task.GetShortDescription())
	}

	numericKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			//tgbotapi.NewInlineKeyboardButtonData("На проверку", pageCode.TaskDetail+constants.ParamsSeparator+string(callbackDataJSON)),
			tgbotapi.NewInlineKeyboardButtonData("Назад", pageCode.TasksAccepted),
		),
	)
	p.KeyBoard = &numericKeyboard
}
