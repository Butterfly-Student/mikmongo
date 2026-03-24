package scheduler

import (
	"github.com/robfig/cron/v3"
	"mikmongo/internal/queue"
	"mikmongo/internal/service"
)

// Registry holds all scheduler instances
type Registry struct {
	cron              *cron.Cron
	billingScheduler  *BillingScheduler
	isolateScheduler  *IsolateScheduler
	reminderScheduler *ReminderScheduler
	syncScheduler     *SyncScheduler
}

// NewRegistry creates a new scheduler registry
func NewRegistry(
	services *service.Registry,
	q *queue.Registry,
) *Registry {
	c := cron.New()
	return &Registry{
		cron:              c,
		billingScheduler:  NewBillingScheduler(c, services.Billing, services.AgentInvoice, q.BillingProducer),
		isolateScheduler:  NewIsolateScheduler(c, services.Billing),
		reminderScheduler: NewReminderScheduler(c, services.Billing),
		syncScheduler:     NewSyncScheduler(c, services.Router, services.Subscription),
	}
}

// Start starts all schedulers
func (r *Registry) Start() {
	r.billingScheduler.Start()
	r.isolateScheduler.Start()
	r.reminderScheduler.Start()
	r.syncScheduler.Start()
	r.cron.Start()
}

// Stop stops all schedulers
func (r *Registry) Stop() {
	r.cron.Stop()
}
