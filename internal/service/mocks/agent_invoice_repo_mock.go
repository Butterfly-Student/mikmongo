package mocks

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"mikmongo/internal/model"
	"mikmongo/internal/repository"
)

// MockAgentInvoiceRepository mocks repository.AgentInvoiceRepository
type MockAgentInvoiceRepository struct {
	mock.Mock
}

func (m *MockAgentInvoiceRepository) Create(ctx context.Context, inv *model.AgentInvoice) error {
	args := m.Called(ctx, inv)
	return args.Error(0)
}

func (m *MockAgentInvoiceRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.AgentInvoice, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.AgentInvoice), args.Error(1)
}

func (m *MockAgentInvoiceRepository) GetByAgentAndPeriod(ctx context.Context, agentID uuid.UUID, periodStart time.Time, cycle string) (*model.AgentInvoice, error) {
	args := m.Called(ctx, agentID, periodStart, cycle)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.AgentInvoice), args.Error(1)
}

func (m *MockAgentInvoiceRepository) Update(ctx context.Context, inv *model.AgentInvoice) error {
	args := m.Called(ctx, inv)
	return args.Error(0)
}

func (m *MockAgentInvoiceRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string, paidAmount float64) error {
	args := m.Called(ctx, id, status, paidAmount)
	return args.Error(0)
}

func (m *MockAgentInvoiceRepository) UpdateStatusAndNotes(ctx context.Context, id uuid.UUID, status string, paidAmount float64, notes string) error {
	args := m.Called(ctx, id, status, paidAmount, notes)
	return args.Error(0)
}

func (m *MockAgentInvoiceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockAgentInvoiceRepository) List(ctx context.Context, filter repository.AgentInvoiceFilter, limit, offset int) ([]model.AgentInvoice, error) {
	args := m.Called(ctx, filter, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.AgentInvoice), args.Error(1)
}

func (m *MockAgentInvoiceRepository) Count(ctx context.Context, filter repository.AgentInvoiceFilter) (int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockAgentInvoiceRepository) ListByAgentID(ctx context.Context, agentID uuid.UUID, limit, offset int) ([]model.AgentInvoice, error) {
	args := m.Called(ctx, agentID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.AgentInvoice), args.Error(1)
}

func (m *MockAgentInvoiceRepository) GetUnpaidOverdue(ctx context.Context) ([]model.AgentInvoice, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.AgentInvoice), args.Error(1)
}
