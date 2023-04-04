package models

import (
	"gorm.io/gorm"
	"table10/pkg/logging"
	"time"
)

type WeekPeriod struct {
	StartDate time.Time
	EndDate   time.Time
}

func generateWeekPeriods(startDate, endDate time.Time) []WeekPeriod {
	weekPeriods := []WeekPeriod{}

	// Нормализация startDate до начала недели (понедельника)
	for startDate.Weekday() != time.Monday {
		startDate = startDate.AddDate(0, 0, -1)
	}
	startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())

	for {
		weekEndDate := startDate.AddDate(0, 0, 6)
		if weekEndDate.After(endDate) {
			break
		}
		weekEndDate = time.Date(weekEndDate.Year(), weekEndDate.Month(), weekEndDate.Day(), 23, 59, 59, 0, weekEndDate.Location())

		weekPeriods = append(weekPeriods, WeekPeriod{
			StartDate: startDate,
			EndDate:   weekEndDate,
		})

		startDate = startDate.AddDate(0, 0, 7)
	}

	return weekPeriods
}

func SeedPeriods(db *gorm.DB, logger *logging.Logger) error {
	// Определите текущую дату и дату через год
	now := time.Now()
	yearFromNow := now.AddDate(1, 0, 0)

	// Генерируем недельные периоды
	weekPeriods := generateWeekPeriods(now, yearFromNow)

	// Добавляем недельные периоды в таблицу Period
	for _, period := range weekPeriods {
		// Получить номер недели
		_, weekNumber := period.StartDate.ISOWeek()

		newPeriod := &Period{
			GameID:     2,
			WeekNumber: weekNumber,
			StartDate:  period.StartDate,
			EndDate:    period.EndDate,
		}

		//logger.Infof("WeekNumber: %v, EndDate: %v", weekNumber, period.EndDate)

		var existingPeriod Period
		err := db.Where("game_id = ? AND week_number = ?", newPeriod.GameID, newPeriod.WeekNumber).First(&existingPeriod).Error

		if err != nil {
			if err == gorm.ErrRecordNotFound {
				// Запись не найдена, создаем новую
				if err := db.Save(newPeriod).Error; err != nil {
					logger.Errorf("failed to create new period: %v", err)
				}
			} else {
				logger.Errorf("failed to query existing period: %v", err)
			}
		} else {
			// Запись найдена, обновляем ее
			existingPeriod.StartDate = newPeriod.StartDate
			existingPeriod.EndDate = newPeriod.EndDate
			if err := db.Save(&existingPeriod).Error; err != nil {
				logger.Errorf("failed to update existing period: %v", err)
			}
		}
	}

	return nil
}
