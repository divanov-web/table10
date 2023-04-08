package base

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
)

type AbstractPage struct {
	db          *gorm.DB
	Name        string
	Description string
	Command     string
	KeyBoard    *tgbotapi.InlineKeyboardMarkup
	Text        string
}

func (bp *AbstractPage) GetName() string {
	return bp.Name
}

func (bp *AbstractPage) GetDescription() string {
	return bp.Description
}

func (bp *AbstractPage) GetCommand() string {
	return bp.Command
}

func (bp *AbstractPage) GetKeyboard() *tgbotapi.InlineKeyboardMarkup {
	return bp.KeyBoard
}

func (bp *AbstractPage) GetText() string {
	return bp.Text
}
