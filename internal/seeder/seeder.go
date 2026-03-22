package seeder

import (
	"context"
	"database/sql"
	"log"
)

// Config holds minimal seed configuration.
type Config struct {
	EncryptionKey string // JWT_SECRET, used for AES-GCM encryption of router passwords
}

// Seeder runs idempotent seed functions against the database.
type Seeder struct {
	db  *sql.DB
	cfg Config
}

// New creates a new Seeder.
func New(db *sql.DB, cfg Config) *Seeder {
	return &Seeder{db: db, cfg: cfg}
}

// Run executes all seed functions in order. All operations are idempotent.
func (s *Seeder) Run(ctx context.Context) error {
	seeds := []struct {
		name string
		fn   func(context.Context) error
	}{
		{"users", s.seedUsers},
		{"routers", s.seedRouters},
		{"customers", s.seedCustomers},
		{"casbin", s.seedCasbin},
		{"system_settings", s.seedSystemSettings},
		{"message_templates", s.seedMessageTemplates},
		{"sequence_counters", s.seedSequenceCounters},
	}

	for _, seed := range seeds {
		if err := seed.fn(ctx); err != nil {
			return err
		}
		log.Printf("[SEED] %s OK", seed.name)
	}

	return nil
}
