package scheduler

import (
	"context"
	"log"
	"time"

	"github.com/robfig/cron/v3"
	"mikmongo/internal/service"
)

// IsolateScheduler checks and isolates subscriptions with overdue invoices
type IsolateScheduler struct {
	cron       *cron.Cron
	billingSvc *service.BillingService
}

// NewIsolateScheduler creates a new isolate scheduler
func NewIsolateScheduler(
	c *cron.Cron,
	billingSvc *service.BillingService,
) *IsolateScheduler {
	return &IsolateScheduler{
		cron:       c,
		billingSvc: billingSvc,
	}
}

// Start schedules the isolate job
func (s *IsolateScheduler) Start() {
	// Run daily at 02:00
	s.cron.AddFunc("0 2 * * *", func() {
		log.Println("Running isolate scheduler...")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()
		if err := s.billingSvc.CheckAndIsolateOverdue(ctx); err != nil {
			log.Printf("Isolate scheduler error: %v", err)
		}
	})
	log.Println("Isolate scheduler started")
}
