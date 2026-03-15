# Arsitektur Sistem go-ros

Dokumen ini menjelaskan arsitektur dan design pattern yang digunakan dalam library go-ros.

## Overview Arsitektur

```
┌─────────────────────────────────────────────────────────────┐
│                    Aplikasi Anda                            │
└──────────────────────┬──────────────────────────────────────┘
                       │
                       │ Menggunakan
                       ▼
┌─────────────────────────────────────────────────────────────┐
│                   Repository Layer                            │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐        │
│  │ Hotspot  │ │   PPP    │ │  System  │ │ Mikhmon  │        │
│  └────┬─────┘ └────┬─────┘ └────┬─────┘ └────┬─────┘        │
└───────┼────────────┼────────────┼────────────┼──────────────┘
        │            │            │            │
        └────────────┴────────────┴────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────────────┐
│                    Client Layer                               │
│  ┌──────────────────────────────────────────────────────┐   │
│  │                    Manager                              │   │
│  │  (Multi-router connection management)                   │   │
│  └──────────────────────┬───────────────────────────────┘   │
│                         │                                    │
│              ┌──────────┴──────────┐                        │
│              │                     │                        │
│              ▼                     ▼                        │
│  ┌──────────────────┐  ┌──────────────────┐                │
│  │  Client (Router1)│  │  Client (Router2)│                │
│  │  - Async mode    │  │  - Async mode    │                │
│  │  - Auto-reconnect│  │  - Auto-reconnect│                │
│  └────────┬─────────┘  └────────┬─────────┘                │
└───────────┼─────────────────────┼──────────────────────────┘
            │                     │
            ▼                     ▼
┌─────────────────────────────────────────────────────────────┐
│                 Low-Level Protocol                            │
│                    pkg/routeros                               │
│         (RouterOS API protocol implementation)                │
└──────────────────────┬──────────────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────────────┐
│              TCP/TLS Connection                               │
│              Port 8728 / 8729                                 │
└─────────────────────────────────────────────────────────────┘
```

## Layer Architecture

### 1. Repository Layer

Repository layer mengimplementasikan **Repository Pattern** yang memisahkan logika bisnis dari akses data.

#### Struktur Repository

```
repository/
├── hotspot/
│   ├── interface.go      # Definisi interface
│   ├── user.go          # Implementasi user repository
│   ├── profile.go       # Implementasi profile repository
│   ├── active.go        # Implementasi active repository
│   ├── repository.go    # Aggregator
│   └── ...
├── ppp/
│   ├── interface.go
│   ├── secret.go
│   ├── profile.go
│   └── repository.go
└── ...
```

#### Repository Pattern Benefits

1. **Abstraction** - Aplikasi tidak perlu tahu detail RouterOS API
2. **Testability** - Mudah di-mock untuk unit testing
3. **Maintainability** - Perubahan di satu tempat, tidak perlu ubah aplikasi
4. **Type Safety** - Return typed structs, bukan raw map

#### Contoh Repository Interface

```go
// domain/hotspot.go
type HotspotUser struct {
    ID       string
    Name     string
    Password string
    Profile  string
    // ... fields lainnya
}

// repository/hotspot/interface.go
type UserRepository interface {
    GetUsers(ctx context.Context, profile string) ([]*domain.HotspotUser, error)
    GetUserByName(ctx context.Context, name string) (*domain.HotspotUser, error)
    AddUser(ctx context.Context, user *domain.HotspotUser) (string, error)
    UpdateUser(ctx context.Context, id string, user *domain.HotspotUser) error
    RemoveUser(ctx context.Context, id string) error
}

// Penggunaan
hotspotRepo := hotspot.NewRepository(client)
users, err := hotspotRepo.User().GetUsers(ctx, "default")
```

### 2. Client Layer

Client layer menangani koneksi ke MikroTik RouterOS.

#### Client

```go
type Client struct {
    conn        *routeros.Client  // Low-level connection
    config      Config            // Configuration
    asyncCtx    context.Context   // Async context
    asyncCancel context.CancelFunc
    mu          sync.RWMutex      // Thread safety
    closed      bool
    logger      *zap.Logger
}
```

**Fitur Client:**

1. **Async Mode** - Single connection untuk banyak concurrent command
2. **Auto-Reconnect** - Reconnect otomatis dengan exponential backoff
3. **Context Support** - Timeout dan cancellation
4. **Thread-Safe** - Safe untuk concurrent use

#### Manager

```go
type Manager struct {
    clients map[string]*Client  // Map nama router ke client
    mu      sync.RWMutex
    logger  *zap.Logger
}
```

**Fitur Manager:**

