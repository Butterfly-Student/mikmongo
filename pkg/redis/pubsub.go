package redis

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
)

// Subscription wraps a Redis Pub/Sub subscription.
type Subscription struct {
	pubsub *redis.PubSub
}

// Channel returns the channel that receives published messages.
func (s *Subscription) Channel() <-chan *redis.Message {
	return s.pubsub.Channel()
}

// Close unsubscribes and releases the underlying connection.
func (s *Subscription) Close() error {
	return s.pubsub.Close()
}

// Publish sends a message to a Redis Pub/Sub channel.
// If message is not a string or []byte, it is JSON-encoded before publishing.
func (c *Client) Publish(ctx context.Context, channel string, message any) error {
	var payload any
	switch v := message.(type) {
	case string:
		payload = v
	case []byte:
		payload = v
	default:
		b, err := json.Marshal(v)
		if err != nil {
			return err
		}
		payload = b
	}
	return c.client.Publish(ctx, channel, payload).Err()
}

// Subscribe subscribes to one or more Redis Pub/Sub channels.
func (c *Client) Subscribe(ctx context.Context, channels ...string) *Subscription {
	return &Subscription{
		pubsub: c.client.Subscribe(ctx, channels...),
	}
}
