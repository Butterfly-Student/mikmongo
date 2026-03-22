package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(up021, down021)
}

func up021(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
		ALTER TABLE payments
		ADD COLUMN IF NOT EXISTS gateway_payment_url TEXT;

		COMMENT ON COLUMN payments.gateway_payment_url IS 'URL halaman pembayaran dari payment gateway (xendit, dll)';
	`)
	return err
}

func down021(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
		ALTER TABLE payments DROP COLUMN IF EXISTS gateway_payment_url;
	`)
	return err
}
