package base

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
	"table10/internal/callbackdata"
	"table10/internal/models"
	"table10/pkg/logging"
)

type AbstractPage struct {
	Db           *gorm.DB
	Logger       *logging.Logger
	Ctx          context.Context
	User         *models.User
	Name         string
	Description  string
	Code         string
	KeyBoard     *tgbotapi.InlineKeyboardMarkup
	UserText     string
	CallbackData *callbackdata.CallbackData
}

func (bp *AbstractPage) GetName() string {
	return bp.Name
}

func (bp *AbstractPage) GetDescription() string {
	return bp.Description
}

func (bp *AbstractPage) GetCode() string {
	return bp.Code
}

func (bp *AbstractPage) GetKeyboard() *tgbotapi.InlineKeyboardMarkup {
	return bp.KeyBoard
}

func (bp *AbstractPage) GetUserText() string {
	return bp.UserText
}

func (bp *AbstractPage) Generate() {

}

func (bp *AbstractPage) SetUserText(text string) {
	bp.UserText = text
}
