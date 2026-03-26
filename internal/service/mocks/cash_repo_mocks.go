package mocks

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"mikmongo/internal/model"
	"mikmongo/internal/repository"
)

// MockCashEntryRepository mocks repository.CashEntryRepository
type MockCashEntryRepository struct {
	mock.Mock
}

func (m *MockCashEntryRepository) Create(ctx context.Context, entry *model.CashEntry) error {
	args := m.Called(ctx, entry)
	return args.Error(0)
}

func (m *MockCashEntryRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.CashEntry, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.CashEntry), args.Error(1)
}

func (m *MockCashEntryRepository) Update(ctx context.Context, entry *model.CashEntry) error {
	args := m.Called(ctx, entry)
	return args.Error(0)
}

func (m *MockCashEntryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCashEntryRepository) List(ctx context.Context, filter repository.CashEntryFilter, limit, offset int) ([]model.CashEntry, error) {
	args := m.Called(ctx, filter, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.CashEntry), args.Error(1)
}

func (m *MockCashEntryRepository) Count(ctx context.Context, filter repository.CashEntryFilter) (int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockCashEntryRepository) GetByReference(ctx context.Context, refType string, refID uuid.UUID) (*model.CashEntry, error) {
	args := m.Called(ctx, refType, refID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.CashEntry), args.Error(1)
}

func (m *MockCashEntryRepository) SumByTypeAndPeriod(ctx context.Context, entryType string, from, to time.Time) (float64, error) {
	args := m.Called(ctx, entryType, from, to)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockCashEntryRepository) SumBySourceAndPeriod(ctx context.Context, from, to time.Time) ([]repository.SourceSum, error) {
	args := m.Called(ctx, from, to)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]repository.SourceSum), args.Error(1)
}

// MockPettyCashFundRepository mocks repository.PettyCashFundRepository
type MockPettyCashFundRepository struct {
	mock.Mock
}

func (m *MockPettyCashFundRepository) Create(ctx context.Context, fund *model.PettyCashFund) error {
	args := m.Called(ctx, fund)
	return args.Error(0)
}

func (m *MockPettyCashFundRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.PettyCashFund, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.PettyCashFund), args.Error(1)
}

func (m *MockPettyCashFundRepository) Update(ctx context.Context, fund *model.PettyCashFund) error {
	args := m.Called(ctx, fund)
	return args.Error(0)
}

func (m *MockPettyCashFundRepository) List(ctx context.Context, limit, offset int) ([]model.PettyCashFund, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.PettyCashFund), args.Error(1)
}

func (m *MockPettyCashFundRepository) Count(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockPettyCashFundRepository) AdjustBalance(ctx context.Context, id uuid.UUID, delta float64) error {
	args := m.Called(ctx, id, delta)
	return args.Error(0)
}
