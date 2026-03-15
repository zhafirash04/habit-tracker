package services

import (
	"time"

	"habitflow/internal/models"

	"gorm.io/gorm"
)

// DailyReport holds the complete daily report data.
type DailyReport struct {
	Date                string          `json:"date"`
	DayName             string          `json:"day_name"`
	Summary             DailySummary    `json:"summary"`
	Habits              []DailyHabit    `json:"habits"`
	Streaks             StreakHighlight `json:"streaks"`
	ComparisonYesterday DailyComparison `json:"comparison_yesterday"`
	GeneratedAt         time.Time       `json:"generated_at"`
}

// DailySummary holds the daily summary counts.
type DailySummary struct {
	TotalHabits    int     `json:"total_habits"`
	Completed      int     `json:"completed"`
	Remaining      int     `json:"remaining"`
	CompletionRate float64 `json:"completion_rate"`
}

// DailyHabit holds info about a single habit for the daily report.
type DailyHabit struct {
	ID            uint       `json:"id"`
	Name          string     `json:"name"`
	Category      string     `json:"category"`
	IsDeletedNow  bool       `json:"is_deleted_now"`
	IsDone        bool       `json:"is_done"`
	CompletedAt   *time.Time `json:"completed_at"`
	Note          *string    `json:"note"`
	CurrentStreak int        `json:"current_streak"`
}

// StreakHighlight holds the best active streak info.
type StreakHighlight struct {
	LongestActive      *ActiveStreak `json:"longest_active"`
	TotalActiveStreaks int           `json:"total_active_streaks"`
}

// ActiveStreak holds info about a single active streak.
type ActiveStreak struct {
	HabitName string `json:"habit_name"`
	Days      int    `json:"days"`
}

// DailyComparison holds the comparison with yesterday.
type DailyComparison struct {
	YesterdayCompleted int    `json:"yesterday_completed"`
	TodayCompleted     int    `json:"today_completed"`
	Trend              string `json:"trend"`
	Message            string `json:"message"`
}

// DailyReportService generates daily reports.
type DailyReportService struct {
	DB *gorm.DB
}

// NewDailyReportService creates a new DailyReportService.
func NewDailyReportService(db *gorm.DB) *DailyReportService {
	return &DailyReportService{DB: db}
}

var dayNamesID = []string{"Minggu", "Senin", "Selasa", "Rabu", "Kamis", "Jumat", "Sabtu"}

// Generate produces a daily report for the given user and date.
func (s *DailyReportService) Generate(userID uint, date time.Time) (*DailyReport, error) {
	day := date.In(WIB)
	dayStart := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, WIB)
	dayEnd := dayStart.Add(24*time.Hour - time.Second)
	dateStr := dayStart.Format("2006-01-02")
	yesterday := dayStart.AddDate(0, 0, -1).Format("2006-01-02")

	// 1. Snapshot habits that existed by end-of-day (frozen report basis).
	// Do not filter by is_active so historical reports stay stable after soft delete.
	var habits []models.Habit
	if err := s.DB.Where("user_id = ? AND created_at <= ?", userID, dayEnd).
		Order("created_at ASC").Find(&habits).Error; err != nil {
		return nil, err
	}

	totalHabits := len(habits)

	// 2. Get habit logs for the target date
	var logs []models.HabitLog
	if err := s.DB.Where("user_id = ? AND date = ? AND is_done = ?", userID, dateStr, true).
		Find(&logs).Error; err != nil {
		return nil, err
	}
	logMap := make(map[uint]models.HabitLog, len(logs))
	for _, l := range logs {
		logMap[l.HabitID] = l
	}

	// 3. Get streaks for all habits
	habitIDs := make([]uint, 0, totalHabits)
	for _, h := range habits {
		habitIDs = append(habitIDs, h.ID)
	}
	var streaks []models.Streak
	if len(habitIDs) > 0 {
		s.DB.Where("habit_id IN ?", habitIDs).Find(&streaks)
	}
	streakMap := make(map[uint]models.Streak, len(streaks))
	for _, st := range streaks {
		streakMap[st.HabitID] = st
	}

	// 4. Build habit list and calculate summary
	completed := 0
	dailyHabits := make([]DailyHabit, 0, totalHabits)
	var bestStreak *ActiveStreak
	activeStreaks := 0

	for _, h := range habits {
		dh := DailyHabit{
			ID:           h.ID,
			Name:         h.Name,
			Category:     h.Category,
			IsDeletedNow: !h.IsActive,
		}
		if log, ok := logMap[h.ID]; ok {
			dh.IsDone = true
			dh.CompletedAt = log.CompletedAt
			dh.Note = log.Note
			completed++
		}
		if st, ok := streakMap[h.ID]; ok {
			dh.CurrentStreak = st.CurrentStreak
			if st.CurrentStreak > 0 {
				activeStreaks++
				if bestStreak == nil || st.CurrentStreak > bestStreak.Days {
					bestStreak = &ActiveStreak{HabitName: h.Name, Days: st.CurrentStreak}
				}
			}
		}
		dailyHabits = append(dailyHabits, dh)
	}

	remaining := totalHabits - completed
	var rate float64
	if totalHabits > 0 {
		rate = float64(completed) / float64(totalHabits) * 100
	}

	// 5. Get yesterday's data for comparison
	var yesterdayCount int64
	s.DB.Model(&models.HabitLog{}).
		Where("user_id = ? AND date = ? AND is_done = ?", userID, yesterday, true).
		Count(&yesterdayCount)

	trend := "same"
	var message string
	yc := int(yesterdayCount)
	if completed > yc {
		trend = "up"
		message = "Lebih baik dari kemarin! Pertahankan!"
	} else if completed < yc {
		trend = "down"
		if remaining > 0 {
			message = "Kemarin kamu menyelesaikan " + itoa(yc) + " habit. Ayo selesaikan " + itoa(remaining) + " habit lagi!"
		} else {
			message = "Kemarin kamu menyelesaikan " + itoa(yc) + " habit."
		}
	} else {
		if completed == totalHabits && totalHabits > 0 {
			message = "Sama dengan kemarin — semua selesai! Konsisten!"
		} else {
			message = "Sama dengan kemarin. Yuk tingkatkan!"
		}
	}

	report := &DailyReport{
		Date:    dateStr,
		DayName: dayNamesID[dayStart.Weekday()],
		Summary: DailySummary{
			TotalHabits:    totalHabits,
			Completed:      completed,
			Remaining:      remaining,
			CompletionRate: rate,
		},
		Habits: dailyHabits,
		Streaks: StreakHighlight{
			LongestActive:      bestStreak,
			TotalActiveStreaks: activeStreaks,
		},
		ComparisonYesterday: DailyComparison{
			YesterdayCompleted: yc,
			TodayCompleted:     completed,
			Trend:              trend,
			Message:            message,
		},
		GeneratedAt: time.Now().In(WIB),
	}

	return report, nil
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	s := ""
	neg := n < 0
	if neg {
		n = -n
	}
	for n > 0 {
		s = string(rune('0'+n%10)) + s
		n /= 10
	}
	if neg {
		s = "-" + s
	}
	return s
}
