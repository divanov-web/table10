package gameAcceptPage

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
	"strconv"
	"table10/internal/callbackdata"
	"table10/internal/constants"
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

func NewPage(db *gorm.DB, logger *logging.Logger, ctx context.Context, user *models.User, callbackData *callbackdata.CallbackData) interfaces.Page {
	numericKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", constants.GamePageCode),
		),
	)

	return &page{
		AbstractPage: base.AbstractPage{
			Db:           db,
			Logger:       logger,
			Ctx:          ctx,
			User:         user,
			Name:         "Вступление в сервер игры",
			Description:  "Вы добавлены к серверу",
			Code:         constants.GameAcceptPageCode,
			KeyBoard:     &numericKeyboard,
			CallbackData: callbackData,
		},
	}
}

func (p *page) Generate() {
	gameIdString, ok := p.CallbackData.Params["id"]
	if !ok {
		p.Logger.Errorf("Ошибка принятия сервера из-за ошибки передачи параметра")
	}
	gameId, err := strconv.Atoi(gameIdString)
	if err != nil {
		p.Logger.Errorf("Ошибка принятия сервера из-за ошибки передачи параметра")
	}
	gameRepo := repository.NewGameRepository(p.Db)
	userRepo := repository.NewUserRepository(p.Db)
	roleRepo := repository.NewRoleRepository(p.Db)
	gameService := services.NewGameService(gameRepo, roleRepo, userRepo, p.Logger, p.Ctx)
	game, err1 := gameService.AddUserToGame(gameId, p.User)
	if err1 != nil {
		p.Logger.Errorf("Ошибка добавления пользователя в игру")
		p.Description = "Ошибка принятия игры"
	} else {
		p.Description = fmt.Sprintf("Вы успешно добавлены в игру %v", game.Name)
		numericKeyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Главное меню", constants.MainPageCode),
			),
		)
		p.KeyBoard = &numericKeyboard
	}

}
