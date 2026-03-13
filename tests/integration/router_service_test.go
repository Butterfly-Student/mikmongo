//go:build integration

package integration

import (
	"testing"
	"time"

	"mikmongo/internal/model"
	"mikmongo/internal/repository/postgres"
	"mikmongo/internal/service"
	mkdomain "mikmongo/pkg/mikrotik/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// TestRouterServiceMikroTikIntegration tests RouterService with real MikroTik
func TestRouterServiceMikroTikIntegration(t *testing.T) {
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
	routerRepo := postgres.NewRouterDeviceRepository(suite.DB)

	// Create logger
	logger := zap.NewNop()

	// Create router service with encryption key
	routerSvc := service.NewRouterService(routerRepo, "test-key-16-bytes", nil, logger)

	t.Run("Create Router and Get MikroTik Client", func(t *testing.T) {
		defer suite.Cleanup(t)

		// Create router with plain text password (like test does)
		router := &model.MikrotikRouter{
			ID:                uuid.New().String(),
			Name:              "Test Router Service",
			Address:           mtHost,
			APIPort:           parsePort(mtPort),
			Username:          mtUser,
			PasswordEncrypted: mtPass, // Plain text password
			IsActive:          true,
			Status:            "online",
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		}
		err := routerRepo.Create(suite.Ctx, router)
		require.NoError(t, err)

		t.Logf("✓ Router created: %s", router.ID)

		// Get MikroTik client via RouterService
		routerID := uuid.MustParse(router.ID)
		mt, err := routerSvc.GetMikrotikClient(suite.Ctx, routerID)
		require.NoError(t, err, "Failed to get MikroTik client")

		t.Logf("✓ Got MikroTik client via RouterService")

		// Test connection by listing profiles
		profiles, err := mt.PPP.GetProfiles(suite.Ctx)
		require.NoError(t, err)
		t.Logf("✓ Successfully connected and got %d PPP profiles", len(profiles))
	})

	t.Run("TestConnection via RouterService", func(t *testing.T) {
		defer suite.Cleanup(t)

		// Create router
		router := &model.MikrotikRouter{
			ID:                uuid.New().String(),
			Name:              "Test Router Connection",
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

		// Test connection
		routerID := uuid.MustParse(router.ID)
		err = routerSvc.TestConnection(suite.Ctx, routerID)
		require.NoError(t, err, "TestConnection should succeed")

		t.Logf("✓ TestConnection succeeded")

		// Verify router status updated
		updatedRouter, err := routerRepo.GetByID(suite.Ctx, routerID)
		require.NoError(t, err)
		assert.Equal(t, "online", updatedRouter.Status)
		t.Logf("✓ Router status: %s", updatedRouter.Status)
	})

	t.Run("Create PPP Secret via RouterService", func(t *testing.T) {
		defer suite.Cleanup(t)

		// Create router
		router := &model.MikrotikRouter{
			ID:                uuid.New().String(),
			Name:              "Test Router PPP",
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

		// Get MikroTik client
		routerID := uuid.MustParse(router.ID)
		mt, err := routerSvc.GetMikrotikClient(suite.Ctx, routerID)
		require.NoError(t, err)

		// Create a test profile first
		profileName := "TEST_ROUTER_SVC_PROFILE"
		profiles, _ := mt.PPP.GetProfiles(suite.Ctx)
		exists := false
		for _, p := range profiles {
			if p.Name == profileName {
				exists = true
				break
			}
		}
		if !exists {
			profile := &mkdomain.PPPProfile{Name: profileName}
			mt.PPP.AddProfile(suite.Ctx, profile)
		}

		// Create PPP secret
		username := "test_routersvc_" + uuid.New().String()[:8]
		secret := &mkdomain.PPPSecret{
			Name:     username,
			Password: "testpass123",
			Profile:  profileName,
		}

		err = mt.PPP.AddSecret(suite.Ctx, secret)
		require.NoError(t, err)

		t.Logf("✓ Created PPP secret: %s", username)

		// Verify it exists
		foundSecret, err := mt.PPP.GetSecretByName(suite.Ctx, username)
		require.NoError(t, err)
		assert.Equal(t, username, foundSecret.Name)
		t.Logf("✓ Verified PPP secret exists")

		// Cleanup
		mt.PPP.RemoveSecret(suite.Ctx, foundSecret.ID)
		mt.PPP.RemoveProfile(suite.Ctx, profileName)
	})
}
