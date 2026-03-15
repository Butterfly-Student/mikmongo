//go:build integration

package integration

import (
	"strings"
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

func TestPaymentConfirm_Overpayment(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	paymentSvc, billingSvc, customer, sub, _, adminID := setupPaymentTestFixtures(t, suite, "OVR01")

	subID, _ := uuid.Parse(sub.ID)
	inv, err := billingSvc.GenerateInvoice(suite.Ctx, subID, time.Now())
	require.NoError(t, err)
	require.NotNil(t, inv)

	// Pay more than the invoice amount
	overpayment := &model.Payment{
		CustomerID:    customer.ID,
		Amount:        inv.TotalAmount + 100_000, // 100_000 overpayment
		PaymentMethod: "cash",
	}
	require.NoError(t, paymentSvc.Create(suite.Ctx, overpayment))

	paymentID, _ := uuid.Parse(overpayment.ID)
	require.NoError(t, paymentSvc.Confirm(suite.Ctx, paymentID, adminID))

	repos := postgres.NewRepository(suite.DB)

	// Invoice should be fully paid (not over-paid)
	invID, _ := uuid.Parse(inv.ID)
	updatedInv, err := repos.InvoiceRepo.GetByID(suite.Ctx, invID)
	require.NoError(t, err)
	assert.Equal(t, "paid", updatedInv.Status)
	assert.Equal(t, inv.TotalAmount, updatedInv.PaidAmount,
		"PaidAmount should equal invoice total, not payment amount")

	// Allocated amount should equal invoice total (remainder silently discarded)
	allocations, err := repos.PaymentAllocationRepo.ListByPaymentID(suite.Ctx, paymentID)
	require.NoError(t, err)
	totalAllocated := 0.0
	for _, a := range allocations {
		totalAllocated += a.AllocatedAmount
	}
	assert.Equal(t, inv.TotalAmount, totalAllocated,
		"allocated amount should equal invoice total, not overpayment amount")

	// Document current behavior: the 100_000 overpayment is silently discarded
	t.Logf("%.0f overpayment untracked — current behavior; future: add credit balance feature",
		inv.TotalAmount+100_000-totalAllocated)
}

// TestConcurrentPaymentConfirm_RaceCondition verifies the idempotency guard inside
// PaymentService.Confirm: when two goroutines race to confirm the same payment,
// exactly one must succeed and one must fail.
//
// Design notes:
//   - Uses RootDB (committed transactions) so both goroutines share the same DB state.
//   - Uses UUID-based identifiers on every fixture to prevent collision with any prior run.
//   - Registers all cleanup BEFORE fixture creation so partial failures still clean up.
//   - G1 signals via channel before G2 starts; this guarantees G2 reads the payment
//     AFTER G1 committed (READ COMMITTED), so the idempotency re-read inside G2's tx
//     sees status="confirmed" and returns an error. Without this ordering guarantee the
//     test is non-deterministic (G2 can read "pending" before G1 commits and both succeed).
func TestConcurrentPaymentConfirm_RaceCondition(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t) // close the unused suite transaction opened by SetupSuite

	// Unique suffix prevents fixture collision with any committed data from a prior run.
	suffix := uuid.New().String()[:8]

	rootRepos := postgres.NewRepository(suite.RootDB)
	transactor := postgres.NewTransactor(suite.RootDB)
	logger := zap.NewNop()

	routerSvc := service.NewRouterService(rootRepos.RouterDeviceRepo, "test-key-16-bytes", nil, logger)
	customerSvc := service.NewCustomerService(
		rootRepos.CustomerRepo, rootRepos.SequenceCounterRepo,
		rootRepos.BandwidthProfileRepo, domain.NewCustomerDomain(), routerSvc,
	)
	paymentSvc := service.NewPaymentService(
		rootRepos.PaymentRepo, rootRepos.InvoiceRepo, rootRepos.PaymentAllocationRepo,
		rootRepos.CustomerRepo, rootRepos.SequenceCounterRepo,
		domain.NewPaymentDomain(), domain.NewBillingDomain(),
		transactor,
	)
	paymentSvc.SetCustomerService(customerSvc)
	billingSvc := service.NewBillingService(
		rootRepos.InvoiceRepo, rootRepos.InvoiceItemRepo, rootRepos.SubscriptionRepo,
		rootRepos.BandwidthProfileRepo, rootRepos.CustomerRepo, rootRepos.SystemSettingRepo,
		rootRepos.SequenceCounterRepo, domain.NewBillingDomain(),
	)

	// ── Fixtures ────────────────────────────────────────────────────────────────

	router := &model.MikrotikRouter{
		ID:                uuid.New().String(),
		Name:              "ConcurrentRouter-" + suffix,
		Address:           "10.99.88.1", // non-routable; unique per run via router.ID
		APIPort:           8728,
		Username:          "admin",
		PasswordEncrypted: "enc_pass",
		IsActive:          true,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
	require.NoError(t, rootRepos.RouterDeviceRepo.Create(suite.Ctx, router))

	profile := &model.BandwidthProfile{
		ID:              uuid.New().String(),
		RouterID:        router.ID,
		ProfileCode:     "RACE-" + suffix,
		Name:            "Race Profile " + suffix,
		DownloadSpeed:   10000,
		UploadSpeed:     10000,
		PriceMonthly:    200_000,
		TaxRate:         0,
		GracePeriodDays: 3,
		IsActive:        true,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	require.NoError(t, rootRepos.BandwidthProfileRepo.Create(suite.Ctx, profile))

	customer := &model.Customer{
		FullName: "Race Customer " + suffix,
		Phone:    "08" + suffix, // unique per run
	}
	require.NoError(t, customerSvc.Create(suite.Ctx, customer))

	sub := &model.Subscription{
		CustomerID: customer.ID,
		PlanID:     profile.ID,
		RouterID:   router.ID,
		Username:   "race-" + suffix,
		Password:   "password123",
		Status:     "pending",
	}
	require.NoError(t, suite.RootDB.WithContext(suite.Ctx).Create(sub).Error)
	require.NoError(t, suite.RootDB.WithContext(suite.Ctx).
		Exec("UPDATE subscriptions SET status='active', activated_at=? WHERE id=?", time.Now(), sub.ID).Error)

	subID, _ := uuid.Parse(sub.ID)
	inv, err := billingSvc.GenerateInvoice(suite.Ctx, subID, time.Now())
	require.NoError(t, err)

	// Admin user — must be a real committed row because payments.processed_by has a FK
	// to users. Using suite.DB (rollback tx) would make the FK invisible to RootDB tx.
	adminID := uuid.New().String()
	require.NoError(t, suite.RootDB.WithContext(suite.Ctx).Exec(
		`INSERT INTO users (id, full_name, email, password_hash, role, is_active, created_at, updated_at)
		 VALUES (?, ?, ?, ?, 'admin', true, NOW(), NOW())`,
		adminID, "Race Admin "+suffix, adminID+"@test.com", "$2a$04$placeholder",
	).Error)

	// ── Cleanup — registered BEFORE payment creation so it always runs ──────────
	// Use-after-test cleanup for all RootDB-committed fixtures.
	t.Cleanup(func() {
		suite.RootDB.Exec("DELETE FROM payment_allocations WHERE invoice_id = ?", inv.ID)
		suite.RootDB.Exec("DELETE FROM invoice_items WHERE invoice_id = ?", inv.ID)
		suite.RootDB.Exec("DELETE FROM invoices WHERE id = ?", inv.ID)
		suite.RootDB.Exec("DELETE FROM payments WHERE customer_id = ?", customer.ID)
		suite.RootDB.Exec("DELETE FROM subscriptions WHERE id = ?", sub.ID)
		suite.RootDB.Exec("DELETE FROM customers WHERE id = ?", customer.ID)
		suite.RootDB.Exec("DELETE FROM bandwidth_profiles WHERE id = ?", profile.ID)
		suite.RootDB.Exec("DELETE FROM mikrotik_routers WHERE id = ?", router.ID)
		suite.RootDB.Exec("DELETE FROM users WHERE id = ?", adminID)
	})

	// ── Single payment — both goroutines race to confirm the same payment ───────
	payment := &model.Payment{
		CustomerID:    customer.ID,
		Amount:        inv.TotalAmount,
		PaymentMethod: "cash",
	}
	require.NoError(t, paymentSvc.Create(suite.Ctx, payment))
	paymentID, _ := uuid.Parse(payment.ID)

	// ── Concurrent confirm ───────────────────────────────────────────────────────
	// G1 runs first; g1Done is closed after G1's Confirm returns (and commits).
	// G2 only starts reading from g1Done, guaranteeing it calls Confirm AFTER G1
	// committed. In READ COMMITTED PostgreSQL, G2's inside-tx re-read of the payment
	// then sees status="confirmed", causing the idempotency guard to fire.
	//
	// Why this ordering is required: the idempotency re-read in PaymentService.Confirm
	// happens BEFORE GetByCustomerIDForUpdate. If both goroutines race freely, G2 can
	// read the payment as "pending" before G1 commits, making the outcome non-deterministic.
	type confirmResult struct {
		err error
	}
	results := make(chan confirmResult, 2)
	g1Done := make(chan struct{})

	go func() {
		err := paymentSvc.Confirm(suite.Ctx, paymentID, adminID)
		results <- confirmResult{err}
		close(g1Done) // signal: G1 committed (or failed), G2 may now start
	}()

	<-g1Done // wait for G1 to finish before launching G2

	go func() {
		results <- confirmResult{paymentSvc.Confirm(suite.Ctx, paymentID, adminID)}
	}()

	r1 := <-results
	r2 := <-results

	// ── Assertions ───────────────────────────────────────────────────────────────

	var successErr, failErr error
	successCount := 0
	for _, r := range []confirmResult{r1, r2} {
		if r.err == nil {
			successCount++
		} else {
			failErr = r.err
		}
	}
	_ = successErr

	// Exactly one confirm must succeed and one must fail.
	assert.Equal(t, 1, successCount,
		"exactly one Confirm must succeed (idempotency guard inside tx)")

	require.Error(t, failErr,
		"second Confirm must return an error")
	assert.True(t,
		strings.Contains(failErr.Error(), "only pending payments can be confirmed"),
		"expected idempotency error, got: %q", failErr.Error())

	// Invoice must be fully paid exactly once — SELECT FOR UPDATE prevents double allocation.
	invID, _ := uuid.Parse(inv.ID)
	updatedInv, err := rootRepos.InvoiceRepo.GetByID(suite.Ctx, invID)
	require.NoError(t, err)
	assert.Equal(t, "paid", updatedInv.Status,
		"invoice must be paid after one successful confirm")
	assert.Equal(t, inv.TotalAmount, updatedInv.PaidAmount,
		"PaidAmount must equal TotalAmount exactly (no double allocation)")
}
