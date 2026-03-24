package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"mikmongo/internal/model"
	"mikmongo/internal/repository"
)

type agentInvoiceRepository struct {
	db *gorm.DB
}

// NewAgentInvoiceRepository creates a new agent invoice repository.
func NewAgentInvoiceRepository(db *gorm.DB) repository.AgentInvoiceRepository {
	return &agentInvoiceRepository{db: db}
}

func (r *agentInvoiceRepository) Create(ctx context.Context, inv *model.AgentInvoice) error {
	return r.db.WithContext(ctx).Create(inv).Error
}

func (r *agentInvoiceRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.AgentInvoice, error) {
	var inv model.AgentInvoice
	err := r.db.WithContext(ctx).Preload("Agent").First(&inv, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &inv, nil
}

func (r *agentInvoiceRepository) GetByAgentAndPeriod(ctx context.Context, agentID uuid.UUID, periodStart time.Time, cycle string) (*model.AgentInvoice, error) {
	var inv model.AgentInvoice
	err := r.db.WithContext(ctx).
		Where("agent_id = ? AND period_start = ? AND billing_cycle = ?", agentID.String(), periodStart, cycle).
		First(&inv).Error
	if err != nil {
		return nil, err
	}
	return &inv, nil
}

func (r *agentInvoiceRepository) Update(ctx context.Context, inv *model.AgentInvoice) error {
	return r.db.WithContext(ctx).Save(inv).Error
}

func (r *agentInvoiceRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string, paidAmount float64) error {
	return r.db.WithContext(ctx).
		Model(&model.AgentInvoice{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     status,
			"paid_amount": paidAmount,
			"updated_at": time.Now(),
		}).Error
}

func (r *agentInvoiceRepository) UpdateStatusAndNotes(ctx context.Context, id uuid.UUID, status string, paidAmount float64, notes string) error {
	return r.db.WithContext(ctx).
		Model(&model.AgentInvoice{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":      status,
			"paid_amount": paidAmount,
			"notes":       notes,
			"updated_at":  time.Now(),
		}).Error
}

func (r *agentInvoiceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.AgentInvoice{}, "id = ?", id).Error
}

func (r *agentInvoiceRepository) List(ctx context.Context, filter repository.AgentInvoiceFilter, limit, offset int) ([]model.AgentInvoice, error) {
	var invs []model.AgentInvoice
	q := r.db.WithContext(ctx).Preload("Agent")
	q = applyAgentInvoiceFilter(q, filter)
	err := q.Order("created_at DESC").Limit(limit).Offset(offset).Find(&invs).Error
	return invs, err
}

func (r *agentInvoiceRepository) Count(ctx context.Context, filter repository.AgentInvoiceFilter) (int64, error) {
	var count int64
	q := r.db.WithContext(ctx).Model(&model.AgentInvoice{})
	q = applyAgentInvoiceFilter(q, filter)
	err := q.Count(&count).Error
	return count, err
}

func (r *agentInvoiceRepository) ListByAgentID(ctx context.Context, agentID uuid.UUID, limit, offset int) ([]model.AgentInvoice, error) {
	var invs []model.AgentInvoice
	err := r.db.WithContext(ctx).
		Where("agent_id = ?", agentID.String()).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&invs).Error
	return invs, err
}

func (r *agentInvoiceRepository) GetUnpaidOverdue(ctx context.Context) ([]model.AgentInvoice, error) {
	var invs []model.AgentInvoice
	err := r.db.WithContext(ctx).
		Where("status = ? AND period_end < ?", "unpaid", time.Now()).
		Find(&invs).Error
	return invs, err
}

func applyAgentInvoiceFilter(q *gorm.DB, f repository.AgentInvoiceFilter) *gorm.DB {
	if f.AgentID != nil {
		q = q.Where("agent_id = ?", f.AgentID.String())
	}
	if f.RouterID != nil {
		q = q.Where("router_id = ?", f.RouterID.String())
	}
	if f.Status != "" {
		q = q.Where("status = ?", f.Status)
	}
	if f.BillingCycle != "" {
		q = q.Where("billing_cycle = ?", f.BillingCycle)
	}
	if f.BillingYear != nil {
		q = q.Where("billing_year = ?", *f.BillingYear)
	}
	if f.BillingMonth != nil {
		q = q.Where("billing_month = ?", *f.BillingMonth)
	}
	if f.BillingWeek != nil {
		q = q.Where("billing_week = ?", *f.BillingWeek)
	}
	return q
}
