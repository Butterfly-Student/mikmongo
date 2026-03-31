package redis

import (
	"context"
	"time"
)

// HSet sets fields in a Redis hash.
func (c *Client) HSet(ctx context.Context, key string, fields map[string]any) error {
	return c.client.HSet(ctx, key, fields).Err()
}

// HGetAll returns all fields and values of a Redis hash.
func (c *Client) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return c.client.HGetAll(ctx, key).Result()
}

// HGet returns the value of a single field in a Redis hash.
func (c *Client) HGet(ctx context.Context, key string, field string) (string, error) {
	return c.client.HGet(ctx, key, field).Result()
}

// HDel deletes one or more fields from a Redis hash.
func (c *Client) HDel(ctx context.Context, key string, fields ...string) error {
	return c.client.HDel(ctx, key, fields...).Err()
}

// HSetWithTTL sets fields in a hash and applies a TTL to the key.
func (c *Client) HSetWithTTL(ctx context.Context, key string, fields map[string]any, ttl time.Duration) error {
	pipe := c.client.Pipeline()
	pipe.HSet(ctx, key, fields)
	pipe.Expire(ctx, key, ttl)
	_, err := pipe.Exec(ctx)
	return err
}
