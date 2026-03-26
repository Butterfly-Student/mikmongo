package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() { goose.AddMigrationContext(up020, down020) }

func up020(ctx context.Context, tx *sql.Tx) error {
	// Seed data moved to internal/seeder/ (seed_settings.go, seed_templates.go, seed_sequences.go)
	return nil
}

func down020(ctx context.Context, tx *sql.Tx) error {
	return nil
}
