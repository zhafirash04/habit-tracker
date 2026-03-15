package services

import (
	"testing"
	"time"

	"habitflow/internal/models"
)

func TestWeeklyBreakdown_IgnoresRetroactiveHabits(t *testing.T) {
	db := setupTestDB(t)
	scoreSvc := NewScoreService(db)
	insightSvc := NewInsightService(db)
	svc := NewReportService(db, scoreSvc, insightSvc)

	user := models.User{Name: "WeeklyUser", Email: "weekly@test.com", PasswordHash: "hash"}
	db.Create(&user)

	// HabitA created on March 10
	habitA := models.Habit{UserID: user.ID, Name: "Habit A", Category: "general", IsActive: true, CreatedAt: time.Date(2026, 3, 10, 8, 0, 0, 0, time.UTC)}
	db.Create(&habitA)
	db.Create(&models.Streak{HabitID: habitA.ID})

	// HabitB created on March 12
	habitB := models.Habit{UserID: user.ID, Name: "Habit B", Category: "general", IsActive: true, CreatedAt: time.Date(2026, 3, 12, 9, 0, 0, 0, time.UTC)}
	db.Create(&habitB)
	db.Create(&models.Streak{HabitID: habitB.ID})

	// Log: HabitA done on March 10 and March 12
	db.Create(&models.HabitLog{HabitID: habitA.ID, UserID: user.ID, Date: "2026-03-10", IsDone: true})
	db.Create(&models.HabitLog{HabitID: habitA.ID, UserID: user.ID, Date: "2026-03-12", IsDone: true})
	// Log: HabitB done on March 12
	db.Create(&models.HabitLog{HabitID: habitB.ID, UserID: user.ID, Date: "2026-03-12", IsDone: true})

	report, err := svc.GenerateWeeklyForPeriod(user.ID, "2026-03-10", "2026-03-12")
	if err != nil {
		t.Fatalf("GenerateWeeklyForPeriod failed: %v", err)
	}

	if len(report.DailyBreakdown) != 3 {
		t.Fatalf("expected 3 days in breakdown, got %d", len(report.DailyBreakdown))
	}

	// March 10: only HabitA existed -> total=1, completed=1 (HabitA done)
	day0 := report.DailyBreakdown[0]
	if day0.Total != 1 {
		t.Errorf("March 10: expected total=1, got %d", day0.Total)
	}
	if day0.Completed != 1 {
		t.Errorf("March 10: expected completed=1, got %d", day0.Completed)
	}

	// March 11: only HabitA existed -> total=1, completed=0
	day1 := report.DailyBreakdown[1]
	if day1.Total != 1 {
		t.Errorf("March 11: expected total=1, got %d", day1.Total)
	}
	if day1.Completed != 0 {
		t.Errorf("March 11: expected completed=0, got %d", day1.Completed)
	}

	// March 12: both HabitA and HabitB exist -> total=2, completed=2
	day2 := report.DailyBreakdown[2]
	if day2.Total != 2 {
		t.Errorf("March 12: expected total=2, got %d", day2.Total)
	}
	if day2.Completed != 2 {
		t.Errorf("March 12: expected completed=2, got %d", day2.Completed)
	}
}
