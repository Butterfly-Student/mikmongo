package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"mikmongo/internal/model"
)

// AgentInvoiceFilter holds optional filters for listing agent invoices.
type AgentInvoiceFilter struct {
	AgentID      *uuid.UUID
	RouterID     *uuid.UUID
	Status       string
	BillingCycle string
	BillingYear  *int
	BillingMonth *int
	BillingWeek  *int
}

// AgentInvoiceRepository defines persistence operations for agent invoices.
type AgentInvoiceRepository interface {
	Create(ctx context.Context, inv *model.AgentInvoice) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.AgentInvoice, error)

	// GetByAgentAndPeriod is used for idempotency — returns an existing invoice
	// for the given agent, period start, and billing cycle if one already exists.
	GetByAgentAndPeriod(ctx context.Context, agentID uuid.UUID, periodStart time.Time, cycle string) (*model.AgentInvoice, error)

	Update(ctx context.Context, inv *model.AgentInvoice) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status string, paidAmount float64) error
	UpdateStatusAndNotes(ctx context.Context, id uuid.UUID, status string, paidAmount float64, notes string) error
	Delete(ctx context.Context, id uuid.UUID) error

	List(ctx context.Context, filter AgentInvoiceFilter, limit, offset int) ([]model.AgentInvoice, error)
	Count(ctx context.Context, filter AgentInvoiceFilter) (int64, error)
	ListByAgentID(ctx context.Context, agentID uuid.UUID, limit, offset int) ([]model.AgentInvoice, error)

	// GetUnpaidOverdue returns unpaid invoices whose period_end is in the past.
	GetUnpaidOverdue(ctx context.Context) ([]model.AgentInvoice, error)
}
