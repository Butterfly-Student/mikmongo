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

func TestSubscriptionManagement_Integration(t *testing.T) {
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

	// Helper function to create customer with subscription
	createCustomerWithSubscription := func(t *testing.T) (*model.Customer, *model.Subscription) {
		router := createTestRouter(t)
		profile := createTestProfile(t, router.ID)

		customer := &model.Customer{
			FullName: "Test Customer",
			Phone:    "08123456789",
		}

		subscription := &model.Subscription{
			PlanID:   profile.ID,
			RouterID: router.ID,
			Username: "test-user",
			Password: "password123",
		}

		createdCustomer, createdSubscription, err := customerSvc.CreateWithSubscription(
			suite.Ctx,
			customer,
			subscription,
		)
		require.NoError(t, err)

		return createdCustomer, createdSubscription
	}

	t.Run("Activate Subscription", func(t *testing.T) {
		defer suite.Cleanup(t)

		_, subscription := createCustomerWithSubscription(t)
		require.Equal(t, "pending", subscription.Status)

		// Activate subscription
		subID, _ := uuid.Parse(subscription.ID)
		err := subSvc.Activate(suite.Ctx, subID)
		require.NoError(t, err)

		// Verify activated
		fetched, err := subSvc.GetByID(suite.Ctx, subID)
		require.NoError(t, err)
		assert.Equal(t, "active", fetched.Status)
		assert.NotNil(t, fetched.ActivatedAt)
	})

	t.Run("Suspend Subscription", func(t *testing.T) {
		defer suite.Cleanup(t)

		_, subscription := createCustomerWithSubscription(t)

		// Activate first
		subID, _ := uuid.Parse(subscription.ID)
		err := subSvc.Activate(suite.Ctx, subID)
		require.NoError(t, err)

		// Suspend
		err = subSvc.Suspend(suite.Ctx, subID, "late_payment")
		require.NoError(t, err)

		// Verify suspended
		fetched, err := subSvc.GetByID(suite.Ctx, subID)
		require.NoError(t, err)
		assert.Equal(t, "suspended", fetched.Status)
		assert.Equal(t, "late_payment", *fetched.SuspendReason)
	})

	t.Run("Isolate Subscription", func(t *testing.T) {
		defer suite.Cleanup(t)

		_, subscription := createCustomerWithSubscription(t)

		// Activate first
		subID, _ := uuid.Parse(subscription.ID)
		err := subSvc.Activate(suite.Ctx, subID)
		require.NoError(t, err)

		// Isolate
		err = subSvc.Isolate(suite.Ctx, subID, "overdue_invoice")
		require.NoError(t, err)

		// Verify isolated
		fetched, err := subSvc.GetByID(suite.Ctx, subID)
		require.NoError(t, err)
		assert.Equal(t, "isolated", fetched.Status)
		assert.Equal(t, "overdue_invoice", *fetched.SuspendReason)
	})

	t.Run("Restore Subscription from Isolated", func(t *testing.T) {
		defer suite.Cleanup(t)

		_, subscription := createCustomerWithSubscription(t)

		// Activate -> Isolate -> Restore
		subID, _ := uuid.Parse(subscription.ID)
		err := subSvc.Activate(suite.Ctx, subID)
		require.NoError(t, err)

		err = subSvc.Isolate(suite.Ctx, subID, "test")
		require.NoError(t, err)

		err = subSvc.Restore(suite.Ctx, subID)
		require.NoError(t, err)

		// Verify restored to active
		fetched, err := subSvc.GetByID(suite.Ctx, subID)
		require.NoError(t, err)
		assert.Equal(t, "active", fetched.Status)
		assert.Nil(t, fetched.SuspendReason)
	})

	t.Run("Terminate Subscription", func(t *testing.T) {
		defer suite.Cleanup(t)

		_, subscription := createCustomerWithSubscription(t)

		// Activate first
		subID, _ := uuid.Parse(subscription.ID)
		err := subSvc.Activate(suite.Ctx, subID)
		require.NoError(t, err)

		// Terminate
		err = subSvc.Terminate(suite.Ctx, subID)
		require.NoError(t, err)

		// Verify terminated
		fetched, err := subSvc.GetByID(suite.Ctx, subID)
		require.NoError(t, err)
		assert.Equal(t, "terminated", fetched.Status)
		assert.NotNil(t, fetched.TerminatedAt)
	})

	t.Run("Invalid Status Transition", func(t *testing.T) {
		defer suite.Cleanup(t)

		_, subscription := createCustomerWithSubscription(t)

		// Try to restore from pending (invalid)
		subID, _ := uuid.Parse(subscription.ID)
		err := subSvc.Restore(suite.Ctx, subID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot transition")
	})

	t.Run("Get Subscriptions by Customer", func(t *testing.T) {
		defer suite.Cleanup(t)

		customer, subscription := createCustomerWithSubscription(t)

		// Get by customer ID
		customerID, _ := uuid.Parse(customer.ID)
		subs, err := subSvc.GetByCustomerID(suite.Ctx, customerID)
		require.NoError(t, err)
		assert.Len(t, subs, 1)
		assert.Equal(t, subscription.ID, subs[0].ID)
	})

	t.Run("Get Subscription by Username", func(t *testing.T) {
		defer suite.Cleanup(t)

		_, subscription := createCustomerWithSubscription(t)

		// Get by username
		fetched, err := subSvc.GetByUsername(suite.Ctx, subscription.Username)
		require.NoError(t, err)
		assert.Equal(t, subscription.ID, fetched.ID)
		assert.Equal(t, subscription.Username, fetched.Username)
	})

	t.Run("Update Subscription", func(t *testing.T) {
		defer suite.Cleanup(t)

		_, subscription := createCustomerWithSubscription(t)

		// Update subscription
		subscription.Password = "newpassword123"
		subscription.StaticIP = strPtr("192.168.1.100")

		err := subSvc.Update(suite.Ctx, subscription)
		require.NoError(t, err)

		// Verify update
		subID, _ := uuid.Parse(subscription.ID)
		fetched, err := subSvc.GetByID(suite.Ctx, subID)
		require.NoError(t, err)
		assert.Equal(t, "newpassword123", fetched.Password)
		assert.Equal(t, "192.168.1.100", *fetched.StaticIP)
	})

	t.Run("Delete Subscription", func(t *testing.T) {
		defer suite.Cleanup(t)

		_, subscription := createCustomerWithSubscription(t)

		// Delete
		subID, _ := uuid.Parse(subscription.ID)
		err := subSvc.Delete(suite.Ctx, subID)
		require.NoError(t, err)

		// Verify deletion
		_, err = subSvc.GetByID(suite.Ctx, subID)
		assert.Error(t, err)
	})
}
