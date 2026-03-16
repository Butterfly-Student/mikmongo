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

func TestBandwidthProfileRepository_Integration(t *testing.T) {
	t.Run("Create and Get Bandwidth Profile", func(t *testing.T) {
		suite := SetupSuite(t)
		defer suite.TearDownSuite(t)
		defer suite.Cleanup(t)

		repo := postgres.NewBandwidthProfileRepository(suite.DB)
		routerRepo := postgres.NewRouterDeviceRepository(suite.DB)

		router := &model.MikrotikRouter{
			ID:                uuid.New().String(),
			Name:              "Test Router",
			Address:           "192.168.88.1",
			APIPort:           8728,
			Username:          "admin",
			PasswordEncrypted: "encrypted_password",
			IsActive:          true,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		}
		require.NoError(t, routerRepo.Create(suite.Ctx, router))

		profile := &model.BandwidthProfile{
			ID:              uuid.New().String(),
			RouterID:        router.ID,
			ProfileCode:     "BASIC10",
			Name:            "Basic 10Mbps",
			Description:     strPtr("Basic internet package"),
			DownloadSpeed:   10000,
			UploadSpeed:     10000,
			PriceMonthly:    100000,
			TaxRate:         0.11,
			BillingCycle:    "monthly",
			IsActive:        true,
			IsVisible:       true,
			SortOrder:       1,
			GracePeriodDays: 3,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}

		require.NoError(t, repo.Create(suite.Ctx, profile))

		id, err := uuid.Parse(profile.ID)
		require.NoError(t, err)

		fetched, err := repo.GetByID(suite.Ctx, id)
		require.NoError(t, err)
		assert.Equal(t, profile.ProfileCode, fetched.ProfileCode)
		assert.Equal(t, profile.Name, fetched.Name)
		assert.Equal(t, profile.DownloadSpeed, fetched.DownloadSpeed)
		assert.Equal(t, profile.UploadSpeed, fetched.UploadSpeed)
		assert.Equal(t, profile.PriceMonthly, fetched.PriceMonthly)
	})

	t.Run("GetByCode", func(t *testing.T) {
		suite := SetupSuite(t)
		defer suite.TearDownSuite(t)
		defer suite.Cleanup(t)

		repo := postgres.NewBandwidthProfileRepository(suite.DB)
		routerRepo := postgres.NewRouterDeviceRepository(suite.DB)

		router := &model.MikrotikRouter{
			ID:                uuid.New().String(),
			Name:              "Test Router",
			Address:           "192.168.88.1",
			APIPort:           8728,
			Username:          "admin",
			PasswordEncrypted: "encrypted_password",
			IsActive:          true,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		}
		require.NoError(t, routerRepo.Create(suite.Ctx, router))

		profile := &model.BandwidthProfile{
			ID:            uuid.New().String(),
			RouterID:      router.ID,
			ProfileCode:   "PREMIUM50",
			Name:          "Premium 50Mbps",
			DownloadSpeed: 50000,
			UploadSpeed:   50000,
			PriceMonthly:  500000,
			IsActive:      true,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		require.NoError(t, repo.Create(suite.Ctx, profile))

		fetched, err := repo.GetByCode(suite.Ctx, "PREMIUM50")
		require.NoError(t, err)
		assert.Equal(t, profile.Name, fetched.Name)
	})

	t.Run("Update Bandwidth Profile", func(t *testing.T) {
		suite := SetupSuite(t)
		defer suite.TearDownSuite(t)
		defer suite.Cleanup(t)

		repo := postgres.NewBandwidthProfileRepository(suite.DB)
		routerRepo := postgres.NewRouterDeviceRepository(suite.DB)

		router := &model.MikrotikRouter{
			ID:                uuid.New().String(),
			Name:              "Test Router",
			Address:           "192.168.88.1",
			APIPort:           8728,
			Username:          "admin",
			PasswordEncrypted: "encrypted_password",
			IsActive:          true,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		}
		require.NoError(t, routerRepo.Create(suite.Ctx, router))

		profile := &model.BandwidthProfile{
			ID:            uuid.New().String(),
			RouterID:      router.ID,
			ProfileCode:   "UPDATE_TEST",
			Name:          "Original Name",
			DownloadSpeed: 20000,
			UploadSpeed:   20000,
			PriceMonthly:  200000,
			IsActive:      true,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
		require.NoError(t, repo.Create(suite.Ctx, profile))

		profile.Name = "Updated Name"
		profile.PriceMonthly = 250000
		profile.IsActive = false

		require.NoError(t, repo.Update(suite.Ctx, profile))

		id, _ := uuid.Parse(profile.ID)
		fetched, err := repo.GetByID(suite.Ctx, id)
		require.NoError(t, err)
		assert.Equal(t, "Updated Name", fetched.Name)
		assert.Equal(t, float64(250000), fetched.PriceMonthly)
		assert.False(t, fetched.IsActive)
	})

	t.Run("Delete Bandwidth Profile", func(t *testing.T) {
		suite := SetupSuite(t)
		defer suite.TearDownSuite(t)
		defer suite.Cleanup(t)

		repo := postgres.NewBandwidthProfileRepository(suite.DB)
		routerRepo := postgres.NewRouterDeviceRepository(suite.DB)

		router := &model.MikrotikRouter{
			ID:                uuid.New().String(),
			Name:              "Test Router",
			Address:           "192.168.88.1",
			APIPort:           8728,
			Username:          "admin",
			PasswordEncrypted: "encrypted_password",
			IsActive:          true,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		}
		require.NoError(t, routerRepo.Create(suite.Ctx, router))

		profile := &model.BandwidthProfile{
			ID:            uuid.New().String(),
			RouterID:      router.ID,
			ProfileCode:   "DELETE_TEST",
			Name:          "To Be Deleted",
			DownloadSpeed: 10000,
			UploadSpeed:   10000,
			PriceMonthly:  100000,
			IsActive:      true,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
		require.NoError(t, repo.Create(suite.Ctx, profile))

		id, _ := uuid.Parse(profile.ID)
		require.NoError(t, repo.Delete(suite.Ctx, id))

		_, err := repo.GetByID(suite.Ctx, id)
		assert.Error(t, err)
	})

	t.Run("ListByRouterID", func(t *testing.T) {
		suite := SetupSuite(t)
		defer suite.TearDownSuite(t)
		defer suite.Cleanup(t)

		repo := postgres.NewBandwidthProfileRepository(suite.DB)
		routerRepo := postgres.NewRouterDeviceRepository(suite.DB)

		router := &model.MikrotikRouter{
			ID:                uuid.New().String(),
			Name:              "Test Router",
			Address:           "192.168.88.1",
			APIPort:           8728,
			Username:          "admin",
			PasswordEncrypted: "encrypted_password",
			IsActive:          true,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		}
		require.NoError(t, routerRepo.Create(suite.Ctx, router))

		for i := 0; i < 3; i++ {
			profile := &model.BandwidthProfile{
				ID:            uuid.New().String(),
				RouterID:      router.ID,
				ProfileCode:   "ROUTER_" + string(rune('A'+i)),
				Name:          "Profile " + string(rune('A'+i)),
				DownloadSpeed: int64(10000 * (i + 1)),
				UploadSpeed:   int64(10000 * (i + 1)),
				PriceMonthly:  float64(100000 * (i + 1)),
				IsActive:      true,
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			}
			require.NoError(t, repo.Create(suite.Ctx, profile))
		}

		routerID, _ := uuid.Parse(router.ID)
		profiles, err := repo.ListByRouterID(suite.Ctx, routerID, 10, 0)
		require.NoError(t, err)
		assert.Len(t, profiles, 3)
	})

	t.Run("ListActive", func(t *testing.T) {
		suite := SetupSuite(t)
		defer suite.TearDownSuite(t)
		defer suite.Cleanup(t)

		repo := postgres.NewBandwidthProfileRepository(suite.DB)
		routerRepo := postgres.NewRouterDeviceRepository(suite.DB)

		router := &model.MikrotikRouter{
			ID:                uuid.New().String(),
			Name:              "Test Router",
			Address:           "192.168.88.1",
			APIPort:           8728,
			Username:          "admin",
			PasswordEncrypted: "encrypted_password",
			IsActive:          true,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		}
		require.NoError(t, routerRepo.Create(suite.Ctx, router))

		activeProfile := &model.BandwidthProfile{
			ID:            uuid.New().String(),
			RouterID:      router.ID,
			ProfileCode:   "ACTIVE",
			Name:          "Active Profile",
			DownloadSpeed: 10000,
			UploadSpeed:   10000,
			PriceMonthly:  100000,
			IsActive:      true,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
		require.NoError(t, repo.Create(suite.Ctx, activeProfile))

		inactiveProfile := &model.BandwidthProfile{
			ID:            uuid.New().String(),
			RouterID:      router.ID,
			ProfileCode:   "INACTIVE",
			Name:          "Inactive Profile",
			DownloadSpeed: 20000,
			UploadSpeed:   20000,
			PriceMonthly:  200000,
			IsActive:      true, // created active, then forced inactive below
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
		require.NoError(t, repo.Create(suite.Ctx, inactiveProfile))
		// GORM skips false zero-values on Create; force is_active=false via SQL.
		require.NoError(t, suite.DB.WithContext(suite.Ctx).
			Exec("UPDATE bandwidth_profiles SET is_active=false WHERE id=?", inactiveProfile.ID).Error)

		profiles, err := repo.ListActive(suite.Ctx)
		require.NoError(t, err)
		assert.Len(t, profiles, 1)
		assert.Equal(t, "Active Profile", profiles[0].Name)
	})
}
