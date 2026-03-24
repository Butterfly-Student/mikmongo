package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(up024, down024)
}

func up024(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
	CREATE TABLE IF NOT EXISTS hotspot_sales (
		id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		router_id      UUID NOT NULL REFERENCES mikrotik_routers(id),
		username       VARCHAR(100) NOT NULL,
		profile        VARCHAR(100) NOT NULL,
		price          DECIMAL(15,2) NOT NULL DEFAULT 0,
		selling_price  DECIMAL(15,2) NOT NULL DEFAULT 0,
		prefix         VARCHAR(20),
		batch_code     VARCHAR(10),
		sales_agent_id UUID REFERENCES sales_agents(id) ON DELETE SET NULL,
		created_at     TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_hotspot_sales_router_id      ON hotspot_sales(router_id);
	CREATE INDEX IF NOT EXISTS idx_hotspot_sales_sales_agent_id ON hotspot_sales(sales_agent_id);
	CREATE INDEX IF NOT EXISTS idx_hotspot_sales_profile        ON hotspot_sales(profile);
	CREATE INDEX IF NOT EXISTS idx_hotspot_sales_batch_code     ON hotspot_sales(batch_code);
	CREATE INDEX IF NOT EXISTS idx_hotspot_sales_created_at     ON hotspot_sales(created_at);

	COMMENT ON TABLE  hotspot_sales IS 'Catatan penjualan voucher hotspot';
	COMMENT ON COLUMN hotspot_sales.price         IS 'Harga modal/base price';
	COMMENT ON COLUMN hotspot_sales.selling_price IS 'Harga jual ke pelanggan';
	COMMENT ON COLUMN hotspot_sales.batch_code    IS 'Kode batch dari VoucherBatch.Code';
	`)
	return err
}

func down024(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `DROP TABLE IF EXISTS hotspot_sales CASCADE;`)
	return err
}
