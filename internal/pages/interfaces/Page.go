package interfaces

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"table10/internal/structs/telegram"
)

type Page interface {
	GetName() string
	GetDescription() string
	GetCode() string
	GetFullCode() string
	GetKeyboard() *tgbotapi.InlineKeyboardMarkup
	GetUserText() string                //Получить текст пользвоателя, который он написал на этой странице
	SetUserText(text string)            //Добавить текст пользвоателя, который он написал на этой странице
	GetUserPhoto() *telegram.Photo      //Получить фотографию пользвоателя, которую он загрузил на эту страницу
	SetUserPhoto(photo *telegram.Photo) //Добавить информацию о фотографии пользователя, которую он загрузил на эту страницу
	Generate()                          //Метод, который вызывается для каждой страницы для её генерации
}
