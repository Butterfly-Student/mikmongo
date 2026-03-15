package service

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	billingDomain "mikmongo/internal/domain/billing"
	paymentDomain "mikmongo/internal/domain/payment"
	"mikmongo/internal/model"
	"mikmongo/internal/repository"
	"mikmongo/internal/service/mocks"
)

// funcTransactor executes fn synchronously with the provided mock repos (no real DB).
type funcTransactor struct {
	paymentRepo    repository.PaymentRepository
	invoiceRepo    repository.InvoiceRepository
	allocationRepo repository.PaymentAllocationRepository
}

func (f *funcTransactor) RunInTx(ctx context.Context, fn func(
	repository.PaymentRepository,
	repository.InvoiceRepository,
	repository.PaymentAllocationRepository,
) error) error {
	return fn(f.paymentRepo, f.invoiceRepo, f.allocationRepo)
}

func newPaymentServiceWithMocks() (
	*PaymentService,
	*mocks.MockPaymentRepository,
	*mocks.MockInvoiceRepository,
	*mocks.MockPaymentAllocationRepository,
	*mocks.MockCustomerRepository,
	*mocks.MockSequenceCounterRepository,
) {
	paymentRepo := &mocks.MockPaymentRepository{}
	invoiceRepo := &mocks.MockInvoiceRepository{}
	allocationRepo := &mocks.MockPaymentAllocationRepository{}
	customerRepo := &mocks.MockCustomerRepository{}
	seqRepo := &mocks.MockSequenceCounterRepository{}
	transactor := &funcTransactor{paymentRepo, invoiceRepo, allocationRepo}

	svc := NewPaymentService(
		paymentRepo,
		invoiceRepo,
		allocationRepo,
		customerRepo,
		seqRepo,
		paymentDomain.NewDomain(),
		billingDomain.NewDomain(),
		transactor,
	)
	return svc, paymentRepo, invoiceRepo, allocationRepo, customerRepo, seqRepo
}

func TestCreatePayment(t *testing.T) {
	ctx := context.Background()
	svc, paymentRepo, _, _, _, seqRepo := newPaymentServiceWithMocks()

	seqRepo.On("NextNumber", ctx, "payment_number").Return(1, nil)
	paymentRepo.On("Create", ctx, mock.AnythingOfType("*model.Payment")).Return(nil)

	p := &model.Payment{
		CustomerID:    uuid.New().String(),
		Amount:        100_000,
		PaymentMethod: "cash",
	}

	err := svc.Create(ctx, p)
	require.NoError(t, err)
	assert.Equal(t, "PAY000001", p.PaymentNumber)
	assert.Equal(t, "pending", p.Status)
	assert.False(t, p.PaymentDate.IsZero())
}

func TestCreatePayment_GeneratesPaymentNumber(t *testing.T) {
	ctx := context.Background()
	svc, paymentRepo, _, _, _, seqRepo := newPaymentServiceWithMocks()

	seqRepo.On("NextNumber", ctx, "payment_number").Return(42, nil)
	paymentRepo.On("Create", ctx, mock.AnythingOfType("*model.Payment")).Return(nil)

	p := &model.Payment{
		CustomerID:    uuid.New().String(),
		Amount:        50_000,
		PaymentMethod: "bank_transfer",
	}

	err := svc.Create(ctx, p)
	require.NoError(t, err)
	assert.Equal(t, "PAY000042", p.PaymentNumber)
}

func TestConfirmPayment_SingleInvoice(t *testing.T) {
	ctx := context.Background()
	svc, paymentRepo, invoiceRepo, allocationRepo, customerRepo, _ := newPaymentServiceWithMocks()

	paymentID := uuid.New()
	customerID := uuid.New()
	invoiceID := uuid.New()

	p := &model.Payment{
		ID:            paymentID.String(),
		CustomerID:    customerID.String(),
		Amount:        100_000,
		PaymentMethod: "cash",
		Status:        "pending",
	}

	inv := model.Invoice{
		ID:          invoiceID.String(),
		CustomerID:  customerID.String(),
		Status:      "unpaid",
		TotalAmount: 100_000,
		PaidAmount:  0,
		DueDate:     time.Now().AddDate(0, 0, 10),
	}

	updatedInv := inv
	updatedInv.PaidAmount = 100_000

	paymentRepo.On("GetByID", ctx, paymentID).Return(p, nil)
	invoiceRepo.On("GetByCustomerIDForUpdate", ctx, customerID).Return([]model.Invoice{inv}, nil)
	invoiceRepo.On("GetByID", ctx, invoiceID).Return(&updatedInv, nil)
	allocationRepo.On("Create", ctx, mock.AnythingOfType("*model.PaymentAllocation")).Return(nil)
	invoiceRepo.On("Update", ctx, mock.AnythingOfType("*model.Invoice")).Return(nil)
	invoiceRepo.On("UpdateStatus", ctx, invoiceID, "paid").Return(nil)
	paymentRepo.On("Update", ctx, mock.AnythingOfType("*model.Payment")).Return(nil)
	customerRepo.On("GetByID", ctx, customerID).Return(
		&model.Customer{ID: customerID.String(), FullName: "Budi", Phone: "081234"}, nil,
	)

	err := svc.Confirm(ctx, paymentID, "admin-001")
	require.NoError(t, err)

	paymentRepo.AssertCalled(t, "Update", ctx, mock.AnythingOfType("*model.Payment"))
}

