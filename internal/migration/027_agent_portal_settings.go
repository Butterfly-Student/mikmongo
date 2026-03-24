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

	// Seed default billing settings for agents
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO system_settings (group_name, key_name, value, value_type, label, description) VALUES
		  ('billing', 'agent_default_billing_cycle',       'monthly', 'string',  'Default Siklus Tagihan Agen', 'Siklus tagihan default untuk agent baru: monthly atau weekly'),
		  ('billing', 'agent_default_billing_day_monthly', '1',       'integer', 'Default Hari Tagihan (Monthly)', 'Tanggal invoice diproses untuk agent monthly (1-28)'),
		  ('billing', 'agent_default_billing_day_weekly',  '1',       'integer', 'Default Hari Tagihan (Weekly)', 'Hari dalam minggu untuk agent weekly (1=Senin, 7=Minggu)')
		ON CONFLICT (group_name, key_name) DO NOTHING
	`); err != nil {
		return err
	}

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

	// Remove seeded settings
	if _, err := tx.ExecContext(ctx, `DELETE FROM system_settings WHERE group_name = 'billing' AND key_name IN ('agent_default_billing_cycle', 'agent_default_billing_day_monthly', 'agent_default_billing_day_weekly')`); err != nil {
		return err
	}

	return nil
}
