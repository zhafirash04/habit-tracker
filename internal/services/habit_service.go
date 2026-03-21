package services

import (
	"errors"
	"regexp"
	"time"

	"habitflow/internal/models"

	"gorm.io/gorm"
)

// HabitService handles CRUD operations for habits.
type HabitService struct {
	DB *gorm.DB
}

// NewHabitService creates a new HabitService.
func NewHabitService(db *gorm.DB) *HabitService {
	return &HabitService{DB: db}
}

// CreateHabitInput holds data for creating a new habit.
type CreateHabitInput struct {
	Name       string  `json:"name" binding:"required,max=100"`
	Category   string  `json:"category"`
	NotifyTime *string `json:"notify_time"`
}

// UpdateHabitInput holds data for updating an existing habit.
type UpdateHabitInput struct {
	Name       *string `json:"name" binding:"omitempty,max=100"`
	Category   *string `json:"category"`
	NotifyTime *string `json:"notify_time"`
}

// HabitResponse is the flattened habit+streak response returned to the client.
type HabitResponse struct {
	ID            uint      `json:"id"`
	Name          string    `json:"name"`
	Category      string    `json:"category"`
	NotifyTime    *string   `json:"notify_time"`
	IsActive      bool      `json:"is_active"`
	CurrentStreak int       `json:"current_streak"`
	LongestStreak int       `json:"longest_streak"`
	CreatedAt     time.Time `json:"created_at"`
}

// notifyTimeRegex validates HH:MM format (00:00 – 23:59).
var notifyTimeRegex = regexp.MustCompile(`^([01]\d|2[0-3]):([0-5]\d)$`)

// validateNotifyTime checks that notify_time is in HH:MM format.
func validateNotifyTime(t *string) error {
	if t == nil || *t == "" {
		return nil
	}
	if !notifyTimeRegex.MatchString(*t) {
		return errors.New("format notify_time harus HH:MM (contoh: 07:00)")
	}
	return nil
}

// toResponse converts a Habit + its Streak into a flat HabitResponse.
func (s *HabitService) toResponse(habit models.Habit) HabitResponse {
	var streak models.Streak
	s.DB.Where("habit_id = ?", habit.ID).First(&streak)

	return HabitResponse{
		ID:            habit.ID,
		Name:          habit.Name,
		Category:      habit.Category,
		NotifyTime:    habit.NotifyTime,
		IsActive:      habit.IsActive,
		CurrentStreak: streak.CurrentStreak,
		LongestStreak: streak.LongestStreak,
		CreatedAt:     habit.CreatedAt,
	}
}

// Create creates a new habit for a user.
func (s *HabitService) Create(userID uint, input CreateHabitInput) (*HabitResponse, error) {
	// Validate notify_time format
	if err := validateNotifyTime(input.NotifyTime); err != nil {
		return nil, err
	}

	category := input.Category
	if category == "" {
		category = "general"
	}

	habit := models.Habit{
		UserID:     userID,
		Name:       input.Name,
		Category:   category,
		NotifyTime: input.NotifyTime,
		IsActive:   true,
	}

	if err := s.DB.Create(&habit).Error; err != nil {
		return nil, errors.New("gagal membuat habit")
	}

	// Initialize streak record
	streak := models.Streak{
		HabitID:       habit.ID,
		CurrentStreak: 0,
		LongestStreak: 0,
	}
	s.DB.Create(&streak)

	resp := HabitResponse{
		ID:            habit.ID,
		Name:          habit.Name,
		Category:      habit.Category,
		NotifyTime:    habit.NotifyTime,
		IsActive:      habit.IsActive,
		CurrentStreak: 0,
		LongestStreak: 0,
		CreatedAt:     habit.CreatedAt,
	}
	return &resp, nil
}

// GetAll returns all active habits for a user, including streak data.
// Uses a single batch-query for streaks to avoid N+1 problem.
// Before: 1 + N queries (1 habits + N streaks). After: 2 queries total.
func (s *HabitService) GetAll(userID uint) ([]HabitResponse, error) {
	var habits []models.Habit
	if err := s.DB.Where("user_id = ? AND is_active = ?", userID, true).
		Order("created_at DESC").
		Find(&habits).Error; err != nil {
		return nil, errors.New("gagal mengambil daftar habit")
	}

	if len(habits) == 0 {
		return []HabitResponse{}, nil
	}

	// Batch-fetch all streaks in one query instead of one per habit (N+1 fix)
	habitIDs := make([]uint, len(habits))
	for i, h := range habits {
		habitIDs[i] = h.ID
	}
	var streaks []models.Streak
	s.DB.Where("habit_id IN ?", habitIDs).Find(&streaks)

	streakMap := make(map[uint]models.Streak, len(streaks))
	for _, st := range streaks {
		streakMap[st.HabitID] = st
	}

	responses := make([]HabitResponse, 0, len(habits))
	for _, habit := range habits {
		st := streakMap[habit.ID] // zero-value if not found
		responses = append(responses, HabitResponse{
			ID:            habit.ID,
			Name:          habit.Name,
			Category:      habit.Category,
			NotifyTime:    habit.NotifyTime,
			IsActive:      habit.IsActive,
			CurrentStreak: st.CurrentStreak,
			LongestStreak: st.LongestStreak,
			CreatedAt:     habit.CreatedAt,
		})
	}

	return responses, nil
}

// GetByID returns a single habit by ID (scoped to user), including streak data.
func (s *HabitService) GetByID(userID, habitID uint) (*HabitResponse, error) {
	var habit models.Habit
	if err := s.DB.Where("id = ? AND user_id = ?", habitID, userID).First(&habit).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("habit tidak ditemukan")
		}
		return nil, errors.New("gagal mengambil habit")
	}

	resp := s.toResponse(habit)
	return &resp, nil
}

// Update modifies an existing habit.
func (s *HabitService) Update(userID, habitID uint, input UpdateHabitInput) (*HabitResponse, error) {
	// Verify ownership
	var habit models.Habit
	if err := s.DB.Where("id = ? AND user_id = ?", habitID, userID).First(&habit).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("habit tidak ditemukan")
		}
		return nil, errors.New("gagal mengambil habit")
	}

	// Validate notify_time if provided
	if input.NotifyTime != nil {
		if err := validateNotifyTime(input.NotifyTime); err != nil {
			return nil, err
		}
	}

	updates := make(map[string]interface{})
	if input.Name != nil {
		updates["name"] = *input.Name
	}
	if input.Category != nil {
		updates["category"] = *input.Category
	}
	if input.NotifyTime != nil {
		updates["notify_time"] = *input.NotifyTime
	}

	if len(updates) > 0 {
		if err := s.DB.Model(&habit).Updates(updates).Error; err != nil {
			return nil, errors.New("gagal mengupdate habit")
		}
	}

	// Reload
	s.DB.First(&habit, habit.ID)
	resp := s.toResponse(habit)
	return &resp, nil
}

// Delete performs a soft delete: sets is_active = false.
func (s *HabitService) Delete(userID, habitID uint) error {
	var habit models.Habit
	if err := s.DB.Where("id = ? AND user_id = ? AND is_active = ?", habitID, userID, true).
		First(&habit).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("habit tidak ditemukan")
		}
		return errors.New("gagal mengambil habit")
	}

	if err := s.DB.Model(&habit).Update("is_active", false).Error; err != nil {
		return errors.New("gagal menghapus habit")
	}

	return nil
}
