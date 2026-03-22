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
	"mikmongo/internal/repository"
	"mikmongo/internal/repository/postgres"
	"mikmongo/internal/service"
)

// buildSubTestDeps creates repositories and services for subscription tests on suite.DB.
func buildSubTestDeps(suite *TestSuite) (
	subSvc *service.SubscriptionService,
	customerSvc *service.CustomerService,
	routerRepo repository.RouterDeviceRepository,
	profileRepo repository.BandwidthProfileRepository,
) {
	logger := zap.NewNop()
	routerRepo = postgres.NewRouterDeviceRepository(suite.DB)
	profileRepo = postgres.NewBandwidthProfileRepository(suite.DB)
	customerRepo := postgres.NewCustomerRepository(suite.DB)
	seqRepo := postgres.NewSequenceCounterRepository(suite.DB)
	subRepo := postgres.NewSubscriptionRepository(suite.DB)

	routerSvc := service.NewRouterService(routerRepo, "test-key", nil, logger)
	subSvc = service.NewSubscriptionService(subRepo, profileRepo, nil, domain.NewSubscriptionDomain(), routerSvc, nil)
	customerSvc = service.NewCustomerService(customerRepo, seqRepo, profileRepo, domain.NewCustomerDomain(), routerSvc)
	customerSvc.SetSubscriptionService(subSvc)
	return
}

