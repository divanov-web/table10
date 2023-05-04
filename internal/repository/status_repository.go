package repository

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"table10/internal/models"
)

type StatusRepositoryInterface interface {
	GetOne(ctx context.Context, code string) (*models.Status, error)
}

type statusRepository struct {
	db *gorm.DB
}

func NewStatusRepository(db *gorm.DB) StatusRepositoryInterface {
	return &statusRepository{
		db: db,
	}
}

func (r *statusRepository) GetOne(ctx context.Context, code string) (*models.Status, error) {
	var existingStatus models.Status

	if err := r.db.WithContext(ctx).Where("code = ?", code).First(&existingStatus).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("no Status found")
		}
		return nil, err
	}

	return &existingStatus, nil
}
