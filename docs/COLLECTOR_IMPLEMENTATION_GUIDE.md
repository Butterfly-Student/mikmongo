# Implementation Guide: 3-Pipeline Collector System

## 📁 File Structure (Final)

```
pkg/mikrotik/
├── mikrotik.go              # ✅ EXISTING - facade client
├── client/
│   └── client.go            # ✅ EXISTING - go-ros client
└── collector/               # 🆕 NEW - 3 Pipeline System
    ├── spec.go              # CommandSpec definitions
    ├── pool/
    │   └── pool.go          # Connection pool per pipeline
    ├── writer/
    │   ├── batch_writer.go  # Shared batch writer
    │   ├── influx_handler.go # InfluxDB writes
    │   └── redis_handler.go  # Redis writes
    ├── pipeline/
    │   ├── time_series/
    │   │   └── collector.go  # Pipeline A → InfluxDB
    │   ├── operational/
    │   │   ├── tier2_collector.go  # Pipeline B Tier 2
    │   │   └── tier3_collector.go  # Pipeline B Tier 3
    │   └── ondemand/
    │       └── runner.go     # Pipeline C direct
    ├── supervisor.go        # Auto-restart supervisor
    └── manager.go           # Main manager
```

## 🔄 3 Pipeline Architecture

### Pipeline A: Time-series → InfluxDB
**Purpose**: Metrics yang perlu histori (traffic, queue, resources)

```go
// Connection pool terpisah (agresif)
poolA := pool.NewConnPool(routerCfg, pool.Config{
    Type:        pool.PoolTimeSeries,
    MaxConns:    5,    // Banyak koneksi untuk concurrent
    MaxCommands: 20,   // RouterOS limit per conn
})

// Collector runs semua specs concurrent
// Fan-in ke single channel
// BatchWriter → InfluxDB (line protocol)
```

**Data Flow**:
```
RouterOS API
    │ /interface/print stats interval=1s
    ▼
[Pool A - Conn 1] ──┐
[Pool A - Conn 2] ──┤
[Pool A - Conn 3] ──┼──► Fan-in Channel
    ...             │
    ▼               ▼
              BatchWriter (100 items or 50ms)
                    │
                    ▼
              InfluxDB (line protocol)
                    │
                    ▼
              Dashboard Charts
```

### Pipeline B: Operational State → Redis  
**Purpose**: Current state (active sessions, static configs)

```go
// Connection pool terpisah (moderate)
poolB := pool.NewConnPool(routerCfg, pool.Config{
    Type:        pool.PoolOperational,
    MaxConns:    5,    // Moderate untuk follow=yes
    MaxCommands: 10,   // Lebih sedikit karena permanen
})
```

**Tier 2 (follow=yes)**: Real-time state changes
```
RouterOS API
    │ /ppp/active/print follow=yes
    ▼
[Pool B - Tier 2 Slot] ──► Redis HSET (last state)
                                    │
                                    ▼
                              Dashboard Tables
```

**Tier 3 (Run + ticker)**: Static data dengan TTL
```
Ticker (every 5m)
    │
    ▼
[Pool B - Tier 3 Slot] ──► /ppp/secret/print (Run, no follow)
    │
    ▼
Redis HSET dengan TTL=10m
    │
    ▼
Dashboard Forms (HGETALL)
```

### Pipeline C: On-demand → Direct
**Purpose**: Langsung ke router, no storage

```go
// Connection pool minimal (on-demand)
poolC := pool.NewConnPool(routerCfg, pool.Config{
    Type:        pool.PoolOnDemand,
    MaxConns:    2,    // Minimal untuk burst
    MaxCommands: 25,
})
```

**Read Operations**:
```
Dashboard Request
    │ POST /api/ping
    ▼
[Pool C] ──► Run("/ping") ──► RouterOS
    │                              │
    └──────────────────────────────┘
                    │
                    ▼
              Response langsung
```

**Write Operations**:
```
Dashboard Action
    │ POST /api/ppp/secrets (add)
    ▼
[Pool C] ──► Run("/ppp/secret/add") ──► RouterOS
    │                                         │
    │                                         ▼
    │                                   Write Success
    │
    ▼
Redis DEL mikrotik:{router}:ppp:secrets
    │
    ▼
(Tier 3 ticker akan refresh dalam 5m)
```

