package mocks

import (
	"context"
	"time"

	mikhmonDomain "github.com/Butterfly-Student/go-ros/domain/mikhmon"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"mikmongo/internal/model"
	"mikmongo/internal/repository"
)

// MockVoucherGenerator mocks service.VoucherGenerator
type MockVoucherGenerator struct {
	mock.Mock
}

func (m *MockVoucherGenerator) GenerateBatch(ctx context.Context, routerID uuid.UUID, req *mikhmonDomain.VoucherGenerateRequest) (*mikhmonDomain.VoucherBatch, error) {
	args := m.Called(ctx, routerID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*mikhmonDomain.VoucherBatch), args.Error(1)
}

// MockHotspotSaleRepository mocks repository.HotspotSaleRepository
type MockHotspotSaleRepository struct {
	mock.Mock
}

func (m *MockHotspotSaleRepository) Create(ctx context.Context, sale *model.HotspotSale) error {
	args := m.Called(ctx, sale)
	return args.Error(0)
}

func (m *MockHotspotSaleRepository) CreateBatch(ctx context.Context, sales []model.HotspotSale) error {
	args := m.Called(ctx, sales)
	return args.Error(0)
}

func (m *MockHotspotSaleRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.HotspotSale, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.HotspotSale), args.Error(1)
}

func (m *MockHotspotSaleRepository) List(ctx context.Context, filter repository.HotspotSaleFilter, limit, offset int) ([]model.HotspotSale, error) {
	args := m.Called(ctx, filter, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.HotspotSale), args.Error(1)
}

func (m *MockHotspotSaleRepository) Count(ctx context.Context, filter repository.HotspotSaleFilter) (int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockHotspotSaleRepository) ListByBatchCode(ctx context.Context, routerID uuid.UUID, batchCode string) ([]model.HotspotSale, error) {
	args := m.Called(ctx, routerID, batchCode)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.HotspotSale), args.Error(1)
}

func (m *MockHotspotSaleRepository) DeleteByBatchCode(ctx context.Context, routerID uuid.UUID, batchCode string) error {
	args := m.Called(ctx, routerID, batchCode)
	return args.Error(0)
}

func (m *MockHotspotSaleRepository) SumByAgentAndPeriod(ctx context.Context, agentID uuid.UUID, from, to time.Time) (count int, subtotal, sellingTotal float64, err error) {
	args := m.Called(ctx, agentID, from, to)
	return args.Int(0), args.Get(1).(float64), args.Get(2).(float64), args.Error(3)
}

// MockSalesAgentRepository mocks repository.SalesAgentRepository
type MockSalesAgentRepository struct {
	mock.Mock
}

func (m *MockSalesAgentRepository) Create(ctx context.Context, agent *model.SalesAgent) error {
	args := m.Called(ctx, agent)
	return args.Error(0)
}

func (m *MockSalesAgentRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.SalesAgent, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.SalesAgent), args.Error(1)
}

func (m *MockSalesAgentRepository) GetByUsername(ctx context.Context, username string) (*model.SalesAgent, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.SalesAgent), args.Error(1)
}

func (m *MockSalesAgentRepository) Update(ctx context.Context, agent *model.SalesAgent) error {
	args := m.Called(ctx, agent)
	return args.Error(0)
}

func (m *MockSalesAgentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockSalesAgentRepository) List(ctx context.Context, routerID *uuid.UUID, limit, offset int) ([]model.SalesAgent, error) {
	args := m.Called(ctx, routerID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.SalesAgent), args.Error(1)
}

func (m *MockSalesAgentRepository) Count(ctx context.Context, routerID *uuid.UUID) (int64, error) {
	args := m.Called(ctx, routerID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockSalesAgentRepository) UpsertProfilePrice(ctx context.Context, price *model.SalesProfilePrice) error {
	args := m.Called(ctx, price)
	return args.Error(0)
}

func (m *MockSalesAgentRepository) GetProfilePrice(ctx context.Context, agentID uuid.UUID, profileName string) (*model.SalesProfilePrice, error) {
	args := m.Called(ctx, agentID, profileName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.SalesProfilePrice), args.Error(1)
}

func (m *MockSalesAgentRepository) ListProfilePrices(ctx context.Context, agentID uuid.UUID) ([]model.SalesProfilePrice, error) {
	args := m.Called(ctx, agentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.SalesProfilePrice), args.Error(1)
}
