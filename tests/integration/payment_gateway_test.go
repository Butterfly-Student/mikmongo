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
	gateway "mikmongo/pkg/payment"
)

// setupPaymentGatewayServices creates the services needed for gateway tests.
func setupPaymentGatewayServices(t *testing.T, suite *TestSuite) *service.PaymentService {
	t.Helper()
	repos := postgres.NewRepository(suite.DB)
	logger := zap.NewNop()

	routerSvc := service.NewRouterService(repos.RouterDeviceRepo, "test-key-16-bytes", nil, logger)
	customerSvc := service.NewCustomerService(
		repos.CustomerRepo, repos.SequenceCounterRepo, repos.BandwidthProfileRepo,
		domain.NewCustomerDomain(), routerSvc,
	)

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
	return paymentSvc
}

// createTestCustomerAndPayment creates a customer and a pending payment in the DB.
func createTestCustomerAndPayment(t *testing.T, suite *TestSuite, repos *postgres.Registry, suffix string) (*model.Customer, *model.Payment) {
	t.Helper()

	customer := &model.Customer{
		ID:       uuid.New().String(),
		FullName: "Gateway Test Customer " + suffix,
		Phone:    "0812000" + suffix,
	}
	require.NoError(t, repos.CustomerRepo.Create(suite.Ctx, customer))

	payment := &model.Payment{
		ID:            uuid.New().String(),
		PaymentNumber: "PAY-GW-" + suffix,
		CustomerID:    customer.ID,
		Amount:        150000,
		PaymentMethod: "gateway",
		PaymentDate:   time.Now(),
		Status:        "pending",
	}
	require.NoError(t, repos.PaymentRepo.Create(suite.Ctx, payment))
	return customer, payment
}

// TestSetGatewayInfo verifies that SetGatewayInfo persists all gateway fields.
func TestSetGatewayInfo(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	repos := postgres.NewRepository(suite.DB)
	paymentSvc := setupPaymentGatewayServices(t, suite)
	_, payment := createTestCustomerAndPayment(t, suite, repos, "001")

	paymentID, err := uuid.Parse(payment.ID)
	require.NoError(t, err)

	err = paymentSvc.SetGatewayInfo(
		suite.Ctx, paymentID,
		"xendit",
		"inv-test-001",
		"https://checkout.xendit.co/web/inv-test-001",
		`{"id":"inv-test-001","status":"PENDING"}`,
	)
	require.NoError(t, err)

	// Reload from DB and verify fields
	updated, err := repos.PaymentRepo.GetByID(suite.Ctx, paymentID)
	require.NoError(t, err)

	require.NotNil(t, updated.GatewayName)
	assert.Equal(t, "xendit", *updated.GatewayName)

	require.NotNil(t, updated.GatewayTrxID)
	assert.Equal(t, "inv-test-001", *updated.GatewayTrxID)

	require.NotNil(t, updated.GatewayPaymentURL)
	assert.Equal(t, "https://checkout.xendit.co/web/inv-test-001", *updated.GatewayPaymentURL)

	require.NotNil(t, updated.GatewayResponse)
	assert.Contains(t, *updated.GatewayResponse, "inv-test-001")
}

// TestHandleGatewayWebhook_Confirmed verifies that a "confirmed" event confirms the payment.
func TestHandleGatewayWebhook_Confirmed(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	repos := postgres.NewRepository(suite.DB)
	paymentSvc := setupPaymentGatewayServices(t, suite)
	_, payment := createTestCustomerAndPayment(t, suite, repos, "002")

	// Pre-set gateway info so GatewayTrxID is set (required by domain for some checks)
	gatewayTrxID := "inv-gw-002"
	paymentID, err := uuid.Parse(payment.ID)
	require.NoError(t, err)
	require.NoError(t, paymentSvc.SetGatewayInfo(suite.Ctx, paymentID, "xendit", gatewayTrxID, "https://checkout.xendit.co/web/inv-gw-002", "{}"))

	event := &gateway.WebhookEvent{
		ExternalID: payment.ID,
		GatewayID:  gatewayTrxID,
		Status:     "confirmed",
	}

	err = paymentSvc.HandleGatewayWebhook(suite.Ctx, event)
	require.NoError(t, err)

	updated, err := repos.PaymentRepo.GetByID(suite.Ctx, paymentID)
	require.NoError(t, err)
	assert.Equal(t, "confirmed", updated.Status)
}

