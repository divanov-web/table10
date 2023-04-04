package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"table10/internal/pages"
)

type Handler interface {
	Register(tgbotapi *tgbotapi.Update) pages.Page
}
