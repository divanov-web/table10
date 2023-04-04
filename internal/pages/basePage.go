package pages

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type BasePage struct {
	Name        string
	Description string
	Command     string
	KeyBoard    *tgbotapi.InlineKeyboardMarkup
}

func (bp *BasePage) GetName() string {
	return bp.Name
}

func (bp *BasePage) GetDescription() string {
	return bp.Description
}

func (bp *BasePage) GetCommand() string {
	return bp.Command
}

func (bp *BasePage) GetKeyboard() *tgbotapi.InlineKeyboardMarkup {
	return bp.KeyBoard
}
