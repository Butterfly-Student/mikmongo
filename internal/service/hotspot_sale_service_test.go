package service

import (
	"context"
	"errors"
	"testing"

	mikhmonDomain "github.com/Butterfly-Student/go-ros/domain/mikhmon"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"mikmongo/internal/model"
	"mikmongo/internal/repository"
	"mikmongo/internal/service/mocks"
)

// newHotspotSaleServiceWithMocks creates a HotspotSaleService backed by mocks.
func newHotspotSaleServiceWithMocks() (
	*HotspotSaleService,
	*mocks.MockVoucherGenerator,
	*mocks.MockHotspotSaleRepository,
	*mocks.MockSalesAgentRepository,
) {
	voucherGen := &mocks.MockVoucherGenerator{}
	saleRepo := &mocks.MockHotspotSaleRepository{}
	agentRepo := &mocks.MockSalesAgentRepository{}
	svc := NewHotspotSaleService(voucherGen, saleRepo, agentRepo)
	return svc, voucherGen, saleRepo, agentRepo
}

// fakeBatch returns a minimal VoucherBatch with n vouchers for testing.
func fakeBatch(profile, code string, n int) *mikhmonDomain.VoucherBatch {
	vouchers := make([]mikhmonDomain.Voucher, n)
	for i := range vouchers {
		vouchers[i] = mikhmonDomain.Voucher{Name: "voucher" + string(rune('A'+i)), Profile: profile}
	}
	return &mikhmonDomain.VoucherBatch{
		Code:     code,
		Quantity: n,
		Profile:  profile,
		Vouchers: vouchers,
	}
}

// --- GenerateBatchAndRecord ---

func TestGenerateBatch_NoAgent(t *testing.T) {
	ctx := context.Background()
	svc, voucherGen, saleRepo, _ := newHotspotSaleServiceWithMocks()

	routerID := uuid.New()
	req := &mikhmonDomain.VoucherGenerateRequest{Profile: "10mb", Prefix: "ts"}
	batch := fakeBatch("10mb", "ABCD", 3)

	voucherGen.On("GenerateBatch", ctx, routerID, req).Return(batch, nil)
	saleRepo.On("CreateBatch", ctx, mock.AnythingOfType("[]model.HotspotSale")).Return(nil)

	result, err := svc.GenerateBatchAndRecord(ctx, routerID, req, nil)
	require.NoError(t, err)
	assert.Equal(t, batch, result)

	// Verify CreateBatch was called once
	saleRepo.AssertCalled(t, "CreateBatch", ctx, mock.AnythingOfType("[]model.HotspotSale"))
}

func TestGenerateBatch_NoAgent_SalesRecordFields(t *testing.T) {
	ctx := context.Background()
	svc, voucherGen, saleRepo, _ := newHotspotSaleServiceWithMocks()

	routerID := uuid.New()
	req := &mikhmonDomain.VoucherGenerateRequest{Profile: "5mb", Prefix: "pf"}
	batch := fakeBatch("5mb", "XY12", 2)

	voucherGen.On("GenerateBatch", ctx, routerID, req).Return(batch, nil)

	var capturedSales []model.HotspotSale
	saleRepo.On("CreateBatch", ctx, mock.AnythingOfType("[]model.HotspotSale")).
		Run(func(args mock.Arguments) {
			capturedSales = args.Get(1).([]model.HotspotSale)
		}).
		Return(nil)

	_, err := svc.GenerateBatchAndRecord(ctx, routerID, req, nil)
	require.NoError(t, err)

	require.Len(t, capturedSales, 2)
	for _, s := range capturedSales {
		assert.Equal(t, routerID.String(), s.RouterID)
		assert.Equal(t, "5mb", s.Profile)
		assert.Equal(t, "XY12", s.BatchCode)
		assert.Equal(t, "pf", s.Prefix)
		assert.Nil(t, s.SalesAgentID)
		assert.Equal(t, 0.0, s.Price)
		assert.Equal(t, 0.0, s.SellingPrice)
	}
}

