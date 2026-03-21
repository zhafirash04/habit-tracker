package services

import (
	"time"

	"habitflow/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// WIB is the Asia/Jakarta timezone (UTC+7).
var WIB = time.FixedZone("WIB", 7*3600)

// TodayWIB returns today's date string in WIB timezone.
func TodayWIB() string {
	return time.Now().In(WIB).Format("2006-01-02")
}

// YesterdayWIB returns yesterday's date string in WIB timezone.
func YesterdayWIB() string {
	return time.Now().In(WIB).AddDate(0, 0, -1).Format("2006-01-02")
}

// StreakService handles streak calculation logic.
type StreakService struct {
	DB *gorm.DB
}

// NewStreakService creates a new StreakService.
func NewStreakService(db *gorm.DB) *StreakService {
	return &StreakService{DB: db}
}

// CheckinResult holds the flat response data after a check-in operation.
type CheckinResult struct {
	HabitID       uint       `json:"habit_id"`
	Date          string     `json:"date"`
	IsDone        bool       `json:"is_done"`
	Note          *string    `json:"note"`
	CurrentStreak int        `json:"current_streak"`
	LongestStreak int        `json:"longest_streak"`
	CompletedAt   *time.Time `json:"completed_at"`
}

// Checkin performs the full check-in logic inside a database transaction:
// creates/updates HabitLog, updates Streak, and returns the flat result.
func (s *StreakService) Checkin(userID, habitID uint, note *string) (*CheckinResult, error) {
	today := TodayWIB()
	yesterday := YesterdayWIB()
	var result CheckinResult

	err := s.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Verify habit belongs to user and is active
		var habit models.Habit
		if err := tx.Where("id = ? AND user_id = ? AND is_active = ?", habitID, userID, true).
			First(&habit).Error; err != nil {
			return gorm.ErrRecordNotFound
		}

		// 2. Check if already checked in today
		var existingLog models.HabitLog
		logExists := tx.Where("habit_id = ? AND date = ?", habitID, today).
			First(&existingLog).Error == nil

		if logExists && existingLog.IsDone {
			return ErrAlreadyCheckedIn
		}

		// 3. Create or update HabitLog
		now := time.Now().In(WIB)
		if logExists {
			existingLog.IsDone = true
			existingLog.CompletedAt = &now
			existingLog.Note = note
			if err := tx.Save(&existingLog).Error; err != nil {
				return err
			}
		} else {
			existingLog = models.HabitLog{
				HabitID:     habitID,
				UserID:      userID,
				Date:        today,
				IsDone:      true,
				Note:        note,
				CompletedAt: &now,
			}
			if err := tx.Create(&existingLog).Error; err != nil {
				return err
			}
		}

		// 4. Lock and update streak (SELECT FOR UPDATE pattern via clause.Locking)
		var streak models.Streak
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("habit_id = ?", habitID).First(&streak).Error; err != nil {
			// Create if not exists
			streak = models.Streak{
				HabitID:       habitID,
				CurrentStreak: 0,
				LongestStreak: 0,
			}
			tx.Create(&streak)
		}

		// 5. Calculate new streak
		if streak.LastDoneDate != nil {
			switch *streak.LastDoneDate {
			case yesterday:
				// Consecutive day — continue streak
				streak.CurrentStreak++
			case today:
				// Same day — already counted, skip
			default:
				// Gap — start new streak
				streak.CurrentStreak = 1
			}
		} else {
			// First-ever check-in
			streak.CurrentStreak = 1
		}

		// 6. Update longest streak (never decreases)
		if streak.CurrentStreak > streak.LongestStreak {
			streak.LongestStreak = streak.CurrentStreak
		}

		streak.LastDoneDate = &today
		if err := tx.Save(&streak).Error; err != nil {
			return err
		}

		// 7. Build result
		result = CheckinResult{
			HabitID:       habitID,
			Date:          today,
			IsDone:        true,
			Note:          note,
			CurrentStreak: streak.CurrentStreak,
			LongestStreak: streak.LongestStreak,
			CompletedAt:   &now,
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return &result, nil
}

// UndoCheckin removes today's check-in and recalculates the streak in a transaction.
func (s *StreakService) UndoCheckin(userID, habitID uint) error {
	today := TodayWIB()

	return s.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Find and delete today's log
		result := tx.Where("habit_id = ? AND user_id = ? AND date = ? AND is_done = ?",
			habitID, userID, today, true).Delete(&models.HabitLog{})
		if result.RowsAffected == 0 {
			return ErrNoCheckinToday
		}

		// 2. Lock streak for update
		var streak models.Streak
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("habit_id = ?", habitID).First(&streak).Error; err != nil {
			return err
		}

		// 3. Recalculate streak from remaining logs
		var logs []models.HabitLog
		tx.Where("habit_id = ? AND is_done = ?", habitID, true).
			Order("date DESC").
			Find(&logs)

		if len(logs) == 0 {
			streak.CurrentStreak = 0
			streak.LastDoneDate = nil
		} else {
			streak.LastDoneDate = &logs[0].Date

			// Count consecutive days backwards from the most recent log
			currentStreak := 1
			for i := 1; i < len(logs); i++ {
				curr, _ := time.Parse("2006-01-02", logs[i-1].Date)
				prev, _ := time.Parse("2006-01-02", logs[i].Date)
				diff := curr.Sub(prev).Hours() / 24

				if diff == 1 {
					currentStreak++
				} else {
					break
				}
			}
			streak.CurrentStreak = currentStreak
		}
		// longest_streak never decreases

		return tx.Save(&streak).Error
	})
}

