package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

// SessionData represents user session data
type SessionData struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

// Key prefixes
const (
	SessionPrefix         = "session:"
	BlacklistPrefix       = "blacklist:"
	PasswordChangedPrefix = "pwd_changed:"
	SelectedRouterPrefix  = "selected_router:"
)

// SetSession stores JWT session
func (c *Client) SetSession(ctx context.Context, token string, data *SessionData, ttl time.Duration) error {
	key := fmt.Sprintf("%s%s", SessionPrefix, token)
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return c.Set(ctx, key, jsonData, ttl)
}

// GetSession retrieves session data
func (c *Client) GetSession(ctx context.Context, token string) (*SessionData, error) {
	key := fmt.Sprintf("%s%s", SessionPrefix, token)
	data, err := c.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	var session SessionData
	if err := json.Unmarshal([]byte(data), &session); err != nil {
		return nil, err
	}
	return &session, nil
}

// DeleteSession removes session
func (c *Client) DeleteSession(ctx context.Context, token string) error {
	key := fmt.Sprintf("%s%s", SessionPrefix, token)
	return c.Del(ctx, key)
}

// BlacklistToken adds a JTI to the blacklist with the given TTL
func (c *Client) BlacklistToken(ctx context.Context, jti string, ttl time.Duration) error {
	key := fmt.Sprintf("%s%s", BlacklistPrefix, jti)
	return c.Set(ctx, key, "1", ttl)
}

// IsBlacklisted checks if a JTI is blacklisted
func (c *Client) IsBlacklisted(ctx context.Context, jti string) (bool, error) {
	key := fmt.Sprintf("%s%s", BlacklistPrefix, jti)
	return c.Exists(ctx, key)
}

// SetPasswordChangedAt stores the timestamp when a user changed their password
func (c *Client) SetPasswordChangedAt(ctx context.Context, userID string, t time.Time, ttl time.Duration) error {
	key := fmt.Sprintf("%s%s", PasswordChangedPrefix, userID)
	return c.Set(ctx, key, t.Unix(), ttl)
}

// GetPasswordChangedAt retrieves the timestamp when a user last changed their password
func (c *Client) GetPasswordChangedAt(ctx context.Context, userID string) (time.Time, error) {
	key := fmt.Sprintf("%s%s", PasswordChangedPrefix, userID)
	data, err := c.Get(ctx, key)
	if err != nil {
		if err == goredis.Nil {
			return time.Time{}, nil
		}
		return time.Time{}, err
	}

	var unix int64
	if _, err := fmt.Sscanf(data, "%d", &unix); err != nil {
		return time.Time{}, err
	}
	return time.Unix(unix, 0), nil
}

// SetSelectedRouter stores the selected router ID for a user
func (c *Client) SetSelectedRouter(ctx context.Context, userID, routerID string, ttl time.Duration) error {
	key := fmt.Sprintf("%s%s", SelectedRouterPrefix, userID)
	return c.Set(ctx, key, routerID, ttl)
}

// GetSelectedRouter retrieves the selected router ID for a user
func (c *Client) GetSelectedRouter(ctx context.Context, userID string) (string, error) {
	key := fmt.Sprintf("%s%s", SelectedRouterPrefix, userID)
	data, err := c.Get(ctx, key)
	if err != nil {
		if err == goredis.Nil {
			return "", nil
		}
		return "", err
	}
	return data, nil
}

// ClearSelectedRouter removes the selected router for a user
func (c *Client) ClearSelectedRouter(ctx context.Context, userID string) error {
	key := fmt.Sprintf("%s%s", SelectedRouterPrefix, userID)
	return c.Del(ctx, key)
}
