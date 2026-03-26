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

func newCashSvcWithMocks() (
	*CashManagementService,
	*mocks.MockCashEntryRepository,
	*mocks.MockPettyCashFundRepository,
	*mocks.MockSequenceCounterRepository,
) {
	cashRepo := &mocks.MockCashEntryRepository{}
	fundRepo := &mocks.MockPettyCashFundRepository{}
	seqRepo := &mocks.MockSequenceCounterRepository{}
	// db is nil — methods that use s.db (CreateEntry with petty cash, UpdateEntry, GetReconciliation)
	// are tested via integration tests.
	svc := NewCashManagementService(cashRepo, fundRepo, seqRepo, nil)
	return svc, cashRepo, fundRepo, seqRepo
}

// --- RecordPaymentIncome tests ---

func TestRecordPaymentIncome_Success(t *testing.T) {
	ctx := context.Background()
	svc, cashRepo, _, seqRepo := newCashSvcWithMocks()

	paymentID := uuid.New()
	payment := &model.Payment{
		ID:            paymentID.String(),
		PaymentNumber: "PAY000001",
		Amount:        100000,
		PaymentMethod: "cash",
		PaymentDate:   time.Now(),
	}

	cashRepo.On("GetByReference", ctx, "payment", paymentID).Return(nil, errors.New("not found"))
	seqRepo.On("NextNumber", ctx, "cash_entry_number").Return(1, nil)
	cashRepo.On("Create", ctx, mock.AnythingOfType("*model.CashEntry")).Return(nil)

	err := svc.RecordPaymentIncome(ctx, payment)
	require.NoError(t, err)

	cashRepo.AssertCalled(t, "Create", ctx, mock.AnythingOfType("*model.CashEntry"))
}

func TestRecordPaymentIncome_IdempotentSkip(t *testing.T) {
	ctx := context.Background()
	svc, cashRepo, _, _ := newCashSvcWithMocks()

	paymentID := uuid.New()
	payment := &model.Payment{
		ID:            paymentID.String(),
		PaymentNumber: "PAY000001",
		Amount:        100000,
	}

	existing := &model.CashEntry{ID: uuid.New().String()}
	cashRepo.On("GetByReference", ctx, "payment", paymentID).Return(existing, nil)

	err := svc.RecordPaymentIncome(ctx, payment)
	require.NoError(t, err)
	cashRepo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestRecordPaymentIncome_SequenceError(t *testing.T) {
	ctx := context.Background()
	svc, cashRepo, _, seqRepo := newCashSvcWithMocks()

	paymentID := uuid.New()
	payment := &model.Payment{
		ID:            paymentID.String(),
		PaymentNumber: "PAY000001",
		Amount:        100000,
	}

	cashRepo.On("GetByReference", ctx, "payment", paymentID).Return(nil, errors.New("not found"))
	seqRepo.On("NextNumber", ctx, "cash_entry_number").Return(0, errors.New("seq error"))

	err := svc.RecordPaymentIncome(ctx, payment)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "entry number")
}

// --- RecordAgentInvoiceIncome tests ---

func TestRecordAgentInvoiceIncome_Success(t *testing.T) {
	ctx := context.Background()
	svc, cashRepo, _, seqRepo := newCashSvcWithMocks()

	invID := uuid.New()
	inv := &model.AgentInvoice{
		ID:            invID.String(),
		InvoiceNumber: "AGT000001",
		TotalAmount:   500000,
	}

	cashRepo.On("GetByReference", ctx, "agent_invoice", invID).Return(nil, errors.New("not found"))
	seqRepo.On("NextNumber", ctx, "cash_entry_number").Return(2, nil)
	cashRepo.On("Create", ctx, mock.AnythingOfType("*model.CashEntry")).Return(nil)

	err := svc.RecordAgentInvoiceIncome(ctx, inv, 500000)
	require.NoError(t, err)
	cashRepo.AssertCalled(t, "Create", ctx, mock.AnythingOfType("*model.CashEntry"))
}

