package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"mikmongo/internal/model"
)

// MockInvoiceRepository is a mock implementation of repository.InvoiceRepository
type MockInvoiceRepository struct {
	mock.Mock
}

func (m *MockInvoiceRepository) Create(ctx context.Context, invoice *model.Invoice) error {
	args := m.Called(ctx, invoice)
	return args.Error(0)
}

func (m *MockInvoiceRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Invoice, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Invoice), args.Error(1)
}

func (m *MockInvoiceRepository) GetByCustomerID(ctx context.Context, customerID uuid.UUID) ([]model.Invoice, error) {
	args := m.Called(ctx, customerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Invoice), args.Error(1)
}

func (m *MockInvoiceRepository) GetByCustomerIDForUpdate(ctx context.Context, customerID uuid.UUID) ([]model.Invoice, error) {
	args := m.Called(ctx, customerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Invoice), args.Error(1)
}

func (m *MockInvoiceRepository) Update(ctx context.Context, invoice *model.Invoice) error {
	args := m.Called(ctx, invoice)
	return args.Error(0)
}

func (m *MockInvoiceRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockInvoiceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockInvoiceRepository) List(ctx context.Context, limit, offset int) ([]model.Invoice, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Invoice), args.Error(1)
}

func (m *MockInvoiceRepository) GetOverdue(ctx context.Context) ([]model.Invoice, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Invoice), args.Error(1)
}

func (m *MockInvoiceRepository) GetBySubscriptionAndPeriod(ctx context.Context, subID uuid.UUID, month, year int) (*model.Invoice, error) {
	args := m.Called(ctx, subID, month, year)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Invoice), args.Error(1)
}

// MockInvoiceItemRepository is a mock implementation of repository.InvoiceItemRepository
type MockInvoiceItemRepository struct {
	mock.Mock
}

func (m *MockInvoiceItemRepository) Create(ctx context.Context, item *model.InvoiceItem) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MockInvoiceItemRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.InvoiceItem, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.InvoiceItem), args.Error(1)
}

func (m *MockInvoiceItemRepository) Update(ctx context.Context, item *model.InvoiceItem) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MockInvoiceItemRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockInvoiceItemRepository) List(ctx context.Context, limit, offset int) ([]model.InvoiceItem, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.InvoiceItem), args.Error(1)
}

func (m *MockInvoiceItemRepository) Count(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockInvoiceItemRepository) ListByInvoiceID(ctx context.Context, invoiceID uuid.UUID) ([]model.InvoiceItem, error) {
	args := m.Called(ctx, invoiceID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.InvoiceItem), args.Error(1)
}

func (m *MockInvoiceItemRepository) DeleteByInvoiceID(ctx context.Context, invoiceID uuid.UUID) error {
	args := m.Called(ctx, invoiceID)
	return args.Error(0)
}
