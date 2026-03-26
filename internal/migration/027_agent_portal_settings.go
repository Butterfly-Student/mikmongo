package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(up027, down027)
}

func up027(ctx context.Context, tx *sql.Tx) error {
	// Add "review" status to agent_invoices CHECK constraint
	if _, err := tx.ExecContext(ctx, `ALTER TABLE agent_invoices DROP CONSTRAINT IF EXISTS agent_invoices_status_check`); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `ALTER TABLE agent_invoices ADD CONSTRAINT agent_invoices_status_check CHECK (status IN ('draft', 'unpaid', 'review', 'paid', 'cancelled'))`); err != nil {
		return err
	}

	// Seed data moved to internal/seeder/seed_settings.go

	return nil
}

func down027(ctx context.Context, tx *sql.Tx) error {
	// Revert CHECK constraint
	if _, err := tx.ExecContext(ctx, `ALTER TABLE agent_invoices DROP CONSTRAINT IF EXISTS agent_invoices_status_check`); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `ALTER TABLE agent_invoices ADD CONSTRAINT agent_invoices_status_check CHECK (status IN ('draft', 'unpaid', 'paid', 'cancelled'))`); err != nil {
		return err
	}

	return nil
}
