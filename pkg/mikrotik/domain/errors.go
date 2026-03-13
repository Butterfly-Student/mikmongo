package domain

import "errors"

// Common Mikrotik errors
var (
	ErrConnectionFailed = errors.New("failed to connect to Mikrotik")
	ErrAuthentication   = errors.New("authentication failed")
	ErrCommandFailed    = errors.New("command execution failed")
	ErrDeviceNotFound   = errors.New("device not found")
	ErrInvalidParameter = errors.New("invalid parameter")
	ErrTimeout          = errors.New("operation timed out")
)

// IsConnectionError checks if error is connection-related
func IsConnectionError(err error) bool {
	if err == nil {
		return false
	}
	return errors.Is(err, ErrConnectionFailed) || errors.Is(err, ErrTimeout)
}
