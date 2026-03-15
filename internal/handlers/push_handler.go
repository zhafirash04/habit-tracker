package handlers

import (
	"errors"
	"net/http"

	"habitflow/internal/services"

	"github.com/gin-gonic/gin"
)

// PushHandler handles push subscription management.
type PushHandler struct {
	Service *services.PushService
}

// NewPushHandler creates a new PushHandler.
func NewPushHandler(service *services.PushService) *PushHandler {
	return &PushHandler{Service: service}
}

// Subscribe handles POST /api/v1/push/subscribe
func (h *PushHandler) Subscribe(c *gin.Context) {
	userID := c.GetUint("user_id")

	var input services.SubscribeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Validasi gagal: " + err.Error(),
			"data":    nil,
		})
		return
	}

	sub, err := h.Service.Subscribe(userID, input)
	if err != nil {
		status := http.StatusInternalServerError
		message := "Gagal menyimpan subscription: " + err.Error()
		if errors.Is(err, services.ErrInvalidSubscription) {
			status = http.StatusBadRequest
			message = err.Error()
		}
		if errors.Is(err, services.ErrSubscriptionLimit) {
			status = http.StatusTooManyRequests
			message = err.Error()
		}
		c.JSON(status, gin.H{
			"success": false,
			"message": message,
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Push subscription berhasil disimpan",
		"data":    sub,
	})
}

// Unsubscribe handles DELETE /api/v1/push/unsubscribe
func (h *PushHandler) Unsubscribe(c *gin.Context) {
	userID := c.GetUint("user_id")

	var body struct {
		Endpoint string `json:"endpoint" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "endpoint diperlukan",
			"data":    nil,
		})
		return
	}

	if err := h.Service.Unsubscribe(userID, body.Endpoint); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal menghapus subscription: " + err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Push subscription berhasil dihapus",
		"data":    nil,
	})
}

// VAPIDKey handles GET /api/v1/push/vapid-key
// Returns the VAPID public key for the frontend to use when subscribing.
func (h *PushHandler) VAPIDKey(c *gin.Context) {
	key := h.Service.GetVAPIDPublicKey()

	if key == "" {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"success": false,
			"message": "VAPID key belum dikonfigurasi",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "VAPID public key",
		"data": gin.H{
			"vapid_public_key": key,
		},
	})
}
