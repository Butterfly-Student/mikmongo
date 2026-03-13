// Package registration contains customer registration domain logic
package registration

import (
	"errors"
	"mikmongo/internal/model"
	"strings"
)

// Domain represents registration business logic
type Domain struct{}

// NewDomain creates a new registration domain
func NewDomain() *Domain {
	return &Domain{}
}

// ValidateRegistration validates required fields on a CustomerRegistration
func (d *Domain) ValidateRegistration(reg *model.CustomerRegistration) error {
	if strings.TrimSpace(reg.FullName) == "" {
		return errors.New("full name is required")
	}
	if strings.TrimSpace(reg.Phone) == "" {
		return errors.New("phone is required")
	}
	return nil
}

// CanApprove checks if registration can be approved (must be pending)
func (d *Domain) CanApprove(reg *model.CustomerRegistration) error {
	if reg.Status != "pending" {
		return errors.New("only pending registrations can be approved")
	}
	return nil
}

// CanReject checks if registration can be rejected (must be pending)
func (d *Domain) CanReject(reg *model.CustomerRegistration) error {
	if reg.Status != "pending" {
		return errors.New("only pending registrations can be rejected")
	}
	return nil
}

// IsAlreadyConverted returns true if the registration has been converted to a customer
func (d *Domain) IsAlreadyConverted(reg *model.CustomerRegistration) bool {
	return reg.CustomerID != nil
}

// NeedsProfileAssignment returns true if no bandwidth profile has been assigned
func (d *Domain) NeedsProfileAssignment(reg *model.CustomerRegistration) bool {
	return reg.BandwidthProfileID == nil
}