func TestConfirmPayment_AlreadyConfirmed(t *testing.T) {
	ctx := context.Background()
	svc, paymentRepo, _, _, _, _ := newPaymentServiceWithMocks()

	paymentID := uuid.New()
	p := &model.Payment{
		ID:     paymentID.String(),
		Status: "confirmed",
	}

	paymentRepo.On("GetByID", ctx, paymentID).Return(p, nil)

	err := svc.Confirm(ctx, paymentID, "admin-001")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "pending")
}

func TestConfirmPayment_MultipleInvoices_FIFO(t *testing.T) {
	ctx := context.Background()
	svc, paymentRepo, invoiceRepo, allocationRepo, customerRepo, _ := newPaymentServiceWithMocks()

	paymentID := uuid.New()
	customerID := uuid.New()
	inv1ID := uuid.New()
	inv2ID := uuid.New()

	p := &model.Payment{
		ID:         paymentID.String(),
		CustomerID: customerID.String(),
		Amount:     120_000,
		Status:     "pending",
	}

	// Older invoice (lower due date = oldest)
	inv1 := model.Invoice{
		ID:          inv1ID.String(),
		CustomerID:  customerID.String(),
		Status:      "unpaid",
		TotalAmount: 100_000,
		PaidAmount:  0,
		DueDate:     time.Now().AddDate(0, -1, 0),
	}
	// Newer invoice
	inv2 := model.Invoice{
		ID:          inv2ID.String(),
		CustomerID:  customerID.String(),
		Status:      "unpaid",
		TotalAmount: 100_000,
		PaidAmount:  0,
		DueDate:     time.Now().AddDate(0, 0, 10),
	}

	inv1Updated := inv1
	inv1Updated.PaidAmount = 100_000
	inv2Updated := inv2
	inv2Updated.PaidAmount = 20_000

	paymentRepo.On("GetByID", ctx, paymentID).Return(p, nil)
	invoiceRepo.On("GetByCustomerIDForUpdate", ctx, customerID).Return([]model.Invoice{inv1, inv2}, nil)
	invoiceRepo.On("GetByID", ctx, inv1ID).Return(&inv1Updated, nil)
	invoiceRepo.On("GetByID", ctx, inv2ID).Return(&inv2Updated, nil)
	allocationRepo.On("Create", ctx, mock.AnythingOfType("*model.PaymentAllocation")).Return(nil)
	invoiceRepo.On("Update", ctx, mock.AnythingOfType("*model.Invoice")).Return(nil)
	invoiceRepo.On("UpdateStatus", ctx, inv1ID, mock.Anything).Return(nil)
	invoiceRepo.On("UpdateStatus", ctx, inv2ID, mock.Anything).Return(nil)
	paymentRepo.On("Update", ctx, mock.AnythingOfType("*model.Payment")).Return(nil)
	customerRepo.On("GetByID", ctx, customerID).Return(
		&model.Customer{ID: customerID.String()}, nil,
	)

	err := svc.Confirm(ctx, paymentID, "admin-001")
	require.NoError(t, err)
}

func TestRejectPayment(t *testing.T) {
	ctx := context.Background()
	svc, paymentRepo, _, _, _, _ := newPaymentServiceWithMocks()

	paymentID := uuid.New()
	p := &model.Payment{
		ID:     paymentID.String(),
		Status: "pending",
	}

	paymentRepo.On("GetByID", ctx, paymentID).Return(p, nil)
	paymentRepo.On("Update", ctx, mock.AnythingOfType("*model.Payment")).Return(nil)

	err := svc.Reject(ctx, paymentID, "bukti tidak valid")
	require.NoError(t, err)
	assert.Equal(t, "rejected", p.Status)
	require.NotNil(t, p.RejectionReason)
	assert.Equal(t, "bukti tidak valid", *p.RejectionReason)
}

func TestRejectPayment_AlreadyConfirmed(t *testing.T) {
	ctx := context.Background()
	svc, paymentRepo, _, _, _, _ := newPaymentServiceWithMocks()

	paymentID := uuid.New()
	p := &model.Payment{
		ID:     paymentID.String(),
		Status: "confirmed",
	}

	paymentRepo.On("GetByID", ctx, paymentID).Return(p, nil)

	err := svc.Reject(ctx, paymentID, "reason")
	assert.Error(t, err)
}

func TestRefundPayment(t *testing.T) {
	ctx := context.Background()
	svc, paymentRepo, _, _, _, _ := newPaymentServiceWithMocks()

	paymentID := uuid.New()
	p := &model.Payment{
		ID:           paymentID.String(),
		Status:       "confirmed",
		RefundAmount: 0,
	}

	paymentRepo.On("GetByID", ctx, paymentID).Return(p, nil)
	paymentRepo.On("Update", ctx, mock.AnythingOfType("*model.Payment")).Return(nil)

	err := svc.Refund(ctx, paymentID, 100_000, "customer request")
	require.NoError(t, err)
	assert.Equal(t, "refunded", p.Status)
	assert.Equal(t, 100_000.0, p.RefundAmount)
	require.NotNil(t, p.RefundReason)
	assert.Equal(t, "customer request", *p.RefundReason)
}

func TestRefundPayment_AlreadyRefunded(t *testing.T) {
	ctx := context.Background()
	svc, paymentRepo, _, _, _, _ := newPaymentServiceWithMocks()

	paymentID := uuid.New()
	p := &model.Payment{
		ID:           paymentID.String(),
		Status:       "confirmed",
		RefundAmount: 50_000,
	}

	paymentRepo.On("GetByID", ctx, paymentID).Return(p, nil)

	err := svc.Refund(ctx, paymentID, 50_000, "reason")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "refunded")
}
