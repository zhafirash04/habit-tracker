package services

import (
	"time"

	"habitflow/internal/models"

	"gorm.io/gorm"
)

// ReportService generates weekly reports.
type ReportService struct {
	DB             *gorm.DB
	ScoreService   *ScoreService
	InsightService *InsightService
}

// NewReportService creates a new ReportService.
func NewReportService(db *gorm.DB, scoreService *ScoreService, insightService *InsightService) *ReportService {
	return &ReportService{
		DB:             db,
		ScoreService:   scoreService,
		InsightService: insightService,
	}
}

// WeeklyReport holds the data for a weekly report.
type WeeklyReport struct {
	Period         string               `json:"period"`
	StartDate      string               `json:"start_date"`
	EndDate        string               `json:"end_date"`
	TotalHabits    int                  `json:"total_habits"`
	TotalCheckin   int                  `json:"total_checkins"`
	Score          *OverallScore        `json:"score"`
	Insights       []Insight            `json:"insights"`
	Streaks        []HabitStreakSummary `json:"streaks"`
	DailyBreakdown []DayBreakdown       `json:"daily_breakdown"`
}

// DayBreakdown holds completion data for a single day.
type DayBreakdown struct {
	Date      string  `json:"date"`
	DayName   string  `json:"day_name"`
	Completed int     `json:"completed"`
	Total     int     `json:"total"`
	Rate      float64 `json:"rate"`
}

// HabitStreakSummary holds streak info for a habit in the report.
type HabitStreakSummary struct {
	HabitID       uint   `json:"habit_id"`
	HabitName     string `json:"habit_name"`
	CurrentStreak int    `json:"current_streak"`
	LongestStreak int    `json:"longest_streak"`
}

// GenerateWeekly generates a weekly report for a user.
func (s *ReportService) GenerateWeekly(userID uint) (*WeeklyReport, error) {
	now := time.Now().In(WIB)
	// Find Monday of the current week
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7 // Sunday = 7
	}
	monday := now.AddDate(0, 0, -(weekday - 1))
	startStr := monday.Format("2006-01-02")
	endStr := now.Format("2006-01-02")
	return s.GenerateWeeklyForPeriod(userID, startStr, endStr)
}

// GenerateWeeklyForPeriod generates a weekly report for a specific date range.
func (s *ReportService) GenerateWeeklyForPeriod(userID uint, startStr, endStr string) (*WeeklyReport, error) {
	// Parse dates to calculate days
	startDate, err := time.Parse("2006-01-02", startStr)
	if err != nil {
		return nil, err
	}
	endDate, err := time.Parse("2006-01-02", endStr)
	if err != nil {
		return nil, err
	}
	days := int(endDate.Sub(startDate).Hours()/24) + 1
	if days < 1 {
		days = 7
	}

	// Count active habits
	var habitCount int64
	var habits []models.Habit
	s.DB.Where("user_id = ? AND is_active = ?", userID, true).Find(&habits)
	habitCount = int64(len(habits))

	// Count total check-ins this week (only for active habits)
	activeIDs := make([]uint, len(habits))
	for i, h := range habits {
		activeIDs[i] = h.ID
	}
	var checkinCount int64
	if len(activeIDs) > 0 {
		s.DB.Model(&models.HabitLog{}).
			Where("user_id = ? AND is_done = ? AND date >= ? AND date <= ? AND habit_id IN ?", userID, true, startStr, endStr, activeIDs).
			Count(&checkinCount)
	}

	// Get consistency score
	score, err := s.ScoreService.Calculate(userID, days)
	if err != nil {
		return nil, err
	}

	// Get insights
	insights, err := s.InsightService.Generate(userID, startStr, endStr)
	if err != nil {
		return nil, err
	}

	// Get streaks for all active habits
	var streaks []HabitStreakSummary
	for _, habit := range habits {
		var streak models.Streak
		if err := s.DB.Where("habit_id = ?", habit.ID).First(&streak).Error; err == nil {
			streaks = append(streaks, HabitStreakSummary{
				HabitID:       habit.ID,
				HabitName:     habit.Name,
				CurrentStreak: streak.CurrentStreak,
				LongestStreak: streak.LongestStreak,
			})
		}
	}

	// Build daily breakdown for the period.
	// Each day only counts habits that existed on that date (created_at <= end of day).
	var breakdown []DayBreakdown
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		ds := d.Format("2006-01-02")
		dayEnd := time.Date(d.Year(), d.Month(), d.Day(), 23, 59, 59, 0, d.Location())

		// Filter to habits that existed on this day
		var dayIDs []uint
		for _, h := range habits {
			if !h.CreatedAt.After(dayEnd) {
				dayIDs = append(dayIDs, h.ID)
			}
		}
		dayTotal := len(dayIDs)

		var cnt int64
		if dayTotal > 0 {
			s.DB.Model(&models.HabitLog{}).
				Where("user_id = ? AND date = ? AND is_done = ? AND habit_id IN ?", userID, ds, true, dayIDs).
				Count(&cnt)
		}
		rate := 0.0
		if dayTotal > 0 {
			rate = float64(cnt) / float64(dayTotal) * 100
		}
		breakdown = append(breakdown, DayBreakdown{
			Date:      ds,
			DayName:   dayNamesID[d.Weekday()],
			Completed: int(cnt),
			Total:     dayTotal,
			Rate:      rate,
		})
	}

	return &WeeklyReport{
		Period:         "weekly",
		StartDate:      startStr,
		EndDate:        endStr,
		TotalHabits:    int(habitCount),
		TotalCheckin:   int(checkinCount),
		Score:          score,
		Insights:       insights,
		Streaks:        streaks,
		DailyBreakdown: breakdown,
	}, nil
}
