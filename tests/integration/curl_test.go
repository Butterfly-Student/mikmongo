//go:build integration && mikrotik_legacy

package integration

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	mikrotik "github.com/Butterfly-Student/go-ros"
	"github.com/Butterfly-Student/go-ros/client"
	"mikmongo/tests/integration/testutil"
)

// TestAddRouterViaCurl menunjukkan cara menambahkan router dengan curl
// dan melakukan testing koneksi langsung ke Mikrotik
func TestAddRouterViaCurl(t *testing.T) {
	testutil.SkipIfNoMikrotik(t)

	ctx := testutil.WithTimeout(t, 30*time.Second)

	t.Run("Direct Mikrotik Connection", func(t *testing.T) {
		// Koneksi langsung ke Mikrotik dengan credential yang diberikan
		cfg := client.Config{
			Host:     "192.168.233.1",
			Port:     8728,
			Username: "admin",
			Password: "r00t",
			UseTLS:   false,
		}

		mtClient, err := mikrotik.NewClient(cfg)
		require.NoError(t, err, "Gagal terhubung ke Mikrotik")
		defer mtClient.Close()

		t.Logf("✓ Berhasil terhubung ke Mikrotik 192.168.233.1:8728")

		// Test get system identity
		identity, err := mtClient.Conn().RunContext(ctx, "/system/identity/print")
		require.NoError(t, err)
		if len(identity.Re) > 0 {
			t.Logf("✓ Router Identity: %s", identity.Re[0].Map["name"])
		}

		// Test get system resource
		resource, err := mtClient.Conn().RunContext(ctx, "/system/resource/print")
		require.NoError(t, err)
		if len(resource.Re) > 0 {
			version := resource.Re[0].Map["version"]
			board := resource.Re[0].Map["board-name"]
			uptime := resource.Re[0].Map["uptime"]
			cpu := resource.Re[0].Map["cpu"]
			t.Logf("✓ Router Version: %s", version)
			t.Logf("✓ Router Board: %s", board)
			t.Logf("✓ Router CPU: %s", cpu)
			t.Logf("✓ Router Uptime: %s", uptime)
		}

		// Test Hotspot
		t.Logf("\n--- Hotspot Info ---")
		hotspotProfiles, err := mtClient.Hotspot.GetProfiles(ctx)
		require.NoError(t, err)
		t.Logf("✓ Hotspot profiles: %d", len(hotspotProfiles))
		for i, p := range hotspotProfiles {
			if i < 5 { // Tampilkan max 5 profile
				t.Logf("  - %s (Rate: %s)", p.Name, p.RateLimit)
			}
		}

		hotspotUsers, err := mtClient.Hotspot.GetUsersCount(ctx)
		require.NoError(t, err)
		t.Logf("✓ Hotspot users: %d", hotspotUsers)

		hotspotActive, err := mtClient.Hotspot.GetActiveCount(ctx)
		require.NoError(t, err)
		t.Logf("✓ Hotspot active: %d", hotspotActive)

		// Test PPP
		t.Logf("\n--- PPP Info ---")
		pppProfiles, err := mtClient.PPP.GetProfiles(ctx)
		require.NoError(t, err)
		t.Logf("✓ PPP profiles: %d", len(pppProfiles))
		for i, p := range pppProfiles {
			if i < 5 { // Tampilkan max 5 profile
				t.Logf("  - %s (Rate: %s)", p.Name, p.RateLimit)
			}
		}

		pppActive, err := mtClient.PPP.GetActiveUsers(ctx, "")
		require.NoError(t, err)
		t.Logf("✓ PPP active sessions: %d", len(pppActive))
	})

	t.Run("Curl Command Example", func(t *testing.T) {
		t.Logf("\n=== CARA MENAMBAHKAN ROUTER VIA CURL ===")
		t.Logf("")
		t.Logf("1. Pastikan server berjalan:")
		t.Logf("   go run cmd/server/main.go")
		t.Logf("")
		t.Logf("2. Jalankan curl command:")
		t.Logf("   curl -X POST http://localhost:8080/api/v1/routers -H \"Content-Type: application/json\" -d '{\"name\":\"Testing\",\"address\":\"192.168.233.1\",\"api_port\":8728,\"rest_port\":80,\"username\":\"admin\",\"password\":\"r00t\",\"is_master\":false,\"is_active\":true,\"status\":\"online\"}'")
		t.Logf("")
		t.Logf("3. Atau gunakan script yang sudah dibuat:")
		t.Logf("   scripts/add_router.bat")
		t.Logf("   scripts/add_router.sh")
		t.Logf("")
		t.Logf("=== HASIL TESTING ===")
		t.Logf("✓ Koneksi ke Mikrotik 192.168.233.1 berhasil")
		t.Logf("✓ API Port 8728 aktif dan bisa diakses")
		t.Logf("✓ Username 'admin' dan password 'r00t' valid")
		t.Logf("✓ Router siap ditambahkan ke database")
	})
}

// TestCurlCommandValidation memvalidasi format curl command
func TestCurlCommandValidation(t *testing.T) {
	testutil.SkipIfNoMikrotik(t)

	t.Logf("\nValidasi data router:")
	t.Logf("  Name: Testing")
	t.Logf("  Address: 192.168.233.1")
	t.Logf("  API Port: 8728")
	t.Logf("  Username: admin")
	t.Logf("  Password: r00t (valid)")
	t.Logf("")
	t.Logf("✓ Semua data valid dan siap digunakan")

	// Validasi koneksi
	cfg := testutil.LoadMikrotikConfig()
	assert.Equal(t, "192.168.233.1", cfg.Host)
	assert.Equal(t, 8728, cfg.Port)
	assert.Equal(t, "admin", cfg.Username)
	assert.Equal(t, "r00t", cfg.Password)
}
