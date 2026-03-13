package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(up007, down007)
}

// up007 creates the customer_registrations table for new customer registration requests.
func up007(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS customer_registrations (
			id                   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			full_name            VARCHAR(100) NOT NULL,
			email                VARCHAR(100),
			phone                VARCHAR(20)  NOT NULL,
			address              TEXT,
			latitude             DECIMAL(10,8),
			longitude            DECIMAL(11,8),
			service_type         VARCHAR(20)  NOT NULL
			                     CHECK (service_type IN ('pppoe', 'hotspot')),
			notes                TEXT,
			bandwidth_profile_id UUID REFERENCES bandwidth_profiles(id) ON DELETE SET NULL,
			status               VARCHAR(20)  NOT NULL DEFAULT 'pending'
			                     CHECK (status IN ('pending', 'approved', 'rejected')),
			rejection_reason     TEXT,
			approved_by          UUID REFERENCES users(id) ON DELETE SET NULL,
			approved_at          TIMESTAMPTZ,
			customer_id          UUID REFERENCES customers(id) ON DELETE SET NULL,
			created_at           TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
			updated_at           TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
			deleted_at           TIMESTAMPTZ
		);
		-- Optimized indexes with partial conditions for soft delete
		CREATE INDEX IF NOT EXISTS idx_customer_registrations_status 
			ON customer_registrations(status) WHERE deleted_at IS NULL;
		CREATE INDEX IF NOT EXISTS idx_customer_registrations_deleted_at 
			ON customer_registrations(deleted_at);
		CREATE INDEX IF NOT EXISTS idx_customer_registrations_customer_id 
			ON customer_registrations(customer_id) WHERE deleted_at IS NULL;
		CREATE INDEX IF NOT EXISTS idx_customer_registrations_service_type
			ON customer_registrations(service_type) WHERE deleted_at IS NULL;
		CREATE INDEX IF NOT EXISTS idx_customer_registrations_bandwidth_profile_id
			ON customer_registrations(bandwidth_profile_id) WHERE bandwidth_profile_id IS NOT NULL AND deleted_at IS NULL;
	`)
	return err
}

func down007(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `DROP TABLE IF EXISTS customer_registrations;`)
	return err
}
