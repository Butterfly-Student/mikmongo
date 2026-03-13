package postgres

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"mikmongo/internal/model"
	"mikmongo/internal/repository"
)

type paymentAllocationRepository struct {
	db *gorm.DB
}

func NewPaymentAllocationRepository(db *gorm.DB) repository.PaymentAllocationRepository {
	return &paymentAllocationRepository{db: db}
}

func (r *paymentAllocationRepository) Create(ctx context.Context, allocation *model.PaymentAllocation) error {
	return r.db.WithContext(ctx).Create(allocation).Error
}

func (r *paymentAllocationRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.PaymentAllocation, error) {
	var allocation model.PaymentAllocation
	err := r.db.WithContext(ctx).First(&allocation, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &allocation, nil
}

func (r *paymentAllocationRepository) Update(ctx context.Context, allocation *model.PaymentAllocation) error {
	return r.db.WithContext(ctx).Save(allocation).Error
}

func (r *paymentAllocationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.PaymentAllocation{}, "id = ?", id).Error
}

func (r *paymentAllocationRepository) List(ctx context.Context, limit, offset int) ([]model.PaymentAllocation, error) {
	var allocations []model.PaymentAllocation
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&allocations).Error
	return allocations, err
}

func (r *paymentAllocationRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.PaymentAllocation{}).Count(&count).Error
	return count, err
}

func (r *paymentAllocationRepository) ListByPaymentID(ctx context.Context, paymentID uuid.UUID) ([]model.PaymentAllocation, error) {
	var allocations []model.PaymentAllocation
	err := r.db.WithContext(ctx).Where("payment_id = ?", paymentID).Find(&allocations).Error
	return allocations, err
}

func (r *paymentAllocationRepository) ListByInvoiceID(ctx context.Context, invoiceID uuid.UUID) ([]model.PaymentAllocation, error) {
	var allocations []model.PaymentAllocation
	err := r.db.WithContext(ctx).Where("invoice_id = ?", invoiceID).Find(&allocations).Error
	return allocations, err
}
