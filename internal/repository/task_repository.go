package repository

import (
	"context"
	"gorm.io/gorm"
	"table10/internal/models"
)

type TaskRepositoryInterface interface {
	GetTasks(ctx context.Context, period *models.Period) ([]models.Task, error)
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepositoryInterface {
	return &taskRepository{
		db: db,
	}
}

func (r *taskRepository) GetTasks(ctx context.Context, period *models.Period) ([]models.Task, error) {
	var tasks []models.Task

	err := r.db.WithContext(ctx).
		Where("period_id = ?", period.ID).
		Preload("TaskType").
		Find(&tasks).Error

	if err != nil {
		return nil, err
	}

	return tasks, nil
}
