//go:build integration

// Package integration contains integration tests using real PostgreSQL and real Mikrotik devices
// These tests require:
// 1. Docker Desktop running (for testcontainers PostgreSQL)
// 2. A real Mikrotik router accessible via environment variables:
//   - TEST_MIKROTIK_HOST (default: 192.168.88.1)
//   - TEST_MIKROTIK_PORT (default: 8728)
//   - TEST_MIKROTIK_USER (default: admin)
//   - TEST_MIKROTIK_PASS (required for Mikrotik tests)
//
// Run tests with: go test -v -tags=integration ./tests/integration/...
package integration

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// TestConfig holds test configuration
type TestConfig struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

// LoadTestConfig loads test configuration from environment
func LoadTestConfig() *TestConfig {
	return &TestConfig{
		DBHost:     getEnv("TEST_DB_HOST", "localhost"),
		DBPort:     getEnv("TEST_DB_PORT", "5432"),
		DBUser:     getEnv("TEST_DB_USER", "postgres"),
		DBPassword: getEnv("TEST_DB_PASSWORD", "postgres"),
		DBName:     getEnv("TEST_DB_NAME", "mikmongo_test"),
	}
}

// TestSuite holds shared resources for integration tests
type TestSuite struct {
	DB     *gorm.DB
	SQLDB  *sql.DB
	Ctx    context.Context
	Config *TestConfig
}

// SetupSuite initializes the test suite with PostgreSQL connection
func SetupSuite(t *testing.T) *TestSuite {
	ctx := context.Background()
	cfg := LoadTestConfig()

	// Build connection string
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	// Connect with GORM
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err, "Failed to connect to database")

	// Get underlying sql.DB
	sqlDB, err := db.DB()
	require.NoError(t, err, "Failed to get sql.DB")

	// Run migrations
	runMigrations(t, sqlDB)

	return &TestSuite{
		DB:     db,
		SQLDB:  sqlDB,
		Ctx:    ctx,
		Config: cfg,
	}
}

// runMigrations runs all database migrations
func runMigrations(t *testing.T, db *sql.DB) {
	// Set goose dialect
	if err := goose.SetDialect("postgres"); err != nil {
		t.Fatalf("Failed to set goose dialect: %v", err)
	}

	// Run migrations from internal/migration
	migrationsDir := "../../internal/migration"
	if err := goose.Up(db, migrationsDir); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}
}

// TearDownSuite cleans up resources
func (s *TestSuite) TearDownSuite(t *testing.T) {
	if s.SQLDB != nil {
		_ = s.SQLDB.Close()
	}
}

// Cleanup truncates all tables between tests (except seed data tables)
func (s *TestSuite) Cleanup(t *testing.T) {
	tables := []string{
		"users", "customers", "mikrotik_routers", "bandwidth_profiles",
		"subscriptions", "invoices", "invoice_items", "payments",
		"payment_allocations", "audit_logs",
		"customer_registrations",
	}

	for _, table := range tables {
		if err := s.DB.WithContext(s.Ctx).Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table)).Error; err != nil {
			t.Logf("Warning: failed to truncate table %s: %v", table, err)
		}
	}

	// Reset sequence counters to 0 but keep the records
	if err := s.DB.WithContext(s.Ctx).Exec("UPDATE sequence_counters SET last_number = 0").Error; err != nil {
		t.Logf("Warning: failed to reset sequence counters: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// WithTimeout creates a context with timeout for tests
func WithTimeout(t *testing.T, timeout time.Duration) context.Context {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	t.Cleanup(cancel)
	return ctx
}
