package menu

import (
	"context"
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
	"strings"
	"table10/internal/callbackdata"
	"table10/internal/constants"
	"table10/internal/handlers"
	"table10/internal/models"
	"table10/internal/pages"
	"table10/internal/pages/interfaces"
	"table10/pkg/logging"
)

type handler struct {
	logger      *logging.Logger
	db          *gorm.DB
	user        *models.User
	pageFactory *pages.PageFactory
	ctx         context.Context
	pages       map[string]interfaces.Page
}

func NewHandler(logger *logging.Logger, db *gorm.DB, user *models.User, ctx context.Context) handlers.Handler {
	pageFactory := pages.NewPageFactory(db)
	return &handler{
		logger:      logger,
		db:          db,
		user:        user,
		pageFactory: pageFactory,
		ctx:         ctx,
		pages:       make(map[string]interfaces.Page),
	}
}

func (h *handler) getPage(pageCode string, callbackdata *callbackdata.CallbackData) interfaces.Page {
	if page, ok := h.pages[pageCode]; ok {
		return page
	}

	newPage := h.pageFactory.CreatePage(pageCode, h.logger, h.user, h.ctx, callbackdata)
	h.pages[pageCode] = newPage
	return newPage
}

// processPageData разбивает строку с адресом на код страницы и её параметры
func (h *handler) processPageData(dataString string) (string, callbackdata.CallbackData) {
	dataParts := strings.Split(dataString, constants.ParamsSeparator)
	pageCode := dataParts[0]
	var callbackData callbackdata.CallbackData

	if len(dataParts) > 1 {
		jsonString := dataParts[1]
		h.logger.Infof("jsonString: %v", jsonString)
		err := json.Unmarshal([]byte(jsonString), &callbackData)
		if err != nil {
			// Обработка ошибки
		}
	}

	return pageCode, callbackData
}

// Register Регистрирует код страницы из сообщения пользователя
func (h *handler) Register(update *tgbotapi.Update) (page interfaces.Page) {
	var dataString string

	if update.CallbackQuery != nil {
		dataString = update.CallbackQuery.Data
	} else {
		dataString = h.user.LastPage
		if dataString == "" {
			dataString = "main"
		}
	}

	pageCode, callbackData := h.processPageData(dataString)
	h.logger.Infof("Current page = %v", pageCode)
	page = h.getPage(pageCode, &callbackData)
	return page
}
