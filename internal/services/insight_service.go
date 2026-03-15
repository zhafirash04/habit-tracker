package services

import (
	"fmt"
	"time"

	"habitflow/internal/models"

	"gorm.io/gorm"
)

// InsightService detects patterns in habit data.
type InsightService struct {
	DB *gorm.DB
}

// NewInsightService creates a new InsightService.
func NewInsightService(db *gorm.DB) *InsightService {
	return &InsightService{DB: db}
}

// Insight represents a single pattern insight.
type Insight struct {
	Type    string `json:"type"`    // e.g. "best_day", "declining", "consistency"
	Message string `json:"message"` // human-readable insight
	HabitID uint   `json:"habit_id,omitempty"`
}

// Generate produces pattern insights for a user's habits within a specific period.
func (s *InsightService) Generate(userID uint, startDate, endDate string) ([]Insight, error) {
	var insights []Insight

	// 1. Best day of the week (within the report period)
	bestDay, err := s.findBestDay(userID, startDate, endDate)
	if err == nil && bestDay != "" {
		insights = append(insights, Insight{
			Type:    "best_day",
			Message: fmt.Sprintf("Hari paling produktifmu adalah %s! Pertahankan!", bestDay),
		})
	}

	// 2. Declining habits (active but no check-in in last 3 days of the period)
	declining, err := s.findDeclining(userID, endDate)
	if err == nil {
		for _, h := range declining {
			insights = append(insights, Insight{
				Type:    "declining",
				Message: fmt.Sprintf("Kebiasaan '%s' belum dilakukan 3 hari terakhir. Yuk mulai lagi!", h.Name),
				HabitID: h.ID,
			})
		}
	}

	// 3. High consistency habits within the period
	consistent, err := s.findHighConsistency(userID, startDate, endDate)
	if err == nil {
		for _, name := range consistent {
			insights = append(insights, Insight{
				Type:    "consistency",
				Message: fmt.Sprintf("Kamu sangat konsisten dengan '%s' minggu ini! Keren!", name),
			})
		}
	}

	if len(insights) == 0 {
		insights = append(insights, Insight{
			Type:    "encouragement",
			Message: "Mulai check-in kebiasaanmu hari ini untuk melihat insight-mu!",
		})
	}

	return insights, nil
}

// findBestDay finds the day of the week with the most completed habits within the period.
func (s *InsightService) findBestDay(userID uint, startDate, endDate string) (string, error) {
	var logs []models.HabitLog
	if err := s.DB.Where("user_id = ? AND is_done = ? AND date >= ? AND date <= ?", userID, true, startDate, endDate).
		Find(&logs).Error; err != nil {
		return "", err
	}

	if len(logs) == 0 {
		return "", nil
	}

	dayCounts := make(map[time.Weekday]int)
	for _, log := range logs {
		date, err := time.Parse("2006-01-02", log.Date)
		if err != nil {
			continue
		}
		dayCounts[date.Weekday()]++
	}

	var bestDay time.Weekday
	maxCount := 0
	for day, count := range dayCounts {
		if count > maxCount {
			maxCount = count
			bestDay = day
		}
	}

	dayNames := map[time.Weekday]string{
		time.Sunday:    "Minggu",
		time.Monday:    "Senin",
		time.Tuesday:   "Selasa",
		time.Wednesday: "Rabu",
		time.Thursday:  "Kamis",
		time.Friday:    "Jumat",
		time.Saturday:  "Sabtu",
	}

	return dayNames[bestDay], nil
}

// findDeclining finds active habits with no check-in in the last 3 days of the period.
func (s *InsightService) findDeclining(userID uint, endDate string) ([]models.Habit, error) {
	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, err
	}
	threeDaysAgo := end.AddDate(0, 0, -3).Format("2006-01-02")

	var habits []models.Habit
	if err := s.DB.Where("user_id = ? AND is_active = ?", userID, true).Find(&habits).Error; err != nil {
		return nil, err
	}

	var declining []models.Habit
	for _, habit := range habits {
		var count int64
		s.DB.Model(&models.HabitLog{}).
			Where("habit_id = ? AND is_done = ? AND date >= ? AND date <= ?", habit.ID, true, threeDaysAgo, endDate).
			Count(&count)

		if count == 0 {
			declining = append(declining, habit)
		}
	}

	return declining, nil
}

// findHighConsistency finds habits with >= 80% consistency within the period.
func (s *InsightService) findHighConsistency(userID uint, startDate, endDate string) ([]string, error) {
	start, _ := time.Parse("2006-01-02", startDate)
	end, _ := time.Parse("2006-01-02", endDate)
	totalDays := int(end.Sub(start).Hours()/24) + 1
	threshold := int(float64(totalDays) * 0.8)
	if threshold < 1 {
		threshold = 1
	}

	var habits []models.Habit
	if err := s.DB.Where("user_id = ? AND is_active = ?", userID, true).Find(&habits).Error; err != nil {
		return nil, err
	}

	var consistent []string
	for _, habit := range habits {
		var doneCount int64
		s.DB.Model(&models.HabitLog{}).
			Where("habit_id = ? AND is_done = ? AND date >= ? AND date <= ?",
				habit.ID, true, startDate, endDate).
			Count(&doneCount)

		if doneCount >= int64(threshold) {
			consistent = append(consistent, habit.Name)
		}
	}

	return consistent, nil
}
