package services

import (
	"strings"
	"testing"
	"time"

	"habitflow/internal/models"
)

// ─── InsightService Tests ─────────────────────────────────────────

// TestInsight_BestDay verifies best day detection based on most completions.
func TestInsight_BestDay(t *testing.T) {
	db := setupTestDB(t)
	svc := NewInsightService(db)
	userID, habitID := seedHabit(t, db, "ins_best", "Minum Air", "health")

	// Insert multiple logs on Mondays in the last 30 days
	now := time.Now()
	for i := 0; i < 30; i++ {
		d := now.AddDate(0, 0, -i)
		if d.Weekday() == time.Monday {
			insertLog(t, db, habitID, userID, d.Format("2006-01-02"), true)
		}
	}

	insights, err := genInsights(t, svc, userID, 30)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found := false
	for _, ins := range insights {
		if ins.Type == "best_day" {
			found = true
			if !strings.Contains(ins.Message, "Senin") {
				t.Errorf("expected best day to mention Senin, got message: %s", ins.Message)
			}
		}
	}
	if !found {
		t.Error("expected best_day insight but not found")
	}
}

// TestInsight_WeakestDay verifies best_day detection correctly identifies the dominant day.
func TestInsight_WeakestDay(t *testing.T) {
	db := setupTestDB(t)
	svc := NewInsightService(db)
	userID, habitID := seedHabit(t, db, "ins_weak", "Olahraga", "health")

	// Insert logs on specific days to make Saturday the most done day
	now := time.Now()
	for i := 0; i < 30; i++ {
		d := now.AddDate(0, 0, -i)
		if d.Weekday() == time.Saturday {
			insertLog(t, db, habitID, userID, d.Format("2006-01-02"), true)
		}
	}

	insights, err := genInsights(t, svc, userID, 30)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found := false
	for _, ins := range insights {
		if ins.Type == "best_day" {
			found = true
			if !strings.Contains(ins.Message, "Sabtu") {
				t.Errorf("expected best day to mention Sabtu, got: %s", ins.Message)
			}
		}
	}
	if !found {
		t.Error("expected best_day insight but not found")
	}
}

// TestInsight_Declining verifies habits with no check-in in 3+ days are flagged.
func TestInsight_Declining(t *testing.T) {
	db := setupTestDB(t)
	svc := NewInsightService(db)

	t.Run("LastDone5DaysAgo", func(t *testing.T) {
		userID, habitID := seedHabit(t, db, "ins_dec1", "Olahraga", "health")

		fiveDaysAgo := time.Now().AddDate(0, 0, -5).Format("2006-01-02")
		insertLog(t, db, habitID, userID, fiveDaysAgo, true)

		insights, err := genInsights(t, svc, userID, 7)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		found := false
		for _, ins := range insights {
			if ins.Type == "declining" {
				found = true
				if !strings.Contains(ins.Message, "Olahraga") {
					t.Errorf("expected declining for 'Olahraga', got: %s", ins.Message)
				}
				if ins.HabitID != habitID {
					t.Errorf("expected habit_id=%d, got %d", habitID, ins.HabitID)
				}
			}
		}
		if !found {
			t.Error("expected declining insight but not found")
		}
	})

	t.Run("NoLogsAtAll", func(t *testing.T) {
		userID, habitID := seedHabit(t, db, "ins_dec2", "Coding", "productivity")

		insights, err := genInsights(t, svc, userID, 7)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// No logs at all → count=0 → should be declining
		found := false
		for _, ins := range insights {
			if ins.Type == "declining" && ins.HabitID == habitID {
				found = true
			}
		}
		if !found {
			t.Error("habit with no logs should appear as declining")
		}
	})
}