## 🔑 Critical Design Principles

### 1. Pool Isolation (Jangan Digabung!)

```go
// ❌ SALAH - Tier 1 traffic akan delay Tier 2 events
sharedPool := NewConnPool(routerCfg, Config{MaxConns: 10})
// Tier 1 + Tier 2 + Tier 3 pakai pool yang sama

// ✅ BENAR - Terpisah per pipeline
timeSeriesPool := NewConnPool(routerCfg, Config{Type: PoolTimeSeries})   // 5 conn
opPool := NewConnPool(routerCfg, Config{Type: PoolOperational})          // 5 conn  
onDemandPool := NewConnPool(routerCfg, Config{Type: PoolOnDemand})       // 2 conn
```

### 2. Tier 3 Jangan Pakai follow=yes

```go
// ❌ SALAH - Membuang slot permanen untuk data yang jarang berubah
{
    Args: []string{"/ppp/secret/print", "follow=yes"}, // JANGAN!
    Tier: Tier3,
}

// ✅ BENAR - Pakai Run() + ticker
{
    Args: []string{"/ppp/secret/print"}, // Run saja
    Tier: Tier3,
    Interval: 5 * time.Minute,  // Fetch periodic
    TTL: 10 * time.Minute,      // Cache dengan TTL
}
```

### 3. Cache Invalidation Flow

```go
// After write operation (Pipeline C)
func (r *Runner) Write(ctx context.Context, args []string, invalidateKeys []string) error {
    // 1. Execute write
    _, err := conn.RunRaw(ctx, args)
    if err != nil {
        return err
    }
    
    // 2. Invalidate cache (Pipeline B Tier 3)
    redisHandler.Invalidate(ctx, routerID, invalidateKeys)
    
    // 3. Success - Tier 3 ticker akan refresh otomatis
    return nil
}
```

### 4. Shared BatchWriter

```go
// Satu BatchWriter untuk Pipeline A & B
batchWriter := writer.NewBatchWriter(
    writer.Config{
        BatchSize:     100,
        FlushInterval: 50 * time.Millisecond,
    },
    influxHandler,  // Pipeline A
    redisHandler,   // Pipeline B
)

// Pipeline A (Time-series) → influxHandler
// Pipeline B (Operational) → redisHandler
```

## 📊 InfluxDB Schema

### Measurements

```sql
-- Interface traffic (real-time)
measurement: interface_stats
tags: router, name (interface name)
fields: rx_byte, tx_byte, rx_packet, tx_packet, rx_drop, tx_drop
time: timestamp

-- Queue statistics
measurement: queue_stats  
tags: router, name (queue name)
fields: bytes, packets, dropped, queued_bytes, rate, limit_at
time: timestamp

-- System resources
measurement: system_resource
tags: router
fields: cpu_load, free_memory, total_memory, uptime
time: timestamp
```

### Query Examples

```sql
-- Traffic chart (last 1 hour)
SELECT mean(rx_byte) as rx, mean(tx_byte) as tx
FROM interface_stats
WHERE router = 'router-1' AND name = 'ether1'
  AND time > now() - 1h
GROUP BY time(1m)

-- Queue usage
SELECT mean(rate) as rate
FROM queue_stats
WHERE router = 'router-1'
  AND time > now() - 24h
GROUP BY time(5m), name
```

## 📦 Redis Schema

### Keys

```
# Pipeline B Tier 2 - Real-time state (HSET)
mikrotik:{router}:ppp:active:{name} -> hash
mikrotik:{router}:hotspot:active:{mac} -> hash
mikrotik:{router}:interface:status:{name} -> hash

# Pipeline B Tier 3 - Static dengan TTL (HSET)
mikrotik:{router}:ppp:secrets:{name} -> hash (TTL=10m)
mikrotik:{router}:ppp:profiles:{name} -> hash (TTL=10m)
mikrotik:{router}:hotspot:users:{name} -> hash (TTL=10m)
mikrotik:{router}:ip:pools:{name} -> hash (TTL=20m)
```

