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

// TestHotspot_Integration tests Hotspot operations against a real Mikrotik router
// These tests require a real Mikrotik router configured via environment variables:
// - TEST_MIKROTIK_HOST (default: 192.168.88.1)
// - TEST_MIKROTIK_PORT (default: 8728)
// - TEST_MIKROTIK_USER (default: admin)
// - TEST_MIKROTIK_PASS (required)
func TestHotspot_Integration(t *testing.T) {
	testutil.SkipIfNoMikrotik(t)

	// Setup Mikrotik client
	cfg := testutil.LoadMikrotikConfig()
	client, err := testutil.NewMikrotikClient(cfg)
	require.NoError(t, err)
	defer client.Close()

	ctx := testutil.WithTimeout(t, 30*time.Second)

	t.Run("GetProfiles", func(t *testing.T) {
		profiles, err := client.Hotspot.GetProfiles(ctx)
		require.NoError(t, err)
		// Just verify we can fetch profiles (router may have 0 or more)
		assert.NotNil(t, profiles)
		t.Logf("Found %d hotspot profiles", len(profiles))
	})

	t.Run("Profile CRUD", func(t *testing.T) {
		// Create a test profile (without address-pool to avoid validation errors)
		profile := &domain.UserProfile{
			Name:        "TEST_PROFILE_" + time.Now().Format("20060102_150405"),
			SharedUsers: 1,
			RateLimit:   "1M/1M",
		}

		// Add profile
		id, err := client.Hotspot.AddProfile(ctx, profile)
		require.NoError(t, err)
		// Note: RouterOS may return empty ret for some operations
		if id == "" {
			t.Log("Profile created but returned empty ID, trying to find by name")
			found, err := client.Hotspot.GetProfileByName(ctx, profile.Name)
			require.NoError(t, err)
			id = found.ID
		}
		assert.NotEmpty(t, id)
		t.Logf("Created profile with ID: %s", id)

		// Cleanup: remove profile after test
		defer func() {
			err := client.Hotspot.RemoveProfile(ctx, id)
			if err != nil {
				t.Logf("Warning: failed to cleanup profile: %v", err)
			}
		}()

		// Get profile by ID
		fetched, err := client.Hotspot.GetProfileByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, profile.Name, fetched.Name)
		assert.Equal(t, profile.RateLimit, fetched.RateLimit)

		// Get profile by name
		byName, err := client.Hotspot.GetProfileByName(ctx, profile.Name)
		require.NoError(t, err)
		assert.Equal(t, id, byName.ID)

		// Update profile
		profile.RateLimit = "2M/2M"
		err = client.Hotspot.UpdateProfile(ctx, id, profile)
		require.NoError(t, err)

		// Verify update
		updated, err := client.Hotspot.GetProfileByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, "2M/2M", updated.RateLimit)
	})

	t.Run("User CRUD", func(t *testing.T) {
		// First ensure we have a profile to use (without address-pool)
		profile := &domain.UserProfile{
			Name:        "TEST_USER_PROFILE_" + time.Now().Format("20060102_150405"),
			SharedUsers: 1,
			RateLimit:   "1M/1M",
		}
		profileID, err := client.Hotspot.AddProfile(ctx, profile)
		require.NoError(t, err)
		// Handle empty ID case
		if profileID == "" {
			found, err := client.Hotspot.GetProfileByName(ctx, profile.Name)
			require.NoError(t, err)
			profileID = found.ID
		}

		// Cleanup profile
		defer func() {
			_ = client.Hotspot.RemoveProfile(ctx, profileID)
		}()

		// Create user
		user := &domain.HotspotUser{
			Name:     "testuser_" + time.Now().Format("150405"),
			Password: "testpass123",
			Profile:  profile.Name,
			Comment:  "Integration test user",
		}

		userID, err := client.Hotspot.AddUser(ctx, user)
		require.NoError(t, err)
		// Note: RouterOS may return empty ret for some operations
		if userID == "" {
			t.Log("User created but returned empty ID, trying to find by name")
			found, err := client.Hotspot.GetUserByName(ctx, user.Name)
			require.NoError(t, err)
			require.NotNil(t, found, "User should exist after creation")
			userID = found.ID
		}
		assert.NotEmpty(t, userID)
		t.Logf("Created user with ID: %s", userID)

		// Cleanup user
		defer func() {
			err := client.Hotspot.RemoveUser(ctx, userID)
			if err != nil {
				t.Logf("Warning: failed to cleanup user: %v", err)
			}
		}()

		// Get user by ID
		fetched, err := client.Hotspot.GetUserByID(ctx, userID)
		require.NoError(t, err)
		require.NotNil(t, fetched)
		assert.Equal(t, user.Name, fetched.Name)
		assert.Equal(t, user.Profile, fetched.Profile)
		assert.Equal(t, user.Comment, fetched.Comment)

		// Get user by name
		byName, err := client.Hotspot.GetUserByName(ctx, user.Name)
		require.NoError(t, err)
		require.NotNil(t, byName)
		assert.Equal(t, userID, byName.ID)

		// Update user
		user.Password = "newpass456"
		user.Comment = "Updated comment"
		err = client.Hotspot.UpdateUser(ctx, userID, user)
		require.NoError(t, err)

		// Verify update
		updated, err := client.Hotspot.GetUserByID(ctx, userID)
		require.NoError(t, err)
		assert.Equal(t, "Updated comment", updated.Comment)

		// Disable user
		err = client.Hotspot.DisableUser(ctx, userID)
		require.NoError(t, err)

		disabled, err := client.Hotspot.GetUserByID(ctx, userID)
		require.NoError(t, err)
		assert.True(t, disabled.Disabled)

		// Enable user
		err = client.Hotspot.EnableUser(ctx, userID)
		require.NoError(t, err)

		enabled, err := client.Hotspot.GetUserByID(ctx, userID)
		require.NoError(t, err)
		assert.False(t, enabled.Disabled)
	})

	t.Run("GetUsersCount", func(t *testing.T) {
		count, err := client.Hotspot.GetUsersCount(ctx)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, count, 0)
		t.Logf("Total hotspot users: %d", count)
	})

	t.Run("GetActive", func(t *testing.T) {
		active, err := client.Hotspot.GetActive(ctx)
		require.NoError(t, err)
		assert.NotNil(t, active)
		t.Logf("Active hotspot sessions: %d", len(active))
	})

	t.Run("GetActiveCount", func(t *testing.T) {
		count, err := client.Hotspot.GetActiveCount(ctx)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, count, 0)
		t.Logf("Active hotspot count: %d", count)
	})

	t.Run("GetHosts", func(t *testing.T) {
		hosts, err := client.Hotspot.GetHosts(ctx)
		require.NoError(t, err)
		assert.NotNil(t, hosts)
		t.Logf("Hotspot hosts: %d", len(hosts))
	})

	t.Run("GetServers", func(t *testing.T) {
		servers, err := client.Hotspot.GetServers(ctx)
		require.NoError(t, err)
		assert.NotNil(t, servers)
		t.Logf("Hotspot servers: %v", servers)
	})

	t.Run("Batch Operations", func(t *testing.T) {
		// Create a profile for batch test users (without address-pool)
		profile := &domain.UserProfile{
			Name:        "TEST_BATCH_PROFILE_" + time.Now().Format("20060102_150405"),
			SharedUsers: 1,
			RateLimit:   "512k/512k",
		}
		profileID, err := client.Hotspot.AddProfile(ctx, profile)
		require.NoError(t, err)
		// Handle empty ID case
		if profileID == "" {
			found, err := client.Hotspot.GetProfileByName(ctx, profile.Name)
			require.NoError(t, err)
			profileID = found.ID
		}
		defer func() {
			_ = client.Hotspot.RemoveProfile(ctx, profileID)
		}()

		// Create multiple users
		var userIDs []string
		for i := 0; i < 3; i++ {
			user := &domain.HotspotUser{
				Name:     "batchuser_" + time.Now().Format("150405") + "_" + string(rune('A'+i)),
				Password: "batchpass",
				Profile:  profile.Name,
				Comment:  "Batch test user",
			}
			id, err := client.Hotspot.AddUser(ctx, user)
			require.NoError(t, err)
			if id == "" {
				found, err := client.Hotspot.GetUserByName(ctx, user.Name)
				require.NoError(t, err)
				id = found.ID
			}
			userIDs = append(userIDs, id)
		}

		// Cleanup users
		defer func() {
			err := client.Hotspot.RemoveUsers(ctx, userIDs)
			if err != nil {
				t.Logf("Warning: failed to cleanup batch users: %v", err)
			}
		}()

		// Test batch disable
		err = client.Hotspot.DisableUsers(ctx, userIDs)
		require.NoError(t, err)

		// Verify all disabled
		for _, id := range userIDs {
			user, err := client.Hotspot.GetUserByID(ctx, id)
			require.NoError(t, err)
			assert.True(t, user.Disabled)
		}

		// Test batch enable
		err = client.Hotspot.EnableUsers(ctx, userIDs)
		require.NoError(t, err)

		// Verify all enabled
		for _, id := range userIDs {
			user, err := client.Hotspot.GetUserByID(ctx, id)
			require.NoError(t, err)
			assert.False(t, user.Disabled)
		}
	})
}
