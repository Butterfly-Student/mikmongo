package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(up005, down005)
}

// up005 creates the bandwidth_profiles table for ISP service packages.
func up005(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec(`
	-- ============================================
	-- BANDWIDTH PROFILES TABLE
	-- ============================================
	CREATE TABLE IF NOT EXISTS bandwidth_profiles (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			
			-- Router scope (multi-router support)
			router_id UUID NOT NULL REFERENCES mikrotik_routers(id),
			
			-- Display info (untuk UI aplikasi)
			profile_code VARCHAR(50) NOT NULL,
			name VARCHAR(100) NOT NULL,
			description TEXT,

			-- Bandwidth info (untuk display & billing)
			download_speed BIGINT NOT NULL,
			upload_speed BIGINT NOT NULL,
			
			-- Pricing
			price_monthly DECIMAL(12,2) NOT NULL,
			tax_rate DECIMAL(5,4) DEFAULT 0.11,
			billing_cycle VARCHAR(20) DEFAULT 'monthly'
				CHECK (billing_cycle IN ('daily', 'weekly', 'monthly', 'yearly')),
			billing_day   INTEGER,
			
		-- Status
		is_active BOOLEAN DEFAULT true,
		is_visible BOOLEAN DEFAULT true,
		sort_order INTEGER DEFAULT 0,

		-- Grace period & isolation
		grace_period_days INTEGER NOT NULL DEFAULT 3,
		isolate_profile_name VARCHAR(100),

		-- Timestamps
			created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
			updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
			deleted_at TIMESTAMPTZ
	);

	-- Indexes (optimized for common queries)
	CREATE INDEX IF NOT EXISTS idx_bandwidth_profiles_router         ON bandwidth_profiles(router_id);
	CREATE INDEX IF NOT EXISTS idx_bandwidth_profiles_is_active    ON bandwidth_profiles(is_active) WHERE deleted_at IS NULL;
	CREATE INDEX IF NOT EXISTS idx_bandwidth_profiles_is_visible   ON bandwidth_profiles(is_visible) WHERE deleted_at IS NULL AND is_active = true;
	CREATE INDEX IF NOT EXISTS idx_bandwidth_profiles_code         ON bandwidth_profiles(profile_code);
	CREATE INDEX IF NOT EXISTS idx_bandwidth_profiles_deleted      ON bandwidth_profiles(deleted_at);
	
	-- Unique constraint: profile name per router
	CREATE UNIQUE INDEX IF NOT EXISTS idx_bandwidth_profiles_router_name 
		ON bandwidth_profiles(router_id, name) WHERE deleted_at IS NULL;
	
	-- Unique constraint: profile code per router  
	CREATE UNIQUE INDEX IF NOT EXISTS idx_bandwidth_profiles_router_code
		ON bandwidth_profiles(router_id, profile_code) WHERE deleted_at IS NULL;

	COMMENT ON TABLE  bandwidth_profiles                IS 'Paket bandwidth untuk layanan PPPoE ISP';
	COMMENT ON COLUMN bandwidth_profiles.billing_cycle   IS 'Siklus tagihan: daily/weekly/monthly/yearly';
	COMMENT ON COLUMN bandwidth_profiles.download_speed  IS 'Kecepatan download dalam kbps (untuk display & billing)';
	COMMENT ON COLUMN bandwidth_profiles.upload_speed    IS 'Kecepatan upload dalam kbps (untuk display & billing)';
	`)
	return err
}

func down005(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE IF EXISTS bandwidth_profiles CASCADE;`)
	return err
}
