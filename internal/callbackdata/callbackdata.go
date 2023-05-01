package callbackdata

import (
	"errors"
	"strconv"
)

// CallbackData Это структура для параметров страницы, которые передаются вместе с кодом страницы в одной строке
type CallbackData struct {
	Params map[string]string `json:"params"`
}

// GetTaskId Метод для получения id из параметров CallbackData
func (cd *CallbackData) GetTaskId() (int, error) {
	taskIdString, ok := cd.Params["id"]
	if !ok {
		return 0, errors.New("ошибка обработки задания: отсутствует параметр 'id'")
	}
	taskId, err := strconv.Atoi(taskIdString)
	if err != nil {
		return 0, errors.New("ошибка обработки задания: неверный формат 'id'")
	}
	return taskId, nil
}