func TestGenerateBatch_WithAgent_WithProfilePrice(t *testing.T) {
	ctx := context.Background()
	svc, voucherGen, saleRepo, agentRepo := newHotspotSaleServiceWithMocks()

	routerID := uuid.New()
	agentID := uuid.New()
	agentIDStr := agentID.String()
	req := &mikhmonDomain.VoucherGenerateRequest{Profile: "20mb"}
	batch := fakeBatch("20mb", "BCDE", 1)

	agent := &model.SalesAgent{ID: agentIDStr}
	pp := &model.SalesProfilePrice{BasePrice: 5000, SellingPrice: 7000}

	voucherGen.On("GenerateBatch", ctx, routerID, req).Return(batch, nil)
	agentRepo.On("GetByID", ctx, agentID).Return(agent, nil)
	agentRepo.On("GetProfilePrice", ctx, agentID, "20mb").Return(pp, nil)

	var capturedSales []model.HotspotSale
	saleRepo.On("CreateBatch", ctx, mock.AnythingOfType("[]model.HotspotSale")).
		Run(func(args mock.Arguments) {
			capturedSales = args.Get(1).([]model.HotspotSale)
		}).
		Return(nil)

	_, err := svc.GenerateBatchAndRecord(ctx, routerID, req, &agentID)
	require.NoError(t, err)

	require.Len(t, capturedSales, 1)
	assert.Equal(t, 5000.0, capturedSales[0].Price)
	assert.Equal(t, 7000.0, capturedSales[0].SellingPrice)
	require.NotNil(t, capturedSales[0].SalesAgentID)
	assert.Equal(t, agentIDStr, *capturedSales[0].SalesAgentID)
}

func TestGenerateBatch_WithAgent_NoProfilePrice(t *testing.T) {
	ctx := context.Background()
	svc, voucherGen, saleRepo, agentRepo := newHotspotSaleServiceWithMocks()

	routerID := uuid.New()
	agentID := uuid.New()
	agentIDStr := agentID.String()
	req := &mikhmonDomain.VoucherGenerateRequest{Profile: "10mb"}
	batch := fakeBatch("10mb", "AA11", 1)

	agent := &model.SalesAgent{ID: agentIDStr}

	voucherGen.On("GenerateBatch", ctx, routerID, req).Return(batch, nil)
	agentRepo.On("GetByID", ctx, agentID).Return(agent, nil)
	agentRepo.On("GetProfilePrice", ctx, agentID, "10mb").Return(nil, errors.New("not found"))

	var capturedSales []model.HotspotSale
	saleRepo.On("CreateBatch", ctx, mock.AnythingOfType("[]model.HotspotSale")).
		Run(func(args mock.Arguments) {
			capturedSales = args.Get(1).([]model.HotspotSale)
		}).
		Return(nil)

	_, err := svc.GenerateBatchAndRecord(ctx, routerID, req, &agentID)
	require.NoError(t, err)

	require.Len(t, capturedSales, 1)
	assert.Equal(t, 0.0, capturedSales[0].Price)
	assert.Equal(t, 0.0, capturedSales[0].SellingPrice)
	require.NotNil(t, capturedSales[0].SalesAgentID)
}

func TestGenerateBatch_MikrotikFails(t *testing.T) {
	ctx := context.Background()
	svc, voucherGen, saleRepo, _ := newHotspotSaleServiceWithMocks()

	routerID := uuid.New()
	req := &mikhmonDomain.VoucherGenerateRequest{Profile: "10mb"}

	voucherGen.On("GenerateBatch", ctx, routerID, req).Return(nil, errors.New("mikrotik error"))

	result, err := svc.GenerateBatchAndRecord(ctx, routerID, req, nil)
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "mikrotik")

	// DB must NOT have been touched
	saleRepo.AssertNotCalled(t, "CreateBatch", mock.Anything, mock.Anything)
}

