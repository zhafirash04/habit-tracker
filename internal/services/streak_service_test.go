package services

import (
	"testing"
	"time"

	"habitflow/internal/models"
)

// ─── StreakService Tests ───────────────────────────────────────────

// TestStreak_FirstCheckIn verifies the very first check-in for a habit.
func TestStreak_FirstCheckIn(t *testing.T) {
	db := setupTestDB(t)
	svc := NewStreakService(db)
	userID, habitID := seedHabit(t, db, "u_first", "Minum Air", "health")

	result, err := svc.Checkin(userID, habitID, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.CurrentStreak != 1 {
		t.Errorf("expected current_streak=1, got %d", result.CurrentStreak)
	}
	if result.LongestStreak != 1 {
		t.Errorf("expected longest_streak=1, got %d", result.LongestStreak)
	}
	if !result.IsDone {
		t.Error("expected is_done=true")
	}
	if result.Date != TodayWIB() {
		t.Errorf("expected date=%s, got %s", TodayWIB(), result.Date)
	}
	if result.CompletedAt == nil {
		t.Error("expected completed_at to be set")
	}

	// Verify streak record in DB
	streak, err := svc.GetStreak(habitID)
	if err != nil {
		t.Fatalf("GetStreak failed: %v", err)
	}
	if streak.CurrentStreak != 1 || streak.LongestStreak != 1 {
		t.Errorf("DB streak mismatch: current=%d, longest=%d", streak.CurrentStreak, streak.LongestStreak)
	}
	if streak.LastDoneDate == nil || *streak.LastDoneDate != TodayWIB() {
		t.Errorf("expected last_done_date=%s, got %v", TodayWIB(), streak.LastDoneDate)
	}
}

// TestStreak_ConsecutiveDay verifies streak continuation on consecutive days.
func TestStreak_ConsecutiveDay(t *testing.T) {
	db := setupTestDB(t)
	svc := NewStreakService(db)

	t.Run("StreakContinues_LongestUnchanged", func(t *testing.T) {
		userID, habitID := seedHabit(t, db, "u_consec1", "Olahraga", "health")
		yesterday := time.Now().In(WIB).AddDate(0, 0, -1).Format("2006-01-02")
		insertLog(t, db, habitID, userID, yesterday, true)
		setStreak(t, db, habitID, 5, 10, &yesterday)

		result, err := svc.Checkin(userID, habitID, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.CurrentStreak != 6 {
			t.Errorf("expected current_streak=6, got %d", result.CurrentStreak)
		}
		if result.LongestStreak != 10 {
			t.Errorf("expected longest_streak=10 (unchanged since 6<10), got %d", result.LongestStreak)
		}
	})

	t.Run("StreakContinues_LongestUpdated", func(t *testing.T) {
		userID, habitID := seedHabit(t, db, "u_consec2", "Meditasi", "spiritual")
		yesterday := time.Now().In(WIB).AddDate(0, 0, -1).Format("2006-01-02")
		insertLog(t, db, habitID, userID, yesterday, true)
		setStreak(t, db, habitID, 9, 9, &yesterday)

		result, err := svc.Checkin(userID, habitID, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.CurrentStreak != 10 {
			t.Errorf("expected current_streak=10, got %d", result.CurrentStreak)
		}
		if result.LongestStreak != 10 {
			t.Errorf("expected longest_streak=10 (updated), got %d", result.LongestStreak)
		}
	})
}

// TestStreak_SkippedOneDay verifies streak resets when one day is skipped.
func TestStreak_SkippedOneDay(t *testing.T) {
	db := setupTestDB(t)
	svc := NewStreakService(db)
	userID, habitID := seedHabit(t, db, "u_skip1", "Baca Buku", "learning")

	twoDaysAgo := time.Now().In(WIB).AddDate(0, 0, -2).Format("2006-01-02")
	insertLog(t, db, habitID, userID, twoDaysAgo, true)
	setStreak(t, db, habitID, 5, 8, &twoDaysAgo)

	result, err := svc.Checkin(userID, habitID, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.CurrentStreak != 1 {
		t.Errorf("expected current_streak=1 (reset), got %d", result.CurrentStreak)
	}
	if result.LongestStreak != 8 {
		t.Errorf("expected longest_streak=8 (unchanged), got %d", result.LongestStreak)
	}
}

// TestStreak_SkippedMultipleDays verifies streak resets when multiple days are skipped.
func TestStreak_SkippedMultipleDays(t *testing.T) {
	db := setupTestDB(t)
	svc := NewStreakService(db)
	userID, habitID := seedHabit(t, db, "u_skip5", "Tidur Awal", "health")

	fiveDaysAgo := time.Now().In(WIB).AddDate(0, 0, -5).Format("2006-01-02")
	insertLog(t, db, habitID, userID, fiveDaysAgo, true)
	setStreak(t, db, habitID, 7, 15, &fiveDaysAgo)

	result, err := svc.Checkin(userID, habitID, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.CurrentStreak != 1 {
		t.Errorf("expected current_streak=1 (reset), got %d", result.CurrentStreak)
	}
	if result.LongestStreak != 15 {
		t.Errorf("expected longest_streak=15 (unchanged), got %d", result.LongestStreak)
	}
}

// TestStreak_SameDayDoubleCheckIn verifies double check-in on same day returns error.
func TestStreak_SameDayDoubleCheckIn(t *testing.T) {
	db := setupTestDB(t)
	svc := NewStreakService(db)
	userID, habitID := seedHabit(t, db, "u_double", "Code", "productivity")

	// First checkin succeeds
	first, err := svc.Checkin(userID, habitID, nil)
	if err != nil {
		t.Fatalf("first checkin should succeed: %v", err)
	}
	if first.CurrentStreak != 1 {
		t.Errorf("first: expected current_streak=1, got %d", first.CurrentStreak)
	}

	// Second checkin same day must fail with ErrAlreadyCheckedIn
	_, err = svc.Checkin(userID, habitID, nil)
	if err == nil {
		t.Fatal("expected ErrAlreadyCheckedIn, got nil")
	}
	if err != ErrAlreadyCheckedIn {
		t.Errorf("expected ErrAlreadyCheckedIn, got %v", err)
	}

	// Verify streak is unchanged after the failed second attempt
	streak, _ := svc.GetStreak(habitID)
	if streak.CurrentStreak != 1 {
		t.Errorf("streak should not change on duplicate: expected 1, got %d", streak.CurrentStreak)
	}
}

// TestStreak_LongestStreakNeverDecreases verifies longest_streak never goes down.
func TestStreak_LongestStreakNeverDecreases(t *testing.T) {
	db := setupTestDB(t)
	svc := NewStreakService(db)
	userID, habitID := seedHabit(t, db, "u_longest", "Meditate", "spiritual")

	// Simulate: had longest=10, then gap, now checking in again
	threeDaysAgo := time.Now().In(WIB).AddDate(0, 0, -3).Format("2006-01-02")
	insertLog(t, db, habitID, userID, threeDaysAgo, true)
	setStreak(t, db, habitID, 3, 10, &threeDaysAgo)

	result, err := svc.Checkin(userID, habitID, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Gap → reset to 1
	if result.CurrentStreak != 1 {
		t.Errorf("expected current_streak=1, got %d", result.CurrentStreak)
	}
	// longest must still be 10
	if result.LongestStreak != 10 {
		t.Errorf("expected longest_streak=10 (never decrease), got %d", result.LongestStreak)
	}
}

// TestStreak_NewLongestStreak verifies longest_streak updates when surpassed.
func TestStreak_NewLongestStreak(t *testing.T) {
	db := setupTestDB(t)
	svc := NewStreakService(db)
	userID, habitID := seedHabit(t, db, "u_newlong", "Run", "health")

	// current=10, longest=9, yesterday done → checkin → current=11, longest=11
	yesterday := time.Now().In(WIB).AddDate(0, 0, -1).Format("2006-01-02")
	insertLog(t, db, habitID, userID, yesterday, true)
	setStreak(t, db, habitID, 10, 9, &yesterday)

	result, err := svc.Checkin(userID, habitID, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.CurrentStreak != 11 {
		t.Errorf("expected current_streak=11, got %d", result.CurrentStreak)
	}
	if result.LongestStreak != 11 {
		t.Errorf("expected longest_streak=11 (surpassed old 9), got %d", result.LongestStreak)
	}
}

// TestStreak_UndoCheckIn verifies undo recalculates streak from remaining logs.
func TestStreak_UndoCheckIn(t *testing.T) {
	db := setupTestDB(t)
	svc := NewStreakService(db)

	t.Run("UndoWithPreviousDays", func(t *testing.T) {
		userID, habitID := seedHabit(t, db, "u_undo1", "Minum Air", "health")

		// Seed logs: yesterday + day-before-yesterday = 2 consecutive days
		yesterday := time.Now().In(WIB).AddDate(0, 0, -1).Format("2006-01-02")
		dayBefore := time.Now().In(WIB).AddDate(0, 0, -2).Format("2006-01-02")
		insertLog(t, db, habitID, userID, dayBefore, true)
		insertLog(t, db, habitID, userID, yesterday, true)
		setStreak(t, db, habitID, 2, 2, &yesterday)

		// Check in today, streak becomes 3
		_, err := svc.Checkin(userID, habitID, nil)
		if err != nil {
			t.Fatalf("checkin failed: %v", err)
		}

		// Undo today's checkin → streak recalculated from yesterday backwards
		err = svc.UndoCheckin(userID, habitID)
		if err != nil {
			t.Fatalf("undo should succeed: %v", err)
		}

		streak, _ := svc.GetStreak(habitID)
		// Yesterday + DayBefore = 2 consecutive days
		if streak.CurrentStreak != 2 {
			t.Errorf("expected recalculated streak=2, got %d", streak.CurrentStreak)
		}
		if streak.LastDoneDate == nil || *streak.LastDoneDate != yesterday {
			t.Errorf("expected last_done_date=%s, got %v", yesterday, streak.LastDoneDate)
		}
		// longest never decreases
		if streak.LongestStreak < 2 {
			t.Errorf("expected longest_streak >= 2, got %d", streak.LongestStreak)
		}
	})

	t.Run("UndoOnlyCheckin_ResetsToZero", func(t *testing.T) {
		userID, habitID := seedHabit(t, db, "u_undo2", "Code", "productivity")

		// Check in today (first ever), then undo
		_, err := svc.Checkin(userID, habitID, nil)
		if err != nil {
			t.Fatalf("checkin failed: %v", err)
		}

		err = svc.UndoCheckin(userID, habitID)
		if err != nil {
			t.Fatalf("undo should succeed: %v", err)
		}

		streak, _ := svc.GetStreak(habitID)
		if streak.CurrentStreak != 0 {
			t.Errorf("expected streak=0 after undo only checkin, got %d", streak.CurrentStreak)
		}
		if streak.LastDoneDate != nil {
			t.Errorf("expected last_done_date=nil, got %v", streak.LastDoneDate)
		}
	})

	t.Run("UndoWithNonConsecutivePrevious", func(t *testing.T) {
		userID, habitID := seedHabit(t, db, "u_undo3", "Walk", "health")

		// Seed: 4 days ago done, 3 days ago done, gap at 2 days ago, yesterday done
		d4 := time.Now().In(WIB).AddDate(0, 0, -4).Format("2006-01-02")
		d3 := time.Now().In(WIB).AddDate(0, 0, -3).Format("2006-01-02")
		d1 := time.Now().In(WIB).AddDate(0, 0, -1).Format("2006-01-02")
		insertLog(t, db, habitID, userID, d4, true)
		insertLog(t, db, habitID, userID, d3, true)
		// d2 skipped
		insertLog(t, db, habitID, userID, d1, true)
		setStreak(t, db, habitID, 1, 3, &d1)

		// Check in today
		_, err := svc.Checkin(userID, habitID, nil)
		if err != nil {
			t.Fatalf("checkin failed: %v", err)
		}

		// Undo today → last log is yesterday, gap before that → streak=1
		err = svc.UndoCheckin(userID, habitID)
		if err != nil {
			t.Fatalf("undo failed: %v", err)
		}

		streak, _ := svc.GetStreak(habitID)
		if streak.CurrentStreak != 1 {
			t.Errorf("expected recalculated streak=1, got %d", streak.CurrentStreak)
		}
	})

	t.Run("UndoNoCheckinToday_Error", func(t *testing.T) {
		userID, habitID := seedHabit(t, db, "u_undo4", "Journal", "creative")

		err := svc.UndoCheckin(userID, habitID)
		if err == nil {
			t.Fatal("expected ErrNoCheckinToday, got nil")
		}
		if err != ErrNoCheckinToday {
			t.Errorf("expected ErrNoCheckinToday, got: %v", err)
		}
	})
}

// TestCheckin_WithNote verifies note is stored with check-in.
func TestCheckin_WithNote(t *testing.T) {
	db := setupTestDB(t)
	svc := NewStreakService(db)
	userID, habitID := seedHabit(t, db, "u_note", "Journaling", "creative")

	note := "Wrote 2 pages today"
	result, err := svc.Checkin(userID, habitID, &note)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Note == nil || *result.Note != note {
		t.Errorf("expected note=%q, got %v", note, result.Note)
	}

	// Verify in DB
	var log models.HabitLog
	db.Where("habit_id = ? AND date = ?", habitID, TodayWIB()).First(&log)
	if log.Note == nil || *log.Note != note {
		t.Errorf("DB note mismatch: expected %q, got %v", note, log.Note)
	}
}

// TestCheckin_HabitNotFound verifies error for non-existent habit.
func TestCheckin_HabitNotFound(t *testing.T) {
	db := setupTestDB(t)
	svc := NewStreakService(db)
	seedHabit(t, db, "u_notfound", "Dummy", "general")

	_, err := svc.Checkin(1, 9999, nil)
	if err == nil {
		t.Fatal("expected error for non-existent habit")
	}
}

// TestCheckin_InactiveHabit verifies error for deactivated habit.
func TestCheckin_InactiveHabit(t *testing.T) {
	db := setupTestDB(t)
	svc := NewStreakService(db)
	userID, habitID := seedHabit(t, db, "u_inactive", "Old Habit", "general")

	// Deactivate the habit
	db.Model(&models.Habit{}).Where("id = ?", habitID).Update("is_active", false)

	_, err := svc.Checkin(userID, habitID, nil)
	if err == nil {
		t.Fatal("expected error for inactive habit")
	}
}

// TestCheckin_WrongUser verifies a user cannot check in another user's habit.
func TestCheckin_WrongUser(t *testing.T) {
	db := setupTestDB(t)
	svc := NewStreakService(db)
	_, habitID := seedHabit(t, db, "u_owner", "My Habit", "general")

	// Create second user
	user2 := models.User{Name: "Other", Email: "other@test.com", PasswordHash: "hash"}
	db.Create(&user2)

	_, err := svc.Checkin(user2.ID, habitID, nil)
	if err == nil {
		t.Fatal("expected error: user2 should not check in user1's habit")
	}
}

// TestGetTodayStatus verifies the today status aggregation.
func TestGetTodayStatus(t *testing.T) {
	db := setupTestDB(t)
	svc := NewStreakService(db)

	t.Run("MixedStatus", func(t *testing.T) {
		user := models.User{Name: "u_status", Email: "status@test.com", PasswordHash: "hash"}
		db.Create(&user)

		h1 := models.Habit{UserID: user.ID, Name: "HabitA", Category: "health", IsActive: true}
		h2 := models.Habit{UserID: user.ID, Name: "HabitB", Category: "learning", IsActive: true}
		h3 := models.Habit{UserID: user.ID, Name: "Inactive", Category: "general", IsActive: true}
		db.Create(&h1)
		db.Create(&h2)
		db.Create(&h3)
		db.Model(&models.Habit{}).Where("id = ?", h3.ID).Update("is_active", false)
		db.Create(&models.Streak{HabitID: h1.ID, CurrentStreak: 3, LongestStreak: 5})
		db.Create(&models.Streak{HabitID: h2.ID, CurrentStreak: 0, LongestStreak: 0})

		// Check in h1 today
		svc.Checkin(user.ID, h1.ID, nil)

		statuses, err := svc.GetTodayStatus(user.ID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Should only include active habits (h1, h2), not h3
		if len(statuses) != 2 {
			t.Fatalf("expected 2 statuses, got %d", len(statuses))
		}

		for _, s := range statuses {
			if s.Name == "HabitA" && !s.IsDoneToday {
				t.Error("HabitA should be done today")
			}
			if s.Name == "HabitB" && s.IsDoneToday {
				t.Error("HabitB should NOT be done today")
			}
		}
	})

	t.Run("NoHabits", func(t *testing.T) {
		user := models.User{Name: "u_empty", Email: "empty@test.com", PasswordHash: "hash"}
		db.Create(&user)

		statuses, err := svc.GetTodayStatus(user.ID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(statuses) != 0 {
			t.Errorf("expected 0 statuses, got %d", len(statuses))
		}
	})
}
