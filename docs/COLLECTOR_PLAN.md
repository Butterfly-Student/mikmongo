# Planning: MikroTik 4-Tier Collector System

## 🎯 Tujuan
Implementasi data collector untuk MikroTik ISP dengan 4 tier system:
- **Tier 1**: Real-time streaming (high freq) - traffic, queue
- **Tier 2**: Event-driven (medium freq) - logs, active sessions
- **Tier 3**: Static data (cache TTL) - configs, profiles
- **Tier 4**: Write operations - dengan cache invalidation

## 📁 Struktur File Baru (tidak mengganggu existing)

```
pkg/mikrotik/
├── mikrotik.go              # existing - tidak diubah
├── collector/               # NEW: collector package
│   ├── specs.go            # CommandSpec definitions
│   ├── tier.go             # Tier constants & types
│   ├── pool.go             # Connection pool management
│   ├── cache.go            # Redis cache wrapper
│   ├── collector.go        # Main collector orchestrator
│   ├── worker.go           # Worker untuk tiap tier
│   └── router.go           # Router instance manager
├── cmd/collector/          # NEW: testing command
│   └── main.go             # Test runner
```

## 🔌 Integration Point dengan Existing Code

### 1. Menggunakan Client Existing
```go
// collector menggunakan Client.ListenRaw() yang sudah ada
// tidak membuat koneksi baru, menggunakan method existing

type RouterCollector struct {
    client     *mikrotik.Client      // existing Client
    specs      []CommandSpec         // tier configs
    pools      *PoolManager          // pool per tier
    cache      *CacheManager         // redis cache
}
```

### 2. Pool Strategy
```
RouterClient (existing)
    │
    ├── Pool Tier 1 (High Freq): Max 5 connections
    │   ├── Traffic Monitor (interface/print follow=yes)
    │   └── Queue Stats (queue/simple/print stats)
    │
    ├── Pool Tier 2 (Event): Max 5 connections  
    │   ├── Hotspot Active (follow=yes)
    │   ├── PPP Active (follow=yes)
    │   └── System Logs (follow=yes)
    │
    └── Pool Tier 3 (On-Demand): Single connection
        └── Run() commands untuk static data
```

## 🏗️ Architecture Detail

### Tier 1 & 2: Listener Pools
```go
// specs.go - CommandSpec yang user berikan
type CommandSpec struct {
    Args     []string      // RouterOS command
    Channel  string        // Redis pub/sub channel (optional)
    Mode     Mode          // ModePublish | ModeCacheHSet | ModeCacheStream
    RedisKey string        // Key untuk cache
    KeyField string        // Field untuk HSET key
    CacheTTL time.Duration // TTL untuk cache
    Tier     Tier          // Tier1 | Tier2 | Tier3
}
```

### Pool Management
```go
// pool.go

const (
    MaxTier1Slots = 5  // High freq, limited slots
    MaxTier2Slots = 5  // Event-driven
    MaxTier3Slots = 1  // On-demand Run()
)

type PoolManager struct {
    tier1 chan *PoolSlot  // Slots untuk Tier 1
    tier2 chan *PoolSlot  // Slots untuk Tier 2
    tier3 *PoolSlot       // Single slot untuk Tier 3
}

type PoolSlot struct {
    id       string
    client   *mikrotik.Client  // reference ke client
    inUse    bool
    cancel   context.CancelFunc
}
```

### Cache Strategy per Tier

| Tier | Redis Strategy | Contoh Data |
|------|---------------|-------------|
| Tier 1 | HSET (last state) | interface:stats, queue:stats |
| Tier 1 | Pub/Sub (stream) | traffic:realtime |
| Tier 2 | HSET (last state) | hotspot:active, ppp:active |
| Tier 2 | Stream (event log) | logs:events |
| Tier 3 | HSET dengan TTL | ppp:secrets, profiles |
| Tier 4 | DEL on write | Invalidate Tier 3 keys |

## 🔄 Data Flow

### Tier 1 & 2: Listener Flow
```
RouterOS API
    │ follow=yes
    ▼
ListenRaw() → Channel
    │
    ├──► HSET redis (last state)
    └──► PUBLISH redis (stream event)
```

