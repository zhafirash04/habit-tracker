package models

import "time"

// RefreshToken stores issued refresh tokens so they can be rotated and revoked.
type RefreshToken struct {
	ID            uint       `json:"id" gorm:"primaryKey"`
	UserID        uint       `json:"user_id" gorm:"not null;index"`
	JTI           string     `json:"jti" gorm:"size:64;not null;uniqueIndex"`
	TokenHash     string     `json:"-" gorm:"size:64;not null;index"`
	ExpiresAt     time.Time  `json:"expires_at" gorm:"not null;index"`
	RevokedAt     *time.Time `json:"revoked_at"`
	ReplacedByJTI *string    `json:"replaced_by_jti" gorm:"size:64"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`

	// Relations
	User User `json:"-" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}
