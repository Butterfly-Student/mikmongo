// spec.go - CommandSpec untuk 3 Pipeline Architecture
package collector

import "time"

// PipelineType menentukan pipeline mana yang handle spec ini
type PipelineType string

const (
	PipelineTimeSeries   PipelineType = "time_series"   // Pipeline A -> InfluxDB
	PipelineOperational  PipelineType = "operational"   // Pipeline B -> Redis
	PipelineOnDemand     PipelineType = "on_demand"     // Pipeline C -> Direct
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
	StorageInfluxDB StorageType = "influxdb"  // Pipeline A
	StorageRedis    StorageType = "redis"     // Pipeline B
	StorageNone     StorageType = "none"      // Pipeline C
)

// CommandSpec mendefinisikan satu command untuk collector
type CommandSpec struct {
	// Identifikasi
	Name        string       // Nama command (e.g., "interface_traffic")
	Pipeline    PipelineType // Pipeline A, B, atau C
	
	// RouterOS Command
	Args        []string     // RouterOS command args
	UseFollow   bool         // true untuk follow=yes (streaming)
	Interval    time.Duration // Untuk tier3 ticker interval
	
	// Tier (h untuk Pipeline Operational)
	Tier        Tier         // Tier2 atau Tier3
	
	// Storage Configuration
	Storage     StorageType  // influxdb, redis, none
	
	// InfluxDB Config (Pipeline A)
	Measurement string       // InfluxDB measurement name
	TagFields   []string     // Fields yang jadi tag (dimensions)
	ValueFields []string     // Fields yang jadi value (metrics)
	
	// Redis Config (Pipeline B)
	RedisKey    string       // Redis key prefix (e.g., "ppp:active")
	KeyField    string       // Field untuk hash key (e.g., "name")
	TTL         time.Duration // TTL untuk Tier3
	
	// Enable/Disable
	Enabled     bool         // Apakah spec ini aktif
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

// DefaultSpecs returns default specs untuk semua pipeline
func DefaultSpecs() []CommandSpec {
	specs := make([]CommandSpec, 0)
	
	// ═══════════════════════════════════════════════════════════
	// PIPELINE A: Time-series → InfluxDB
	// Connection pool terpisah, agresif, concurrent
	// ═══════════════════════════════════════════════════════════
	
	// Interface traffic - measurement per interface
	specs = append(specs, CommandSpec{
		Name:        "interface_traffic",
		Pipeline:    PipelineTimeSeries,
		Args:        []string{"/interface/print", "stats", "=interval=1s", "=.proplist=name,rx-byte,tx-byte,rx-packet,tx-packet,rx-drop,tx-drop"},
		UseFollow:   false, // pakai interval
		Storage:     StorageInfluxDB,
		Measurement: "interface_stats",
		TagFields:   []string{"name"},
		ValueFields: []string{"rx-byte", "tx-byte", "rx-packet", "tx-packet", "rx-drop", "tx-drop"},
		Enabled:     true,
	})
	
	// Queue stats - untuk ISP throttling monitoring
	specs = append(specs, CommandSpec{
		Name:        "queue_stats",
		Pipeline:    PipelineTimeSeries,
		Args:        []string{"/queue/simple/print", "stats", "=interval=5s", "=.proplist=name,bytes,packets,dropped,queued-bytes,rate,limit-at"},
		UseFollow:   false,
		Storage:     StorageInfluxDB,
		Measurement: "queue_stats",
		TagFields:   []string{"name"},
		ValueFields: []string{"bytes", "packets", "dropped", "queued-bytes", "rate", "limit-at"},
		Enabled:     true,
	})
	
	// System resources
	specs = append(specs, CommandSpec{
		Name:        "system_resource",
		Pipeline:    PipelineTimeSeries,
		Args:        []string{"/system/resource/print", "=interval=10s", "=.proplist=cpu-load,free-memory,total-memory,uptime"},
		UseFollow:   false,
		Storage:     StorageInfluxDB,
		Measurement: "system_resource",
		TagFields:   []string{},
		ValueFields: []string{"cpu-load", "free-memory", "total-memory", "uptime"},
		Enabled:     true,
	})
	
	// ═══════════════════════════════════════════════════════════
	// PIPELINE B: Operational State → Redis
	// Tier 2: follow=yes untuk event-driven updates
	// Tier 3: Run() + ticker untuk static data
	// ═══════════════════════════════════════════════════════════
	
	// Tier 2: PPP Active - event driven
	specs = append(specs, CommandSpec{
		Name:        "ppp_active",
		Pipeline:    PipelineOperational,
		Args:        []string{"/ppp/active/print", "follow=yes", "=.proplist=name,address,caller-id,uptime,service"},
		UseFollow:   true,
		Tier:        Tier2,
		Storage:     StorageRedis,
		RedisKey:    "ppp:active",
		KeyField:    "name",
		Enabled:     true,
	})
	
	// Tier 2: Hotspot Active - event driven
	specs = append(specs, CommandSpec{
		Name:        "hotspot_active",
		Pipeline:    PipelineOperational,
		Args:        []string{"/ip/hotspot/active/print", "follow=yes", "=.proplist=mac-address,address,user,uptime,bytes-in,bytes-out"},
		UseFollow:   true,
		Tier:        Tier2,
		Storage:     StorageRedis,
		RedisKey:    "hotspot:active",
		KeyField:    "mac-address",
		Enabled:     true,
	})
	
	// Tier 2: Interface status - event driven
	specs = append(specs, CommandSpec{
		Name:        "interface_status",
		Pipeline:    PipelineOperational,
		Args:        []string{"/interface/print", "follow=yes", "=.proplist=name,type,running,disabled"},
		UseFollow:   true,
		Tier:        Tier2,
		Storage:     StorageRedis,
		RedisKey:    "interface:status",
		KeyField:    "name",
		Enabled:     true,
	})
	
	// Tier 3: PPP Secrets - static, fetched periodically
	specs = append(specs, CommandSpec{
		Name:        "ppp_secrets",
		Pipeline:    PipelineOperational,
		Args:        []string{"/ppp/secret/print", "=.proplist=.id,name,service,profile,disabled,comment"},
		UseFollow:   false,
		Tier:        Tier3,
		Interval:    5 * time.Minute,
		Storage:     StorageRedis,
		RedisKey:    "ppp:secrets",
		KeyField:    "name",
		TTL:         10 * time.Minute,
		Enabled:     true,
	})
	
	// Tier 3: PPP Profiles
	specs = append(specs, CommandSpec{
		Name:        "ppp_profiles",
		Pipeline:    PipelineOperational,
		Args:        []string{"/ppp/profile/print", "=.proplist=.id,name,local-address,remote-address,rate-limit"},
		UseFollow:   false,
		Tier:        Tier3,
		Interval:    5 * time.Minute,
		Storage:     StorageRedis,
		RedisKey:    "ppp:profiles",
		KeyField:    "name",
		TTL:         10 * time.Minute,
		Enabled:     true,
	})
	
	// Tier 3: Hotspot Users
	specs = append(specs, CommandSpec{
		Name:        "hotspot_users",
		Pipeline:    PipelineOperational,
		Args:        []string{"/ip/hotspot/user/print", "=.proplist=.id,name,profile,disabled,comment"},
		UseFollow:   false,
		Tier:        Tier3,
		Interval:    5 * time.Minute,
		Storage:     StorageRedis,
		RedisKey:    "hotspot:users",
		KeyField:    "name",
		TTL:         10 * time.Minute,
		Enabled:     true,
	})
	
	// Tier 3: IP Pools
	specs = append(specs, CommandSpec{
		Name:        "ip_pools",
		Pipeline:    PipelineOperational,
		Args:        []string{"/ip/pool/print", "=.proplist=.id,name,ranges"},
		UseFollow:   false,
		Tier:        Tier3,
		Interval:    10 * time.Minute,
		Storage:     StorageRedis,
		RedisKey:    "ip:pools",
		KeyField:    "name",
		TTL:         20 * time.Minute,
		Enabled:     true,
	})
	
	return specs
}

// FilterByPipeline returns specs untuk pipeline tertentu
func FilterByPipeline(specs []CommandSpec, pipeline PipelineType) []CommandSpec {
	var result []CommandSpec
	for _, spec := range specs {
		if spec.Pipeline == pipeline && spec.Enabled {
			result = append(result, spec)
		}
	}
	return result
}

// FilterByTier returns specs untuk tier tertentu (Pipeline B only)
func FilterByTier(specs []CommandSpec, tier Tier) []CommandSpec {
	var result []CommandSpec
	for _, spec := range specs {
		if spec.Pipeline == PipelineOperational && spec.Tier == tier && spec.Enabled {
			result = append(result, spec)
		}
	}
	return result
}
