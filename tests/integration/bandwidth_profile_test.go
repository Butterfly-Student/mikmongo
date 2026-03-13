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
	// Setup test suite
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)

	// Create repositories
	repo := postgres.NewBandwidthProfileRepository(suite.DB)
	routerRepo := postgres.NewRouterDeviceRepository(suite.DB)

	// Helper function to create a test router
	createTestRouter := func(t *testing.T) *model.MikrotikRouter {
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
		err := routerRepo.Create(suite.Ctx, router)
		require.NoError(t, err)
		return router
	}

	t.Run("Create and Get Bandwidth Profile", func(t *testing.T) {
		defer suite.Cleanup(t)

		// Create test router first
		router := createTestRouter(t)

		// Create test bandwidth profile
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

		// Test Create
		err := repo.Create(suite.Ctx, profile)
		require.NoError(t, err)

		// Test GetByID
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
		defer suite.Cleanup(t)

		router := createTestRouter(t)

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

		err := repo.Create(suite.Ctx, profile)
		require.NoError(t, err)

		fetched, err := repo.GetByCode(suite.Ctx, "PREMIUM50")
		require.NoError(t, err)
		assert.Equal(t, profile.Name, fetched.Name)
	})

	t.Run("Update Bandwidth Profile", func(t *testing.T) {
		defer suite.Cleanup(t)

		router := createTestRouter(t)

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

		err := repo.Create(suite.Ctx, profile)
		require.NoError(t, err)

		// Update profile
		profile.Name = "Updated Name"
		profile.PriceMonthly = 250000
		profile.IsActive = false

		err = repo.Update(suite.Ctx, profile)
		require.NoError(t, err)

		// Verify update
		id, _ := uuid.Parse(profile.ID)
		fetched, err := repo.GetByID(suite.Ctx, id)
		require.NoError(t, err)
		assert.Equal(t, "Updated Name", fetched.Name)
		assert.Equal(t, float64(250000), fetched.PriceMonthly)
		assert.False(t, fetched.IsActive)
	})

	t.Run("Delete Bandwidth Profile", func(t *testing.T) {
		defer suite.Cleanup(t)

		router := createTestRouter(t)

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

		err := repo.Create(suite.Ctx, profile)
		require.NoError(t, err)

		// Delete
		id, _ := uuid.Parse(profile.ID)
		err = repo.Delete(suite.Ctx, id)
		require.NoError(t, err)

		// Verify deletion (should return error)
		_, err = repo.GetByID(suite.Ctx, id)
		assert.Error(t, err)
	})

	t.Run("ListByRouterID", func(t *testing.T) {
		defer suite.Cleanup(t)

		router := createTestRouter(t)

		// Create multiple profiles for the same router
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
			err := repo.Create(suite.Ctx, profile)
			require.NoError(t, err)
		}

		// Test ListByRouterID
		routerID, _ := uuid.Parse(router.ID)
		profiles, err := repo.ListByRouterID(suite.Ctx, routerID, 10, 0)
		require.NoError(t, err)
		assert.Len(t, profiles, 3)
	})

	t.Run("ListActive", func(t *testing.T) {
		defer suite.Cleanup(t)

		router := createTestRouter(t)

		// Create active and inactive profiles
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
		err := repo.Create(suite.Ctx, activeProfile)
		require.NoError(t, err)

		inactiveProfile := &model.BandwidthProfile{
			ID:            uuid.New().String(),
			RouterID:      router.ID,
			ProfileCode:   "INACTIVE",
			Name:          "Inactive Profile",
			DownloadSpeed: 20000,
			UploadSpeed:   20000,
			PriceMonthly:  200000,
			IsActive:      false,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
		err = repo.Create(suite.Ctx, inactiveProfile)
		require.NoError(t, err)

		// Test ListActive
		profiles, err := repo.ListActive(suite.Ctx)
		require.NoError(t, err)
		assert.Len(t, profiles, 1)
		assert.Equal(t, "Active Profile", profiles[0].Name)
	})
}
