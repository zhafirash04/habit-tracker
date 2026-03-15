package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// BodySizeLimitMiddleware rejects requests over the configured size and
// wraps the body reader so JSON binders cannot read beyond that limit.
func BodySizeLimitMiddleware(maxBytes int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.ContentLength > maxBytes {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{
				"success": false,
				"message": "Payload terlalu besar",
				"data":    nil,
			})
			c.Abort()
			return
		}

		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes)
		c.Next()
	}
}

// RequireJSONMiddleware enforces application/json for API write operations.
// IMPORTANT: OPTIONS requests are skipped to allow CORS preflight to succeed.
func RequireJSONMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		// Skip OPTIONS (preflight) - handled by CORSMiddleware
		if method == http.MethodOptions {
			c.Next()
			return
		}

		if method != http.MethodPost && method != http.MethodPut && method != http.MethodPatch && method != http.MethodDelete {
			c.Next()
			return
		}

		contentType := strings.ToLower(strings.TrimSpace(c.GetHeader("Content-Type")))
		contentLength := c.Request.ContentLength
		hasBody := contentLength > 0
		if contentLength == -1 {
			hasBody = strings.TrimSpace(c.GetHeader("Transfer-Encoding")) != ""
		}

		if hasBody && !strings.HasPrefix(contentType, "application/json") {
			c.JSON(http.StatusUnsupportedMediaType, gin.H{
				"success": false,
				"message": "Content-Type harus application/json",
				"data":    nil,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
