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

// habitLogCount is a helper struct for batch GROUP BY count queries.
type habitLogCount struct {
	HabitID uint
	Count   int64
}

// findDeclining finds active habits with no check-in in the last 3 days of the period.
// Uses a single GROUP BY query instead of per-habit count (N+1 fix).
// Before: 1 + N queries. After: 2 queries total.
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

	if len(habits) == 0 {
		return nil, nil
	}

	// Batch-fetch check-in counts for all habits in one query (N+1 fix)
	habitIDs := make([]uint, len(habits))
	for i, h := range habits {
		habitIDs[i] = h.ID
	}
	var counts []habitLogCount
	s.DB.Model(&models.HabitLog{}).
		Select("habit_id, COUNT(*) as count").
		Where("habit_id IN ? AND is_done = ? AND date >= ? AND date <= ?", habitIDs, true, threeDaysAgo, endDate).
		Group("habit_id").
		Scan(&counts)

	activeSet := make(map[uint]bool, len(counts))
	for _, c := range counts {
		if c.Count > 0 {
			activeSet[c.HabitID] = true
		}
	}

	var declining []models.Habit
	for _, habit := range habits {
		if !activeSet[habit.ID] {
			declining = append(declining, habit)
		}
	}

	return declining, nil
}

// findHighConsistency finds habits with >= 80% consistency within the period.
// Uses a single GROUP BY query instead of per-habit count (N+1 fix).
// Before: 1 + N queries. After: 2 queries total.
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

	if len(habits) == 0 {
		return nil, nil
	}

	// Batch-fetch done counts for all habits in one query (N+1 fix)
	habitIDs := make([]uint, len(habits))
	habitNames := make(map[uint]string, len(habits))
	for i, h := range habits {
		habitIDs[i] = h.ID
		habitNames[h.ID] = h.Name
	}
	var counts []habitLogCount
	s.DB.Model(&models.HabitLog{}).
		Select("habit_id, COUNT(*) as count").
		Where("habit_id IN ? AND is_done = ? AND date >= ? AND date <= ?",
			habitIDs, true, startDate, endDate).
		Group("habit_id").
		Scan(&counts)

	countMap := make(map[uint]int64, len(counts))
	for _, c := range counts {
		countMap[c.HabitID] = c.Count
	}

	var consistent []string
	for _, habit := range habits {
		if countMap[habit.ID] >= int64(threshold) {
			consistent = append(consistent, habit.Name)
		}
	}

	return consistent, nil
}
