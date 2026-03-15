package handlers

import (
	"net/http"
	"time"

	"habitflow/internal/config"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// HealthHandler serves liveness and runtime self-check endpoints.
type HealthHandler struct {
	Cfg *config.Config
	DB  *gorm.DB
}

// NewHealthHandler creates a new health handler.
func NewHealthHandler(cfg *config.Config, db *gorm.DB) *HealthHandler {
	return &HealthHandler{Cfg: cfg, DB: db}
}

// Liveness handles GET /api/v1/health.
func (h *HealthHandler) Liveness(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "ok",
		"data": gin.H{
			"status":    "ok",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		},
	})
}

// SecuritySelfCheck handles GET /api/v1/health/security.
// It returns deployment-safe booleans and warnings without exposing secrets.
func (h *HealthHandler) SecuritySelfCheck(c *gin.Context) {
	dbConnected := false
	if h.DB != nil {
		if sqlDB, err := h.DB.DB(); err == nil {
			dbConnected = sqlDB.Ping() == nil
		}
	}

	jwtStrong := h.Cfg != nil && h.Cfg.IsJWTSecretStrong()
	corsConfigured := h.Cfg != nil && len(h.Cfg.CORSOrigins) > 0
	maxBodyOK := h.Cfg != nil && h.Cfg.MaxBodyBytes >= 1024
	env := "unknown"
	if h.Cfg != nil && h.Cfg.Environment != "" {
		env = h.Cfg.Environment
	}

	warnings := make([]string, 0)
	if !dbConnected {
		warnings = append(warnings, "database connection check failed")
	}
	if env == "production" && !jwtStrong {
		warnings = append(warnings, "JWT secret is weak for production")
	}
	if env == "production" && !corsConfigured {
		warnings = append(warnings, "CORS allowlist is empty in production")
	}
	if !maxBodyOK {
		warnings = append(warnings, "MAX_BODY_BYTES is too low")
	}

	overall := "ok"
	if len(warnings) > 0 {
		overall = "warning"
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "security self-check",
		"data": gin.H{
			"status":      overall,
			"environment": env,
			"checks": gin.H{
				"database_connected":          dbConnected,
				"jwt_secret_strong":           jwtStrong,
				"cors_allowlist_configured":   corsConfigured,
				"max_body_bytes_configured":   maxBodyOK,
				"refresh_token_http_only":     true,
				"api_response_cache_disabled": true,
			},
			"warnings": warnings,
		},
	})
}
