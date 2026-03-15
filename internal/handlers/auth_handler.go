package handlers

import (
	"net/http"
	"strings"
	"time"

	"habitflow/internal/services"

	"github.com/gin-gonic/gin"
)

const refreshCookieName = "hf_refresh_token"

// AuthHandler handles authentication-related requests.
type AuthHandler struct {
	Service *services.AuthService
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(service *services.AuthService) *AuthHandler {
	return &AuthHandler{Service: service}
}

// Register handles POST /api/v1/auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var input services.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Validasi gagal: " + err.Error(),
			"data":    nil,
		})
		return
	}

	result, err := h.Service.Register(input)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	h.setRefreshCookie(c, result.Tokens.RefreshToken)

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Registrasi berhasil",
		"data": gin.H{
			"user":         result.User,
			"access_token": result.Tokens.AccessToken,
			"expires_in":   result.Tokens.ExpiresIn,
		},
	})
}

// Login handles POST /api/v1/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var input services.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Validasi gagal: " + err.Error(),
			"data":    nil,
		})
		return
	}

	result, err := h.Service.Login(input)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	h.setRefreshCookie(c, result.Tokens.RefreshToken)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Login berhasil",
		"data": gin.H{
			"user":         result.User,
			"access_token": result.Tokens.AccessToken,
			"expires_in":   result.Tokens.ExpiresIn,
		},
	})
}

// Refresh handles POST /api/v1/auth/refresh
// Expects the refresh token in HttpOnly cookie or Authorization header.
func (h *AuthHandler) Refresh(c *gin.Context) {
	refreshToken, ok := h.extractRefreshToken(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Refresh token diperlukan",
			"data":    nil,
		})
		return
	}
	tokens, err := h.Service.RefreshToken(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	h.setRefreshCookie(c, tokens.RefreshToken)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Token berhasil diperbarui",
		"data": gin.H{
			"access_token": tokens.AccessToken,
			"expires_in":   tokens.ExpiresIn,
		},
	})
}

// Logout handles POST /api/v1/auth/logout
// Expects refresh token in HttpOnly cookie or Authorization header.
func (h *AuthHandler) Logout(c *gin.Context) {
	refreshToken, ok := h.extractRefreshToken(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Refresh token diperlukan",
			"data":    nil,
		})
		return
	}

	if err := h.Service.RevokeRefreshToken(refreshToken); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	h.clearRefreshCookie(c)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Logout berhasil",
		"data":    nil,
	})
}

func (h *AuthHandler) extractRefreshToken(c *gin.Context) (string, bool) {
	if cookieToken, err := c.Cookie(refreshCookieName); err == nil && strings.TrimSpace(cookieToken) != "" {
		return cookieToken, true
	}

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", false
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" || strings.TrimSpace(parts[1]) == "" {
		return "", false
	}

	return parts[1], true
}

func (h *AuthHandler) setRefreshCookie(c *gin.Context, refreshToken string) {
	secure := h.Service.Cfg.Environment == "production"
	// Use SameSiteNoneMode for cross-origin requests (PWA on different domain).
	// SameSiteLaxMode blocks cookies on cross-origin POST, causing "Failed to fetch".
	sameSite := http.SameSiteNoneMode
	if !secure {
		// SameSite=None requires Secure=true; fallback to Lax for local HTTP dev.
		sameSite = http.SameSiteLaxMode
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     refreshCookieName,
		Value:    refreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   secure,
		SameSite: sameSite,
		MaxAge:   int((7 * 24 * time.Hour).Seconds()),
	})
}

func (h *AuthHandler) clearRefreshCookie(c *gin.Context) {
	secure := h.Service.Cfg.Environment == "production"
	sameSite := http.SameSiteNoneMode
	if !secure {
		sameSite = http.SameSiteLaxMode
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     refreshCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   secure,
		SameSite: sameSite,
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
	})
}
