package repository

import (
	"context"

	"github.com/google/uuid"
	"mikmongo/internal/model"
)

// PaymentAllocationRepository defines the interface for payment allocation data access
type PaymentAllocationRepository interface {
	Create(ctx context.Context, allocation *model.PaymentAllocation) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.PaymentAllocation, error)
	Update(ctx context.Context, allocation *model.PaymentAllocation) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]model.PaymentAllocation, error)
	Count(ctx context.Context) (int64, error)
	ListByPaymentID(ctx context.Context, paymentID uuid.UUID) ([]model.PaymentAllocation, error)
	ListByInvoiceID(ctx context.Context, invoiceID uuid.UUID) ([]model.PaymentAllocation, error)
}
