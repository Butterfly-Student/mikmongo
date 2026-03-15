package customer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"mikmongo/internal/model"
)

func TestValidateCustomer(t *testing.T) {
	d := NewDomain()

	t.Run("valid customer", func(t *testing.T) {
		c := &model.Customer{FullName: "Budi Santoso", Phone: "081234567890"}
		assert.NoError(t, d.ValidateCustomer(c))
	})

	t.Run("full name empty → error", func(t *testing.T) {
		c := &model.Customer{FullName: "", Phone: "081234567890"}
		assert.Error(t, d.ValidateCustomer(c))
	})

	t.Run("full name whitespace only → error", func(t *testing.T) {
		c := &model.Customer{FullName: "   ", Phone: "081234567890"}
		assert.Error(t, d.ValidateCustomer(c))
	})

	t.Run("phone empty → error", func(t *testing.T) {
		c := &model.Customer{FullName: "Budi Santoso", Phone: ""}
		assert.Error(t, d.ValidateCustomer(c))
	})

	t.Run("phone whitespace only → error", func(t *testing.T) {
		c := &model.Customer{FullName: "Budi Santoso", Phone: "   "}
		assert.Error(t, d.ValidateCustomer(c))
	})

	t.Run("both empty → error", func(t *testing.T) {
		c := &model.Customer{FullName: "", Phone: ""}
		assert.Error(t, d.ValidateCustomer(c))
	})
}

func TestCanDeactivate(t *testing.T) {
	d := NewDomain()

	t.Run("active customer → can deactivate", func(t *testing.T) {
		c := &model.Customer{IsActive: true}
		assert.NoError(t, d.CanDeactivate(c))
	})

	t.Run("already inactive → error", func(t *testing.T) {
		c := &model.Customer{IsActive: false}
		assert.Error(t, d.CanDeactivate(c))
	})
}

func TestCanActivate(t *testing.T) {
	d := NewDomain()

	t.Run("inactive customer → can activate", func(t *testing.T) {
		c := &model.Customer{IsActive: false}
		assert.NoError(t, d.CanActivate(c))
	})

	t.Run("already active → error", func(t *testing.T) {
		c := &model.Customer{IsActive: true}
		assert.Error(t, d.CanActivate(c))
	})
}
