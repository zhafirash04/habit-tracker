package middleware

import (
	"net/http"
	"strings"

	"habitflow/internal/config"

	"github.com/gin-gonic/gin"
)

// CORSMiddleware handles Cross-Origin Resource Sharing.
// In development, it allows all localhost origins.
func CORSMiddleware(cfg ...*config.Config) gin.HandlerFunc {
	var conf *config.Config
	if len(cfg) > 0 {
		conf = cfg[0]
	}

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		allowed := isOriginAllowed(origin, conf)

		if allowed {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Vary", "Origin")
		}

		allowMethods := "GET, POST, PUT, PATCH, DELETE, OPTIONS"
		allowHeaders := "Origin, Content-Type, Accept, Authorization"
		if reqHeaders := strings.TrimSpace(c.GetHeader("Access-Control-Request-Headers")); reqHeaders != "" {
			allowHeaders = reqHeaders
		}

		c.Header("Access-Control-Allow-Methods", allowMethods)
		c.Header("Access-Control-Allow-Headers", allowHeaders)
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == http.MethodOptions {
			// Always short-circuit preflight so OPTIONS is never blocked by downstream middleware.
			// Browser will still enforce origin by presence/absence of Allow-Origin header.
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func isOriginAllowed(origin string, cfg *config.Config) bool {
	if origin == "" {
		return false
	}

	if cfg != nil && len(cfg.CORSOrigins) > 0 {
		for _, allowed := range cfg.CORSOrigins {
			if strings.EqualFold(strings.TrimSpace(allowed), origin) {
				return true
			}
		}
		return false
	}

	return strings.HasPrefix(origin, "http://localhost") ||
		strings.HasPrefix(origin, "http://127.0.0.1") ||
		strings.EqualFold(origin, "https://zhafirash04.github.io")
}
