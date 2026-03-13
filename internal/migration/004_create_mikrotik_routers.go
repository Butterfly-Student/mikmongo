package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(up004, down004)
}

// up004 creates the mikrotik_routers table for managed MikroTik routers.
func up004(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec(`
	CREATE TABLE IF NOT EXISTS mikrotik_routers (
		id                UUID      PRIMARY KEY DEFAULT gen_random_uuid(),
		name              VARCHAR(100) NOT NULL,
		address           VARCHAR(100) NOT NULL,
		area              VARCHAR(100),
		api_port          INTEGER   DEFAULT 8728,
		rest_port         INTEGER   DEFAULT 80,
		username          VARCHAR(100) NOT NULL,
		password_encrypted TEXT NOT NULL,
		use_ssl           BOOLEAN   DEFAULT false,
		router_os_version VARCHAR(50),
		identity          VARCHAR(255),
		is_master         BOOLEAN   DEFAULT false,
		is_active         BOOLEAN   DEFAULT true,
		status            VARCHAR(20)  DEFAULT 'unknown'
		                  CHECK (status IN ('online', 'offline', 'unknown')),
		last_seen_at      TIMESTAMPTZ,
		last_ping         TIMESTAMPTZ,
		notes             TEXT,
		created_at        TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
		updated_at        TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
		deleted_at        TIMESTAMPTZ
	);

	CREATE INDEX IF NOT EXISTS idx_mikrotik_routers_name      ON mikrotik_routers(name);
	CREATE INDEX IF NOT EXISTS idx_mikrotik_routers_address   ON mikrotik_routers(address);
	CREATE INDEX IF NOT EXISTS idx_mikrotik_routers_is_active ON mikrotik_routers(is_active)  WHERE deleted_at IS NULL;
	CREATE INDEX IF NOT EXISTS idx_mikrotik_routers_status    ON mikrotik_routers(status)      WHERE deleted_at IS NULL;
	CREATE INDEX IF NOT EXISTS idx_mikrotik_routers_area      ON mikrotik_routers(area)        WHERE area IS NOT NULL AND deleted_at IS NULL;
	CREATE INDEX IF NOT EXISTS idx_mikrotik_routers_deleted   ON mikrotik_routers(deleted_at);

	COMMENT ON TABLE  mikrotik_routers                   IS 'Daftar router MikroTik yang dikelola sistem';
	COMMENT ON COLUMN mikrotik_routers.area               IS 'Wilayah/zona cakupan router (mis: Timur, Barat)';
	COMMENT ON COLUMN mikrotik_routers.is_master          IS 'TRUE = router utama untuk area tersebut';
	COMMENT ON COLUMN mikrotik_routers.password_encrypted IS 'Password terenkripsi AES-256 oleh aplikasi';
	COMMENT ON COLUMN mikrotik_routers.status             IS 'Status koneksi: online/offline/unknown (hasil health check)';
	`)
	return err
}

func down004(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE IF EXISTS mikrotik_routers CASCADE;`)
	return err
}

