package repository

import (
	"context"

	"github.com/google/uuid"
	"mikmongo/internal/model"
)

// InvoiceRepository defines the interface for invoice data access
type InvoiceRepository interface {
	Create(ctx context.Context, invoice *model.Invoice) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Invoice, error)
	GetByCustomerID(ctx context.Context, customerID uuid.UUID) ([]model.Invoice, error)
	GetByCustomerIDForUpdate(ctx context.Context, customerID uuid.UUID) ([]model.Invoice, error)
	Update(ctx context.Context, invoice *model.Invoice) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]model.Invoice, error)
	GetOverdue(ctx context.Context) ([]model.Invoice, error)
}