1. **Named Connections** - Register router dengan nama
2. **Lazy Connection** - Connect hanya saat diperlukan
3. **Connection Caching** - Reuse existing connections
4. **Health Check** - Deteksi disconnect dan reconnect

### 3. Low-Level Protocol Layer

Package `pkg/routeros` mengimplementasikan protokol RouterOS API:

- **Sentence** - Representasi satu baris dari RouterOS
- **Reply** - Response dari command (multiple sentences)
- **ListenReply** - Streaming response untuk follow commands
- **Tag Multiplexing** - Async command execution

## Data Flow

### 1. Read Operation Flow

```
Aplikasi
    │
    │ hotspotRepo.User().GetUsers(ctx, "")
    ▼
Repository
    │
    │ client.RunContext(ctx, "/ip/hotspot/user/print")
    ▼
Client
    │
    │ conn.RunContext(ctx, sentence...)
    ▼
routeros.Client
    │
    │ Kirim command ke MikroTik
    ▼
MikroTik RouterOS
    │
    │ Response: !re sentences
    ▼
routeros.Client
    │
    │ Parse ke []map[string]string
    ▼
Client
    │
    │ Return *routeros.Reply
    ▼
Repository
    │
    │ Parse map ke []*domain.HotspotUser
    ▼
Aplikasi
    │
    │ Terima []*domain.HotspotUser
```

### 2. Write Operation Flow

```
Aplikasi
    │
    │ hotspotRepo.User().AddUser(ctx, user)
    ▼
Repository
    │
    │ Build command: /ip/hotspot/user/add
    │ dengan parameters dari struct user
    ▼
Client
    │
    │ client.RunContext(ctx, args...)
    ▼
routeros.Client
    │
    │ Kirim command dengan tag
    ▼
MikroTik RouterOS
    │
    │ Response: !done dengan ret attribute
    ▼
routeros.Client
    │
    │ Return ID yang baru dibuat
    ▼
Repository
    │
    │ Return ID ke aplikasi
    ▼
Aplikasi
```

### 3. Streaming/Monitoring Flow

```
Aplikasi
    │
    │ monitorRepo.System().StartSystemResourceMonitorListen(ctx, ch)
    ▼
Repository
    │
    │ client.ListenArgs(args)
    ▼
Client
    │
    │ conn.ListenArgsContext(ctx, args)
    ▼
routeros.Client
    │
    │ Kirim follow command
    │ Buka streaming channel
    ▼
MikroTik RouterOS
    │
    │ Kirim !re setiap interval
    │ (misal: setiap 1 detik)
    ▼
routeros.Client
    │
    │ Terima sentences
    │ Kirim ke channel
    ▼
Client
    │
    │ ListenBatches() - debounce rapid updates
    ▼
Repository
    │
    │ Parse ke struct
    │ Kirim ke resultChan
    ▼
Aplikasi
    │
    │ Terima data real-time dari channel
```

## Domain Models

### Struktur Domain

```
domain/
├── hotspot.go      # HotspotUser, HotspotActive, UserProfile, dll
├── ppp.go          # PPPSecret, PPPProfile, PPPActive
├── ip.go           # IPAddress, IPPool
├── firewall.go     # NATRule, FirewallRule, AddressList
├── queue.go        # QueueStats, SimpleQueue, TreeQueue
├── system.go       # SystemResource, SystemIdentity, Scheduler, dll
├── interface.go    # Interface, TrafficStats
├── ping.go         # PingConfig, PingResult
├── voucher.go      # VoucherGenerateRequest, Voucher
├── report.go       # SalesReport, ReportSummary
└── mikhmon/        # Models khusus Mikhmon
    ├── generator.go
    ├── profile.go
    ├── report.go
    └── voucher.go
```

### Design Principles

1. **JSON Tags** - Semua field memiliki json tag untuk serialization
2. **Validation Tags** - Request structs menggunakan validation tags
3. **Pointer Slices** - `[]*DomainType` untuk konsistensi dan efisiensi
4. **Optional Fields** - Gunakan `omitempty` untuk field opsional

## Connection Management

### Lifecycle Koneksi

```
┌─────────────┐
│   Created   │ ← NewClient() atau GetOrConnect()
└──────┬──────┘
       │
       │ Connect()
       ▼
┌─────────────┐
│  Connected  │ ← Async mode enabled
└──────┬──────┘
       │
       │ Connection lost
       ▼
┌─────────────┐
│Disconnected │
└──────┬──────┘
       │
       │ Auto-reconnect
       ▼
┌─────────────┐
│ Reconnecting│ ← Exponential backoff
└──────┬──────┘
       │
       │ Success
       ▼
┌─────────────┐
│  Connected  │ ← Back to connected state
└─────────────┘
       │
       │ Close() atau Unregister()
       ▼
┌─────────────┐
│   Closed    │
└─────────────┘
```

