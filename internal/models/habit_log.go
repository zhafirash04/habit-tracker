package models

import "time"

// HabitLog represents a daily check-in record for a habit.
type HabitLog struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	HabitID     uint       `json:"habit_id" gorm:"not null;uniqueIndex:idx_habit_date"`
	UserID      uint       `json:"user_id" gorm:"not null;index"`
	Date        string     `json:"date" gorm:"not null;uniqueIndex:idx_habit_date;size:10"` // YYYY-MM-DD
	IsDone      bool       `json:"is_done" gorm:"default:false"`
	Note        *string    `json:"note" gorm:"size:200"` // optional one-line note
	CompletedAt *time.Time `json:"completed_at"`

	// Relations
	Habit Habit `json:"-" gorm:"foreignKey:HabitID;constraint:OnDelete:CASCADE"`
	User  User  `json:"-" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}
