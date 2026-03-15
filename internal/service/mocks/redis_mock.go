package mocks

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"
)

// MockRedisClient is a mock for service.RedisClientInterface
type MockRedisClient struct {
	mock.Mock
}

func (m *MockRedisClient) BlacklistToken(ctx context.Context, jti string, ttl time.Duration) error {
	args := m.Called(ctx, jti, ttl)
	return args.Error(0)
}

func (m *MockRedisClient) IsBlacklisted(ctx context.Context, jti string) (bool, error) {
	args := m.Called(ctx, jti)
	return args.Bool(0), args.Error(1)
}

func (m *MockRedisClient) SetPasswordChangedAt(ctx context.Context, userID string, t time.Time, ttl time.Duration) error {
	args := m.Called(ctx, userID, t, ttl)
	return args.Error(0)
}
