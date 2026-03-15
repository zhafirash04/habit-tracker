package models

import "time"

// PushSubscription stores a user's Web Push subscription details.
type PushSubscription struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null;index"`
	Endpoint  string    `json:"endpoint" gorm:"type:text;not null"`
	P256dh    string    `json:"p256dh" gorm:"type:text;not null"`
	Auth      string    `json:"auth" gorm:"type:text;not null"`
	CreatedAt time.Time `json:"created_at"`

	// Relations
	User User `json:"-" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}
