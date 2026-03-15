package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// visitor tracks the last seen time and request count for rate limiting.
type visitor struct {
	lastSeen time.Time
	tokens   float64
}

// RateLimiter implements a per-IP token bucket rate limiter.
type RateLimiter struct {
	mu         sync.Mutex
	visitors   map[string]*visitor
	rate       float64 // tokens per second
	burst      int     // max tokens
	cleanupInt time.Duration
}

// NewRateLimiter creates a rate limiter with the given rate (requests/sec) and burst size.
func NewRateLimiter(rate float64, burst int) *RateLimiter {
	rl := &RateLimiter{
		visitors:   make(map[string]*visitor),
		rate:       rate,
		burst:      burst,
		cleanupInt: 5 * time.Minute,
	}

	go rl.cleanup()
	return rl
}

// cleanup removes stale visitors periodically.
func (rl *RateLimiter) cleanup() {
	for {
		time.Sleep(rl.cleanupInt)
		rl.mu.Lock()
		for ip, v := range rl.visitors {
			if time.Since(v.lastSeen) > rl.cleanupInt {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// allow checks whether a given IP is allowed to make a request.
func (rl *RateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[ip]
	now := time.Now()

	if !exists {
		rl.visitors[ip] = &visitor{
			lastSeen: now,
			tokens:   float64(rl.burst) - 1,
		}
		return true
	}

	// Refill tokens based on elapsed time
	elapsed := now.Sub(v.lastSeen).Seconds()
	v.tokens += elapsed * rl.rate
	if v.tokens > float64(rl.burst) {
		v.tokens = float64(rl.burst)
	}
	v.lastSeen = now

	if v.tokens < 1 {
		return false
	}

	v.tokens--
	return true
}

// RateLimitMiddleware returns a Gin middleware that enforces general rate limiting.
// Default: 10 requests/second with a burst of 20.
func RateLimitMiddleware() gin.HandlerFunc {
	limiter := NewRateLimiter(10, 20)

	return func(c *gin.Context) {
		ip := c.ClientIP()

		if !limiter.allow(ip) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"message": "Terlalu banyak request, coba lagi nanti",
				"data":    nil,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// LoginRateLimitMiddleware returns a Gin middleware that enforces strict rate limiting
// specifically for login attempts: max 5 attempts per IP per minute.
func LoginRateLimitMiddleware() gin.HandlerFunc {
	// 5 tokens per 60 seconds = ~0.0833 tokens/sec, burst of 5
	limiter := NewRateLimiter(5.0/60.0, 5)

	return func(c *gin.Context) {
		ip := c.ClientIP()

		if !limiter.allow(ip) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"message": "Terlalu banyak percobaan login, coba lagi dalam 1 menit",
				"data":    nil,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
