package repository

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"table10/internal/models"
)

type UserGameInfo struct {
	Game models.Game
	Role models.Role
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) AddOrUpdateUser(ctx context.Context, user *models.User) error {
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

func (r *UserRepository) GetOneById(ctx context.Context, user *models.User) (*models.User, error) {
	var existingUser models.User

	if err := r.db.WithContext(ctx).Where("telegram_id = ?", user.TelegramID).First(&existingUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("no user found")
		}
		return nil, err
	}

	return &existingUser, nil
}

func (r *UserRepository) AddUserToGame(ctx context.Context, user *models.User, game *models.Game) error {
	userGame := &models.UserGame{
		UserID: user.ID,
		GameID: game.ID,
	}

	return r.db.WithContext(ctx).Create(userGame).Error
}

func (r *UserRepository) GetUserGames(ctx context.Context, user *models.User) ([]UserGameInfo, error) {
	var userGames []models.UserGame
	err := r.db.WithContext(ctx).Preload("Game").Where("user_id = ?", user.ID).Find(&userGames).Error
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
