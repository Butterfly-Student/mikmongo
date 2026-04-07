// writer/redis_handler.go - Redis Batch Handler
package writer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisHandler writes batches to Redis
type RedisHandler struct {
	client *redis.Client
	prefix string
}

// RedisConfig untuk Redis connection
type RedisConfig struct {
	Addr   string
	DB     int
	Prefix string // Key prefix: "mikrotik"
}

// NewRedisHandler creates new Redis handler
func NewRedisHandler(config RedisConfig) (*RedisHandler, error) {
	client := redis.NewClient(&redis.Options{
		Addr: config.Addr,
		DB:   config.DB,
	})
	
	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping failed: %w", err)
	}
	
	return &RedisHandler{
		client: client,
		prefix: config.Prefix,
	}, nil
}

// Close closes Redis client
func (h *RedisHandler) Close() error {
	return h.client.Close()
}

// WriteBatch writes batch of items to Redis
func (h *RedisHandler) WriteBatch(ctx context.Context, items []WriteItem) error {
	pipe := h.client.Pipeline()
	
	for _, item := range items {
		// Build full key: prefix:router:key
		fullKey := fmt.Sprintf("%s:%s:%s", h.prefix, item.RouterID, item.Key)
		
		// Convert value to JSON
		jsonValue, err := json.Marshal(item.Value)
		if err != nil {
			continue // Skip invalid items
		}
		
		// HSET
		pipe.HSet(ctx, fullKey, item.Field, jsonValue)
		
		// Set TTL jika ada
		if item.TTL > 0 {
			pipe.Expire(ctx, fullKey, item.TTL)
		}
	}
	
	_, err := pipe.Exec(ctx)
	return err
}

// Get retrieves single field dari Redis
func (h *RedisHandler) Get(ctx context.Context, routerID, key, field string) (map[string]string, error) {
	fullKey := fmt.Sprintf("%s:%s:%s", h.prefix, routerID, key)
	
	data, err := h.client.HGet(ctx, fullKey, field).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	
	var result map[string]string
	if err := json.Unmarshal([]byte(data), &result); err != nil {
		return nil, err
	}
	
	return result, nil
}

// GetAll retrieves all fields dari Redis key
func (h *RedisHandler) GetAll(ctx context.Context, routerID, key string) (map[string]map[string]string, error) {
	fullKey := fmt.Sprintf("%s:%s:%s", h.prefix, routerID, key)
	
	allData, err := h.client.HGetAll(ctx, fullKey).Result()
	if err != nil {
		return nil, err
	}
	
	result := make(map[string]map[string]string)
	for field, data := range allData {
		var item map[string]string
		if err := json.Unmarshal([]byte(data), &item); err != nil {
			continue
		}
		result[field] = item
	}
	
	return result, nil
}

// Delete menghapus key dari Redis
func (h *RedisHandler) Delete(ctx context.Context, routerID, key string) error {
	fullKey := fmt.Sprintf("%s:%s:%s", h.prefix, routerID, key)
	return h.client.Del(ctx, fullKey).Err()
}

// DeleteField menghapus field dari hash
func (h *RedisHandler) DeleteField(ctx context.Context, routerID, key, field string) error {
	fullKey := fmt.Sprintf("%s:%s:%s", h.prefix, routerID, key)
	return h.client.HDel(ctx, fullKey, field).Err()
}

// Invalidate menghapus cache setelah write operation (Pipeline C)
func (h *RedisHandler) Invalidate(ctx context.Context, routerID string, keys []string) error {
	pipe := h.client.Pipeline()
	
	for _, key := range keys {
		fullKey := fmt.Sprintf("%s:%s:%s", h.prefix, routerID, key)
		pipe.Del(ctx, fullKey)
	}
	
	_, err := pipe.Exec(ctx)
	return err
}