// TestHandleGatewayWebhook_Rejected verifies that a "rejected" event rejects the payment.
func TestHandleGatewayWebhook_Rejected(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	repos := postgres.NewRepository(suite.DB)
	paymentSvc := setupPaymentGatewayServices(t, suite)
	_, payment := createTestCustomerAndPayment(t, suite, repos, "003")

	event := &gateway.WebhookEvent{
		ExternalID: payment.ID,
		GatewayID:  "inv-gw-003",
		Status:     "rejected",
	}

	err := paymentSvc.HandleGatewayWebhook(suite.Ctx, event)
	require.NoError(t, err)

	paymentID, _ := uuid.Parse(payment.ID)
	updated, err := repos.PaymentRepo.GetByID(suite.Ctx, paymentID)
	require.NoError(t, err)
	assert.Equal(t, "rejected", updated.Status)
}

// TestHandleGatewayWebhook_InvalidExternalID verifies that an invalid UUID returns an error.
func TestHandleGatewayWebhook_InvalidExternalID(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	paymentSvc := setupPaymentGatewayServices(t, suite)

	event := &gateway.WebhookEvent{
		ExternalID: "not-a-uuid",
		GatewayID:  "inv-x",
		Status:     "confirmed",
	}

	err := paymentSvc.HandleGatewayWebhook(suite.Ctx, event)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid external_id")
}

// TestHandleGatewayWebhook_PendingStatus verifies that "pending" status is a no-op.
func TestHandleGatewayWebhook_PendingStatus(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	repos := postgres.NewRepository(suite.DB)
	paymentSvc := setupPaymentGatewayServices(t, suite)
	_, payment := createTestCustomerAndPayment(t, suite, repos, "004")

	event := &gateway.WebhookEvent{
		ExternalID: payment.ID,
		GatewayID:  "inv-gw-004",
		Status:     "pending",
	}

	err := paymentSvc.HandleGatewayWebhook(suite.Ctx, event)
	require.NoError(t, err)

	// Status should remain pending
	paymentID, _ := uuid.Parse(payment.ID)
	updated, err := repos.PaymentRepo.GetByID(suite.Ctx, paymentID)
	require.NoError(t, err)
	assert.Equal(t, "pending", updated.Status)
}

// TestHandleGatewayWebhook_UnknownStatus verifies unknown statuses return an error (so gateway retries).
func TestHandleGatewayWebhook_UnknownStatus(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	repos := postgres.NewRepository(suite.DB)
	paymentSvc := setupPaymentGatewayServices(t, suite)
	_, payment := createTestCustomerAndPayment(t, suite, repos, "005")

	event := &gateway.WebhookEvent{
		ExternalID: payment.ID,
		GatewayID:  "inv-gw-005",
		Status:     "some_future_status",
	}

	err := paymentSvc.HandleGatewayWebhook(suite.Ctx, event)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unhandled status")
}

// TestSetGatewayInfo_IdempotencyGuard verifies that calling SetGatewayInfo twice returns an error.
func TestSetGatewayInfo_IdempotencyGuard(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	repos := postgres.NewRepository(suite.DB)
	paymentSvc := setupPaymentGatewayServices(t, suite)
	_, payment := createTestCustomerAndPayment(t, suite, repos, "006")

	paymentID, _ := uuid.Parse(payment.ID)

	// First call — should succeed
	err := paymentSvc.SetGatewayInfo(suite.Ctx, paymentID, "xendit", "inv-006", "https://checkout.xendit.co/inv-006", "{}")
	require.NoError(t, err)

	// Second call — should fail because GatewayTrxID is already set
	err = paymentSvc.SetGatewayInfo(suite.Ctx, paymentID, "xendit", "inv-006-new", "https://checkout.xendit.co/inv-006-new", "{}")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "already has a gateway invoice")
}
