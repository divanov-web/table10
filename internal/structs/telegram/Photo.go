package telegram

// Photo Структура для хранения инфомрации о фотографии, полученной от телеграмм
type Photo struct {
	FileId   string //ID файла в телеграмм, можно получить путь bot.GetFileDirectURL(TelegramFileId)
	UniqueID string //Уникальный ID выданный телеграммом, будет именем файла
	Caption  string //Надпись к фотографии от пользователя
	Url      string //Url где хранится фотография
}
