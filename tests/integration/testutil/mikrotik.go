// Package testutil provides test utilities for integration tests
package testutil

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	mikrotik "github.com/Butterfly-Student/go-ros"
	"github.com/Butterfly-Student/go-ros/client"
)

// MikrotikConfig holds configuration for real Mikrotik connection
type MikrotikConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	UseTLS   bool
}

// LoadMikrotikConfig loads Mikrotik configuration from environment variables
func LoadMikrotikConfig() *MikrotikConfig {
	return &MikrotikConfig{
		Host:     getEnv("TEST_MIKROTIK_HOST", "192.168.233.1"),
		Port:     getEnvInt("TEST_MIKROTIK_PORT", 8728),
		Username: getEnv("TEST_MIKROTIK_USER", "admin"),
		Password: getEnv("TEST_MIKROTIK_PASS", "r00t"),
		UseTLS:   getEnvBool("TEST_MIKROTIK_TLS", false),
	}
}

// NewMikrotikClient creates a new Mikrotik client for testing
func NewMikrotikClient(cfg *MikrotikConfig) (*mikrotik.Client, error) {
	clientCfg := client.Config{
		Host:     cfg.Host,
		Port:     cfg.Port,
		Username: cfg.Username,
		Password: cfg.Password,
		UseTLS:   cfg.UseTLS,
	}

	return mikrotik.NewClient(clientCfg)
}

// SkipIfNoMikrotik skips the test if no Mikrotik connection is configured
func SkipIfNoMikrotik(t *testing.T) {
	if os.Getenv("TEST_MIKROTIK_HOST") == "" {
		t.Skip("Skipping test: TEST_MIKROTIK_HOST not set")
	}
}

// getEnv gets environment variable with default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt gets environment variable as int with default value
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		var result int
		fmt.Sscanf(value, "%d", &result)
		return result
	}
	return defaultValue
}

// getEnvBool gets environment variable as bool with default value
func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		return value == "true" || value == "1" || value == "yes"
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
