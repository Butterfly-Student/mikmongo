package rabbitmq

import (
	"context"
	"log"
)

// Handler is the message handler function type
type Handler func(ctx context.Context, body []byte) error

// Subscribe subscribes to a queue
func (c *Client) Subscribe(ctx context.Context, queueName string, handler Handler) error {
	msgs, err := c.channel.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return err
	}
	
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-msgs:
				if !ok {
					return
				}
				
				if err := handler(ctx, msg.Body); err != nil {
					log.Printf("Message handler error: %v", err)
					msg.Nack(false, true) // requeue
				} else {
					msg.Ack(false)
				}
			}
		}
	}()
	
	return nil
}

// SubscribeWithQOS subscribes with quality of service settings
func (c *Client) SubscribeWithQOS(ctx context.Context, queueName string, prefetchCount int, handler Handler) error {
	err := c.channel.Qos(prefetchCount, 0, false)
	if err != nil {
		return err
	}
	return c.Subscribe(ctx, queueName, handler)
}
