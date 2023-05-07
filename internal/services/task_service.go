package services

import (
	"context"
	"table10/internal/models"
	"table10/internal/repository"
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

func (s *TaskService) GetTasks(period *models.Period) ([]models.Task, error) {
	return s.taskRepo.GetTasks(s.ctx, period)
}

func (s *TaskService) GetOneById(id int) (*models.Task, error) {
	return s.taskRepo.GetOneById(s.ctx, id)
}

// AddUserToTask Добавляет юзера в выбранное задание
func (s *TaskService) AddUserToTask(task *models.Task, user *models.User) error {
	defaultStatus, err := s.statusRepo.GetOne(s.ctx, "new")
	if err != nil {
		return err
	}

	err1 := s.taskRepo.AddUserToTask(s.ctx, user, task, defaultStatus)
	if err1 != nil {
		return err1
	}

	return nil
}
