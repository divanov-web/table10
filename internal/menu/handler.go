package menu

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
	"table10/internal/handlers"
	"table10/internal/pages"
	"table10/internal/pages/interfaces"
	"table10/pkg/logging"
)

type handler struct {
	logger      *logging.Logger
	db          *gorm.DB
	pageFactory *pages.PageFactory
	pages       map[string]interfaces.Page
}

func NewHandler(logger *logging.Logger, db *gorm.DB) handlers.Handler {
	pageFactory := pages.NewPageFactory(db)
	return &handler{
		logger:      logger,
		db:          db,
		pageFactory: pageFactory,
		pages:       make(map[string]interfaces.Page),
	}
}

func (h *handler) getPage(pageName string) interfaces.Page {
	if page, ok := h.pages[pageName]; ok {
		return page
	}

	newPage := h.pageFactory.CreatePage(pageName, h.logger)
	h.pages[pageName] = newPage
	return newPage
}

func (h *handler) Register(tgbotapi *tgbotapi.Update) (page interfaces.Page) {
	var pageName string
	if tgbotapi.CallbackQuery != nil {
		pageName = tgbotapi.CallbackQuery.Data
	} else {
		pageName = "main"
	}
	page = h.getPage(pageName)
	return page
}
