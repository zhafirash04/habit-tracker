package services

import "errors"

// Sentinel errors for check-in operations.
var (
	ErrAlreadyCheckedIn    = errors.New("habit sudah dicheckin hari ini")
	ErrNoCheckinToday      = errors.New("tidak ada checkin hari ini untuk di-undo")
	ErrInvalidSubscription = errors.New("subscription push tidak valid")
	ErrSubscriptionLimit   = errors.New("batas subscription perangkat tercapai")
)
