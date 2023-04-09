package utils

import (
	"encoding/json"
	"table10/internal/callbackdata"
)

func CreateCallbackDataJSON(params map[string]string) ([]byte, error) {
	callbackData := callbackdata.CallbackData{
		Params: params,
	}

	callbackDataJSON, err := json.Marshal(callbackData)
	if err != nil {
		return nil, err
	}

	return callbackDataJSON, nil
}
