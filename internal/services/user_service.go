package services

import (
	"context"
	"gorm.io/gorm"
	"table10/internal/models"
	"table10/internal/repository"
	"table10/pkg/logging"
	"time"
)

// UserService предоставляет методы для работы с пользователями
type UserService struct {
	repo   repository.UserRepositoryInterface
	logger *logging.Logger
}

// NewUserService создает и возвращает новый экземпляр UserService с заданным репозиторием и логгером.
func NewUserService(repo repository.UserRepositoryInterface, logger *logging.Logger) *UserService {
	return &UserService{
		repo:   repo,
		logger: logger,
	}
}

// AddOrUpdateUser добавляет нового пользователя или обновляет существующего в базе данных.
func (s *UserService) AddOrUpdateUser(user *models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	return s.repo.AddOrUpdateUser(ctx, user)
}

// GetUser возвращает существующего пользователя по его идентификатору или создает нового пользователя, если он не найден в базе данных.
func (s *UserService) GetUser(user *models.User) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	existingUser, err := s.repo.GetOneById(ctx, user)

	if err != nil {
		// Здесь мы обрабатываем только ошибку "пользователь не найден"
		if err == gorm.ErrRecordNotFound {
			// Создаем нового пользователя
			if err = s.repo.AddOrUpdateUser(ctx, user); err != nil {
				return nil, err
			}
			// Получаем созданного пользователя
			existingUser, err = s.repo.GetOneById(ctx, user)
			if err != nil {
				return nil, err
			}
		} else {
			// Обрабатываем другие ошибки, связанные с получением пользователя
			return nil, err
		}
	}

	return existingUser, nil
}