// TestInsight_HighConsistency verifies consistency detection (>=6/7 days).
func TestInsight_HighConsistency(t *testing.T) {
	db := setupTestDB(t)
	svc := NewInsightService(db)

	t.Run("7of7Days", func(t *testing.T) {
		userID, habitID := seedHabit(t, db, "ins_con1", "Meditasi", "spiritual")

		for i := 0; i < 7; i++ {
			date := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
			insertLog(t, db, habitID, userID, date, true)
		}

		insights, err := genInsights(t, svc, userID, 7)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		found := false
		for _, ins := range insights {
			if ins.Type == "consistency" {
				found = true
				if !strings.Contains(ins.Message, "Meditasi") {
					t.Errorf("expected consistency for 'Meditasi', got: %s", ins.Message)
				}
			}
		}
		if !found {
			t.Error("expected consistency insight for 7/7 days")
		}
	})

	t.Run("6of7Days", func(t *testing.T) {
		userID, habitID := seedHabit(t, db, "ins_con2", "Jogging", "health")

		// Insert 6 out of 7 days (skip one day)
		for i := 0; i < 7; i++ {
			if i == 3 { // skip one day
				continue
			}
			date := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
			insertLog(t, db, habitID, userID, date, true)
		}

		insights, err := genInsights(t, svc, userID, 7)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		found := false
		for _, ins := range insights {
			if ins.Type == "consistency" && strings.Contains(ins.Message, "Jogging") {
				found = true
			}
		}
		if !found {
			t.Error("expected consistency insight for 6/7 days (should meet >=6 threshold)")
		}
	})

	t.Run("5of7Days_NoConsistency", func(t *testing.T) {
		userID, habitID := seedHabit(t, db, "ins_con3", "Reading", "learning")

		// Insert only 5 out of 7 days
		for i := 0; i < 5; i++ {
			date := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
			insertLog(t, db, habitID, userID, date, true)
		}

		insights, err := genInsights(t, svc, userID, 7)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		for _, ins := range insights {
			if ins.Type == "consistency" && ins.HabitID == habitID {
				t.Error("5/7 days should NOT trigger consistency insight")
			}
		}
	})
}

// TestInsight_Encouragement_WhenNoData verifies fallback when no insights found.
func TestInsight_Encouragement_WhenNoData(t *testing.T) {
	db := setupTestDB(t)
	svc := NewInsightService(db)

	// User with no habits → no best_day, no declining, no consistency
	user := models.User{Name: "ins_enc", Email: "ins_enc@test.com", PasswordHash: "hash"}
	db.Create(&user)

	insights, err := genInsights(t, svc, user.ID, 7)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(insights) != 1 {
		t.Fatalf("expected exactly 1 encouragement insight, got %d", len(insights))
	}
	if insights[0].Type != "encouragement" {
		t.Errorf("expected type='encouragement', got %q", insights[0].Type)
	}
	if insights[0].Message == "" {
		t.Error("encouragement message should not be empty")
	}
}

// TestInsight_NoDeclineIfRecentlyDone verifies no declining for recently active habits.
func TestInsight_NoDeclineIfRecentlyDone(t *testing.T) {
	db := setupTestDB(t)
	svc := NewInsightService(db)

	t.Run("DoneToday", func(t *testing.T) {
		userID, habitID := seedHabit(t, db, "ins_nodec1", "Code", "productivity")

		today := time.Now().Format("2006-01-02")
		insertLog(t, db, habitID, userID, today, true)

		insights, err := genInsights(t, svc, userID, 7)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		for _, ins := range insights {
			if ins.Type == "declining" && ins.HabitID == habitID {
				t.Error("habit done today should NOT be declining")
			}
		}
	})

	t.Run("DoneYesterday", func(t *testing.T) {
		userID, habitID := seedHabit(t, db, "ins_nodec2", "Walk", "health")

		yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
		insertLog(t, db, habitID, userID, yesterday, true)

		insights, err := genInsights(t, svc, userID, 7)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		for _, ins := range insights {
			if ins.Type == "declining" && ins.HabitID == habitID {
				t.Error("habit done yesterday should NOT be declining")
			}
		}
	})

	t.Run("Done2DaysAgo", func(t *testing.T) {
		userID, habitID := seedHabit(t, db, "ins_nodec3", "Meditate", "spiritual")

		twoDaysAgo := time.Now().AddDate(0, 0, -2).Format("2006-01-02")
		insertLog(t, db, habitID, userID, twoDaysAgo, true)

		insights, err := genInsights(t, svc, userID, 7)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		for _, ins := range insights {
			if ins.Type == "declining" && ins.HabitID == habitID {
				t.Error("habit done 2 days ago should NOT be declining (threshold is 3 days)")
			}
		}
	})
}

