package repository

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"table10/internal/models"
	"time"
)

type PeriodRepository struct {
	db *gorm.DB
}

func NewPeriodRepository(db *gorm.DB) *PeriodRepository {
	return &PeriodRepository{
		db: db,
	}
}

type PeriodRepositoryInterface interface {
	ShowCurrent(ctx context.Context) (*models.Period, error)
}

func (r *PeriodRepository) ShowCurrent(ctx context.Context) (*models.Period, error) {
	var currentPeriod models.Period
	now := time.Now()

	if err := r.db.WithContext(ctx).Where("start_date <= ? AND end_date >= ?", now, now).First(&currentPeriod).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("no active period found")
		}
		return nil, err
	}

	return &currentPeriod, nil
}
