package services

import (
	"context"
	"fmt"
	"path/filepath"
	"table10/internal/config"
	"table10/internal/models"
	"table10/internal/repository"
	"table10/internal/structs/telegram"
	"table10/pkg/logging"
	"table10/pkg/utils/file"
)

type AnswerService struct {
	answerRepo repository.AnswerRepositoryInterface
	logger     *logging.Logger
	ctx        context.Context
}

func NewAnswerService(
	answerRepo repository.AnswerRepositoryInterface,
	logger *logging.Logger,
	ctx context.Context) *AnswerService {
	return &AnswerService{
		answerRepo: answerRepo,
		logger:     logger,
		ctx:        ctx,
	}
}

// AddAnswer Добавляет ответ пользователя
func (s *AnswerService) AddAnswer(userText string, userPhoto *telegram.Photo, user *models.User, task *models.Task) error {
	answer := models.Answer{}
	if userText != "" {
		answer.Text = userText
	}
	if userPhoto != nil {
		imagePath, err := s.CopyFile(userPhoto, task)
		if err != nil {
			return err
		}
		answer.ImagePath = imagePath
		answer.TelegramFileId = userPhoto.FileId
		if userText != "" {
			answer.Text = userText
		} else {
			answer.Text = userPhoto.Caption
		}
	}
	return s.answerRepo.AddAnswer(s.ctx, &answer, user, &task.UserTasks[0])
}

func (s *AnswerService) GetAnswers(filter *repository.AnswerFilter) ([]models.Answer, error) {
	return s.answerRepo.GetAnswers(s.ctx, filter)
}

// CopyFile Метод копирует фото, присланные пользователем
func (s *AnswerService) CopyFile(userPhoto *telegram.Photo, task *models.Task) (string, error) {
	cfg := config.GetConfig()
	uploadPath := cfg.Storage.UploadPath

	// Получение расширения файла из URL
	fileExtension := filepath.Ext(userPhoto.Url)

	// Формирование пути сохранения файла
	relSavePath := filepath.Join("answers", fmt.Sprintf("task_id_%d", task.ID), fmt.Sprintf("user_task_id_%d", task.UserTasks[0].ID), userPhoto.UniqueID+fileExtension)

	// Скачивание и сохранение файла
	err := file.DownloadAndSaveFile(userPhoto.Url, filepath.Join(uploadPath, relSavePath))
	if err != nil {
		return "", err
	}

	return relSavePath, nil
}
