package postgres

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"mikmongo/internal/model"
	"mikmongo/internal/repository"
)

// paymentRepository implements repository.PaymentRepository
type paymentRepository struct {
	db *gorm.DB
}

// NewPaymentRepository creates a new payment repository
func NewPaymentRepository(db *gorm.DB) repository.PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) Create(ctx context.Context, payment *model.Payment) error {
	return r.db.WithContext(ctx).Create(payment).Error
}

func (r *paymentRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Payment, error) {
	var payment model.Payment
	err := r.db.WithContext(ctx).Preload("Customer").First(&payment, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *paymentRepository) GetByInvoiceID(ctx context.Context, invoiceID uuid.UUID) ([]model.Payment, error) {
	var payments []model.Payment
	err := r.db.WithContext(ctx).
		Joins("JOIN payment_allocations ON payment_allocations.payment_id = payments.id").
		Where("payment_allocations.invoice_id = ?", invoiceID).
		Find(&payments).Error
	return payments, err
}

func (r *paymentRepository) GetByTransactionID(ctx context.Context, transactionID string) (*model.Payment, error) {
	var payment model.Payment
	err := r.db.WithContext(ctx).First(&payment, "transaction_reference = ?", transactionID).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *paymentRepository) Update(ctx context.Context, payment *model.Payment) error {
	return r.db.WithContext(ctx).Save(payment).Error
}

func (r *paymentRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	return r.db.WithContext(ctx).Model(&model.Payment{}).Where("id = ?", id).Update("status", status).Error
}

func (r *paymentRepository) GetByCustomerID(ctx context.Context, customerID uuid.UUID) ([]model.Payment, error) {
	var payments []model.Payment
	err := r.db.WithContext(ctx).
		Where("customer_id = ? AND deleted_at IS NULL", customerID).
		Order("payment_date DESC").
		Find(&payments).Error
	return payments, err
}

func (r *paymentRepository) List(ctx context.Context, limit, offset int) ([]model.Payment, error) {
	var payments []model.Payment
	err := r.db.WithContext(ctx).Preload("Customer").Limit(limit).Offset(offset).Find(&payments).Error
	return payments, err
}
