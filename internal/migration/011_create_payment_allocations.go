package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(up011, down011)
}

// up011 creates the payment_allocations table for linking payments to invoices.
func up011(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec(`
	-- ============================================
	-- PAYMENT ALLOCATIONS TABLE
	-- ============================================
	CREATE TABLE IF NOT EXISTS payment_allocations (
		id               UUID      PRIMARY KEY DEFAULT gen_random_uuid(),
		payment_id       UUID      NOT NULL REFERENCES payments(id) ON DELETE CASCADE,
		invoice_id       UUID      NOT NULL REFERENCES invoices(id) ON DELETE CASCADE,
		allocated_amount DECIMAL(12,2) NOT NULL,
		created_at       TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
		UNIQUE(payment_id, invoice_id)
	);

	CREATE INDEX IF NOT EXISTS idx_payment_allocations_payment ON payment_allocations(payment_id);
	CREATE INDEX IF NOT EXISTS idx_payment_allocations_invoice ON payment_allocations(invoice_id);

	COMMENT ON TABLE payment_allocations IS 'Alokasi 1 pembayaran ke multiple invoice (partial/advance payment)';
	`)
	return err
}

func down011(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE IF EXISTS payment_allocations CASCADE;`)
	return err
}
