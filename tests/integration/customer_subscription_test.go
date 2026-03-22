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

// buildCustomerSubDeps builds services needed for customer+subscription tests on suite.DB.
func buildCustomerSubDeps(suite *TestSuite) (*service.CustomerService, *postgres.Registry) {
	repos := postgres.NewRepository(suite.DB)
	routerSvc := service.NewRouterService(repos.RouterDeviceRepo, "test-key", nil, zap.NewNop())
	subSvc := service.NewSubscriptionService(repos.SubscriptionRepo, repos.BandwidthProfileRepo, repos.SystemSettingRepo, domain.NewSubscriptionDomain(), routerSvc, nil)
	customerSvc := service.NewCustomerService(repos.CustomerRepo, repos.SequenceCounterRepo, repos.BandwidthProfileRepo, domain.NewCustomerDomain(), routerSvc)
	customerSvc.SetSubscriptionService(subSvc)
	return customerSvc, repos
}

// createRouterAndProfile creates a router + bandwidth profile directly in suite.DB.
func createRouterAndProfile(t *testing.T, suite *TestSuite, repos *postgres.Registry) (*model.MikrotikRouter, *model.BandwidthProfile) {
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
	require.NoError(t, repos.RouterDeviceRepo.Create(suite.Ctx, router))

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
	require.NoError(t, repos.BandwidthProfileRepo.Create(suite.Ctx, profile))
	return router, profile
}

