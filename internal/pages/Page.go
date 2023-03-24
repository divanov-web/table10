package pages

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Page interface {
	GetKeyboard(tgbotapi *tgbotapi.Update) *tgbotapi.InlineKeyboardMarkup
}
