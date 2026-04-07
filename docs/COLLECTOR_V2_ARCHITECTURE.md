# Arsitektur Collector V2: 3 Pipeline Design

## 🎯 Overview

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           COLLECTOR SUPERVISOR                               │
│                    (Monitor & Auto-restart Collectors)                       │
└─────────────────────────────────────────────────────────────────────────────┘
                                      │
        ┌─────────────────────────────┼─────────────────────────────┐
        │                             │                             │
        ▼                             ▼                             ▼
┌───────────────┐           ┌───────────────┐           ┌───────────────┐
│  PIPELINE A   │           │  PIPELINE B   │           │  PIPELINE C   │
│  Time-series  │           │  Operational  │           │   On-demand   │
│   → InfluxDB  │           │    → Redis    │           │    → Direct   │
└───────┬───────┘           └───────┬───────┘           └───────┬───────┘
        │                           │                           │
        ▼                           ▼                           ▼
┌───────────────┐           ┌───────────────┐           ┌───────────────┐
│  ConnPool A   │           │  ConnPool B   │           │  ConnPool C   │
│  (dedicated)  │           │  (dedicated)  │           │  (on-demand)  │
│  Aggressive   │           │  follow=yes   │           │  Run() only   │
│  short interval│          │  state changes│           │  no storage   │
└───────┬───────┘           └───────┬───────┘           └───────┬───────┘
        │                           │                           │
        ▼                           ▼                           ▼
┌───────────────┐           ┌───────────────┐           ┌───────────────┐
│   Collector   │           │   Collector   │           │   No collector│
│   (concurrent)│           │   Tier 2 & 3  │           │   Direct call │
└───────┬───────┘           └───────┬───────┘           └───────┬───────┘
        │                           │                           │
        ▼                           ▼                           ▼
┌───────────────┐           ┌───────────────┐                   │
│  Fan-in Chan  │           │  Fan-in Chan  │                   │
└───────┬───────┘           └───────┬───────┘                   │
        │                           │                           │
        ▼                           ▼                           ▼
┌───────────────┐           ┌───────────────┐           ┌───────────────┐
│  BatchWriter  │◄──────────┤  BatchWriter  │           │    Run()      │
│  (shared)     │           │  (shared)     │           │               │
└───────┬───────┘           └───────┬───────┘           └───────┬───────┘
        │                           │                           │
        ▼                           ▼                           ▼
┌───────────────┐           ┌───────────────┐           ┌───────────────┐
│   InfluxDB    │           │    Redis      │           │    Router     │
│  (historical) │           │  (last state) │           │  (live data)  │
│  line protocol│           │  HSET + TTL   │           │               │
└───────────────┘           └───────────────┘           └───────────────┘
        │                           │                           │
        └───────────────────────────┼───────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                         CUSTOM DASHBOARD                                     │
├─────────────────────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        │
│  │   Charts    │  │   Tables    │  │    Logs     │  │   Actions   │        │
│  │  (InfluxDB) │  │  (Redis)    │  │  (InfluxDB) │  │  (Pipeline C)│       │
│  │  Traffic    │  │  Active     │  │  Historical │  │  Write      │        │
│  │  Queue Stats│  │  Users      │  │  Events     │  │  Invalidate │        │
│  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘        │
└─────────────────────────────────────────────────────────────────────────────┘
```

## 📊 Pipeline Detail

### Pipeline A: Time-series → InfluxDB
**Purpose**: Data yang berubah terus-menerus, butuh histori

**Characteristics**:
- Connection pool terpisah (dedicated)
- Command concurrent, semua di-fan-in ke satu channel
- BatchWriter writes ke InfluxDB dalam line protocol
- Retention policy untuk data lifecycle

**Data Types**:
- Interface traffic (rx_byte, tx_byte, packets)
- Queue stats (rate, bytes, dropped)
- System resources (CPU, memory, uptime)

**InfluxDB Schema**:
```
measurement: interface_stats
tags: router, interface
fields: rx_byte, tx_byte, rx_packet, tx_packet, rx_drop, tx_drop

measurement: queue_stats
tags: router, queue_name
fields: bytes, packets, dropped, queued_bytes, rate

measurement: system_resource
tags: router
fields: cpu_load, memory_used, memory_total, uptime
```

### Pipeline B: Operational State → Redis
**Purpose**: State terakhir yang perlu diketahui, jarang query histori

**Tier 2 (follow=yes)**:
- PPPoE active connections
- Hotspot active sessions
- Interface status
- Data berubah sewaktu-waktu, perlu real-time update

**Tier 3 (Run() + ticker)**:
- PPPoE secrets (user accounts)
- PPP profiles
- Hotspot users
- IP pools
- Data jarang berubah, di-fetch periodic dengan TTL

**Redis Schema**:
```
# Tier 2 - Last state (HSET)
mikrotik:{router}:ppp:active:{name} -> hash
mikrotik:{router}:hotspot:active:{mac} -> hash
mikrotik:{router}:interface:status:{name} -> hash

