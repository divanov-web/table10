package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"gorm.io/gorm"
	StatusCode "table10/internal/constants/statusCode"
	"table10/internal/models"
	"time"
)

// TaskFilter структура содержащая в себе фильтры для выборки из базы
type TaskFilter struct {
	Current           bool         // Фильтр по датам, когда задание доступно для принятия
	Available         bool         // Фильтр по датам, когда задание доступно для сдачи
	User              *models.User // Фильтр по привязке пользователя
	NotAssignedToUser bool         // Флаг, определяющий, исключать ли задания с пользователем User
	Limit             int          // Ограничить выборку
	IsActive          *bool        // Стоит флаг "Ативно"
}

// UserTaskFilter фильтры для взятых заданий выборки из базы
type UserTaskFilter struct {
	PlayWithYou bool         //Вместе с тобой играют
	StatusCode  string       //Фильтр по статусу
	Task        *models.Task //Фильтр по заданию
	GameId      uint         //Фильтр по играм
}

type TaskRepositoryInterface interface {
	GetTasks(ctx context.Context, game *models.Game, filter *TaskFilter) ([]models.Task, error)
	GetUserTasks(ctx context.Context, filter *UserTaskFilter) ([]models.UserTask, error)
	GetOneById(ctx context.Context, id int, filter *TaskFilter) (*models.Task, error)
	GetUserTaskById(ctx context.Context, id int) (*models.UserTask, error)
	AddUserToTask(ctx context.Context, user *models.User, task *models.Task, status *models.Status) error
	UpdateUserTaskStatus(ctx context.Context, userTask *models.UserTask, newStatus *models.Status) error
	ChangeActive(ctx context.Context, task *models.Task, isActive bool) error
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

func (r *taskRepository) GetUserTaskById(ctx context.Context, id int) (*models.UserTask, error) {
	var existingUserTask models.UserTask

	query := r.db.WithContext(ctx).Where("user_tasks.id = ?", id).
		Preload("Task").
		Preload("User").
		Preload("Status")

	if err := query.First(&existingUserTask).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(fmt.Sprintf("no userTask found id=%v", id))
		}
		return nil, err
	}

	return &existingUserTask, nil
}

func (r *taskRepository) GetTasks(ctx context.Context, game *models.Game, filter *TaskFilter) ([]models.Task, error) {
	var tasks []models.Task

	query := r.db.WithContext(ctx).
		Joins("LEFT JOIN user_tasks ON user_tasks.task_id = tasks.id").
		Where("game_id = ?", game.ID).
		Preload("TaskType")

	if filter != nil {
		if filter.Current {
			now := time.Now()
			query = query.Where("start_date <= ? AND end_date >= ?", now, now)
		}

		if filter.IsActive != nil {
			query = query.Where("is_active = ?", filter.IsActive)
		}

		if filter.Available {
			now := time.Now()
			query = query.
				Joins("JOIN statuses ON user_tasks.status_id = statuses.id").
				Where("start_date <= ? AND close_date >= ?", now, now).
				Where("statuses.code = ?", StatusCode.InProgress)
		}

		if filter.User != nil {
			if filter.NotAssignedToUser {
				query = query.Where("NOT EXISTS (SELECT 1 FROM user_tasks WHERE user_tasks.task_id = tasks.id AND user_tasks.user_id = ?)", filter.User.ID)
			} else {
				query = query.
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
func (r *taskRepository) GetUserTasks(ctx context.Context, filter *UserTaskFilter) ([]models.UserTask, error) {
	var userTasks []models.UserTask

	query := r.db.WithContext(ctx).
		Joins("JOIN tasks ON user_tasks.task_id = tasks.id").
		Preload("User")

	if filter != nil {
		if filter.Task != nil {
			query = query.Where("task_id = ?", filter.Task.ID)
		}

		if filter.GameId != 0 {
			query = query.Where("tasks.game_id = ?", filter.GameId)
		}

		if filter.PlayWithYou {
			query = query.Joins("JOIN statuses ON user_tasks.status_id = statuses.id").
				Where("statuses.code = ? AND user_tasks.user_id <> ? ", StatusCode.InProgress, filter.Task.UserTasks[0].UserID)
		}

		if filter.StatusCode != "" {
			query = query.Joins("JOIN statuses ON user_tasks.status_id = statuses.id").
				Where("statuses.code = ?", filter.StatusCode).
				Preload("Task")
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

	result := r.db.WithContext(ctx).Save(userTask)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

// ChangeActive change active status for task
func (r *taskRepository) ChangeActive(ctx context.Context, task *models.Task, isActive bool) error {
	task.IsActive = isActive

	result := r.db.WithContext(ctx).Save(task)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
