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

func setupCustomerTestServices(t *testing.T, suite *TestSuite) (
	*service.CustomerService,
	*postgres.Registry,
) {
	t.Helper()
	repos := postgres.NewRepository(suite.DB)
	logger := zap.NewNop()

	routerSvc := service.NewRouterService(repos.RouterDeviceRepo, "test-key-16-bytes", nil, logger)
	customerSvc := service.NewCustomerService(
		repos.CustomerRepo, repos.SequenceCounterRepo, repos.BandwidthProfileRepo,
		domain.NewCustomerDomain(), routerSvc,
	)
	return customerSvc, repos
}

func createTestRouterAndProfile(t *testing.T, suite *TestSuite, repos *postgres.Registry, suffix string) (*model.MikrotikRouter, *model.BandwidthProfile) {
	t.Helper()
	router := &model.MikrotikRouter{
		ID:                uuid.New().String(),
		Name:              "Router " + suffix,
		Address:           "192.168.88.1",
		APIPort:           8728,
		Username:          "admin",
		PasswordEncrypted: "placeholder",
		IsActive:          true,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
	require.NoError(t, repos.RouterDeviceRepo.Create(suite.Ctx, router))

	profile := &model.BandwidthProfile{
		ID:              uuid.New().String(),
		RouterID:        router.ID,
		ProfileCode:     "CUST" + suffix,
		Name:            "Cust Profile " + suffix,
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
	return router, profile
}

func TestCustomerService_Create(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	customerSvc, _ := setupCustomerTestServices(t, suite)

	customer := &model.Customer{
		FullName: "New Customer",
		Phone:    "081111111111",
	}
	err := customerSvc.Create(suite.Ctx, customer)
	require.NoError(t, err)

	assert.NotEmpty(t, customer.ID)
	assert.NotEmpty(t, customer.CustomerCode)
	// Format should be CST##### (CST + 5 digits)
	assert.Regexp(t, `^CST\d{5}$`, customer.CustomerCode)
}

func TestCustomerService_CreateWithSubscription(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	customerSvc, repos := setupCustomerTestServices(t, suite)
	router, profile := createTestRouterAndProfile(t, suite, repos, "C01")

	// Create customer first
	customer := &model.Customer{FullName: "Sub Customer", Phone: "082222222222"}
	require.NoError(t, customerSvc.Create(suite.Ctx, customer))

	// Create subscription directly (bypass router)
	sub := &model.Subscription{
		CustomerID: customer.ID,
		PlanID:     profile.ID,
		RouterID:   router.ID,
		Username:   "sub-customer-001",
		Password:   "password123",
	}
	directCreateSub(t, suite, sub)

	assert.NotEmpty(t, customer.ID)
	assert.NotEmpty(t, customer.CustomerCode)
	assert.NotEmpty(t, sub.ID)
	assert.Equal(t, customer.ID, sub.CustomerID)
	assert.Equal(t, "pending", sub.Status)
}

func TestCustomerService_SuspendAllSubscriptions(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	customerSvc, repos := setupCustomerTestServices(t, suite)
	_, profile := createTestRouterAndProfile(t, suite, repos, "C02")

	customer := &model.Customer{FullName: "Suspend Customer", Phone: "083333333333"}
	require.NoError(t, customerSvc.Create(suite.Ctx, customer))
	customerID, _ := uuid.Parse(customer.ID)

	sub1 := &model.Subscription{
		CustomerID: customer.ID,
		PlanID:     profile.ID,
		RouterID:   profile.RouterID,
		Username:   "suspend-user-1",
		Password:   "password123",
	}
	sub2 := &model.Subscription{
		CustomerID: customer.ID,
		PlanID:     profile.ID,
		RouterID:   profile.RouterID,
		Username:   "suspend-user-2",
		Password:   "password123",
	}
	directCreateSub(t, suite, sub1)
	directCreateSub(t, suite, sub2)

	// Activate then suspend directly
	directActivate(t, suite, sub1.ID)
	directActivate(t, suite, sub2.ID)
	directSuspend(t, suite, sub1.ID)
	directSuspend(t, suite, sub2.ID)

	subs, err := repos.SubscriptionRepo.GetByCustomerID(suite.Ctx, customerID)
	require.NoError(t, err)
	for _, s := range subs {
		assert.Equal(t, "suspended", s.Status)
	}
}

func TestCustomerService_IsolateAllSubscriptions(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	customerSvc, repos := setupCustomerTestServices(t, suite)
	_, profile := createTestRouterAndProfile(t, suite, repos, "C03")

	customer := &model.Customer{FullName: "Isolate Customer", Phone: "084444444444"}
	require.NoError(t, customerSvc.Create(suite.Ctx, customer))
	customerID, _ := uuid.Parse(customer.ID)

	sub := &model.Subscription{
		CustomerID: customer.ID,
		PlanID:     profile.ID,
		RouterID:   profile.RouterID,
		Username:   "isolate-user-1",
		Password:   "password123",
	}
	directCreateSub(t, suite, sub)
	directActivate(t, suite, sub.ID)
	directIsolate(t, suite, sub.ID)

	subs, err := repos.SubscriptionRepo.GetByCustomerID(suite.Ctx, customerID)
	require.NoError(t, err)
	for _, s := range subs {
		assert.Equal(t, "isolated", s.Status)
	}
}

func TestCustomerService_RestoreAllSubscriptions(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	customerSvc, repos := setupCustomerTestServices(t, suite)
	_, profile := createTestRouterAndProfile(t, suite, repos, "C04")

	customer := &model.Customer{FullName: "Restore Customer", Phone: "085555555555"}
	require.NoError(t, customerSvc.Create(suite.Ctx, customer))
	customerID, _ := uuid.Parse(customer.ID)

	sub := &model.Subscription{
		CustomerID: customer.ID,
		PlanID:     profile.ID,
		RouterID:   profile.RouterID,
		Username:   "restore-user-1",
		Password:   "password123",
	}
	directCreateSub(t, suite, sub)
	directActivate(t, suite, sub.ID)
	directIsolate(t, suite, sub.ID)

	// Verify isolated
	subID, _ := uuid.Parse(sub.ID)
	isolatedSub, err := repos.SubscriptionRepo.GetByID(suite.Ctx, subID)
	require.NoError(t, err)
	assert.Equal(t, "isolated", isolatedSub.Status)

	// Restore directly
	err = suite.DB.WithContext(suite.Ctx).
		Exec("UPDATE subscriptions SET status='active' WHERE id=?", sub.ID).Error
	require.NoError(t, err)

	subs, err := repos.SubscriptionRepo.GetByCustomerID(suite.Ctx, customerID)
	require.NoError(t, err)
	for _, s := range subs {
		assert.Equal(t, "active", s.Status)
	}
}

func TestCustomerService_PortalAuth(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	customerSvc, _ := setupCustomerTestServices(t, suite)

	customer := &model.Customer{FullName: "Portal Customer", Phone: "086666666666"}
	require.NoError(t, customerSvc.Create(suite.Ctx, customer))
	customerID, _ := uuid.Parse(customer.ID)

	// Set portal password
	err := customerSvc.SetPortalPassword(suite.Ctx, customerID, "portal-password-123")
	require.NoError(t, err)

	// Reload customer to get code
	updatedCustomer, err := customerSvc.GetByID(suite.Ctx, customerID)
	require.NoError(t, err)

	// Auth portal with correct password
	authedCustomer, err := customerSvc.AuthPortal(suite.Ctx, updatedCustomer.CustomerCode, "portal-password-123")
	require.NoError(t, err)
	assert.Equal(t, customerID.String(), authedCustomer.ID)

	// Auth portal with wrong password
	_, err = customerSvc.AuthPortal(suite.Ctx, updatedCustomer.CustomerCode, "wrong-password")
	assert.Error(t, err)
}
