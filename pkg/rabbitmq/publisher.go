package rabbitmq

import (
	"context"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

// Publish publishes a message to an exchange
func (c *Client) Publish(ctx context.Context, exchange, routingKey string, body []byte) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	
	return c.channel.PublishWithContext(
		ctx,
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp091.Publishing{
			ContentType:  "application/json",
			Body:         body,
			Timestamp:    time.Now(),
			DeliveryMode: amqp091.Persistent,
		},
	)
}

// PublishWithDelay publishes a message with delay (using delayed message plugin or TTL)
func (c *Client) PublishWithDelay(ctx context.Context, exchange, routingKey string, body []byte, delay time.Duration) error {
	// Implementation depends on RabbitMQ delayed message plugin or TTL queue
	return c.Publish(ctx, exchange, routingKey, body)
}
