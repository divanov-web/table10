package services

import (
	"context"
	"table10/internal/models"
	"table10/internal/repository"
	"table10/pkg/logging"
)

type TaskService struct {
	taskRepo repository.TaskRepositoryInterface
	logger   *logging.Logger
	ctx      context.Context
}

func NewTaskService(
	taskRepo repository.TaskRepositoryInterface,
	logger *logging.Logger,
	ctx context.Context) *TaskService {
	return &TaskService{
		taskRepo: taskRepo,
		logger:   logger,
		ctx:      ctx,
	}
}

func (s *TaskService) GetTasks(period *models.Period) ([]models.Task, error) {
	return s.taskRepo.GetTasks(s.ctx, period)
}

func (s *TaskService) GetOneById(id int) (*models.Task, error) {
	return s.taskRepo.GetOneById(s.ctx, id)
}
