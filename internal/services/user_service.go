package services

import (
	"context"
	"table10/internal/models"
	"table10/internal/repository"
	"time"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) AddOrUpdateUser(user *models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	return s.repo.AddOrUpdateUser(ctx, user)
}
