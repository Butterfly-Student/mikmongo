package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"mikmongo/internal/model"
	"mikmongo/internal/repository"
	"mikmongo/internal/service/mocks"
)

func newAgentInvoiceSvcWithMocks() (
	*AgentInvoiceService,
	*mocks.MockAgentInvoiceRepository,
	*mocks.MockHotspotSaleRepository,
	*mocks.MockSalesAgentRepository,
	*mocks.MockSequenceCounterRepository,
) {
	invoiceRepo := &mocks.MockAgentInvoiceRepository{}
	saleRepo := &mocks.MockHotspotSaleRepository{}
	agentRepo := &mocks.MockSalesAgentRepository{}
	seqRepo := &mocks.MockSequenceCounterRepository{}
	svc := NewAgentInvoiceService(invoiceRepo, saleRepo, agentRepo, seqRepo)
	return svc, invoiceRepo, saleRepo, agentRepo, seqRepo
}

// --- GenerateForAgent tests ---

func TestGenerateForAgent_AgentNotFound(t *testing.T) {
	ctx := context.Background()
	svc, _, _, agentRepo, _ := newAgentInvoiceSvcWithMocks()
	agentID := uuid.New()

	agentRepo.On("GetByID", ctx, agentID).Return(nil, errors.New("not found"))

	_, err := svc.GenerateForAgent(ctx, agentID, time.Now(), time.Now(), "monthly")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "agent not found")
}

func TestGenerateForAgent_AgentInactive(t *testing.T) {
	ctx := context.Background()
	svc, _, _, agentRepo, _ := newAgentInvoiceSvcWithMocks()
	agentID := uuid.New()

	agent := &model.SalesAgent{ID: agentID.String(), Status: "inactive"}
	agentRepo.On("GetByID", ctx, agentID).Return(agent, nil)

	_, err := svc.GenerateForAgent(ctx, agentID, time.Now(), time.Now(), "monthly")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not active")
}

func TestGenerateForAgent_IdempotentExisting(t *testing.T) {
	ctx := context.Background()
	svc, invoiceRepo, _, agentRepo, _ := newAgentInvoiceSvcWithMocks()
	agentID := uuid.New()
	periodStart := time.Date(2024, time.March, 1, 0, 0, 0, 0, time.UTC)
	periodEnd := time.Date(2024, time.April, 1, 0, 0, 0, 0, time.UTC)

	agent := &model.SalesAgent{ID: agentID.String(), Status: "active"}
	agentRepo.On("GetByID", ctx, agentID).Return(agent, nil)

	existing := &model.AgentInvoice{
		ID:            uuid.New().String(),
		AgentID:       agentID.String(),
		InvoiceNumber: "AGT000001",
		Status:        "unpaid",
	}
	invoiceRepo.On("GetByAgentAndPeriod", ctx, agentID, periodStart, "monthly").Return(existing, nil)

	inv, err := svc.GenerateForAgent(ctx, agentID, periodStart, periodEnd, "monthly")
	require.NoError(t, err)
	assert.Equal(t, existing.InvoiceNumber, inv.InvoiceNumber)
}

func TestGenerateForAgent_Success_Monthly(t *testing.T) {
	ctx := context.Background()
	svc, invoiceRepo, saleRepo, agentRepo, seqRepo := newAgentInvoiceSvcWithMocks()
	agentID := uuid.New()
	routerID := uuid.New().String()
	periodStart := time.Date(2024, time.March, 1, 0, 0, 0, 0, time.UTC)
	periodEnd := time.Date(2024, time.April, 1, 0, 0, 0, 0, time.UTC)

	agent := &model.SalesAgent{ID: agentID.String(), Status: "active", RouterID: routerID}
	agentRepo.On("GetByID", ctx, agentID).Return(agent, nil)
	invoiceRepo.On("GetByAgentAndPeriod", ctx, agentID, periodStart, "monthly").Return(nil, errors.New("not found"))
	saleRepo.On("SumByAgentAndPeriod", ctx, agentID, periodStart, periodEnd).Return(10, 500000.0, 750000.0, nil)
	seqRepo.On("NextNumber", ctx, "agent_invoice").Return(42, nil)
	invoiceRepo.On("Create", ctx, mock.AnythingOfType("*model.AgentInvoice")).Return(nil)

	inv, err := svc.GenerateForAgent(ctx, agentID, periodStart, periodEnd, "monthly")
	require.NoError(t, err)
	assert.Equal(t, "AGT000042", inv.InvoiceNumber)
	assert.Equal(t, "monthly", inv.BillingCycle)
	assert.Equal(t, 10, inv.VoucherCount)
	assert.Equal(t, 500000.0, inv.Subtotal)
	assert.Equal(t, 750000.0, inv.SellingTotal)
	assert.Equal(t, 750000.0, inv.TotalAmount)
	assert.Equal(t, "unpaid", inv.Status)
	assert.NotNil(t, inv.BillingMonth)
	assert.Equal(t, 3, *inv.BillingMonth)
}

