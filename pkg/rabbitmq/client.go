// Package rabbitmq provides RabbitMQ client wrapper
package rabbitmq

import (
	"context"
	"fmt"

	"github.com/rabbitmq/amqp091-go"
)

// Client wraps amqp.Connection
type Client struct {
	conn    *amqp091.Connection
	channel *amqp091.Channel
	url     string
}

// NewClient creates a new RabbitMQ client
func NewClient(host string, port int, user, password string) (*Client, error) {
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/", user, password, host, port)
	
	conn, err := amqp091.Dial(url)
	if err != nil {
		return nil, err
	}
	
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}
	
	return &Client{
		conn:    conn,
		channel: ch,
		url:     url,
	}, nil
}

// Close closes the connection
func (c *Client) Close() error {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// Reconnect reconnects to RabbitMQ
func (c *Client) Reconnect() error {
	c.Close()
	
	conn, err := amqp091.Dial(c.url)
	if err != nil {
		return err
	}
	
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return err
	}
	
	c.conn = conn
	c.channel = ch
	return nil
}

// Channel returns the current channel
func (c *Client) Channel() *amqp091.Channel {
	return c.channel
}

// HealthCheck performs health check
func (c *Client) HealthCheck(ctx context.Context) error {
	_, err := c.channel.QueueDeclare(
		"health_check",
		false, false, true, false, nil,
	)
	return err
}