// GetStreak returns the streak for a habit.
func (s *StreakService) GetStreak(habitID uint) (*models.Streak, error) {
	var streak models.Streak
	if err := s.DB.Where("habit_id = ?", habitID).First(&streak).Error; err != nil {
		return nil, err
	}
	return &streak, nil
}

// TodayStatus holds a habit's check-in status for today.
type TodayStatus struct {
	HabitID       uint    `json:"habit_id"`
	Name          string  `json:"name"`
	Category      string  `json:"category"`
	IsDoneToday   bool    `json:"is_done_today"`
	CurrentStreak int     `json:"current_streak"`
	Note          *string `json:"note"`
}

// GetTodayStatus returns all active habits for a user with today's check-in status.
// Uses batch queries for logs and streaks to avoid N+1 problem.
// Before: 1 + 2N queries (1 habits + N logs + N streaks). After: 3 queries total.
func (s *StreakService) GetTodayStatus(userID uint) ([]TodayStatus, error) {
	today := TodayWIB()

	var habits []models.Habit
	if err := s.DB.Where("user_id = ? AND is_active = ?", userID, true).
		Order("created_at ASC").
		Find(&habits).Error; err != nil {
		return nil, err
	}

	if len(habits) == 0 {
		return []TodayStatus{}, nil
	}

	// Collect habit IDs for batch queries
	habitIDs := make([]uint, len(habits))
	for i, h := range habits {
		habitIDs[i] = h.ID
	}

	// Batch-fetch today's logs for all habits in one query (N+1 fix)
	var logs []models.HabitLog
	s.DB.Where("habit_id IN ? AND date = ? AND is_done = ?", habitIDs, today, true).Find(&logs)

	logMap := make(map[uint]models.HabitLog, len(logs))
	for _, log := range logs {
		logMap[log.HabitID] = log
	}

	// Batch-fetch all streaks in one query (N+1 fix)
	var streaks []models.Streak
	s.DB.Where("habit_id IN ?", habitIDs).Find(&streaks)

	streakMap := make(map[uint]models.Streak, len(streaks))
	for _, st := range streaks {
		streakMap[st.HabitID] = st
	}

	// Build statuses using pre-fetched maps — no more per-habit queries
	statuses := make([]TodayStatus, 0, len(habits))
	for _, habit := range habits {
		status := TodayStatus{
			HabitID:  habit.ID,
			Name:     habit.Name,
			Category: habit.Category,
		}

		if log, ok := logMap[habit.ID]; ok {
			status.IsDoneToday = true
			status.Note = log.Note
		}

		if st, ok := streakMap[habit.ID]; ok {
			status.CurrentStreak = st.CurrentStreak
		}

		statuses = append(statuses, status)
	}

	return statuses, nil
}
