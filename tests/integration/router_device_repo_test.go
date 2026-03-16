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
	t.Run("Create and Get Router Device", func(t *testing.T) {
		suite := SetupSuite(t)
		defer suite.TearDownSuite(t)
		defer suite.Cleanup(t)

		repo := postgres.NewRouterDeviceRepository(suite.DB)

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

		require.NoError(t, repo.Create(suite.Ctx, router))

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
		suite := SetupSuite(t)
		defer suite.TearDownSuite(t)
		defer suite.Cleanup(t)

		repo := postgres.NewRouterDeviceRepository(suite.DB)

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
			require.NoError(t, repo.Create(suite.Ctx, router))
		}

		inactiveRouter := &model.MikrotikRouter{
			ID:                uuid.New().String(),
			Name:              "Inactive Router",
			Address:           "192.168.88.100",
			APIPort:           8728,
			Username:          "admin",
			PasswordEncrypted: "encrypted",
			IsActive:          true, // create active, then force inactive via SQL
			Status:            "offline",
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		}
		require.NoError(t, repo.Create(suite.Ctx, inactiveRouter))
		// GORM skips false zero-values on Create; force is_active=false via SQL.
		require.NoError(t, suite.DB.WithContext(suite.Ctx).
			Exec("UPDATE mikrotik_routers SET is_active=false WHERE id=?", inactiveRouter.ID).Error)

		activeRouters, err := repo.GetActive(suite.Ctx)
		require.NoError(t, err)
		assert.Len(t, activeRouters, 3)
		for _, r := range activeRouters {
			assert.True(t, r.IsActive)
		}
	})

	t.Run("Update Router Device", func(t *testing.T) {
		suite := SetupSuite(t)
		defer suite.TearDownSuite(t)
		defer suite.Cleanup(t)

		repo := postgres.NewRouterDeviceRepository(suite.DB)

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

		require.NoError(t, repo.Create(suite.Ctx, router))

		router.Name = "Updated Router Name"
		router.Status = "offline"
		router.APIPort = 8729

		require.NoError(t, repo.Update(suite.Ctx, router))

		id, _ := uuid.Parse(router.ID)
		fetched, err := repo.GetByID(suite.Ctx, id)
		require.NoError(t, err)
		assert.Equal(t, "Updated Router Name", fetched.Name)
		assert.Equal(t, "offline", fetched.Status)
		assert.Equal(t, 8729, fetched.APIPort)
	})

	t.Run("UpdateLastSync", func(t *testing.T) {
		suite := SetupSuite(t)
		defer suite.TearDownSuite(t)
		defer suite.Cleanup(t)

		repo := postgres.NewRouterDeviceRepository(suite.DB)

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

		require.NoError(t, repo.Create(suite.Ctx, router))

		id, _ := uuid.Parse(router.ID)
		require.NoError(t, repo.UpdateLastSync(suite.Ctx, id))

		fetched, err := repo.GetByID(suite.Ctx, id)
		require.NoError(t, err)
		assert.NotNil(t, fetched.LastSeenAt)
		assert.WithinDuration(t, time.Now(), *fetched.LastSeenAt, 5*time.Second)
	})

	t.Run("Delete Router Device", func(t *testing.T) {
		suite := SetupSuite(t)
		defer suite.TearDownSuite(t)
		defer suite.Cleanup(t)

		repo := postgres.NewRouterDeviceRepository(suite.DB)

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

		require.NoError(t, repo.Create(suite.Ctx, router))

		id, _ := uuid.Parse(router.ID)
		require.NoError(t, repo.Delete(suite.Ctx, id))

		_, err := repo.GetByID(suite.Ctx, id)
		assert.Error(t, err)
	})

	t.Run("List Router Devices", func(t *testing.T) {
		suite := SetupSuite(t)
		defer suite.TearDownSuite(t)
		defer suite.Cleanup(t)

		repo := postgres.NewRouterDeviceRepository(suite.DB)

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
			require.NoError(t, repo.Create(suite.Ctx, router))
		}

		routers, err := repo.List(suite.Ctx, 3, 0)
		require.NoError(t, err)
		assert.Len(t, routers, 3)

		routers, err = repo.List(suite.Ctx, 3, 3)
		require.NoError(t, err)
		assert.Len(t, routers, 2)
	})
}
