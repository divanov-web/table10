package pages

import (
	"context"
	"gorm.io/gorm"
	"table10/internal/models"
	"table10/internal/pages/cabinet"
	gamePage "table10/internal/pages/game"
	gamePageInput "table10/internal/pages/gameInput"
	"table10/internal/pages/interfaces"
	mainpage "table10/internal/pages/main"
	tasksPage "table10/internal/pages/tasks"
	"table10/pkg/logging"
)

type PageFactory struct {
	db *gorm.DB
}

func NewPageFactory(db *gorm.DB) *PageFactory {
	return &PageFactory{
		db: db,
	}
}

func (pf *PageFactory) CreatePage(pageName string, logger *logging.Logger, user *models.User, ctx context.Context) interfaces.Page {

	var page interfaces.Page
	switch pageName {
	case "main":
		page = mainpage.NewPage(pf.db, logger, ctx, user)
	case "cabinet":
		page = cabinetPage.NewPage(pf.db, logger, ctx, user)
	case "game":
		page = gamePage.NewPage(pf.db, logger, ctx, user)
	case "game_input":
		page = gamePageInput.NewPage(pf.db, logger, ctx, user)
	case "tasks":
		page = tasksPage.NewPage(pf.db, logger, ctx, user)
	default:
		page = mainpage.NewPage(pf.db, logger, ctx, user)
	}

	return page
}
