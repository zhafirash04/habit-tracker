package services

import (
	"testing"
	"time"

	"habitflow/internal/models"
)

func TestDailyReport_FrozenSnapshot_IgnoresFutureHabits(t *testing.T) {
	db := setupTestDB(t)
	svc := NewDailyReportService(db)

	user := models.User{Name: "SnapshotUser", Email: "snapshot@test.com", PasswordHash: "hash"}
	db.Create(&user)

	habitA := models.Habit{UserID: user.ID, Name: "Habit A", Category: "general", IsActive: true, CreatedAt: time.Date(2026, 3, 10, 8, 0, 0, 0, WIB)}
	habitB := models.Habit{UserID: user.ID, Name: "Habit B", Category: "general", IsActive: true, CreatedAt: time.Date(2026, 3, 11, 9, 0, 0, 0, WIB)}
	habitFuture := models.Habit{UserID: user.ID, Name: "Habit Future", Category: "general", IsActive: true, CreatedAt: time.Date(2026, 3, 12, 9, 0, 0, 0, WIB)}
	db.Create(&habitA)
	db.Create(&habitB)
	db.Create(&habitFuture)

	db.Create(&models.Streak{HabitID: habitA.ID})
	db.Create(&models.Streak{HabitID: habitB.ID})
	db.Create(&models.Streak{HabitID: habitFuture.ID})

	db.Create(&models.HabitLog{HabitID: habitA.ID, UserID: user.ID, Date: "2026-03-11", IsDone: true})

	report, err := svc.Generate(user.ID, time.Date(2026, 3, 11, 12, 0, 0, 0, WIB))
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	if report.Summary.TotalHabits != 2 {
		t.Fatalf("expected total_habits=2, got %d", report.Summary.TotalHabits)
	}
	if report.Summary.Completed != 1 {
		t.Fatalf("expected completed=1, got %d", report.Summary.Completed)
	}
	if report.Summary.Remaining != 1 {
		t.Fatalf("expected remaining=1, got %d", report.Summary.Remaining)
	}
}

func TestDailyReport_FrozenSnapshot_KeepsSoftDeletedHabits(t *testing.T) {
	db := setupTestDB(t)
	svc := NewDailyReportService(db)

	user := models.User{Name: "DeletedUser", Email: "deleted@test.com", PasswordHash: "hash"}
	db.Create(&user)

	habit := models.Habit{UserID: user.ID, Name: "Old Habit", Category: "general", IsActive: true, CreatedAt: time.Date(2026, 3, 10, 8, 0, 0, 0, WIB)}
	db.Create(&habit)
	db.Create(&models.Streak{HabitID: habit.ID})
	db.Create(&models.HabitLog{HabitID: habit.ID, UserID: user.ID, Date: "2026-03-11", IsDone: true})

	// Soft delete happens later, report date must remain stable.
	db.Model(&habit).Update("is_active", false)

	report, err := svc.Generate(user.ID, time.Date(2026, 3, 11, 12, 0, 0, 0, WIB))
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	if report.Summary.TotalHabits != 1 {
		t.Fatalf("expected total_habits=1, got %d", report.Summary.TotalHabits)
	}
	if len(report.Habits) != 1 {
		t.Fatalf("expected 1 habit row, got %d", len(report.Habits))
	}
	if !report.Habits[0].IsDone {
		t.Fatal("expected historical habit to stay done")
	}
	if !report.Habits[0].IsDeletedNow {
		t.Fatal("expected is_deleted_now=true for currently inactive habit")
	}
}
