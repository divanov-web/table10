package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"table10/internal/config"
	"table10/pkg/logging"
)

var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Задания", "task"),
		tgbotapi.NewInlineKeyboardButtonData("Личный кабинет", "cabinet"),
	),
)

func main() {
	logger := logging.GetLogger()
	logger.Info("create telegram connection")

	cfg := config.GetConfig()

	telegramStart(cfg, logger)
}
