package repository

import (
	"context"

	"github.com/google/uuid"
	"mikmongo/internal/model"
)

// PaymentRepository defines the interface for payment data access
type PaymentRepository interface {
	Create(ctx context.Context, payment *model.Payment) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Payment, error)
	GetByInvoiceID(ctx context.Context, invoiceID uuid.UUID) ([]model.Payment, error)
	GetByTransactionID(ctx context.Context, transactionID string) (*model.Payment, error)
	GetByCustomerID(ctx context.Context, customerID uuid.UUID) ([]model.Payment, error)
	Update(ctx context.Context, payment *model.Payment) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
	List(ctx context.Context, limit, offset int) ([]model.Payment, error)
}
