package scheduler

import (
	"context"
	"log"
	"time"

	"github.com/robfig/cron/v3"
	"mikmongo/internal/service"
)

// SyncScheduler handles MikroTik sync cron jobs
type SyncScheduler struct {
	cron            *cron.Cron
	routerService   *service.RouterService
	subscriptionSvc *service.SubscriptionService
}

// NewSyncScheduler creates a new sync scheduler
func NewSyncScheduler(c *cron.Cron, routerService *service.RouterService, subscriptionSvc *service.SubscriptionService) *SyncScheduler {
	return &SyncScheduler{
		cron:            c,
		routerService:   routerService,
		subscriptionSvc: subscriptionSvc,
	}
}

// Start schedules sync jobs
func (s *SyncScheduler) Start() {
	// Run every 5 minutes
	s.cron.AddFunc("*/5 * * * *", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Minute)
		defer cancel()

		// Sync router health
		if err := s.routerService.SyncAllDevices(ctx); err != nil {
			log.Printf("Router sync error: %v", err)
		}

		// Note: ProcessPendingSync removed - all operations now sync automatically
	})
	log.Println("Sync scheduler started")
}
