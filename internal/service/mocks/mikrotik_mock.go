package mocks

import (
	"context"

	mkdomain "github.com/Butterfly-Student/go-ros/domain"
	"github.com/stretchr/testify/mock"
)

// MockMikrotikClientAdapter is a mock implementation of service.MikrotikClientAdapter.
// MockMikrotikProvider is defined locally in subscription_service_test.go to avoid
// the import cycle (service/mocks → service → service/mocks).
type MockMikrotikClientAdapter struct {
	mock.Mock
}

func (m *MockMikrotikClientAdapter) AddSecret(ctx context.Context, secret *mkdomain.PPPSecret) error {
	args := m.Called(ctx, secret)
	return args.Error(0)
}

func (m *MockMikrotikClientAdapter) UpdateSecret(ctx context.Context, id string, secret *mkdomain.PPPSecret) error {
	args := m.Called(ctx, id, secret)
	return args.Error(0)
}

func (m *MockMikrotikClientAdapter) RemoveSecret(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockMikrotikClientAdapter) GetSecretByName(ctx context.Context, name string) (*mkdomain.PPPSecret, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*mkdomain.PPPSecret), args.Error(1)
}

func (m *MockMikrotikClientAdapter) DisableSecret(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockMikrotikClientAdapter) EnableSecret(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
