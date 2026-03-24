package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(up023, down023)
}

func up023(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
	CREATE TABLE IF NOT EXISTS sales_agents (
		id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		router_id      UUID NOT NULL REFERENCES mikrotik_routers(id),
		name           VARCHAR(100) NOT NULL,
		phone          VARCHAR(20),
		username       VARCHAR(50) UNIQUE NOT NULL,
		password_hash  VARCHAR(255) NOT NULL,
		status         VARCHAR(20) NOT NULL DEFAULT 'active'
		               CHECK (status IN ('active', 'inactive')),
		voucher_mode   VARCHAR(20) NOT NULL DEFAULT 'mix'
		               CHECK (voucher_mode IN ('mix', 'num', 'alp')),
		voucher_length INT NOT NULL DEFAULT 6,
		voucher_type   VARCHAR(10) NOT NULL DEFAULT 'upp'
		               CHECK (voucher_type IN ('upp', 'up')),
		bill_discount  DECIMAL(15,2) NOT NULL DEFAULT 0,
		created_at     TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at     TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
		deleted_at     TIMESTAMPTZ
	);

	CREATE INDEX IF NOT EXISTS idx_sales_agents_router_id ON sales_agents(router_id);
	CREATE INDEX IF NOT EXISTS idx_sales_agents_username  ON sales_agents(username);
	CREATE INDEX IF NOT EXISTS idx_sales_agents_status    ON sales_agents(status);
	CREATE INDEX IF NOT EXISTS idx_sales_agents_deleted   ON sales_agents(deleted_at);

	CREATE TABLE IF NOT EXISTS sales_profile_prices (
		id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		sales_agent_id UUID NOT NULL REFERENCES sales_agents(id) ON DELETE CASCADE,
		profile_name   VARCHAR(100) NOT NULL,
		base_price     DECIMAL(15,2) NOT NULL DEFAULT 0,
		selling_price  DECIMAL(15,2) NOT NULL DEFAULT 0,
		voucher_length INT,
		is_active      BOOLEAN NOT NULL DEFAULT true,
		created_at     TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(sales_agent_id, profile_name)
	);

	CREATE INDEX IF NOT EXISTS idx_spp_sales_agent_id ON sales_profile_prices(sales_agent_id);

	COMMENT ON TABLE  sales_agents IS 'Agen penjualan voucher hotspot';
	COMMENT ON COLUMN sales_agents.voucher_mode   IS 'mix=alfanumerik, num=numerik, alp=huruf';
	COMMENT ON COLUMN sales_agents.voucher_type   IS 'upp=username=password, up=username berbeda dengan password';
	COMMENT ON COLUMN sales_agents.bill_discount  IS 'Diskon per tagihan pelanggan PPPoE (IDR)';
	COMMENT ON TABLE  sales_profile_prices IS 'Override harga profil hotspot per agen penjualan';
	COMMENT ON COLUMN sales_profile_prices.voucher_length IS 'NULL = pakai default dari sales_agents.voucher_length';
	`)
	return err
}

func down023(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
	DROP TABLE IF EXISTS sales_profile_prices CASCADE;
	DROP TABLE IF EXISTS sales_agents CASCADE;
	`)
	return err
}
