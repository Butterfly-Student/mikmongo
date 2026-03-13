//go:build integration

package integration

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"mikmongo/pkg/mikrotik"
	"mikmongo/pkg/mikrotik/client"
	mkdomain "mikmongo/pkg/mikrotik/domain"
)

const defaultTimeout = 30 * time.Second

// TestDirectPPPOperations tests PPP operations directly on MikroTik
func TestDirectPPPOperations(t *testing.T) {
	// Get MikroTik connection details from env
	mtHost := getEnv("TEST_MIKROTIK_HOST", "192.168.233.1")
	mtPort := getEnv("TEST_MIKROTIK_PORT", "8728")
	mtUser := getEnv("TEST_MIKROTIK_USER", "admin")
	mtPass := getEnv("TEST_MIKROTIK_PASS", "")

	if mtPass == "" {
		t.Skip("TEST_MIKROTIK_PASS not set, skipping MikroTik integration test")
	}

	ctx := WithTimeout(t, defaultTimeout)

	t.Run("Connect to MikroTik", func(t *testing.T) {
		cfg := client.Config{
			Host:     mtHost,
			Port:     parsePort(mtPort),
			Username: mtUser,
			Password: mtPass,
		}
		client, err := mikrotik.NewClient(cfg)
		require.NoError(t, err)
		defer client.Close()

		t.Logf("✓ Connected to MikroTik at %s", mtHost)
	})

	t.Run("List PPP Profiles", func(t *testing.T) {
		cfg := client.Config{
			Host:     mtHost,
			Port:     parsePort(mtPort),
			Username: mtUser,
			Password: mtPass,
		}
		client, err := mikrotik.NewClient(cfg)
		require.NoError(t, err)
		defer client.Close()

		profiles, err := client.PPP.GetProfiles(ctx)
		require.NoError(t, err)

		t.Logf("✓ Found %d PPP profiles:", len(profiles))
		for _, p := range profiles {
			t.Logf("  - %s (ID: %s)", p.Name, p.ID)
		}
	})

	t.Run("Create PPP Profile", func(t *testing.T) {
		cfg := client.Config{
			Host:     mtHost,
			Port:     parsePort(mtPort),
			Username: mtUser,
			Password: mtPass,
		}
		client, err := mikrotik.NewClient(cfg)
		require.NoError(t, err)
		defer client.Close()

		profileName := "TEST_PROFILE_INTEGRATION"

		// Check if exists
		profiles, _ := client.PPP.GetProfiles(ctx)
		exists := false
		for _, p := range profiles {
			if p.Name == profileName {
				exists = true
				break
			}
		}

		if !exists {
			profile := &mkdomain.PPPProfile{
				Name: profileName,
			}
			err := client.PPP.AddProfile(ctx, profile)
			require.NoError(t, err)
			t.Logf("✓ Created PPP profile: %s", profileName)
		} else {
			t.Logf("✓ PPP profile already exists: %s", profileName)
		}
	})

	t.Run("Create PPP Secret", func(t *testing.T) {
		cfg := client.Config{
			Host:     mtHost,
			Port:     parsePort(mtPort),
			Username: mtUser,
			Password: mtPass,
		}
		client, err := mikrotik.NewClient(cfg)
		require.NoError(t, err)
		defer client.Close()

		// First ensure profile exists
		profileName := "TEST_PROFILE_INTEGRATION"
		profiles, _ := client.PPP.GetProfiles(ctx)
		profileExists := false
		for _, p := range profiles {
			if p.Name == profileName {
				profileExists = true
				break
			}
		}

		if !profileExists {
			profile := &mkdomain.PPPProfile{
				Name: profileName,
			}
			client.PPP.AddProfile(ctx, profile)
		}

		// Create secret
		secret := &mkdomain.PPPSecret{
			Name:     "test_user_integration",
			Password: "testpass123",
			Profile:  profileName,
			Comment:  "Integration test secret",
		}

		err = client.PPP.AddSecret(ctx, secret)
		require.NoError(t, err)

		t.Logf("✓ Created PPP secret: %s", secret.Name)

		// Verify it exists
		found, err := client.PPP.GetSecretByName(ctx, secret.Name)
		require.NoError(t, err)
		assert.Equal(t, secret.Name, found.Name)
		assert.Equal(t, secret.Password, found.Password)
		assert.Equal(t, secret.Profile, found.Profile)

		t.Logf("✓ Verified PPP secret exists with correct data")
	})

	t.Run("Enable/Disable PPP Secret", func(t *testing.T) {
		cfg := client.Config{
			Host:     mtHost,
			Port:     parsePort(mtPort),
			Username: mtUser,
			Password: mtPass,
		}
		client, err := mikrotik.NewClient(cfg)
		require.NoError(t, err)
		defer client.Close()

		secretName := "test_user_integration"

		// Get secret
		secret, err := client.PPP.GetSecretByName(ctx, secretName)
		require.NoError(t, err)

		t.Logf("✓ Found secret: %s (Disabled: %v)", secret.Name, secret.Disabled)

		// Disable it
		err = client.PPP.DisableSecret(ctx, secret.ID)
		require.NoError(t, err)

		// Verify disabled
		secret, _ = client.PPP.GetSecretByName(ctx, secretName)
		assert.True(t, secret.Disabled)
		t.Logf("✓ Secret disabled")

		// Enable it
		err = client.PPP.EnableSecret(ctx, secret.ID)
		require.NoError(t, err)

		// Verify enabled
		secret, _ = client.PPP.GetSecretByName(ctx, secretName)
		assert.False(t, secret.Disabled)
		t.Logf("✓ Secret enabled")
	})

	t.Run("Update PPP Secret Profile", func(t *testing.T) {
		cfg := client.Config{
			Host:     mtHost,
			Port:     parsePort(mtPort),
			Username: mtUser,
			Password: mtPass,
		}
		client, err := mikrotik.NewClient(cfg)
		require.NoError(t, err)
		defer client.Close()

		// Create isolate profile if not exists
		profiles, _ := client.PPP.GetProfiles(ctx)
		isolateExists := false
		for _, p := range profiles {
			if p.Name == "isolate" {
				isolateExists = true
				break
			}
		}
		if !isolateExists {
			isolateProfile := &mkdomain.PPPProfile{Name: "isolate"}
			client.PPP.AddProfile(ctx, isolateProfile)
			t.Logf("✓ Created isolate profile")
		}

		secretName := "test_user_integration"
		secret, err := client.PPP.GetSecretByName(ctx, secretName)
		require.NoError(t, err)

		originalProfile := secret.Profile
		t.Logf("✓ Original profile: %s", originalProfile)

		// Change to isolate
		update := &mkdomain.PPPSecret{
			Profile: "isolate",
		}
		err = client.PPP.UpdateSecret(ctx, secret.ID, update)
		require.NoError(t, err)

		// Verify changed
		secret, _ = client.PPP.GetSecretByName(ctx, secretName)
		assert.Equal(t, "isolate", secret.Profile)
		t.Logf("✓ Profile changed to: isolate")

		// Restore original
		update.Profile = originalProfile
		err = client.PPP.UpdateSecret(ctx, secret.ID, update)
		require.NoError(t, err)

		// Verify restored
		secret, _ = client.PPP.GetSecretByName(ctx, secretName)
		assert.Equal(t, originalProfile, secret.Profile)
		t.Logf("✓ Profile restored to: %s", originalProfile)
	})

	t.Run("Remove PPP Secret", func(t *testing.T) {
		cfg := client.Config{
			Host:     mtHost,
			Port:     parsePort(mtPort),
			Username: mtUser,
			Password: mtPass,
		}
		client, err := mikrotik.NewClient(cfg)
		require.NoError(t, err)
		defer client.Close()

		secretName := "test_user_integration"

		// Get secret
		secret, err := client.PPP.GetSecretByName(ctx, secretName)
		require.NoError(t, err)

		// Remove it
		err = client.PPP.RemoveSecret(ctx, secret.ID)
		require.NoError(t, err)

		t.Logf("✓ Removed PPP secret: %s", secretName)

		// Verify removed
		_, err = client.PPP.GetSecretByName(ctx, secretName)
		assert.Error(t, err)
		t.Logf("✓ Verified secret removed")
	})

	t.Run("Cleanup Test Profile", func(t *testing.T) {
		cfg := client.Config{
			Host:     mtHost,
			Port:     parsePort(mtPort),
			Username: mtUser,
			Password: mtPass,
		}
		client, err := mikrotik.NewClient(cfg)
		require.NoError(t, err)
		defer client.Close()

		profileName := "TEST_PROFILE_INTEGRATION"

		// Find and remove profile
		profiles, _ := client.PPP.GetProfiles(ctx)
		for _, p := range profiles {
			if p.Name == profileName {
				client.PPP.RemoveProfile(ctx, p.ID)
				t.Logf("✓ Removed test profile: %s", profileName)
				break
			}
		}
	})
}
