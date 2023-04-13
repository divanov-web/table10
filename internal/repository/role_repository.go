package repository

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"table10/internal/models"
)

type RoleRepositoryInterface interface {
	GetOne(ctx context.Context, code string) (*models.Role, error)
}

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepositoryInterface {
	return &roleRepository{
		db: db,
	}
}

func (r *roleRepository) GetOne(ctx context.Context, code string) (*models.Role, error) {
	var existingRole models.Role

	if err := r.db.WithContext(ctx).Where("code = ?", code).First(&existingRole).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("no Role found")
		}
		return nil, err
	}

	return &existingRole, nil
}