### Auto-Reconnect Mechanism

```go
func (c *Client) reconnect() {
    backoff := reconnectBaseDelay  // 1 detik
    
    for {
        // Coba reconnect
        conn, err := c.dial(ctx)
        
        if err == nil {
            // Success - enable async dan return
            c.conn = conn
            go c.watchAsync(errCh)
            return
        }
        
        // Failed - tunggu dengan backoff
        time.Sleep(backoff)
        
        // Exponential backoff (max 30 detik)
        if backoff < reconnectMaxDelay {
            backoff *= 2
        }
    }
}
```

## Thread Safety

Semua komponen go-ros **thread-safe**:

1. **Client** - Menggunakan `sync.RWMutex` untuk protect connection state
2. **Manager** - Menggunakan `sync.RWMutex` untuk protect clients map
3. **Repositories** - Stateless, hanya menggunakan client methods yang sudah thread-safe

### Concurrent Access Pattern

```go
// Boleh dijalankan concurrently
go func() {
    users, _ := hotspotRepo.User().GetUsers(ctx, "")
}()

go func() {
    profiles, _ := hotspotRepo.Profile().GetProfiles(ctx)
}()

go func() {
    active, _ := hotspotRepo.Active().GetActive(ctx)
}()
```

## Error Handling

### Error Types

1. **Connection Errors** - Network issues, authentication failure
2. **API Errors** - RouterOS returned error (!trap sentence)
3. **Parse Errors** - Failed to parse RouterOS response
4. **Timeout Errors** - Context deadline exceeded

### Error Propagation

```go
// Repository layer menambahkan context
return fmt.Errorf("failed to get users: %w", err)

// Client layer menambahkan router info  
return fmt.Errorf("connect mikrotik %s: %w", cfg.Host, err)

// Aplikasi bisa menggunakan errors.Is()
if errors.Is(err, context.DeadlineExceeded) {
    // Handle timeout
}
```

## Best Practices

### 1. Selalu Gunakan Context

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

users, err := repo.User().GetUsers(ctx, "")
```

### 2. Close Resources

```go
defer client.Close()
// atau
defer manager.CloseAll()
```

### 3. Handle Errors

```go
users, err := repo.User().GetUsers(ctx, "")
if err != nil {
    log.Printf("Failed to get users: %v", err)
    return err
}
```

### 4. Gunakan Manager untuk Multi-Router

```go
manager := client.NewManager(logger)
defer manager.CloseAll()

// Register multiple routers
manager.GetOrConnect(ctx, "router-1", cfg1)
manager.GetOrConnect(ctx, "router-2", cfg2)
```

## Performance Considerations

### 1. Async Mode

- Single TCP connection untuk semua command
- Tag multiplexing untuk concurrent execution
- Tidak perlu buat connection per command

### 2. Proplist Optimization

```go
// Hanya request field yang diperlukan
const ProplistHotspotUserDefault = ".id,name,profile,disabled"

users, _ := repo.User().GetUsers(ctx, "", ProplistHotspotUserDefault)
```

### 3. Batch Operations

```go
// Lebih efisien daripada loop single operations
repo.User().RemoveUsers(ctx, []string{"id1", "id2", "id3"})
```

### 4. Connection Pooling

- Gunakan Manager untuk maintain persistent connections
- Avoid create/close connection repeatedly

## Security

### 1. TLS/SSL

```go
cfg := client.Config{
    Host:     "192.168.88.1",
    Port:     8729,  // API-SSL port
    Username: "admin",
    Password: "password",
    UseTLS:   true,
}
```

### 2. Credential Management

❌ **Jangan:**
```go
cfg := client.Config{
    Password: "hardcoded-password",  // ❌ Bahaya!
}
```

✅ **Lakukan:**
```go
password := os.Getenv("MIKROTIK_PASSWORD")
cfg := client.Config{
    Password: password,
}
```

### 3. Firewall Rules

```bash
# Batasi akses API
/ip firewall filter add chain=input protocol=tcp dst-port=8728 src-address=192.168.88.0/24 action=accept
/ip firewall filter add chain=input protocol=tcp dst-port=8728 action=drop
```

## Summary

go-ros menggunakan layered architecture dengan:

1. **Repository Pattern** - Clean data access
2. **Async Client** - High-performance connections
3. **Manager** - Multi-router support
4. **Type-Safe Domain** - Compile-time safety
5. **Context Support** - Proper cancellation and timeouts

Arsitektur ini memudahkan:
- Testing dan maintenance
- Scaling ke multiple routers
- Implementasi fitur baru
- Error handling yang robust
