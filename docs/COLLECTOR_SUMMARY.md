# Ringkasan Implementasi: 4-Tier Collector

## 🎯 Gambaran Arsitektur

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              COLLECTOR                                       │
├─────────────────────────────────────────────────────────────────────────────┤
│  Router Instance (per router)                                                │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │  Pool Manager (Separated by Tier)                                    │    │
│  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐               │    │
│  │  │  Tier 1 Pool │  │  Tier 2 Pool │  │  Tier 3 Slot │               │    │
│  │  │  [5 slots]   │  │  [5 slots]   │  │   [1 slot]   │               │    │
│  │  │              │  │              │  │              │               │    │
│  │  │ • traffic    │  │ • hotspot    │  │ • Run() cmd  │               │    │
│  │  │ • queue      │  │ • ppp active │  │ • periodic   │               │    │
│  │  │   (1s)       │  │ • sysres     │  │   fetch      │               │    │
│  │  │              │  │   (5s)       │  │              │               │    │
│  │  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘               │    │
│  │         │                 │                 │                        │    │
│  │    Listener          Listener          Ticker                       │    │
│  │    Workers           Workers           Worker                       │    │
│  │         │                 │                 │                        │    │
│  │         └─────────────────┴─────────────────┘                        │    │
│  │                           │                                          │    │
│  │                    Cache Manager                                     │    │
│  │                    (Redis)                                           │    │
│  │                           │                                          │    │
│  │    ┌──────────────────────┼──────────────────────┐                   │    │
│  │    ▼                      ▼                      ▼                   │    │
│  │ ┌────────┐           ┌────────┐           ┌────────┐                │    │
│  │ │  HSET  │           │ PUB/SUB│           │ STREAM │                │    │
│  │ │(state) │           │(events)│           │ (logs) │                │    │
│  │ └────────┘           └────────┘           └────────┘                │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                           │                                                  │
│                           ▼                                                  │
│  Key Pattern: mikrotik:{router_id}:{resource}                                │
└─────────────────────────────────────────────────────────────────────────────┘
```

## 📦 File Structure

```
pkg/mikrotik/
├── mikrotik.go              # ✅ EXISTING - tidak diubah
├── client/
│   └── client.go            # ✅ EXISTING - digunakan via mikrotik.Client
└── collector/               # 🆕 NEW - 4 tier system
    ├── specs.go             # CommandSpec, Tier, Mode definitions
    ├── pool.go              # PoolManager dengan tier separation
    ├── cache.go             # Redis CacheManager
    ├── worker.go            # ListenerWorker & TickerWorker
    └── collector.go         # Main Collector orchestrator

cmd/collector/
└── main.go                  # 🆕 NEW - Test runner
```

## 🔑 Key Design Decisions

### 1. Tier Separation (Pool Isolation)
```go
// Tier 1 & 2 harus terpisah!
// Alasan: Traffic stats (1s) sangat agresif, bisa delay event listeners

const (
    MaxTier1Slots = 5  // Traffic, Queue - high freq
    MaxTier2Slots = 5  // Active sessions, Logs - events
    MaxTier3Slots = 1  // Run() commands - on demand
)
```

### 2. Cache Strategy per Tier

| Tier | Redis Type | Use Case | Example Key |
|------|-----------|----------|-------------|
| Tier 1 | HSET + PUB/SUB | Last state + Real-time stream | `interface:stats` |
| Tier 2 | HSET | Last state dashboard | `ppp:active` |
| Tier 3 | HSET dengan TTL | Static config cache | `ppp:secrets` (TTL=5m) |
| Tier 4 | DEL | Cache invalidation | Setelah write |

### 3. Integration dengan Existing Code

```go
// Menggunakan existing mikrotik.Client tanpa modifikasi!
client *mikrotik.Client

// Listen menggunakan existing method:
stopListen, err := client.ListenRaw(ctx, args, resultChan)

