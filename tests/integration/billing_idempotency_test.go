//go:build integration

package integration

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"mikmongo/internal/domain"
	"mikmongo/internal/model"
	"mikmongo/internal/repository/postgres"
	"mikmongo/internal/service"
)

func TestBillingIdempotency_GenerateInvoice_Twice(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	repos := postgres.NewRepository(suite.DB)
	billingSvc := service.NewBillingService(
		repos.InvoiceRepo, repos.InvoiceItemRepo, repos.SubscriptionRepo,
		repos.BandwidthProfileRepo, repos.CustomerRepo, repos.SystemSettingRepo,
		repos.SequenceCounterRepo, domain.NewBillingDomain(),
	)

	customer, sub, _ := billingTestSetup(t, suite, repos,
		"Idempotency Test 1", "088001000001", "IDEM01", 200_000, 0, nil)

	// Activate at start of month to avoid proration
	startOfMonth := time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.Local)
	err := suite.DB.WithContext(suite.Ctx).
		Exec("UPDATE subscriptions SET status='active', activated_at=? WHERE id=?", startOfMonth, sub.ID).Error
	require.NoError(t, err)

	subID, _ := uuid.Parse(sub.ID)
	now := time.Now()

	inv1, err := billingSvc.GenerateInvoice(suite.Ctx, subID, now)
	require.NoError(t, err)
	require.NotNil(t, inv1)

	inv2, err := billingSvc.GenerateInvoice(suite.Ctx, subID, now)
	require.NoError(t, err)
	require.NotNil(t, inv2)

	assert.Equal(t, inv1.ID, inv2.ID, "second call should return same invoice")

	customerID, _ := uuid.Parse(customer.ID)
	invoices, err := repos.InvoiceRepo.GetByCustomerID(suite.Ctx, customerID)
	require.NoError(t, err)
	assert.Len(t, invoices, 1, "only one invoice should exist for the period")
}

func TestBillingIdempotency_ProcessDailyBilling_Twice(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	repos := postgres.NewRepository(suite.DB)
	billingSvc := service.NewBillingService(
		repos.InvoiceRepo, repos.InvoiceItemRepo, repos.SubscriptionRepo,
		repos.BandwidthProfileRepo, repos.CustomerRepo, repos.SystemSettingRepo,
		repos.SequenceCounterRepo, domain.NewBillingDomain(),
	)

	today := time.Now().Day()

	type fixture struct {
		customer *model.Customer
		sub      *model.Subscription
	}
	fixtures := make([]fixture, 3)

	for i := 0; i < 3; i++ {
		customer, sub, _ := billingTestSetup(t, suite, repos,
			"Idem Customer "+string(rune('A'+i)),
			"08800200000"+string(rune('1'+i)),
			"IDEM2"+string(rune('A'+i)),
			200_000, 0, &today)
		directActivate(t, suite, sub.ID)
		fixtures[i] = fixture{customer: customer, sub: sub}
	}

	err := billingSvc.ProcessDailyBilling(suite.Ctx)
	require.NoError(t, err)

	err = billingSvc.ProcessDailyBilling(suite.Ctx)
	require.NoError(t, err)

	for _, f := range fixtures {
		customerID, _ := uuid.Parse(f.customer.ID)
		invoices, err := repos.InvoiceRepo.GetByCustomerID(suite.Ctx, customerID)
		require.NoError(t, err)
		assert.Len(t, invoices, 1, "customer %s should have exactly 1 invoice", f.customer.FullName)
	}
}

func TestBillingIdempotency_TriggerMonthlyAPI_Twice(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	repos := postgres.NewRepository(suite.DB)

	today := time.Now().Day()
	customer, sub, _ := billingTestSetup(t, suite, repos,
		"Idem API Customer", "088003000001", "IDEM3A", 200_000, 0, &today)
	directActivate(t, suite, sub.ID)

	r := buildTestRouter(t, suite)
	email, password, _ := createAPIUser(t, suite, "admin")
	adminToken := loginAs(t, r, email, password)

	w1 := makeRequest(t, r, http.MethodPost, "/api/v1/invoices/trigger-monthly", adminToken, nil)
	assert.Equal(t, http.StatusOK, w1.Code)

	w2 := makeRequest(t, r, http.MethodPost, "/api/v1/invoices/trigger-monthly", adminToken, nil)
	assert.Equal(t, http.StatusOK, w2.Code)

	customerID, _ := uuid.Parse(customer.ID)
	invoices, err := repos.InvoiceRepo.GetByCustomerID(suite.Ctx, customerID)
	require.NoError(t, err)
	assert.Len(t, invoices, 1, "trigger-monthly called twice should create only 1 invoice")
}

