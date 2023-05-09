package repository

import (
	"context"
	"gorm.io/gorm"
	"table10/internal/models"
)

type AnswerRepositoryInterface interface {
	AddAnswer(ctx context.Context, answer *models.Answer, user *models.User, task *models.Task) error
}

type answerRepository struct {
	db *gorm.DB
}

func NewAnswerRepository(db *gorm.DB) AnswerRepositoryInterface {
	return &answerRepository{
		db: db,
	}
}

// AddAnswer Добавляет ответ пользователя
func (r *answerRepository) AddAnswer(ctx context.Context, answer *models.Answer, user *models.User, task *models.Task) error {
	answer.TaskID = task.ID
	answer.UserID = user.ID
	result := r.db.WithContext(ctx).Create(answer)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
