//go:build integration && mikrotik_legacy

package integration

import (
	"fmt"
	"testing"
	"time"

	"mikmongo/internal/domain"
	"mikmongo/internal/model"
	"mikmongo/internal/repository/postgres"
	"mikmongo/internal/service"
	"mikmongo/pkg/mikrotik"
	"mikmongo/pkg/mikrotik/client"
	mkdomain "mikmongo/pkg/mikrotik/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// TestSubscriptionMikroTikIntegration tests subscription lifecycle with real MikroTik
func TestSubscriptionMikroTikIntegration(t *testing.T) {
	// Setup test suite
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)

	// Get MikroTik connection details from env
	mtHost := getEnv("TEST_MIKROTIK_HOST", "192.168.233.1")
	mtPort := getEnv("TEST_MIKROTIK_PORT", "8728")
	mtUser := getEnv("TEST_MIKROTIK_USER", "admin")
	mtPass := getEnv("TEST_MIKROTIK_PASS", "")

	if mtPass == "" {
		t.Skip("TEST_MIKROTIK_PASS not set, skipping MikroTik integration test")
	}

	// Create repositories
	customerRepo := postgres.NewCustomerRepository(suite.DB)
	seqRepo := postgres.NewSequenceCounterRepository(suite.DB)
	profileRepo := postgres.NewBandwidthProfileRepository(suite.DB)
	routerRepo := postgres.NewRouterDeviceRepository(suite.DB)
	subRepo := postgres.NewSubscriptionRepository(suite.DB)
	settingRepo := postgres.NewSystemSettingRepository(suite.DB)

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
		settingRepo,
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

	// Helper function to create test router with real MikroTik connection
	createTestRouter := func(t *testing.T) *model.MikrotikRouter {
		router := &model.MikrotikRouter{
			ID:                uuid.New().String(),
			Name:              "Integration Test Router",
			Address:           mtHost,
			APIPort:           parsePort(mtPort),
			Username:          mtUser,
			PasswordEncrypted: mtPass,
			IsActive:          true,
			Status:            "online",
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
			ProfileCode:   "TEST10MBPS",
			Name:          "Test-10Mbps",
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

	// Helper function to create PPP profile in MikroTik
	createPPPProfileInMikroTik := func(t *testing.T, mt *mikrotik.Client, profileName string) {
		// Check if profile exists
		profiles, err := mt.PPP.GetProfiles(suite.Ctx)
		require.NoError(t, err)

		exists := false
		for _, p := range profiles {
			if p.Name == profileName {
				exists = true
				break
			}
		}

		if !exists {
			// Create profile
			profile := &mkdomain.PPPProfile{
				Name: profileName,
			}
			err := mt.PPP.AddProfile(suite.Ctx, profile)
			require.NoError(t, err)
			t.Logf("✓ Created PPP profile '%s' in MikroTik", profileName)
		} else {
			t.Logf("✓ PPP profile '%s' already exists in MikroTik", profileName)
		}
	}

	// Helper function to create isolate profile
	createIsolateProfile := func(t *testing.T, mt *mikrotik.Client) {
		profiles, err := mt.PPP.GetProfiles(suite.Ctx)
		require.NoError(t, err)

		exists := false
		for _, p := range profiles {
			if p.Name == "isolate" {
				exists = true
				break
			}
		}

		if !exists {
			profile := &mkdomain.PPPProfile{
				Name: "isolate",
			}
			err := mt.PPP.AddProfile(suite.Ctx, profile)
			require.NoError(t, err)
			t.Logf("✓ Created isolate profile in MikroTik")
		}
	}

	// Helper function to verify PPP secret exists in MikroTik
	verifyPPPSecretExists := func(t *testing.T, mt *mikrotik.Client, username string) bool {
		secret, err := mt.PPP.GetSecretByName(suite.Ctx, username)
		if err != nil {
			return false
		}
		return secret != nil && secret.Name == username
	}

	// Helper function to verify PPP secret is enabled
	verifyPPPSecretEnabled := func(t *testing.T, mt *mikrotik.Client, username string) bool {
		secret, err := mt.PPP.GetSecretByName(suite.Ctx, username)
		if err != nil {
			return false
		}
		return secret != nil && !secret.Disabled
	}

	// Helper function to verify PPP secret profile
	verifyPPPSecretProfile := func(t *testing.T, mt *mikrotik.Client, username, expectedProfile string) bool {
		secret, err := mt.PPP.GetSecretByName(suite.Ctx, username)
		if err != nil {
			return false
		}
		return secret != nil && secret.Profile == expectedProfile
	}

	// Helper function to create direct MikroTik client for verification
	createDirectMikroTikClient := func(t *testing.T) *mikrotik.Client {
		cfg := client.Config{
			Host:     mtHost,
			Port:     parsePort(mtPort),
			Username: mtUser,
			Password: mtPass,
		}
		mt, err := mikrotik.NewClient(cfg)
		require.NoError(t, err)
		return mt
	}

	t.Run("Create Subscription - PPP Secret Created in MikroTik", func(t *testing.T) {
		defer suite.Cleanup(t)

		// Create router and profile
		router := createTestRouter(t)
		profile := createTestProfile(t, router.ID)

		// Create MikroTik client using direct connection (not from Manager to avoid Close issues)
		mt := createDirectMikroTikClient(t)

		// Ensure PPP profile exists in MikroTik
		createPPPProfileInMikroTik(t, mt, profile.Name)

		// Create customer
		customer := &model.Customer{
			FullName: "Test Customer MikroTik",
			Phone:    "081234567890",
		}

		// Create subscription
		subscription := &model.Subscription{
			PlanID:   profile.ID,
			RouterID: router.ID,
			Username: "testuser_mt_" + uuid.New().String()[:8],
			Password: "testpass123",
		}

		// Create customer with subscription
		createdCustomer, createdSub, err := customerSvc.CreateWithSubscription(
			suite.Ctx,
			customer,
			subscription,
		)
		require.NoError(t, err)
		require.NotNil(t, createdCustomer)
		require.NotNil(t, createdSub)

		t.Logf("✓ Customer created: %s", createdCustomer.CustomerCode)
		t.Logf("✓ Subscription created: %s (status: %s)", createdSub.ID, createdSub.Status)
		t.Logf("✓ Subscription username: %s", createdSub.Username)

		// List all PPP secrets in MikroTik for debugging
		secrets, err := mt.PPP.GetSecrets(suite.Ctx, "")
		if err != nil {
			t.Logf("Warning: failed to list PPP secrets: %v", err)
		} else {
			t.Logf("Total PPP secrets in MikroTik: %d", len(secrets))
			for _, s := range secrets {
				t.Logf("  - %s (Profile: %s, Disabled: %v)", s.Name, s.Profile, s.Disabled)
			}
		}

		// Verify PPP secret exists in MikroTik
		exists := verifyPPPSecretExists(t, mt, createdSub.Username)
		assert.True(t, exists, "PPP secret should exist in MikroTik")

		if exists {
			t.Logf("✓ PPP secret '%s' exists in MikroTik", createdSub.Username)

			// Verify secret details
			secret, _ := mt.PPP.GetSecretByName(suite.Ctx, createdSub.Username)
			assert.Equal(t, createdSub.Password, secret.Password)
			assert.Equal(t, profile.Name, secret.Profile)
			t.Logf("✓ PPP secret has correct password and profile")
		}
	})

	t.Run("Activate Subscription - PPP Secret Enabled", func(t *testing.T) {
		defer suite.Cleanup(t)

		// Setup
		router := createTestRouter(t)
		profile := createTestProfile(t, router.ID)

		mt, err := routerSvc.GetMikrotikClient(suite.Ctx, uuid.MustParse(router.ID))
		require.NoError(t, err)

		createPPPProfileInMikroTik(t, mt, profile.Name)

		// Create customer and subscription
		customer := &model.Customer{
			FullName: "Test Customer Activate",
			Phone:    "081234567891",
		}
		subscription := &model.Subscription{
			PlanID:   profile.ID,
			RouterID: router.ID,
			Username: "testuser_act_" + uuid.New().String()[:8],
			Password: "testpass123",
		}

		_, createdSub, err := customerSvc.CreateWithSubscription(suite.Ctx, customer, subscription)
		require.NoError(t, err)

		t.Logf("✓ Subscription created with status: %s", createdSub.Status)

		// Activate subscription
		subID := uuid.MustParse(createdSub.ID)
		err = subSvc.Activate(suite.Ctx, subID)
		require.NoError(t, err)

		t.Logf("✓ Subscription activated")

		// Verify PPP secret is enabled
		enabled := verifyPPPSecretEnabled(t, mt, createdSub.Username)
		assert.True(t, enabled, "PPP secret should be enabled in MikroTik")

		if enabled {
			t.Logf("✓ PPP secret '%s' is enabled in MikroTik", createdSub.Username)
		}

		// Verify subscription status in DB
		fetchedSub, err := subSvc.GetByID(suite.Ctx, subID)
		require.NoError(t, err)
		assert.Equal(t, "active", fetchedSub.Status)
		assert.NotNil(t, fetchedSub.ActivatedAt)
		t.Logf("✓ Subscription status in DB: %s", fetchedSub.Status)
	})

	t.Run("Suspend Subscription - PPP Secret Disabled", func(t *testing.T) {
		defer suite.Cleanup(t)

		// Setup
		router := createTestRouter(t)
		profile := createTestProfile(t, router.ID)

		mt, err := routerSvc.GetMikrotikClient(suite.Ctx, uuid.MustParse(router.ID))
		require.NoError(t, err)

		createPPPProfileInMikroTik(t, mt, profile.Name)

		// Create customer and subscription
		customer := &model.Customer{
			FullName: "Test Customer Suspend",
			Phone:    "081234567892",
		}
		subscription := &model.Subscription{
			PlanID:   profile.ID,
			RouterID: router.ID,
			Username: "testuser_sus_" + uuid.New().String()[:8],
			Password: "testpass123",
		}

		_, createdSub, err := customerSvc.CreateWithSubscription(suite.Ctx, customer, subscription)
		require.NoError(t, err)

		// Activate first
		subID := uuid.MustParse(createdSub.ID)
		err = subSvc.Activate(suite.Ctx, subID)
		require.NoError(t, err)

		t.Logf("✓ Subscription activated")

		// Verify enabled
		enabled := verifyPPPSecretEnabled(t, mt, createdSub.Username)
		require.True(t, enabled, "PPP secret should be enabled before suspend")
		t.Logf("✓ PPP secret is enabled")

		// Suspend subscription
		err = subSvc.Suspend(suite.Ctx, subID, "late_payment")
		require.NoError(t, err)

		t.Logf("✓ Subscription suspended")

		// Verify PPP secret is disabled
		enabled = verifyPPPSecretEnabled(t, mt, createdSub.Username)
		assert.False(t, enabled, "PPP secret should be disabled in MikroTik")

		if !enabled {
			t.Logf("✓ PPP secret '%s' is disabled in MikroTik", createdSub.Username)
		}

		// Verify subscription status in DB
		fetchedSub, err := subSvc.GetByID(suite.Ctx, subID)
		require.NoError(t, err)
		assert.Equal(t, "suspended", fetchedSub.Status)
		assert.Equal(t, "late_payment", *fetchedSub.SuspendReason)
		t.Logf("✓ Subscription status in DB: %s (reason: %s)", fetchedSub.Status, *fetchedSub.SuspendReason)
	})

	t.Run("Isolate Subscription - PPP Secret Profile Changed", func(t *testing.T) {
		defer suite.Cleanup(t)

		// Setup
		router := createTestRouter(t)
		profile := createTestProfile(t, router.ID)

		mt, err := routerSvc.GetMikrotikClient(suite.Ctx, uuid.MustParse(router.ID))
		require.NoError(t, err)

		createPPPProfileInMikroTik(t, mt, profile.Name)
		createIsolateProfile(t, mt)

		// Create customer and subscription
		customer := &model.Customer{
			FullName: "Test Customer Isolate",
			Phone:    "081234567893",
		}
		subscription := &model.Subscription{
			PlanID:   profile.ID,
			RouterID: router.ID,
			Username: "testuser_iso_" + uuid.New().String()[:8],
			Password: "testpass123",
		}

		_, createdSub, err := customerSvc.CreateWithSubscription(suite.Ctx, customer, subscription)
		require.NoError(t, err)

		// Activate first
		subID := uuid.MustParse(createdSub.ID)
		err = subSvc.Activate(suite.Ctx, subID)
		require.NoError(t, err)

		t.Logf("✓ Subscription activated with profile: %s", profile.Name)

		// Verify original profile
		hasProfile := verifyPPPSecretProfile(t, mt, createdSub.Username, profile.Name)
		require.True(t, hasProfile, "PPP secret should have original profile")
		t.Logf("✓ PPP secret has original profile: %s", profile.Name)

		// Isolate subscription
		err = subSvc.Isolate(suite.Ctx, subID, "overdue_invoice")
		require.NoError(t, err)

		t.Logf("✓ Subscription isolated")

		// Verify PPP secret profile changed to isolate
		hasIsolateProfile := verifyPPPSecretProfile(t, mt, createdSub.Username, "isolate")
		assert.True(t, hasIsolateProfile, "PPP secret should have isolate profile")

		if hasIsolateProfile {
			t.Logf("✓ PPP secret profile changed to 'isolate'")
		}

		// Verify subscription status in DB
		fetchedSub, err := subSvc.GetByID(suite.Ctx, subID)
		require.NoError(t, err)
		assert.Equal(t, "isolated", fetchedSub.Status)
		assert.Equal(t, "overdue_invoice", *fetchedSub.SuspendReason)
		t.Logf("✓ Subscription status in DB: %s", fetchedSub.Status)
	})

	t.Run("Restore Subscription - PPP Secret Profile Restored", func(t *testing.T) {
		defer suite.Cleanup(t)

		// Setup
		router := createTestRouter(t)
		profile := createTestProfile(t, router.ID)

		mt, err := routerSvc.GetMikrotikClient(suite.Ctx, uuid.MustParse(router.ID))
		require.NoError(t, err)

		createPPPProfileInMikroTik(t, mt, profile.Name)
		createIsolateProfile(t, mt)

		// Create customer and subscription
		customer := &model.Customer{
			FullName: "Test Customer Restore",
			Phone:    "081234567894",
		}
		subscription := &model.Subscription{
			PlanID:   profile.ID,
			RouterID: router.ID,
			Username: "testuser_rst_" + uuid.New().String()[:8],
			Password: "testpass123",
		}

		_, createdSub, err := customerSvc.CreateWithSubscription(suite.Ctx, customer, subscription)
		require.NoError(t, err)

		// Activate -> Isolate -> Restore
		subID := uuid.MustParse(createdSub.ID)
		err = subSvc.Activate(suite.Ctx, subID)
		require.NoError(t, err)
		err = subSvc.Isolate(suite.Ctx, subID, "test")
		require.NoError(t, err)

		t.Logf("✓ Subscription isolated")

		// Verify isolate profile
		hasIsolate := verifyPPPSecretProfile(t, mt, createdSub.Username, "isolate")
		require.True(t, hasIsolate, "PPP secret should have isolate profile")
		t.Logf("✓ PPP secret has isolate profile")

		// Restore subscription
		err = subSvc.Restore(suite.Ctx, subID)
		require.NoError(t, err)

		t.Logf("✓ Subscription restored")

		// Verify PPP secret profile restored to original
		hasOriginalProfile := verifyPPPSecretProfile(t, mt, createdSub.Username, profile.Name)
		assert.True(t, hasOriginalProfile, "PPP secret should have original profile restored")

		if hasOriginalProfile {
			t.Logf("✓ PPP secret profile restored to: %s", profile.Name)
		}

		// Verify subscription status in DB
		fetchedSub, err := subSvc.GetByID(suite.Ctx, subID)
		require.NoError(t, err)
		assert.Equal(t, "active", fetchedSub.Status)
		assert.Nil(t, fetchedSub.SuspendReason)
		t.Logf("✓ Subscription status in DB: %s", fetchedSub.Status)
	})

	t.Run("Terminate Subscription - PPP Secret Removed", func(t *testing.T) {
		defer suite.Cleanup(t)

		// Setup
		router := createTestRouter(t)
		profile := createTestProfile(t, router.ID)

		mt, err := routerSvc.GetMikrotikClient(suite.Ctx, uuid.MustParse(router.ID))
		require.NoError(t, err)

		createPPPProfileInMikroTik(t, mt, profile.Name)

		// Create customer and subscription
		customer := &model.Customer{
			FullName: "Test Customer Terminate",
			Phone:    "081234567895",
		}
		subscription := &model.Subscription{
			PlanID:   profile.ID,
			RouterID: router.ID,
			Username: "testuser_trm_" + uuid.New().String()[:8],
			Password: "testpass123",
		}

		_, createdSub, err := customerSvc.CreateWithSubscription(suite.Ctx, customer, subscription)
		require.NoError(t, err)

		// Activate first
		subID := uuid.MustParse(createdSub.ID)
		err = subSvc.Activate(suite.Ctx, subID)
		require.NoError(t, err)

		t.Logf("✓ Subscription activated")

		// Verify PPP secret exists
		exists := verifyPPPSecretExists(t, mt, createdSub.Username)
		require.True(t, exists, "PPP secret should exist before termination")
		t.Logf("✓ PPP secret exists: %s", createdSub.Username)

		// Terminate subscription
		err = subSvc.Terminate(suite.Ctx, subID)
		require.NoError(t, err)

		t.Logf("✓ Subscription terminated")

		// Verify PPP secret removed from MikroTik
		exists = verifyPPPSecretExists(t, mt, createdSub.Username)
		assert.False(t, exists, "PPP secret should be removed from MikroTik")

		if !exists {
			t.Logf("✓ PPP secret '%s' removed from MikroTik", createdSub.Username)
		}

		// Verify subscription status in DB
		fetchedSub, err := subSvc.GetByID(suite.Ctx, subID)
		require.NoError(t, err)
		assert.Equal(t, "terminated", fetchedSub.Status)
		assert.NotNil(t, fetchedSub.TerminatedAt)
		t.Logf("✓ Subscription status in DB: %s", fetchedSub.Status)
	})

	t.Run("Full Lifecycle Test", func(t *testing.T) {
		defer suite.Cleanup(t)

		// Setup
		router := createTestRouter(t)
		profile := createTestProfile(t, router.ID)

		mt, err := routerSvc.GetMikrotikClient(suite.Ctx, uuid.MustParse(router.ID))
		require.NoError(t, err)

		createPPPProfileInMikroTik(t, mt, profile.Name)
		createIsolateProfile(t, mt)

		username := "testuser_full_" + uuid.New().String()[:8]

		// Create customer and subscription
		customer := &model.Customer{
			FullName: "Test Customer Full Lifecycle",
			Phone:    "081234567896",
		}
		subscription := &model.Subscription{
			PlanID:   profile.ID,
			RouterID: router.ID,
			Username: username,
			Password: "testpass123",
		}

		_, createdSub, err := customerSvc.CreateWithSubscription(suite.Ctx, customer, subscription)
		require.NoError(t, err)
		subID := uuid.MustParse(createdSub.ID)

		t.Logf("\n=== FULL LIFECYCLE TEST ===")
		t.Logf("1. Created subscription: %s (status: %s)", createdSub.ID, createdSub.Status)

		// Verify created
		exists := verifyPPPSecretExists(t, mt, username)
		require.True(t, exists)
		t.Logf("   ✓ PPP secret created in MikroTik")

		// Activate
		err = subSvc.Activate(suite.Ctx, subID)
		require.NoError(t, err)
		t.Logf("2. Activated subscription")

		enabled := verifyPPPSecretEnabled(t, mt, username)
		require.True(t, enabled)
		t.Logf("   ✓ PPP secret enabled")

		// Suspend
		err = subSvc.Suspend(suite.Ctx, subID, "test_suspend")
		require.NoError(t, err)
		t.Logf("3. Suspended subscription")

		enabled = verifyPPPSecretEnabled(t, mt, username)
		require.False(t, enabled)
		t.Logf("   ✓ PPP secret disabled")

		// Activate again
		err = subSvc.Activate(suite.Ctx, subID)
		require.NoError(t, err)
		t.Logf("4. Re-activated subscription")

		enabled = verifyPPPSecretEnabled(t, mt, username)
		require.True(t, enabled)
		t.Logf("   ✓ PPP secret re-enabled")

		// Isolate
		err = subSvc.Isolate(suite.Ctx, subID, "test_isolate")
		require.NoError(t, err)
		t.Logf("5. Isolated subscription")

		hasIsolate := verifyPPPSecretProfile(t, mt, username, "isolate")
		require.True(t, hasIsolate)
		t.Logf("   ✓ PPP secret profile changed to isolate")

		// Restore
		err = subSvc.Restore(suite.Ctx, subID)
		require.NoError(t, err)
		t.Logf("6. Restored subscription")

		hasOriginal := verifyPPPSecretProfile(t, mt, username, profile.Name)
		require.True(t, hasOriginal)
		t.Logf("   ✓ PPP secret profile restored to %s", profile.Name)

		// Terminate
		err = subSvc.Terminate(suite.Ctx, subID)
		require.NoError(t, err)
		t.Logf("7. Terminated subscription")

		exists = verifyPPPSecretExists(t, mt, username)
		require.False(t, exists)
		t.Logf("   ✓ PPP secret removed from MikroTik")

		t.Logf("\n=== ALL LIFECYCLE STEPS PASSED ===")
	})
}

func parsePort(port string) int {
	var p int
	fmt.Sscanf(port, "%d", &p)
	if p == 0 {
		return 8728
	}
	return p
}
