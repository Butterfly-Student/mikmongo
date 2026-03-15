//go:build integration

package integration

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"mikmongo/internal/domain"
	"mikmongo/internal/model"
	"mikmongo/internal/repository/postgres"
	"mikmongo/internal/service"
)

// setupPaymentTestFixtures creates common test resources for payment tests.
// Also returns a valid admin user ID for use in Confirm/Reject/Refund calls.
func setupPaymentTestFixtures(t *testing.T, suite *TestSuite, suffix string) (
	*service.PaymentService,
	*service.BillingService,
	*model.Customer,
	*model.Subscription,
	*model.BandwidthProfile,
	string, // adminUserID
) {
	t.Helper()
	repos := postgres.NewRepository(suite.DB)
	logger := zap.NewNop()

	routerSvc := service.NewRouterService(repos.RouterDeviceRepo, "test-key-16-bytes", nil, logger)
	customerSvc := service.NewCustomerService(repos.CustomerRepo, repos.SequenceCounterRepo, repos.BandwidthProfileRepo, domain.NewCustomerDomain(), routerSvc)

	transactor := postgres.NewTransactor(suite.DB)

	paymentSvc := service.NewPaymentService(
		repos.PaymentRepo,
		repos.InvoiceRepo,
		repos.PaymentAllocationRepo,
		repos.CustomerRepo,
		repos.SequenceCounterRepo,
		domain.NewPaymentDomain(),
		domain.NewBillingDomain(),
		transactor,
	)
	paymentSvc.SetCustomerService(customerSvc)

	billingSvc := service.NewBillingService(
		repos.InvoiceRepo,
		repos.InvoiceItemRepo,
		repos.SubscriptionRepo,
		repos.BandwidthProfileRepo,
		repos.CustomerRepo,
		repos.SystemSettingRepo,
		repos.SequenceCounterRepo,
		domain.NewBillingDomain(),
	)

	// Create test router
	router := &model.MikrotikRouter{
		ID:                uuid.New().String(),
		Name:              "Pay Router " + suffix,
		Address:           "192.168.88.1",
		APIPort:           8728,
		Username:          "admin",
		PasswordEncrypted: "enc_pass",
		IsActive:          true,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
	require.NoError(t, repos.RouterDeviceRepo.Create(suite.Ctx, router))

	profile := &model.BandwidthProfile{
		ID:              uuid.New().String(),
		RouterID:        router.ID,
		ProfileCode:     "PAY10" + suffix,
		Name:            "Payment Test " + suffix,
		DownloadSpeed:   10000,
		UploadSpeed:     10000,
		PriceMonthly:    200_000,
		TaxRate:         0,
		GracePeriodDays: 3,
		IsActive:        true,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	require.NoError(t, repos.BandwidthProfileRepo.Create(suite.Ctx, profile))

	customer := &model.Customer{FullName: "Pay Customer " + suffix, Phone: "08" + suffix}
	require.NoError(t, customerSvc.Create(suite.Ctx, customer))

	sub := &model.Subscription{
		CustomerID: customer.ID,
		PlanID:     profile.ID,
		RouterID:   router.ID,
		Username:   "pay-user-" + suffix,
		Password:   "password123",
	}
	directCreateSub(t, suite, sub)
	directActivate(t, suite, sub.ID)

	// Refresh sub
	subID, _ := uuid.Parse(sub.ID)
	activatedSub, err := repos.SubscriptionRepo.GetByID(suite.Ctx, subID)
	require.NoError(t, err)

	adminID := createTestUser(t, suite)
	return paymentSvc, billingSvc, customer, activatedSub, profile, adminID
}

func TestPaymentLifecycle_ConfirmSingle(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	repos := postgres.NewRepository(suite.DB)
	paymentSvc, billingSvc, customer, sub, _, adminID := setupPaymentTestFixtures(t, suite, "001")

	// Generate invoice
	subID, _ := uuid.Parse(sub.ID)
	inv, err := billingSvc.GenerateInvoice(suite.Ctx, subID, time.Now())
	require.NoError(t, err)

	// Create payment
	payment := &model.Payment{
		CustomerID:    customer.ID,
		Amount:        inv.TotalAmount,
		PaymentMethod: "cash",
	}
	err = paymentSvc.Create(suite.Ctx, payment)
	require.NoError(t, err)
	assert.Equal(t, "pending", payment.Status)
	assert.NotEmpty(t, payment.PaymentNumber)

	// Confirm payment
	paymentID, _ := uuid.Parse(payment.ID)
	err = paymentSvc.Confirm(suite.Ctx, paymentID, adminID)
	require.NoError(t, err)

	// Verify invoice is now paid
	invID, _ := uuid.Parse(inv.ID)
	updatedInv, err := repos.InvoiceRepo.GetByID(suite.Ctx, invID)
	require.NoError(t, err)
	assert.Equal(t, "paid", updatedInv.Status)
	assert.Equal(t, inv.TotalAmount, updatedInv.PaidAmount)

	// Verify payment allocation created
	allocations, err := repos.PaymentAllocationRepo.ListByPaymentID(suite.Ctx, paymentID)
	require.NoError(t, err)
	assert.Len(t, allocations, 1)
	assert.Equal(t, inv.TotalAmount, allocations[0].AllocatedAmount)
}

func TestPaymentLifecycle_PartialPayment(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	repos := postgres.NewRepository(suite.DB)
	paymentSvc, billingSvc, customer, sub, _, adminID := setupPaymentTestFixtures(t, suite, "002")

	subID, _ := uuid.Parse(sub.ID)
	inv, err := billingSvc.GenerateInvoice(suite.Ctx, subID, time.Now())
	require.NoError(t, err)

	// First partial payment
	partialAmount := inv.TotalAmount / 2
	payment1 := &model.Payment{
		CustomerID:    customer.ID,
		Amount:        partialAmount,
		PaymentMethod: "bank_transfer",
	}
	require.NoError(t, paymentSvc.Create(suite.Ctx, payment1))
	paymentID1, _ := uuid.Parse(payment1.ID)
	require.NoError(t, paymentSvc.Confirm(suite.Ctx, paymentID1, adminID))

	// Check invoice is partial
	invID, _ := uuid.Parse(inv.ID)
	updatedInv, err := repos.InvoiceRepo.GetByID(suite.Ctx, invID)
	require.NoError(t, err)
	assert.Equal(t, "partial", updatedInv.Status)

	// Second payment to complete
	remaining := inv.TotalAmount - partialAmount
	payment2 := &model.Payment{
		CustomerID:    customer.ID,
		Amount:        remaining,
		PaymentMethod: "cash",
	}
	require.NoError(t, paymentSvc.Create(suite.Ctx, payment2))
	paymentID2, _ := uuid.Parse(payment2.ID)
	require.NoError(t, paymentSvc.Confirm(suite.Ctx, paymentID2, adminID))

	// Invoice should now be paid
	updatedInv2, err := repos.InvoiceRepo.GetByID(suite.Ctx, invID)
	require.NoError(t, err)
	assert.Equal(t, "paid", updatedInv2.Status)
}

func TestPaymentLifecycle_FIFOMultipleInvoices(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	repos := postgres.NewRepository(suite.DB)
	paymentSvc, _, customer, sub, _, adminID := setupPaymentTestFixtures(t, suite, "003")

	// Create two invoices with different due dates (simulating two months of billing)
	subIDStr := sub.ID
	olderDue := time.Now().AddDate(0, -1, 0)
	newerDue := time.Now().AddDate(0, 0, 10)

	inv1 := &model.Invoice{
		InvoiceNumber:      "INV-FIFO-OLD-001",
		CustomerID:         customer.ID,
		SubscriptionID:     &subIDStr,
		BillingPeriodStart: olderDue.AddDate(0, -1, 0),
		BillingPeriodEnd:   olderDue,
		IssueDate:          olderDue.AddDate(0, 0, -10),
		DueDate:            olderDue,
		TotalAmount:        200_000,
		Status:             "unpaid",
		InvoiceType:        "recurring",
	}
	inv2 := &model.Invoice{
		InvoiceNumber:      "INV-FIFO-NEW-001",
		CustomerID:         customer.ID,
		SubscriptionID:     &subIDStr,
		BillingPeriodStart: newerDue.AddDate(0, -1, 0),
		BillingPeriodEnd:   newerDue,
		IssueDate:          newerDue.AddDate(0, 0, -10),
		DueDate:            newerDue,
		TotalAmount:        200_000,
		Status:             "unpaid",
		InvoiceType:        "recurring",
	}
	require.NoError(t, repos.InvoiceRepo.Create(suite.Ctx, inv1))
	require.NoError(t, repos.InvoiceRepo.Create(suite.Ctx, inv2))

	// Payment that covers oldest invoice + part of newer
	payment := &model.Payment{
		CustomerID:    customer.ID,
		Amount:        250_000, // enough for inv1 + partial inv2
		PaymentMethod: "bank_transfer",
	}
	require.NoError(t, paymentSvc.Create(suite.Ctx, payment))
	paymentID, _ := uuid.Parse(payment.ID)
	require.NoError(t, paymentSvc.Confirm(suite.Ctx, paymentID, adminID))

	// Older invoice should be paid
	inv1ID, _ := uuid.Parse(inv1.ID)
	updatedInv1, err := repos.InvoiceRepo.GetByID(suite.Ctx, inv1ID)
	require.NoError(t, err)
	assert.Equal(t, "paid", updatedInv1.Status)

	// Newer invoice should be partial
	inv2ID, _ := uuid.Parse(inv2.ID)
	updatedInv2, err := repos.InvoiceRepo.GetByID(suite.Ctx, inv2ID)
	require.NoError(t, err)
	assert.Equal(t, "partial", updatedInv2.Status)
}

func TestPaymentLifecycle_RejectPayment(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	repos := postgres.NewRepository(suite.DB)
	paymentSvc, billingSvc, customer, sub, _, _ := setupPaymentTestFixtures(t, suite, "004")

	subID, _ := uuid.Parse(sub.ID)
	inv, err := billingSvc.GenerateInvoice(suite.Ctx, subID, time.Now())
	require.NoError(t, err)

	payment := &model.Payment{
		CustomerID:    customer.ID,
		Amount:        inv.TotalAmount,
		PaymentMethod: "bank_transfer",
	}
	require.NoError(t, paymentSvc.Create(suite.Ctx, payment))

	paymentID, _ := uuid.Parse(payment.ID)
	err = paymentSvc.Reject(suite.Ctx, paymentID, "bukti tidak valid")
	require.NoError(t, err)

	// Verify payment is rejected
	updatedPayment, err := repos.PaymentRepo.GetByID(suite.Ctx, paymentID)
	require.NoError(t, err)
	assert.Equal(t, "rejected", updatedPayment.Status)
	require.NotNil(t, updatedPayment.RejectionReason)
	assert.Equal(t, "bukti tidak valid", *updatedPayment.RejectionReason)

	// Invoice should still be unpaid
	invID, _ := uuid.Parse(inv.ID)
	updatedInv, err := repos.InvoiceRepo.GetByID(suite.Ctx, invID)
	require.NoError(t, err)
	assert.Equal(t, "unpaid", updatedInv.Status)
}

func TestPaymentLifecycle_Refund(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	repos := postgres.NewRepository(suite.DB)
	paymentSvc, billingSvc, customer, sub, _, adminID := setupPaymentTestFixtures(t, suite, "005")

	subID, _ := uuid.Parse(sub.ID)
	inv, err := billingSvc.GenerateInvoice(suite.Ctx, subID, time.Now())
	require.NoError(t, err)

	payment := &model.Payment{
		CustomerID:    customer.ID,
		Amount:        inv.TotalAmount,
		PaymentMethod: "cash",
	}
	require.NoError(t, paymentSvc.Create(suite.Ctx, payment))

	paymentID, _ := uuid.Parse(payment.ID)
	require.NoError(t, paymentSvc.Confirm(suite.Ctx, paymentID, adminID))

	// Refund the payment
	err = paymentSvc.Refund(suite.Ctx, paymentID, inv.TotalAmount, "customer request")
	require.NoError(t, err)

	updatedPayment, err := repos.PaymentRepo.GetByID(suite.Ctx, paymentID)
	require.NoError(t, err)
	assert.Equal(t, "refunded", updatedPayment.Status)
	assert.Equal(t, inv.TotalAmount, updatedPayment.RefundAmount)
}
