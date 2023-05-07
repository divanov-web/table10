package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"table10/internal/models"
)

type TaskRepositoryInterface interface {
	GetTasks(ctx context.Context, period *models.Period) ([]models.Task, error)
	GetOneById(ctx context.Context, id int) (*models.Task, error)
	AddUserToTask(ctx context.Context, user *models.User, task *models.Task, status *models.Status) error
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

func (r *taskRepository) AddUserToTask(ctx context.Context, user *models.User, task *models.Task, status *models.Status) error {
	userTask := &models.UserTask{
		UserID:   user.ID,
		TaskID:   task.ID,
		StatusID: status.ID,
	}

	result := r.db.WithContext(ctx).Create(userTask)
	if result.Error != nil {
		// Проверяем, является ли ошибка ошибкой дублирования ключа
		var pqErr *pq.Error
		if errors.As(result.Error, &pqErr) && pqErr.Code.Name() == "unique_violation" {
			return errors.New("duplicated key not allowed")
		}
		return result.Error
	}

	return nil
}
