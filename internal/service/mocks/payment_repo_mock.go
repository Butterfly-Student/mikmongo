package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"mikmongo/internal/model"
)

// MockPaymentRepository is a mock implementation of repository.PaymentRepository
type MockPaymentRepository struct {
	mock.Mock
}

func (m *MockPaymentRepository) Create(ctx context.Context, payment *model.Payment) error {
	args := m.Called(ctx, payment)
	return args.Error(0)
}

func (m *MockPaymentRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Payment, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Payment), args.Error(1)
}

func (m *MockPaymentRepository) GetByInvoiceID(ctx context.Context, invoiceID uuid.UUID) ([]model.Payment, error) {
	args := m.Called(ctx, invoiceID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Payment), args.Error(1)
}

func (m *MockPaymentRepository) GetByTransactionID(ctx context.Context, transactionID string) (*model.Payment, error) {
	args := m.Called(ctx, transactionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Payment), args.Error(1)
}

func (m *MockPaymentRepository) Update(ctx context.Context, payment *model.Payment) error {
	args := m.Called(ctx, payment)
	return args.Error(0)
}

func (m *MockPaymentRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockPaymentRepository) GetByCustomerID(ctx context.Context, customerID uuid.UUID) ([]model.Payment, error) {
	args := m.Called(ctx, customerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Payment), args.Error(1)
}

func (m *MockPaymentRepository) List(ctx context.Context, limit, offset int) ([]model.Payment, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Payment), args.Error(1)
}

// MockPaymentAllocationRepository is a mock implementation of repository.PaymentAllocationRepository
type MockPaymentAllocationRepository struct {
	mock.Mock
}

func (m *MockPaymentAllocationRepository) Create(ctx context.Context, allocation *model.PaymentAllocation) error {
	args := m.Called(ctx, allocation)
	return args.Error(0)
}

func (m *MockPaymentAllocationRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.PaymentAllocation, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.PaymentAllocation), args.Error(1)
}

func (m *MockPaymentAllocationRepository) Update(ctx context.Context, allocation *model.PaymentAllocation) error {
	args := m.Called(ctx, allocation)
	return args.Error(0)
}

func (m *MockPaymentAllocationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPaymentAllocationRepository) List(ctx context.Context, limit, offset int) ([]model.PaymentAllocation, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.PaymentAllocation), args.Error(1)
}

func (m *MockPaymentAllocationRepository) Count(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockPaymentAllocationRepository) ListByPaymentID(ctx context.Context, paymentID uuid.UUID) ([]model.PaymentAllocation, error) {
	args := m.Called(ctx, paymentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.PaymentAllocation), args.Error(1)
}

func (m *MockPaymentAllocationRepository) ListByInvoiceID(ctx context.Context, invoiceID uuid.UUID) ([]model.PaymentAllocation, error) {
	args := m.Called(ctx, invoiceID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.PaymentAllocation), args.Error(1)
}
