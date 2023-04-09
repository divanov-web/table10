package main

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
	"log"
	"table10/internal/config"
	"table10/internal/menu"
	"table10/internal/models"
	"table10/internal/pages/interfaces"
	"table10/internal/repository"
	"table10/internal/services"
	"table10/pkg/logging"
	contextUtils "table10/pkg/utils/context"
	"time"
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
		userService := services.NewUserService(userRepo, logger)

		if update.Message != nil || update.CallbackQuery != nil {
			var userTelegram *tgbotapi.User
			var page interfaces.Page

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)

			if update.Message != nil {
				userTelegram = update.Message.From
			} else if update.CallbackQuery != nil {
				userTelegram = update.CallbackQuery.From
			}

			user := models.User{
				TelegramID:   int(userTelegram.ID),
				IsBot:        userTelegram.IsBot,
				FirstName:    userTelegram.FirstName,
				Username:     userTelegram.UserName,
				LanguageCode: userTelegram.LanguageCode,
			}
			var existingUser *models.User
			if existingUser, err = userService.GetUser(&user); err != nil {
				logger.Errorf("Can't find user: %v", err)
			}

			//Если пришло текстовое сообение
			if update.Message != nil {

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите пункт меню:")

				menuHandler := menu.NewHandler(logger, db, existingUser, ctx)
				page = menuHandler.Register(&update)
				page.SetUserText(update.Message.Text)
				page.Generate()
				if errContext := contextUtils.CheckContext(ctx); errContext != nil {
					msg.Text = "Произошел таймаут операции"
				} else {
					pageText := page.GetDescription()
					if pageText != "" {
						msg.Text = pageText
					}
					msg.ReplyMarkup = page.GetKeyboard()
				}
				msg.ParseMode = tgbotapi.ModeHTML

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

				menuHandler := menu.NewHandler(logger, db, existingUser, ctx)
				page = menuHandler.Register(&update)
				page.Generate()
				if errContext := contextUtils.CheckContext(ctx); errContext != nil {
					msg.Text = "Произошел таймаут операции"
				} else {
					msg.ReplyMarkup = page.GetKeyboard()
					msg.Text = page.GetDescription() + " (" + page.GetCommand() + ")"
				}
				msg.ParseMode = tgbotapi.ModeHTML

				if _, err := bot.Send(msg); err != nil {
					panic(err)
				}
			}

			user.LastPage = page.GetCommand()

			if err = userService.AddOrUpdateUser(&user); err != nil {
				logger.Errorf("Failed to add or update user: %v", err)
			}
			cancel()
		}
	}
}
