# go-ros - RouterOS API Client untuk Go

Library Go untuk mengakses MikroTik RouterOS API dengan mudah, cepat, dan reliable.

## Fitur Utama

- 🔌 **Koneksi Async** - Dukungan mode asynchronous untuk performa tinggi
- 🔄 **Auto-Reconnect** - Reconnect otomatis dengan exponential backoff
- 🏗️ **Repository Pattern** - Clean architecture dengan pattern repository
- 📦 **Multi-Router** - Manage multiple MikroTik router dalam satu aplikasi
- 🎯 **Type-Safe** - Domain models dengan tipe data yang jelas
- 🚀 **Context Support** - Full support untuk context cancellation dan timeout
- 📊 **Mikhmon Support** - Built-in support untuk sistem voucher hotspot Mikhmon
- 🔒 **TLS Support** - Koneksi aman dengan TLS/SSL

## Instalasi

```bash
go get github.com/Butterfly-Student/go-ros
```

## Quick Start

### Koneksi Dasar

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
    
    // Konfigurasi koneksi
    cfg := client.Config{
        Host:     "192.168.88.1",
        Port:     8728,
        Username: "admin",
        Password: "password",
        Timeout:  10 * time.Second,
    }
    
    // Buat client dan connect
    c, err := client.New(cfg)
    if err != nil {
        panic(err)
    }
    defer c.Close()
    
    // Buat repository
    hotspotRepo := hotspot.NewRepository(c)
    
    // Ambil daftar user hotspot
    users, err := hotspotRepo.User().GetUsers(ctx, "")
    if err != nil {
        panic(err)
    }
    
    for _, user := range users {
        fmt.Printf("User: %s, Profile: %s\n", user.Name, user.Profile)
    }
}
```

### Menggunakan Manager (Multi-Router)

```go
// Create manager
manager := client.NewManager(nil)
defer manager.CloseAll()

// Connect ke router dengan nama
cfg := client.Config{
    Host:     "192.168.88.1",
    Port:     8728,
    Username: "admin",
    Password: "password",
}

c, err := manager.GetOrConnect(ctx, "router-1", cfg)
if err != nil {
    panic(err)
}

// Gunakan client
hotspotRepo := hotspot.NewRepository(c)
users, _ := hotspotRepo.User().GetUsers(ctx, "")
```

## Struktur Project

```
go-ros/
├── client/          # Manajemen koneksi (Client & Manager)
├── domain/          # Domain models dan struct
│   └── mikhmon/     # Models khusus Mikhmon
├── repository/      # Repository pattern implementation
│   ├── hotspot/     # Hotspot repositories
│   ├── ppp/         # PPP repositories
│   ├── system/      # System repositories
│   ├── firewall/    # Firewall repositories
│   ├── queue/       # Queue repositories
│   ├── monitor/     # Monitoring repositories
│   └── mikhmon/     # Mikhmon repositories
├── utils/           # Utility functions
└── examples/        # Contoh penggunaan
```

## Modul yang Tersedia

### 1. Hotspot
- User management (CRUD, enable/disable, reset counters)
- Profile management
- Active sessions monitoring
- IP bindings
- Host management

### 2. PPP
- Secret/PPP user management
- PPP profile management
- Active PPP sessions

### 3. System
- System resources (CPU, memory, disk)
- System identity
- Scripts management
- Scheduler management
- Routerboard info

### 4. Firewall
- NAT rules
- Filter rules
- Address lists

### 5. IP Address
- IP address management
- IP pool management

### 6. Queue
- Simple queues
- Queue statistics monitoring

### 7. Monitor
- Real-time system monitoring
- Interface traffic monitoring
- Ping monitoring
- Log streaming

### 8. Mikhmon
- Voucher generation
- Profile dengan on-login script
- Sales reports
- Expire monitoring

## Dokumentasi Lengkap

- [Arsitektur Sistem](docs/ARCHITECTURE.md) - Memahami arsitektur go-ros
- [Referensi API](docs/API.md) - Daftar lengkap API dan models
- [Repository Guide](docs/REPOSITORIES.md) - Panduan penggunaan repository
- [Panduan Mikhmon](docs/MIKHMON.md) - Dokumentasi lengkap Mikhmon
- [Contoh Penggunaan](docs/EXAMPLES.md) - Berbagai contoh kode
- [Konfigurasi](docs/CONFIGURATION.md) - Opsi konfigurasi lengkap
- [Troubleshooting](docs/TROUBLESHOOTING.md) - Solusi masalah umum

## Contoh Penggunaan

### Generate Voucher Mikhmon

```go
import (
    "github.com/Butterfly-Student/go-ros/repository/hotspot"
    mikhmonRepo "github.com/Butterfly-Student/go-ros/repository/mikhmon"
    mikhmonDomain "github.com/Butterfly-Student/go-ros/domain/mikhmon"
)

// Setup repositories
hotspotRepo := hotspot.NewRepository(c)
generatorRepo := mikhmonRepo.NewGeneratorRepository()
voucherRepo := mikhmonRepo.NewVoucherRepository(c, hotspotRepo, generatorRepo)

// Generate voucher
req := &mikhmonDomain.VoucherGenerateRequest{
    Quantity:   10,
    Profile:    "default",
    Mode:       mikhmonDomain.VoucherModeVoucher, // vc = username = password
    NameLength: 6,
    CharSet:    mikhmonDomain.CharSetUpplow1,
    TimeLimit:  "1h",
    DataLimit:  "1G",
}

batch, err := voucherRepo.GenerateBatch(ctx, req)
if err != nil {
    panic(err)
}

for _, v := range batch.Vouchers {
    fmt.Printf("Voucher: %s / %s\n", v.Name, v.Password)
}
```

### Monitor Active Users

```go
// Get active hotspot users
active, err := hotspotRepo.Active().GetActive(ctx)
if err != nil {
    panic(err)
}

for _, a := range active {
    fmt.Printf("User: %s, IP: %s, MAC: %s, Uptime: %s\n",
        a.User, a.Address, a.MACAddress, a.Uptime)
}
```

### Real-time Monitoring

```go
// Monitor system resources
monitorRepo := monitor.NewRepository(c)

resultChan := make(chan *domain.SystemResourceMonitorStats)
stop, err := monitorRepo.System().StartSystemResourceMonitorListen(ctx, resultChan)
if err != nil {
    panic(err)
}
defer stop()

for stats := range resultChan {
    fmt.Printf("CPU: %.2f%%, Memory: %.2f%%\n", 
        stats.CPULoad, stats.MemoryUsage)
}
```

## Persyaratan

- Go 1.21 atau lebih baru
- MikroTik RouterOS dengan API service enabled
- Akses network ke MikroTik (port 8728 untuk API, 8729 untuk API-SSL)

## Enable API Service di MikroTik

```bash
/ip service enable api
/ip service set api port=8728
```

Untuk API dengan SSL:
```bash
/ip service enable api-ssl
/ip service set api-ssl port=8729
```

## Keamanan

⚠️ **Peringatan Keamanan:**
- Selalu gunakan koneksi TLS/SSL untuk production
- Jangan hardcode password di kode
- Gunakan environment variables atau secret management
- Batasi akses API dengan firewall rules

## Lisensi

MIT License - Lihat [LICENSE](LICENSE) untuk detail.

## Kontribusi

Kontribusi sangat diterima! Silakan buat issue atau pull request.

## Support

Jika ada pertanyaan atau masalah, silakan buat issue di GitHub.

---

**Selamat menggunakan go-ros!** 🚀
