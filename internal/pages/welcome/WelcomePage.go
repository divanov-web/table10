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
	gameDescription += fmt.Sprintf("Вы участвуете в игре *%v*\\. \nОписание: \n%v\n", currentGame.GetName(), currentGame.GetShortDescription())

	//Если пользователь не участвует ни в одной игре, то добавляем его в игру tashkent
	if len(user.Games) == 0 {
		game, err1 := gameService.AddUserToGame(int(currentGame.ID), user)
		if err1 != nil {
			logger.Errorf("Ошибка добавления пользователя в дефолтную игру")
			gameDescription = "Ошибка принятия игры"
		} else {
			gameDescription += fmt.Sprintf("Вы успешно добавлены в игру %v", game.Name)
		}
	}

	return &page{
		AbstractPage: base.AbstractPage{
			Name:        "Добро пожаловать",
			Description: "Добро пожаловать в игру Table 10\\.\nТут будет описание игры\\.\n\n" + gameDescription,
			Code:        pageCode.Main,
			KeyBoard:    &numericKeyboard,
			UserText:    "",
		},
	}
}
