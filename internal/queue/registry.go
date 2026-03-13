// Package queue contains RabbitMQ producers and consumers
package queue

import (
	"mikmongo/internal/queue/consumer"
	"mikmongo/internal/queue/producer"
	"mikmongo/pkg/rabbitmq"
)

// Registry holds queue producers and consumers
type Registry struct {
	Client               *rabbitmq.Client
	BillingProducer      *producer.BillingProducer
	SuspendProducer      *producer.SuspendProducer
	NotificationProducer *producer.NotificationProducer
	BillingConsumer      *consumer.BillingConsumer
	SuspendConsumer      *consumer.SuspendConsumer
	NotificationConsumer *consumer.NotificationConsumer
}

// NewRegistry creates a new queue registry
func NewRegistry(client *rabbitmq.Client) *Registry {
	return &Registry{
		Client:               client,
		BillingProducer:      producer.NewBillingProducer(client),
		SuspendProducer:      producer.NewSuspendProducer(client),
		NotificationProducer: producer.NewNotificationProducer(client),
		BillingConsumer:      consumer.NewBillingConsumer(client),
		SuspendConsumer:      consumer.NewSuspendConsumer(client),
		NotificationConsumer: consumer.NewNotificationConsumer(client),
	}
}

// Setup initializes exchanges and queues
func (r *Registry) Setup() error {
	return nil
}
