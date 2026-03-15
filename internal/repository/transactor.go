package repository

import "context"

// Transactor runs a function atomically within a database transaction.
// txPayment, txInvoice, and txAlloc are transaction-scoped repo instances.
type Transactor interface {
	RunInTx(ctx context.Context, fn func(
		txPayment PaymentRepository,
		txInvoice InvoiceRepository,
		txAlloc   PaymentAllocationRepository,
	) error) error
}
