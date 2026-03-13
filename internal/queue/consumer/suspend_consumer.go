package consumer

import (
	"context"
	"encoding/json"
	"log"

	"github.com/google/uuid"
	"mikmongo/pkg/rabbitmq"
)

// SuspendHandler defines the callback for suspend events
type SuspendHandler func(ctx context.Context, customerID uuid.UUID) error

// SuspendConsumer handles customer suspend/isolate messages
type SuspendConsumer struct {
	client  *rabbitmq.Client
	handler SuspendHandler
}

// NewSuspendConsumer creates a new suspend consumer
func NewSuspendConsumer(client *rabbitmq.Client) *SuspendConsumer {
	return &SuspendConsumer{client: client}
}

// SetHandler sets the suspend event handler
func (c *SuspendConsumer) SetHandler(h SuspendHandler) {
	c.handler = h
}

// Start starts consuming suspend messages
func (c *SuspendConsumer) Start(ctx context.Context) error {
	if c.client == nil {
		return nil
	}
	return c.client.Subscribe(ctx, "suspend.queue", c.handleMessage)
}

type suspendEvent struct {
	CustomerID string `json:"customer_id"`
	RouterID   string `json:"router_id"`
	Reason     string `json:"reason"`
}

func (c *SuspendConsumer) handleMessage(ctx context.Context, body []byte) error {
	var event suspendEvent
	if err := json.Unmarshal(body, &event); err != nil {
		log.Printf("Failed to unmarshal suspend event: %v", err)
		return err
	}

	if c.handler == nil {
		log.Printf("Suspend consumer: no handler set, skipping")
		return nil
	}

	customerID, err := uuid.Parse(event.CustomerID)
	if err != nil {
		return nil
	}

	if err := c.handler(ctx, customerID); err != nil {
		log.Printf("Failed to isolate customer %s: %v", event.CustomerID, err)
		return err
	}
	return nil
}
