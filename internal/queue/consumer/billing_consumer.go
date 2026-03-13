package consumer

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
	"mikmongo/pkg/rabbitmq"
)

// BillingHandler defines the callback for billing events
type BillingHandler func(ctx context.Context, subscriptionID uuid.UUID, period time.Time) error

// BillingConsumer handles billing-related messages
type BillingConsumer struct {
	client  *rabbitmq.Client
	handler BillingHandler
}

// NewBillingConsumer creates a new billing consumer
func NewBillingConsumer(client *rabbitmq.Client) *BillingConsumer {
	return &BillingConsumer{client: client}
}

// SetHandler sets the billing event handler
func (c *BillingConsumer) SetHandler(h BillingHandler) {
	c.handler = h
}

// Start starts consuming billing messages
func (c *BillingConsumer) Start(ctx context.Context) error {
	if c.client == nil {
		return nil
	}
	return c.client.Subscribe(ctx, "billing.queue", c.handleMessage)
}

type billingEvent struct {
	SubscriptionID string `json:"subscription_id"`
	Period         string `json:"period"`
}

func (c *BillingConsumer) handleMessage(ctx context.Context, body []byte) error {
	var event billingEvent
	if err := json.Unmarshal(body, &event); err != nil {
		log.Printf("Failed to unmarshal billing event: %v", err)
		return err
	}

	if c.handler == nil {
		log.Printf("Billing consumer: no handler set, skipping event")
		return nil
	}

	subID, err := uuid.Parse(event.SubscriptionID)
	if err != nil {
		return nil
	}

	period := time.Now()
	if event.Period != "" {
		if t, err := time.Parse("2006-01", event.Period); err == nil {
			period = t
		}
	}

	if err := c.handler(ctx, subID, period); err != nil {
		log.Printf("Failed to process billing event for sub %s: %v", event.SubscriptionID, err)
		return err
	}
	return nil
}
