package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(up026, down026)
}

func up026(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
	CREATE TABLE IF NOT EXISTS agent_invoices (
		id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		agent_id        UUID NOT NULL REFERENCES sales_agents(id) ON DELETE RESTRICT,
		router_id       UUID NOT NULL REFERENCES mikrotik_routers(id),
		invoice_number  VARCHAR(20) UNIQUE NOT NULL,
		billing_cycle   VARCHAR(20) NOT NULL CHECK (billing_cycle IN ('weekly', 'monthly')),
		period_start    TIMESTAMPTZ NOT NULL,
		period_end      TIMESTAMPTZ NOT NULL,
		billing_month   INT,
		billing_week    INT,
		billing_year    INT NOT NULL,
		voucher_count   INT NOT NULL DEFAULT 0,
		subtotal        DECIMAL(15,2) NOT NULL DEFAULT 0,
		selling_total   DECIMAL(15,2) NOT NULL DEFAULT 0,
		profit          DECIMAL(15,2) GENERATED ALWAYS AS (selling_total - subtotal) STORED,
		discount_amount DECIMAL(15,2) NOT NULL DEFAULT 0,
		total_amount    DECIMAL(15,2) NOT NULL DEFAULT 0,
		paid_amount     DECIMAL(15,2) NOT NULL DEFAULT 0,
		balance         DECIMAL(15,2) GENERATED ALWAYS AS (total_amount - paid_amount) STORED,
		status          VARCHAR(20) NOT NULL DEFAULT 'unpaid'
		                CHECK (status IN ('draft', 'unpaid', 'paid', 'cancelled')),
		notes           TEXT,
		created_at      TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at      TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
		deleted_at      TIMESTAMPTZ,

		UNIQUE (agent_id, period_start, billing_cycle)
	);

	CREATE INDEX IF NOT EXISTS idx_agent_invoices_agent_id      ON agent_invoices(agent_id);
	CREATE INDEX IF NOT EXISTS idx_agent_invoices_status        ON agent_invoices(status) WHERE deleted_at IS NULL;
	CREATE INDEX IF NOT EXISTS idx_agent_invoices_billing_year  ON agent_invoices(billing_year, billing_week, billing_month);
	CREATE INDEX IF NOT EXISTS idx_agent_invoices_period_end    ON agent_invoices(period_end) WHERE status = 'unpaid' AND deleted_at IS NULL;
	CREATE INDEX IF NOT EXISTS idx_agent_invoices_deleted       ON agent_invoices(deleted_at);

	COMMENT ON TABLE  agent_invoices IS 'Invoice periodik penjualan voucher hotspot per agen';
	COMMENT ON COLUMN agent_invoices.subtotal      IS 'Total harga modal (SUM price dari hotspot_sales)';
	COMMENT ON COLUMN agent_invoices.selling_total IS 'Total harga jual (SUM selling_price dari hotspot_sales)';
	COMMENT ON COLUMN agent_invoices.profit        IS 'Keuntungan agen: selling_total - subtotal (generated)';
	COMMENT ON COLUMN agent_invoices.total_amount  IS 'Jumlah tagihan: selling_total - discount';
	COMMENT ON COLUMN agent_invoices.balance       IS 'Sisa tagihan: total_amount - paid_amount (generated)';
	COMMENT ON COLUMN agent_invoices.billing_month IS 'Bulan tagihan (1-12), NULL untuk weekly';
	COMMENT ON COLUMN agent_invoices.billing_week  IS 'Nomor minggu ISO (1-53), NULL untuk monthly';
	`)
	return err
}

func down026(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `DROP TABLE IF EXISTS agent_invoices CASCADE;`)
	return err
}