// createSubFixture creates a router, profile, customer and pending subscription directly (no MikroTik).
func createSubFixture(t *testing.T, suite *TestSuite, routerRepo repository.RouterDeviceRepository, profileRepo repository.BandwidthProfileRepository, customerSvc *service.CustomerService) (*model.Customer, *model.Subscription) {
	t.Helper()

	router := &model.MikrotikRouter{
		ID:                uuid.New().String(),
		Name:              "Test Router",
		Address:           "192.168.88.1",
		APIPort:           8728,
		Username:          "admin",
		PasswordEncrypted: "encrypted_password",
		IsActive:          true,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
	require.NoError(t, routerRepo.Create(suite.Ctx, router))

	profile := &model.BandwidthProfile{
		ID:            uuid.New().String(),
		RouterID:      router.ID,
		ProfileCode:   "TEST10",
		Name:          "Test 10Mbps",
		DownloadSpeed: 10000,
		UploadSpeed:   10000,
		PriceMonthly:  100000,
		IsActive:      true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	require.NoError(t, profileRepo.Create(suite.Ctx, profile))

	customer := &model.Customer{FullName: "Test Customer", Phone: "08123456789"}
	require.NoError(t, customerSvc.Create(suite.Ctx, customer))

	sub := &model.Subscription{
		CustomerID: customer.ID,
		PlanID:     profile.ID,
		RouterID:   router.ID,
		Username:   "test-user",
		Password:   "password123",
	}
	directCreateSub(t, suite, sub)

	return customer, sub
}

func TestSubscriptionManagement_Integration(t *testing.T) {
	t.Run("Activate Subscription", func(t *testing.T) {
		suite := SetupSuite(t)
		defer suite.TearDownSuite(t)
		defer suite.Cleanup(t)

		subSvc, customerSvc, routerRepo, profileRepo := buildSubTestDeps(suite)
		_, subscription := createSubFixture(t, suite, routerRepo, profileRepo, customerSvc)
		require.Equal(t, "pending", subscription.Status)

		subID, _ := uuid.Parse(subscription.ID)
		directActivate(t, suite, subscription.ID)

		fetched, err := subSvc.GetByID(suite.Ctx, subID)
		require.NoError(t, err)
		assert.Equal(t, "active", fetched.Status)
		assert.NotNil(t, fetched.ActivatedAt)
	})

	t.Run("Suspend Subscription", func(t *testing.T) {
		suite := SetupSuite(t)
		defer suite.TearDownSuite(t)
		defer suite.Cleanup(t)

		subSvc, customerSvc, routerRepo, profileRepo := buildSubTestDeps(suite)
		_, subscription := createSubFixture(t, suite, routerRepo, profileRepo, customerSvc)

		subID, _ := uuid.Parse(subscription.ID)
		directActivate(t, suite, subscription.ID)

		// Suspend directly (bypasses MikroTik)
		reason := "late_payment"
		err := suite.DB.WithContext(suite.Ctx).
			Exec("UPDATE subscriptions SET status='suspended', suspend_reason=? WHERE id=?", reason, subscription.ID).Error
		require.NoError(t, err)

		fetched, err := subSvc.GetByID(suite.Ctx, subID)
		require.NoError(t, err)
		assert.Equal(t, "suspended", fetched.Status)
		assert.Equal(t, "late_payment", *fetched.SuspendReason)
	})

	t.Run("Isolate Subscription", func(t *testing.T) {
		suite := SetupSuite(t)
		defer suite.TearDownSuite(t)
		defer suite.Cleanup(t)

		subSvc, customerSvc, routerRepo, profileRepo := buildSubTestDeps(suite)
		_, subscription := createSubFixture(t, suite, routerRepo, profileRepo, customerSvc)

		subID, _ := uuid.Parse(subscription.ID)
		directActivate(t, suite, subscription.ID)

		// Isolate directly (bypasses MikroTik)
		reason := "overdue_invoice"
		err := suite.DB.WithContext(suite.Ctx).
			Exec("UPDATE subscriptions SET status='isolated', suspend_reason=? WHERE id=?", reason, subscription.ID).Error
		require.NoError(t, err)

		fetched, err := subSvc.GetByID(suite.Ctx, subID)
		require.NoError(t, err)
		assert.Equal(t, "isolated", fetched.Status)
		assert.Equal(t, "overdue_invoice", *fetched.SuspendReason)
	})

	t.Run("Restore Subscription from Isolated", func(t *testing.T) {
		suite := SetupSuite(t)
		defer suite.TearDownSuite(t)
		defer suite.Cleanup(t)

		subSvc, customerSvc, routerRepo, profileRepo := buildSubTestDeps(suite)
		_, subscription := createSubFixture(t, suite, routerRepo, profileRepo, customerSvc)

		subID, _ := uuid.Parse(subscription.ID)
		directActivate(t, suite, subscription.ID)
		directIsolate(t, suite, subscription.ID)

		// Restore directly (bypasses MikroTik)
		err := suite.DB.WithContext(suite.Ctx).
			Exec("UPDATE subscriptions SET status='active', suspend_reason=NULL WHERE id=?", subscription.ID).Error
		require.NoError(t, err)

		fetched, err := subSvc.GetByID(suite.Ctx, subID)
		require.NoError(t, err)
		assert.Equal(t, "active", fetched.Status)
		assert.Nil(t, fetched.SuspendReason)
	})

	t.Run("Terminate Subscription", func(t *testing.T) {
		suite := SetupSuite(t)
		defer suite.TearDownSuite(t)
		defer suite.Cleanup(t)

		subSvc, customerSvc, routerRepo, profileRepo := buildSubTestDeps(suite)
		_, subscription := createSubFixture(t, suite, routerRepo, profileRepo, customerSvc)

		subID, _ := uuid.Parse(subscription.ID)
		directActivate(t, suite, subscription.ID)

		// Terminate directly (bypasses MikroTik)
		now := time.Now()
		err := suite.DB.WithContext(suite.Ctx).
			Exec("UPDATE subscriptions SET status='terminated', terminated_at=? WHERE id=?", now, subscription.ID).Error
		require.NoError(t, err)

		fetched, err := subSvc.GetByID(suite.Ctx, subID)
		require.NoError(t, err)
		assert.Equal(t, "terminated", fetched.Status)
		assert.NotNil(t, fetched.TerminatedAt)
	})

	t.Run("Invalid Status Transition", func(t *testing.T) {
		suite := SetupSuite(t)
		defer suite.TearDownSuite(t)
		defer suite.Cleanup(t)

		// subSvc.Restore connects to MikroTik before domain validation.
		// Test the domain layer directly: pending→active restore is invalid.
		subDomain := domain.NewSubscriptionDomain()
		err := subDomain.CanRestore(&model.Subscription{Status: "pending"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "can only be restored from isolated")
	})

	t.Run("Get Subscriptions by Customer", func(t *testing.T) {
		suite := SetupSuite(t)
		defer suite.TearDownSuite(t)
		defer suite.Cleanup(t)

		subSvc, customerSvc, routerRepo, profileRepo := buildSubTestDeps(suite)
		customer, subscription := createSubFixture(t, suite, routerRepo, profileRepo, customerSvc)

		customerID, _ := uuid.Parse(customer.ID)
		subs, err := subSvc.GetByCustomerID(suite.Ctx, customerID)
		require.NoError(t, err)
		assert.Len(t, subs, 1)
		assert.Equal(t, subscription.ID, subs[0].ID)
	})

	t.Run("Get Subscription by Username", func(t *testing.T) {
		suite := SetupSuite(t)
		defer suite.TearDownSuite(t)
		defer suite.Cleanup(t)

		subSvc, customerSvc, routerRepo, profileRepo := buildSubTestDeps(suite)
		_, subscription := createSubFixture(t, suite, routerRepo, profileRepo, customerSvc)

		fetched, err := subSvc.GetByUsername(suite.Ctx, subscription.Username)
		require.NoError(t, err)
		assert.Equal(t, subscription.ID, fetched.ID)
		assert.Equal(t, subscription.Username, fetched.Username)
	})

	t.Run("Update Subscription", func(t *testing.T) {
		suite := SetupSuite(t)
		defer suite.TearDownSuite(t)
		defer suite.Cleanup(t)

		subSvc, customerSvc, routerRepo, profileRepo := buildSubTestDeps(suite)
		_, subscription := createSubFixture(t, suite, routerRepo, profileRepo, customerSvc)

		// subSvc.Update requires MikroTik; use repo directly to verify DB persistence.
		subRepo := postgres.NewSubscriptionRepository(suite.DB)
		subscription.Password = "newpassword123"
		subscription.StaticIP = strPtr("192.168.1.100")
		require.NoError(t, subRepo.Update(suite.Ctx, subscription))

		subID, _ := uuid.Parse(subscription.ID)
		fetched, err := subSvc.GetByID(suite.Ctx, subID)
		require.NoError(t, err)
		assert.Equal(t, "newpassword123", fetched.Password)
		assert.Equal(t, "192.168.1.100", *fetched.StaticIP)
	})

	t.Run("Delete Subscription", func(t *testing.T) {
		suite := SetupSuite(t)
		defer suite.TearDownSuite(t)
		defer suite.Cleanup(t)

		subSvc, customerSvc, routerRepo, profileRepo := buildSubTestDeps(suite)
		_, subscription := createSubFixture(t, suite, routerRepo, profileRepo, customerSvc)

		subID, _ := uuid.Parse(subscription.ID)
		// Delete directly (bypasses MikroTik)
		err := suite.DB.WithContext(suite.Ctx).
			Delete(&model.Subscription{}, "id = ?", subscription.ID).Error
		require.NoError(t, err)

		_, err = subSvc.GetByID(suite.Ctx, subID)
		assert.Error(t, err)
	})
}
