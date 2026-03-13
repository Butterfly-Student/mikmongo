// Package customer contains customer domain logic
package customer

import (
	"errors"
	"mikmongo/internal/model"
	"strings"
)

// Domain represents customer business logic
type Domain struct{}

// NewDomain creates a new customer domain
func NewDomain() *Domain {
	return &Domain{}
}

// ValidateCustomer validates required customer fields
func (d *Domain) ValidateCustomer(c *model.Customer) error {
	if strings.TrimSpace(c.FullName) == "" {
		return errors.New("full name is required")
	}
	if strings.TrimSpace(c.Phone) == "" {
		return errors.New("phone is required")
	}
	return nil
}

// CanDeactivate checks if customer account can be deactivated
func (d *Domain) CanDeactivate(c *model.Customer) error {
	if !c.IsActive {
		return errors.New("customer is already deactivated")
	}
	return nil
}

// CanActivate checks if customer account can be activated
func (d *Domain) CanActivate(c *model.Customer) error {
	if c.IsActive {
		return errors.New("customer is already active")
	}
	return nil
}
