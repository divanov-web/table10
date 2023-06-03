package services

import (
	"context"
	"errors"
	"table10/internal/models"
	"table10/internal/repository"
	"table10/internal/services/task_straregy"
	"table10/pkg/logging"
)

type TaskService struct {
	taskRepo   repository.TaskRepositoryInterface
	userRepo   repository.UserRepositoryInterface
	statusRepo repository.StatusRepositoryInterface
	logger     *logging.Logger
	ctx        context.Context
}

func NewTaskService(
	taskRepo repository.TaskRepositoryInterface,
	userRepo repository.UserRepositoryInterface,
	statusRepo repository.StatusRepositoryInterface,
	logger *logging.Logger,
	ctx context.Context) *TaskService {
	return &TaskService{
		taskRepo:   taskRepo,
		userRepo:   userRepo,
		statusRepo: statusRepo,
		logger:     logger,
		ctx:        ctx,
	}
}

// GetTasks список заданий
func (s *TaskService) GetTasks(game *models.Game, filter *repository.TaskFilter) ([]models.Task, error) {
	return s.taskRepo.GetTasks(s.ctx, game, filter)
}

// GetUserTasks список взятых заданий
func (s *TaskService) GetUserTasks(filter *repository.UserTaskFilter) ([]models.UserTask, error) {
	return s.taskRepo.GetUserTasks(s.ctx, filter)
}

// GetOneById одно задание с фильтром
func (s *TaskService) GetOneById(id int, filter *repository.TaskFilter) (*models.Task, task_straregy.TaskProgressionStrategy, error) {
	task, err := s.taskRepo.GetOneById(s.ctx, id, filter)
	if err != nil {
		return nil, nil, err
	}
	taskStrategy, err1 := s.GetTaskProgressionStrategy(task)
	if err1 != nil {
		return nil, nil, err1
	}
	return task, taskStrategy, nil
}

// GetUserTaskById get userTask by id
func (s *TaskService) GetUserTaskById(id int) (*models.UserTask, error) {
	return s.taskRepo.GetUserTaskById(s.ctx, id)
}

// AddUserToTask Добавляет юзера в выбранное задание
func (s *TaskService) AddUserToTask(task *models.Task, user *models.User, taskStrategy task_straregy.TaskProgressionStrategy) error {
	statusCode, err := taskStrategy.GetFirstStatusCode()
	if err != nil {
		return err
	}
	defaultStatus, err := s.statusRepo.GetOne(s.ctx, statusCode)
	if err != nil {
		return err
	}

	err1 := s.taskRepo.AddUserToTask(s.ctx, user, task, defaultStatus)
	if err1 != nil {
		return err1
	}

	return nil
}

// ChangeStatus Change status to userTask
func (s *TaskService) ChangeStatus(userTask *models.UserTask, statusCode string) error {
	newStatus, err := s.statusRepo.GetOne(s.ctx, statusCode)
	if err != nil {
		return err
	}

	err = s.taskRepo.UpdateUserTaskStatus(s.ctx, userTask, newStatus)
	if err != nil {
		return err
	}
	return nil
}

// ChangeActive change active status for task
func (s *TaskService) ChangeActive(task *models.Task, isActive bool) error {
	err := s.taskRepo.ChangeActive(s.ctx, task, isActive)
	if err != nil {
		return err
	}
	return nil
}

// GetTaskProgressionStrategy Status strategy with different task types
func (s *TaskService) GetTaskProgressionStrategy(task *models.Task) (task_straregy.TaskProgressionStrategy, error) {
	switch task.TaskType.Code {
	case "common":
		return &task_straregy.CommonTaskProgression{}, nil
	case "buddy":
		return &task_straregy.BuddyTaskProgression{}, nil
	case "solo":
		return &task_straregy.SoloTaskProgression{}, nil
	default:
		return nil, errors.New("unsupported task type")
	}
}
