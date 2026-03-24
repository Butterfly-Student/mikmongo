package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(up028, down028)
}

func up028(ctx context.Context, tx *sql.Tx) error {
	// Create petty_cash_funds first (referenced by cash_entries FK)
	if _, err := tx.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS petty_cash_funds (
			id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			fund_name       VARCHAR(100) NOT NULL,
			initial_balance DECIMAL(15,2) NOT NULL DEFAULT 0,
			current_balance DECIMAL(15,2) NOT NULL DEFAULT 0,
			custodian_id    UUID NOT NULL REFERENCES users(id),
			status          VARCHAR(20) NOT NULL DEFAULT 'active'
			                CHECK (status IN ('active', 'closed')),
			created_at      TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at      TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
			deleted_at      TIMESTAMPTZ
		)
	`); err != nil {
		return err
	}

	// Create cash_entries
	if _, err := tx.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS cash_entries (
			id                 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			entry_number       VARCHAR(50) UNIQUE NOT NULL,
			type               VARCHAR(10) NOT NULL CHECK (type IN ('income', 'expense')),
			source             VARCHAR(20) NOT NULL CHECK (source IN (
				'invoice','agent_invoice','installation','penalty','other',
				'operational','upstream','purchase','salary'
			)),
			amount             DECIMAL(15,2) NOT NULL CHECK (amount > 0),
			description        TEXT NOT NULL,
			reference_type     VARCHAR(20),
			reference_id       UUID,
			payment_method     VARCHAR(20) NOT NULL DEFAULT 'cash'
			                   CHECK (payment_method IN ('cash','bank_transfer','e-wallet','qris','gateway')),
			bank_name          VARCHAR(100),
			account_number     VARCHAR(50),
			petty_cash_fund_id UUID REFERENCES petty_cash_funds(id),
			entry_date         TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
			created_by         UUID NOT NULL REFERENCES users(id),
			approved_by        UUID REFERENCES users(id),
			approved_at        TIMESTAMPTZ,
			status             VARCHAR(20) NOT NULL DEFAULT 'pending'
			                   CHECK (status IN ('pending', 'approved', 'rejected')),
			notes              TEXT,
			receipt_image      TEXT,
			created_at         TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at         TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
			deleted_at         TIMESTAMPTZ
		)
	`); err != nil {
		return err
	}

	// Indexes
	for _, ddl := range []string{
		`CREATE INDEX IF NOT EXISTS idx_cash_entries_type ON cash_entries(type) WHERE deleted_at IS NULL`,
		`CREATE INDEX IF NOT EXISTS idx_cash_entries_source ON cash_entries(source) WHERE deleted_at IS NULL`,
		`CREATE INDEX IF NOT EXISTS idx_cash_entries_status ON cash_entries(status) WHERE deleted_at IS NULL`,
		`CREATE INDEX IF NOT EXISTS idx_cash_entries_entry_date ON cash_entries(entry_date)`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_cash_entries_reference ON cash_entries(reference_type, reference_id) WHERE reference_type IS NOT NULL AND deleted_at IS NULL`,
	} {
		if _, err := tx.ExecContext(ctx, ddl); err != nil {
			return err
		}
	}

	// Seed sequence counter for KAS000001
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO sequence_counters (name, prefix, padding, last_number)
		VALUES ('cash_entry_number', 'KAS', 6, 0)
		ON CONFLICT (name) DO NOTHING
	`); err != nil {
		return err
	}

	return nil
}

func down028(ctx context.Context, tx *sql.Tx) error {
	if _, err := tx.ExecContext(ctx, `DROP TABLE IF EXISTS cash_entries`); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `DROP TABLE IF EXISTS petty_cash_funds`); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM sequence_counters WHERE name = 'cash_entry_number'`); err != nil {
		return err
	}
	return nil
}
