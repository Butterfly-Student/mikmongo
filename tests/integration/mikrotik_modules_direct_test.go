//go:build integration && mikrotik_legacy

package integration

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"mikmongo/internal/model"
	"mikmongo/internal/repository/postgres"
	"mikmongo/internal/service"
	mikrotiksvc "mikmongo/internal/service/mikrotik"
	mkdomain "mikmongo/pkg/mikrotik/domain"
)

// TestMikrotikModulesDirect tests all MikroTik modules directly using services
func TestMikrotikModulesDirect(t *testing.T) {
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

	// Create router service
	routerSvc := service.NewRouterService(routerRepo, "test-key-16-bytes", nil, logger)

	// Create MikroTik service registry
	mikrotikSvc := mikrotiksvc.NewRegistry(routerSvc)

	// Create test router in DB
	createTestRouter := func(t *testing.T) *model.MikrotikRouter {
		router := &model.MikrotikRouter{
			ID:                uuid.New().String(),
			Name:              "Test Router Modules",
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

	t.Run("PPP Module - Complete CRUD", func(t *testing.T) {
		defer suite.Cleanup(t)
		router := createTestRouter(t)
		routerUUID := uuid.MustParse(router.ID)

		t.Logf("\n=== PPP MODULE TEST ===")
		t.Logf("Router: %s", router.Name)

		// Test Get Profiles
		profiles, err := mikrotikSvc.PPP.GetProfiles(suite.Ctx, routerUUID)
		require.NoError(t, err)
		t.Logf("✓ Get Profiles: %d profiles found", len(profiles))

		// Test Add Profile
		profileReq := &mkdomain.PPPProfile{
			Name:      "TEST_PPP_PROFILE_" + uuid.New().String()[:8],
			RateLimit: "10M/10M",
		}
		err = mikrotikSvc.PPP.AddProfile(suite.Ctx, routerUUID, profileReq)
		require.NoError(t, err)
		t.Logf("✓ Add Profile: %s", profileReq.Name)

		// Get profile to find ID
		addedProfile, err := mikrotikSvc.PPP.GetProfileByName(suite.Ctx, routerUUID, profileReq.Name)
		require.NoError(t, err)
		t.Logf("✓ Get Profile by Name: %s (ID: %s)", addedProfile.Name, addedProfile.ID)

		// Test Get Secrets
		secrets, err := mikrotikSvc.PPP.GetSecrets(suite.Ctx, routerUUID, "")
		require.NoError(t, err)
		t.Logf("✓ Get Secrets: %d secrets found", len(secrets))

		// Test Add Secret
		secretReq := &mkdomain.PPPSecret{
			Name:     "testppp_" + uuid.New().String()[:8],
			Password: "testpass123",
			Profile:  profileReq.Name,
		}
		err = mikrotikSvc.PPP.AddSecret(suite.Ctx, routerUUID, secretReq)
		require.NoError(t, err)
		t.Logf("✓ Add Secret: %s", secretReq.Name)

		// Get secret to find ID
		addedSecret, err := mikrotikSvc.PPP.GetSecretByName(suite.Ctx, routerUUID, secretReq.Name)
		require.NoError(t, err)
		t.Logf("✓ Get Secret by Name: %s (ID: %s)", addedSecret.Name, addedSecret.ID)

		// Test Get Active Users
		active, err := mikrotikSvc.PPP.GetActiveUsers(suite.Ctx, routerUUID, "")
		require.NoError(t, err)
		t.Logf("✓ Get Active Users: %d active sessions", len(active))

		// Cleanup
		_ = mikrotikSvc.PPP.RemoveSecret(suite.Ctx, routerUUID, addedSecret.ID)
		_ = mikrotikSvc.PPP.RemoveProfile(suite.Ctx, routerUUID, addedProfile.ID)
		t.Logf("✓ Cleanup completed")
	})

	t.Run("Hotspot Module - Complete CRUD", func(t *testing.T) {
		defer suite.Cleanup(t)
		router := createTestRouter(t)
		routerUUID := uuid.MustParse(router.ID)

		t.Logf("\n=== HOTSPOT MODULE TEST ===")
		t.Logf("Router: %s", router.Name)

		// Test Get Profiles
		profiles, err := mikrotikSvc.Hotspot.GetProfiles(suite.Ctx, routerUUID)
		require.NoError(t, err)
		t.Logf("✓ Get Profiles: %d profiles found", len(profiles))

		// Test Create Profile
		profileReq := &mkdomain.UserProfile{
			Name:        "TEST_HS_PROFILE_" + uuid.New().String()[:8],
			SharedUsers: 1,
			RateLimit:   "5M/5M",
		}
		profileID, err := mikrotikSvc.Hotspot.AddProfile(suite.Ctx, routerUUID, profileReq)
		require.NoError(t, err)
		t.Logf("✓ Create Profile: %s (ID: %s)", profileReq.Name, profileID)

		// If profileID is empty, try to find by name
		if profileID == "" {
			profile, err := mikrotikSvc.Hotspot.GetProfileByName(suite.Ctx, routerUUID, profileReq.Name)
			require.NoError(t, err)
			require.NotNil(t, profile)
			profileID = profile.ID
			t.Logf("✓ Found Profile by Name: %s (ID: %s)", profile.Name, profileID)
		}

		// Test Get Profile
		profile, err := mikrotikSvc.Hotspot.GetProfileByID(suite.Ctx, routerUUID, profileID)
		require.NoError(t, err)
		require.NotNil(t, profile)
		assert.Equal(t, profileReq.Name, profile.Name)
		t.Logf("✓ Get Profile: %s", profile.Name)

		// Test List Users
		users, err := mikrotikSvc.Hotspot.GetUsers(suite.Ctx, routerUUID, "")
		require.NoError(t, err)
		t.Logf("✓ List Users: %d users found", len(users))

		// Test Create User
		userReq := &mkdomain.HotspotUser{
			Name:     "tesths_" + uuid.New().String()[:8],
			Password: "hspass123",
			Profile:  profileReq.Name,
		}
		userID, err := mikrotikSvc.Hotspot.AddUser(suite.Ctx, routerUUID, userReq)
		require.NoError(t, err)
		t.Logf("✓ Create User: %s (ID: %s)", userReq.Name, userID)

		// If userID is empty, try to find by name
		if userID == "" {
			user, err := mikrotikSvc.Hotspot.GetUserByName(suite.Ctx, routerUUID, userReq.Name)
			require.NoError(t, err)
			require.NotNil(t, user)
			userID = user.ID
			t.Logf("✓ Found User by Name: %s (ID: %s)", user.Name, userID)
		}

		// Test Get User
		user, err := mikrotikSvc.Hotspot.GetUserByID(suite.Ctx, routerUUID, userID)
		require.NoError(t, err)
		require.NotNil(t, user)
		assert.Equal(t, userReq.Name, user.Name)
		t.Logf("✓ Get User: %s", user.Name)

		// Test List Active
		active, err := mikrotikSvc.Hotspot.GetActive(suite.Ctx, routerUUID)
		require.NoError(t, err)
		t.Logf("✓ List Active: %d active sessions", len(active))

		// Test List Hosts
		hosts, err := mikrotikSvc.Hotspot.GetHosts(suite.Ctx, routerUUID)
		require.NoError(t, err)
		t.Logf("✓ List Hosts: %d hosts", len(hosts))

		// Test List Servers
		servers, err := mikrotikSvc.Hotspot.GetServers(suite.Ctx, routerUUID)
		require.NoError(t, err)
		t.Logf("✓ List Servers: %v", servers)

		// Cleanup
		_ = mikrotikSvc.Hotspot.RemoveUser(suite.Ctx, routerUUID, userID)
		_ = mikrotikSvc.Hotspot.RemoveProfile(suite.Ctx, routerUUID, profileID)
		t.Logf("✓ Cleanup completed")
	})

	t.Run("Queue Module", func(t *testing.T) {
		defer suite.Cleanup(t)
		router := createTestRouter(t)
		routerUUID := uuid.MustParse(router.ID)

		t.Logf("\n=== QUEUE MODULE TEST ===")
		t.Logf("Router: %s", router.Name)

		// Test Get Simple Queues
		queues, err := mikrotikSvc.Queue.GetSimpleQueues(suite.Ctx, routerUUID)
		require.NoError(t, err)
		t.Logf("✓ Get Simple Queues: %d queues found", len(queues))

		t.Logf("✓ Queue Module working correctly")
	})

	t.Run("Firewall Module", func(t *testing.T) {
		defer suite.Cleanup(t)
		router := createTestRouter(t)
		routerUUID := uuid.MustParse(router.ID)

		t.Logf("\n=== FIREWALL MODULE TEST ===")
		t.Logf("Router: %s", router.Name)

		// Test Get Filter Rules
		filters, err := mikrotikSvc.Firewall.GetFilterRules(suite.Ctx, routerUUID)
		require.NoError(t, err)
		t.Logf("✓ Get Filter Rules: %d rules found", len(filters))

		// Test Get NAT Rules
		nats, err := mikrotikSvc.Firewall.GetNATRules(suite.Ctx, routerUUID)
		require.NoError(t, err)
		t.Logf("✓ Get NAT Rules: %d rules found", len(nats))

		// Test Get Address Lists
		addrLists, err := mikrotikSvc.Firewall.GetAddressLists(suite.Ctx, routerUUID)
		require.NoError(t, err)
		t.Logf("✓ Get Address Lists: %d lists found", len(addrLists))

		t.Logf("✓ Firewall Module working correctly")
	})

	t.Run("IP Pool Module", func(t *testing.T) {
		defer suite.Cleanup(t)
		router := createTestRouter(t)
		routerUUID := uuid.MustParse(router.ID)

		t.Logf("\n=== IP POOL MODULE TEST ===")
		t.Logf("Router: %s", router.Name)

		// Test Get Pools
		pools, err := mikrotikSvc.IPPool.GetPools(suite.Ctx, routerUUID)
		require.NoError(t, err)
		t.Logf("✓ Get Pools: %d pools found", len(pools))

		// Test Add Pool
		poolReq := &mkdomain.IPPool{
			Name:   "TEST_POOL_" + uuid.New().String()[:8],
			Ranges: "192.168.200.10-192.168.200.100",
		}
		poolID, err := mikrotikSvc.IPPool.AddPool(suite.Ctx, routerUUID, poolReq)
		require.NoError(t, err)
		t.Logf("✓ Add Pool: %s (ID: %s)", poolReq.Name, poolID)

		// If poolID is empty, try to find by name
		if poolID == "" {
			pool, err := mikrotikSvc.IPPool.GetPoolByName(suite.Ctx, routerUUID, poolReq.Name)
			require.NoError(t, err)
			require.NotNil(t, pool)
			poolID = pool.ID
			t.Logf("✓ Found Pool by Name: %s (ID: %s)", pool.Name, poolID)
		}

		// Test Get Pool by ID
		pool, err := mikrotikSvc.IPPool.GetPoolByID(suite.Ctx, routerUUID, poolID)
		require.NoError(t, err)
		require.NotNil(t, pool)
		assert.Equal(t, poolReq.Name, pool.Name)
		t.Logf("✓ Get Pool by ID: %s", pool.Name)

		// Cleanup
		_ = mikrotikSvc.IPPool.RemovePool(suite.Ctx, routerUUID, poolID)
		t.Logf("✓ Cleanup completed")
	})

	t.Run("IP Address Module", func(t *testing.T) {
		defer suite.Cleanup(t)
		router := createTestRouter(t)
		routerUUID := uuid.MustParse(router.ID)

		t.Logf("\n=== IP ADDRESS MODULE TEST ===")
		t.Logf("Router: %s", router.Name)

		// Test Get Addresses
		addresses, err := mikrotikSvc.IPAddress.GetAddresses(suite.Ctx, routerUUID)
		require.NoError(t, err)
		t.Logf("✓ Get Addresses: %d addresses found", len(addresses))

		t.Logf("✓ IP Address Module working correctly")
	})

	t.Run("Monitor Module", func(t *testing.T) {
		defer suite.Cleanup(t)
		router := createTestRouter(t)
		routerUUID := uuid.MustParse(router.ID)

		t.Logf("\n=== MONITOR MODULE TEST ===")
		t.Logf("Router: %s", router.Name)

		// Test Get System Resource
		resources, err := mikrotikSvc.Monitor.GetSystemResource(suite.Ctx, routerUUID)
		require.NoError(t, err)
		t.Logf("✓ System Resource:")
		t.Logf("  - CPU Load: %d", resources.CpuLoad)
		t.Logf("  - Free Memory: %d", resources.FreeMemory)
		t.Logf("  - Total Memory: %d", resources.TotalMemory)
		t.Logf("  - Uptime: %s", resources.Uptime)

		// Test Get Interfaces
		interfaces, err := mikrotikSvc.Monitor.GetInterfaces(suite.Ctx, routerUUID)
		require.NoError(t, err)
		t.Logf("✓ Interfaces: %d interfaces found", len(interfaces))

		t.Logf("✓ Monitor Module working correctly")
	})

	t.Run("Report Module", func(t *testing.T) {
		defer suite.Cleanup(t)
		router := createTestRouter(t)
		routerUUID := uuid.MustParse(router.ID)

		t.Logf("\n=== REPORT MODULE TEST ===")
		t.Logf("Router: %s", router.Name)

		// Test Get Sales Reports
		reports, err := mikrotikSvc.Report.GetSalesReports(suite.Ctx, routerUUID, "")
		require.NoError(t, err)
		t.Logf("✓ Sales Reports: %d reports found", len(reports))

		t.Logf("✓ Report Module working correctly")
	})

	t.Run("Script Module", func(t *testing.T) {
		// Script module tidak memerlukan router ID

		t.Logf("\n=== SCRIPT MODULE TEST ===")

		// Test Generate OnLogin Script
		scriptReq := &mkdomain.ProfileRequest{
			Name:     "TEST_PROFILE",
			Validity: "1h",
		}
		script := mikrotikSvc.Script.GenerateOnLoginScript(scriptReq)
		assert.NotEmpty(t, script)
		t.Logf("✓ Generate OnLogin Script: %d characters", len(script))

		// Test Generate Expired Action
		expiredScript := mikrotikSvc.Script.GenerateExpiredAction("isolate")
		assert.NotEmpty(t, expiredScript)
		t.Logf("✓ Generate Expired Action: %d characters", len(expiredScript))

		t.Logf("✓ Script Module working correctly")
	})

	t.Run("All Modules Integration", func(t *testing.T) {
		defer suite.Cleanup(t)
		router := createTestRouter(t)
		routerUUID := uuid.MustParse(router.ID)

		t.Logf("\n=== ALL MODULES INTEGRATION TEST ===")
		t.Logf("Router: %s", router.Name)

		// Create PPP Profile
		pppProfile := &mkdomain.PPPProfile{
			Name:      "INTEGRATION_PPP_" + uuid.New().String()[:8],
			RateLimit: "5M/5M",
		}
		err := mikrotikSvc.PPP.AddProfile(suite.Ctx, routerUUID, pppProfile)
		require.NoError(t, err)
		t.Logf("✓ PPP Profile created: %s", pppProfile.Name)

		// Create Hotspot Profile
		hsProfile := &mkdomain.UserProfile{
			Name:        "INTEGRATION_HS_" + uuid.New().String()[:8],
			SharedUsers: 1,
			RateLimit:   "3M/3M",
		}
		hsProfileID, err := mikrotikSvc.Hotspot.AddProfile(suite.Ctx, routerUUID, hsProfile)
		require.NoError(t, err)
		t.Logf("✓ Hotspot Profile created: %s (ID: %s)", hsProfile.Name, hsProfileID)

		// Create IP Pool
		pool := &mkdomain.IPPool{
			Name:   "INTEGRATION_POOL_" + uuid.New().String()[:8],
			Ranges: "192.168.250.10-192.168.250.50",
		}
		poolID, err := mikrotikSvc.IPPool.AddPool(suite.Ctx, routerUUID, pool)
		require.NoError(t, err)
		t.Logf("✓ IP Pool created: %s (ID: %s)", pool.Name, poolID)

		// Get System Resource
		resources, err := mikrotikSvc.Monitor.GetSystemResource(suite.Ctx, routerUUID)
		require.NoError(t, err)
		t.Logf("✓ System Resource: CPU %d, Memory %d", resources.CpuLoad, resources.FreeMemory)

		// Get Sales Reports
		reports, err := mikrotikSvc.Report.GetSalesReports(suite.Ctx, routerUUID, "")
		require.NoError(t, err)
		t.Logf("✓ Sales Reports: %d reports found", len(reports))

		// Cleanup PPP Profile
		pppProfileData, _ := mikrotikSvc.PPP.GetProfileByName(suite.Ctx, routerUUID, pppProfile.Name)
		if pppProfileData != nil {
			_ = mikrotikSvc.PPP.RemoveProfile(suite.Ctx, routerUUID, pppProfileData.ID)
		}
		_ = mikrotikSvc.Hotspot.RemoveProfile(suite.Ctx, routerUUID, hsProfileID)
		_ = mikrotikSvc.IPPool.RemovePool(suite.Ctx, routerUUID, poolID)
		t.Logf("✓ All resources cleaned up")

		t.Logf("\n=== ALL MODULES INTEGRATION SUCCESSFUL ===")
	})
}
