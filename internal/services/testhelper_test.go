package services

import (
	"testing"
	"time"

	"habitflow/internal/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// setupTestDB creates an in-memory SQLite database for testing.
func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	err = db.AutoMigrate(
		&models.User{},
		&models.Habit{},
		&models.HabitLog{},
		&models.Streak{},
	)
	if err != nil {
		t.Fatalf("failed to migrate test db: %v", err)
	}

	return db
}

// seedHabit inserts a user, habit, and its streak record, returning the habit ID.
func seedHabit(t *testing.T, db *gorm.DB, userName, habitName, category string) (userID uint, habitID uint) {
	t.Helper()
	user := models.User{Name: userName, Email: userName + "@test.com", PasswordHash: "hash"}
	db.Create(&user)

	habit := models.Habit{
		UserID:   user.ID,
		Name:     habitName,
		Category: category,
		IsActive: true,
	}
	db.Create(&habit)

	streak := models.Streak{
		HabitID:       habit.ID,
		CurrentStreak: 0,
		LongestStreak: 0,
	}
	db.Create(&streak)

	return user.ID, habit.ID
}

// insertLog inserts a HabitLog for the given habit on the given date.
func insertLog(t *testing.T, db *gorm.DB, habitID, userID uint, date string, isDone bool) {
	t.Helper()
	log := models.HabitLog{
		HabitID: habitID,
		UserID:  userID,
		Date:    date,
		IsDone:  isDone,
	}
	if err := db.Create(&log).Error; err != nil {
		t.Fatalf("failed to insert log for date %s: %v", date, err)
	}
}

// setStreak directly sets streak values in the database.
func setStreak(t *testing.T, db *gorm.DB, habitID uint, current, longest int, lastDate *string) {
	t.Helper()
	db.Model(&models.Streak{}).Where("habit_id = ?", habitID).Updates(map[string]interface{}{
		"current_streak": current,
		"longest_streak": longest,
		"last_done_date": lastDate,
	})
}

// genInsights calls InsightService.Generate with a date range ending today.
func genInsights(t *testing.T, svc *InsightService, userID uint, daysBack int) ([]Insight, error) {
	t.Helper()
	end := time.Now().Format("2006-01-02")
	start := time.Now().AddDate(0, 0, -daysBack).Format("2006-01-02")
	return svc.Generate(userID, start, end)
}
