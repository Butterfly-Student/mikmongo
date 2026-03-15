# Contoh Penggunaan

Dokumen ini berisi berbagai contoh penggunaan go-ros untuk berbagai skenario.

## Table of Contents

1. [Basic Examples](#basic-examples)
2. [Hotspot Management](#hotspot-management)
3. [Monitoring](#monitoring)
4. [Multi-Router](#multi-router)
5. [Advanced Examples](#advanced-examples)

---

## Basic Examples

### Koneksi Dasar

```go
package main

import (
    "context"
    "fmt"
    "time"
    
    "github.com/Butterfly-Student/go-ros/client"
)

func main() {
    ctx := context.Background()
    
    cfg := client.Config{
        Host:     "192.168.88.1",
        Port:     8728,
        Username: "admin",
        Password: "password",
        Timeout:  10 * time.Second,
    }
    
    c, err := client.New(cfg)
    if err != nil {
        panic(err)
    }
    defer c.Close()
    
    // Test connection
    reply, err := c.Run("/system/identity/print")
    if err != nil {
        panic(err)
    }
    
    for _, re := range reply.Re {
        fmt.Printf("Router Identity: %s\n", re.Map["name"])
    }
}
```

### Menggunakan Manager

```go
package main

import (
    "context"
    "fmt"
    "time"
    
    "github.com/Butterfly-Student/go-ros/client"
    "github.com/Butterfly-Student/go-ros/repository/hotspot"
)

func main() {
    ctx := context.Background()
    
    // Create manager
    manager := client.NewManager(nil)
    defer manager.CloseAll()
    
    // Router configurations
    routers := []struct {
        name string
        cfg  client.Config
    }{
        {
            name: "router-1",
            cfg: client.Config{
                Host:     "192.168.88.1",
                Port:     8728,
                Username: "admin",
                Password: "pass1",
                Timeout:  10 * time.Second,
            },
        },
        {
            name: "router-2",
            cfg: client.Config{
                Host:     "192.168.89.1",
                Port:     8728,
                Username: "admin",
                Password: "pass2",
                Timeout:  10 * time.Second,
            },
        },
    }
    
    // Connect to all routers
    for _, r := range routers {
        c, err := manager.GetOrConnect(ctx, r.name, r.cfg)
        if err != nil {
            fmt.Printf("Failed to connect to %s: %v\n", r.name, err)
            continue
        }
        
        // Use the client
        hotspotRepo := hotspot.NewRepository(c)
        users, _ := hotspotRepo.User().GetUsers(ctx, "")
        
        fmt.Printf("Router %s: %d hotspot users\n", r.name, len(users))
    }
    
    fmt.Printf("Connected routers: %v\n", manager.Names())
}
```

---

## Hotspot Management

### User Management Dashboard

```go
package main

import (
    "context"
    "fmt"
    "time"
    
    "github.com/Butterfly-Student/go-ros/client"
    "github.com/Butterfly-Student/go-ros/repository/hotspot"
)

type HotspotDashboard struct {
    client      *client.Client
    hotspotRepo hotspot.Repository
}

func NewDashboard(cfg client.Config) (*HotspotDashboard, error) {
    c, err := client.New(cfg)
    if err != nil {
        return nil, err
    }
    
    return &HotspotDashboard{
        client:      c,
        hotspotRepo: hotspot.NewRepository(c),
    }, nil
}

func (d *HotspotDashboard) GetStats(ctx context.Context) (*DashboardStats, error) {
    users, err := d.hotspotRepo.User().GetUsers(ctx, "")
    if err != nil {
        return nil, err
    }
    
    active, err := d.hotspotRepo.Active().GetActive(ctx)
    if err != nil {
        return nil, err
    }
    
    profiles, err := d.hotspotRepo.Profile().GetProfiles(ctx)
    if err != nil {
        return nil, err
    }
    
    // Count disabled users
    disabledCount := 0
    for _, u := range users {
        if u.Disabled {
            disabledCount++
        }
    }
    
    return &DashboardStats{
        TotalUsers:    len(users),
        ActiveUsers:   len(active),
        DisabledUsers: disabledCount,
        TotalProfiles: len(profiles),
    }, nil
}

func (d *HotspotDashboard) PrintStats(ctx context.Context) {
    stats, err := d.GetStats(ctx)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    fmt.Println("╔════════════════════════════════════╗")
    fmt.Println("║      HOTSPOT DASHBOARD             ║")
    fmt.Println("╠════════════════════════════════════╣")
    fmt.Printf("║ Total Users:    %4d              ║\n", stats.TotalUsers)
    fmt.Printf("║ Active Users:   %4d              ║\n", stats.ActiveUsers)
    fmt.Printf("║ Disabled Users: %4d              ║\n", stats.DisabledUsers)
    fmt.Printf("║ Total Profiles: %4d              ║\n", stats.TotalProfiles)
    fmt.Println("╚════════════════════════════════════╝")
}

type DashboardStats struct {
    TotalUsers    int
    ActiveUsers   int
    DisabledUsers int
    TotalProfiles int
}

func main() {
    cfg := client.Config{
        Host:     "192.168.88.1",
        Port:     8728,
        Username: "admin",
        Password: "password",
    }
    
    dashboard, err := NewDashboard(cfg)
    if err != nil {
        panic(err)
    }
    defer dashboard.client.Close()
    
    ctx := context.Background()
    dashboard.PrintStats(ctx)
}
```

### Batch User Operations

```go
package main

import (
    "context"
    "fmt"
    "time"
    
    "github.com/Butterfly-Student/go-ros/client"
    "github.com/Butterfly-Student/go-ros/domain"
    "github.com/Butterfly-Student/go-ros/repository/hotspot"
)

func createBatchUsers(ctx context.Context, repo hotspot.Repository, count int) error {
    baseName := "user"
    basePass := "pass"
    
    for i := 1; i <= count; i++ {
        user := &domain.HotspotUser{
            Name:        fmt.Sprintf("%s%d", baseName, i),
            Password:    fmt.Sprintf("%s%d", basePass, i),
            Profile:     "default",
            Comment:     "Batch created",
            LimitUptime: "1h",
        }
        
        _, err := repo.User().AddUser(ctx, user)
        if err != nil {
            return fmt.Errorf("failed to create user %d: %w", i, err)
        }
    }
    
    return nil
}

func removeBatchUsers(ctx context.Context, repo hotspot.Repository, prefix string) error {
    users, err := repo.User().GetUsers(ctx, "")
    if err != nil {
        return err
    }
    
    var idsToRemove []string
    for _, user := range users {
        if len(user.Name) >= len(prefix) && user.Name[:len(prefix)] == prefix {
            idsToRemove = append(idsToRemove, user.ID)
        }
    }
    
    if len(idsToRemove) > 0 {
        return repo.User().RemoveUsers(ctx, idsToRemove)
    }
    
    return nil
}

func main() {
    cfg := client.Config{
        Host:     "192.168.88.1",
        Port:     8728,
        Username: "admin",
        Password: "password",
    }
    
    c, err := client.New(cfg)
    if err != nil {
        panic(err)
    }
    defer c.Close()
    
    ctx := context.Background()
    repo := hotspot.NewRepository(c)
    
    // Create 10 users
    fmt.Println("Creating 10 users...")
    if err := createBatchUsers(ctx, repo, 10); err != nil {
        panic(err)
    }
    
    // Remove users with prefix "user"
    fmt.Println("Removing batch users...")
    if err := removeBatchUsers(ctx, repo, "user"); err != nil {
        panic(err)
    }
    
    fmt.Println("Done!")
}
```

---

## Monitoring

### Real-time System Monitor

```go
package main

import (
    "context"
    "fmt"
    "os"
    "os/signal"
    "syscall"
    "time"
    
    "github.com/Butterfly-Student/go-ros/client"
    "github.com/Butterfly-Student/go-ros/repository/monitor"
)

func main() {
    cfg := client.Config{
        Host:     "192.168.88.1",
        Port:     8728,
        Username: "admin",
        Password: "password",
    }
    
    c, err := client.New(cfg)
    if err != nil {
        panic(err)
    }
    defer c.Close()
    
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    
    // Handle Ctrl+C
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    
    go func() {
        <-sigChan
        fmt.Println("\nStopping monitor...")
        cancel()
    }()
    
    monRepo := monitor.NewRepository(c)
    
    // Start system resource monitoring
    resultChan := make(chan *domain.SystemResourceMonitorStats)
    stop, err := monRepo.System().StartSystemResourceMonitorListen(ctx, resultChan)
    if err != nil {
        panic(err)
    }
    defer stop()
    
    fmt.Println("System Monitor Started (Press Ctrl+C to stop)")
    fmt.Println("─────────────────────────────────────────")
    
    for stats := range resultChan {
        fmt.Printf("\rCPU: %5.1f%% | Memory: %5.1f%% | Uptime: %s",
            stats.CPULoad,
            stats.MemoryUsage,
            stats.Uptime)
    }
}
```

### Interface Traffic Monitor

```go
package main

import (
    "context"
    "fmt"
    "time"
    
    "github.com/Butterfly-Student/go-ros/client"
    "github.com/Butterfly-Student/go-ros/repository/monitor"
)

func monitorInterface(ctx context.Context, c *client.Client, ifaceName string) {
    monRepo := monitor.NewRepository(c)
    
    trafficChan := make(chan *domain.TrafficMonitorStats)
    stop, err := monRepo.Interface().StartTrafficMonitorListen(ctx, ifaceName, trafficChan)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    defer stop()
    
    fmt.Printf("Monitoring interface: %s\n", ifaceName)
    fmt.Println("Time          RX Rate      TX Rate      RX Total     TX Total")
    fmt.Println("─────────────────────────────────────────────────────────────")
    
    for stats := range trafficChan {
        fmt.Printf("%s  %10s  %10s  %10s  %10s\n",
            time.Now().Format("15:04:05"),
            formatBytes(stats.RXRate)+"/s",
            formatBytes(stats.TXRate)+"/s",
            formatBytes(stats.RXBytes),
            formatBytes(stats.TXBytes))
    }
}

func formatBytes(bytes int64) string {
    const unit = 1024
    if bytes < unit {
        return fmt.Sprintf("%d B", bytes)
    }
    div, exp := int64(unit), 0
    for n := bytes / unit; n >= unit; n /= unit {
        div *= unit
        exp++
    }
    return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func main() {
    cfg := client.Config{
        Host:     "192.168.88.1",
        Port:     8728,
        Username: "admin",
        Password: "password",
    }
    
    c, err := client.New(cfg)
    if err != nil {
        panic(err)
    }
    defer c.Close()
    
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
    defer cancel()
    
    monitorInterface(ctx, c, "ether1")
}
```

### Log Watcher

```go
package main

import (
    "context"
    "fmt"
    "strings"
    "time"
    
    "github.com/Butterfly-Student/go-ros/client"
    "github.com/Butterfly-Student/go-ros/repository/monitor"
)

func watchLogs(ctx context.Context, c *client.Client, filter string) {
    monRepo := monitor.NewRepository(c)
    
    logChan := make(chan *domain.LogEntry)
    stop, err := monRepo.Log().ListenLogs(ctx, logChan)
    if err != nil {
        panic(err)
    }
    defer stop()
    
    fmt.Printf("Watching logs (filter: %s)...\n", filter)
    fmt.Println("─────────────────────────────────────────────────")
    
    for entry := range logChan {
        // Apply filter
        if filter != "" && !strings.Contains(entry.Message, filter) {
            continue
        }
        
        fmt.Printf("[%s] %s: %s\n", 
            entry.Time, 
            entry.Topics, 
            entry.Message)
    }
}

func main() {
    cfg := client.Config{
        Host:     "192.168.88.1",
        Port:     8728,
        Username: "admin",
        Password: "password",
    }
    
    c, err := client.New(cfg)
    if err != nil {
        panic(err)
    }
    defer c.Close()
    
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
    defer cancel()
    
    // Watch hotspot logs only
    watchLogs(ctx, c, "hotspot")
}
```

---

## Multi-Router

### Sync Users Across Routers

```go
package main

import (
    "context"
    "fmt"
    "time"
    
    "github.com/Butterfly-Student/go-ros/client"
    "github.com/Butterfly-Student/go-ros/domain"
    "github.com/Butterfly-Student/go-ros/repository/hotspot"
)

type RouterConfig struct {
    Name   string
    Config client.Config
}

func syncUsers(ctx context.Context, source, target hotspot.Repository) error {
    // Get users from source
    users, err := source.User().GetUsers(ctx, "")
    if err != nil {
        return fmt.Errorf("failed to get source users: %w", err)
    }
    
    fmt.Printf("Found %d users in source\n", len(users))
    
    // Add users to target
    for _, user := range users {
        newUser := &domain.HotspotUser{
            Name:            user.Name,
            Password:        user.Password,
            Profile:         user.Profile,
            Comment:         user.Comment,
            LimitUptime:     user.LimitUptime,
            LimitBytesTotal: user.LimitBytesTotal,
        }
        
        _, err := target.User().AddUser(ctx, newUser)
        if err != nil {
            fmt.Printf("Failed to add user %s: %v\n", user.Name, err)
            continue
        }
        
        fmt.Printf("Synced user: %s\n", user.Name)
    }
    
    return nil
}

func main() {
    routers := []RouterConfig{
        {
            Name: "source",
            Config: client.Config{
                Host:     "192.168.88.1",
                Port:     8728,
                Username: "admin",
                Password: "pass1",
            },
        },
        {
            Name: "target",
            Config: client.Config{
                Host:     "192.168.89.1",
                Port:     8728,
                Username: "admin",
                Password: "pass2",
            },
        },
    }
    
    ctx := context.Background()
    
    // Connect to routers
    var clients []*client.Client
    var repos []hotspot.Repository
    
    for _, r := range routers {
        c, err := client.New(r.Config)
        if err != nil {
            fmt.Printf("Failed to connect to %s: %v\n", r.Name, err)
            continue
        }
        defer c.Close()
        
        clients = append(clients, c)
        repos = append(repos, hotspot.NewRepository(c))
        
        fmt.Printf("Connected to %s\n", r.Name)
    }
    
    if len(repos) >= 2 {
        // Sync from first to second
        fmt.Println("\nSyncing users...")
        if err := syncUsers(ctx, repos[0], repos[1]); err != nil {
            fmt.Printf("Sync failed: %v\n", err)
        }
    }
}
```

### Router Health Check

```go
package main

import (
    "context"
    "fmt"
    "time"
    
    "github.com/Butterfly-Student/go-ros/client"
    "github.com/Butterfly-Student/go-ros/repository/system"
)

type HealthStatus struct {
    Name      string
    Online    bool
    Uptime    string
    CPU       int
    Memory    float64
    Version   string
    Error     string
}

func checkRouterHealth(ctx context.Context, name string, cfg client.Config) HealthStatus {
    status := HealthStatus{Name: name}
    
    c, err := client.New(cfg)
    if err != nil {
        status.Error = err.Error()
        return status
    }
    defer c.Close()
    
    status.Online = true
    
    sysRepo := system.NewRepository(c)
    
    resources, err := sysRepo.Resources().GetResources(ctx)
    if err != nil {
        status.Error = err.Error()
        return status
    }
    
    status.Uptime = resources.Uptime
    status.CPU = resources.CpuLoad
    status.Memory = float64(resources.TotalMemory-resources.FreeMemory) / float64(resources.TotalMemory) * 100
    status.Version = resources.Version
    
    return status
}

func printHealthStatus(statuses []HealthStatus) {
    fmt.Println("╔═══════════════════════════════════════════════════════════╗")
    fmt.Println("║                 ROUTER HEALTH CHECK                       ║")
    fmt.Println("╠═══════════════════════════════════════════════════════════╣")
    
    for _, s := range statuses {
        status := "✓ ONLINE"
        if !s.Online {
            status = "✗ OFFLINE"
        }
        
        fmt.Printf("║ %-10s %s\n", s.Name, status)
        
        if s.Online {
            fmt.Printf("║   Version: %-20s Uptime: %-15s\n", s.Version, s.Uptime)
            fmt.Printf("║   CPU: %3d%%                    Memory: %5.1f%%\n", s.CPU, s.Memory)
        } else {
            fmt.Printf("║   Error: %s\n", s.Error)
        }
        fmt.Println("║")
    }
    
    fmt.Println("╚═══════════════════════════════════════════════════════════╝")
}

func main() {
    routers := []struct {
        name string
        cfg  client.Config
    }{
        {
            name: "Router-1",
            cfg: client.Config{
                Host:     "192.168.88.1",
                Port:     8728,
                Username: "admin",
                Password: "pass1",
            },
        },
        {
            name: "Router-2",
            cfg: client.Config{
                Host:     "192.168.89.1",
                Port:     8728,
                Username: "admin",
                Password: "pass2",
            },
        },
    }
    
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    var statuses []HealthStatus
    
    for _, r := range routers {
        fmt.Printf("Checking %s...\n", r.name)
        status := checkRouterHealth(ctx, r.name, r.cfg)
        statuses = append(statuses, status)
    }
    
    fmt.Println()
    printHealthStatus(statuses)
}
```

---

## Advanced Examples

### Complete Voucher System

```go
package main

import (
    "context"
    "fmt"
    "os"
    "time"
    
    "github.com/Butterfly-Student/go-ros/client"
    "github.com/Butterfly-Student/go-ros/repository/hotspot"
    "github.com/Butterfly-Student/go-ros/repository/system"
    mikhmonRepo "github.com/Butterfly-Student/go-ros/repository/mikhmon"
    mikhmonDomain "github.com/Butterfly-Student/go-ros/domain/mikhmon"
)

type VoucherSystem struct {
    ctx         context.Context
    client      *client.Client
    hotspotRepo hotspot.Repository
    sysRepo     system.Repository
    voucherRepo mikhmonRepo.VoucherRepository
    profileRepo mikhmonRepo.ProfileRepository
    expireRepo  mikhmonRepo.ExpireRepository
}

func NewVoucherSystem(cfg client.Config) (*VoucherSystem, error) {
    c, err := client.New(cfg)
    if err != nil {
        return nil, err
    }
    
    hotspotRepo := hotspot.NewRepository(c)
    sysRepo := system.NewRepository(c)
    generatorRepo := mikhmonRepo.NewGeneratorRepository()
    
    return &VoucherSystem{
        ctx:         context.Background(),
        client:      c,
        hotspotRepo: hotspotRepo,
        sysRepo:     sysRepo,
        voucherRepo: mikhmonRepo.NewVoucherRepository(c, hotspotRepo, generatorRepo),
        profileRepo: mikhmonRepo.NewProfileRepository(hotspotRepo),
        expireRepo:  mikhmonRepo.NewExpireRepository(c, sysRepo),
    }, nil
}

func (s *VoucherSystem) Close() {
    s.client.Close()
}

func (s *VoucherSystem) GenerateVouchers(quantity int, profile string) (*mikhmonDomain.VoucherBatch, error) {
    req := &mikhmonDomain.VoucherGenerateRequest{
        Quantity:   quantity,
        Profile:    profile,
        Mode:       mikhmonDomain.VoucherModeVoucher,
        NameLength: 6,
        CharSet:    mikhmonDomain.CharSetUpplow1,
        TimeLimit:  "1h",
    }
    
    return s.voucherRepo.GenerateBatch(s.ctx, req)
}

func (s *VoucherSystem) PrintVouchers(batch *mikhmonDomain.VoucherBatch) {
    fmt.Println("╔════════════════════════════════════════╗")
    fmt.Println("║         VOUCHER HOTSPOT                ║")
    fmt.Println("╠════════════════════════════════════════╣")
    fmt.Printf("║ Kode: %-32s ║\n", batch.Code)
    fmt.Printf("║ Profile: %-29s ║\n", batch.Profile)
    fmt.Printf("║ Jumlah: %-30d ║\n", batch.Quantity)
    fmt.Println("╠════════════════════════════════════════╣")
    
    for i, v := range batch.Vouchers {
        fmt.Printf("║ %2d. %-10s / %-10s      ║\n", i+1, v.Name, v.Password)
    }
    
    fmt.Println("╚════════════════════════════════════════╝")
}

func (s *VoucherSystem) SetupExpireMonitor() error {
    return s.expireRepo.SetupExpireMonitor(s.ctx)
}

func (s *VoucherSystem) CreateProfile(name string, rate string, price int64) error {
    req := &mikhmonDomain.ProfileRequest{
        Name:        name,
        AddressPool: "hs-pool",
        RateLimit:   rate,
        SharedUsers: 1,
        Config: mikhmonDomain.ProfileConfig{
            Name:         name,
            Price:        price,
            SellingPrice: price + 2000,
            Validity:     "1h",
            ExpireMode:   mikhmonDomain.ExpireModeRemove,
        },
    }
    
    return s.profileRepo.CreateProfile(s.ctx, req)
}

func main() {
    cfg := client.Config{
        Host:     "192.168.88.1",
        Port:     8728,
        Username: "admin",
        Password: os.Getenv("MIKROTIK_PASSWORD"),
    }
    
    system, err := NewVoucherSystem(cfg)
    if err != nil {
        fmt.Printf("Failed to connect: %v\n", err)
        os.Exit(1)
    }
    defer system.Close()
    
    fmt.Println("✓ Connected to MikroTik")
    
    // Setup expire monitor
    fmt.Println("Setting up expire monitor...")
    if err := system.SetupExpireMonitor(); err != nil {
        fmt.Printf("Warning: %v\n", err)
    }
    
    // Generate vouchers
    fmt.Println("\nGenerating vouchers...")
    batch, err := system.GenerateVouchers(5, "default")
    if err != nil {
        fmt.Printf("Failed: %v\n", err)
        os.Exit(1)
    }
    
    system.PrintVouchers(batch)
}
```

### Backup Configuration

```go
package main

import (
    "context"
    "fmt"
    "time"
    
    "github.com/Butterfly-Student/go-ros/client"
    "github.com/Butterfly-Student/go-ros/domain"
    "github.com/Butterfly-Student/go-ros/repository/system"
)

func createBackup(ctx context.Context, sysRepo system.Repository) error {
    // Create backup script
    backupName := fmt.Sprintf("backup-%s", time.Now().Format("20060102-150405"))
    
    script := &domain.SystemScript{
        Name:    "auto-backup",
        Owner:   "admin",
        Policy:  "read,write,policy,test",
        Source:  fmt.Sprintf("/system backup save name=%s", backupName),
        Comment: "Auto backup",
    }
    
    _, err := sysRepo.Scripts().AddScript(ctx, script)
    if err != nil {
        return err
    }
    
    // Run the script
    err = sysRepo.Scripts().RunScript(ctx, "auto-backup")
    if err != nil {
        return err
    }
    
    fmt.Printf("Backup created: %s\n", backupName)
    return nil
}

func main() {
    cfg := client.Config{
        Host:     "192.168.88.1",
        Port:     8728,
        Username: "admin",
        Password: "password",
    }
    
    c, err := client.New(cfg)
    if err != nil {
        panic(err)
    }
    defer c.Close()
    
    ctx := context.Background()
    sysRepo := system.NewRepository(c)
    
    if err := createBackup(ctx, sysRepo); err != nil {
        panic(err)
    }
}
```

---

**Selamat mencoba contoh-contoh di atas!** 🚀
