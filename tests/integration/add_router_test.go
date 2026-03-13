//go:build integration

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
	"mikmongo/pkg/mikrotik"
	"mikmongo/pkg/mikrotik/client"
	"mikmongo/tests/integration/testutil"
)

func TestAddRouterAndConnect_Integration(t *testing.T) {
	// Skip jika tidak ada database
	testutil.SkipIfNoMikrotik(t)

	// Setup test suite
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	// Create repository
	repo := postgres.NewRouterDeviceRepository(suite.DB)

	// Create router device sesuai permintaan
	router := &model.MikrotikRouter{
		ID:        uuid.New().String(),
		Name:      "Testing",
		Address:   "192.168.233.1",
		APIPort:   8728,
		RESTPort:  80,
		Username:  "admin",
		IsMaster:  false,
		IsActive:  true,
		Status:    "online",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Create service dengan encryption key
	encKey := "test-encryption-key-32-bytes-long!!"
	logger := zap.NewNop()
	routerService := service.NewRouterService(repo, encKey, nil, logger)

	// Test Create dengan password
	t.Run("Add Router via Service", func(t *testing.T) {
		err := routerService.Create(suite.Ctx, router, "r00t")
		require.NoError(t, err)
		t.Logf("✓ Router berhasil ditambahkan dengan ID: %s", router.ID)

		// Verify router was created
		id, _ := uuid.Parse(router.ID)
		fetched, err := repo.GetByID(suite.Ctx, id)
		require.NoError(t, err)

		// Assertions
		assert.Equal(t, "Testing", fetched.Name)
		assert.Equal(t, "192.168.233.1", fetched.Address)
		assert.Equal(t, 8728, fetched.APIPort)
		assert.Equal(t, "admin", fetched.Username)
		assert.True(t, fetched.IsActive)
		assert.Equal(t, "online", fetched.Status)

		t.Logf("✓ Router berhasil diverifikasi di database:")
		t.Logf("  Name: %s", fetched.Name)
		t.Logf("  Address: %s", fetched.Address)
		t.Logf("  API Port: %d", fetched.APIPort)
		t.Logf("  Username: %s", fetched.Username)
		t.Logf("  Status: %s", fetched.Status)
		t.Logf("  IsActive: %v", fetched.IsActive)
	})

	// Test koneksi ke Mikrotik
	t.Run("Connect to Mikrotik", func(t *testing.T) {
		id, _ := uuid.Parse(router.ID)
		client, err := routerService.Connect(suite.Ctx, id)
		require.NoError(t, err, "Gagal terhubung ke Mikrotik")
		defer client.Close()

		t.Logf("✓ Berhasil terhubung ke Mikrotik %s", router.Address)

		// Test get identity
		identity, err := client.Conn().Run("/system/identity/print")
		require.NoError(t, err)
		t.Logf("✓ Router Identity: %v", identity.Re)

		// Test get resource
		resource, err := client.Conn().Run("/system/resource/print")
		require.NoError(t, err)
		if len(resource.Re) > 0 {
			version := resource.Re[0].Map["version"]
			board := resource.Re[0].Map["board-name"]
			uptime := resource.Re[0].Map["uptime"]
			t.Logf("✓ Router Version: %s", version)
			t.Logf("✓ Router Board: %s", board)
			t.Logf("✓ Router Uptime: %s", uptime)
		}
	})

	// Test langsung dengan client config
	t.Run("Direct Connection Test", func(t *testing.T) {
		cfg := client.Config{
			Host:     "192.168.233.1",
			Port:     8728,
			Username: "admin",
			Password: "r00t",
			UseTLS:   false,
		}

		mtClient, err := mikrotik.NewClient(cfg)
		require.NoError(t, err, "Gagal membuat koneksi langsung")
		defer mtClient.Close()

		t.Logf("✓ Direct connection berhasil")

		// Test hotspot profile count
		hotspotProfiles, err := mtClient.Hotspot.GetProfiles(suite.Ctx)
		require.NoError(t, err)
		t.Logf("✓ Hotspot profiles: %d", len(hotspotProfiles))

		// Test PPP profile count
		pppProfiles, err := mtClient.PPP.GetProfiles(suite.Ctx)
		require.NoError(t, err)
		t.Logf("✓ PPP profiles: %d", len(pppProfiles))
	})
}

// TestCurlEquivalent menunjukkan equivalent curl command
func TestCurlEquivalent(t *testing.T) {
	t.Logf("Equivalent curl command untuk menambahkan router:")
	t.Logf("")
	t.Logf("curl -X POST http://localhost:8080/api/v1/routers \\")
	t.Logf("  -H \"Content-Type: application/json\" \\")
	t.Logf("  -d '{\"name\":\"Testing\",\"address\":\"192.168.233.1\",\"api_port\":8728,\"rest_port\":80,\"username\":\"admin\",\"password\":\"r00t\",\"is_master\":false,\"is_active\":true,\"status\":\"online\"}'")
	t.Logf("")
	t.Logf("Note: Server harus berjalan dan memiliki database PostgreSQL")
}
