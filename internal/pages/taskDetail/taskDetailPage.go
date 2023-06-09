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
	StatusCode "table10/internal/constants/statusCode"
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
	//case "to_review":
	//	p.ToReview()
	case "under_review":
		p.UnderReview()
	default:
		p.Detail()
	}
}

// Detail детальная страница задания
func (p *page) Detail() {
	task := p.task
	p.Description = fmt.Sprintf("*%v*\n%v\n\n", task.GetName(), task.GetShortDescription())

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
	p.Description = fmt.Sprintf("*%v*\n%v\n\n", task.GetName(), task.GetShortDescription())

	howToAnswer := fmt.Sprintf("_Под этим сообщением_ напиши ответы на вопросы из задания и\\/или прикрепи фото \\(прикрепи его именно как фото, а не файл\\)\\. Сообщения можно отправлять несколько раз\\. Когда будешь готов, нажми \\'Сдать задание\\', чтобы все ответы отправть на проверку\\.")
	p.Description += howToAnswer

	userText := p.GetUserText()
	userPhoto := p.GetUserPhoto()
	//Если были отправлены ответы или файлы
	if userText != "" || userPhoto != nil {

		answerRepo := repository.NewAnswerRepository(p.Db)
		answerService := services.NewAnswerService(answerRepo, p.Logger, p.Ctx)
		err := answerService.AddAnswer(userText, userPhoto, p.User, &p.task.UserTasks[0])
		if err != nil {
			p.Logger.Errorf("Ошибка добавления ответа пользователя к заданию: %v", err)
			p.Description = "Ошибка добавления ответа"
			return
		} else {
			p.Description = fmt.Sprintf("Ответ записан\n\n")
		}
		p.Description += howToAnswer

	} else {
		userTasks, err := p.taskService.GetUserTasks(&repository.UserTaskFilter{Task: task, PlayWithYou: true})
		if err != nil {
			//
		} else {
			var userTasksDescriptions string
			if len(userTasks) > 0 {
				userTasksDescriptions += fmt.Sprintf("\n_Вместе с тобой взяли это задание\\:_ ")
				for _, userTask := range userTasks {
					userTasksDescriptions += fmt.Sprintf("[@%s](tg://user?id=%d) ", userTask.User.Username, userTask.User.TelegramID)
				}
			} else {
				userTasksDescriptions += fmt.Sprintf("\n_Ты первый взял это задание_")
			}
			p.Description += userTasksDescriptions
		}

	}

	callbackDataJSON, err := utils.CreateCallbackDataJSON(map[string]string{"id": strconv.Itoa(int(p.task.ID)), "action": "under_review"})
	if err != nil {
		// Обработка ошибки
	}
	numericKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Сдать задание", pageCode.TaskDetail+constants.ParamsSeparator+string(callbackDataJSON)),
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
		p.Description = fmt.Sprintf("Ошибка получения списка ответов")
		p.Logger.Errorf("Ошибка получения списка ответов. userTask id = %v", p.task.UserTasks[0])
		return
	}
	//Если у пользователя нет ответов на задание
	if len(answers) == 0 {
		p.Description = fmt.Sprintf("Ошибка\\: ты не отправил ответы на задание\\.\nПерейди по кнопке \\'Вернуться к отправке ответа\\' и отправь ответы на задание\\.")
		callbackDataJSON, err := utils.CreateCallbackDataJSON(map[string]string{"id": strconv.Itoa(int(p.task.ID)), "action": "in_progress"})
		if err != nil {
			// Обработка ошибки
		}
		numericKeyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Вернуться к отправке ответа", pageCode.TaskDetail+constants.ParamsSeparator+string(callbackDataJSON)),
			),
		)
		p.KeyBoard = &numericKeyboard
		return
	}

	p.Description = fmt.Sprintf("Задание было отправлено на проверку\\.\nПодожди, пока наши модераторы проверят задание\\.")
	if p.task.UserTasks[0].Status.Code == StatusCode.InProgress {
		err := p.taskService.ChangeStatus(&p.task.UserTasks[0], StatusCode.UnderReview)
		if err != nil {
			p.Description = fmt.Sprintf("Ошибка отправки задания на проверку")
		}
	} else {
		p.Description += fmt.Sprintf("\n\n*%v*\n%v", p.task.GetName(), p.task.GetShortDescription())
	}

	numericKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			//tgbotapi.NewInlineKeyboardButtonData("На проверку", pageCode.TaskDetail+constants.ParamsSeparator+string(callbackDataJSON)),
			tgbotapi.NewInlineKeyboardButtonData("Назад", pageCode.TasksAccepted),
		),
	)
	p.KeyBoard = &numericKeyboard
}
