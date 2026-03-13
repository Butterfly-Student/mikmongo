package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(up002, down002)
}

// up002 creates the customers table for ISP customer identity data.
func up002(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec(`
	-- ============================================
	-- CUSTOMERS TABLE (identitas pelanggan)
	-- ============================================
	CREATE TABLE IF NOT EXISTS customers (
		id            UUID      PRIMARY KEY DEFAULT gen_random_uuid(),
		customer_code VARCHAR(50) UNIQUE NOT NULL,

		-- Identitas
		full_name  VARCHAR(100) NOT NULL,
		email      VARCHAR(100) UNIQUE,
		phone      VARCHAR(20)  NOT NULL,
		id_card_number VARCHAR(30),

		-- Alamat
		address   TEXT,
		latitude  DECIMAL(10,8),
		longitude DECIMAL(11,8),

		-- Status Account
		is_active BOOLEAN DEFAULT true,

		-- Portal Self-Service
		portal_password_hash VARCHAR(255),
		portal_last_login TIMESTAMPTZ,

		notes TEXT,
		tags  JSONB,

		created_by UUID REFERENCES users(id) ON DELETE SET NULL,
		updated_by UUID REFERENCES users(id) ON DELETE SET NULL,
		created_at TIMESTAMPTZ   DEFAULT CURRENT_TIMESTAMP NOT NULL,
		updated_at TIMESTAMPTZ   DEFAULT CURRENT_TIMESTAMP NOT NULL,
		deleted_at TIMESTAMPTZ
	);

	CREATE INDEX IF NOT EXISTS idx_customers_code     ON customers(customer_code);
	CREATE INDEX IF NOT EXISTS idx_customers_is_active ON customers(is_active) WHERE deleted_at IS NULL;
	CREATE INDEX IF NOT EXISTS idx_customers_phone    ON customers(phone);
	CREATE INDEX IF NOT EXISTS idx_customers_email    ON customers(email)   WHERE email IS NOT NULL;
	CREATE INDEX IF NOT EXISTS idx_customers_deleted  ON customers(deleted_at);

	COMMENT ON TABLE  customers                IS 'Data identitas pelanggan ISP (konfigurasi layanan di tabel subscriptions)';
	COMMENT ON COLUMN customers.is_active      IS 'true=account aktif, false=account non-aktif (soft delete)';
	COMMENT ON COLUMN customers.created_by     IS 'Admin yang membuat data pelanggan';
	COMMENT ON COLUMN customers.updated_by     IS 'Admin yang terakhir mengupdate data pelanggan';
	`)
	return err
}

func down002(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE IF EXISTS customers CASCADE;`)
	return err
}
