package pages

import (
	"context"
	"gorm.io/gorm"
	"table10/internal/callbackdata"
	"table10/internal/constants/pageCode"
	"table10/internal/models"
	adminPage "table10/internal/pages/admin"
	adminReviewPage "table10/internal/pages/adminReview"
	adminTaskDetailPage "table10/internal/pages/adminTaskDetail"
	adminTasksPage "table10/internal/pages/adminTasks"
	adminUserTaskDetailPage "table10/internal/pages/adminUserTaskDetail"
	"table10/internal/pages/cabinet"
	gamePage "table10/internal/pages/game"
	gameAcceptPage "table10/internal/pages/gameAccept"
	gameInputPage "table10/internal/pages/gameInput"
	"table10/internal/pages/interfaces"
	mainpage "table10/internal/pages/main"
	taskDetailPage "table10/internal/pages/taskDetail"
	tasksPage "table10/internal/pages/tasks"
	tasksAcceptedPage "table10/internal/pages/tasksAccepted"
	welcomePage "table10/internal/pages/welcome"
	"table10/pkg/logging"
)

type PageFactory struct {
	db *gorm.DB
}

func NewPageFactory(db *gorm.DB) *PageFactory {
	return &PageFactory{
		db: db,
	}
}

// CreatePage Выбор класса страницы в зависимости от кода страницы
func (pf *PageFactory) CreatePage(pageName string, logger *logging.Logger, user *models.User, ctx context.Context, callbackdata *callbackdata.CallbackData) interfaces.Page {

	var page interfaces.Page
	switch pageName {
	case pageCode.Main:
		page = mainpage.NewPage(pf.db, logger, ctx, user)
	case pageCode.Welcome:
		page = welcomePage.NewPage(pf.db, logger, ctx, user)
	case pageCode.Cabinet:
		page = cabinetPage.NewPage(pf.db, logger, ctx, user)
	case pageCode.Game:
		page = gamePage.NewPage(pf.db, logger, ctx, user)
	case pageCode.GameInput:
		page = gameInputPage.NewPage(pf.db, logger, ctx, user)
	case pageCode.GameAccept:
		page = gameAcceptPage.NewPage(pf.db, logger, ctx, user, callbackdata)
	case pageCode.Tasks:
		page = tasksPage.NewPage(pf.db, logger, ctx, user, callbackdata)
	case pageCode.TasksAccepted:
		page = tasksAcceptedPage.NewPage(pf.db, logger, ctx, user, callbackdata)
	case pageCode.TaskDetail:
		page = taskDetailPage.NewPage(pf.db, logger, ctx, user, callbackdata)
	//admin pages
	case pageCode.Admin:
		page = adminPage.NewPage(pf.db, logger, ctx, user)
	case pageCode.AdminReview:
		page = adminReviewPage.NewPage(pf.db, logger, ctx, user, callbackdata)
	case pageCode.AdminUserTaskDetail:
		page = adminUserTaskDetailPage.NewPage(pf.db, logger, ctx, user, callbackdata)
	case pageCode.AdminTasks:
		page = adminTasksPage.NewPage(pf.db, logger, ctx, user, callbackdata)
	case pageCode.AdminTaskDetail:
		page = adminTaskDetailPage.NewPage(pf.db, logger, ctx, user, callbackdata)
	default:
		page = mainpage.NewPage(pf.db, logger, ctx, user)
	}

	return page
}
