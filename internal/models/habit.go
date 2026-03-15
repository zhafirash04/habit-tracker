package models

import "time"

// Habit represents a user's tracked habit.
type Habit struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	UserID     uint      `json:"user_id" gorm:"not null;index"`
	Name       string    `json:"name" gorm:"not null;size:100"`
	Category   string    `json:"category" gorm:"default:'general'"`
	NotifyTime *string   `json:"notify_time" gorm:"size:5"` // format "HH:MM", nullable
	IsActive   bool      `json:"is_active" gorm:"default:true"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	// Relations
	User User `json:"-" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}
