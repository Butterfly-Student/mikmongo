package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"mikmongo/internal/model"
)

// MockSubscriptionRepository is a mock implementation of repository.SubscriptionRepository
type MockSubscriptionRepository struct {
	mock.Mock
}

func (m *MockSubscriptionRepository) Create(ctx context.Context, subscription *model.Subscription) error {
	args := m.Called(ctx, subscription)
	return args.Error(0)
}

func (m *MockSubscriptionRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Subscription, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Subscription), args.Error(1)
}

func (m *MockSubscriptionRepository) GetByCustomerID(ctx context.Context, customerID uuid.UUID) ([]model.Subscription, error) {
	args := m.Called(ctx, customerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Subscription), args.Error(1)
}

func (m *MockSubscriptionRepository) GetByUsername(ctx context.Context, username string) (*model.Subscription, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Subscription), args.Error(1)
}

func (m *MockSubscriptionRepository) Update(ctx context.Context, subscription *model.Subscription) error {
	args := m.Called(ctx, subscription)
	return args.Error(0)
}

func (m *MockSubscriptionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockSubscriptionRepository) List(ctx context.Context, limit, offset int) ([]model.Subscription, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Subscription), args.Error(1)
}

func (m *MockSubscriptionRepository) Count(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockSubscriptionRepository) ListByRouterID(ctx context.Context, routerID uuid.UUID, limit, offset int) ([]model.Subscription, error) {
	args := m.Called(ctx, routerID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Subscription), args.Error(1)
}

func (m *MockSubscriptionRepository) CountByRouterID(ctx context.Context, routerID uuid.UUID) (int64, error) {
	args := m.Called(ctx, routerID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockSubscriptionRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockSubscriptionRepository) ListByStatus(ctx context.Context, status string) ([]model.Subscription, error) {
	args := m.Called(ctx, status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Subscription), args.Error(1)
}
