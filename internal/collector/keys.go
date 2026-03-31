package collector

import (
	"time"

	"github.com/google/uuid"
)

const (
	// Prefix for all collector Redis keys and channels.
	Prefix = "monitor:"

	// TTLRealtime is the cache expiry for real-time data.
	// Short so stale data disappears quickly if the collector stops.
	TTLRealtime = 45 * time.Second

	// TTLSlowChanging is the cache expiry for infrequently changing data.
	TTLSlowChanging = 5 * time.Minute

	// PollInterval is the default interval for polling slow-changing data.
	PollInterval = 5 * time.Minute
)

// CacheKey returns the Redis key for a router's monitoring topic.
// Format: "monitor:{routerID}:{topic}"
func CacheKey(routerID uuid.UUID, topic string) string {
	return Prefix + routerID.String() + ":" + topic
}

// PubSubChannel returns the Redis Pub/Sub channel for a topic.
// Uses the same format as CacheKey for simplicity.
func PubSubChannel(routerID uuid.UUID, topic string) string {
	return Prefix + routerID.String() + ":" + topic
}

// TTLFor returns the appropriate TTL for a data category.
func TTLFor(category DataCategory) time.Duration {
	if category == RealTime {
		return TTLRealtime
	}
	return TTLSlowChanging
}