func TestRecordAgentInvoiceIncome_IdempotentSkip(t *testing.T) {
	ctx := context.Background()
	svc, cashRepo, _, _ := newCashSvcWithMocks()

	invID := uuid.New()
	inv := &model.AgentInvoice{
		ID:            invID.String(),
		InvoiceNumber: "AGT000001",
	}

	existing := &model.CashEntry{ID: uuid.New().String()}
	cashRepo.On("GetByReference", ctx, "agent_invoice", invID).Return(existing, nil)

	err := svc.RecordAgentInvoiceIncome(ctx, inv, 500000)
	require.NoError(t, err)
	cashRepo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

// --- GetEntry tests ---

func TestGetEntry_Found(t *testing.T) {
	ctx := context.Background()
	svc, cashRepo, _, _ := newCashSvcWithMocks()
	entryID := uuid.New()

	entry := &model.CashEntry{ID: entryID.String(), EntryNumber: "KAS000001"}
	cashRepo.On("GetByID", ctx, entryID).Return(entry, nil)

	result, err := svc.GetEntry(ctx, entryID)
	require.NoError(t, err)
	assert.Equal(t, "KAS000001", result.EntryNumber)
}

func TestGetEntry_NotFound(t *testing.T) {
	ctx := context.Background()
	svc, cashRepo, _, _ := newCashSvcWithMocks()
	entryID := uuid.New()

	cashRepo.On("GetByID", ctx, entryID).Return(nil, errors.New("not found"))

	_, err := svc.GetEntry(ctx, entryID)
	require.Error(t, err)
}

// --- ListEntries tests ---

func TestListEntries_Success(t *testing.T) {
	ctx := context.Background()
	svc, cashRepo, _, _ := newCashSvcWithMocks()
	filter := repository.CashEntryFilter{}

	entries := []model.CashEntry{
		{ID: uuid.New().String(), EntryNumber: "KAS000001"},
		{ID: uuid.New().String(), EntryNumber: "KAS000002"},
	}
	cashRepo.On("List", ctx, filter, 10, 0).Return(entries, nil)
	cashRepo.On("Count", ctx, filter).Return(int64(2), nil)

	result, count, err := svc.ListEntries(ctx, filter, 10, 0)
	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, int64(2), count)
}

func TestListEntries_Empty(t *testing.T) {
	ctx := context.Background()
	svc, cashRepo, _, _ := newCashSvcWithMocks()
	filter := repository.CashEntryFilter{}

	cashRepo.On("List", ctx, filter, 10, 0).Return([]model.CashEntry{}, nil)
	cashRepo.On("Count", ctx, filter).Return(int64(0), nil)

	result, count, err := svc.ListEntries(ctx, filter, 10, 0)
	require.NoError(t, err)
	assert.Empty(t, result)
	assert.Equal(t, int64(0), count)
}

func TestListEntries_RepoError(t *testing.T) {
	ctx := context.Background()
	svc, cashRepo, _, _ := newCashSvcWithMocks()
	filter := repository.CashEntryFilter{}

	cashRepo.On("List", ctx, filter, 10, 0).Return(nil, errors.New("db error"))

	_, _, err := svc.ListEntries(ctx, filter, 10, 0)
	require.Error(t, err)
}

// --- Approve tests ---

func TestApprove_Success(t *testing.T) {
	ctx := context.Background()
	svc, cashRepo, _, _ := newCashSvcWithMocks()
	entryID := uuid.New()
	approverID := uuid.New().String()

	entry := &model.CashEntry{ID: entryID.String(), Status: "pending"}
	cashRepo.On("GetByID", ctx, entryID).Return(entry, nil)
	cashRepo.On("Update", ctx, mock.AnythingOfType("*model.CashEntry")).Return(nil)

	result, err := svc.Approve(ctx, entryID, approverID)
	require.NoError(t, err)
	assert.Equal(t, "approved", result.Status)
	assert.NotNil(t, result.ApprovedBy)
	assert.Equal(t, approverID, *result.ApprovedBy)
	assert.NotNil(t, result.ApprovedAt)
}

func TestApprove_NotPending(t *testing.T) {
	ctx := context.Background()
	svc, cashRepo, _, _ := newCashSvcWithMocks()
	entryID := uuid.New()

	entry := &model.CashEntry{ID: entryID.String(), Status: "approved"}
	cashRepo.On("GetByID", ctx, entryID).Return(entry, nil)

	_, err := svc.Approve(ctx, entryID, "approver-id")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "pending")
}

func TestApprove_NotFound(t *testing.T) {
	ctx := context.Background()
	svc, cashRepo, _, _ := newCashSvcWithMocks()
	entryID := uuid.New()

	cashRepo.On("GetByID", ctx, entryID).Return(nil, errors.New("not found"))

	_, err := svc.Approve(ctx, entryID, "approver-id")
	require.Error(t, err)
}

// --- Reject tests ---

func TestReject_Success(t *testing.T) {
	ctx := context.Background()
	svc, cashRepo, _, _ := newCashSvcWithMocks()
	entryID := uuid.New()

	entry := &model.CashEntry{ID: entryID.String(), Status: "pending", Type: "income"}
	cashRepo.On("GetByID", ctx, entryID).Return(entry, nil)
	cashRepo.On("Update", ctx, mock.AnythingOfType("*model.CashEntry")).Return(nil)

	result, err := svc.Reject(ctx, entryID, "invalid receipt")
	require.NoError(t, err)
	assert.Equal(t, "rejected", result.Status)
	assert.NotNil(t, result.Notes)
	assert.Equal(t, "invalid receipt", *result.Notes)
}

func TestReject_NotPending(t *testing.T) {
	ctx := context.Background()
	svc, cashRepo, _, _ := newCashSvcWithMocks()
	entryID := uuid.New()

	entry := &model.CashEntry{ID: entryID.String(), Status: "approved"}
	cashRepo.On("GetByID", ctx, entryID).Return(entry, nil)

	_, err := svc.Reject(ctx, entryID, "reason")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "pending")
}

