// spec/types.go - Shared types untuk collector (mencegah import cycle)
package spec

import "time"

// PipelineType menentukan pipeline mana yang handle spec ini
type PipelineType string

const (
	PipelineTimeSeries  PipelineType = "time_series"  // Pipeline A -> InfluxDB
	PipelineOperational PipelineType = "operational"  // Pipeline B -> Redis
	PipelineOnDemand    PipelineType = "on_demand"    // Pipeline C -> Direct
)

// Tier dalam Pipeline Operational
type Tier int

const (
	Tier2 Tier = 2 // Event-driven: follow=yes
	Tier3 Tier = 3 // Static: Run() + ticker + TTL
)

// StorageType menentukan storage backend
type StorageType string

const (
	StorageInfluxDB StorageType = "influxdb" // Pipeline A
	StorageRedis    StorageType = "redis"    // Pipeline B
	StorageNone     StorageType = "none"     // Pipeline C
)

// CommandSpec mendefinisikan satu command untuk collector
type CommandSpec struct {
	// Identifikasi
	Name     string       // Nama command (e.g., "interface_traffic")
	Pipeline PipelineType // Pipeline A, B, atau C

	// RouterOS Command
	Args      []string      // RouterOS command args
	UseFollow bool          // true untuk follow=yes (streaming)
	Interval  time.Duration // Untuk tier3 ticker interval

	// Tier (hanya untuk Pipeline Operational)
	Tier Tier // Tier2 atau Tier3

	// Storage Configuration
	Storage StorageType // influxdb, redis, none

	// InfluxDB Config (Pipeline A)
	Measurement string   // InfluxDB measurement name
	TagFields   []string // Fields yang jadi tag (dimensions)
	ValueFields []string // Fields yang jadi value (metrics)

	// Redis Config (Pipeline B)
	RedisKey string        // Redis key prefix (e.g., "ppp:active")
	KeyField string        // Field untuk hash key (e.g., "name")
	TTL      time.Duration // TTL untuk Tier3

	// Enable/Disable
	Enabled bool // Apakah spec ini aktif
	
	// Legacy fields untuk backward compatibility
	ToInflux bool // Legacy: true untuk InfluxDB
}

// IsTimeSeries returns true jika spec untuk Pipeline A
func (s CommandSpec) IsTimeSeries() bool {
	return s.Pipeline == PipelineTimeSeries
}

// IsOperational returns true jika spec untuk Pipeline B
func (s CommandSpec) IsOperational() bool {
	return s.Pipeline == PipelineOperational
}

// IsOnDemand returns true jika spec untuk Pipeline C
func (s CommandSpec) IsOnDemand() bool {
	return s.Pipeline == PipelineOnDemand
}

// IsStreaming returns true jika menggunakan follow=yes
func (s CommandSpec) IsStreaming() bool {
	return s.UseFollow
}
