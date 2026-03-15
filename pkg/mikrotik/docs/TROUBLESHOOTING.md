# Troubleshooting

Dokumen ini berisi solusi untuk masalah umum yang mungkin dihadapi saat menggunakan go-ros.

## Table of Contents

1. [Connection Issues](#connection-issues)
2. [Authentication Issues](#authentication-issues)
3. [API Errors](#api-errors)
4. [Performance Issues](#performance-issues)
5. [Common Error Messages](#common-error-messages)
6. [Debug Tips](#debug-tips)

---

## Connection Issues

### Error: "connection refused"

**Error Message:**
```
dial tcp 192.168.88.1:8728: connectex: No connection could be made because the target machine actively refused it.
```

**Penyebab:**
1. API service tidak aktif di MikroTik
2. Port API berbeda
3. Firewall memblok koneksi
4. IP address salah

**Solusi:**

```bash
# 1. Cek API service di MikroTik
/ip service print

# 2. Enable API jika belum aktif
/ip service enable api

# 3. Cek port API
/ip service set api port=8728

# 4. Cek firewall rules
/ip firewall filter print

# 5. Tambahkan rule untuk allow API (jika perlu)
/ip firewall filter add chain=input protocol=tcp dst-port=8728 src-address=192.168.88.0/24 action=accept place-before=0
```

### Error: "connection timeout"

**Error Message:**
```
dial tcp 192.168.88.1:8728: i/o timeout
```

**Penyebab:**
1. Network unreachable
2. Router down
3. Firewall block
4. Wrong IP address

**Solusi:**

```bash
# 1. Test ping ke router
ping 192.168.88.1

# 2. Test port dengan telnet/nc
telnet 192.168.88.1 8728
# atau
nc -zv 192.168.88.1 8728

# 3. Cek IP address di MikroTik
/ip address print

# 4. Increase timeout di kode
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
```

### Error: "no route to host"

**Penyebab:**
1. Network configuration salah
2. Router di network berbeda
3. Gateway tidak tersedia

**Solusi:**

```bash
# Cek routing table
/ip route print

# Cek apakah router reachable
ping 192.168.88.1

# Cek network interface
/interface print
```

---

## Authentication Issues

### Error: "cannot log in"

**Error Message:**
```
from RouterOS device: cannot log in
```

**Penyebab:**
1. Username/password salah
2. User tidak aktif
3. User tidak punya permission

**Solusi:**

```bash
# 1. Cek user di MikroTik
/user print

# 2. Cek apakah user disabled
/user print where name="admin"

# 3. Reset password jika perlu
/user set admin password=newpassword

# 4. Cek group permission
/user group print

# 5. Pastikan user punya policy yang cukup
/user group set full policy=local,telnet,ssh,ftp,reboot,read,write,policy,test,winbox,password,web,sniff,sensitive,api,romon,dude,tikapp
```

### Error: "invalid user name or password"

**Solusi:**

```go
// Pastikan password benar
// Perhatikan: password case-sensitive!

cfg := client.Config{
    Username: "admin",        // pastikan lowercase/uppercase benar
    Password: "Password123",  // pastikan case benar
}
```

### Error: "user is not allowed"

**Penyebab:** User tidak punya permission untuk API.

**Solusi:**

```bash
# Cek group user
/user print detail

# Tambahkan permission API
/user group set [find name=yourgroup] policy=api,read,write

# Atau buat user baru dengan full permission
/user add name=apiuser group=full password=securepass
```

---

## API Errors

### Error: "no such command"

**Error Message:**
```
from RouterOS device: no such command
```

**Penyebab:**
1. Package tidak terinstall (misal: hotspot, ppp)
2. Command typo
3. RouterOS version tidak support

**Solusi:**

```bash
# 1. Cek installed packages
/system package print

# 2. Install package yang diperlukan (jika perlu)
# Download dari MikroTik website dan upload ke router

# 3. Cek RouterOS version
/system resource print
```

### Error: "input does not match any value"

**Error Message:**
```
from RouterOS device: input does not match any value of address-pool
```

**Penyebab:** Value yang diberikan tidak valid (misal: pool tidak ada).

**Solusi:**

```bash
# Cek value yang tersedia
/ip pool print

# Buat pool jika belum ada
/ip pool add name=hs-pool ranges=192.168.88.10-192.168.88.254
```

### Error: "bad command name"

**Penyebab:** Command path salah.

**Solusi:**

```bash
# Cek available commands dengan '?'
/ip ?
/ip hotspot ?
/ip hotspot user ?
```

### Error: "already have"

**Error Message:**
```
from RouterOS device: already have such name
```

**Penyebab:** Mencoba membuat entry dengan nama yang sudah ada.

**Solusi:**

```go
// Check existence terlebih dahulu
existing, err := repo.User().GetUserByName(ctx, "testuser")
if err == nil && existing != nil {
    // Update instead of create
    err = repo.User().UpdateUser(ctx, existing.ID, user)
} else {
    // Create new
    id, err = repo.User().AddUser(ctx, user)
}
```

### Error: "no such item"

**Error Message:**
```
from RouterOS device: no such item
```

**Penyebab:** ID yang diberikan tidak ditemukan.

**Solusi:**

```go
// Pastikan ID benar
users, err := repo.User().GetUsers(ctx, "")
for _, u := range users {
    fmt.Printf("ID: %s, Name: %s\n", u.ID, u.Name)
}

// Gunakan ID yang valid
err = repo.User().RemoveUser(ctx, "*1F")  // Pastikan ID benar
```

---

## Performance Issues

### Slow Response

**Penyebab:**
1. Network latency tinggi
2. Router CPU/Memory tinggi
3. Too many concurrent connections
4. Large data retrieval

**Solusi:**

```go
// 1. Gunakan proplist untuk limit fields
users, err := repo.User().GetUsers(ctx, "", ".id,name,profile")

// 2. Gunakan pagination (jika tersedia)
// Filter dengan profile
users, err := repo.User().GetUsers(ctx, "default")

// 3. Increase timeout untuk large operations
ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
defer cancel()

// 4. Gunakan goroutines untuk concurrent operations
var wg sync.WaitGroup
wg.Add(2)

go func() {
    defer wg.Done()
    users, _ := repo.User().GetUsers(ctx, "")
    // process
}()

go func() {
    defer wg.Done()
    profiles, _ := repo.Profile().GetProfiles(ctx)
    // process
}()

wg.Wait()
```

### Memory Leaks

**Penyebab:** Tidak menutup client/manager.

**Solusi:**

```go
// Selalu close resources
c, err := client.New(cfg)
if err != nil {
    return err
}
defer c.Close()  // Jangan lupa!

// Untuk manager
manager := client.NewManager(nil)
defer manager.CloseAll()  // Jangan lupa!
```

### High CPU Usage

**Penyebab:**
1. Terlalu banyak reconnect
2. Streaming tanpa debounce
3. Busy waiting

**Solusi:**

```go
// 1. Gunakan ListenBatches untuk debounce
batches := client.ListenBatches(ctx, sentences, 200*time.Millisecond)

// 2. Gunakan proper context cancellation
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

// 3. Jangan reconnect terlalu sering
// Library sudah handle dengan exponential backoff
```

---

## Common Error Messages

### "context deadline exceeded"

**Penyebab:** Operation timeout.

**Solusi:**

```go
// Increase timeout
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
```

### "context canceled"

**Penyebab:** Context di-cancel.

**Solusi:**

```go
// Jangan cancel context terlalu cepat
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
// ... do work
defer cancel()  // Cancel di akhir
```

### "not connected to mikrotik"

**Penyebab:** Client tidak terhubung.

**Solusi:**

```go
c, err := client.New(cfg)
if err != nil {
    return err
}

// Cek status
if !c.IsConnected() {
    return errors.New("not connected")
}
```

### "router %q not registered"

**Penyebab:** Mengakses router yang belum di-register di Manager.

**Solusi:**

```go
manager := client.NewManager(nil)

// Register terlebih dahulu
c, err := manager.GetOrConnect(ctx, "router-1", cfg)
if err != nil {
    return err
}

// Baru akses
c, err = manager.Get("router-1")  // atau MustGet
```

---

## Debug Tips

### Enable Debug Logging

```go
import (
    "go.uber.org/zap"
    "github.com/Butterfly-Student/go-ros/client"
)

// Create development logger
logger, _ := zap.NewDevelopment()

// Use with client
cfg := client.Config{...}
c := client.NewClient(cfg, logger)

// Atau dengan manager
manager := client.NewManager(logger)
```

### Print Raw Responses

```go
// Gunakan Raw methods untuk debug
raw, err := repo.User().PrintUsersRaw(ctx, ".id,name,profile")
for _, r := range raw {
    fmt.Printf("Raw: %+v\n", r)
}
```

### Test Connection

```go
// Test dengan command sederhana
reply, err := c.Run("/system/identity/print")
if err != nil {
    log.Printf("Connection test failed: %v", err)
    return
}

for _, re := range reply.Re {
    log.Printf("Router identity: %s", re.Map["name"])
}
```

### Check RouterOS Version

```go
sysRepo := system.NewRepository(c)
resources, err := sysRepo.Resources().GetResources(ctx)
if err != nil {
    return err
}

log.Printf("RouterOS Version: %s", resources.Version)
log.Printf("Platform: %s", resources.Platform)
```

### Monitor Connection State

```go
// Check connection state
log.Printf("Connected: %v", c.IsConnected())
log.Printf("Async mode: %v", c.IsAsync())

// Dengan manager
log.Printf("Registered routers: %v", manager.Names())
```

### Wireshark Capture

Untuk debug network issues:

```bash
# Capture traffic ke MikroTik
sudo tcpdump -i any host 192.168.88.1 and port 8728 -w mikrotik.pcap

# Analisis dengan Wireshark
# Filter: tcp.port == 8728
```

### Enable MikroTik Debug

```bash
# Enable debug logging di MikroTik
/system logging add topics=api action=memory

# Lihat log
/log print where topics~"api"
```

---

## FAQ

### Q: Bagaimana cara handle reconnect?

**A:** Library sudah handle auto-reconnect. Anda hanya perlu memastikan menggunakan Manager atau Client dengan benar.

```go
manager := client.NewManager(nil)
c, err := manager.GetOrConnect(ctx, "router-1", cfg)
// Auto-reconnect akan aktif jika koneksi terputus
```

### Q: Kenapa sering timeout?

**A:** Beberapa kemungkinan:
1. Network latency tinggi → increase timeout
2. Router sibuk → cek CPU/memory
3. Query terlalu besar → gunakan proplist

### Q: Bagaimana handle concurrent requests?

**A:** Client sudah thread-safe. Bisa digunakan concurrently.

```go
go func() { users, _ := repo.User().GetUsers(ctx, "") }()
go func() { profiles, _ := repo.Profile().GetProfiles(ctx) }()
```

### Q: Bagaimana stop streaming/monitoring?

**A:** Gunakan return value dari Listen methods.

```go
stop, err := repo.Active().ListenActive(ctx, resultChan)
if err != nil {
    return err
}
defer stop()  // Stop saat selesai

// Atau stop manual
stop()
```

### Q: Kenapa data tidak update real-time?

**A:** Gunakan Listen methods untuk real-time data.

```go
// Bukan real-time
users, _ := repo.User().GetUsers(ctx, "")

// Real-time
resultChan := make(chan []*domain.HotspotActive)
stop, _ := repo.Active().ListenActive(ctx, resultChan)
```

---

## Getting Help

Jika masalah masih belum terpecahkan:

1. **Cek Documentation:**
   - [MikroTik API Wiki](https://wiki.mikrotik.com/wiki/Manual:API)
   - [Go Documentation](https://golang.org/doc/)

2. **Enable Debug Logging:**
   ```go
   logger, _ := zap.NewDevelopment()
   manager := client.NewManager(logger)
   ```

3. **Test dengan RouterOS CLI:**
   ```bash
   # Cek apakah command berfungsi di CLI
   /ip hotspot user print
   ```

4. **Buat Issue:**
   - Sertakan error message lengkap
   - Sertakan RouterOS version
   - Sertakan code snippet
   - Sertakan debug logs

---

**Semoga troubleshooting guide ini membantu!** 🔧
