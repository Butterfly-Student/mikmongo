package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RateLimiter implements sliding window rate limiting
type RateLimiter struct {
	client *Client
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(client *Client) *RateLimiter {
	return &RateLimiter{client: client}
}

// IsAllowed checks if request is within rate limit
func (r *RateLimiter) IsAllowed(ctx context.Context, key string, limit int, window time.Duration) (bool, error) {
	now := time.Now().Unix()
	windowStart := now - int64(window.Seconds())
	
	redisKey := fmt.Sprintf("ratelimit:%s", key)
	
	// Remove old entries
	pipe := r.client.client.Pipeline()
	pipe.ZRemRangeByScore(ctx, redisKey, "0", fmt.Sprintf("%d", windowStart))
	
	// Count current entries
	countCmd := pipe.ZCard(ctx, redisKey)
	
	// Add current request
	pipe.ZAdd(ctx, redisKey, redis.Z{
		Score:  float64(now),
		Member: now,
	})
	
	// Set expiry on the key
	pipe.Expire(ctx, redisKey, window)
	
	_, err := pipe.Exec(ctx)
	if err != nil {
		return false, err
	}
	
	count := countCmd.Val()
	return count < int64(limit), nil
}

// Reset resets rate limit for key
func (r *RateLimiter) Reset(ctx context.Context, key string) error {
	redisKey := fmt.Sprintf("ratelimit:%s", key)
	return r.client.Del(ctx, redisKey)
}
