package base

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
	"table10/internal/callbackdata"
	"table10/internal/constants"
	"table10/internal/models"
	"table10/internal/structs/telegram"
	"table10/pkg/logging"
	"table10/pkg/utils"
)

type AbstractPage struct {
	Db           *gorm.DB
	Logger       *logging.Logger
	Ctx          context.Context
	User         *models.User
	Name         string                         //Имя страницы, вроде нигде не выводится
	Description  string                         //Выводимый текст
	Messages     []telegram.Message             //Дополнительные выводимые сообщения
	Code         string                         //код страницы, возможно тоже нигде не выводится
	KeyBoard     *tgbotapi.InlineKeyboardMarkup //Выводимые пункты меню
	UserText     string                         //Текст, полученный от пользователя
	UserPhoto    *telegram.Photo                //Несохранённые фотографии от пользователя
	CallbackData *callbackdata.CallbackData     //Параметры страницы в json
}

func (bp *AbstractPage) GetName() string {
	return bp.Name
}

func (bp *AbstractPage) GetDescription() string {
	return bp.Description
}

func (bp *AbstractPage) GetMessages() []telegram.Message {
	return bp.Messages
}

func (bp *AbstractPage) GetCode() string {
	return bp.Code
}

// GetFullCode Формирует полный адрес страницы с учётом параметров
func (bp *AbstractPage) GetFullCode() string {
	if bp.CallbackData == nil || len(bp.CallbackData.Params) == 0 {
		return bp.Code
	}

	callbackDataJSON, err := utils.CreateCallbackDataJSON(bp.CallbackData.Params)
	if err != nil {
		// Обработка ошибки
	}

	return bp.Code + constants.ParamsSeparator + string(callbackDataJSON)
}

func (bp *AbstractPage) GetKeyboard() *tgbotapi.InlineKeyboardMarkup {
	return bp.KeyBoard
}

func (bp *AbstractPage) GetUserText() string {
	return bp.UserText
}

func (bp *AbstractPage) SetUserText(text string) {
	bp.UserText = text
}

func (bp *AbstractPage) GetUserPhoto() *telegram.Photo {
	return bp.UserPhoto
}

func (bp *AbstractPage) SetUserPhoto(photo *telegram.Photo) {
	bp.UserPhoto = photo
}

func (bp *AbstractPage) Generate() {

}

// CanModerate Может ли пользователь модерировать эту страницу
func (bp *AbstractPage) CanModerate() bool {
	return bp.User.Games[0].Role.CanModerate()
}
