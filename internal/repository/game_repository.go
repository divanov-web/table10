package repository

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"table10/internal/models"
)

type GameRepositoryInterface interface {
	GetOne(ctx context.Context, text string) (*models.Game, error)
	GetOneById(ctx context.Context, id int) (*models.Game, error)
	AddUserToGame(ctx context.Context, user *models.User, game *models.Game) error
}

type gameRepository struct {
	db *gorm.DB
}

func NewGameRepository(db *gorm.DB) GameRepositoryInterface {
	return &gameRepository{
		db: db,
	}
}

func (r *gameRepository) GetOne(ctx context.Context, text string) (*models.Game, error) {
	var existingGame models.Game

	if err := r.db.WithContext(ctx).Where("code = ?", text).First(&existingGame).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("no game found")
		}
		return nil, err
	}

	return &existingGame, nil
}

func (r *gameRepository) GetOneById(ctx context.Context, id int) (*models.Game, error) {
	var existingGame models.Game

	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&existingGame).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("no game found")
		}
		return nil, err
	}

	return &existingGame, nil
}

func (r *gameRepository) AddUserToGame(ctx context.Context, user *models.User, game *models.Game) error {
	userGame := &models.UserGame{
		UserID: user.ID,
		GameID: game.ID,
	}

	return r.db.WithContext(ctx).Create(userGame).Error
}
