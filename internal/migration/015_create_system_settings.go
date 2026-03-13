package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(up015, down015)
}

// up015 creates the system_settings table for application configuration.
func up015(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec(`
	-- ============================================
	-- SYSTEM SETTINGS TABLE
	-- ============================================
	CREATE TABLE IF NOT EXISTS system_settings (
		id           UUID     PRIMARY KEY DEFAULT gen_random_uuid(),
		group_name   VARCHAR(50) NOT NULL,
		key_name     VARCHAR(100) NOT NULL,
		value        TEXT,
		type         VARCHAR(20) CHECK (type IN ('string', 'integer', 'boolean', 'json', 'password')) DEFAULT 'string',
		label        VARCHAR(150),
		description  TEXT,
		is_encrypted BOOLEAN DEFAULT false,
		is_public    BOOLEAN DEFAULT false,
		updated_at   TIMESTAMPTZ,
		updated_by   UUID REFERENCES users(id) ON DELETE SET NULL,
		UNIQUE (group_name, key_name)
	);

	CREATE INDEX IF NOT EXISTS idx_system_settings_group ON system_settings(group_name);
	CREATE INDEX IF NOT EXISTS idx_system_settings_public ON system_settings(is_public) WHERE is_public = true;

	COMMENT ON TABLE  system_settings              IS 'Konfigurasi sistem dalam format key-value per group';
	COMMENT ON COLUMN system_settings.is_encrypted IS 'jika TRUE, value disimpan terenkripsi AES-256';
	COMMENT ON COLUMN system_settings.is_public    IS 'jika TRUE, bisa dibaca tanpa auth (untuk frontend config)';
	`)
	return err
}

func down015(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE IF EXISTS system_settings CASCADE;`)
	return err
}
