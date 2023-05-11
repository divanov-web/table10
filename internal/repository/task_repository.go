package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"table10/internal/models"
	"time"
)

// TaskFilter структура содержащая в себе фильтры для выборки из базы
type TaskFilter struct {
	Current           bool         // Фильтр по датам, когда задание доступно для принятия
	Active            bool         // Фильтр по датам, когда задание доступно для сдачи
	User              *models.User // Фильтр по привязке пользователя
	NotAssignedToUser bool         // Флаг, определяющий, исключать ли задания с пользователем User
}

// UserTaskFilter фильтры для взятых заданий выборки из базы
type UserTaskFilter struct {
	PlayWithYou bool //Вместе с тобой играют
}

type TaskRepositoryInterface interface {
	GetTasks(ctx context.Context, game *models.Game, filter *TaskFilter) ([]models.Task, error)
	GetUserTasks(ctx context.Context, task *models.Task, filter *UserTaskFilter) ([]models.UserTask, error)
	GetOneById(ctx context.Context, id int, filter *TaskFilter) (*models.Task, error)
	AddUserToTask(ctx context.Context, user *models.User, task *models.Task, status *models.Status) error
	UpdateUserTaskStatus(ctx context.Context, userTask *models.UserTask, newStatus *models.Status) error
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepositoryInterface {
	return &taskRepository{
		db: db,
	}
}

func (r *taskRepository) GetOneById(ctx context.Context, id int, filter *TaskFilter) (*models.Task, error) {
	var existingTask models.Task

	query := r.db.WithContext(ctx).Where("tasks.id = ?", id).Preload("TaskType")

	if filter != nil && filter.User != nil {
		query = query.Joins("LEFT JOIN user_tasks ON user_tasks.task_id = tasks.id AND user_tasks.user_id = ?", filter.User.ID).
			Preload("UserTasks", "user_tasks.user_id = ?", filter.User.ID).
			Preload("UserTasks.User").
			Preload("UserTasks.Status")
	}

	if err := query.First(&existingTask).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(fmt.Sprintf("no task found id=%v", id))
		}
		return nil, err
	}

	return &existingTask, nil
}

func (r *taskRepository) GetTasks(ctx context.Context, game *models.Game, filter *TaskFilter) ([]models.Task, error) {
	var tasks []models.Task

	query := r.db.WithContext(ctx).
		Where("game_id = ?", game.ID).
		Preload("TaskType")

	if filter != nil {
		if filter.Current {
			now := time.Now()
			query = query.Where("start_date <= ? AND end_date >= ?", now, now)
		}

		if filter.Active {
			now := time.Now()
			query = query.Where("start_date <= ? AND close_date >= ?", now, now)
		}

		if filter.User != nil {
			if filter.NotAssignedToUser {
				query = query.Where("NOT EXISTS (SELECT 1 FROM user_tasks WHERE user_tasks.task_id = tasks.id AND user_tasks.user_id = ?)", filter.User.ID)
			} else {
				query = query.Joins("JOIN user_tasks ON user_tasks.task_id = tasks.id").
					Where("user_tasks.user_id = ?", filter.User.ID).
					Preload("UserTasks", "user_tasks.user_id = ?", filter.User.ID).
					Preload("UserTasks.User").
					Preload("UserTasks.Status")
			}
		}
	}

	err := query.Find(&tasks).Error

	if err != nil {
		return nil, err
	}

	return tasks, nil
}

// GetUserTasks Список взятых заданий
func (r *taskRepository) GetUserTasks(ctx context.Context, task *models.Task, filter *UserTaskFilter) ([]models.UserTask, error) {
	var userTasks []models.UserTask

	query := r.db.WithContext(ctx).
		Where("task_id = ?", task.ID).
		Preload("User")

	if filter != nil {
		if filter.PlayWithYou {
			query = query.Joins("JOIN statuses ON user_tasks.status_id = statuses.id").
				Where("statuses.code = ? AND user_tasks.user_id <> ? ", "in_progress", task.UserTasks[0].UserID)
		}
	}

	err := query.Find(&userTasks).Error

	if err != nil {
		return nil, err
	}

	return userTasks, nil
}

// AddUserToTask Добавляет пользователя к задаче. Принятие задачи пользователем.
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

func (r *taskRepository) UpdateUserTaskStatus(ctx context.Context, userTask *models.UserTask, newStatus *models.Status) error {
	userTask.Status = *newStatus

	result := r.db.Save(userTask)

	if result.Error != nil {
		return result.Error
	}

	return nil
}