// TestInsight_MultipleTypes verifies multiple insight types can coexist.
func TestInsight_MultipleTypes(t *testing.T) {
	db := setupTestDB(t)
	svc := NewInsightService(db)

	user := models.User{Name: "ins_multi", Email: "ins_multi@test.com", PasswordHash: "hash"}
	db.Create(&user)

	// Habit A: highly consistent (7/7 days) → consistency + best_day
	hA := models.Habit{UserID: user.ID, Name: "HabitA", Category: "health", IsActive: true}
	db.Create(&hA)
	db.Create(&models.Streak{HabitID: hA.ID})
	for i := 0; i < 7; i++ {
		insertLog(t, db, hA.ID, user.ID, time.Now().AddDate(0, 0, -i).Format("2006-01-02"), true)
	}

	// Habit B: declining (last done 5 days ago) → declining
	hB := models.Habit{UserID: user.ID, Name: "HabitB", Category: "learning", IsActive: true}
	db.Create(&hB)
	db.Create(&models.Streak{HabitID: hB.ID})
	insertLog(t, db, hB.ID, user.ID, time.Now().AddDate(0, 0, -5).Format("2006-01-02"), true)

	insights, err := genInsights(t, svc, user.ID, 7)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	types := make(map[string]bool)
	for _, ins := range insights {
		types[ins.Type] = true
	}

	if !types["best_day"] {
		t.Error("expected best_day insight")
	}
	if !types["consistency"] {
		t.Error("expected consistency insight")
	}
	if !types["declining"] {
		t.Error("expected declining insight")
	}
	// Encouragement should NOT appear when other insights exist
	if types["encouragement"] {
		t.Error("encouragement should not appear when other insights exist")
	}
}

// TestInsight_AllDaysPerfect verifies insights when all days are equally active.
func TestInsight_AllDaysPerfect(t *testing.T) {
	db := setupTestDB(t)
	svc := NewInsightService(db)

	user := models.User{Name: "ins_allperf", Email: "ins_allperf@test.com", PasswordHash: "hash"}
	db.Create(&user)

	h1 := models.Habit{UserID: user.ID, Name: "Daily", Category: "health", IsActive: true}
	db.Create(&h1)
	db.Create(&models.Streak{HabitID: h1.ID})

	// Insert logs every day for 28 days (4 full weeks → all days equal)
	for i := 0; i < 28; i++ {
		insertLog(t, db, h1.ID, user.ID, time.Now().AddDate(0, 0, -i).Format("2006-01-02"), true)
	}

	insights, err := genInsights(t, svc, user.ID, 28)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should still generate best_day (one of the 7 days will appear)
	foundBestDay := false
	foundConsistency := false
	for _, ins := range insights {
		if ins.Type == "best_day" {
			foundBestDay = true
		}
		if ins.Type == "consistency" {
			foundConsistency = true
		}
	}
	if !foundBestDay {
		t.Error("expected best_day insight even with all days equal")
	}
	if !foundConsistency {
		t.Error("expected consistency insight for 28/28 days")
	}

	// Encouragement should NOT appear
	for _, ins := range insights {
		if ins.Type == "encouragement" {
			t.Error("encouragement should not appear when other insights exist")
		}
	}
}

// TestInsight_DayNameTranslation verifies all day names are in Indonesian.
func TestInsight_DayNameTranslation(t *testing.T) {
	db := setupTestDB(t)
	svc := NewInsightService(db)

	validDayNames := []string{"Senin", "Selasa", "Rabu", "Kamis", "Jumat", "Sabtu", "Minggu"}

	// For each weekday, verify the Indonesian name appears in insights
	for _, dayName := range validDayNames {
		t.Run(dayName, func(t *testing.T) {
			user := models.User{
				Name:         "ins_day_" + dayName,
				Email:        "ins_day_" + dayName + "@test.com",
				PasswordHash: "hash",
			}
			db.Create(&user)

			h := models.Habit{UserID: user.ID, Name: "Test_" + dayName, Category: "health", IsActive: true}
			db.Create(&h)
			db.Create(&models.Streak{HabitID: h.ID})

			// Find the target weekday in the last 30 days and add logs
			targetDay := map[string]time.Weekday{
				"Senin": time.Monday, "Selasa": time.Tuesday, "Rabu": time.Wednesday,
				"Kamis": time.Thursday, "Jumat": time.Friday, "Sabtu": time.Saturday,
				"Minggu": time.Sunday,
			}[dayName]

			now := time.Now()
			for i := 0; i < 30; i++ {
				d := now.AddDate(0, 0, -i)
				if d.Weekday() == targetDay {
					insertLog(t, db, h.ID, user.ID, d.Format("2006-01-02"), true)
				}
			}

			insights, _ := genInsights(t, svc, user.ID, 30)

			for _, ins := range insights {
				if ins.Type == "best_day" {
					if !strings.Contains(ins.Message, dayName) {
						t.Errorf("expected day name %q in message, got: %s", dayName, ins.Message)
					}
				}
			}
		})
	}
}
