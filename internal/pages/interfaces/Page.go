package interfaces

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Page interface {
	GetName() string
	GetDescription() string
	GetCode() string
	GetKeyboard() *tgbotapi.InlineKeyboardMarkup
	GetUserText() string
	SetUserText(text string)
	Generate()
}
