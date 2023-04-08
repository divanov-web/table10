package pages

import (
	"gorm.io/gorm"
	"table10/internal/pages/cabinet"
	"table10/internal/pages/interfaces"
	mainpage "table10/internal/pages/main"
	"table10/internal/pages/tasks"
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

func (pf *PageFactory) CreatePage(pageName string, logger *logging.Logger) interfaces.Page {
	switch pageName {
	case "cabinet":
		return cabinet.NewPage(pf.db, logger)
	case "tasks":
		return tasks.NewPage(pf.db, logger)
	case "main":
		return mainpage.NewPage(pf.db, logger)
	default:
		return mainpage.NewPage(pf.db, logger)
	}
}
