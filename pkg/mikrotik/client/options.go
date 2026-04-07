package client

import "time"

// Config holds connection parameters for a single MikroTik router.
type Config struct {
	Host              string
	Port              int
	Username          string
	Password          string
	UseTLS            bool
	ReconnectInterval time.Duration
	Timeout           time.Duration // per-command timeout (default 10s)
}
