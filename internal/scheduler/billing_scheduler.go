package scheduler

import (
	"context"
	"log"
	"time"

	"github.com/robfig/cron/v3"
	"mikmongo/internal/queue/producer"
	"mikmongo/internal/service"
)

// BillingScheduler handles billing cron jobs
type BillingScheduler struct {
	cron            *cron.Cron
	billingService  *service.BillingService
	billingProducer *producer.BillingProducer
}

// NewBillingScheduler creates a new billing scheduler
func NewBillingScheduler(
	c *cron.Cron,
	billingService *service.BillingService,
	billingProducer *producer.BillingProducer,
) *BillingScheduler {
	return &BillingScheduler{
		cron:            c,
		billingService:  billingService,
		billingProducer: billingProducer,
	}
}

// Start schedules billing jobs
func (s *BillingScheduler) Start() {
	// Run daily at 00:00 — generates invoices for subscriptions whose billing day is today
	s.cron.AddFunc("0 0 * * *", func() {
		log.Println("Running daily billing scheduler...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		if err := s.billingService.ProcessDailyBilling(ctx); err != nil {
			log.Printf("Billing scheduler error: %v", err)
		}
	})
	log.Println("Billing scheduler started")
}
