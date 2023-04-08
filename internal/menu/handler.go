package menu

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
	"table10/internal/handlers"
	"table10/internal/models"
	"table10/internal/pages"
	"table10/internal/pages/interfaces"
	"table10/pkg/logging"
)

type handler struct {
	logger      *logging.Logger
	db          *gorm.DB
	user        *models.User
	pageFactory *pages.PageFactory
	ctx         context.Context
	pages       map[string]interfaces.Page
}

func NewHandler(logger *logging.Logger, db *gorm.DB, user *models.User, ctx context.Context) handlers.Handler {
	pageFactory := pages.NewPageFactory(db)
	return &handler{
		logger:      logger,
		db:          db,
		user:        user,
		pageFactory: pageFactory,
		ctx:         ctx,
		pages:       make(map[string]interfaces.Page),
	}
}

func (h *handler) getPage(pageName string) interfaces.Page {
	if page, ok := h.pages[pageName]; ok {
		return page
	}

	newPage := h.pageFactory.CreatePage(pageName, h.logger, h.user, h.ctx)
	h.pages[pageName] = newPage
	return newPage
}

func (h *handler) Register(tgbotapi *tgbotapi.Update) (page interfaces.Page) {
	var pageName string
	if tgbotapi.CallbackQuery != nil {
		pageName = tgbotapi.CallbackQuery.Data
	} else {
		if h.user.LastPage != "" {
			pageName = h.user.LastPage
		} else {
			pageName = "main"
		}
	}
	h.logger.Infof("Current page = %v", pageName)
	page = h.getPage(pageName)
	return page
}
