package services

import (
	"testing"
	"time"

	"habitflow/internal/models"
)

// ─── ScoreService Tests ───────────────────────────────────────────

// TestScore_PerfectScore verifies 100% score when all days are done.
func TestScore_PerfectScore(t *testing.T) {
	db := setupTestDB(t)
	svc := NewScoreService(db)

	user := models.User{Name: "score_perf", Email: "score_perf@test.com", PasswordHash: "hash"}
	db.Create(&user)

	habit := models.Habit{UserID: user.ID, Name: "Minum Air", Category: "health", IsActive: true}
	db.Create(&habit)
	db.Model(&habit).Update("created_at", time.Now().AddDate(0, 0, -10))
	db.Create(&models.Streak{HabitID: habit.ID})

	// Insert logs for all 7 days in range
	for i := 0; i < 7; i++ {
		date := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
		insertLog(t, db, habit.ID, user.ID, date, true)
	}

	result, err := svc.Calculate(user.ID, 7)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Overall < 95 {
		t.Errorf("expected overall score ~100, got %.2f", result.Overall)
	}
	if len(result.HabitScores) != 1 {
		t.Fatalf("expected 1 habit score, got %d", len(result.HabitScores))
	}
	if result.HabitScores[0].Score < 95 {
		t.Errorf("expected habit score ~100, got %.2f", result.HabitScores[0].Score)
	}
	if result.HabitScores[0].Done != 7 {
		t.Errorf("expected done=7, got %d", result.HabitScores[0].Done)
	}
}

// TestScore_ZeroScore verifies 0% score when no days are done.
func TestScore_ZeroScore(t *testing.T) {
	db := setupTestDB(t)
	svc := NewScoreService(db)
	userID, _ := seedHabit(t, db, "score_zero", "Olahraga", "health")

	result, err := svc.Calculate(userID, 7)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Overall != 0 {
		t.Errorf("expected overall=0, got %.2f", result.Overall)
	}
	if result.HabitScores[0].Done != 0 {
		t.Errorf("expected done=0, got %d", result.HabitScores[0].Done)
	}
	if result.HabitScores[0].Score != 0 {
		t.Errorf("expected score=0, got %.2f", result.HabitScores[0].Score)
	}
}

// TestScore_NoHabits verifies empty result for user with no habits.
func TestScore_NoHabits(t *testing.T) {
	db := setupTestDB(t)
	svc := NewScoreService(db)

	user := models.User{Name: "score_empty", Email: "score_empty@test.com", PasswordHash: "hash"}
	db.Create(&user)

	result, err := svc.Calculate(user.ID, 7)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Overall != 0 {
		t.Errorf("expected overall=0 with no habits, got %.2f", result.Overall)
	}
	if len(result.HabitScores) != 0 {
		t.Errorf("expected 0 habit scores, got %d", len(result.HabitScores))
	}
}

// TestScore_PartialCompletion verifies score for partial completion.
func TestScore_PartialCompletion(t *testing.T) {
	db := setupTestDB(t)
	svc := NewScoreService(db)

	user := models.User{Name: "score_half", Email: "score_half@test.com", PasswordHash: "hash"}
	db.Create(&user)

	habit := models.Habit{UserID: user.ID, Name: "Baca", Category: "learning", IsActive: true}
	db.Create(&habit)
	db.Model(&habit).Update("created_at", time.Now().AddDate(0, 0, -10))
	db.Create(&models.Streak{HabitID: habit.ID})

	// Insert logs for 3 out of 7 days
	for i := 0; i < 3; i++ {
		date := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
		insertLog(t, db, habit.ID, user.ID, date, true)
	}

	result, err := svc.Calculate(user.ID, 7)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 3/7 ≈ 42.86%
	if result.Overall < 30 || result.Overall > 60 {
		t.Errorf("expected overall score ~43%%, got %.2f", result.Overall)
	}
	if result.HabitScores[0].Done != 3 {
		t.Errorf("expected done=3, got %d", result.HabitScores[0].Done)
	}
}

// TestScore_MultipleHabits verifies average calculation across multiple habits.
func TestScore_MultipleHabits(t *testing.T) {
	db := setupTestDB(t)
	svc := NewScoreService(db)

	user := models.User{Name: "score_multi", Email: "score_multi@test.com", PasswordHash: "hash"}
	db.Create(&user)

	// Habit 1: all 7 days done → ~100%
	h1 := models.Habit{UserID: user.ID, Name: "HabitA", Category: "health", IsActive: true}
	db.Create(&h1)
	db.Model(&h1).Update("created_at", time.Now().AddDate(0, 0, -10))
	db.Create(&models.Streak{HabitID: h1.ID})
	for i := 0; i < 7; i++ {
		insertLog(t, db, h1.ID, user.ID, time.Now().AddDate(0, 0, -i).Format("2006-01-02"), true)
	}

	// Habit 2: 0 days done → 0%
	h2 := models.Habit{UserID: user.ID, Name: "HabitB", Category: "learning", IsActive: true}
	db.Create(&h2)
	db.Model(&h2).Update("created_at", time.Now().AddDate(0, 0, -10))
	db.Create(&models.Streak{HabitID: h2.ID})

	result, err := svc.Calculate(user.ID, 7)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.HabitScores) != 2 {
		t.Fatalf("expected 2 habit scores, got %d", len(result.HabitScores))
	}

	// Average: (~100 + 0) / 2 ≈ ~50
	if result.Overall < 40 || result.Overall > 60 {
		t.Errorf("expected overall ~50%%, got %.2f", result.Overall)
	}
}