# Tier 3 - Static dengan TTL (HSET)
mikrotik:{router}:ppp:secrets:{name} -> hash (TTL=5m)
mikrotik:{router}:ppp:profiles:{name} -> hash (TTL=5m)
mikrotik:{router}:ip:pools:{name} -> hash (TTL=10m)
```

### Pipeline C: On-demand Run()
**Purpose**: Direct query ke router, no storage

**Characteristics**:
- Tidak ada collector yang running
- Setiap request langsung Run() ke router
- Untuk operasi read: ping, log, diagnostics
- Untuk operasi write: add/edit/delete dengan cache invalidation

**Cache Invalidation**:
```go
// After write operation:
DEL mikrotik:{router}:ppp:secrets:{name}
// Pipeline B Tier 3 ticker akan refresh otomatis
```

## 🔧 Shared Components

### 1. Manager
- Manages lifecycle semua pipeline
- Handles router registration/unregistration
- Coordinated shutdown

### 2. ConnPool
- Shared pool management pattern
- Setiap pipeline punya pool instance tersendiri
- Auto-reconnect pada connection failure

### 3. BatchWriter
- Shared component untuk batch writes
- Pipeline A: Write ke InfluxDB (line protocol)
- Pipeline B: Write ke Redis (HSET pipeline)
- Configurable: batch size, flush interval

### 4. CollectorSupervisor
- Monitor Pipeline A & B collectors
- Restart dengan pool baru kalau koneksi putus
- Health checks periodic

## 🏗️ Implementation Structure

```
pkg/mikrotik/
├── mikrotik.go                    # Existing - tidak diubah
├── collector/
│   ├── spec.go                    # CommandSpec (shared)
│   ├── manager.go                 # Manager (shared)
│   ├── supervisor.go              # CollectorSupervisor (shared)
│   ├── pool/
│   │   ├── pool.go                # ConnPool (shared pattern)
│   │   └── options.go             # Pool configuration
│   ├── writer/
│   │   ├── batch_writer.go        # BatchWriter (shared)
│   │   ├── influx_writer.go       # InfluxDB specific
│   │   └── redis_writer.go        # Redis specific
│   ├── pipeline/
│   │   ├── pipeline.go            # Pipeline interface
│   │   ├── time_series/           # Pipeline A
│   │   │   ├── collector.go
│   │   │   └── metrics.go
│   │   ├── operational/           # Pipeline B
│   │   │   ├── tier2_collector.go # follow=yes
│   │   │   └── tier3_collector.go # ticker + Run()
│   │   └── ondemand/              # Pipeline C
│   │       └── runner.go
│   └── influxdb/
│       └── client.go              # InfluxDB client wrapper

cmd/collector/
└── main.go                        # Test runner
```

## 📈 Data Flow Examples

### Pipeline A: Traffic Monitoring
```
RouterOS API
    │ /interface/print follow=yes interval=1s
    ▼
[ConnPool A - Slot 1]
    │
    ▼
Fan-in Channel
    │
    ▼
BatchWriter (batch=100 or 50ms)
    │ line protocol
    ▼
InfluxDB
    │ SELECT mean(rx_byte) FROM interface_stats
    ▼
Dashboard Chart
```

### Pipeline B Tier 2: Active Sessions
```
RouterOS API
    │ /ppp/active/print follow=yes
    ▼
[ConnPool B - Slot 1]
    │
    ▼
Fan-in Channel
    │
    ▼
BatchWriter (batch=10 or 100ms)
    │ HSET pipeline
    ▼
Redis
    │ HGETALL mikrotik:router-1:ppp:active
    ▼
Dashboard Table
```

### Pipeline B Tier 3: Static Data
```
Ticker (every 5m)
    │
    ▼
[ConnPool B - Slot 5]
    │ /ppp/secret/print (Run, no follow)
    ▼
Parse Results
    │
    ▼
Redis HSET dengan TTL
    │ HGET mikrotik:router-1:ppp:secrets:user1
    ▼
Dashboard Form
```

### Pipeline C: Write Operation
```
Dashboard Action
    │ POST /api/ppp/secrets (add user)
    ▼
Pipeline C Runner
    │ /ppp/secret/add =name=newuser...
    ▼
[ConnPool C - On-demand]
    │
    ▼
RouterOS API (write)
    │
    ▼
Cache Invalidation
    │ DEL mikrotik:router-1:ppp:secrets:newuser
    ▼
Response Success
    (Tier 3 ticker akan refresh dalam 5m)
```

## 🎯 Design Principles

1. **Pipeline Isolation**: Setiap pipeline punya connection pool sendiri
   - Pipeline A (aggressive) tidak mengganggu Pipeline B (events)
   - Pipeline C (on-demand) tidak mengganggu yang lain

2. **Shared Efficiency**: BatchWriter dan Manager dipakai bersama
   - Reduces resource usage
   - Centralized configuration

3. **Fault Tolerance**: Supervisor restart otomatis
   - Connection failure → new pool → restart collector
   - No manual intervention

4. **Data Separation**: InfluxDB untuk histori, Redis untuk state
   - Query pattern yang berbeda di-handle oleh storage yang tepat
   - Dashboard query optimal untuk masing-masing use case