func TestGenerateBatch_AgentNotFound(t *testing.T) {
	ctx := context.Background()
	svc, voucherGen, saleRepo, agentRepo := newHotspotSaleServiceWithMocks()

	routerID := uuid.New()
	agentID := uuid.New()
	req := &mikhmonDomain.VoucherGenerateRequest{Profile: "10mb"}
	batch := fakeBatch("10mb", "ZZ99", 1)

	voucherGen.On("GenerateBatch", ctx, routerID, req).Return(batch, nil)
	agentRepo.On("GetByID", ctx, agentID).Return(nil, errors.New("not found"))

	result, err := svc.GenerateBatchAndRecord(ctx, routerID, req, &agentID)
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "sales agent not found")

	saleRepo.AssertNotCalled(t, "CreateBatch", mock.Anything, mock.Anything)
}

func TestGenerateBatch_DBFails_ReturnsPartialError(t *testing.T) {
	ctx := context.Background()
	svc, voucherGen, saleRepo, _ := newHotspotSaleServiceWithMocks()

	routerID := uuid.New()
	req := &mikhmonDomain.VoucherGenerateRequest{Profile: "10mb"}
	batch := fakeBatch("10mb", "PP22", 2)

	voucherGen.On("GenerateBatch", ctx, routerID, req).Return(batch, nil)
	saleRepo.On("CreateBatch", ctx, mock.AnythingOfType("[]model.HotspotSale")).Return(errors.New("db error"))

	result, err := svc.GenerateBatchAndRecord(ctx, routerID, req, nil)
	// MikroTik-first: batch is returned even on DB failure
	require.Error(t, err)
	assert.NotNil(t, result, "batch must be returned even when DB write fails")
	assert.Contains(t, err.Error(), "failed to record sales")
}

// --- ListSales ---

func TestListSales_Empty(t *testing.T) {
	ctx := context.Background()
	svc, _, saleRepo, _ := newHotspotSaleServiceWithMocks()

	filter := repository.HotspotSaleFilter{}
	saleRepo.On("List", ctx, filter, 10, 0).Return([]model.HotspotSale{}, nil)
	saleRepo.On("Count", ctx, filter).Return(int64(0), nil)

	sales, count, err := svc.ListSales(ctx, filter, 10, 0)
	require.NoError(t, err)
	assert.Empty(t, sales)
	assert.Equal(t, int64(0), count)
}

func TestListSales_WithFilter(t *testing.T) {
	ctx := context.Background()
	svc, _, saleRepo, _ := newHotspotSaleServiceWithMocks()

	routerID := uuid.New()
	filter := repository.HotspotSaleFilter{RouterID: &routerID, Profile: "10mb"}
	expected := []model.HotspotSale{{RouterID: routerID.String(), Profile: "10mb"}}

	saleRepo.On("List", ctx, filter, 20, 0).Return(expected, nil)
	saleRepo.On("Count", ctx, filter).Return(int64(1), nil)

	sales, count, err := svc.ListSales(ctx, filter, 20, 0)
	require.NoError(t, err)
	assert.Len(t, sales, 1)
	assert.Equal(t, int64(1), count)
}

func TestListSales_Pagination(t *testing.T) {
	ctx := context.Background()
	svc, _, saleRepo, _ := newHotspotSaleServiceWithMocks()

	filter := repository.HotspotSaleFilter{}
	saleRepo.On("List", ctx, filter, 5, 10).Return([]model.HotspotSale{{}}, nil)
	saleRepo.On("Count", ctx, filter).Return(int64(15), nil)

	sales, count, err := svc.ListSales(ctx, filter, 5, 10)
	require.NoError(t, err)
	assert.Len(t, sales, 1)
	assert.Equal(t, int64(15), count)

	saleRepo.AssertCalled(t, "List", ctx, filter, 5, 10)
}

func TestListSales_CountError(t *testing.T) {
	ctx := context.Background()
	svc, _, saleRepo, _ := newHotspotSaleServiceWithMocks()

	filter := repository.HotspotSaleFilter{}
	saleRepo.On("List", ctx, filter, 10, 0).Return([]model.HotspotSale{}, nil)
	saleRepo.On("Count", ctx, filter).Return(int64(0), errors.New("count failed"))

	_, _, err := svc.ListSales(ctx, filter, 10, 0)
	require.Error(t, err)
}
