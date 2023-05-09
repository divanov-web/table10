package welcomePage

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
	"table10/internal/constants/pageCode"
	"table10/internal/models"
	"table10/internal/pages/base"
	"table10/internal/pages/interfaces"
	"table10/internal/repository"
	"table10/internal/services"
	"table10/pkg/logging"
)

type page struct {
	base.AbstractPage
}

func NewPage(db *gorm.DB, logger *logging.Logger, ctx context.Context, user *models.User) interfaces.Page {
	numericKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Меню", pageCode.Main),
		),
	)

	var gameDescription string
	gameCode := "tashkent"
	gameRepo := repository.NewGameRepository(db)
	userRepo := repository.NewUserRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	gameService := services.NewGameService(gameRepo, roleRepo, userRepo, logger, ctx)
	currentGame, err := gameService.GetOneByCode(gameCode)
	if err != nil {
		logger.Errorf("Can't add default game: %v", err)
	}
	gameDescription += fmt.Sprintf("Сервер игры: *%v*\\. \nОписание: \n%v\n", currentGame.GetName(), currentGame.GetLongDescription())

	//Если пользователь не участвует ни в одной игре, то добавляем его в игру tashkent
	if len(user.Games) == 0 {
		game, err1 := gameService.AddUserToGame(int(currentGame.ID), user)
		if err1 != nil {
			logger.Errorf("Ошибка добавления пользователя в дефолтную игру")
			gameDescription = "Ошибка принятия игры"
		} else {
			gameDescription += fmt.Sprintf("\nВы успешно добавлены в игру %v", game.Name)
		}
	}

	return &page{
		AbstractPage: base.AbstractPage{
			Name: "Добро пожаловать",
			Description: "*Table 10* \\- социальная игра, которая мотивирует к личностному развитию и общению с новыми людьми\\.\n" +
				"Участники получают от телеграм\\-бота [Table10 Bot](tg://user?id=bigtable10_bot) еженедельные задания, предназначенные для развлечения, посещения интересных мест и других занятий, " +
				"которые в обычной жизни ты можешь проигнорировать или полениться посетить\\.\n" +
				"Задания выполняются в команде или паре, и фокусируются на увлечениях, развлечениях и расширении кругозора\\.\n" +
				"Игра не предназначена для романтических знакомств\\.\n\n" + gameDescription,
			Code:     pageCode.Main,
			KeyBoard: &numericKeyboard,
			UserText: "",
		},
	}
}
