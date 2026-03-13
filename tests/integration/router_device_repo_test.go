//go:build integration

package integration

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"mikmongo/internal/model"
	"mikmongo/internal/repository/postgres"
)

func TestRouterDeviceRepository_Integration(t *testing.T) {
	// Setup test suite
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)

	// Create repository
	repo := postgres.NewRouterDeviceRepository(suite.DB)

	t.Run("Create and Get Router Device", func(t *testing.T) {
		defer suite.Cleanup(t)

		// Create test router
		router := &model.MikrotikRouter{
			ID:                uuid.New().String(),
			Name:              "Router Test 01",
			Address:           "192.168.88.1",
			APIPort:           8728,
			RESTPort:          80,
			Username:          "admin",
			PasswordEncrypted: "encrypted_password_here",
			IsMaster:          true,
			IsActive:          true,
			Status:            "online",
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		}

		// Test Create
		err := repo.Create(suite.Ctx, router)
		require.NoError(t, err)

		// Test GetByID
		id, err := uuid.Parse(router.ID)
		require.NoError(t, err)

		fetched, err := repo.GetByID(suite.Ctx, id)
		require.NoError(t, err)
		assert.Equal(t, router.Name, fetched.Name)
		assert.Equal(t, router.Address, fetched.Address)
		assert.Equal(t, router.APIPort, fetched.APIPort)
		assert.Equal(t, router.Username, fetched.Username)
		assert.Equal(t, router.IsMaster, fetched.IsMaster)
		assert.Equal(t, router.IsActive, fetched.IsActive)
	})

	t.Run("GetActive Routers", func(t *testing.T) {
		defer suite.Cleanup(t)

		// Create active routers
		for i := 0; i < 3; i++ {
			router := &model.MikrotikRouter{
				ID:                uuid.New().String(),
				Name:              "Active Router " + string(rune('A'+i)),
				Address:           "192.168.88." + string(rune('1'+i)),
				APIPort:           8728,
				Username:          "admin",
				PasswordEncrypted: "encrypted",
				IsActive:          true,
				Status:            "online",
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
			}
			err := repo.Create(suite.Ctx, router)
			require.NoError(t, err)
		}

		// Create inactive router
		inactiveRouter := &model.MikrotikRouter{
			ID:                uuid.New().String(),
			Name:              "Inactive Router",
			Address:           "192.168.88.100",
			APIPort:           8728,
			Username:          "admin",
			PasswordEncrypted: "encrypted",
			IsActive:          false,
			Status:            "offline",
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		}
		err := repo.Create(suite.Ctx, inactiveRouter)
		require.NoError(t, err)

		// Test GetActive
		activeRouters, err := repo.GetActive(suite.Ctx)
		require.NoError(t, err)
		assert.Len(t, activeRouters, 3)

		// Verify all returned routers are active
		for _, r := range activeRouters {
			assert.True(t, r.IsActive)
		}
	})

	t.Run("Update Router Device", func(t *testing.T) {
		defer suite.Cleanup(t)

		router := &model.MikrotikRouter{
			ID:                uuid.New().String(),
			Name:              "Original Router Name",
			Address:           "192.168.88.50",
			APIPort:           8728,
			Username:          "admin",
			PasswordEncrypted: "encrypted",
			IsActive:          true,
			Status:            "online",
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		}

		err := repo.Create(suite.Ctx, router)
		require.NoError(t, err)

		// Update router
		router.Name = "Updated Router Name"
		router.Status = "offline"
		router.APIPort = 8729

		err = repo.Update(suite.Ctx, router)
		require.NoError(t, err)

		// Verify update
		id, _ := uuid.Parse(router.ID)
		fetched, err := repo.GetByID(suite.Ctx, id)
		require.NoError(t, err)
		assert.Equal(t, "Updated Router Name", fetched.Name)
		assert.Equal(t, "offline", fetched.Status)
		assert.Equal(t, 8729, fetched.APIPort)
	})

	t.Run("UpdateLastSync", func(t *testing.T) {
		defer suite.Cleanup(t)

		router := &model.MikrotikRouter{
			ID:                uuid.New().String(),
			Name:              "Sync Test Router",
			Address:           "192.168.88.60",
			APIPort:           8728,
			Username:          "admin",
			PasswordEncrypted: "encrypted",
			IsActive:          true,
			Status:            "online",
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		}

		err := repo.Create(suite.Ctx, router)
		require.NoError(t, err)

		// Update last sync
		id, _ := uuid.Parse(router.ID)
		err = repo.UpdateLastSync(suite.Ctx, id)
		require.NoError(t, err)

		// Verify update
		fetched, err := repo.GetByID(suite.Ctx, id)
		require.NoError(t, err)
		assert.NotNil(t, fetched.LastSeenAt)
		assert.WithinDuration(t, time.Now(), *fetched.LastSeenAt, 5*time.Second)
	})

	t.Run("Delete Router Device", func(t *testing.T) {
		defer suite.Cleanup(t)

		router := &model.MikrotikRouter{
			ID:                uuid.New().String(),
			Name:              "To Be Deleted",
			Address:           "192.168.88.70",
			APIPort:           8728,
			Username:          "admin",
			PasswordEncrypted: "encrypted",
			IsActive:          true,
			Status:            "online",
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		}

		err := repo.Create(suite.Ctx, router)
		require.NoError(t, err)

		// Delete
		id, _ := uuid.Parse(router.ID)
		err = repo.Delete(suite.Ctx, id)
		require.NoError(t, err)

		// Verify deletion
		_, err = repo.GetByID(suite.Ctx, id)
		assert.Error(t, err)
	})

	t.Run("List Router Devices", func(t *testing.T) {
		defer suite.Cleanup(t)

		// Create multiple routers
		for i := 0; i < 5; i++ {
			router := &model.MikrotikRouter{
				ID:                uuid.New().String(),
				Name:              "List Router " + string(rune('A'+i)),
				Address:           "10.0.0." + string(rune('1'+i)),
				APIPort:           8728,
				Username:          "admin",
				PasswordEncrypted: "encrypted",
				IsActive:          true,
				Status:            "online",
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
			}
			err := repo.Create(suite.Ctx, router)
			require.NoError(t, err)
		}

		// Test List with pagination
		routers, err := repo.List(suite.Ctx, 3, 0)
		require.NoError(t, err)
		assert.Len(t, routers, 3)

		// Test List with offset
		routers, err = repo.List(suite.Ctx, 3, 3)
		require.NoError(t, err)
		assert.Len(t, routers, 2)
	})
}
