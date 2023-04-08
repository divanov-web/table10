package repository

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"table10/internal/models"
)

type GameRepository struct {
	db *gorm.DB
}

func NewGameRepository(db *gorm.DB) *GameRepository {
	return &GameRepository{
		db: db,
	}
}

func (r *GameRepository) GetOne(ctx context.Context, text string) (*models.Game, error) {
	var existingGame models.Game

	if err := r.db.WithContext(ctx).Where("code = ?", text).First(&existingGame).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("no game found")
		}
		return nil, err
	}

	return &existingGame, nil
}
