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

func TestCustomerWithSubscription_Integration(t *testing.T) {
	// Setup test suite
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)

	// Create repositories
	customerRepo := postgres.NewCustomerRepository(suite.DB)
	seqRepo := postgres.NewSequenceCounterRepository(suite.DB)
	profileRepo := postgres.NewBandwidthProfileRepository(suite.DB)
	routerRepo := postgres.NewRouterDeviceRepository(suite.DB)
	subRepo := postgres.NewSubscriptionRepository(suite.DB)

	// Create domains
	customerDomain := domain.NewCustomerDomain()
	subDomain := domain.NewSubscriptionDomain()

	// Create logger
	logger := zap.NewNop()

	// Create router service
	routerSvc := service.NewRouterService(routerRepo, "test-key", nil, logger)

	// Create subscription service
	subSvc := service.NewSubscriptionService(
		subRepo,
		profileRepo,
		nil, // settingRepo not needed for this test
		subDomain,
		routerSvc,
	)

	// Create customer service
	customerSvc := service.NewCustomerService(
		customerRepo,
		seqRepo,
		profileRepo,
		customerDomain,
		routerSvc,
	)
	customerSvc.SetSubscriptionService(subSvc)

	// Helper function to create test router
	createTestRouter := func(t *testing.T) *model.MikrotikRouter {
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
		err := routerRepo.Create(suite.Ctx, router)
		require.NoError(t, err)
		return router
	}

	// Helper function to create test bandwidth profile
	createTestProfile := func(t *testing.T, routerID string) *model.BandwidthProfile {
		profile := &model.BandwidthProfile{
			ID:            uuid.New().String(),
			RouterID:      routerID,
			ProfileCode:   "TEST10",
			Name:          "Test 10Mbps",
			DownloadSpeed: 10000,
			UploadSpeed:   10000,
			PriceMonthly:  100000,
			IsActive:      true,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
		err := profileRepo.Create(suite.Ctx, profile)
		require.NoError(t, err)
		return profile
	}

	t.Run("Create Customer with Auto Subscription", func(t *testing.T) {
		defer suite.Cleanup(t)

		// Create router and profile
		router := createTestRouter(t)
		profile := createTestProfile(t, router.ID)

		// Create customer
		customer := &model.Customer{
			FullName: "John Doe",
			Phone:    "08123456789",
			Email:    strPtr("john@example.com"),
			Address:  strPtr("Jl. Test No. 1"),
		}

		// Create subscription data
		subscription := &model.Subscription{
			PlanID:   profile.ID,
			RouterID: router.ID,
			Username: "john-doe",
			Password: "password123",
		}

		// Create both customer and subscription
		createdCustomer, createdSubscription, err := customerSvc.CreateWithSubscription(
			suite.Ctx,
			customer,
			subscription,
		)
		require.NoError(t, err)

		// Verify customer created
		assert.NotEmpty(t, createdCustomer.ID)
		assert.NotEmpty(t, createdCustomer.CustomerCode)
		assert.Equal(t, "John Doe", createdCustomer.FullName)
		assert.Equal(t, "08123456789", createdCustomer.Phone)
		assert.True(t, createdCustomer.IsActive)

		// Verify subscription created
		assert.NotEmpty(t, createdSubscription.ID)
		assert.Equal(t, createdCustomer.ID, createdSubscription.CustomerID)
		assert.Equal(t, profile.ID, createdSubscription.PlanID)
		assert.Equal(t, router.ID, createdSubscription.RouterID)
		assert.Equal(t, "john-doe", createdSubscription.Username)
		assert.Equal(t, "password123", createdSubscription.Password)
		assert.Equal(t, "pending", createdSubscription.Status)

		// Verify customer has subscription
		customerID, _ := uuid.Parse(createdCustomer.ID)
		subs, err := subSvc.GetByCustomerID(suite.Ctx, customerID)
		require.NoError(t, err)
		assert.Len(t, subs, 1)
		assert.Equal(t, createdSubscription.ID, subs[0].ID)
	})

	t.Run("Create Customer with Auto Generated Username Password", func(t *testing.T) {
		defer suite.Cleanup(t)

		// Create router and profile
		router := createTestRouter(t)
		profile := createTestProfile(t, router.ID)

		// Create customer with FullName containing spaces
		customer := &model.Customer{
			FullName: "Jane Doe Smith",
			Phone:    "08123456790",
		}

		// Create subscription without username/password (will be auto-generated)
		subscription := &model.Subscription{
			PlanID:   profile.ID,
			RouterID: router.ID,
			Username: "jane-doe-smith", // Auto-generated from FullName
			Password: "jane-doe-smith", // Auto-generated from FullName
		}

		// Create both
		_, createdSubscription, err := customerSvc.CreateWithSubscription(
			suite.Ctx,
			customer,
			subscription,
		)
		require.NoError(t, err)

		// Verify username/password generated correctly
		assert.Equal(t, "jane-doe-smith", createdSubscription.Username)
		assert.Equal(t, "jane-doe-smith", createdSubscription.Password)
	})

	t.Run("Create Customer with Invalid Plan", func(t *testing.T) {
		defer suite.Cleanup(t)

		// Create router but no profile
		router := createTestRouter(t)

		customer := &model.Customer{
			FullName: "Invalid Plan Test",
			Phone:    "08123456791",
		}

		subscription := &model.Subscription{
			PlanID:   uuid.New().String(), // Non-existent plan
			RouterID: router.ID,
			Username: "invalid",
			Password: "password",
		}

		// Should fail because plan doesn't exist
		_, _, err := customerSvc.CreateWithSubscription(
			suite.Ctx,
			customer,
			subscription,
		)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "plan not found")
	})

	t.Run("Create Customer with Plan from Different Router", func(t *testing.T) {
		defer suite.Cleanup(t)

		// Create two routers
		router1 := createTestRouter(t)
		router2 := createTestRouter(t)

		// Create profile for router1
		profile := createTestProfile(t, router1.ID)

		customer := &model.Customer{
			FullName: "Wrong Router Test",
			Phone:    "08123456792",
		}

		// Try to use profile with router2
		subscription := &model.Subscription{
			PlanID:   profile.ID,
			RouterID: router2.ID, // Different router
			Username: "wrong",
			Password: "password",
		}

		// Should fail
		_, _, err := customerSvc.CreateWithSubscription(
			suite.Ctx,
			customer,
			subscription,
		)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "does not belong")
	})

	t.Run("Activate and Deactivate Customer Account", func(t *testing.T) {
		defer suite.Cleanup(t)

		// Create router and profile
		router := createTestRouter(t)
		profile := createTestProfile(t, router.ID)

		customer := &model.Customer{
			FullName: "Activation Test",
			Phone:    "08123456793",
		}

		subscription := &model.Subscription{
			PlanID:   profile.ID,
			RouterID: router.ID,
			Username: "activation-test",
			Password: "password",
		}

		createdCustomer, _, err := customerSvc.CreateWithSubscription(
			suite.Ctx,
			customer,
			subscription,
		)
		require.NoError(t, err)

		// Verify initially active
		assert.True(t, createdCustomer.IsActive)

		// Deactivate account
		customerID, _ := uuid.Parse(createdCustomer.ID)
		err = customerSvc.DeactivateAccount(suite.Ctx, customerID)
		require.NoError(t, err)

		// Verify deactivated
		fetched, err := customerSvc.GetByID(suite.Ctx, customerID)
		require.NoError(t, err)
		assert.False(t, fetched.IsActive)

		// Reactivate account
		err = customerSvc.ActivateAccount(suite.Ctx, customerID)
		require.NoError(t, err)

		// Verify reactivated
		fetched, err = customerSvc.GetByID(suite.Ctx, customerID)
		require.NoError(t, err)
		assert.True(t, fetched.IsActive)
	})
}
