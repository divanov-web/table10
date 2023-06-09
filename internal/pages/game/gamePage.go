package gamePage

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
	"strings"
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
			tgbotapi.NewInlineKeyboardButtonData("Поиск по коду", pageCode.GameInput),
			tgbotapi.NewInlineKeyboardButtonData("Меню", pageCode.Main),
		),
	)
	return &page{
		AbstractPage: base.AbstractPage{
			Db:          db,
			Logger:      logger,
			Ctx:         ctx,
			User:        user,
			Name:        "Список серверов",
			Description: "Просмотр список серверов, в которых вы участвуете",
			Code:        pageCode.Game,
			KeyBoard:    &numericKeyboard,
		},
	}
}

func (p *page) Generate() {
	gameRepo := repository.NewGameRepository(p.Db)
	userRepo := repository.NewUserRepository(p.Db)
	roleRepo := repository.NewRoleRepository(p.Db)
	gameService := services.NewGameService(gameRepo, roleRepo, userRepo, p.Logger, p.Ctx)
	games, err := gameService.GetUserGames(p.User)
	if err != nil {
		p.Logger.Errorf("Ошибка поиска игр пользователя")
		p.Description = fmt.Sprintf("Ошибка поиска игр")
	}
	if len(games) == 0 {
		p.Description = fmt.Sprintf("У вас нет активных игр")
	} else {
		gameDescriptions := make([]string, len(games))
		for i, userGame := range games {
			gameDescriptions[i] = userGame.Game.GetName()
		}
		p.Description = fmt.Sprintf("Список активных игр:\n%s", strings.Join(gameDescriptions, "\n"))
	}

}
