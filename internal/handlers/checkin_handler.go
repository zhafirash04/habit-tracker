package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"habitflow/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CheckinHandler handles daily habit check-in requests.
type CheckinHandler struct {
	StreakService *services.StreakService
}

// NewCheckinHandler creates a new CheckinHandler.
func NewCheckinHandler(streakService *services.StreakService) *CheckinHandler {
	return &CheckinHandler{
		StreakService: streakService,
	}
}

// checkInput holds the optional request body for a check-in.
type checkInput struct {
	Note *string `json:"note"` // optional one-line note
}

// Check handles POST /api/v1/habits/:id/check
func (h *CheckinHandler) Check(c *gin.Context) {
	userID := c.GetUint("user_id")
	habitID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID habit tidak valid",
			"data":    nil,
		})
		return
	}

	var input checkInput
	c.ShouldBindJSON(&input) // optional body

	result, err := h.StreakService.Checkin(userID, uint(habitID), input.Note)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrAlreadyCheckedIn):
			c.JSON(http.StatusConflict, gin.H{
				"success": false,
				"message": err.Error(),
				"data":    nil,
			})
		case errors.Is(err, gorm.ErrRecordNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Habit tidak ditemukan",
				"data":    nil,
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Gagal melakukan checkin: " + err.Error(),
				"data":    nil,
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Habit berhasil dicheckin!",
		"data":    result,
	})
}

// Undo handles DELETE /api/v1/habits/:id/check
func (h *CheckinHandler) Undo(c *gin.Context) {
	userID := c.GetUint("user_id")
	habitID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID habit tidak valid",
			"data":    nil,
		})
		return
	}

	err = h.StreakService.UndoCheckin(userID, uint(habitID))
	if err != nil {
		if errors.Is(err, services.ErrNoCheckinToday) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": err.Error(),
				"data":    nil,
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Gagal membatalkan checkin: " + err.Error(),
				"data":    nil,
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Checkin hari ini berhasil dibatalkan",
		"data":    nil,
	})
}

// Today handles GET /api/v1/habits/today
func (h *CheckinHandler) Today(c *gin.Context) {
	userID := c.GetUint("user_id")

	statuses, err := h.StreakService.GetTodayStatus(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil status hari ini: " + err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Status habit hari ini",
		"data": gin.H{
			"date":   services.TodayWIB(),
			"habits": statuses,
		},
	})
}
