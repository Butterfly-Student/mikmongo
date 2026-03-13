package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(up019, down019)
}

// up019 creates the audit_logs table for tracking system changes.
func up019(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec(`
	-- ============================================
	-- AUDIT LOGS TABLE
	-- ============================================
	CREATE TABLE IF NOT EXISTS audit_logs (
		id          UUID     PRIMARY KEY DEFAULT gen_random_uuid(),
		admin_id    UUID     REFERENCES users(id) ON DELETE SET NULL,
		action      VARCHAR(100) NOT NULL,
		entity_type VARCHAR(50)  NOT NULL,
		entity_id   UUID         NOT NULL,
		old_value   JSONB,
		new_value   JSONB,
		ip_address  VARCHAR(45),
		user_agent  VARCHAR(255),
		notes       TEXT,
		created_at  TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
	);

	CREATE INDEX IF NOT EXISTS idx_audit_logs_admin      ON audit_logs(admin_id);
	CREATE INDEX IF NOT EXISTS idx_audit_logs_entity     ON audit_logs(entity_type, entity_id);
	CREATE INDEX IF NOT EXISTS idx_audit_logs_action     ON audit_logs(action);
	CREATE INDEX IF NOT EXISTS idx_audit_logs_created_at ON audit_logs(created_at DESC);

	COMMENT ON TABLE  audit_logs          IS 'Audit trail semua perubahan penting di sistem';
	COMMENT ON COLUMN audit_logs.admin_id  IS 'NULL = dilakukan oleh sistem/scheduler, bukan admin';
	COMMENT ON COLUMN audit_logs.old_value IS 'Snapshot data sebelum perubahan (JSON)';
	COMMENT ON COLUMN audit_logs.new_value IS 'Snapshot data sesudah perubahan (JSON)';
	`)
	return err
}

func down019(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE IF EXISTS audit_logs CASCADE;`)
	return err
}
