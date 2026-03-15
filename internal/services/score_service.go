package services

import (
	"math"
	"time"

	"habitflow/internal/models"

	"gorm.io/gorm"
)

// ScoreService calculates Habit Consistency Scores.
type ScoreService struct {
	DB *gorm.DB
}

// NewScoreService creates a new ScoreService.
func NewScoreService(db *gorm.DB) *ScoreService {
	return &ScoreService{DB: db}
}

// HabitScore holds the consistency score for a single habit.
type HabitScore struct {
	HabitID   uint    `json:"habit_id"`
	HabitName string  `json:"habit_name"`
	Score     float64 `json:"score"`      // 0-100
	Done      int     `json:"done"`       // days completed
	Total     int     `json:"total_days"` // total days in range
}

// OverallScore holds the overall consistency score across all habits.
type OverallScore struct {
	Overall     float64      `json:"overall_score"` // average across all habits
	HabitScores []HabitScore `json:"habit_scores"`
	Period      string       `json:"period"` // e.g. "last_7_days"
}

// Calculate computes the consistency score for all user habits over a date range.
func (s *ScoreService) Calculate(userID uint, days int) (*OverallScore, error) {
	var habits []models.Habit
	if err := s.DB.Where("user_id = ? AND is_active = ?", userID, true).Find(&habits).Error; err != nil {
		return nil, err
	}

	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days)
	startStr := startDate.Format("2006-01-02")
	endStr := endDate.Format("2006-01-02")

	var scores []HabitScore
	var totalScore float64

	for _, habit := range habits {
		// Count how many days the habit was done in the range
		var doneCount int64
		s.DB.Model(&models.HabitLog{}).
			Where("habit_id = ? AND user_id = ? AND is_done = ? AND date >= ? AND date <= ?",
				habit.ID, userID, true, startStr, endStr).
			Count(&doneCount)

		// Calculate total possible days (from habit creation or start of range, whichever is later)
		habitStart := habit.CreatedAt
		if habitStart.Before(startDate) {
			habitStart = startDate
		}
		// Truncate to calendar dates to avoid timestamp-based miscounting
		habitStartDay := time.Date(habitStart.Year(), habitStart.Month(), habitStart.Day(), 0, 0, 0, 0, habitStart.Location())
		endDay := time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 0, 0, 0, 0, endDate.Location())
		totalDays := int(endDay.Sub(habitStartDay).Hours()/24) + 1
		if totalDays <= 0 {
			totalDays = 1
		}
		if totalDays > days {
			totalDays = days
		}

		score := float64(doneCount) / float64(totalDays) * 100
		if score > 100 {
			score = 100
		}

		scores = append(scores, HabitScore{
			HabitID:   habit.ID,
			HabitName: habit.Name,
			Score:     math.Round(score*100) / 100,
			Done:      int(doneCount),
			Total:     totalDays,
		})

		totalScore += score
	}

	overall := 0.0
	if len(scores) > 0 {
		overall = math.Round(totalScore/float64(len(scores))*100) / 100
	}

	period := "last_7_days"
	if days == 30 {
		period = "last_30_days"
	} else if days != 7 {
		period = time.Now().AddDate(0, 0, -days).Format("2006-01-02") + "_to_" + time.Now().Format("2006-01-02")
	}

	return &OverallScore{
		Overall:     overall,
		HabitScores: scores,
		Period:      period,
	}, nil
}
