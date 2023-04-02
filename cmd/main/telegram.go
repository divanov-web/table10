package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
	"log"
	"table10/internal/config"
	"table10/internal/menu"
	"table10/internal/models"
	"table10/internal/repository"
	"table10/internal/services"
	"table10/pkg/logging"
)

var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Задания", "task"),
		tgbotapi.NewInlineKeyboardButtonData("Личный кабинет", "cabinet"),
	),
)

func telegramStart(cfg *config.Config, logger *logging.Logger, db *gorm.DB) {
	bot, err := tgbotapi.NewBotAPI(cfg.Keys.Telegram)
	if err != nil {
		log.Panic(err)
	}

	//bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	// Создаем UserRepository
	userRepo := repository.NewUserRepository(db)

	// Loop through each update.
	for update := range updates {
		// Check if we've gotten a message update.
		if update.Message != nil {
			// Log user information
			logger.Infof("User: ID: %d, FirstName: %s, UserName: %s, ChatID: %d, Text: %s",
				update.Message.From.ID,
				update.Message.From.FirstName,
				update.Message.From.UserName,
				update.Message.Chat.ID,
				update.Message.Text,
			)

			// Создаем и сохраняем пользователя в базе данных
			user := models.User{
				TelegramID:   int(update.Message.From.ID),
				IsBot:        update.Message.From.IsBot,
				FirstName:    update.Message.From.FirstName,
				Username:     update.Message.From.UserName,
				LanguageCode: update.Message.From.LanguageCode,
			}
			userService := services.NewUserService(userRepo)
			if err = userService.AddOrUpdateUser(&user); err != nil {
				logger.Errorf("Failed to add or update user: %v", err)
			}

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
