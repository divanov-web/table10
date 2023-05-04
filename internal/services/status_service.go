package services

import (
	"context"
	"table10/internal/models"
	"table10/internal/repository"
	"table10/pkg/logging"
)

type StatusService struct {
	statusRepo repository.StatusRepositoryInterface
	logger     *logging.Logger
	ctx        context.Context
}

func NewStatusService(
	statusRepo repository.StatusRepositoryInterface,
	logger *logging.Logger,
	ctx context.Context) *StatusService {
	return &StatusService{
		statusRepo: statusRepo,
		logger:     logger,
		ctx:        ctx,
	}
}

func (s *StatusService) GetOneByCode(code string) (*models.Status, error) {
	return s.statusRepo.GetOne(s.ctx, code)
}
