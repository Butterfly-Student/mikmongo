package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"mikmongo/internal/model"
	"mikmongo/internal/repository"
)

// invoiceRepository implements repository.InvoiceRepository
type invoiceRepository struct {
	db *gorm.DB
}

// NewInvoiceRepository creates a new invoice repository
func NewInvoiceRepository(db *gorm.DB) repository.InvoiceRepository {
	return &invoiceRepository{db: db}
}

func (r *invoiceRepository) Create(ctx context.Context, invoice *model.Invoice) error {
	return r.db.WithContext(ctx).Create(invoice).Error
}

func (r *invoiceRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Invoice, error) {
	var invoice model.Invoice
	err := r.db.WithContext(ctx).Preload("Customer").Preload("Subscription").First(&invoice, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &invoice, nil
}

func (r *invoiceRepository) GetByCustomerID(ctx context.Context, customerID uuid.UUID) ([]model.Invoice, error) {
	var invoices []model.Invoice
	err := r.db.WithContext(ctx).Where("customer_id = ?", customerID).Find(&invoices).Error
	return invoices, err
}

func (r *invoiceRepository) GetByCustomerIDForUpdate(ctx context.Context, customerID uuid.UUID) ([]model.Invoice, error) {
	var invoices []model.Invoice
	err := r.db.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("customer_id = ?", customerID).
		Find(&invoices).Error
	return invoices, err
}

func (r *invoiceRepository) Update(ctx context.Context, invoice *model.Invoice) error {
	return r.db.WithContext(ctx).Save(invoice).Error
}

func (r *invoiceRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	return r.db.WithContext(ctx).Model(&model.Invoice{}).Where("id = ?", id).Update("status", status).Error
}

func (r *invoiceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.Invoice{}, "id = ?", id).Error
}

func (r *invoiceRepository) List(ctx context.Context, limit, offset int) ([]model.Invoice, error) {
	var invoices []model.Invoice
	err := r.db.WithContext(ctx).Preload("Customer").Preload("Subscription").Limit(limit).Offset(offset).Find(&invoices).Error
	return invoices, err
}

func (r *invoiceRepository) GetOverdue(ctx context.Context) ([]model.Invoice, error) {
	var invoices []model.Invoice
	err := r.db.WithContext(ctx).
		Where("status = ? AND due_date < ?", "unpaid", time.Now()).
		Find(&invoices).Error
	return invoices, err
}

func (r *invoiceRepository) GetBySubscriptionAndPeriod(ctx context.Context, subID uuid.UUID, month, year int) (*model.Invoice, error) {
	var inv model.Invoice
	err := r.db.WithContext(ctx).
		Where("subscription_id = ? AND billing_month = ? AND billing_year = ?",
			subID.String(), month, year).
		First(&inv).Error
	if err != nil {
		return nil, err
	}
	return &inv, nil
}
