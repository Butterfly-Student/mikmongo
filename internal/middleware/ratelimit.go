package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"mikmongo/pkg/redis"
	"mikmongo/pkg/response"
)

// RateLimitMiddleware handles rate limiting
type RateLimitMiddleware struct {
	limiter *redis.RateLimiter
}

// NewRateLimitMiddleware creates a new rate limit middleware
func NewRateLimitMiddleware(client *redis.Client) *RateLimitMiddleware {
	return &RateLimitMiddleware{
		limiter: redis.NewRateLimiter(client),
	}
}

// Limit creates a rate limit middleware with specified limit and window
func (m *RateLimitMiddleware) Limit(limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.ClientIP()
		if userID, exists := c.Get("user_id"); exists {
			key = userID.(string)
		}

		allowed, err := m.limiter.IsAllowed(c.Request.Context(), key, limit, window)
		if err != nil || !allowed {
			response.Error(c, http.StatusTooManyRequests, "rate limit exceeded")
			c.Abort()
			return
		}

		c.Next()
	}
}

// LimitByIP creates a rate limit middleware based on IP address
func (m *RateLimitMiddleware) LimitByIP(limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.ClientIP()

		allowed, err := m.limiter.IsAllowed(c.Request.Context(), key, limit, window)
		if err != nil || !allowed {
			response.Error(c, http.StatusTooManyRequests, "rate limit exceeded")
			c.Abort()
			return
		}

		c.Next()
	}
}
