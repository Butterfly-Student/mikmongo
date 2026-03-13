// Package router contains router domain logic
package router

import (
	"errors"
	"mikmongo/internal/model"
	"time"
)

// Domain represents router business logic
type Domain struct{}

// NewDomain creates a new router domain
func NewDomain() *Domain {
	return &Domain{}
}

// ValidateConnection validates router host and port
func (d *Domain) ValidateConnection(host string, port int) error {
	if host == "" {
		return errors.New("router host is required")
	}
	if port < 1 || port > 65535 {
		return errors.New("router port must be between 1 and 65535")
	}
	return nil
}

// IsOnline returns true if router status is "online"
func (d *Domain) IsOnline(r *model.MikrotikRouter) bool {
	return r.Status == "online"
}

// CanConnect returns true if router is active and not offline
func (d *Domain) CanConnect(r *model.MikrotikRouter) bool {
	return r.IsActive && r.Status != "offline"
}

// ShouldSync returns true if LastSeenAt is nil or exceeded intervalMinutes ago
func (d *Domain) ShouldSync(lastSeenAt *time.Time, intervalMinutes int) bool {
	if lastSeenAt == nil {
		return true
	}
	return time.Since(*lastSeenAt) > time.Duration(intervalMinutes)*time.Minute
}

// IsStale returns true if router has not been seen within staleMinutes
func (d *Domain) IsStale(r *model.MikrotikRouter, staleMinutes int) bool {
	if r.LastSeenAt == nil {
		return true
	}
	return time.Since(*r.LastSeenAt) > time.Duration(staleMinutes)*time.Minute
}