// Run menggunakan existing method:
results, err := client.RunRaw(ctx, args)
```

## 🚀 Usage Flow

### Step 1: Create Client (existing)
```go
cfg := client.Config{Host: "192.168.27.1", ...}
c, _ := client.New(cfg)
mtClient := mikrotik.NewClientFromConnection(c)
```

### Step 2: Create Collector
```go
collCfg := collector.Config{RedisAddr: "localhost:6379"}
coll, _ := collector.NewCollector(collCfg)
```

### Step 3: Register Router dengan Specs
```go
// Gunakan default specs atau custom
specs := collector.DefaultISPSpecs()  // Tier 1 & 2
tier3Specs := collector.Tier3Specs()  // Tier 3
allSpecs := append(specs, tier3Specs...)

// Register
coll.RegisterRouter("router-1", mtClient, allSpecs)
```

### Step 4: Start Collection
```go
coll.StartRouter("router-1")
```

### Step 5: Read Cache (Tier 3)
```go
data, _ := coll.GetCachedDataAll("router-1", "ppp:secrets")
```

### Step 6: Write (Tier 4) dengan Auto-Invalidate
```go
// Write dan invalidate cache
coll.Write("router-1", 
    []string{"/ppp/secret/add", "=name=test", "=password=123"},
    []string{"ppp:secrets"})  // Cache keys yang di-invalidate
```

## 🧪 Testing

```bash
# Run test dengan real router
go run cmd/collector/main.go \
    -host=192.168.27.1 \
    -user=admin \
    -pass=r00t \
    -redis=localhost:6379 \
    -duration=60
```

Output yang diharapkan:
```
[15:04:05] Pool Usage: T1[2/5] T2[3/5] T3[1/1]
  [Cache] interface:stats: 4 entries
  [Cache] ppp:active: 12 entries
  [Cache] ppp:secrets: 25 entries
```

## ⚠️ Critical Rules

### 1. Tier 1 & 2 Pool Terpisah
```go
// ❌ Jangan digabung!
// Tier 1 traffic (1s interval) bisa starve Tier 2 events

// ✅ Pisahkan:
tier1Pool := make([]*PoolSlot, MaxTier1Slots)  // 5 slots
tier2Pool := make([]*PoolSlot, MaxTier2Slots)  // 5 slots
```

### 2. Tier 3 Jangan Pakai Listener
```go
// ❌ SALAH - membuang slot connection permanen
{
    Args: []string{"/ppp/secret/print", "follow=yes"},
    Tier: Tier3,
}

// ✅ BENAR - pakai Run() + TTL
{
    Args: []string{"/ppp/secret/print"},
    Tier: Tier3,
    CacheTTL: 5 * time.Minute,
}
```

### 3. Tier 4 Invalidate Cache
```go
// Setelah write, hapus cache Tier 3
// Tier 3 ticker akan refresh otomatis dalam interval berikutnya

func (c *Collector) Write(routerID string, args []string, invalidateKeys []string) {
    // 1. Execute write
    client.RunRaw(ctx, args)
    
    // 2. Invalidate cache
    for _, key := range invalidateKeys {
        cache.Delete(ctx, key)
    }
}
```

## 📊 Redis Key Convention

```
mikrotik:{router_id}:interface:stats      # Tier 1 - HSET
mikrotik:{router_id}:queue:stats          # Tier 1 - HSET
mikrotik:{router_id}:pubsub:traffic       # Tier 1 - PUB/SUB

mikrotik:{router_id}:hotspot:active       # Tier 2 - HSET
mikrotik:{router_id}:ppp:active           # Tier 2 - HSET
mikrotik:{router_id}:system:resource      # Tier 2 - HSET

mikrotik:{router_id}:ppp:secrets          # Tier 3 - HSET (TTL=5m)
mikrotik:{router_id}:ppp:profiles         # Tier 3 - HSET (TTL=5m)
mikrotik:{router_id}:hotspot:users        # Tier 3 - HSET (TTL=5m)
mikrotik:{router_id}:ip:pools             # Tier 3 - HSET (TTL=10m)
```

## 🔧 Next Steps

1. **Implementasi** - File sudah lengkap, tinggal di-test
2. **Error Handling** - Add retry logic untuk connection failure
3. **Metrics** - Add prometheus metrics untuk monitoring
4. **Scaling** - Support multiple routers dengan goroutine per router
