package menu

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"table10/internal/handlers"
	"table10/internal/pages/cabinet"
	"table10/internal/pages/tasks"
	"table10/pkg/logging"
)

type handler struct {
	logger *logging.Logger
}

func NewHandler(logger *logging.Logger) handlers.Handler {
	return &handler{
		logger: logger,
	}
}

func (h *handler) Register(tgbotapi *tgbotapi.Update) (keyBoard *tgbotapi.InlineKeyboardMarkup) {
	page := tasks.NewPage()
	switch tgbotapi.CallbackQuery.Data {
	case "cabinet":
		page = cabinet.NewPage()
	case "tasks":
		page = tasks.NewPage()
	}
	return page.GetKeyboard(tgbotapi)
}
