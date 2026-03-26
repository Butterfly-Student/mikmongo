package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(up029, down029)
}

func up029(ctx context.Context, tx *sql.Tx) error {
	// Seed data moved to internal/seeder/seed_templates.go
	return nil
}

func down029(ctx context.Context, tx *sql.Tx) error {
	return nil
}