func TestBillingIdempotency_CheckAndIsolateOverdue_Twice(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	repos := postgres.NewRepository(suite.DB)
	billingSvc := service.NewBillingService(
		repos.InvoiceRepo, repos.InvoiceItemRepo, repos.SubscriptionRepo,
		repos.BandwidthProfileRepo, repos.CustomerRepo, repos.SystemSettingRepo,
		repos.SequenceCounterRepo, domain.NewBillingDomain(),
	)

	_, sub, _ := billingTestSetup(t, suite, repos,
		"Idem Overdue Customer", "088004000001", "IDEM4A", 200_000, 0, nil)

	directActivate(t, suite, sub.ID)
	directIsolate(t, suite, sub.ID)

	graceDays := 3
	dueDate := time.Now().AddDate(0, 0, -(graceDays + 5))
	subIDStr := sub.ID
	invoice := &model.Invoice{
		InvoiceNumber:      "INV-IDEM-OVERDUE-001",
		CustomerID:         sub.CustomerID,
		SubscriptionID:     &subIDStr,
		BillingPeriodStart: dueDate.AddDate(0, -1, 0),
		BillingPeriodEnd:   dueDate,
		IssueDate:          dueDate.AddDate(0, 0, -30),
		DueDate:            dueDate,
		TotalAmount:        200_000,
		Status:             "unpaid",
		InvoiceType:        "recurring",
		IsAutoGenerated:    true,
	}
	require.NoError(t, repos.InvoiceRepo.Create(suite.Ctx, invoice))

	err := billingSvc.CheckAndIsolateOverdue(suite.Ctx)
	require.NoError(t, err)

	err = billingSvc.CheckAndIsolateOverdue(suite.Ctx)
	require.NoError(t, err)

	subID, _ := uuid.Parse(sub.ID)
	updatedSub, err := repos.SubscriptionRepo.GetByID(suite.Ctx, subID)
	require.NoError(t, err)
	assert.Equal(t, "isolated", updatedSub.Status, "subscription should remain isolated, not change state")
}

func TestBillingIdempotency_GenerateInvoice_DifferentPeriods(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	repos := postgres.NewRepository(suite.DB)
	billingSvc := service.NewBillingService(
		repos.InvoiceRepo, repos.InvoiceItemRepo, repos.SubscriptionRepo,
		repos.BandwidthProfileRepo, repos.CustomerRepo, repos.SystemSettingRepo,
		repos.SequenceCounterRepo, domain.NewBillingDomain(),
	)

	customer, sub, _ := billingTestSetup(t, suite, repos,
		"Idem DiffPeriod Customer", "088005000001", "IDEM5A", 200_000, 0, nil)

	// Activate 2 months ago so neither thisMonth nor lastMonth triggers proration
	twoMonthsAgo := time.Date(time.Now().Year(), time.Now().Month()-2, 1, 0, 0, 0, 0, time.Local)
	err := suite.DB.WithContext(suite.Ctx).
		Exec("UPDATE subscriptions SET status='active', activated_at=? WHERE id=?", twoMonthsAgo, sub.ID).Error
	require.NoError(t, err)

	subID, _ := uuid.Parse(sub.ID)
	now := time.Now()
	thisMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
	lastMonth := thisMonth.AddDate(0, -1, 0)

	inv1, err := billingSvc.GenerateInvoice(suite.Ctx, subID, thisMonth)
	require.NoError(t, err)
	require.NotNil(t, inv1)

	inv2, err := billingSvc.GenerateInvoice(suite.Ctx, subID, lastMonth)
	require.NoError(t, err)
	require.NotNil(t, inv2)

	assert.NotEqual(t, inv1.ID, inv2.ID, "different periods should create different invoices")

	customerID, _ := uuid.Parse(customer.ID)
	invoices, err := repos.InvoiceRepo.GetByCustomerID(suite.Ctx, customerID)
	require.NoError(t, err)
	assert.Len(t, invoices, 2, "two invoices should exist for two different periods")
}
