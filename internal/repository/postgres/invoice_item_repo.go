package postgres

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"mikmongo/internal/model"
	"mikmongo/internal/repository"
)

type invoiceItemRepository struct {
	db *gorm.DB
}

func NewInvoiceItemRepository(db *gorm.DB) repository.InvoiceItemRepository {
	return &invoiceItemRepository{db: db}
}

func (r *invoiceItemRepository) Create(ctx context.Context, item *model.InvoiceItem) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *invoiceItemRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.InvoiceItem, error) {
	var item model.InvoiceItem
	err := r.db.WithContext(ctx).First(&item, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *invoiceItemRepository) Update(ctx context.Context, item *model.InvoiceItem) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *invoiceItemRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.InvoiceItem{}, "id = ?", id).Error
}

func (r *invoiceItemRepository) List(ctx context.Context, limit, offset int) ([]model.InvoiceItem, error) {
	var items []model.InvoiceItem
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&items).Error
	return items, err
}

func (r *invoiceItemRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.InvoiceItem{}).Count(&count).Error
	return count, err
}

func (r *invoiceItemRepository) ListByInvoiceID(ctx context.Context, invoiceID uuid.UUID) ([]model.InvoiceItem, error) {
	var items []model.InvoiceItem
	err := r.db.WithContext(ctx).Where("invoice_id = ?", invoiceID).Order("sort_order").Find(&items).Error
	return items, err
}

func (r *invoiceItemRepository) DeleteByInvoiceID(ctx context.Context, invoiceID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.InvoiceItem{}, "invoice_id = ?", invoiceID).Error
}
