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

func TestCustomerRepository_Integration(t *testing.T) {
	t.Run("Create and Get Customer", func(t *testing.T) {
		suite := SetupSuite(t)
		defer suite.TearDownSuite(t)
		defer suite.Cleanup(t)

		repo := postgres.NewCustomerRepository(suite.DB)

		customer := &model.Customer{
			ID:           uuid.New().String(),
			CustomerCode: "CUST001",
			FullName:     "John Doe",
			Email:        strPtr("john@example.com"),
			Phone:        "08123456789",
			IsActive:     true,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		require.NoError(t, repo.Create(suite.Ctx, customer))

		id, err := uuid.Parse(customer.ID)
		require.NoError(t, err)

		fetched, err := repo.GetByID(suite.Ctx, id)
		require.NoError(t, err)
		assert.Equal(t, customer.CustomerCode, fetched.CustomerCode)
		assert.Equal(t, customer.FullName, fetched.FullName)
		assert.Equal(t, *customer.Email, *fetched.Email)
		assert.Equal(t, customer.Phone, fetched.Phone)
		assert.True(t, fetched.IsActive)
	})

	t.Run("GetByEmail", func(t *testing.T) {
		suite := SetupSuite(t)
		defer suite.TearDownSuite(t)
		defer suite.Cleanup(t)

		repo := postgres.NewCustomerRepository(suite.DB)

		customer := &model.Customer{
			ID:           uuid.New().String(),
			CustomerCode: "CUST002",
			FullName:     "Jane Doe",
			Email:        strPtr("jane@example.com"),
			Phone:        "08123456790",
			IsActive:     true,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		require.NoError(t, repo.Create(suite.Ctx, customer))

		fetched, err := repo.GetByEmail(suite.Ctx, "jane@example.com")
		require.NoError(t, err)
		assert.Equal(t, customer.FullName, fetched.FullName)
	})

	t.Run("Update Customer", func(t *testing.T) {
		suite := SetupSuite(t)
		defer suite.TearDownSuite(t)
		defer suite.Cleanup(t)

		repo := postgres.NewCustomerRepository(suite.DB)

		customer := &model.Customer{
			ID:           uuid.New().String(),
			CustomerCode: "CUST003",
			FullName:     "Original Name",
			Phone:        "08123456791",
			IsActive:     true,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		require.NoError(t, repo.Create(suite.Ctx, customer))

		customer.FullName = "Updated Name"
		customer.IsActive = false

		require.NoError(t, repo.Update(suite.Ctx, customer))

		id, _ := uuid.Parse(customer.ID)
		fetched, err := repo.GetByID(suite.Ctx, id)
		require.NoError(t, err)
		assert.Equal(t, "Updated Name", fetched.FullName)
		assert.False(t, fetched.IsActive)
	})

	t.Run("Delete Customer", func(t *testing.T) {
		suite := SetupSuite(t)
		defer suite.TearDownSuite(t)
		defer suite.Cleanup(t)

		repo := postgres.NewCustomerRepository(suite.DB)

		customer := &model.Customer{
			ID:           uuid.New().String(),
			CustomerCode: "CUST004",
			FullName:     "To Be Deleted",
			Phone:        "08123456792",
			IsActive:     true,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		require.NoError(t, repo.Create(suite.Ctx, customer))

		id, _ := uuid.Parse(customer.ID)
		require.NoError(t, repo.Delete(suite.Ctx, id))

		_, err := repo.GetByID(suite.Ctx, id)
		assert.Error(t, err)
	})

	t.Run("List and Count Customers", func(t *testing.T) {
		suite := SetupSuite(t)
		defer suite.TearDownSuite(t)
		defer suite.Cleanup(t)

		repo := postgres.NewCustomerRepository(suite.DB)

		// Record count before to handle pre-existing committed records in DB.
		beforeCount, err := repo.Count(suite.Ctx)
		require.NoError(t, err)

		for i := 0; i < 5; i++ {
			customer := &model.Customer{
				ID:           uuid.New().String(),
				CustomerCode: "CUST" + string(rune('A'+i)),
				FullName:     "Customer " + string(rune('A'+i)),
				Phone:        "0812345679" + string(rune('0'+i)),
				IsActive:     true,
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			}
			require.NoError(t, repo.Create(suite.Ctx, customer))
		}

		afterCount, err := repo.Count(suite.Ctx)
		require.NoError(t, err)
		assert.Equal(t, beforeCount+5, afterCount)

		// Limit=3 should return exactly 3 (total > 3).
		customers, err := repo.List(suite.Ctx, 3, 0)
		require.NoError(t, err)
		assert.Len(t, customers, 3)

		// offset=beforeCount+3 skips pre-existing + first 3 of ours → 2 remain.
		customers, err = repo.List(suite.Ctx, 3, int(beforeCount)+3)
		require.NoError(t, err)
		assert.Len(t, customers, 2)
	})

	t.Run("Customer IsActive Filter", func(t *testing.T) {
		suite := SetupSuite(t)
		defer suite.TearDownSuite(t)
		defer suite.Cleanup(t)

		repo := postgres.NewCustomerRepository(suite.DB)

		activeCustomer := &model.Customer{
			ID:           uuid.New().String(),
			CustomerCode: "ACTIVE001",
			FullName:     "Active Customer",
			Phone:        "08123456793",
			IsActive:     true,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		require.NoError(t, repo.Create(suite.Ctx, activeCustomer))

		inactiveCustomer := &model.Customer{
			ID:           uuid.New().String(),
			CustomerCode: "INACTIVE001",
			FullName:     "Inactive Customer",
			Phone:        "08123456794",
			IsActive:     true, // create active, then force inactive via SQL
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		require.NoError(t, repo.Create(suite.Ctx, inactiveCustomer))
		// GORM skips false zero-values on Create; force is_active=false via SQL.
		require.NoError(t, suite.DB.WithContext(suite.Ctx).
			Exec("UPDATE customers SET is_active=false WHERE id=?", inactiveCustomer.ID).Error)

		// repo.List returns all customers; verify both our customers are present.
		customers, err := repo.List(suite.Ctx, 100, 0)
		require.NoError(t, err)
		found := map[string]bool{}
		for _, c := range customers {
			if c.ID == activeCustomer.ID || c.ID == inactiveCustomer.ID {
				found[c.ID] = true
			}
		}
		assert.True(t, found[activeCustomer.ID], "active customer should be in list")
		assert.True(t, found[inactiveCustomer.ID], "inactive customer should be in list")
	})
}

func strPtr(s string) *string {
	return &s
}