func TestGenerateForAgent_Success_Weekly(t *testing.T) {
	ctx := context.Background()
	svc, invoiceRepo, saleRepo, agentRepo, seqRepo := newAgentInvoiceSvcWithMocks()
	agentID := uuid.New()
	// A Monday
	periodStart := time.Date(2024, time.March, 4, 0, 0, 0, 0, time.UTC)
	periodEnd := time.Date(2024, time.March, 11, 0, 0, 0, 0, time.UTC)

	agent := &model.SalesAgent{ID: agentID.String(), Status: "active", RouterID: uuid.New().String()}
	agentRepo.On("GetByID", ctx, agentID).Return(agent, nil)
	invoiceRepo.On("GetByAgentAndPeriod", ctx, agentID, periodStart, "weekly").Return(nil, errors.New("not found"))
	saleRepo.On("SumByAgentAndPeriod", ctx, agentID, periodStart, periodEnd).Return(5, 200000.0, 350000.0, nil)
	seqRepo.On("NextNumber", ctx, "agent_invoice").Return(7, nil)
	invoiceRepo.On("Create", ctx, mock.AnythingOfType("*model.AgentInvoice")).Return(nil)

	inv, err := svc.GenerateForAgent(ctx, agentID, periodStart, periodEnd, "weekly")
	require.NoError(t, err)
	assert.Equal(t, "AGT000007", inv.InvoiceNumber)
	assert.Equal(t, "weekly", inv.BillingCycle)
	assert.NotNil(t, inv.BillingWeek)
	assert.Nil(t, inv.BillingMonth)
}

func TestGenerateForAgent_SequenceError(t *testing.T) {
	ctx := context.Background()
	svc, invoiceRepo, saleRepo, agentRepo, seqRepo := newAgentInvoiceSvcWithMocks()
	agentID := uuid.New()
	periodStart := time.Date(2024, time.March, 1, 0, 0, 0, 0, time.UTC)
	periodEnd := time.Date(2024, time.April, 1, 0, 0, 0, 0, time.UTC)

	agent := &model.SalesAgent{ID: agentID.String(), Status: "active"}
	agentRepo.On("GetByID", ctx, agentID).Return(agent, nil)
	invoiceRepo.On("GetByAgentAndPeriod", ctx, agentID, periodStart, "monthly").Return(nil, errors.New("not found"))
	saleRepo.On("SumByAgentAndPeriod", ctx, agentID, periodStart, periodEnd).Return(3, 100000.0, 150000.0, nil)
	seqRepo.On("NextNumber", ctx, "agent_invoice").Return(0, errors.New("sequence error"))

	_, err := svc.GenerateForAgent(ctx, agentID, periodStart, periodEnd, "monthly")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invoice number")
}

// --- RequestPayment tests ---

func TestRequestPayment_Success(t *testing.T) {
	ctx := context.Background()
	svc, invoiceRepo, _, _, _ := newAgentInvoiceSvcWithMocks()
	invID := uuid.New()
	agentID := uuid.New().String()

	inv := &model.AgentInvoice{ID: invID.String(), AgentID: agentID, Status: "unpaid", TotalAmount: 500000}
	invoiceRepo.On("GetByID", ctx, invID).Return(inv, nil)
	invoiceRepo.On("UpdateStatusAndNotes", ctx, invID, "review", 500000.0, "bukti transfer").Return(nil)

	result, err := svc.RequestPayment(ctx, invID, agentID, 500000, "bukti transfer")
	require.NoError(t, err)
	assert.Equal(t, "review", result.Status)
	assert.Equal(t, 500000.0, result.PaidAmount)
}

