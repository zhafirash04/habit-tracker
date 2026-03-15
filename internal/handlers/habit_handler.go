package handlers

import (
	"net/http"
	"strconv"

	"habitflow/internal/services"

	"github.com/gin-gonic/gin"
)

// HabitHandler handles habit CRUD requests.
type HabitHandler struct {
	Service *services.HabitService
}

// NewHabitHandler creates a new HabitHandler.
func NewHabitHandler(service *services.HabitService) *HabitHandler {
	return &HabitHandler{Service: service}
}

// Create handles POST /api/v1/habits
func (h *HabitHandler) Create(c *gin.Context) {
	userID := c.GetUint("user_id")

	var input services.CreateHabitInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Validasi gagal: " + err.Error(),
			"data":    nil,
		})
		return
	}

	habit, err := h.Service.Create(userID, input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Habit berhasil dibuat",
		"data":    habit,
	})
}

// GetAll handles GET /api/v1/habits
func (h *HabitHandler) GetAll(c *gin.Context) {
	userID := c.GetUint("user_id")

	habits, err := h.Service.GetAll(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Daftar habit berhasil diambil",
		"data":    habits,
	})
}

// GetByID handles GET /api/v1/habits/:id
func (h *HabitHandler) GetByID(c *gin.Context) {
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

	habit, err := h.Service.GetByID(userID, uint(habitID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Detail habit berhasil diambil",
		"data":    habit,
	})
}

// Update handles PUT /api/v1/habits/:id
func (h *HabitHandler) Update(c *gin.Context) {
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

	var input services.UpdateHabitInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Validasi gagal: " + err.Error(),
			"data":    nil,
		})
		return
	}

	habit, err := h.Service.Update(userID, uint(habitID), input)
	if err != nil {
		status := http.StatusNotFound
		if err.Error() == "format notify_time harus HH:MM (contoh: 07:00)" {
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Habit berhasil diupdate",
		"data":    habit,
	})
}

// Delete handles DELETE /api/v1/habits/:id
// Soft delete: sets is_active = false
func (h *HabitHandler) Delete(c *gin.Context) {
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

	if err := h.Service.Delete(userID, uint(habitID)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Habit berhasil dihapus",
		"data":    nil,
	})
}
