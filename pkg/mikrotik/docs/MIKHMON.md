# Panduan Lengkap Mikhmon

Dokumen ini menjelaskan cara menggunakan fitur Mikhmon dalam go-ros untuk manajemen voucher hotspot.

## Table of Contents

1. [Apa itu Mikhmon?](#apa-itu-mikhmon)
2. [Fitur Mikhmon](#fitur-mikhmon)
3. [Setup Awal](#setup-awal)
4. [Voucher Generation](#voucher-generation)
5. [Profile Management](#profile-management)
6. [Sales Reports](#sales-reports)
7. [Expire Monitoring](#expire-monitoring)
8. [Contoh Lengkap](#contoh-lengkap)

---

## Apa itu Mikhmon?

Mikhmon adalah sistem manajemen voucher hotspot untuk MikroTik RouterOS. Fitur-fiturnya meliputi:

- **Voucher Generation** - Generate voucher dengan berbagai format
- **User Expiration** - Otomatis disable/remove user saat expired
- **Sales Reporting** - Tracking penjualan voucher
- **MAC Locking** - Lock user ke MAC address tertentu
- **Multi-server** - Support multiple hotspot server

## Fitur Mikhmon

### 1. Voucher Generation

Generate voucher dalam dua mode:

- **VC (Voucher Card)** - Username = Password
- **UP (User/Password)** - Username ≠ Password

### 2. Expire Modes

- **REM (Remove)** - Hapus user saat expired
- **NTF (Notify)** - Disable user saat expired
- **REMC (Remove + Record)** - Hapus + catat ke report
- **NTFC (Notify + Record)** - Disable + catat ke report
- **0 (No Expire)** - Tidak ada expiration

### 3. Sales Reports

Otomatis catat setiap penjualan ke `/system/script` dengan format:
```
date-|-time-|-user-|-price-|-ip-|-mac-|-validity-|-profile-|-comment
```

### 4. Expire Monitor

Scheduler yang berjalan setiap menit untuk check dan proses user expired.

---

## Setup Awal

### 1. Enable API Service

```bash
/ip service enable api
/ip service set api port=8728
```

### 2. Buat Hotspot Profile

```bash
/ip hotspot profile add name=default hotspot-address=192.168.88.1
```

### 3. Buat IP Pool

```bash
/ip pool add name=hs-pool ranges=192.168.88.10-192.168.88.254
```

### 4. Buat Hotspot Server

```bash
/ip hotspot add name=hotspot1 interface=ether2 address-pool=hs-pool profile=default
```

---

## Voucher Generation

### Basic Voucher Generation

```go
package main

import (
    "context"
    "fmt"
    "time"
    
    "github.com/Butterfly-Student/go-ros/client"
    "github.com/Butterfly-Student/go-ros/repository/hotspot"
    mikhmonRepo "github.com/Butterfly-Student/go-ros/repository/mikhmon"
    mikhmonDomain "github.com/Butterfly-Student/go-ros/domain/mikhmon"
)

func main() {
    ctx := context.Background()
    
    // Connect ke MikroTik
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
    
    // Setup repositories
    hotspotRepo := hotspot.NewRepository(c)
    generatorRepo := mikhmonRepo.NewGeneratorRepository()
    voucherRepo := mikhmonRepo.NewVoucherRepository(c, hotspotRepo, generatorRepo)
    
    // Generate voucher VC mode
    req := &mikhmonDomain.VoucherGenerateRequest{
        Quantity:   10,
        Profile:    "default",
        Mode:       mikhmonDomain.VoucherModeVoucher, // "vc"
        NameLength: 6,
        CharSet:    mikhmonDomain.CharSetUpplow1,
        TimeLimit:  "1h",
        DataLimit:  "1G",
    }
    
    batch, err := voucherRepo.GenerateBatch(ctx, req)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Generated %d vouchers with code: %s\n", batch.Quantity, batch.Code)
    for _, v := range batch.Vouchers {
        fmt.Printf("  %s / %s\n", v.Name, v.Password)
    }
}
```

### Voucher dengan Prefix

```go
req := &mikhmonDomain.VoucherGenerateRequest{
    Quantity:   5,
    Profile:    "default",
    Mode:       mikhmonDomain.VoucherModeVoucher,
    NameLength: 4,
    Prefix:     "VC",
    CharSet:    mikhmonDomain.CharSetNumeric,
    TimeLimit:  "30m",
}

// Hasil: VC1234, VC5678, dll
```

### User/Password Mode

```go
req := &mikhmonDomain.VoucherGenerateRequest{
    Quantity:   5,
    Profile:    "default",
    Mode:       mikhmonDomain.VoucherModeUserPassword, // "up"
    NameLength: 6,
    CharSet:    mikhmonDomain.CharSetUpplow1,
    TimeLimit:  "2h",
}

// Hasil: Username dan password berbeda
// user1: SAnUzo / JB4nyS
// user2: TfLssG / 75lDjN
```

### CharSet Options

```go
// Lowercase only
CharSet: mikhmonDomain.CharSetLower   // abcdef

// Uppercase only  
CharSet: mikhmonDomain.CharSetUpper   // ABCDEF

// Mixed case
CharSet: mikhmonDomain.CharSetUpplow  // AbCdEf

// With numbers
CharSet: mikhmonDomain.CharSetUpplow1 // AbC1d2

// With special chars
CharSet: mikhmonDomain.CharSetMix1    // Ab1@cd
```

### List dan Remove Vouchers

```go
// List vouchers by comment
vouchers, err := voucherRepo.GetVouchersByComment(ctx, "vc-123-01.02.26")
if err != nil {
    panic(err)
}

for _, v := range vouchers {
    fmt.Printf("Voucher: %s, Profile: %s\n", v.Name, v.Profile)
}

// Remove voucher batch
err = voucherRepo.RemoveVoucherBatch(ctx, "vc-123-01.02.26")
if err != nil {
    panic(err)
}
```

---

## Profile Management

### Create Profile dengan Mikhmon Config

```go
profileRepo := mikhmonRepo.NewProfileRepository(hotspotRepo)

// Buat request
req := &mikhmonDomain.ProfileRequest{
    Name:        "Paket-1Jam",
    AddressPool: "hs-pool",
    RateLimit:   "1M/2M",
    SharedUsers: 1,
    Config: mikhmonDomain.ProfileConfig{
        Name:         "Paket-1Jam",
        Price:        5000,
        SellingPrice: 7000,
        Validity:     "1h",
        ExpireMode:   mikhmonDomain.ExpireModeRemove, // "rem"
        LockUser:     false,
        LockServer:   false,
    },
}

// Create profile
err := profileRepo.CreateProfile(ctx, req)
if err != nil {
    panic(err)
}

fmt.Println("Profile created successfully!")
```

### Generate On-Login Script Manual

```go
scriptData := &mikhmonDomain.OnLoginScriptData{
    Mode:         mikhmonDomain.ExpireModeRemove,
    Price:        5000,
    Validity:     "1h",
    SellingPrice: 7000,
    NoExp:        false,
    LockUser:     "Disable",
    LockServer:   "Disable",
}

script := profileRepo.GenerateOnLoginScript(scriptData)
fmt.Println(script)
```

### Profile dengan MAC Locking

```go
req := &mikhmonDomain.ProfileRequest{
    Name:        "Paket-1Jam-Locked",
    AddressPool: "hs-pool",
    RateLimit:   "1M/2M",
    SharedUsers: 1,
    Config: mikhmonDomain.ProfileConfig{
        Name:         "Paket-1Jam-Locked",
        Price:        5000,
        SellingPrice: 7000,
        Validity:     "1h",
        ExpireMode:   mikhmonDomain.ExpireModeRemove,
        LockUser:     true,  // Enable MAC locking
        LockServer:   false,
    },
}

err := profileRepo.CreateProfile(ctx, req)
```

### Update Profile

```go
// Get profile ID terlebih dahulu
profiles, _ := hotspotRepo.Profile().GetProfiles(ctx)
var profileID string
for _, p := range profiles {
    if p.Name == "Paket-1Jam" {
        profileID = p.ID
        break
    }
}

// Update dengan config baru
req := &mikhmonDomain.ProfileRequest{
    Name:        "Paket-1Jam",
    AddressPool: "hs-pool",
    RateLimit:   "2M/4M",  // Update speed
    SharedUsers: 2,
    Config: mikhmonDomain.ProfileConfig{
        Name:         "Paket-1Jam",
        Price:        6000,  // Update price
        SellingPrice: 8000,
        Validity:     "1h",
        ExpireMode:   mikhmonDomain.ExpireModeRemove,
    },
}

err := profileRepo.UpdateProfile(ctx, profileID, req)
```

---

## Sales Reports

### Struktur Report

Report disimpan di `/system/script` dengan nama format:
```
date-|-time-|-user-|-price-|-ip-|-mac-|-validity-|-profile-|-comment
```

Contoh:
```
jan/15/2024-|-14:30:25-|-user123-|-5000-|-192.168.88.100-|-AA:BB:CC:DD:EE:FF-|-1h-|-default-|-vc-123-01.15.24
```

### Membaca Reports

```go
sysRepo := system.NewRepository(c)

// Get all scripts
scripts, err := sysRepo.Scripts().GetScripts(ctx)
if err != nil {
    panic(err)
}

// Filter mikhmon reports
for _, script := range scripts {
    if script.Comment == "mikhmon" {
        fmt.Printf("Report: %s\n", script.Name)
        fmt.Printf("  Owner: %s\n", script.Owner)
        fmt.Printf("  Source: %s\n", script.Source)
    }
}
```

### Parse Report Data

```go
// Parse dari nama script
func parseReport(name string) *SalesReport {
    parts := strings.Split(name, "-|-")
    if len(parts) >= 9 {
        return &SalesReport{
            Date:     parts[0],
            Time:     parts[1],
            User:     parts[2],
            Price:    parseInt(parts[3]),
            IP:       parts[4],
            MAC:      parts[5],
            Validity: parts[6],
            Profile:  parts[7],
            Comment:  parts[8],
        }
    }
    return nil
}
```

### Report Summary

```go
func getReportSummary(scripts []*domain.SystemScript) {
    var totalSales int64
    var totalCount int
    
    for _, script := range scripts {
        if script.Comment == "mikhmon" {
            report := parseReport(script.Name)
            if report != nil {
                totalSales += report.Price
                totalCount++
            }
        }
    }
    
    fmt.Printf("Total Sales: %d\n", totalCount)
    fmt.Printf("Total Revenue: Rp %d\n", totalSales)
}
```

---

## Expire Monitoring

### Setup Expire Monitor

```go
expireRepo := mikhmonRepo.NewExpireRepository(c, sysRepo)

// Check status
enabled, err := expireRepo.IsExpireMonitorEnabled(ctx)
if err != nil {
    panic(err)
}

if enabled {
    fmt.Println("Expire monitor sudah aktif")
} else {
    fmt.Println("Expire monitor belum aktif")
}

// Setup expire monitor
err = expireRepo.SetupExpireMonitor(ctx)
if err != nil {
    panic(err)
}

fmt.Println("Expire monitor berhasil di-setup!")
```

### Disable Expire Monitor

```go
err := expireRepo.DisableExpireMonitor(ctx)
if err != nil {
    panic(err)
}
fmt.Println("Expire monitor dinonaktifkan")
```

### Cara Kerja Expire Monitor

1. **Scheduler** berjalan setiap 1 menit
2. **Check** semua user dengan comment mengandung tanggal
3. **Parse** tanggal expire dari comment
4. **Compare** dengan tanggal sekarang
5. **Action** berdasarkan mode:
   - Mode "N" (Notify): Set limit-uptime=1s (disable)
   - Mode "X" (Remove): Hapus user

### Format Comment

Comment user setelah login:
```
jan/15/2024 14:30:25 X
```

- `jan/15/2024` - Tanggal expire
- `14:30:25` - Waktu expire
- `X` - Mode (X=Remove, N=Notify)

---

## Contoh Lengkap

### Aplikasi Voucher Generator Lengkap

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

type VoucherApp struct {
    client      *client.Client
    hotspotRepo hotspot.Repository
    sysRepo     system.Repository
    voucherRepo mikhmonRepo.VoucherRepository
    profileRepo mikhmonRepo.ProfileRepository
    expireRepo  mikhmonRepo.ExpireRepository
}

func NewVoucherApp(cfg client.Config) (*VoucherApp, error) {
    c, err := client.New(cfg)
    if err != nil {
        return nil, err
    }
    
    hotspotRepo := hotspot.NewRepository(c)
    sysRepo := system.NewRepository(c)
    generatorRepo := mikhmonRepo.NewGeneratorRepository()
    
    return &VoucherApp{
        client:      c,
        hotspotRepo: hotspotRepo,
        sysRepo:     sysRepo,
        voucherRepo: mikhmonRepo.NewVoucherRepository(c, hotspotRepo, generatorRepo),
        profileRepo: mikhmonRepo.NewProfileRepository(hotspotRepo),
        expireRepo:  mikhmonRepo.NewExpireRepository(c, sysRepo),
    }, nil
}

func (app *VoucherApp) Close() {
    app.client.Close()
}

func (app *VoucherApp) GenerateVouchers(quantity int, profile string) (*mikhmonDomain.VoucherBatch, error) {
    ctx := context.Background()
    
    req := &mikhmonDomain.VoucherGenerateRequest{
        Quantity:   quantity,
        Profile:    profile,
        Mode:       mikhmonDomain.VoucherModeVoucher,
        NameLength: 6,
        CharSet:    mikhmonDomain.CharSetUpplow1,
        TimeLimit:  "1h",
    }
    
    return app.voucherRepo.GenerateBatch(ctx, req)
}

func (app *VoucherApp) PrintVouchers(batch *mikhmonDomain.VoucherBatch) {
    fmt.Println("╔════════════════════════════════════════╗")
    fmt.Println("║         VOUCHER HOTSPOT               ║")
    fmt.Println("╠════════════════════════════════════════╣")
    fmt.Printf("║ Kode: %-32s ║\n", batch.Code)
    fmt.Printf("║ Profile: %-29s ║\n", batch.Profile)
    fmt.Printf("║ Jumlah: %-30d ║\n", batch.Quantity)
    fmt.Printf("║ Time Limit: %-26s ║\n", batch.TimeLimit)
    fmt.Println("╠════════════════════════════════════════╣")
    
    for i, v := range batch.Vouchers {
        fmt.Printf("║ %2d. %-10s / %-10s     ║\n", i+1, v.Name, v.Password)
    }
    
    fmt.Println("╚════════════════════════════════════════╝")
}

func (app *VoucherApp) SetupExpireMonitor() error {
    ctx := context.Background()
    return app.expireRepo.SetupExpireMonitor(ctx)
}

func (app *VoucherApp) GetSalesReport() error {
    ctx := context.Background()
    
    scripts, err := app.sysRepo.Scripts().GetScripts(ctx)
    if err != nil {
        return err
    }
    
    fmt.Println("\n=== SALES REPORT ===")
    total := int64(0)
    count := 0
    
    for _, script := range scripts {
        if script.Comment == "mikhmon" {
            fmt.Printf("• %s\n", script.Name)
            count++
        }
    }
    
    fmt.Printf("\nTotal Transaksi: %d\n", count)
    fmt.Printf("Total Pendapatan: Rp %d\n", total)
    
    return nil
}

func main() {
    cfg := client.Config{
        Host:     "192.168.88.1",
        Port:     8728,
        Username: "admin",
        Password: os.Getenv("MIKROTIK_PASSWORD"),
        Timeout:  10 * time.Second,
    }
    
    app, err := NewVoucherApp(cfg)
    if err != nil {
        fmt.Printf("Failed to connect: %v\n", err)
        os.Exit(1)
    }
    defer app.Close()
    
    fmt.Println("✓ Connected to MikroTik")
    
    // Setup expire monitor
    fmt.Println("Setting up expire monitor...")
    if err := app.SetupExpireMonitor(); err != nil {
        fmt.Printf("Warning: %v\n", err)
    } else {
        fmt.Println("✓ Expire monitor active")
    }
    
    // Generate vouchers
    fmt.Println("\nGenerating vouchers...")
    batch, err := app.GenerateVouchers(5, "default")
    if err != nil {
        fmt.Printf("Failed to generate: %v\n", err)
        os.Exit(1)
    }
    
    // Print vouchers
    app.PrintVouchers(batch)
    
    // Show reports
    app.GetSalesReport()
}
```

### Best Practices

1. **Selalu cleanup test data:**
```go
defer func() {
    voucherRepo.RemoveVoucherBatch(ctx, batch.Vouchers[0].Comment)
}()
```

2. **Gunakan context dengan timeout:**
```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
```

3. **Validasi input:**
```go
if req.Quantity < 1 || req.Quantity > 1000 {
    return errors.New("quantity must be between 1-1000")
}
```

4. **Handle errors dengan baik:**
```go
batch, err := voucherRepo.GenerateBatch(ctx, req)
if err != nil {
    log.Printf("Failed to generate vouchers: %v", err)
    return err
}
```

5. **Gunakan environment variables untuk credentials:**
```go
password := os.Getenv("MIKROTIK_PASSWORD")
```

---

## Troubleshooting

### Error: "input does not match any value of address-pool"

**Penyebab:** Address pool yang diminta tidak ada di MikroTik.

**Solusi:**
```bash
# Cek pool yang tersedia
/ip pool print

# Buat pool jika belum ada
/ip pool add name=hs-pool ranges=192.168.88.10-192.168.88.254
```

### Error: "no such command"

**Penyebab:** Hotspot package tidak terinstall.

**Solusi:**
```bash
# Cek package
/system package print

# Install hotspot jika belum ada
```

### Voucher tidak muncul di Hotspot

**Penyebab:** 
1. Profile tidak sesuai
2. Server hotspot tidak aktif

**Solusi:**
```bash
# Cek hotspot server
/ip hotspot print

# Cek profile
/ip hotspot user profile print
```

### Expire monitor tidak berjalan

**Penyebab:** Scheduler tidak aktif.

**Solusi:**
```bash
# Cek scheduler
/system scheduler print

# Enable scheduler
/system scheduler enable [find name="Mikhmon-Expire-Monitor"]
```

---

## Tips dan Trik

### 1. Batch Voucher dengan Comment Custom

```go
// Generate dengan comment custom
req.Comment = "Promo-Weekend"
// Hasil comment: vc-123-01.02.26-Promo-Weekend
```

### 2. Generate Voucher dengan Validity Berbeda

```go
// 30 menit
req.TimeLimit = "30m"

// 1 hari
req.TimeLimit = "1d"

// 1 minggu
req.TimeLimit = "1w"

// 1 bulan
req.TimeLimit = "30d"
```

### 3. Data Limit

```go
// 500 MB
req.DataLimit = "500M"

// 1 GB
req.DataLimit = "1G"

// 10 GB
req.DataLimit = "10G"
```

### 4. Export Voucher ke CSV

```go
func exportToCSV(batch *mikhmonDomain.VoucherBatch, filename string) error {
    file, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer file.Close()
    
    writer := csv.NewWriter(file)
    defer writer.Flush()
    
    // Header
    writer.Write([]string{"No", "Username", "Password", "Profile"})
    
    // Data
    for i, v := range batch.Vouchers {
        writer.Write([]string{
            fmt.Sprintf("%d", i+1),
            v.Name,
            v.Password,
            v.Profile,
        })
    }
    
    return nil
}
```

---

**Selamat menggunakan Mikhmon dengan go-ros!** 🎉
