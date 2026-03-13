//go:build integration

package integration

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"mikmongo/pkg/mikrotik/domain"
	"mikmongo/tests/integration/testutil"
)

// TestPPP_Integration tests PPP operations against a real Mikrotik router
// These tests require a real Mikrotik router configured via environment variables:
// - TEST_MIKROTIK_HOST (default: 192.168.88.1)
// - TEST_MIKROTIK_PORT (default: 8728)
// - TEST_MIKROTIK_USER (default: admin)
// - TEST_MIKROTIK_PASS (required)
func TestPPP_Integration(t *testing.T) {
	testutil.SkipIfNoMikrotik(t)

	// Setup Mikrotik client
	cfg := testutil.LoadMikrotikConfig()
	client, err := testutil.NewMikrotikClient(cfg)
	require.NoError(t, err)
	defer client.Close()

	ctx := testutil.WithTimeout(t, 30*time.Second)

	t.Run("GetProfiles", func(t *testing.T) {
		profiles, err := client.PPP.GetProfiles(ctx)
		require.NoError(t, err)
		assert.NotNil(t, profiles)
		t.Logf("Found %d PPP profiles", len(profiles))
	})

	t.Run("Profile CRUD", func(t *testing.T) {
		// Create a test profile (without local/remote address to avoid pool validation)
		profile := &domain.PPPProfile{
			Name:      "TEST_PPP_PROFILE_" + time.Now().Format("20060102_150405"),
			RateLimit: "1M/1M",
			Comment:   "Integration test profile",
		}

		// Add profile
		err := client.PPP.AddProfile(ctx, profile)
		require.NoError(t, err)
		t.Logf("Created PPP profile: %s", profile.Name)

		// Cleanup: remove profile after test
		defer func() {
			// Find the profile by name to get ID
			profiles, _ := client.PPP.GetProfiles(ctx)
			for _, p := range profiles {
				if p.Name == profile.Name {
					_ = client.PPP.RemoveProfile(ctx, p.ID)
					break
				}
			}
		}()

		// Get profile by name
		fetched, err := client.PPP.GetProfileByName(ctx, profile.Name)
		require.NoError(t, err)
		assert.Equal(t, profile.Name, fetched.Name)
		assert.Equal(t, profile.LocalAddress, fetched.LocalAddress)
		assert.Equal(t, profile.RateLimit, fetched.RateLimit)

		// Get profile by ID
		byID, err := client.PPP.GetProfileByID(ctx, fetched.ID)
		require.NoError(t, err)
		assert.Equal(t, profile.Name, byID.Name)

		// Update profile
		profile.RateLimit = "2M/2M"
		err = client.PPP.UpdateProfile(ctx, fetched.ID, profile)
		require.NoError(t, err)

		// Verify update
		updated, err := client.PPP.GetProfileByID(ctx, fetched.ID)
		require.NoError(t, err)
		assert.Equal(t, "2M/2M", updated.RateLimit)
	})

	t.Run("Secret CRUD", func(t *testing.T) {
		// First ensure we have a profile to use (without local/remote address)
		profile := &domain.PPPProfile{
			Name:      "TEST_SECRET_PROFILE_" + time.Now().Format("20060102_150405"),
			RateLimit: "1M/1M",
		}
		err := client.PPP.AddProfile(ctx, profile)
		require.NoError(t, err)

		// Cleanup profile
		defer func() {
			profiles, _ := client.PPP.GetProfiles(ctx)
			for _, p := range profiles {
				if p.Name == profile.Name {
					_ = client.PPP.RemoveProfile(ctx, p.ID)
					break
				}
			}
		}()

		// Create secret
		secret := &domain.PPPSecret{
			Name:     "testsecret_" + time.Now().Format("150405"),
			Password: "testpass123",
			Profile:  profile.Name,
			Service:  "pppoe",
			Comment:  "Integration test secret",
		}

		err = client.PPP.AddSecret(ctx, secret)
		require.NoError(t, err)
		t.Logf("Created PPP secret: %s", secret.Name)

		// Cleanup secret
		defer func() {
			secrets, _ := client.PPP.GetSecrets(ctx, "")
			for _, s := range secrets {
				if s.Name == secret.Name {
					_ = client.PPP.RemoveSecret(ctx, s.ID)
					break
				}
			}
		}()

		// Get secret by name
		fetched, err := client.PPP.GetSecretByName(ctx, secret.Name)
		require.NoError(t, err)
		assert.Equal(t, secret.Name, fetched.Name)
		assert.Equal(t, secret.Profile, fetched.Profile)
		assert.Equal(t, secret.Service, fetched.Service)
		assert.Equal(t, secret.Comment, fetched.Comment)

		// Get secret by ID
		byID, err := client.PPP.GetSecretByID(ctx, fetched.ID)
		require.NoError(t, err)
		assert.Equal(t, secret.Name, byID.Name)

		// Update secret
		secret.Password = "newpass456"
		secret.Comment = "Updated comment"
		err = client.PPP.UpdateSecret(ctx, fetched.ID, secret)
		require.NoError(t, err)

		// Verify update
		updated, err := client.PPP.GetSecretByID(ctx, fetched.ID)
		require.NoError(t, err)
		assert.Equal(t, "Updated comment", updated.Comment)

		// Disable secret
		err = client.PPP.DisableSecret(ctx, fetched.ID)
		require.NoError(t, err)

		disabled, err := client.PPP.GetSecretByID(ctx, fetched.ID)
		require.NoError(t, err)
		assert.True(t, disabled.Disabled)

		// Enable secret
		err = client.PPP.EnableSecret(ctx, fetched.ID)
		require.NoError(t, err)

		enabled, err := client.PPP.GetSecretByID(ctx, fetched.ID)
		require.NoError(t, err)
		assert.False(t, enabled.Disabled)
	})

	t.Run("GetSecrets with Profile Filter", func(t *testing.T) {
		// Create a profile (without local/remote address)
		profile := &domain.PPPProfile{
			Name: "TEST_FILTER_PROFILE_" + time.Now().Format("20060102_150405"),
		}
		err := client.PPP.AddProfile(ctx, profile)
		require.NoError(t, err)

		// Cleanup
		defer func() {
			profiles, _ := client.PPP.GetProfiles(ctx)
			for _, p := range profiles {
				if p.Name == profile.Name {
					_ = client.PPP.RemoveProfile(ctx, p.ID)
					break
				}
			}
		}()

		// Create secrets with this profile
		for i := 0; i < 3; i++ {
			secret := &domain.PPPSecret{
				Name:     "filteruser_" + time.Now().Format("150405") + "_" + string(rune('A'+i)),
				Password: "testpass",
				Profile:  profile.Name,
				Service:  "pppoe",
			}
			err := client.PPP.AddSecret(ctx, secret)
			require.NoError(t, err)
		}

		// Get secrets by profile
		secrets, err := client.PPP.GetSecrets(ctx, profile.Name)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(secrets), 3)

		// Cleanup secrets
		defer func() {
			allSecrets, _ := client.PPP.GetSecrets(ctx, "")
			for _, s := range allSecrets {
				if s.Profile == profile.Name {
					_ = client.PPP.RemoveSecret(ctx, s.ID)
				}
			}
		}()
	})

	t.Run("GetActiveUsers", func(t *testing.T) {
		active, err := client.PPP.GetActiveUsers(ctx, "")
		require.NoError(t, err)
		assert.NotNil(t, active)
		t.Logf("Active PPP sessions: %d", len(active))
	})

	t.Run("Batch Operations", func(t *testing.T) {
		// Create a profile for batch test (without local/remote address)
		profile := &domain.PPPProfile{
			Name:      "TEST_BATCH_PPP_PROFILE_" + time.Now().Format("20060102_150405"),
			RateLimit: "512k/512k",
		}
		err := client.PPP.AddProfile(ctx, profile)
		require.NoError(t, err)

		// Cleanup profile
		defer func() {
			profiles, _ := client.PPP.GetProfiles(ctx)
			for _, p := range profiles {
				if p.Name == profile.Name {
					_ = client.PPP.RemoveProfile(ctx, p.ID)
					break
				}
			}
		}()

		// Create multiple secrets
		var secretIDs []string
		for i := 0; i < 3; i++ {
			secret := &domain.PPPSecret{
				Name:     "batchsecret_" + time.Now().Format("150405") + "_" + string(rune('A'+i)),
				Password: "batchpass",
				Profile:  profile.Name,
				Service:  "pppoe",
			}
			err := client.PPP.AddSecret(ctx, secret)
			require.NoError(t, err)

			// Get the ID
			created, _ := client.PPP.GetSecretByName(ctx, secret.Name)
			if created != nil {
				secretIDs = append(secretIDs, created.ID)
			}
		}

		// Cleanup secrets
		defer func() {
			_ = client.PPP.RemoveSecrets(ctx, secretIDs)
		}()

		if len(secretIDs) == 0 {
			t.Skip("No secrets created, skipping batch operations test")
		}

		// Test batch disable
		err = client.PPP.DisableSecrets(ctx, secretIDs)
		require.NoError(t, err)

		// Verify all disabled
		for _, id := range secretIDs {
			secret, err := client.PPP.GetSecretByID(ctx, id)
			if err == nil {
				assert.True(t, secret.Disabled)
			}
		}

		// Test batch enable
		err = client.PPP.EnableSecrets(ctx, secretIDs)
		require.NoError(t, err)

		// Verify all enabled
		for _, id := range secretIDs {
			secret, err := client.PPP.GetSecretByID(ctx, id)
			if err == nil {
				assert.False(t, secret.Disabled)
			}
		}
	})
}
