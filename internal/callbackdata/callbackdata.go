package callbackdata

import (
	"errors"
	"strconv"
)

// CallbackData Это структура для параметров страницы, которые передаются вместе с кодом страницы в одной строке
type CallbackData struct {
	Params map[string]string `json:"params"`
}

// GetId Метод для получения id из параметров CallbackData
func (cd *CallbackData) GetId() (int, error) {
	idString, ok := cd.Params["id"]
	if !ok {
		return 0, errors.New("ошибка обработки задания: отсутствует параметр 'id'")
	}
	id, err := strconv.Atoi(idString)
	if err != nil {
		return 0, errors.New("ошибка обработки задания: неверный формат 'id'")
	}
	return id, nil
}

// GetAction Метод для получения action из параметров CallbackData
func (cd *CallbackData) GetAction() string {
	actionString, ok := cd.Params["action"]
	if !ok {
		return "default"
	}
	return actionString
}
