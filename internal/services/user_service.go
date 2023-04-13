package services

import (
	"context"
	"table10/internal/models"
	"table10/internal/repository"
	"table10/pkg/logging"
	"time"
)

type UserService struct {
	repo   repository.UserRepositoryInterface
	logger *logging.Logger
}

func NewUserService(repo repository.UserRepositoryInterface, logger *logging.Logger) *UserService {
	return &UserService{
		repo:   repo,
		logger: logger,
	}
}

func (s *UserService) AddOrUpdateUser(user *models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	return s.repo.AddOrUpdateUser(ctx, user)
}

func (s *UserService) GetUser(user *models.User) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var existingUser *models.User
	var err error

	existingUser, err = s.repo.GetOneById(ctx, user)
	if err != nil {
		existingUser = user
	}

	return existingUser, nil
}
