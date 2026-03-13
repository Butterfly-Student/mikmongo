package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(up006, down006)
}

// up006 creates the subscriptions table linking customers to service plans and routers.
func up006(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec(`
	-- ============================================
	-- SUBSCRIPTIONS TABLE
	-- ============================================
	CREATE TABLE IF NOT EXISTS subscriptions (
		id           UUID      PRIMARY KEY DEFAULT gen_random_uuid(),
		customer_id  UUID      NOT NULL REFERENCES customers(id) ON DELETE RESTRICT,
		plan_id      UUID      NOT NULL REFERENCES bandwidth_profiles(id) ON DELETE RESTRICT,
		router_id    UUID      NOT NULL REFERENCES mikrotik_routers(id) ON DELETE RESTRICT,

		-- Kredensial di MikroTik
		username VARCHAR(100) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL,

		-- Network Config
		static_ip   VARCHAR(45),
		gateway     VARCHAR(15),
		mac_address VARCHAR(17),

		-- Status & Periode
		status       VARCHAR(20) NOT NULL
		             CHECK (status IN ('pending', 'active', 'suspended', 'isolated', 'expired', 'terminated'))
		             DEFAULT 'pending',
		activated_at TIMESTAMPTZ,
		expiry_date  DATE,

		-- Billing & Isolation Config
		billing_day       INTEGER CHECK (billing_day >= 1 AND billing_day <= 31),
		auto_isolate      BOOLEAN DEFAULT true,
		grace_period_days INTEGER,

		-- Alasan suspend/terminate
		suspend_reason   TEXT,
		terminated_at    TIMESTAMPTZ,

		-- Profil sebelumnya (untuk restore setelah isolasi)
		previous_plan_id UUID REFERENCES bandwidth_profiles(id) ON DELETE SET NULL,

		notes      TEXT,
		created_by UUID REFERENCES users(id) ON DELETE SET NULL,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
		deleted_at TIMESTAMPTZ
	);

	-- Core indexes with soft delete filtering
	CREATE INDEX IF NOT EXISTS idx_subscriptions_customer    ON subscriptions(customer_id)   WHERE deleted_at IS NULL;
	CREATE INDEX IF NOT EXISTS idx_subscriptions_plan        ON subscriptions(plan_id)        WHERE deleted_at IS NULL;
	CREATE INDEX IF NOT EXISTS idx_subscriptions_router      ON subscriptions(router_id)      WHERE deleted_at IS NULL;
	CREATE INDEX IF NOT EXISTS idx_subscriptions_username    ON subscriptions(username)       WHERE deleted_at IS NULL;
	CREATE INDEX IF NOT EXISTS idx_subscriptions_status      ON subscriptions(status)         WHERE deleted_at IS NULL;
	CREATE INDEX IF NOT EXISTS idx_subscriptions_expiry      ON subscriptions(expiry_date)    WHERE deleted_at IS NULL;
	CREATE INDEX IF NOT EXISTS idx_subscriptions_deleted     ON subscriptions(deleted_at);
	
	-- Composite indexes untuk query yang sering digunakan
	CREATE INDEX IF NOT EXISTS idx_subscriptions_customer_status ON subscriptions(customer_id, status) WHERE deleted_at IS NULL;

	COMMENT ON TABLE  subscriptions              IS 'Layanan aktif pelanggan — jembatan antara customer, plan, dan router MikroTik';
	COMMENT ON COLUMN subscriptions.username      IS 'PPP username di MikroTik (harus unik per router)';
	COMMENT ON COLUMN subscriptions.previous_plan_id IS 'Plan sebelum isolasi, dikembalikan saat isolasi dicabut';
	`)
	return err
}

func down006(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE IF EXISTS subscriptions CASCADE;`)
	return err
}