func TestRequestPayment_WrongAgent(t *testing.T) {
	ctx := context.Background()
	svc, invoiceRepo, _, _, _ := newAgentInvoiceSvcWithMocks()
	invID := uuid.New()

	inv := &model.AgentInvoice{ID: invID.String(), AgentID: "agent-A", Status: "unpaid"}
	invoiceRepo.On("GetByID", ctx, invID).Return(inv, nil)

	_, err := svc.RequestPayment(ctx, invID, "agent-B", 100000, "")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "access denied")
}

func TestRequestPayment_AlreadyPaid(t *testing.T) {
	ctx := context.Background()
	svc, invoiceRepo, _, _, _ := newAgentInvoiceSvcWithMocks()
	invID := uuid.New()
	agentID := "agent-A"

	inv := &model.AgentInvoice{ID: invID.String(), AgentID: agentID, Status: "paid"}
	invoiceRepo.On("GetByID", ctx, invID).Return(inv, nil)

	_, err := svc.RequestPayment(ctx, invID, agentID, 100000, "")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "paid")
}

func TestRequestPayment_AlreadyCancelled(t *testing.T) {
	ctx := context.Background()
	svc, invoiceRepo, _, _, _ := newAgentInvoiceSvcWithMocks()
	invID := uuid.New()
	agentID := "agent-A"

	inv := &model.AgentInvoice{ID: invID.String(), AgentID: agentID, Status: "cancelled"}
	invoiceRepo.On("GetByID", ctx, invID).Return(inv, nil)

	_, err := svc.RequestPayment(ctx, invID, agentID, 100000, "")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cancelled")
}

// --- MarkPaid tests ---

func TestMarkPaid_FullPayment_Paid(t *testing.T) {
	ctx := context.Background()
	svc, invoiceRepo, _, _, _ := newAgentInvoiceSvcWithMocks()
	invID := uuid.New()

	inv := &model.AgentInvoice{ID: invID.String(), Status: "unpaid", TotalAmount: 500000}
	invoiceRepo.On("GetByID", ctx, invID).Return(inv, nil)
	invoiceRepo.On("UpdateStatus", ctx, invID, "paid", 500000.0).Return(nil)

	result, err := svc.MarkPaid(ctx, invID, 500000)
	require.NoError(t, err)
	assert.Equal(t, "paid", result.Status)
	assert.Equal(t, 500000.0, result.PaidAmount)
}

func TestMarkPaid_PartialPayment_StaysUnpaid(t *testing.T) {
	ctx := context.Background()
	svc, invoiceRepo, _, _, _ := newAgentInvoiceSvcWithMocks()
	invID := uuid.New()

	inv := &model.AgentInvoice{ID: invID.String(), Status: "unpaid", TotalAmount: 500000}
	invoiceRepo.On("GetByID", ctx, invID).Return(inv, nil)
	invoiceRepo.On("UpdateStatus", ctx, invID, "unpaid", 250000.0).Return(nil)

	result, err := svc.MarkPaid(ctx, invID, 250000)
	require.NoError(t, err)
	assert.Equal(t, "unpaid", result.Status)
}

func TestMarkPaid_CancelledInvoice_Error(t *testing.T) {
	ctx := context.Background()
	svc, invoiceRepo, _, _, _ := newAgentInvoiceSvcWithMocks()
	invID := uuid.New()

	inv := &model.AgentInvoice{ID: invID.String(), Status: "cancelled", TotalAmount: 500000}
	invoiceRepo.On("GetByID", ctx, invID).Return(inv, nil)

	_, err := svc.MarkPaid(ctx, invID, 500000)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cancelled")
}

// --- Cancel tests ---