### Tier 3: Cache Flow
```
HTTP Request
    │
    ▼
Check Redis Cache
    │
    ├── Cache Hit ──► Return data
    │
    └── Cache Miss ──► Run() ke RouterOS
              │
              ▼
        HSET redis dengan TTL
              │
              ▼
        Return data
```

### Tier 4: Write + Invalidate
```
Write Request
    │
    ▼
Run() ke RouterOS (add/remove/set)
    │
    ▼
DEL redis key (Tier 3 cache)
    │
    ▼
Return success
(Tier 3 ticker akan refresh otomatis)
```

## 🛠️ Implementation Plan

### Phase 1: Core Types (1 file)
- `specs.go`: CommandSpec, Tier constants, Mode flags

### Phase 2: Pool Management (2 files)
- `pool.go`: PoolSlot, PoolManager dengan tier separation
- `worker.go`: Worker goroutine untuk handle listener

### Phase 3: Cache Layer (1 file)
- `cache.go`: Redis wrapper dengan HSET, PUBLISH, Stream support

### Phase 4: Collector (2 files)
- `collector.go`: Main Collector struct, Start(), Stop()
- `router.go`: Router instance management

### Phase 5: Testing (1 file)
- `cmd/collector/main.go`: Test dengan real/mock router

## 📋 Redis Key Convention

```
# Tier 1 - Real-time stats (HSET)
mikrotik:{router_id}:interface:stats    # HSET, field=name
mikrotik:{router_id}:queue:stats        # HSET, field=name

# Tier 2 - Active sessions (HSET)
mikrotik:{router_id}:hotspot:active     # HSET, field=mac-address
mikrotik:{router_id}:ppp:active         # HSET, field=name

# Tier 3 - Static data dengan TTL (HSET)
mikrotik:{router_id}:ppp:secrets        # HSET, field=name, TTL=5m
mikrotik:{router_id}:ppp:profiles       # HSET, field=name, TTL=5m
mikrotik:{router_id}:hotspot:users      # HSET, field=name, TTL=5m
mikrotik:{router_id}:ip:pools           # HSET, field=name, TTL=10m

# Tier 1 & 2 - Pub/Sub channels
mikrotik:{router_id}:pubsub:traffic
mikrotik:{router_id}:pubsub:queue
mikrotik:{router_id}:pubsub:logs

# Tier 1 & 2 - Event Stream (Redis Stream)
mikrotik:{router_id}:stream:events
```

## ⚠️ Constraints & Considerations

### 1. Connection Limit
- RouterOS API: ~25 commands per connection (safe limit)
- Tier 1: Max 5 slot (untuk 5 high-freq monitors)
- Tier 2: Max 5 slot (untuk event listeners)
- Tier 3: 1 slot (sequential Run() commands)

### 2. Resource Management
- Setiap slot punya context dengan cancel
- Graceful shutdown: cancel all contexts → close connections
- Auto-reconnect: handled by existing Client

### 3. Thread Safety
- PoolManager: thread-safe dengan channel
- CacheManager: Redis atomic operations
- Collector: 1 goroutine per slot

## 🧪 Testing Strategy

### Unit Test
- Mock Redis
- Mock MikroTik client
- Test pool allocation/deallocation
- Test cache hit/miss

### Integration Test
- Real Redis instance
- Real MikroTik router (RB750Gr3 di 192.168.27.1)
- Test dengan 3-5 concurrent specs

### Load Test
- Monitor memory usage
- Monitor goroutine count
- Monitor Redis connection count

## 🚀 Usage Example (Target)

```go
// Initialize collector
cfg := collector.Config{
    RedisAddr: "localhost:6379",
}

coll := collector.New(cfg)

// Register router dengan specs
routerCfg := mikrotik.Config{Host: "192.168.27.1", ...}
specs := collector.DefaultISPSpecs()

err := coll.RegisterRouter("router-1", routerCfg, specs)

// Start collecting
coll.Start()

// Later: get cached data
data, err := coll.GetCached("router-1", "ppp:secrets", "user1")

// Write operation dengan auto-invalidate
err := coll.Write("router-1", "/ppp/secret/add", params)
// Auto: DEL mikrotik:router-1:ppp:secrets

// Stop
coll.Stop()
```
