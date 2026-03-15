package models

// Streak tracks current and longest streaks for a habit.
type Streak struct {
	ID            uint    `json:"id" gorm:"primaryKey"`
	HabitID       uint    `json:"habit_id" gorm:"uniqueIndex;not null"`
	CurrentStreak int     `json:"current_streak" gorm:"default:0"`
	LongestStreak int     `json:"longest_streak" gorm:"default:0"`
	LastDoneDate  *string `json:"last_done_date" gorm:"size:10"` // YYYY-MM-DD, nullable

	// Relations
	Habit Habit `json:"-" gorm:"foreignKey:HabitID;constraint:OnDelete:CASCADE"`
}
