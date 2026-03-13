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
	// Setup test suite
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)

	// Create repository
	repo := postgres.NewCustomerRepository(suite.DB)

	t.Run("Create and Get Customer", func(t *testing.T) {
		defer suite.Cleanup(t)

		// Create test customer
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

		// Test Create
		err := repo.Create(suite.Ctx, customer)
		require.NoError(t, err)

		// Test GetByID
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
		defer suite.Cleanup(t)

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

		err := repo.Create(suite.Ctx, customer)
		require.NoError(t, err)

		fetched, err := repo.GetByEmail(suite.Ctx, "jane@example.com")
		require.NoError(t, err)
		assert.Equal(t, customer.FullName, fetched.FullName)
	})

	t.Run("Update Customer", func(t *testing.T) {
		defer suite.Cleanup(t)

		customer := &model.Customer{
			ID:           uuid.New().String(),
			CustomerCode: "CUST003",
			FullName:     "Original Name",
			Phone:        "08123456791",
			IsActive:     true,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		err := repo.Create(suite.Ctx, customer)
		require.NoError(t, err)

		// Update customer
		customer.FullName = "Updated Name"
		customer.IsActive = false

		err = repo.Update(suite.Ctx, customer)
		require.NoError(t, err)

		// Verify update
		id, _ := uuid.Parse(customer.ID)
		fetched, err := repo.GetByID(suite.Ctx, id)
		require.NoError(t, err)
		assert.Equal(t, "Updated Name", fetched.FullName)
		assert.False(t, fetched.IsActive)
	})

	t.Run("Delete Customer", func(t *testing.T) {
		defer suite.Cleanup(t)

		customer := &model.Customer{
			ID:           uuid.New().String(),
			CustomerCode: "CUST004",
			FullName:     "To Be Deleted",
			Phone:        "08123456792",
			IsActive:     true,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		err := repo.Create(suite.Ctx, customer)
		require.NoError(t, err)

		// Delete
		id, _ := uuid.Parse(customer.ID)
		err = repo.Delete(suite.Ctx, id)
		require.NoError(t, err)

		// Verify deletion (should return error)
		_, err = repo.GetByID(suite.Ctx, id)
		assert.Error(t, err)
	})

	t.Run("List and Count Customers", func(t *testing.T) {
		defer suite.Cleanup(t)

		// Create multiple customers
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
			err := repo.Create(suite.Ctx, customer)
			require.NoError(t, err)
		}

		// Test Count
		count, err := repo.Count(suite.Ctx)
		require.NoError(t, err)
		assert.Equal(t, int64(5), count)

		// Test List with pagination
		customers, err := repo.List(suite.Ctx, 3, 0)
		require.NoError(t, err)
		assert.Len(t, customers, 3)

		// Test List with offset
		customers, err = repo.List(suite.Ctx, 3, 3)
		require.NoError(t, err)
		assert.Len(t, customers, 2)
	})

	t.Run("Customer IsActive Filter", func(t *testing.T) {
		defer suite.Cleanup(t)

		// Create active customer
		activeCustomer := &model.Customer{
			ID:           uuid.New().String(),
			CustomerCode: "ACTIVE001",
			FullName:     "Active Customer",
			Phone:        "08123456793",
			IsActive:     true,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		err := repo.Create(suite.Ctx, activeCustomer)
		require.NoError(t, err)

		// Create inactive customer
		inactiveCustomer := &model.Customer{
			ID:           uuid.New().String(),
			CustomerCode: "INACTIVE001",
			FullName:     "Inactive Customer",
			Phone:        "08123456794",
			IsActive:     false,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		err = repo.Create(suite.Ctx, inactiveCustomer)
		require.NoError(t, err)

		// List all customers
		customers, err := repo.List(suite.Ctx, 10, 0)
		require.NoError(t, err)
		assert.Len(t, customers, 2)
	})
}

func strPtr(s string) *string {
	return &s
}
