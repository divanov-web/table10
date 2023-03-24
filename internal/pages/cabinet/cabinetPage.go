package cabinet

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"table10/internal/pages"
)

type page struct {
	keyBoard *tgbotapi.InlineKeyboardMarkup
}

func NewPage() pages.Page {
	numericKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Счёт", "score"),
			tgbotapi.NewInlineKeyboardButtonData("Настройки", "config"),
		),
	)
	return &page{
		keyBoard: &numericKeyboard,
	}
}

func (p *page) GetKeyboard(tgbotapi *tgbotapi.Update) (keyboard *tgbotapi.InlineKeyboardMarkup) {
	return p.keyBoard
}
