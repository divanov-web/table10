package gamePageInput

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"strings"
	"table10/internal/models"
	"table10/internal/pages/base"
	"table10/internal/pages/interfaces"
	"table10/internal/repository"
	"table10/pkg/logging"
)

const Command = "game_input"

type page struct {
	base.AbstractPage
}

func NewPage(db *gorm.DB, logger *logging.Logger, ctx context.Context, user *models.User) interfaces.Page {

	return &page{
		AbstractPage: base.AbstractPage{
			Db:          db,
			Logger:      logger,
			Ctx:         ctx,
			User:        user,
			Name:        "Поиск ко коду",
			Description: "Введите код сервера игры:",
			Command:     Command,
			KeyBoard:    nil,
		},
	}
}

func (p *page) Generate() {
	userText := p.GetUserText()
	words := strings.Split(userText, " ")
	searchCode := words[0]
	var descriptionText string
	if userText != "" {
		gameRepo := repository.NewGameRepository(p.Db)
		currentGame, err := gameRepo.GetOne(p.Ctx, searchCode)
		if err != nil {
			descriptionText = fmt.Sprintf("Игра с кодом %v не найдена", searchCode)
		} else {
			descriptionText = fmt.Sprintf("Найден сервер игры <b>%v</b>. \nОписание: \n%v", currentGame.Name, currentGame.GetShortDescription())
		}

		p.Description = descriptionText
	}
}
