package repository

import (
	"context"

	"github.com/google/uuid"
	"mikmongo/internal/model"
)

// InvoiceItemRepository defines the interface for invoice item data access
type InvoiceItemRepository interface {
	Create(ctx context.Context, item *model.InvoiceItem) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.InvoiceItem, error)
	Update(ctx context.Context, item *model.InvoiceItem) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]model.InvoiceItem, error)
	Count(ctx context.Context) (int64, error)
	ListByInvoiceID(ctx context.Context, invoiceID uuid.UUID) ([]model.InvoiceItem, error)
	DeleteByInvoiceID(ctx context.Context, invoiceID uuid.UUID) error
}
