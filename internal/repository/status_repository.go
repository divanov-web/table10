package repository

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"table10/internal/models"
)

type StatusRepositoryInterface interface {
	GetOne(ctx context.Context, code string) (*models.Status, error)
	GetNextStatus(ctx context.Context, currentStatus *models.Status) (*models.Status, error)
	GetPreviousStatus(ctx context.Context, currentStatus *models.Status) (*models.Status, error)
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

func (r *statusRepository) GetNextStatus(ctx context.Context, currentStatus *models.Status) (*models.Status, error) {
	var nextStatus models.Status

	switch currentStatus.Code {
	case "in_progress":
		err := r.db.WithContext(ctx).Where("code = ?", "under_review").First(&nextStatus).Error
		if err != nil {
			return nil, err
		}
	case "under_review":
		err := r.db.WithContext(ctx).Where("code = ?", "accepted").First(&nextStatus).Error
		if err != nil {
			return nil, err
		}
	case "accepted":
		return nil, errors.New("cannot move to the next status from accepted")
	default:
		return nil, errors.New("unsupported status code")
	}

	return &nextStatus, nil
}

func (r *statusRepository) GetPreviousStatus(ctx context.Context, currentStatus *models.Status) (*models.Status, error) {
	var previousStatus models.Status

	switch currentStatus.Code {
	case "in_progress":
		return nil, errors.New("cannot move to the previous status from in_progress")
	case "under_review":
		err := r.db.WithContext(ctx).Where("code = ?", "in_progress").First(&previousStatus).Error
		if err != nil {
			return nil, err
		}
	case "accepted":
		err := r.db.WithContext(ctx).Where("code = ?", "under_review").First(&previousStatus).Error
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("unsupported status code")
	}

	return &previousStatus, nil
}
