package seed

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"table10/internal/models"
	"table10/internal/repository"
	"table10/pkg/logging"
	"table10/pkg/utils/formtating"
)

func AddTask(db *gorm.DB, logger *logging.Logger, ctx context.Context) error {
	var taskType models.TaskType
	typeCode := "common"
	if err := db.Where("code = ?", typeCode).First(&taskType).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Errorf("failed to find game tasktype: %v", typeCode)
		}
		logger.Errorf("failed to find tasktype code: %v, unknown error", typeCode)
	}

	periodRepo := repository.NewPeriodRepository(db)
	currentPeriod, err := periodRepo.ShowCurrent(ctx)
	if err != nil {
		logger.Errorf("Текущая неделя не найдена в миграции")
	}

	// Создаем список тестовых заданий
	itemsToAdd := []models.Task{
		{
			Name:             fmt.Sprintf("Задание 1 для недели %v %v", currentPeriod.WeekNumber, currentPeriod.StartDate.Format("02.01.2006")),
			PeriodID:         currentPeriod.ID,
			TaskTypeID:       taskType.ID,
			ShortDescription: formtating.StrPtr("Короткое описание задания 1"),
			LongDescription:  formtating.StrPtr("Полное описание задания 1"),
		},
		{
			Name:             fmt.Sprintf("Задание 2 для недели %v %v", currentPeriod.WeekNumber, currentPeriod.StartDate.Format("02.01.2006")),
			PeriodID:         currentPeriod.ID,
			TaskTypeID:       taskType.ID,
			ShortDescription: formtating.StrPtr("Короткое описание задания 2"),
			LongDescription:  formtating.StrPtr("Полное описание задания 2"),
		},
	}

	for _, itemToAdd := range itemsToAdd {
		var item models.Task
		if err := db.Where("name = ? AND period_id = ?", itemToAdd.Name, currentPeriod.ID).First(&item).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				// Если запись не существует, создаем новую запись и сохраняем ее в таблице
				logger.Infof("Задание %s для недели %d не существует , добавляем", itemToAdd.Name, currentPeriod.WeekNumber)

				if err = db.Create(&itemToAdd).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		} else {
			if err = db.Model(&item).Updates(&itemToAdd).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
