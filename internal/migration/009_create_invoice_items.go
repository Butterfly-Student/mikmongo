package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(up009, down009)
}

// up009 creates the invoice_items table for invoice line items.
func up009(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec(`
	-- ============================================
	-- INVOICE ITEMS TABLE
	-- ============================================
	CREATE TABLE IF NOT EXISTS invoice_items (
		id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		invoice_id UUID NOT NULL REFERENCES invoices(id) ON DELETE CASCADE,

		item_type   VARCHAR(20)
		            CHECK (item_type IN ('subscription', 'installation', 'equipment', 'other')),
		description VARCHAR(255) NOT NULL,
		profile_id  UUID REFERENCES bandwidth_profiles(id) ON DELETE SET NULL,

		quantity   INTEGER      NOT NULL DEFAULT 1,
		unit_price DECIMAL(12,2) NOT NULL,
		subtotal   DECIMAL(12,2) NOT NULL,
		tax_rate   DECIMAL(5,4)  DEFAULT 0,
		tax_amount DECIMAL(12,2) DEFAULT 0,
		total      DECIMAL(12,2) NOT NULL,

		-- Proration
		is_prorated         BOOLEAN DEFAULT false,
		proration_days      INTEGER,
		proration_percentage DECIMAL(5,2),

		sort_order INTEGER DEFAULT 0,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
	);

	CREATE INDEX IF NOT EXISTS idx_invoice_items_invoice ON invoice_items(invoice_id);
	CREATE INDEX IF NOT EXISTS idx_invoice_items_profile ON invoice_items(profile_id) WHERE profile_id IS NOT NULL;

	COMMENT ON TABLE invoice_items IS 'Detail item dalam invoice (mendukung proration)';
	`)
	return err
}

func down009(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE IF EXISTS invoice_items CASCADE;`)
	return err
}
