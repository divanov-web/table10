package handlers

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Handler interface {
	Register(tgbotapi *tgbotapi.Update) *tgbotapi.InlineKeyboardMarkup
}