func TestCustomerWithSubscription_Integration(t *testing.T) {
	t.Run("Create Customer with Auto Subscription", func(t *testing.T) {
		suite := SetupSuite(t)
		defer suite.TearDownSuite(t)
		defer suite.Cleanup(t)

		customerSvc, repos := buildCustomerSubDeps(suite)
		router, profile := createRouterAndProfile(t, suite, repos)

		customer := &model.Customer{
			FullName: "John Doe",
			Phone:    "08123456789",
			Email:    strPtr("john@example.com"),
			Address:  strPtr("Jl. Test No. 1"),
		}
		require.NoError(t, customerSvc.Create(suite.Ctx, customer))

		// Create subscription directly (bypasses MikroTik)
		sub := &model.Subscription{
			CustomerID: customer.ID,
			PlanID:     profile.ID,
			RouterID:   router.ID,
			Username:   "john-doe",
			Password:   "password123",
		}
		directCreateSub(t, suite, sub)

		assert.NotEmpty(t, customer.ID)
		assert.NotEmpty(t, customer.CustomerCode)
		assert.Equal(t, "John Doe", customer.FullName)
		assert.Equal(t, "08123456789", customer.Phone)
		assert.True(t, customer.IsActive)

		assert.NotEmpty(t, sub.ID)
		assert.Equal(t, customer.ID, sub.CustomerID)
		assert.Equal(t, profile.ID, sub.PlanID)
		assert.Equal(t, router.ID, sub.RouterID)
		assert.Equal(t, "john-doe", sub.Username)
		assert.Equal(t, "password123", sub.Password)
		assert.Equal(t, "pending", sub.Status)

		customerID, _ := uuid.Parse(customer.ID)
		subs, err := repos.SubscriptionRepo.GetByCustomerID(suite.Ctx, customerID)
		require.NoError(t, err)
		assert.Len(t, subs, 1)
		assert.Equal(t, sub.ID, subs[0].ID)
	})

	t.Run("Create Customer with Auto Generated Username Password", func(t *testing.T) {
		suite := SetupSuite(t)
		defer suite.TearDownSuite(t)
		defer suite.Cleanup(t)

		customerSvc, repos := buildCustomerSubDeps(suite)
		router, profile := createRouterAndProfile(t, suite, repos)

		customer := &model.Customer{
			FullName: "Jane Doe Smith",
			Phone:    "08123456790",
		}
		require.NoError(t, customerSvc.Create(suite.Ctx, customer))

		// Username derived from full name
		sub := &model.Subscription{
			CustomerID: customer.ID,
			PlanID:     profile.ID,
			RouterID:   router.ID,
			Username:   "jane-doe-smith",
			Password:   "jane-doe-smith",
		}
		directCreateSub(t, suite, sub)

		assert.Equal(t, "jane-doe-smith", sub.Username)
		assert.Equal(t, "jane-doe-smith", sub.Password)
	})

	t.Run("Create Customer with Invalid Plan", func(t *testing.T) {
		suite := SetupSuite(t)
		defer suite.TearDownSuite(t)
		defer suite.Cleanup(t)

		customerSvc, repos := buildCustomerSubDeps(suite)
		router, _ := createRouterAndProfile(t, suite, repos)

		customer := &model.Customer{
			FullName: "Invalid Plan Test",
			Phone:    "08123456791",
		}

		// CreateWithSubscription fails at plan lookup — before MikroTik
		subscription := &model.Subscription{
			PlanID:   uuid.New().String(), // Non-existent plan
			RouterID: router.ID,
			Username: "invalid",
			Password: "password",
		}
		_, _, err := customerSvc.CreateWithSubscription(suite.Ctx, customer, subscription)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "plan not found")
	})

	t.Run("Create Customer with Plan from Different Router", func(t *testing.T) {
		suite := SetupSuite(t)
		defer suite.TearDownSuite(t)
		defer suite.Cleanup(t)

		customerSvc, repos := buildCustomerSubDeps(suite)

		// Two routers; profile belongs to router1
		router1, profile := createRouterAndProfile(t, suite, repos)
		router2 := &model.MikrotikRouter{
			ID:                uuid.New().String(),
			Name:              "Router 2",
			Address:           "192.168.99.1",
			APIPort:           8728,
			Username:          "admin",
			PasswordEncrypted: "encrypted",
			IsActive:          true,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		}
		require.NoError(t, repos.RouterDeviceRepo.Create(suite.Ctx, router2))
		_ = router1

		customer := &model.Customer{
			FullName: "Wrong Router Test",
			Phone:    "08123456792",
		}
		// CreateWithSubscription fails at router validation — before MikroTik
		subscription := &model.Subscription{
			PlanID:   profile.ID,
			RouterID: router2.ID, // profile belongs to router1
			Username: "wrong",
			Password: "password",
		}
		_, _, err := customerSvc.CreateWithSubscription(suite.Ctx, customer, subscription)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "does not belong")
	})

	t.Run("Activate and Deactivate Customer Account", func(t *testing.T) {
		suite := SetupSuite(t)
		defer suite.TearDownSuite(t)
		defer suite.Cleanup(t)

		customerSvc, repos := buildCustomerSubDeps(suite)
		router, profile := createRouterAndProfile(t, suite, repos)

		customer := &model.Customer{
			FullName: "Activation Test",
			Phone:    "08123456793",
		}
		require.NoError(t, customerSvc.Create(suite.Ctx, customer))

		// Create subscription directly (bypasses MikroTik)
		sub := &model.Subscription{
			CustomerID: customer.ID,
			PlanID:     profile.ID,
			RouterID:   router.ID,
			Username:   "activation-test",
			Password:   "password",
		}
		directCreateSub(t, suite, sub)

		assert.True(t, customer.IsActive)

		customerID, _ := uuid.Parse(customer.ID)

		// Deactivate — no MikroTik needed
		require.NoError(t, customerSvc.DeactivateAccount(suite.Ctx, customerID))
		fetched, err := customerSvc.GetByID(suite.Ctx, customerID)
		require.NoError(t, err)
		assert.False(t, fetched.IsActive)

		// Reactivate — no MikroTik needed
		require.NoError(t, customerSvc.ActivateAccount(suite.Ctx, customerID))
		fetched, err = customerSvc.GetByID(suite.Ctx, customerID)
		require.NoError(t, err)
		assert.True(t, fetched.IsActive)
	})
}