### Operations

```bash
# Tier 2 - Update state (real-time)
HSET mikrotik:router-1:ppp:active:user1 name user1 address 10.0.0.2

# Tier 3 - Update dengan TTL
HSET mikrotik:router-1:ppp:secrets:user1 name user1 profile default
EXPIRE mikrotik:router-1:ppp:secrets:user1 600

# Dashboard read
HGETALL mikrotik:router-1:ppp:active

# Pipeline C - Invalidate after write
DEL mikrotik:router-1:ppp:secrets:newuser
```

## 🚀 Usage Example

```go
// 1. Setup handlers
influxHandler, _ := writer.NewInfluxHandler(writer.InfluxConfig{
    URL:    "http://localhost:8086",
    Token:  "my-token",
    Org:    "my-org",
    Bucket: "mikrotik",
})

redisHandler, _ := writer.NewRedisHandler(writer.RedisConfig{
    Addr:   "localhost:6379",
    DB:     0,
    Prefix: "mikrotik",
})

// 2. Create shared BatchWriter
batchWriter := writer.NewBatchWriter(
    writer.DefaultConfig(),
    influxHandler,
    redisHandler,
)
batchWriter.Start()
defer batchWriter.Stop()

// 3. Create Supervisor (Pipeline A & B)
supervisor, _ := collector.NewSupervisor(
    "router-1",
    routerCfg,                    // MikroTik connection config
    timeSeriesSpecs,              // Pipeline A specs
    operationalSpecs,             // Pipeline B specs
    batchWriter,
    collector.DefaultSupervisorConfig(),
)
supervisor.Start()
defer supervisor.Stop()

// 4. Create On-demand Runner (Pipeline C)
onDemandPool, _ := pool.NewConnPool(routerCfg, pool.DefaultConfig(pool.PoolOnDemand))
runner := ondemand.NewRunner("router-1", onDemandPool, redisHandler)

// 5. Execute on-demand commands
// Read
results, _ := runner.Run(ctx, []string{"/ping", "=address=8.8.8.8", "=count=3"})

// Write dengan cache invalidation
runner.Write(ctx, 
    []string{"/ppp/secret/add", "=name=newuser", "=password=secret", "=service=pppoe"},
    []string{"ppp:secrets"},  // Invalidate cache keys
)

// 6. Query data
// InfluxDB (historical)
// ... query untuk charts

// Redis (current state)
users, _ := redisHandler.GetAll(ctx, "router-1", "ppp:secrets")
```

## 🛡️ Supervisor Auto-Restart

```go
// Supervisor monitors collectors
supervisor := collector.NewSupervisor(...)

// Health check setiap 30 detik
// Jika pool unhealthy atau collector mati:
// 1. Close old pool
// 2. Create new pool  
// 3. Restart collector
// 4. Continue operation

// Manual restart jika perlu
supervisor.RestartCollector("time_series")
supervisor.RestartCollector("operational")
```

## 📈 Scaling

### Multiple Routers
```go
manager := collector.NewManager(config)

for _, router := range routers {
    manager.AddRouter(router.ID, router.Config, specs)
}

manager.StartAll()
```

### Resource Limits
```go
// Per router:
// - Pipeline A: 5 conn (time-series)
// - Pipeline B: 5 conn (operational)
// - Pipeline C: 2 conn (on-demand)
// Total: 12 connections per router

// Untuk 100 routers:
// - 500 conn untuk time-series
// - 500 conn untuk operational
// - 200 conn untuk on-demand (burst)
```

## 🧪 Testing Strategy

### Unit Tests
```bash
go test ./pkg/mikrotik/collector/pool/...
go test ./pkg/mikrotik/collector/writer/...
```

### Integration Test
```bash
# Run dengan real router
go run cmd/collector/main.go \
    -host=192.168.27.1 \
    -influx=http://localhost:8086 \
    -redis=localhost:6379
```

### Load Test
```bash
# Monitor:
# - Pool usage
# - BatchWriter queue depth
# - InfluxDB write performance
# - Redis memory usage
```
