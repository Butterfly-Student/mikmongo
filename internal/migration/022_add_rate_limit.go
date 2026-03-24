package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(up022, down022)
}

func up022(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec(`
	ALTER TABLE bandwidth_profiles
		ADD COLUMN IF NOT EXISTS rate_limit VARCHAR(100);

	COMMENT ON COLUMN bandwidth_profiles.rate_limit IS
		'MikroTik rate-limit string (contoh: 10M/10M atau 10240k/10240k). '
		'Jika kosong, dihitung otomatis dari upload_speed/download_speed dalam kbps.';
	`)
	return err
}

func down022(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec(`ALTER TABLE bandwidth_profiles DROP COLUMN IF EXISTS rate_limit;`)
	return err
}