// --- CreateFund tests ---

func TestCreateFund_Success(t *testing.T) {
	ctx := context.Background()
	svc, _, fundRepo, _ := newCashSvcWithMocks()

	fund := &model.PettyCashFund{
		FundName:       "Operasional Kantor",
		InitialBalance: 1000000,
		CustodianID:    uuid.New().String(),
	}
	fundRepo.On("Create", ctx, fund).Return(nil)

	err := svc.CreateFund(ctx, fund)
	require.NoError(t, err)
	assert.Equal(t, 1000000.0, fund.CurrentBalance, "CurrentBalance should be set to InitialBalance")
}

// --- TopUpFund tests ---

func TestTopUpFund_Success(t *testing.T) {
	ctx := context.Background()
	svc, _, fundRepo, _ := newCashSvcWithMocks()
	fundID := uuid.New()

	fundRepo.On("AdjustBalance", ctx, fundID, 500000.0).Return(nil)

	err := svc.TopUpFund(ctx, fundID, 500000)
	require.NoError(t, err)
}

func TestTopUpFund_NegativeAmount(t *testing.T) {
	ctx := context.Background()
	svc, _, _, _ := newCashSvcWithMocks()
	fundID := uuid.New()

	err := svc.TopUpFund(ctx, fundID, -100)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "positive")
}

// --- GetCashBalance tests ---

func TestGetCashBalance_Success(t *testing.T) {
	ctx := context.Background()
	svc, cashRepo, _, _ := newCashSvcWithMocks()

	cashRepo.On("SumByTypeAndPeriod", ctx, "income", mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).Return(5000000.0, nil)
	cashRepo.On("SumByTypeAndPeriod", ctx, "expense", mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).Return(2000000.0, nil)

	result, err := svc.GetCashBalance(ctx)
	require.NoError(t, err)
	assert.Equal(t, 5000000.0, result.TotalIncome)
	assert.Equal(t, 2000000.0, result.TotalExpense)
	assert.Equal(t, 3000000.0, result.Balance)
}

// --- GetCashFlow tests ---

func TestGetCashFlow_Success(t *testing.T) {
	ctx := context.Background()
	svc, cashRepo, _, _ := newCashSvcWithMocks()

	from := time.Date(2024, time.March, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, time.April, 1, 0, 0, 0, 0, time.UTC)

	cashRepo.On("SumByTypeAndPeriod", ctx, "income", from, to).Return(3000000.0, nil)
	cashRepo.On("SumByTypeAndPeriod", ctx, "expense", from, to).Return(1000000.0, nil)
	cashRepo.On("SumBySourceAndPeriod", ctx, from, to).Return([]repository.SourceSum{
		{Source: "invoice", Total: 2500000},
		{Source: "agent_invoice", Total: 500000},
	}, nil)

	result, err := svc.GetCashFlow(ctx, from, to)
	require.NoError(t, err)
	assert.Equal(t, 3000000.0, result.TotalIncome)
	assert.Equal(t, 1000000.0, result.TotalExpense)
	assert.Equal(t, 2000000.0, result.NetCashFlow)
	assert.Len(t, result.Breakdown, 2)
}

// --- DeleteEntry tests ---

func TestDeleteEntry_Success(t *testing.T) {
	ctx := context.Background()
	svc, cashRepo, _, _ := newCashSvcWithMocks()
	entryID := uuid.New()

	entry := &model.CashEntry{ID: entryID.String(), Status: "pending", Type: "income"}
	cashRepo.On("GetByID", ctx, entryID).Return(entry, nil)
	cashRepo.On("Delete", ctx, entryID).Return(nil)

	err := svc.DeleteEntry(ctx, entryID)
	require.NoError(t, err)
}

func TestDeleteEntry_NotPending(t *testing.T) {
	ctx := context.Background()
	svc, cashRepo, _, _ := newCashSvcWithMocks()
	entryID := uuid.New()

	entry := &model.CashEntry{ID: entryID.String(), Status: "approved"}
	cashRepo.On("GetByID", ctx, entryID).Return(entry, nil)

	err := svc.DeleteEntry(ctx, entryID)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "pending")
}

// --- ListFunds tests ---

func TestListFunds_Success(t *testing.T) {
	ctx := context.Background()
	svc, _, fundRepo, _ := newCashSvcWithMocks()

	funds := []model.PettyCashFund{
		{ID: uuid.New().String(), FundName: "Fund A"},
		{ID: uuid.New().String(), FundName: "Fund B"},
	}
	fundRepo.On("List", ctx, 10, 0).Return(funds, nil)
	fundRepo.On("Count", ctx).Return(int64(2), nil)

	result, count, err := svc.ListFunds(ctx, 10, 0)
	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, int64(2), count)
}
