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
//
// NOTE: Tests in this package must NOT use t.Parallel().
// All tests share a single connection pool and use per-test transactions for isolation.
// Running tests in parallel would interleave transactions and cause false failures.
package integration

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	_ "mikmongo/internal/migration"
	"mikmongo/internal/model"
	pkgredis "mikmongo/pkg/redis"
)

// TestConfig holds test configuration
type TestConfig struct {
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       int
}

// LoadTestConfig loads test configuration from environment
func LoadTestConfig() *TestConfig {
	return &TestConfig{
		DBHost:        getEnv("TEST_DB_HOST", "localhost"),
		DBPort:        getEnv("TEST_DB_PORT", "5432"),
		DBUser:        getEnv("TEST_DB_USER", "postgres"),
		DBPassword:    getEnv("TEST_DB_PASSWORD", "postgres"),
		DBName:        getEnv("TEST_DB_NAME", "mikmongo_test"),
		RedisHost:     getEnv("TEST_REDIS_HOST", "localhost"),
		RedisPort:     getEnv("TEST_REDIS_PORT", "6379"),
		RedisPassword: getEnv("TEST_REDIS_PASSWORD", ""),
		RedisDB:       getEnvInt("TEST_REDIS_DB", 15),
	}
}

// Shared connection — opened once in TestMain, reused across all tests.
var (
	sharedDB    *gorm.DB
	sharedSQL   *sql.DB
	sharedRedis *pkgredis.Client
)

// TestMain runs once for the entire test binary.
// It opens the database connection, runs migrations, then hands off to m.Run().
// Each test gets its own transaction via SetupSuite; TestMain closes the connection at the end.
func TestMain(m *testing.M) {
	cfg := LoadTestConfig()
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	var err error
	sharedDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic("failed to connect to test database: " + err.Error())
	}

	sharedSQL, err = sharedDB.DB()
	if err != nil {
		panic("failed to get sql.DB: " + err.Error())
	}

	redisPort, _ := strconv.Atoi(cfg.RedisPort)
	sharedRedis = pkgredis.NewClient(pkgredis.Options{
		Host:     cfg.RedisHost,
		Port:     redisPort,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})
	if err := sharedRedis.Ping(context.Background()); err != nil {
		panic("cannot connect to test Redis: " + err.Error())
	}
	defer sharedRedis.FlushDB(context.Background())
	defer sharedRedis.Close()

	runMigrationsOnce(sharedSQL)

	code := m.Run()
	sharedSQL.Close()
	os.Exit(code)
}

// TestSuite holds shared resources for integration tests
type TestSuite struct {
	DB          *gorm.DB        // transaction-scoped: all test writes are rolled back in Cleanup
	RootDB      *gorm.DB        // non-transactional shared DB (for setup that must be committed)
	SQLDB       *sql.DB         // underlying sql.DB (shared, do NOT close per-test)
	RedisClient *pkgredis.Client // real Redis, DB=15 (isolated test DB)
	Ctx         context.Context
	Config      *TestConfig
}

// SetupSuite begins a transaction on the shared DB and returns a TestSuite scoped to it.
// Every INSERT/UPDATE in the test runs inside this transaction.
// Call defer suite.Cleanup(t) to roll it all back at the end of the test.
func SetupSuite(t *testing.T) *TestSuite {
	t.Helper()
	tx := sharedDB.Begin()
	require.NoError(t, tx.Error, "failed to begin test transaction")

	// Flush test Redis DB before each test to guarantee clean state
	_ = sharedRedis.FlushDB(context.Background())

	sqlDB, _ := sharedDB.DB()
	return &TestSuite{
		DB:          tx,
		RootDB:      sharedDB,
		SQLDB:       sqlDB,
		RedisClient: sharedRedis,
		Ctx:         context.Background(),
		Config:      LoadTestConfig(),
	}
}

// runMigrationsOnce runs all database migrations exactly once (called from TestMain).
func runMigrationsOnce(db *sql.DB) {
	if err := goose.SetDialect("postgres"); err != nil {
		panic("failed to set goose dialect: " + err.Error())
	}
	migrationsDir := "../../internal/migration"
	if err := goose.Up(db, migrationsDir); err != nil {
		panic("failed to run migrations: " + err.Error())
	}
}

// TearDownSuite is a no-op in the transaction-rollback model.
// Connection lifecycle is managed by TestMain.
// Kept for backward compatibility with existing "defer suite.TearDownSuite(t)" call sites.
func (s *TestSuite) TearDownSuite(t *testing.T) {
	// no-op: connection is closed in TestMain after all tests complete
}

// Cleanup rolls back the test transaction, undoing all inserts/updates made during the test.
// This is ~100x faster than TRUNCATE and also reverts sequence_counter increments.
func (s *TestSuite) Cleanup(t *testing.T) {
	if s.DB != nil {
		s.DB.Rollback()
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
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

// createTestUser creates a minimal user record in DB and returns its ID.
// Used for tests that need a valid user UUID for FK constraints.
func createTestUser(t *testing.T, suite *TestSuite) string {
	t.Helper()
	id := uuid.New().String()
	err := suite.DB.WithContext(suite.Ctx).Exec(
		`INSERT INTO users (id, full_name, email, password_hash, role, is_active, created_at, updated_at)
		 VALUES (?, ?, ?, ?, 'admin', true, NOW(), NOW())`,
		id, "Test Admin", id+"@test.com", "$2a$04$placeholder",
	).Error
	require.NoError(t, err)
	return id
}

// directCreateSub creates a subscription directly in the DB,
// bypassing router connectivity (for billing/payment tests).
func directCreateSub(t *testing.T, suite *TestSuite, sub *model.Subscription) {
	t.Helper()
	sub.Status = "pending"
	err := suite.DB.WithContext(suite.Ctx).Create(sub).Error
	require.NoError(t, err)
}

// directActivate sets subscription status=active directly in DB,
// bypassing router connectivity (for billing/payment tests).
func directActivate(t *testing.T, suite *TestSuite, subID string) {
	t.Helper()
	now := time.Now()
	err := suite.DB.WithContext(suite.Ctx).
		Exec("UPDATE subscriptions SET status='active', activated_at=? WHERE id=?", now, subID).Error
	require.NoError(t, err)
}

// directIsolate sets subscription status=isolated directly in DB.
func directIsolate(t *testing.T, suite *TestSuite, subID string) {
	t.Helper()
	err := suite.DB.WithContext(suite.Ctx).
		Exec("UPDATE subscriptions SET status='isolated' WHERE id=?", subID).Error
	require.NoError(t, err)
}

// directSuspend sets subscription status=suspended directly in DB.
func directSuspend(t *testing.T, suite *TestSuite, subID string) {
	t.Helper()
	err := suite.DB.WithContext(suite.Ctx).
		Exec("UPDATE subscriptions SET status='suspended' WHERE id=?", subID).Error
	require.NoError(t, err)
}
