package repository

import (
	"context"
	"gorm.io/gorm"
	"table10/internal/models"
)

// AnswerFilter структура содержащая в себе фильтры для выборки из базы
type AnswerFilter struct {
	UserTask *models.UserTask // Фильтр по заданию
}

type AnswerRepositoryInterface interface {
	GetAnswers(ctx context.Context, filter *AnswerFilter) ([]models.Answer, error)
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

func (r *answerRepository) GetAnswers(ctx context.Context, filter *AnswerFilter) ([]models.Answer, error) {
	var answers []models.Answer

	query := r.db.WithContext(ctx)

	if filter != nil {
		if filter.UserTask != nil {
			query = query.Where("task_id = ?", filter.UserTask.TaskID)
		}
	}

	err := query.Find(&answers).Error

	if err != nil {
		return nil, err
	}

	return answers, nil
}
