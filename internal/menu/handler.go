package menu

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
	"table10/internal/handlers"
	"table10/internal/pages"
	"table10/internal/pages/cabinet"
	mainpage "table10/internal/pages/main"
	"table10/internal/pages/tasks"
	"table10/pkg/logging"
)

type handler struct {
	logger *logging.Logger
	db     *gorm.DB
	pages  map[string]pages.Page
}

func NewHandler(logger *logging.Logger, db *gorm.DB) handlers.Handler {
	return &handler{
		logger: logger,
		db:     db,
		pages:  make(map[string]pages.Page),
	}
}

func (h *handler) getPage(pageName string) pages.Page {
	if page, ok := h.pages[pageName]; ok {
		return page
	}

	availablePages := []pages.Page{
		cabinet.NewPage(),
		tasks.NewPage(),
		mainpage.NewPage(),
	}

	var newPage pages.Page
	for _, page := range availablePages {
		if page.GetCommand() == pageName {
			newPage = page
			break
		}
	}

	if newPage == nil {
		newPage = mainpage.NewPage()
	}

	h.pages[pageName] = newPage
	return newPage
}

func (h *handler) Register(tgbotapi *tgbotapi.Update) (page pages.Page) {
	var pageName string
	if tgbotapi.CallbackQuery != nil {
		pageName = tgbotapi.CallbackQuery.Data
	} else {
		pageName = "main"
	}
	page = h.getPage(pageName)
	return page
}
