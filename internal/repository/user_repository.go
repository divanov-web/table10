package repository

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"table10/internal/models"
)

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
			return r.db.Create(user).Error
		}
		return result.Error
	}

	user.ID = existingUser.ID
	return r.db.Save(user).Error
}
