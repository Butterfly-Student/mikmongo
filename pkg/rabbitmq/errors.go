package rabbitmq

import "errors"

// Common RabbitMQ errors
var (
	ErrConnectionClosed = errors.New("rabbitmq connection is closed")
	ErrChannelClosed    = errors.New("rabbitmq channel is closed")
	ErrPublishFailed    = errors.New("failed to publish message")
	ErrConsumeFailed    = errors.New("failed to consume messages")
	ErrQueueNotFound    = errors.New("queue not found")
	ErrExchangeNotFound = errors.New("exchange not found")
)

// IsConnectionError checks if error is connection-related
func IsConnectionError(err error) bool {
	if err == nil {
		return false
	}
	return errors.Is(err, ErrConnectionClosed) || errors.Is(err, ErrChannelClosed)
}
