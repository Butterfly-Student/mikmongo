package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(up010, down010)
}

// up010 creates the payments table for customer payments.
func up010(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec(`
	-- ============================================
	-- PAYMENTS TABLE
	-- ============================================
	CREATE TABLE IF NOT EXISTS payments (
		id             UUID      PRIMARY KEY DEFAULT gen_random_uuid(),
		payment_number VARCHAR(50) UNIQUE NOT NULL,

		customer_id UUID NOT NULL REFERENCES customers(id) ON DELETE RESTRICT,

		-- Nominal
		amount           DECIMAL(12,2) NOT NULL,
		allocated_amount DECIMAL(12,2) DEFAULT 0,
		remaining_amount DECIMAL(12,2) GENERATED ALWAYS AS (amount - allocated_amount) STORED,

		-- Metode Pembayaran
		payment_method VARCHAR(20) NOT NULL
		               CHECK (payment_method IN ('cash', 'bank_transfer', 'e-wallet', 'credit_card', 'debit_card', 'check', 'qris', 'gateway')),
		payment_date   TIMESTAMPTZ  NOT NULL,

		-- Detail Transfer Bank
		bank_name           VARCHAR(100),
		bank_account_number VARCHAR(50),
		bank_account_name   VARCHAR(100),
		transaction_reference VARCHAR(100),

		-- Detail E-Wallet
		ewallet_provider VARCHAR(50),
		ewallet_number   VARCHAR(50),

		-- Detail Payment Gateway
		gateway_name     VARCHAR(50),
		gateway_trx_id   VARCHAR(150),
		gateway_response JSONB,

		-- Backward-compat Xendit fields
		xendit_invoice_id      VARCHAR(100),
		xendit_external_id     VARCHAR(100),
		xendit_payment_channel VARCHAR(50),

		-- Bukti Bayar
		proof_image    TEXT,
		receipt_number VARCHAR(50),

		-- Status
		status VARCHAR(20) NOT NULL
		       CHECK (status IN ('pending', 'confirmed', 'rejected', 'refunded'))
		       DEFAULT 'pending',

		-- Verifikasi
		processed_by UUID REFERENCES users(id) ON DELETE SET NULL,
		processed_at TIMESTAMPTZ,
		rejection_reason TEXT,

		-- Refund
		refund_amount DECIMAL(12,2) DEFAULT 0,
		refund_date   TIMESTAMPTZ,
		refund_reason TEXT,
		refunded_by   UUID REFERENCES users(id) ON DELETE SET NULL,

		notes      TEXT,
		created_by UUID REFERENCES users(id) ON DELETE SET NULL,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
		deleted_at TIMESTAMPTZ
	);

	CREATE INDEX IF NOT EXISTS idx_payments_customer       ON payments(customer_id)    WHERE deleted_at IS NULL;
	CREATE INDEX IF NOT EXISTS idx_payments_status         ON payments(status)         WHERE deleted_at IS NULL;
	CREATE INDEX IF NOT EXISTS idx_payments_date           ON payments(payment_date)   WHERE deleted_at IS NULL;
	CREATE INDEX IF NOT EXISTS idx_payments_method         ON payments(payment_method) WHERE deleted_at IS NULL;
	CREATE INDEX IF NOT EXISTS idx_payments_number         ON payments(payment_number);
	CREATE INDEX IF NOT EXISTS idx_payments_gateway        ON payments(gateway_name, gateway_trx_id) WHERE gateway_trx_id IS NOT NULL;
	CREATE INDEX IF NOT EXISTS idx_payments_xendit_invoice ON payments(xendit_invoice_id) WHERE xendit_invoice_id IS NOT NULL;
	CREATE INDEX IF NOT EXISTS idx_payments_deleted        ON payments(deleted_at);

	COMMENT ON TABLE  payments                IS 'Pembayaran dari pelanggan';
	COMMENT ON COLUMN payments.remaining_amount IS 'Calculated: amount - allocated_amount (untuk advance payment)';
	COMMENT ON COLUMN payments.gateway_name    IS 'Nama gateway: xendit, midtrans, dll';
	COMMENT ON COLUMN payments.gateway_response IS 'Raw response JSON dari payment gateway';
	`)
	return err
}

func down010(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE IF EXISTS payments CASCADE;`)
	return err
}
