package repository

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"table10/internal/models"
)

type TaskRepositoryInterface interface {
	GetTasks(ctx context.Context, period *models.Period) ([]models.Task, error)
	GetOneById(ctx context.Context, id int) (*models.Task, error)
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepositoryInterface {
	return &taskRepository{
		db: db,
	}
}

func (r *taskRepository) GetOneById(ctx context.Context, id int) (*models.Task, error) {
	var existingTask models.Task

	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&existingTask).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(fmt.Sprintf("no task found id=%v", id))
		}
		return nil, err
	}

	return &existingTask, nil
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
