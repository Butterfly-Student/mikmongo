package service

import (
	"context"
	"time"
)

// CacheClient is a minimal interface for cache operations.
// *redis.Client from pkg/redis satisfies this automatically.
type CacheClient interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Del(ctx context.Context, keys ...string) error
}

// Cache TTLs
const (
	ttlBWProfile = 5 * time.Minute  // DB model cache for bandwidth profiles
	ttlSub       = 2 * time.Minute  // DB model cache for subscriptions
	ttlMtProfile = 30 * time.Second // Live MikroTik PPPProfile data
	ttlMtSecret  = 30 * time.Second // Live MikroTik PPPSecret data
)

// Cache key helpers
func keyBWProfile(id string) string              { return "bwprofile:" + id }
func keySubscription(id string) string           { return "sub:" + id }
func keyMtProfile(routerID, name string) string  { return "mt:profile:" + routerID + ":" + name }
func keyMtSecret(routerID, username string) string { return "mt:secret:" + routerID + ":" + username }
