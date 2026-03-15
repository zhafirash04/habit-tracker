package middleware

import (
	"net/http"
	"strings"

	"habitflow/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware validates JWT tokens and injects user_id and email into the Gin context.
func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Header Authorization diperlukan",
				"data":    nil,
			})
			c.Abort()
			return
		}

		// Expect "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Format header harus 'Bearer <token>'",
				"data":    nil,
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(cfg.JWTSecret), nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Token tidak valid atau sudah kadaluarsa",
				"data":    nil,
			})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Token claims tidak valid",
				"data":    nil,
			})
			c.Abort()
			return
		}

		// Ensure this is an access token, not a refresh token
		tokenType, _ := claims["type"].(string)
		if tokenType != "access" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Gunakan access token, bukan refresh token",
				"data":    nil,
			})
			c.Abort()
			return
		}

		userID, ok := claims["user_id"].(float64)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "user_id tidak valid dalam token",
				"data":    nil,
			})
			c.Abort()
			return
		}

		email, _ := claims["email"].(string)

		c.Set("user_id", uint(userID))
		c.Set("email", email)
		c.Next()
	}
}
