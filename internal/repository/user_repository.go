package repository

import (
	"context"
	"errors"
	"github.com/lib/pq" // Для обработки ошибок PostgreSQL
	"gorm.io/gorm"
	"table10/internal/models"
)

type UserGameInfo struct {
	Game models.Game
	Role models.Role
}

type UserRepositoryInterface interface {
	AddOrUpdateUser(ctx context.Context, user *models.User) error
	GetOneById(ctx context.Context, user *models.User) (*models.User, error)
	AddUserToGameWithRole(ctx context.Context, user *models.User, game *models.Game, role *models.Role) error
	GetUserGames(ctx context.Context, user *models.User) ([]UserGameInfo, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepositoryInterface {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) AddOrUpdateUser(ctx context.Context, user *models.User) error {
	var existingUser models.User
	result := r.db.WithContext(ctx).Where("telegram_id = ?", user.TelegramID).First(&existingUser)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return r.db.WithContext(ctx).Create(user).Error
		}
		return result.Error
	}

	user.ID = existingUser.ID
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *userRepository) GetOneById(ctx context.Context, user *models.User) (*models.User, error) {
	var existingUser models.User

	if err := r.db.WithContext(ctx).
		Preload("Games.Game").
		Preload("Games.Role").
		Joins("LEFT JOIN user_games ON user_games.user_id = users.id").
		Joins("LEFT JOIN games ON games.id = user_games.game_id AND user_games.is_main = ?", true).
		Where("telegram_id = ?", user.TelegramID).
		First(&existingUser).Error; err != nil {

		return nil, err
	}

	return &existingUser, nil
}

func (r *userRepository) AddUserToGameWithRole(ctx context.Context, user *models.User, game *models.Game, role *models.Role) error {
	userGame := &models.UserGame{
		UserID: user.ID,
		GameID: game.ID,
		RoleID: role.ID,
		IsMain: true,
	}

	result := r.db.WithContext(ctx).Create(userGame)
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

func (r *userRepository) GetUserGames(ctx context.Context, user *models.User) ([]UserGameInfo, error) {
	var userGames []models.UserGame
	err := r.db.WithContext(ctx).Preload("Game").Preload("Role").Where("user_id = ?", user.ID).Find(&userGames).Error
	if err != nil {
		return nil, err
	}

	gameInfos := make([]UserGameInfo, len(userGames))
	for i, userGame := range userGames {
		gameInfos[i] = UserGameInfo{
			Game: userGame.Game,
			Role: userGame.Role,
		}
	}

	return gameInfos, nil
}
