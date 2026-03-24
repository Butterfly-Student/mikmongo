package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(up025, down025)
}

func up025(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
	ALTER TABLE sales_agents
		ADD COLUMN IF NOT EXISTS billing_cycle VARCHAR(20) NOT NULL DEFAULT 'monthly'
		    CHECK (billing_cycle IN ('weekly', 'monthly')),
		ADD COLUMN IF NOT EXISTS billing_day   INT NOT NULL DEFAULT 1
		    CHECK (billing_day >= 1 AND billing_day <= 31);

	COMMENT ON COLUMN sales_agents.billing_cycle IS 'Siklus invoice: weekly atau monthly';
	COMMENT ON COLUMN sales_agents.billing_day   IS 'Hari penagihan: untuk monthly = hari ke-N (1-28), untuk weekly = hari dalam minggu (1=Senin, 7=Minggu)';
	`)
	return err
}

func down025(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
	ALTER TABLE sales_agents
		DROP COLUMN IF EXISTS billing_cycle,
		DROP COLUMN IF EXISTS billing_day;
	`)
	return err
}
