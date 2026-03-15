# Panduan Repository

Dokumen ini menjelaskan cara menggunakan setiap repository dalam go-ros dengan contoh lengkap.

## Table of Contents

1. [Hotspot Repository](#hotspot-repository)
2. [PPP Repository](#ppp-repository)
3. [System Repository](#system-repository)
4. [IP Address Repository](#ip-address-repository)
5. [Firewall Repository](#firewall-repository)
6. [Queue Repository](#queue-repository)
7. [Monitor Repository](#monitor-repository)
8. [Mikhmon Repository](#mikhmon-repository)

---

## Hotspot Repository

Repository untuk manajemen Hotspot di MikroTik.

### Setup

```go
import (
    "github.com/Butterfly-Student/go-ros/client"
    "github.com/Butterfly-Student/go-ros/repository/hotspot"
)

// Create client
c, err := client.New(cfg)
if err != nil {
    panic(err)
}
defer c.Close()

// Create repository
hotspotRepo := hotspot.NewRepository(c)
```

### User Management

#### Get Users

```go
ctx := context.Background()

// Get all users
users, err := hotspotRepo.User().GetUsers(ctx, "")
if err != nil {
    log.Fatal(err)
}

for _, user := range users {
    fmt.Printf("User: %s, Profile: %s, Disabled: %v\n", 
        user.Name, user.Profile, user.Disabled)
}

// Get users by profile
users, err = hotspotRepo.User().GetUsers(ctx, "default")

// Get users with specific fields only (optimization)
users, err = hotspotRepo.User().GetUsers(ctx, "", ".id,name,profile,disabled")
```

#### Get User by Name/ID

```go
// By name
user, err := hotspotRepo.User().GetUserByName(ctx, "testuser")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("User: %+v\n", user)

// By ID
user, err = hotspotRepo.User().GetUserByID(ctx, "*1F")
```

#### Add User

```go
user := &domain.HotspotUser{
    Name:            "newuser",
    Password:        "securepass123",
    Profile:         "default",
    Comment:         "Created via API",
    LimitUptime:     "1h",
    LimitBytesTotal: 1073741824, // 1GB
}

id, err := hotspotRepo.User().AddUser(ctx, user)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("User created with ID: %s\n", id)
```

#### Update User

```go
// Get user ID first
user, _ := hotspotRepo.User().GetUserByName(ctx, "newuser")

// Update fields
user.Password = "newpassword"
user.Profile = "premium"
user.Comment = "Updated profile"

err := hotspotRepo.User().UpdateUser(ctx, user.ID, user)
if err != nil {
    log.Fatal(err)
}
```

#### Remove User

```go
// Remove single user
err := hotspotRepo.User().RemoveUser(ctx, "*1F")

// Remove multiple users
err = hotspotRepo.User().RemoveUsers(ctx, []string{"*1F", "*20", "*21"})

// Remove users by comment
err = hotspotRepo.User().RemoveUsersByComment(ctx, "test-users")
```

#### Enable/Disable User

```go
// Disable single user
err := hotspotRepo.User().DisableUser(ctx, "*1F")

// Enable single user
err = hotspotRepo.User().EnableUser(ctx, "*1F")

// Batch disable
err = hotspotRepo.User().DisableUsers(ctx, []string{"*1F", "*20", "*21"})

// Batch enable
err = hotspotRepo.User().EnableUsers(ctx, []string{"*1F", "*20", "*21"})
```

#### Reset Counters

```go
// Reset single user
err := hotspotRepo.User().ResetUserCounters(ctx, "*1F")

// Reset multiple users
err = hotspotRepo.User().ResetUserCountersMultiple(ctx, []string{"*1F", "*20"})
```

### Profile Management

#### Get Profiles

```go
// Get all profiles
profiles, err := hotspotRepo.Profile().GetProfiles(ctx)
if err != nil {
    log.Fatal(err)
}

for _, profile := range profiles {
    fmt.Printf("Profile: %s, Pool: %s, Rate: %s\n",
        profile.Name, profile.AddressPool, profile.RateLimit)
}

// Get by name
profile, err := hotspotRepo.Profile().GetProfileByName(ctx, "default")

// Get by ID
profile, err = hotspotRepo.Profile().GetProfileByID(ctx, "*1")
```

#### Create Profile

```go
profile := &domain.UserProfile{
    Name:              "Premium-10M",
    AddressPool:       "hs-pool",
    RateLimit:         "10M/10M",
    SharedUsers:       3,
    StatusAutorefresh: "1m",
}

id, err := hotspotRepo.Profile().AddProfile(ctx, profile)
if err != nil {
    log.Fatal(err)
}
```

#### Update Profile

```go
profile.RateLimit = "20M/20M"
err := hotspotRepo.Profile().UpdateProfile(ctx, "*1", profile)
```

#### Remove Profile

```go
err := hotspotRepo.Profile().RemoveProfile(ctx, "*1")

// Batch remove
err = hotspotRepo.Profile().RemoveProfiles(ctx, []string{"*1", "*2"})
```

#### Enable/Disable Profile

```go
err := hotspotRepo.Profile().DisableProfile(ctx, "*1")
err = hotspotRepo.Profile().EnableProfile(ctx, "*1")

// Batch
err = hotspotRepo.Profile().DisableProfiles(ctx, []string{"*1", "*2"})
err = hotspotRepo.Profile().EnableProfiles(ctx, []string{"*1", "*2"})
```

### Active Sessions

#### Get Active Users

```go
active, err := hotspotRepo.Active().GetActive(ctx)
if err != nil {
    log.Fatal(err)
}

for _, a := range active {
    fmt.Printf("User: %s, IP: %s, MAC: %s, Uptime: %s\n",
        a.User, a.Address, a.MACAddress, a.Uptime)
}

// Get count
count, err := hotspotRepo.Active().GetActiveCount(ctx)
fmt.Printf("Active users: %d\n", count)
```

#### Remove Active Session

```go
// Kick user
err := hotspotRepo.Active().RemoveActive(ctx, "*1F")

// Batch kick
err = hotspotRepo.Active().RemoveActives(ctx, []string{"*1F", "*20"})
```

#### Monitor Active Users (Real-time)

```go
resultChan := make(chan []*domain.HotspotActive)

stop, err := hotspotRepo.Active().ListenActive(ctx, resultChan)
if err != nil {
    log.Fatal(err)
}
defer stop()

for active := range resultChan {
    fmt.Printf("Active users: %d\n", len(active))
    for _, a := range active {
        fmt.Printf("  - %s (%s)\n", a.User, a.Address)
    }
}
```

### Host Management

```go
// Get all hosts
hosts, err := hotspotRepo.Host().GetHosts(ctx)
if err != nil {
    log.Fatal(err)
}

for _, host := range hosts {
    fmt.Printf("Host: %s, MAC: %s, Server: %s\n",
        host.Address, host.MACAddress, host.Server)
}

// Remove host
err = hotspotRepo.Host().RemoveHost(ctx, "*1F")
```

### IP Binding

```go
// Get all bindings
bindings, err := hotspotRepo.IPBinding().GetIPBindings(ctx)

// Add binding
binding := &domain.HotspotIPBinding{
    MACAddress: "AA:BB:CC:DD:EE:FF",
    Address:    "192.168.88.100",
    Server:     "hotspot1",
    Type:       "bypassed",
    Comment:    "VIP User",
}
id, err := hotspotRepo.IPBinding().AddIPBinding(ctx, binding)

// Remove binding
err = hotspotRepo.IPBinding().RemoveIPBinding(ctx, "*1F")

// Enable/Disable
err = hotspotRepo.IPBinding().EnableIPBinding(ctx, "*1F")
err = hotspotRepo.IPBinding().DisableIPBinding(ctx, "*1F")
```

### Server Management

```go
// Get all hotspot servers
servers, err := hotspotRepo.Server().GetServers(ctx)
if err != nil {
    log.Fatal(err)
}

for _, server := range servers {
    fmt.Printf("Server: %s\n", server)
}
```

---

## PPP Repository

Repository untuk manajemen PPP (VPN) connections.

### Setup

```go
import "github.com/Butterfly-Student/go-ros/repository/ppp"

pppRepo := ppp.NewRepository(c)
```

### Secret Management

#### Get Secrets

```go
// Get all secrets
secrets, err := pppRepo.Secret().GetSecrets(ctx, "")

// Get by profile
secrets, err = pppRepo.Secret().GetSecrets(ctx, "default")

// Get by name
secret, err := pppRepo.Secret().GetSecretByName(ctx, "vpnuser")

// Get by ID
secret, err = pppRepo.Secret().GetSecretByID(ctx, "*1F")
```

#### Add Secret

```go
secret := &domain.PPPSecret{
    Name:          "vpnuser",
    Password:      "vpnpass123",
    Profile:       "default",
    Service:       "any",
    RemoteAddress: "10.0.0.100",
    Comment:       "VPN User",
}

err := pppRepo.Secret().AddSecret(ctx, secret)
```

#### Update Secret

```go
secret.Password = "newpass"
err := pppRepo.Secret().UpdateSecret(ctx, "*1F", secret)
```

#### Remove Secret

```go
err := pppRepo.Secret().RemoveSecret(ctx, "*1F")
err = pppRepo.Secret().RemoveSecrets(ctx, []string{"*1F", "*20"})
```

#### Enable/Disable

```go
err := pppRepo.Secret().DisableSecret(ctx, "*1F")
err = pppRepo.Secret().EnableSecret(ctx, "*1F")

// Batch
err = pppRepo.Secret().DisableSecrets(ctx, []string{"*1F", "*20"})
err = pppRepo.Secret().EnableSecrets(ctx, []string{"*1F", "*20"})
```

### Profile Management

```go
// Get profiles
profiles, err := pppRepo.Profile().GetProfiles(ctx)

// Get by name
profile, err := pppRepo.Profile().GetProfileByName(ctx, "default")

// Create profile
newProfile := &domain.PPPProfile{
    Name:           "Premium-VPN",
    LocalAddress:   "10.0.0.1",
    RemoteAddress:  "10.0.0.100-10.0.0.200",
    RateLimit:      "10M/10M",
    UseCompression: true,
    UseEncryption:  true,
}
err := pppRepo.Profile().AddProfile(ctx, newProfile)

// Update
err = pppRepo.Profile().UpdateProfile(ctx, "*1", newProfile)

// Remove
err = pppRepo.Profile().RemoveProfile(ctx, "*1")
```

### Active Sessions

```go
// Get active
active, err := pppRepo.Active().GetActive(ctx, "any")

// Get by name
session, err := pppRepo.Active().GetActiveByName(ctx, "vpnuser")

// Remove
err = pppRepo.Active().RemoveActive(ctx, "*1F")
err = pppRepo.Active().RemoveActives(ctx, []string{"*1F", "*20"})
```

---

## System Repository

Repository untuk manajemen sistem MikroTik.

### Setup

```go
import "github.com/Butterfly-Student/go-ros/repository/system"

sysRepo := system.NewRepository(c)
```

### Identity

```go
// Get identity
identity, err := sysRepo.Identity().GetIdentity(ctx)
fmt.Printf("Router Name: %s\n", identity.Name)

// Set identity
err = sysRepo.Identity().SetIdentity(ctx, "MyRouter")
```

### Resources

```go
// Get system resources
resources, err := sysRepo.Resources().GetResources(ctx)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Uptime: %s\n", resources.Uptime)
fmt.Printf("Version: %s\n", resources.Version)
fmt.Printf("Free Memory: %d MB\n", resources.FreeMemory/1024/1024)
fmt.Printf("Total Memory: %d MB\n", resources.TotalMemory/1024/1024)
fmt.Printf("CPU Load: %d%%\n", resources.CpuLoad)
fmt.Printf("Platform: %s\n", resources.Platform)
fmt.Printf("Board: %s\n", resources.BoardName)
```

### Routerboard

```go
// Get routerboard info
rb, err := sysRepo.RouterBoard().GetRouterBoardInfo(ctx)
fmt.Printf("Model: %s\n", rb.Model)
fmt.Printf("Serial: %s\n", rb.SerialNumber)

// Get firmware version
current, upgrade, err := sysRepo.RouterBoard().GetFirmware(ctx)
fmt.Printf("Current: %s, Upgrade: %s\n", current, upgrade)
```

### Scripts

```go
// Get all scripts
scripts, err := sysRepo.Scripts().GetScripts(ctx)
for _, script := range scripts {
    fmt.Printf("Script: %s, Owner: %s\n", script.Name, script.Owner)
}

// Get by name
script, err := sysRepo.Scripts().GetScriptByName(ctx, "myscript")

// Add script
newScript := &domain.SystemScript{
    Name:    "backup-script",
    Owner:   "admin",
    Policy:  "read,write,policy,test",
    Source:  "/system backup save name=auto-backup",
    Comment: "Auto backup script",
}
id, err := sysRepo.Scripts().AddScript(ctx, newScript)

// Update
err = sysRepo.Scripts().UpdateScript(ctx, "*1", newScript)

// Remove
err = sysRepo.Scripts().RemoveScript(ctx, "*1")

// Run script
err = sysRepo.Scripts().RunScript(ctx, "backup-script")
```

### Scheduler

```go
// Get all schedulers
schedulers, err := sysRepo.Scheduler().GetSchedulers(ctx)
for _, s := range schedulers {
    fmt.Printf("Scheduler: %s, Interval: %s\n", s.Name, s.Interval)
}

// Get by name
sched, err := sysRepo.Scheduler().GetSchedulerByName(ctx, "auto-backup")

// Add scheduler
newSched := &domain.Scheduler{
    Name:      "daily-backup",
    StartTime: "00:00:00",
    Interval:  "1d",
    OnEvent:   "/system backup save name=daily",
    Comment:   "Daily backup",
}
id, err := sysRepo.Scheduler().AddScheduler(ctx, newSched)

// Update
err = sysRepo.Scheduler().UpdateScheduler(ctx, "*1", newSched)

// Remove
err = sysRepo.Scheduler().RemoveScheduler(ctx, "*1")

// Enable/Disable
err = sysRepo.Scheduler().EnableScheduler(ctx, "*1")
err = sysRepo.Scheduler().DisableScheduler(ctx, "*1")
```

---

## IP Address Repository

### Setup

```go
import "github.com/Butterfly-Student/go-ros/repository/ipaddress"

ipRepo := ipaddress.NewRepository(c)
```

### Address Management

```go
// Get all addresses
addrs, err := ipRepo.Address().GetAddresses(ctx)
for _, addr := range addrs {
    fmt.Printf("%s/%s on %s\n", addr.Address, addr.Network, addr.Interface)
}

// Get by interface
addrs, err = ipRepo.Address().GetAddressesByInterface(ctx, "ether1")

// Get by ID
addr, err := ipRepo.Address().GetAddressByID(ctx, "*1F")

// Add address
newAddr := &domain.IPAddress{
    Address:   "192.168.88.1/24",
    Network:   "192.168.88.0",
    Interface: "ether2",
    Comment:   "LAN",
}
id, err := ipRepo.Address().AddAddress(ctx, newAddr)

// Update
err = ipRepo.Address().UpdateAddress(ctx, "*1F", newAddr)

// Remove
err = ipRepo.Address().RemoveAddress(ctx, "*1F")

// Enable/Disable
err = ipRepo.Address().EnableAddress(ctx, "*1F")
err = ipRepo.Address().DisableAddress(ctx, "*1F")
```

### Pool Management

```go
// Get all pools
pools, err := ipRepo.Pool().GetPools(ctx)
for _, pool := range pools {
    fmt.Printf("Pool: %s, Ranges: %s\n", pool.Name, pool.Ranges)
}

// Get by name
pool, err := ipRepo.Pool().GetPoolByName(ctx, "hs-pool")

// Get names only
names, err := ipRepo.Pool().GetPoolNames(ctx)

// Add pool
newPool := &domain.IPPool{
    Name:     "new-pool",
    Ranges:   "192.168.10.10-192.168.10.254",
    NextPool: "none",
}
id, err := ipRepo.Pool().AddPool(ctx, newPool)

// Update
err = ipRepo.Pool().UpdatePool(ctx, "*1", newPool)

// Remove
err = ipRepo.Pool().RemovePool(ctx, "*1")

// Get used IPs
used, err := ipRepo.Pool().GetPoolUsed(ctx, "hs-pool")
```

---

## Firewall Repository

### Setup

```go
import "github.com/Butterfly-Student/go-ros/repository/firewall"

fwRepo := firewall.NewRepository(c)
```

### NAT Rules

```go
// Get all NAT rules
rules, err := fwRepo.NAT().GetNATRules(ctx)
for _, rule := range rules {
    fmt.Printf("Chain: %s, Action: %s\n", rule.Chain, rule.Action)
}
```

### Filter Rules

```go
// Get all filter rules
rules, err := fwRepo.Filter().GetRules(ctx)
for _, rule := range rules {
    fmt.Printf("Chain: %s, Action: %s\n", rule.Chain, rule.Action)
}
```

### Address Lists

```go
// Get all address lists
lists, err := fwRepo.AddressList().GetAddressLists(ctx)
for _, list := range lists {
    fmt.Printf("List: %s, Address: %s\n", list.List, list.Address)
}
```

---

## Queue Repository

### Setup

```go
import "github.com/Butterfly-Student/go-ros/repository/queue"

queueRepo := queue.NewRepository(c)
```

### Simple Queues

```go
// Get all simple queues
queues, err := queueRepo.Simple().GetSimpleQueues(ctx)
for _, q := range queues {
    fmt.Printf("Queue: %s, Target: %s\n", q.Name, q.Target)
}
```

### Queue Statistics

```go
// Get queue stats
stats, err := queueRepo.Queue().GetQueueStats(ctx, "*1F")
fmt.Printf("Rate: %d/%d bps\n", stats.RateIn, stats.RateOut)
```

---

## Monitor Repository

### Setup

```go
import "github.com/Butterfly-Student/go-ros/repository/monitor"

monRepo := monitor.NewRepository(c)
```

### System Monitoring

```go
// Get system resources
resources, err := monRepo.System().GetSystemResource(ctx)

// Get system health
health, err := monRepo.System().GetSystemHealth(ctx)
fmt.Printf("Voltage: %s, Temp: %s\n", health.Voltage, health.Temperature)

// Real-time monitoring
resultChan := make(chan *domain.SystemResourceMonitorStats)
stop, err := monRepo.System().StartSystemResourceMonitorListen(ctx, resultChan)
if err != nil {
    log.Fatal(err)
}
defer stop()

for stats := range resultChan {
    fmt.Printf("CPU: %.2f%%, Memory: %.2f%%\n", 
        stats.CPULoad, stats.MemoryUsage)
}
```

### Interface Monitoring

```go
// Get all interfaces
ifaces, err := monRepo.Interface().GetInterfaces(ctx)
for _, iface := range ifaces {
    fmt.Printf("Interface: %s, Type: %s\n", iface.Name, iface.Type)
}

// Real-time traffic monitoring
trafficChan := make(chan *domain.TrafficMonitorStats)
stop, err := monRepo.Interface().StartTrafficMonitorListen(ctx, "ether1", trafficChan)
if err != nil {
    log.Fatal(err)
}
defer stop()

for stats := range trafficChan {
    fmt.Printf("RX: %d bytes, TX: %d bytes\n", 
        stats.RXBytes, stats.TXBytes)
}
```

### Ping Monitoring

```go
// Start ping
pingChan := make(chan *domain.PingResult)
stop, err := monRepo.Ping().StartPingListen(ctx, "8.8.8.8", pingChan)
if err != nil {
    log.Fatal(err)
}
defer stop()

for result := range pingChan {
    fmt.Printf("Ping: %s - %s (%.2f ms)\n", 
        result.Host, result.Status, result.Time)
}
```

### Log Monitoring

```go
// Get logs
logs, err := monRepo.Log().GetLogs(ctx, 100)
for _, log := range logs {
    fmt.Printf("[%s] %s: %s\n", log.Time, log.Topics, log.Message)
}

// Get hotspot logs
logs, err = monRepo.Log().GetHotspotLogs(ctx, 50)

// Get PPP logs
logs, err = monRepo.Log().GetPPPLogs(ctx, 50)

// Real-time log streaming
logChan := make(chan *domain.LogEntry)
stop, err := monRepo.Log().ListenLogs(ctx, logChan)
if err != nil {
    log.Fatal(err)
}
defer stop()

for entry := range logChan {
    fmt.Printf("[%s] %s\n", entry.Topics, entry.Message)
}
```

---

## Mikhmon Repository

Lihat [Panduan Mikhmon](MIKHMON.md) untuk dokumentasi lengkap.

### Quick Reference

```go
import (
    mikhmonRepo "github.com/Butterfly-Student/go-ros/repository/mikhmon"
    mikhmonDomain "github.com/Butterfly-Student/go-ros/domain/mikhmon"
)

// Setup
hotspotRepo := hotspot.NewRepository(c)
generatorRepo := mikhmonRepo.NewGeneratorRepository()
voucherRepo := mikhmonRepo.NewVoucherRepository(c, hotspotRepo, generatorRepo)
profileRepo := mikhmonRepo.NewProfileRepository(hotspotRepo)
expireRepo := mikhmonRepo.NewExpireRepository(c, sysRepo)

// Generate vouchers
req := &mikhmonDomain.VoucherGenerateRequest{
    Quantity:   10,
    Profile:    "default",
    Mode:       mikhmonDomain.VoucherModeVoucher,
    NameLength: 6,
    CharSet:    mikhmonDomain.CharSetUpplow1,
    TimeLimit:  "1h",
}
batch, err := voucherRepo.GenerateBatch(ctx, req)

// Create profile
profileReq := &mikhmonDomain.ProfileRequest{
    Name:        "Paket-1Jam",
    AddressPool: "hs-pool",
    RateLimit:   "1M/2M",
    Config: mikhmonDomain.ProfileConfig{
        Price:      5000,
        Validity:   "1h",
        ExpireMode: mikhmonDomain.ExpireModeRemove,
    },
}
err = profileRepo.CreateProfile(ctx, profileReq)

// Setup expire monitor
err = expireRepo.SetupExpireMonitor(ctx)
```

---

## Best Practices

### 1. Context Management

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

users, err := repo.User().GetUsers(ctx, "")
```

### 2. Error Handling

```go
users, err := repo.User().GetUsers(ctx, "")
if err != nil {
    if errors.Is(err, context.DeadlineExceeded) {
        // Handle timeout
        return
    }
    log.Printf("Error: %v", err)
    return
}
```

### 3. Resource Cleanup

```go
defer client.Close()
// atau
defer manager.CloseAll()
```

### 4. Concurrent Operations

```go
var wg sync.WaitGroup

wg.Add(3)

go func() {
    defer wg.Done()
    users, _ := hotspotRepo.User().GetUsers(ctx, "")
    // process users
}()

go func() {
    defer wg.Done()
    profiles, _ := hotspotRepo.Profile().GetProfiles(ctx)
    // process profiles
}()

go func() {
    defer wg.Done()
    active, _ := hotspotRepo.Active().GetActive(ctx)
    // process active
}()

wg.Wait()
```

### 5. Proplist Optimization

```go
// Hanya request field yang diperlukan
const ProplistHotspotUserDefault = ".id,name,profile,disabled"
users, _ := repo.User().GetUsers(ctx, "", ProplistHotspotUserDefault)
```

---

**Selamat menggunakan repositories!** 🚀
