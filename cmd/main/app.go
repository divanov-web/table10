package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"table10/internal/config"
	"table10/internal/menu"
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

	bot, err := tgbotapi.NewBotAPI(cfg.Keys.Telegram)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	// Loop through each update.
	for update := range updates {
		// Check if we've gotten a message update.
		if update.Message != nil {
			// Construct a new message from the given chat ID and containing
			// the text that we received.
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

			// If the message was open, add a copy of our numeric keyboard.
			switch update.Message.Text {
			case "open":
				msg.ReplyMarkup = numericKeyboard
				msg.Text = "Меню"
			}

			// Send the message.
			if _, err = bot.Send(msg); err != nil {
				panic(err)
			}
		} else if update.CallbackQuery != nil {
			// Respond to the callback query, telling Telegram to show the user
			// a message with the data received.
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
			if _, err := bot.Request(callback); err != nil {
				panic(err)
			}

			// And finally, send a message containing the data received.
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
			menuHandler := menu.NewHandler(logger)
			msg.ReplyMarkup = menuHandler.Register(&update)
			if _, err := bot.Send(msg); err != nil {
				panic(err)
			}
		}
	}
}
