package scheduler

import (
	"context"
	"log"
	"time"

	"github.com/robfig/cron/v3"
	"mikmongo/internal/service"
)

// ReminderScheduler sends payment reminders for upcoming due invoices
type ReminderScheduler struct {
	cron       *cron.Cron
	billingSvc *service.BillingService
}

// NewReminderScheduler creates a new reminder scheduler
func NewReminderScheduler(c *cron.Cron, billingSvc *service.BillingService) *ReminderScheduler {
	return &ReminderScheduler{
		cron:       c,
		billingSvc: billingSvc,
	}
}

// Start schedules the reminder job
func (s *ReminderScheduler) Start() {
	// Run daily at 08:00
	s.cron.AddFunc("0 8 * * *", func() {
		log.Println("Running payment reminder scheduler...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		if err := s.billingSvc.CheckAndSendReminders(ctx); err != nil {
			log.Printf("Reminder scheduler error: %v", err)
		}
	})
	log.Println("Reminder scheduler started")
}
