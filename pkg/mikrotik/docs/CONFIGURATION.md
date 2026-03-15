# Konfigurasi

Dokumen ini menjelaskan semua opsi konfigurasi yang tersedia dalam go-ros.

## Table of Contents

1. [Config Struct](#config-struct)
2. [Connection Options](#connection-options)
3. [TLS/SSL Configuration](#tlssl-configuration)
4. [Timeout Settings](#timeout-settings)
5. [Reconnection Settings](#reconnection-settings)
6. [Environment Variables](#environment-variables)
7. [Configuration Examples](#configuration-examples)

---

## Config Struct

Struct utama untuk konfigurasi koneksi ke MikroTik:

```go
type Config struct {
    Host              string
    Port              int
    Username          string
    Password          string
    UseTLS            bool
    ReconnectInterval time.Duration
    Timeout           time.Duration
    PoolSize          int
}
```

## Connection Options

### Host

Alamat IP atau hostname MikroTik router.

```go
cfg := client.Config{
    Host: "192.168.88.1",      // IPv4
    // atau
    Host: "router.local",       // Hostname
    // atau
    Host: "2001:db8::1",       // IPv6 (jika didukung)
}
```

### Port

Port untuk koneksi API MikroTik.

```go
cfg := client.Config{
    Port: 8728,    // API port (default)
    // atau
    Port: 8729,    // API-SSL port
}
```

**Port Standar MikroTik:**
- `8728` - API (plaintext)
- `8729` - API-SSL (encrypted)

### Username & Password

Kredensial untuk autentikasi.

```go
// Hardcoded (tidak direkomendasikan untuk production)
cfg := client.Config{
    Username: "admin",
    Password: "password123",
}

// Dari environment variables (direkomendasikan)
cfg := client.Config{
    Username: os.Getenv("MIKROTIK_USER"),
    Password: os.Getenv("MIKROTIK_PASS"),
}
```

---

## TLS/SSL Configuration

### Enable TLS

```go
cfg := client.Config{
    Host:     "192.168.88.1",
    Port:     8729,        // API-SSL port
    Username: "admin",
    Password: "password",
    UseTLS:   true,        // Enable TLS
}

c, err := client.New(cfg)
```

### TLS dengan Custom Config (Future)

```go
// Saat ini belum diimplementasikan
// TODO: Support untuk custom TLS config
```

### Enable API-SSL di MikroTik

```bash
# Enable API-SSL service
/ip service enable api-ssl

# Set port (default 8729)
/ip service set api-ssl port=8729

# Generate certificate (jika belum ada)
/certificate add name=api-cert common-name=api
/certificate sign api-cert
/ip service set api-ssl certificate=api-cert
```

---

## Timeout Settings

### Per-Command Timeout

Timeout untuk setiap command yang dijalankan.

```go
cfg := client.Config{
    Host:     "192.168.88.1",
    Port:     8728,
    Username: "admin",
    Password: "password",
    Timeout:  10 * time.Second,    // Default: 10 detik
}
```

**Rekomendasi Timeout:**

| Operasi | Timeout |
|---------|---------|
| Simple commands | 5-10 detik |
| List operations | 10-30 detik |
| Bulk operations | 30-60 detik |
| Streaming/monitoring | 0 (no timeout) |

### Context Timeout

Gunakan context untuk kontrol timeout yang lebih granular:

```go
// Short timeout untuk simple operations
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
users, err := repo.User().GetUsers(ctx, "")

// Longer timeout untuk bulk operations
ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
defer cancel()
err := voucherRepo.GenerateBatch(ctx, req)

// No timeout untuk streaming
ctx := context.Background()
stop, err := monitorRepo.System().StartSystemResourceMonitorListen(ctx, resultChan)
```

---

## Reconnection Settings

### Reconnect Interval

Interval awal untuk reconnect attempt.

```go
cfg := client.Config{
    Host:              "192.168.88.1",
    Port:              8728,
    Username:          "admin",
    Password:          "password",
    ReconnectInterval: 1 * time.Second,    // Default: 1 detik
}
```

### Exponential Backoff

Library menggunakan exponential backoff untuk reconnection:

```
Attempt 1: 1 detik
Attempt 2: 2 detik
Attempt 3: 4 detik
Attempt 4: 8 detik
Attempt 5: 16 detik
Attempt 6+: 30 detik (max)
```

**Konfigurasi Internal:**

```go
const (
    reconnectBaseDelay = 1 * time.Second   // Interval awal
    reconnectMaxDelay  = 30 * time.Second  // Interval maksimum
)
```

### Disable Auto-Reconnect

Saat ini tidak ada opsi untuk disable auto-reconnect. Namun, Anda bisa menggunakan `context.Cancel()` untuk menghentikan reconnect attempts.

---

## Environment Variables

### Best Practice untuk Credentials

```go
package main

import (
    "os"
    "strconv"
    "time"
    
    "github.com/Butterfly-Student/go-ros/client"
)

func loadConfigFromEnv() client.Config {
    // Required
    host := getEnvOrDefault("MIKROTIK_HOST", "192.168.88.1")
    username := getEnvOrDefault("MIKROTIK_USER", "admin")
    password := os.Getenv("MIKROTIK_PASS")
    
    if password == "" {
        panic("MIKROTIK_PASS environment variable is required")
    }
    
    // Optional dengan default values
    port, _ := strconv.Atoi(getEnvOrDefault("MIKROTIK_PORT", "8728"))
    timeout, _ := strconv.Atoi(getEnvOrDefault("MIKROTIK_TIMEOUT", "10"))
    useTLS, _ := strconv.ParseBool(getEnvOrDefault("MIKROTIK_TLS", "false"))
    
    return client.Config{
        Host:     host,
        Port:     port,
        Username: username,
        Password: password,
        UseTLS:   useTLS,
        Timeout:  time.Duration(timeout) * time.Second,
    }
}

func getEnvOrDefault(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func main() {
    cfg := loadConfigFromEnv()
    c, err := client.New(cfg)
    if err != nil {
        panic(err)
    }
    defer c.Close()
}
```

### Environment Variables List

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `MIKROTIK_HOST` | No | 192.168.88.1 | Router IP/hostname |
| `MIKROTIK_PORT` | No | 8728 | API port |
| `MIKROTIK_USER` | No | admin | Username |
| `MIKROTIK_PASS` | Yes | - | Password |
| `MIKROTIK_TLS` | No | false | Enable TLS |
| `MIKROTIK_TIMEOUT` | No | 10 | Timeout in seconds |

### .env File Example

```bash
# .env file
MIKROTIK_HOST=192.168.88.1
MIKROTIK_PORT=8728
MIKROTIK_USER=admin
MIKROTIK_PASS=your_secure_password
MIKROTIK_TLS=false
MIKROTIK_TIMEOUT=10
```

Load dengan library seperti `godotenv`:

```go
import "github.com/joho/godotenv"

func init() {
    godotenv.Load()
}
```

---

## Configuration Examples

### Development Configuration

```go
// config/development.go
package config

import (
    "time"
    "github.com/Butterfly-Student/go-ros/client"
)

func DevelopmentConfig() client.Config {
    return client.Config{
        Host:     "192.168.88.1",
        Port:     8728,
        Username: "admin",
        Password: "devpassword",
        Timeout:  30 * time.Second,  // Longer timeout untuk development
    }
}
```

### Production Configuration

```go
// config/production.go
package config

import (
    "os"
    "time"
    "github.com/Butterfly-Student/go-ros/client"
)

func ProductionConfig() client.Config {
    return client.Config{
        Host:     os.Getenv("MIKROTIK_HOST"),
        Port:     8729,  // API-SSL
        Username: os.Getenv("MIKROTIK_USER"),
        Password: os.Getenv("MIKROTIK_PASS"),
        UseTLS:   true,
        Timeout:  10 * time.Second,
    }
}
```

### Multi-Router Configuration

```go
// config/routers.go
package config

import (
    "os"
    "time"
    "github.com/Butterfly-Student/go-ros/client"
)

type RouterConfig struct {
    Name string
    Config client.Config
}

func LoadRouterConfigs() []RouterConfig {
    return []RouterConfig{
        {
            Name: "router-main",
            Config: client.Config{
                Host:     os.Getenv("ROUTER_MAIN_HOST"),
                Port:     8729,
                Username: os.Getenv("ROUTER_MAIN_USER"),
                Password: os.Getenv("ROUTER_MAIN_PASS"),
                UseTLS:   true,
                Timeout:  10 * time.Second,
            },
        },
        {
            Name: "router-backup",
            Config: client.Config{
                Host:     os.Getenv("ROUTER_BACKUP_HOST"),
                Port:     8729,
                Username: os.Getenv("ROUTER_BACKUP_USER"),
                Password: os.Getenv("ROUTER_BACKUP_PASS"),
                UseTLS:   true,
                Timeout:  10 * time.Second,
            },
        },
    }
}
```

### Configuration dengan Validation

```go
package config

import (
    "errors"
    "fmt"
    "time"
    "github.com/Butterfly-Student/go-ros/client"
)

func ValidateConfig(cfg client.Config) error {
    if cfg.Host == "" {
        return errors.New("host is required")
    }
    
    if cfg.Port <= 0 || cfg.Port > 65535 {
        return errors.New("invalid port number")
    }
    
    if cfg.Username == "" {
        return errors.New("username is required")
    }
    
    if cfg.Password == "" {
        return errors.New("password is required")
    }
    
    if cfg.Timeout <= 0 {
        return errors.New("timeout must be positive")
    }
    
    if cfg.UseTLS && cfg.Port != 8729 {
        fmt.Println("Warning: TLS enabled but port is not 8729")
    }
    
    return nil
}

func SafeNewClient(cfg client.Config) (*client.Client, error) {
    if err := ValidateConfig(cfg); err != nil {
        return nil, err
    }
    
    return client.New(cfg)
}
```

---

## Advanced Configuration

### Custom Logger

```go
import (
    "go.uber.org/zap"
    "github.com/Butterfly-Student/go-ros/client"
)

// Create custom logger
logger, _ := zap.NewDevelopment()

// Use with Manager
manager := client.NewManager(logger)

// Use with Client
cfg := client.Config{...}
c := client.NewClient(cfg, logger)
```

### Connection Pool (Future)

```go
// Saat ini PoolSize field ada tapi belum diimplementasikan
// TODO: Implement connection pooling

cfg := client.Config{
    Host:     "192.168.88.1",
    Port:     8728,
    PoolSize: 10,  // Future use
}
```

---

## Security Best Practices

### 1. Jangan Hardcode Password

❌ **Jangan:**
```go
cfg := client.Config{
    Password: "admin123",  // Bahaya!
}
```

✅ **Lakukan:**
```go
cfg := client.Config{
    Password: os.Getenv("MIKROTIK_PASS"),
}
```

### 2. Gunakan TLS untuk Production

```go
// Production
prodCfg := client.Config{
    Host:     "router.company.com",
    Port:     8729,
    UseTLS:   true,
    Username: os.Getenv("MIKROTIK_USER"),
    Password: os.Getenv("MIKROTIK_PASS"),
}

// Development
devCfg := client.Config{
    Host:     "192.168.88.1",
    Port:     8728,
    UseTLS:   false,
    Username: "admin",
    Password: "devpass",
}
```

### 3. Batasi Akses dengan Firewall

```bash
# Di MikroTik - hanya allow dari network tertentu
/ip firewall filter add chain=input protocol=tcp dst-port=8728 src-address=192.168.88.0/24 action=accept
/ip firewall filter add chain=input protocol=tcp dst-port=8728 action=drop
```

### 4. Gunakan User dengan Privilege Minimal

```bash
# Buat user khusus untuk API
/user add name=apiuser group=read password=securepass

# Atau buat group custom dengan policy yang dibutuhkan saja
/user group add name=apigroup policy=read,write
/user add name=apiuser group=apigroup password=securepass
```

---

## Troubleshooting Configuration

### Connection Refused

**Penyebab:**
- API service tidak aktif
- Port salah
- Firewall block

**Solusi:**
```bash
# Cek API service
/ip service print

# Enable API
/ip service enable api

# Cek firewall
/ip firewall filter print
```

### Authentication Failed

**Penyebab:**
- Username/password salah
- User tidak punya permission

**Solusi:**
```bash
# Cek user
/user print

# Cek group permission
/user group print
```

### Timeout

**Penyebab:**
- Network latency tinggi
- Router sibuk
- Timeout terlalu pendek

**Solusi:**
```go
// Increase timeout
cfg := client.Config{
    Timeout: 30 * time.Second,
}
```

---

**Referensi:**
- [MikroTik API Documentation](https://wiki.mikrotik.com/wiki/Manual:API)
- [Go Context Package](https://golang.org/pkg/context/)
- [Go Time Package](https://golang.org/pkg/time/)
