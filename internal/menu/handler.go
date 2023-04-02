package menu

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"table10/internal/handlers"
	"table10/internal/pages"
	"table10/internal/pages/cabinet"
	"table10/internal/pages/tasks"
	"table10/pkg/logging"
)

type handler struct {
	logger *logging.Logger
	pages  map[string]pages.Page
}

func NewHandler(logger *logging.Logger) handlers.Handler {
	return &handler{
		logger: logger,
		pages:  make(map[string]pages.Page),
	}
}

func (h *handler) getPage(pageName string) pages.Page {
	if page, ok := h.pages[pageName]; ok {
		return page
	}

	var newPage pages.Page
	switch pageName {
	case "cabinet":
		newPage = cabinet.NewPage()
	case "tasks":
		newPage = tasks.NewPage()
	default:
		newPage = tasks.NewPage()
	}

	h.pages[pageName] = newPage
	return newPage
}

func (h *handler) Register(tgbotapi *tgbotapi.Update) (keyBoard *tgbotapi.InlineKeyboardMarkup) {
	pageName := tgbotapi.CallbackQuery.Data
	page := h.getPage(pageName)
	return page.GetKeyboard(tgbotapi)
}