func TestCancel_Success(t *testing.T) {
	ctx := context.Background()
	svc, invoiceRepo, _, _, _ := newAgentInvoiceSvcWithMocks()
	invID := uuid.New()

	inv := &model.AgentInvoice{ID: invID.String(), Status: "unpaid"}
	invoiceRepo.On("GetByID", ctx, invID).Return(inv, nil)
	invoiceRepo.On("UpdateStatus", ctx, invID, "cancelled", 0.0).Return(nil)

	err := svc.Cancel(ctx, invID)
	require.NoError(t, err)
}

func TestCancel_PaidInvoice_Error(t *testing.T) {
	ctx := context.Background()
	svc, invoiceRepo, _, _, _ := newAgentInvoiceSvcWithMocks()
	invID := uuid.New()

	inv := &model.AgentInvoice{ID: invID.String(), Status: "paid"}
	invoiceRepo.On("GetByID", ctx, invID).Return(inv, nil)

	err := svc.Cancel(ctx, invID)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cannot cancel a paid invoice")
}

// --- GenerateManual tests ---

func TestGenerateManual_ShortPeriod_WeeklyCycle(t *testing.T) {
	ctx := context.Background()
	svc, invoiceRepo, saleRepo, agentRepo, seqRepo := newAgentInvoiceSvcWithMocks()
	agentID := uuid.New()
	periodStart := time.Date(2024, time.March, 4, 0, 0, 0, 0, time.UTC)
	periodEnd := time.Date(2024, time.March, 11, 0, 0, 0, 0, time.UTC) // 7 days

	agent := &model.SalesAgent{ID: agentID.String(), Status: "active", RouterID: uuid.New().String()}
	agentRepo.On("GetByID", ctx, agentID).Return(agent, nil)
	invoiceRepo.On("GetByAgentAndPeriod", ctx, agentID, periodStart, "weekly").Return(nil, errors.New("not found"))
	saleRepo.On("SumByAgentAndPeriod", ctx, agentID, periodStart, periodEnd).Return(3, 100000.0, 150000.0, nil)
	seqRepo.On("NextNumber", ctx, "agent_invoice").Return(1, nil)
	invoiceRepo.On("Create", ctx, mock.AnythingOfType("*model.AgentInvoice")).Return(nil)

	inv, err := svc.GenerateManual(ctx, agentID, periodStart, periodEnd)
	require.NoError(t, err)
	assert.Equal(t, "weekly", inv.BillingCycle)
}

// --- ListInvoices tests ---

func TestListInvoices_Success(t *testing.T) {
	ctx := context.Background()
	svc, invoiceRepo, _, _, _ := newAgentInvoiceSvcWithMocks()
	filter := repository.AgentInvoiceFilter{}

	invoices := []model.AgentInvoice{
		{ID: uuid.New().String(), InvoiceNumber: "AGT000001"},
		{ID: uuid.New().String(), InvoiceNumber: "AGT000002"},
	}
	invoiceRepo.On("List", ctx, filter, 10, 0).Return(invoices, nil)
	invoiceRepo.On("Count", ctx, filter).Return(int64(2), nil)

	result, count, err := svc.ListInvoices(ctx, filter, 10, 0)
	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, int64(2), count)
}

// --- GetInvoice tests ---

func TestGetInvoice_Found(t *testing.T) {
	ctx := context.Background()
	svc, invoiceRepo, _, _, _ := newAgentInvoiceSvcWithMocks()
	invID := uuid.New()

	inv := &model.AgentInvoice{ID: invID.String(), InvoiceNumber: "AGT000001"}
	invoiceRepo.On("GetByID", ctx, invID).Return(inv, nil)

	result, err := svc.GetInvoice(ctx, invID)
	require.NoError(t, err)
	assert.Equal(t, "AGT000001", result.InvoiceNumber)
}

func TestGetInvoice_NotFound(t *testing.T) {
	ctx := context.Background()
	svc, invoiceRepo, _, _, _ := newAgentInvoiceSvcWithMocks()
	invID := uuid.New()

	invoiceRepo.On("GetByID", ctx, invID).Return(nil, errors.New("not found"))

	_, err := svc.GetInvoice(ctx, invID)
	require.Error(t, err)
}
