package scheduler

import (
	"context"
	"log"
	"time"

	"github.com/robfig/cron/v3"
	"mikmongo/internal/queue/producer"
	"mikmongo/internal/service"
)

// SuspendScheduler handles suspension cron jobs
type SuspendScheduler struct {
	cron            *cron.Cron
	billingService  *service.BillingService
	suspendProducer *producer.SuspendProducer
}

// NewSuspendScheduler creates a new suspend scheduler
func NewSuspendScheduler(
	c *cron.Cron,
	billingService *service.BillingService,
	suspendProducer *producer.SuspendProducer,
) *SuspendScheduler {
	return &SuspendScheduler{
		cron:            c,
		billingService:  billingService,
		suspendProducer: suspendProducer,
	}
}

// Start schedules suspension check jobs
func (s *SuspendScheduler) Start() {
	// Run daily at 01:00 to mark invoices as overdue
	s.cron.AddFunc("0 1 * * *", func() {
		log.Println("Running suspend/overdue check scheduler...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		if err := s.billingService.CheckAndIsolateOverdue(ctx); err != nil {
			log.Printf("Suspend scheduler error: %v", err)
		}
	})
	log.Println("Suspend scheduler started")
}
