package consumer

import (
	"context"
	"encoding/json"
	"log"

	"mikmongo/pkg/rabbitmq"
)

// NotificationHandler defines callbacks for notification events
type NotificationHandler struct {
	SendWhatsApp func(ctx context.Context, phone, message string) error
	SendEmail    func(ctx context.Context, to, subject, body string) error
}

// NotificationConsumer handles notification messages
type NotificationConsumer struct {
	client  *rabbitmq.Client
	handler *NotificationHandler
}

// NewNotificationConsumer creates a new notification consumer
func NewNotificationConsumer(client *rabbitmq.Client) *NotificationConsumer {
	return &NotificationConsumer{client: client}
}

// SetHandler sets the notification handler
func (c *NotificationConsumer) SetHandler(h *NotificationHandler) {
	c.handler = h
}

// Start starts consuming notification messages
func (c *NotificationConsumer) Start(ctx context.Context) error {
	if c.client == nil {
		return nil
	}
	return c.client.Subscribe(ctx, "notification.queue", c.handleMessage)
}

type notificationEvent struct {
	Type    string `json:"type"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func (c *NotificationConsumer) handleMessage(ctx context.Context, body []byte) error {
	var event notificationEvent
	if err := json.Unmarshal(body, &event); err != nil {
		log.Printf("Failed to unmarshal notification event: %v", err)
		return err
	}

	if c.handler == nil {
		log.Printf("Notification consumer: no handler set, skipping")
		return nil
	}

	switch event.Type {
	case "whatsapp", "wa":
		if c.handler.SendWhatsApp != nil {
			return c.handler.SendWhatsApp(ctx, event.To, event.Body)
		}
	case "email":
		if c.handler.SendEmail != nil {
			return c.handler.SendEmail(ctx, event.To, event.Subject, event.Body)
		}
	default:
		log.Printf("Unknown notification type: %s", event.Type)
	}
	return nil
}
