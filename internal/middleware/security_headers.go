package middleware

import (
	"strings"

	"habitflow/internal/config"

	"github.com/gin-gonic/gin"
)

// SecurityHeadersMiddleware sets baseline security headers for browser clients.
func SecurityHeadersMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Permissions-Policy", "camera=(), microphone=(), geolocation=()")

		csp := strings.Join([]string{
			"default-src 'self'",
			"base-uri 'self'",
			"frame-ancestors 'none'",
			"object-src 'none'",
			"img-src 'self' data:",
			"style-src 'self' 'unsafe-inline' https://fonts.googleapis.com https://fonts.gstatic.com https://cdn.tailwindcss.com",
			"font-src 'self' https://fonts.gstatic.com",
			"script-src 'self' 'unsafe-inline' https://cdn.tailwindcss.com https://fonts.googleapis.com https://fonts.gstatic.com",
			"connect-src 'self' https://fonts.googleapis.com https://fonts.gstatic.com https://cdn.tailwindcss.com",
			"worker-src 'self'",
		}, "; ")
		c.Header("Content-Security-Policy", csp)

		if cfg != nil && cfg.Environment == "production" {
			if strings.EqualFold(c.GetHeader("X-Forwarded-Proto"), "https") || c.Request.TLS != nil {
				c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
			}
		}

		c.Next()
	}
}
