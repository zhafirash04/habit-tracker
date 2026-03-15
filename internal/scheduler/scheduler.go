package scheduler

import (
	"fmt"
	"log"
	"time"

	"habitflow/internal/models"
	"habitflow/internal/services"

	"github.com/go-co-op/gocron/v2"
	"gorm.io/gorm"
)

// Scheduler manages scheduled jobs.
type Scheduler struct {
	scheduler     gocron.Scheduler
	db            *gorm.DB
	pushService   *services.PushService
	reportService *services.ReportService
}

// New creates and configures a new Scheduler.
func New(db *gorm.DB, pushService *services.PushService, reportService *services.ReportService) *Scheduler {
	s, err := gocron.NewScheduler()
	if err != nil {
		log.Fatalf("Failed to create scheduler: %v", err)
	}

	return &Scheduler{
		scheduler:     s,
		db:            db,
		pushService:   pushService,
		reportService: reportService,
	}
}

// Start begins all scheduled jobs.
func (s *Scheduler) Start() {
	// Habit reminder — runs every minute, checks for notify_time matches in WIB
	_, err := s.scheduler.NewJob(
		gocron.DurationJob(1*time.Minute),
		gocron.NewTask(s.sendHabitReminders),
	)
	if err != nil {
		log.Printf("Failed to schedule habit reminders: %v", err)
	}

	// Weekly report — every Monday at 08:00 WIB
	_, err = s.scheduler.NewJob(
		gocron.CronJob("0 1 * * 1", false), // 01:00 UTC = 08:00 WIB
		gocron.NewTask(s.generateWeeklyReports),
	)
	if err != nil {
		log.Printf("Failed to schedule weekly reports: %v", err)
	}

	s.scheduler.Start()
	log.Println("Scheduler started — habit reminders (every 1m) + weekly reports (Mon 08:00 WIB)")
}

// Stop gracefully shuts down the scheduler.
func (s *Scheduler) Stop() {
	if err := s.scheduler.Shutdown(); err != nil {
		log.Printf("Error shutting down scheduler: %v", err)
	}
}

// sendHabitReminders checks habits with notify_time matching current WIB HH:MM
// and sends push notifications for habits not yet checked in today.
func (s *Scheduler) sendHabitReminders() {
	now := time.Now().In(services.WIB)
	currentTime := now.Format("15:04")
	today := now.Format("2006-01-02")

	// Find active habits with matching notify_time
	var habits []models.Habit
	s.db.Where("notify_time = ? AND is_active = ?", currentTime, true).Find(&habits)

	if len(habits) == 0 {
		return
	}

	for _, habit := range habits {
		// Skip if already done today
		var count int64
		s.db.Model(&models.HabitLog{}).
			Where("habit_id = ? AND date = ? AND is_done = ?", habit.ID, today, true).
			Count(&count)

		if count > 0 {
			continue
		}

		// Get current streak for the notification body
		var streak models.Streak
		streakText := ""
		if err := s.db.Where("habit_id = ?", habit.ID).First(&streak).Error; err == nil && streak.CurrentStreak > 0 {
			streakText = fmt.Sprintf(" Streak kamu: %d hari", streak.CurrentStreak)
		}

		// Send push notification
		s.pushService.SendToUser(habit.UserID, services.PushPayload{
			Title: "HabitFlow",
			Body:  fmt.Sprintf("Waktunya: %s!%s", habit.Name, streakText),
			Icon:  "/icons/icon-192.png",
			Badge: "/icons/badge-72.png",
			Data: services.PushData{
				HabitID: habit.ID,
				URL:     "/habits",
			},
		})

		log.Printf("Reminder sent: habit=%d user=%d name=%s", habit.ID, habit.UserID, habit.Name)
	}
}

// generateWeeklyReports sends weekly report notifications to all users.
func (s *Scheduler) generateWeeklyReports() {
	var users []models.User
	s.db.Find(&users)

	for _, user := range users {
		report, err := s.reportService.GenerateWeekly(user.ID)
		if err != nil {
			log.Printf("Failed to generate weekly report for user %d: %v", user.ID, err)
			continue
		}

		s.pushService.SendToUser(user.ID, services.PushPayload{
			Title: "Laporan Mingguan HabitFlow",
			Body:  fmt.Sprintf("Skor konsistensi minggu ini: %.0f%%. Kamu menyelesaikan %d check-in!", report.Score.Overall, report.TotalCheckin),
			Icon:  "/icons/icon-192.png",
			Badge: "/icons/badge-72.png",
			Data: services.PushData{
				URL: "/reports",
			},
		})
	}

	log.Printf("Weekly reports sent to %d users", len(users))
}