// TestScore_ScoreNeverExceeds100 verifies score is capped at 100%.
func TestScore_ScoreNeverExceeds100(t *testing.T) {
	db := setupTestDB(t)
	svc := NewScoreService(db)

	user := models.User{Name: "score_cap", Email: "score_cap@test.com", PasswordHash: "hash"}
	db.Create(&user)

	// Create habit with CreatedAt=today (totalDays=1)
	habit := models.Habit{UserID: user.ID, Name: "New", Category: "health", IsActive: true}
	db.Create(&habit)
	db.Create(&models.Streak{HabitID: habit.ID})

	// Insert 2 logs for today (simulating edge case)
	today := time.Now().Format("2006-01-02")
	insertLog(t, db, habit.ID, user.ID, today, true)

	result, err := svc.Calculate(user.ID, 7)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Overall > 100 {
		t.Errorf("score should not exceed 100, got %.2f", result.Overall)
	}
	for _, hs := range result.HabitScores {
		if hs.Score > 100 {
			t.Errorf("habit score should not exceed 100, got %.2f", hs.Score)
		}
	}
}

// TestScore_HabitBaruSehari verifies score for a brand new habit (created today).
func TestScore_HabitBaruSehari(t *testing.T) {
	db := setupTestDB(t)
	svc := NewScoreService(db)

	user := models.User{Name: "score_baru", Email: "score_baru@test.com", PasswordHash: "hash"}
	db.Create(&user)

	// Habit created today (defaulting to now)
	habit := models.Habit{UserID: user.ID, Name: "BrandNew", Category: "productivity", IsActive: true}
	db.Create(&habit)
	db.Create(&models.Streak{HabitID: habit.ID})

	// Check in today
	today := time.Now().Format("2006-01-02")
	insertLog(t, db, habit.ID, user.ID, today, true)

	result, err := svc.Calculate(user.ID, 7)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Brand new habit: totalDays may be 1 or 2 due to ceil rounding
	if result.Overall <= 0 {
		t.Errorf("brand new habit with 1 done should have positive score, got %.2f", result.Overall)
	}
	if result.HabitScores[0].Done != 1 {
		t.Errorf("expected done=1, got %d", result.HabitScores[0].Done)
	}
	if result.HabitScores[0].Total < 1 {
		t.Errorf("expected total >= 1 day, got %d", result.HabitScores[0].Total)
	}
}

// TestScore_InactiveHabitsExcluded verifies inactive habits are excluded.
func TestScore_InactiveHabitsExcluded(t *testing.T) {
	db := setupTestDB(t)
	svc := NewScoreService(db)

	user := models.User{Name: "score_inact", Email: "score_inact@test.com", PasswordHash: "hash"}
	db.Create(&user)

	// Active habit with all days done
	active := models.Habit{UserID: user.ID, Name: "Active", Category: "health", IsActive: true}
	db.Create(&active)
	db.Model(&active).Update("created_at", time.Now().AddDate(0, 0, -10))
	db.Create(&models.Streak{HabitID: active.ID})
	for i := 0; i < 7; i++ {
		insertLog(t, db, active.ID, user.ID, time.Now().AddDate(0, 0, -i).Format("2006-01-02"), true)
	}

	// Inactive habit with no logs
	inactive := models.Habit{UserID: user.ID, Name: "Inactive", Category: "health", IsActive: true}
	db.Create(&inactive)
	db.Model(&inactive).Update("is_active", false)
	db.Create(&models.Streak{HabitID: inactive.ID})

	result, err := svc.Calculate(user.ID, 7)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.HabitScores) != 1 {
		t.Fatalf("expected 1 habit score (inactive excluded), got %d", len(result.HabitScores))
	}
	if result.HabitScores[0].HabitName != "Active" {
		t.Errorf("expected only 'Active', got %q", result.HabitScores[0].HabitName)
	}
}

// TestScore_Periods verifies correct period strings.
func TestScore_Periods(t *testing.T) {
	db := setupTestDB(t)
	svc := NewScoreService(db)

	user := models.User{Name: "score_per", Email: "score_per@test.com", PasswordHash: "hash"}
	db.Create(&user)

	tests := []struct {
		name     string
		days     int
		expected string
	}{
		{"7days", 7, "last_7_days"},
		{"30days", 30, "last_30_days"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := svc.Calculate(user.ID, tc.days)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result.Period != tc.expected {
				t.Errorf("expected period=%q, got %q", tc.expected, result.Period)
			}
		})
	}

	// Custom period (neither 7 nor 30)
	t.Run("CustomPeriod", func(t *testing.T) {
		result, err := svc.Calculate(user.ID, 14)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Period == "last_7_days" || result.Period == "last_30_days" {
			t.Errorf("expected custom period string, got %q", result.Period)
		}
		// Should contain a date range format
		if len(result.Period) < 10 {
			t.Errorf("expected date range period, got %q", result.Period)
		}
	})
}

// TestScore_HabitNameInResult verifies habit name is returned in result.
func TestScore_HabitNameInResult(t *testing.T) {
	db := setupTestDB(t)
	svc := NewScoreService(db)

	user := models.User{Name: "score_name", Email: "score_name@test.com", PasswordHash: "hash"}
	db.Create(&user)

	habit := models.Habit{UserID: user.ID, Name: "Jogging Pagi", Category: "health", IsActive: true}
	db.Create(&habit)
	db.Create(&models.Streak{HabitID: habit.ID})

	result, err := svc.Calculate(user.ID, 7)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.HabitScores) != 1 {
		t.Fatalf("expected 1 habit score, got %d", len(result.HabitScores))
	}
	if result.HabitScores[0].HabitName != "Jogging Pagi" {
		t.Errorf("expected habit name 'Jogging Pagi', got %q", result.HabitScores[0].HabitName)
	}
	if result.HabitScores[0].HabitID != habit.ID {
		t.Errorf("expected habit_id=%d, got %d", habit.ID, result.HabitScores[0].HabitID)
	}
}
