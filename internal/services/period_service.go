package services

import (
	"context"
	"table10/internal/models"
	"table10/internal/repository"
	"table10/pkg/logging"
)

type PeriodService struct {
	periodRepo repository.PeriodRepositoryInterface
	logger     *logging.Logger
	ctx        context.Context
}

func NewPeriodService(
	periodRepo repository.PeriodRepositoryInterface,
	logger *logging.Logger,
	ctx context.Context) *PeriodService {
	return &PeriodService{
		periodRepo: periodRepo,
		logger:     logger,
		ctx:        ctx,
	}
}

func (s *PeriodService) ShowCurrent() (*models.Period, error) {
	return s.periodRepo.ShowCurrent(s.ctx)
}
