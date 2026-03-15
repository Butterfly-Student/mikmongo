# Referensi API Lengkap go-ros

Dokumen ini berisi referensi lengkap untuk semua domain models, repository interfaces, dan methods yang tersedia dalam go-ros.

## Table of Contents

1. [Domain Models](#domain-models)
2. [Repository Interfaces](#repository-interfaces)
3. [Client & Manager](#client--manager)
4. [Utility Functions](#utility-functions)

---

## Domain Models

### Hotspot Domain (`domain/hotspot.go`)

#### HotspotUser

Representasi user hotspot di MikroTik.

```go
type HotspotUser struct {
    ID              string `json:".id,omitempty"`
    Server          string `json:"server,omitempty"`
    Name            string `json:"name"`
    Password        string `json:"password,omitempty"`
    Profile         string `json:"profile,omitempty"`
    MACAddress      string `json:"mac-address,omitempty"`
    IPAddress       string `json:"address,omitempty"`
    Routes          string `json:"routes,omitempty"`
    Uptime          string `json:"uptime,omitempty"`
    BytesIn         int64  `json:"bytes-in,omitempty"`
    BytesOut        int64  `json:"bytes-out,omitempty"`
    PacketsIn       int64  `json:"packets-in,omitempty"`
    PacketsOut      int64  `json:"packets-out,omitempty"`
    LimitUptime     string `json:"limit-uptime,omitempty"`
    LimitBytesIn    int64  `json:"limit-bytes-in,omitempty"`
    LimitBytesOut   int64  `json:"limit-bytes-out,omitempty"`
    LimitBytesTotal int64  `json:"limit-bytes-total,omitempty"`
    Comment         string `json:"comment,omitempty"`
    Disabled        bool   `json:"disabled,omitempty"`
    Email           string `json:"email,omitempty"`
}
```

**Field Descriptions:**

| Field | Type | Description |
|-------|------|-------------|
| ID | string | Unique identifier dari MikroTik (misal: "*1F") |
| Server | string | Nama hotspot server |
| Name | string | Username untuk login |
| Password | string | Password user |
| Profile | string | Nama profile yang digunakan |
| MACAddress | string | MAC address yang ter-binding |
| IPAddress | string | IP address yang ter-binding |
| Routes | string | Static routes untuk user |
| Uptime | string | Total waktu online |
| BytesIn | int64 | Total bytes received |
| BytesOut | int64 | Total bytes sent |
| LimitUptime | string | Batas waktu online (format: 1d2h3m) |
| LimitBytesTotal | int64 | Batas total bytes |
| Comment | string | Komentar/catatan |
| Disabled | bool | Status disable/enable |

#### UserProfile

Profile hotspot yang menentukan karakteristik user.

```go
type UserProfile struct {
    ID                 string `json:".id,omitempty"`
    Name               string `json:"name"`
    AddressPool        string `json:"address-pool,omitempty"`
    SharedUsers        int    `json:"shared-users,omitempty"`
    RateLimit          string `json:"rate-limit,omitempty"`
    ParentQueue        string `json:"parent-queue,omitempty"`
    QueueType          string `json:"queue-type,omitempty"`
    StatusAutorefresh  string `json:"status-autorefresh,omitempty"`
    OnLogin            string `json:"on-login,omitempty"`
    OnLogout           string `json:"on-logout,omitempty"`
    OnUp               string `json:"on-up,omitempty"`
    OpenStatusPage     string `json:"open-status-page,omitempty"`
    TransparentProxy   bool   `json:"transparent-proxy,omitempty"`
    Advertise          bool   `json:"advertise,omitempty"`
    Disabled           bool   `json:"disabled,omitempty"`
}
```

**Field Descriptions:**

| Field | Type | Description |
|-------|------|-------------|
| Name | string | Nama profile |
| AddressPool | string | IP pool untuk user |
| SharedUsers | int | Jumlah user yang bisa sharing |
| RateLimit | string | Limit bandwidth (format: 1M/2M) |
| ParentQueue | string | Parent queue untuk queue tree |
| OnLogin | string | Script yang dijalankan saat login |
| OnLogout | string | Script yang dijalankan saat logout |

#### HotspotActive

User yang sedang aktif/login di hotspot.

```go
type HotspotActive struct {
    ID              string `json:".id,omitempty"`
    Server          string `json:"server,omitempty"`
    User            string `json:"user,omitempty"`
    Address         string `json:"address,omitempty"`
    MACAddress      string `json:"mac-address,omitempty"`
    LoginBy         string `json:"login-by,omitempty"`
    Uptime          string `json:"uptime,omitempty"`
    SessionTimeLeft string `json:"session-time-left,omitempty"`
    IdleTime        string `json:"idle-time,omitempty"`
    IdleTimeout     string `json:"idle-timeout,omitempty"`
    KeepaliveTimeout string `json:"keepalive-timeout,omitempty"`
    BytesIn         int64  `json:"bytes-in,omitempty"`
    BytesOut        int64  `json:"bytes-out,omitempty"`
    PacketsIn       int64  `json:"packets-in,omitempty"`
    PacketsOut      int64  `json:"packets-out,omitempty"`
}
```

### PPP Domain (`domain/ppp.go`)

#### PPPSecret

User/secret untuk PPP (VPN) connections.

```go
type PPPSecret struct {
    ID              string `json:".id,omitempty"`
    Name            string `json:"name"`
    Password        string `json:"password,omitempty"`
    Profile         string `json:"profile,omitempty"`
    Service         string `json:"service,omitempty"`
    CallerID        string `json:"caller-id,omitempty"`
    RemoteAddress   string `json:"remote-address,omitempty"`
    Routes          string `json:"routes,omitempty"`
    LimitBytesIn    int64  `json:"limit-bytes-in,omitempty"`
    LimitBytesOut   int64  `json:"limit-bytes-out,omitempty"`
    Comment         string `json:"comment,omitempty"`
    Disabled        bool   `json:"disabled,omitempty"`
}
```

#### PPPProfile

Profile untuk PPP connections.

```go
type PPPProfile struct {
    ID               string `json:".id,omitempty"`
    Name             string `json:"name"`
    LocalAddress     string `json:"local-address,omitempty"`
    RemoteAddress    string `json:"remote-address,omitempty"`
    RateLimit        string `json:"rate-limit,omitempty"`
    UseCompression   bool   `json:"use-compression,omitempty"`
    UseEncryption    bool   `json:"use-encryption,omitempty"`
    UseMPLS          bool   `json:"use-mpls,omitempty"`
}
```

### System Domain (`domain/system.go`)

#### SystemResource

Informasi resource sistem MikroTik.

```go
type SystemResource struct {
    Uptime        string `json:"uptime,omitempty"`
    Version       string `json:"version,omitempty"`
    BuildTime     string `json:"build-time,omitempty"`
    FreeMemory    int64  `json:"free-memory,omitempty"`
    TotalMemory   int64  `json:"total-memory,omitempty"`
    FreeHddSpace  int64  `json:"free-hdd-space,omitempty"`
    TotalHddSpace int64  `json:"total-hdd-space,omitempty"`
    CpuCount      int    `json:"cpu-count,omitempty"`
    CpuFrequency  int64  `json:"cpu-frequency,omitempty"`
    CpuLoad       int    `json:"cpu-load,omitempty"`
    Platform      string `json:"platform,omitempty"`
    BoardName     string `json:"board-name,omitempty"`
}
```

#### SystemScript

Script yang tersimpan di MikroTik.

```go
type SystemScript struct {
    ID                     string `json:".id,omitempty"`
    Name                   string `json:"name"`
    Owner                  string `json:"owner,omitempty"`
    Policy                 string `json:"policy,omitempty"`
    Source                 string `json:"source,omitempty"`
    Comment                string `json:"comment,omitempty"`
    DontRequirePermissions bool   `json:"dont-require-permissions,omitempty"`
    RunCount               string `json:"run-count,omitempty"`
}
```

#### Scheduler

Task scheduler di MikroTik.

```go
type Scheduler struct {
    ID        string `json:".id,omitempty"`
    Name      string `json:"name"`
    StartTime string `json:"start-time,omitempty"`
    Interval  string `json:"interval,omitempty"`
    OnEvent   string `json:"on-event,omitempty"`
    Comment   string `json:"comment,omitempty"`
    Disabled  bool   `json:"disabled,omitempty"`
    RunCount  int    `json:"run-count,omitempty"`
    NextRun   string `json:"next-run,omitempty"`
}
```

### Mikhmon Domain (`domain/mikhmon/`)

#### VoucherGenerateRequest

Request untuk generate voucher.

```go
type VoucherGenerateRequest struct {
    Quantity   int    `json:"quantity" validate:"required,min=1,max=1000"`
    Server     string `json:"server,omitempty"`
    Profile    string `json:"profile" validate:"required"`
    Mode       string `json:"mode" validate:"required,oneof=vc up"`
    NameLength int    `json:"nameLength" validate:"min=3,max=12"`
    Prefix     string `json:"prefix,omitempty"`
    CharSet    string `json:"charSet" validate:"required"`
    TimeLimit  string `json:"timeLimit,omitempty"`
    DataLimit  string `json:"dataLimit,omitempty"`
    Comment    string `json:"comment,omitempty"`
}
```

**Mode Values:**
- `"vc"` - Voucher Card (username = password)
- `"up"` - User/Password (username != password)

**CharSet Values:**
- `"lower"` - a-z
- `"upper"` - A-Z  
- `"upplow"` - a-zA-Z
- `"lower1"` - a-z0-9
- `"upper1"` - A-Z0-9
- `"upplow1"` - a-zA-Z0-9
- `"mix"` - a-zA-Z0-9
- `"mix1"` - a-zA-Z0-9!@#$%
- `"mix2"` - a-zA-Z0-9!@#$%^&*
- `"num"` - 0-9

#### VoucherBatch

Hasil generate batch voucher.

```go
type VoucherBatch struct {
    Code      string    `json:"code"`
    Quantity  int       `json:"quantity"`
    Profile   string    `json:"profile"`
    Server    string    `json:"server"`
    TimeLimit string    `json:"timeLimit,omitempty"`
    DataLimit string    `json:"dataLimit,omitempty"`
    Vouchers  []Voucher `json:"vouchers"`
}
```

#### ProfileConfig

Konfigurasi profile Mikhmon.

```go
type ProfileConfig struct {
    Name          string `json:"name" validate:"required"`
    AddressPool   string `json:"addressPool,omitempty"`
    RateLimit     string `json:"rateLimit,omitempty"`
    SharedUsers   int    `json:"sharedUsers,omitempty"`
    ParentQueue   string `json:"parentQueue,omitempty"`
    Price         int64  `json:"price,omitempty"`
    SellingPrice  int64  `json:"sellingPrice,omitempty"`
    Validity      string `json:"validity,omitempty"`
    ExpireMode    string `json:"expireMode,omitempty"`
    LockUser      bool   `json:"lockUser,omitempty"`
    LockServer    bool   `json:"lockServer,omitempty"`
    OnLoginScript string `json:"onLoginScript,omitempty"`
}
```

**ExpireMode Values:**
- `"rem"` - Remove user saat expired
- `"ntf"` - Disable user saat expired (limit-uptime=1s)
- `"remc"` - Remove + record ke report
- `"ntfc"` - Disable + record ke report
- `"0"` - Tidak ada expiration

#### SalesReport

Report penjualan yang disimpan di /system/script.

```go
type SalesReport struct {
    ID       string `json:"id,omitempty"`
    Name     string `json:"name"`
    Date     string `json:"date"`
    Time     string `json:"time"`
    User     string `json:"user"`
    Price    int64  `json:"price"`
    IP       string `json:"ip"`
    MAC      string `json:"mac"`
    Validity string `json:"validity"`
    Profile  string `json:"profile"`
    Comment  string `json:"comment"`
    Owner    string `json:"owner"`
    Source   string `json:"source"`
}
```

---

## Repository Interfaces

### Hotspot Repository (`repository/hotspot/`)

#### UserRepository

```go
type UserRepository interface {
    // Read operations
    GetUsers(ctx context.Context, profile string, proplist ...string) ([]*domain.HotspotUser, error)
    GetUsersByComment(ctx context.Context, comment string, proplist ...string) ([]*domain.HotspotUser, error)
    GetUserByID(ctx context.Context, id string, proplist ...string) (*domain.HotspotUser, error)
    GetUserByName(ctx context.Context, name string, proplist ...string) (*domain.HotspotUser, error)
    GetUsersCount(ctx context.Context) (int, error)
    
    // Write operations
    AddUser(ctx context.Context, user *domain.HotspotUser) (string, error)
    UpdateUser(ctx context.Context, id string, user *domain.HotspotUser) error
    RemoveUser(ctx context.Context, id string) error
    RemoveUsersByComment(ctx context.Context, comment string) error
    RemoveUsers(ctx context.Context, ids []string) error
    
    // Status operations
    DisableUser(ctx context.Context, id string) error
    EnableUser(ctx context.Context, id string) error
    DisableUsers(ctx context.Context, ids []string) error
    EnableUsers(ctx context.Context, ids []string) error
    
    // Counter operations
    ResetUserCounters(ctx context.Context, id string) error
    ResetUserCountersMultiple(ctx context.Context, ids []string) error
    
    // Raw operations
    PrintUsersRaw(ctx context.Context, args ...string) ([]map[string]string, error)
    ListenUsersRaw(ctx context.Context, args []string, resultChan chan<- map[string]string) (func() error, error)
}
```

**Usage Example:**

```go
hotspotRepo := hotspot.NewRepository(client)

// Get all users
users, err := hotspotRepo.User().GetUsers(ctx, "")

// Get users by profile
users, err := hotspotRepo.User().GetUsers(ctx, "default")

// Add new user
user := &domain.HotspotUser{
    Name:     "testuser",
    Password: "testpass",
    Profile:  "default",
}
id, err := hotspotRepo.User().AddUser(ctx, user)

// Update user
user.Password = "newpassword"
err = hotspotRepo.User().UpdateUser(ctx, id, user)

// Remove user
err = hotspotRepo.User().RemoveUser(ctx, id)

// Disable/Enable
err = hotspotRepo.User().DisableUser(ctx, id)
err = hotspotRepo.User().EnableUser(ctx, id)

// Batch operations
err = hotspotRepo.User().RemoveUsers(ctx, []string{"*1", "*2", "*3"})
```

#### ProfileRepository

```go
type ProfileRepository interface {
    GetProfiles(ctx context.Context, proplist ...string) ([]*domain.UserProfile, error)
    GetProfileByID(ctx context.Context, id string, proplist ...string) (*domain.UserProfile, error)
    GetProfileByName(ctx context.Context, name string, proplist ...string) (*domain.UserProfile, error)
    AddProfile(ctx context.Context, profile *domain.UserProfile) (string, error)
    UpdateProfile(ctx context.Context, id string, profile *domain.UserProfile) error
    RemoveProfile(ctx context.Context, id string) error
    DisableProfile(ctx context.Context, id string) error
    EnableProfile(ctx context.Context, id string) error
    RemoveProfiles(ctx context.Context, ids []string) error
    DisableProfiles(ctx context.Context, ids []string) error
    EnableProfiles(ctx context.Context, ids []string) error
}
```

#### ActiveRepository

```go
type ActiveRepository interface {
    GetActive(ctx context.Context) ([]*domain.HotspotActive, error)
    GetActiveCount(ctx context.Context) (int, error)
    RemoveActive(ctx context.Context, id string) error
    RemoveActives(ctx context.Context, ids []string) error
    ListenActive(ctx context.Context, resultChan chan<- []*domain.HotspotActive) (func() error, error)
    ListenInactive(ctx context.Context, resultChan chan<- []*domain.HotspotUser) (func() error, error)
}
```

### PPP Repository (`repository/ppp/`)

#### SecretRepository

```go
type SecretRepository interface {
    GetSecrets(ctx context.Context, profile string, proplist ...string) ([]*domain.PPPSecret, error)
    GetSecretByID(ctx context.Context, id string, proplist ...string) (*domain.PPPSecret, error)
    GetSecretByName(ctx context.Context, name string, proplist ...string) (*domain.PPPSecret, error)
    AddSecret(ctx context.Context, secret *domain.PPPSecret) error
    UpdateSecret(ctx context.Context, id string, secret *domain.PPPSecret) error
    RemoveSecret(ctx context.Context, id string) error
    DisableSecret(ctx context.Context, id string) error
    EnableSecret(ctx context.Context, id string) error
    RemoveSecrets(ctx context.Context, ids []string) error
    DisableSecrets(ctx context.Context, ids []string) error
    EnableSecrets(ctx context.Context, ids []string) error
    PrintSecretsRaw(ctx context.Context, args ...string) ([]map[string]string, error)
    ListenSecretsRaw(ctx context.Context, args []string, resultChan chan<- map[string]string) (func() error, error)
}
```

### System Repository (`repository/system/`)

#### ScriptsRepository

```go
type ScriptsRepository interface {
    GetScripts(ctx context.Context) ([]*domain.SystemScript, error)
    GetScriptByName(ctx context.Context, name string) (*domain.SystemScript, error)
    AddScript(ctx context.Context, script *domain.SystemScript) (string, error)
    UpdateScript(ctx context.Context, id string, script *domain.SystemScript) error
    RemoveScript(ctx context.Context, id string) error
    RunScript(ctx context.Context, name string) error
}
```

#### SchedulerRepository

```go
type SchedulerRepository interface {
    GetSchedulers(ctx context.Context) ([]*domain.Scheduler, error)
    GetSchedulerByName(ctx context.Context, name string) (*domain.Scheduler, error)
    AddScheduler(ctx context.Context, scheduler *domain.Scheduler) (string, error)
    UpdateScheduler(ctx context.Context, id string, scheduler *domain.Scheduler) error
    RemoveScheduler(ctx context.Context, id string) error
    EnableScheduler(ctx context.Context, id string) error
    DisableScheduler(ctx context.Context, id string) error
}
```

### Mikhmon Repository (`repository/mikhmon/`)

#### VoucherRepository

```go
type VoucherRepository interface {
    GenerateBatch(ctx context.Context, req *mikhmon.VoucherGenerateRequest) (*mikhmon.VoucherBatch, error)
    GetVouchersByComment(ctx context.Context, comment string) ([]*mikhmon.Voucher, error)
    GetVouchersByCode(ctx context.Context, code string) ([]*mikhmon.Voucher, error)
    RemoveVoucherBatch(ctx context.Context, comment string) error
}
```

**Usage Example:**

```go
hotspotRepo := hotspot.NewRepository(client)
generatorRepo := mikhmonRepo.NewGeneratorRepository()
voucherRepo := mikhmonRepo.NewVoucherRepository(client, hotspotRepo, generatorRepo)

// Generate vouchers
req := &mikhmon.VoucherGenerateRequest{
    Quantity:   10,
    Profile:    "default",
    Mode:       mikhmon.VoucherModeVoucher,
    NameLength: 6,
    CharSet:    mikhmon.CharSetUpplow1,
    TimeLimit:  "1h",
}

batch, err := voucherRepo.GenerateBatch(ctx, req)
if err != nil {
    log.Fatal(err)
}

for _, v := range batch.Vouchers {
    fmt.Printf("Voucher: %s / %s\n", v.Name, v.Password)
}
```

#### ProfileRepository

```go
type ProfileRepository interface {
    CreateProfile(ctx context.Context, req *mikhmon.ProfileRequest) error
    UpdateProfile(ctx context.Context, id string, req *mikhmon.ProfileRequest) error
    GenerateOnLoginScript(data *mikhmon.OnLoginScriptData) string
}
```

#### ExpireRepository

```go
type ExpireRepository interface {
    SetupExpireMonitor(ctx context.Context) error
    DisableExpireMonitor(ctx context.Context) error
    IsExpireMonitorEnabled(ctx context.Context) (bool, error)
    GenerateExpireMonitorScript() string
}
```

---

## Client & Manager

### Client

```go
type Client struct {
    // ... internal fields
}

// Constructor
func New(cfg Config) (*Client, error)
func NewClient(cfg Config, logger *zap.Logger) *Client

// Connection
func (c *Client) Connect(ctx context.Context) error
func (c *Client) Close()
func (c *Client) IsConnected() bool
func (c *Client) IsAsync() bool

// Command execution
func (c *Client) Run(sentence ...string) (*routeros.Reply, error)
func (c *Client) RunContext(ctx context.Context, sentence ...string) (*routeros.Reply, error)
func (c *Client) RunMany(ctx context.Context, commands [][]string) ([]*routeros.Reply, []error)

// Streaming
func (c *Client) ListenArgs(args []string) (*routeros.ListenReply, error)
func (c *Client) ListenArgsContext(ctx context.Context, args []string) (*routeros.ListenReply, error)
```

### Manager

```go
type Manager struct {
    // ... internal fields
}

// Constructor
func NewManager(logger *zap.Logger) *Manager

// Router management
func (m *Manager) Register(ctx context.Context, name string, cfg Config) error
func (m *Manager) Get(name string) (*Client, error)
func (m *Manager) GetOrConnect(ctx context.Context, name string, cfg Config) (*Client, error)
func (m *Manager) MustGet(name string) *Client
func (m *Manager) Unregister(name string)
func (m *Manager) Names() []string
func (m *Manager) TestConnection(ctx context.Context, cfg Config) error
func (m *Manager) CloseAll()
```

### Config

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

---

## Utility Functions

### Parse Functions (`client/helpers.go`)

```go
// Parse integer dari string RouterOS
func parseInt(s string) int64

// Parse float dari string RouterOS
func parseFloat(s string) float64

// Parse bool dari string RouterOS ("yes"/"true")
func parseBool(s string) bool

// Parse rate dengan unit (bps, kbps, Mbps, Gbps)
// Contoh: "74.3kbps" -> 74300
func ParseRate(s string) int64

// Split value dengan slash
// Contoh: "11519824317/96401664078" -> (11519824317, 96401664078)
func SplitSlashValue(value string) (int64, int64)

// Parse byte size (KiB, MiB, GiB, dll)
// Contoh: "6.4MiB" -> 6710886
func ParseByteSize(s string) int64
```

### Batch Listener (`client/helpers.go`)

```go
// ListenBatches membaca sentences dan mengelompokkan menjadi batches
// Berguna untuk monitoring dengan interval
func ListenBatches(
    ctx context.Context,
    sentences <-chan *proto.Sentence,
    debounce time.Duration,
) <-chan []*proto.Sentence
```

---

## Constants

### Proplist Constants

```go
// Hotspot
const ProplistHotspotUserDefault = ".id,name,profile,disabled"
const ProplistHotspotProfileDefault = ".id,name,address-pool,shared-users,rate-limit"
const ProplistHotspotActiveDefault = ".id,user,address,mac-address,uptime"

// PPP
const ProplistPPPSecretDefault = ".id,name,profile,service,disabled"
const ProplistPPPProfileDefault = ".id,name,local-address,remote-address"
```

### CharSet Constants (Mikhmon)

```go
const (
    CharSetLower   = "lower"   // a-z
    CharSetUpper   = "upper"   // A-Z
    CharSetUpplow  = "upplow"  // a-zA-Z
    CharSetLower1  = "lower1"  // a-z0-9
    CharSetUpper1  = "upper1"  // A-Z0-9
    CharSetUpplow1 = "upplow1" // a-zA-Z0-9
    CharSetMix     = "mix"     // a-zA-Z0-9
    CharSetMix1    = "mix1"    // a-zA-Z0-9!@#$%
    CharSetMix2    = "mix2"    // a-zA-Z0-9!@#$%^&*
    CharSetNumeric = "num"     // 0-9
)
```

### Expire Mode Constants (Mikhmon)

```go
const (
    ExpireModeRemove       = "rem"  // Remove user on expire
    ExpireModeNotify       = "ntf"  // Disable user on expire
    ExpireModeRemoveRecord = "remc" // Remove + record
    ExpireModeNotifyRecord = "ntfc" // Disable + record
    ExpireModeNoExpire     = "0"    // No expiration
)
```

---

## Error Handling

Semua repository methods mengembalikan error dengan konteks:

```go
// Contoh error messages:
"failed to get users: from RouterOS device: ..."
"failed to add user: from RouterOS device: ..."
"connect mikrotik 192.168.88.1: ..."
```

Gunakan `errors.Is()` atau `errors.As()` untuk error checking:

```go
users, err := repo.User().GetUsers(ctx, "")
if err != nil {
    if errors.Is(err, context.DeadlineExceeded) {
        // Handle timeout
    }
    if errors.Is(err, context.Canceled) {
        // Handle cancellation
    }
    // Handle other errors
}
```

---

## Best Practices

1. **Selalu gunakan context dengan timeout:**
```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
```

2. **Reuse client/repository instances:**
```go
// Jangan buat client baru untuk setiap operasi
client, _ := client.New(cfg)
repo := hotspot.NewRepository(client)

// Gunakan repo berkali-kali
users, _ := repo.User().GetUsers(ctx, "")
profiles, _ := repo.Profile().GetProfiles(ctx)
```

3. **Handle errors dengan baik:**
```go
users, err := repo.User().GetUsers(ctx, "")
if err != nil {
    log.Printf("Failed to get users: %v", err)
    return err
}
```

4. **Gunakan proplist untuk optimasi:**
```go
// Hanya request field yang diperlukan
users, _ := repo.User().GetUsers(ctx, "", ".id,name,profile")
```

5. **Close resources:**
```go
defer client.Close()
// atau
defer manager.CloseAll()
```
