// Package subscription contains subscription domain logic
package subscription

import (
	"crypto/rand"
	"errors"
	"math/big"
	"mikmongo/internal/model"
	"time"
)

const alphanumeric = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// validTransitions defines allowed status transitions
var validTransitions = map[string][]string{
	"pending":    {"active", "terminated"},
	"active":     {"suspended", "isolated", "expired", "terminated"},
	"suspended":  {"active", "isolated", "terminated"},
	"isolated":   {"active", "suspended", "terminated"},
	"expired":    {"active", "terminated"},
	"terminated": {},
}

// Domain represents subscription business logic
type Domain struct{}

// NewDomain creates a new subscription domain
func NewDomain() *Domain {
	return &Domain{}
}

// ValidateStatusTransition checks if transitioning from current to next is allowed
func (d *Domain) ValidateStatusTransition(current, next string) error {
	allowed, ok := validTransitions[current]
	if !ok {
		return errors.New("unknown status: " + current)
	}
	for _, s := range allowed {
		if s == next {
			return nil
		}
	}
	return errors.New("cannot transition subscription from " + current + " to " + next)
}

// CanActivate checks if subscription can be activated
func (d *Domain) CanActivate(sub *model.Subscription) error {
	if sub.Status == "pending" || sub.Status == "suspended" {
		return nil
	}
	return errors.New("subscription can only be activated from pending or suspended status")
}

// CanSuspend checks if subscription can be suspended
func (d *Domain) CanSuspend(sub *model.Subscription) error {
	if sub.Status == "active" || sub.Status == "isolated" {
		return nil
	}
	return errors.New("subscription can only be suspended from active or isolated status")
}

// CanIsolate checks if subscription can be isolated
func (d *Domain) CanIsolate(sub *model.Subscription) error {
	if sub.Status == "active" || sub.Status == "suspended" {
		return nil
	}
	return errors.New("subscription can only be isolated from active or suspended status")
}

// CanRestore checks if subscription can be restored from isolated
func (d *Domain) CanRestore(sub *model.Subscription) error {
	if sub.Status == "isolated" {
		return nil
	}
	return errors.New("subscription can only be restored from isolated status")
}

// CanTerminate checks if subscription can be terminated
func (d *Domain) CanTerminate(sub *model.Subscription) error {
	if sub.Status == "terminated" {
		return errors.New("subscription is already terminated")
	}
	return nil
}

// IsExpired returns true if ExpiryDate is set and now is after it
func (d *Domain) IsExpired(sub *model.Subscription, now time.Time) bool {
	return sub.ExpiryDate != nil && now.After(*sub.ExpiryDate)
}

// NeedsSync always returns false - all operations now sync automatically
func (d *Domain) NeedsSync(sub *model.Subscription) bool {
	return false
}

// ValidateCredentials checks username length (3-100) and password minimum length (6)
func (d *Domain) ValidateCredentials(username, password string) error {
	if len(username) < 3 || len(username) > 100 {
		return errors.New("username must be between 3 and 100 characters")
	}
	if len(password) < 6 {
		return errors.New("password must be at least 6 characters")
	}
	return nil
}

// GeneratePassword generates a random alphanumeric password of given length using crypto/rand
func (d *Domain) GeneratePassword(length int) (string, error) {
	if length <= 0 {
		return "", errors.New("password length must be positive")
	}
	result := make([]byte, length)
	charLen := big.NewInt(int64(len(alphanumeric)))
	for i := range result {
		n, err := rand.Int(rand.Reader, charLen)
		if err != nil {
			return "", err
		}
		result[i] = alphanumeric[n.Int64()]
	}
	return string(result), nil
}
