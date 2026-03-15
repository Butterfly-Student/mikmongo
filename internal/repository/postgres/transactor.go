package postgres

import (
	"context"

	"gorm.io/gorm"
	"mikmongo/internal/repository"
)

type gormTransactor struct{ db *gorm.DB }

// NewTransactor creates a Transactor backed by GORM.
// When called on an existing GORM transaction (e.g. inside a test tx), GORM
// automatically uses SAVEPOINTs so nesting works without extra code.
func NewTransactor(db *gorm.DB) repository.Transactor {
	return &gormTransactor{db: db}
}

func (t *gormTransactor) RunInTx(ctx context.Context, fn func(
	repository.PaymentRepository,
	repository.InvoiceRepository,
	repository.PaymentAllocationRepository,
) error) error {
	return t.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(
			NewPaymentRepository(tx),
			NewInvoiceRepository(tx),
			NewPaymentAllocationRepository(tx),
		)
	})
}
