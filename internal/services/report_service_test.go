package services

import (
	"testing"
	"time"

	"habitflow/internal/models"
)

// TestWeeklyReport_DailyBreakdown_RespectsCreatedAt verifies that the daily
// breakdown only counts habits that existed on each specific day.
// A habit created mid-week should NOT inflate the Total for earlier days.
func TestWeeklyReport_DailyBreakdown_RespectsCreatedAt(t *testing.T) {
	db := setupTestDB(t)
	scoreSvc := NewScoreService(db)
	insightSvc := NewInsightService(db)
	svc := NewReportService(db, scoreSvc, insightSvc)

	user := models.User{Name: "report_user", Email: "report_user@test.com", PasswordHash: "hash"}
	db.Create(&user)

	// Habit A: created on March 9 — should appear in all 3 days (Mar 9-11)
	habitA := models.Habit{
		UserID:    user.ID,
		Name:      "Habit A",
		Category:  "health",
		IsActive:  true,
		CreatedAt: time.Date(2026, 3, 9, 8, 0, 0, 0, WIB),
	}
	db.Create(&habitA)
	db.Create(&models.Streak{HabitID: habitA.ID})

	// Habit B: created on March 11 — should only appear on Mar 11
	habitB := models.Habit{
		UserID:    user.ID,
		Name:      "Habit B",
		Category:  "learning",
		IsActive:  true,
		CreatedAt: time.Date(2026, 3, 11, 9, 0, 0, 0, WIB),
	}
	db.Create(&habitB)
	db.Create(&models.Streak{HabitID: habitB.ID})

	// Check-in Habit A on all 3 days
	insertLog(t, db, habitA.ID, user.ID, "2026-03-09", true)
	insertLog(t, db, habitA.ID, user.ID, "2026-03-10", true)
	insertLog(t, db, habitA.ID, user.ID, "2026-03-11", true)

	// Check-in Habit B on March 11 only
	insertLog(t, db, habitB.ID, user.ID, "2026-03-11", true)

	report, err := svc.GenerateWeeklyForPeriod(user.ID, "2026-03-09", "2026-03-11")
	if err != nil {
		t.Fatalf("GenerateWeeklyForPeriod failed: %v", err)
	}

	if len(report.DailyBreakdown) != 3 {
		t.Fatalf("expected 3 days in breakdown, got %d", len(report.DailyBreakdown))
	}

	// March 9: only Habit A existed → Total=1, Completed=1
	day0 := report.DailyBreakdown[0]
	if day0.Date != "2026-03-09" {
		t.Fatalf("expected first day=2026-03-09, got %s", day0.Date)
	}
	if day0.Total != 1 {
		t.Errorf("March 9: expected total=1, got %d", day0.Total)
	}
	if day0.Completed != 1 {
		t.Errorf("March 9: expected completed=1, got %d", day0.Completed)
	}
	if day0.Rate != 100 {
		t.Errorf("March 9: expected rate=100, got %.2f", day0.Rate)
	}

	// March 10: only Habit A existed → Total=1, Completed=1
	day1 := report.DailyBreakdown[1]
	if day1.Date != "2026-03-10" {
		t.Fatalf("expected second day=2026-03-10, got %s", day1.Date)
	}
	if day1.Total != 1 {
		t.Errorf("March 10: expected total=1, got %d", day1.Total)
	}
	if day1.Completed != 1 {
		t.Errorf("March 10: expected completed=1, got %d", day1.Completed)
	}

	// March 11: both habits existed → Total=2, Completed=2
	day2 := report.DailyBreakdown[2]
	if day2.Date != "2026-03-11" {
		t.Fatalf("expected third day=2026-03-11, got %s", day2.Date)
	}
	if day2.Total != 2 {
		t.Errorf("March 11: expected total=2, got %d", day2.Total)
	}
	if day2.Completed != 2 {
		t.Errorf("March 11: expected completed=2, got %d", day2.Completed)
	}
	if day2.Rate != 100 {
		t.Errorf("March 11: expected rate=100, got %.2f", day2.Rate)
	}
}

// TestWeeklyReport_DailyBreakdown_NoHabitsOnDay verifies that days before any
// habit was created show Total=0 and Rate=0.
func TestWeeklyReport_DailyBreakdown_NoHabitsOnDay(t *testing.T) {
	db := setupTestDB(t)
	scoreSvc := NewScoreService(db)
	insightSvc := NewInsightService(db)
	svc := NewReportService(db, scoreSvc, insightSvc)

	user := models.User{Name: "no_habit_day", Email: "no_habit_day@test.com", PasswordHash: "hash"}
	db.Create(&user)

	// Habit created on March 12
	habit := models.Habit{
		UserID:    user.ID,
		Name:      "Late Habit",
		Category:  "health",
		IsActive:  true,
		CreatedAt: time.Date(2026, 3, 12, 10, 0, 0, 0, WIB),
	}
	db.Create(&habit)
	db.Create(&models.Streak{HabitID: habit.ID})

	report, err := svc.GenerateWeeklyForPeriod(user.ID, "2026-03-10", "2026-03-12")
	if err != nil {
		t.Fatalf("GenerateWeeklyForPeriod failed: %v", err)
	}

	// March 10: no habits existed yet → Total=0
	if report.DailyBreakdown[0].Total != 0 {
		t.Errorf("March 10: expected total=0, got %d", report.DailyBreakdown[0].Total)
	}
	if report.DailyBreakdown[0].Rate != 0 {
		t.Errorf("March 10: expected rate=0, got %.2f", report.DailyBreakdown[0].Rate)
	}

	// March 11: still no habits → Total=0
	if report.DailyBreakdown[1].Total != 0 {
		t.Errorf("March 11: expected total=0, got %d", report.DailyBreakdown[1].Total)
	}

	// March 12: habit exists → Total=1
	if report.DailyBreakdown[2].Total != 1 {
		t.Errorf("March 12: expected total=1, got %d", report.DailyBreakdown[2].Total)
	}
}
