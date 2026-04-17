package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"bn-mobile/server/internal/delivery/data/response"

	"github.com/gin-gonic/gin"
)

type clientWindow struct {
	Count       int
	WindowStart time.Time
}

type IPRateLimiter struct {
	window time.Duration
	limit  int
	mu     sync.Mutex
	store  map[string]clientWindow
}

func NewIPRateLimiter(limit int, window time.Duration) *IPRateLimiter {
	if limit <= 0 {
		limit = 1
	}
	if window <= 0 {
		window = time.Minute
	}

	return &IPRateLimiter{
		window: window,
		limit:  limit,
		store:  make(map[string]clientWindow),
	}
}

func (l *IPRateLimiter) Middleware(scope string) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := fmt.Sprintf("%s:%s", scope, c.ClientIP())
		now := time.Now()

		l.mu.Lock()
		entry, exists := l.store[key]
		if !exists || now.Sub(entry.WindowStart) >= l.window {
			entry = clientWindow{Count: 0, WindowStart: now}
		}

		entry.Count++
		l.store[key] = entry
		l.cleanupStale(now)
		remaining := l.limit - entry.Count
		windowEnd := entry.WindowStart.Add(l.window)
		resetAfter := maxInt(0, int(time.Until(windowEnd).Seconds()))
		l.mu.Unlock()

		c.Header("X-RateLimit-Limit", strconv.Itoa(l.limit))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(maxInt(0, remaining)))
		c.Header("X-RateLimit-Reset", strconv.Itoa(resetAfter))

		if entry.Count > l.limit {
			c.Header("Retry-After", strconv.Itoa(resetAfter))
			response.Failed(c, http.StatusTooManyRequests, "too many requests", "rate limit exceeded")
			c.Abort()
			return
		}

		c.Next()
	}
}

func (l *IPRateLimiter) cleanupStale(now time.Time) {
	if len(l.store) <= 1024 {
		return
	}

	for key, value := range l.store {
		if now.Sub(value.WindowStart) > 2*l.window {
			delete(l.store, key)
		}
	}
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
