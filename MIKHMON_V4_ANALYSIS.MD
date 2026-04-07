# 📊 ANALISIS KOMPREHENSIF MIKHMON V4

> **Tanggal Analisis:** 6 April 2026  
> **Analyst:** Multiple AI Agents  
> **Project Path:** `/home/butterfly-student/Programing/Mikrotik/mikhmon_v4/`

---

## 📋 DAFTAR ISI

1. [Overview Project](#1-overview-project)
2. [Struktur Direktori](#2-struktur-direktori)
3. [Arsitektur dan Alur Kerja](#3-arsitektur-dan-alur-kerja)
4. [Menu dan Fitur Frontend](#4-menu-dan-fitur-frontend)
5. [Komunikasi dengan MikroTik](#5-komunikasi-dengan-mikrotik)
6. [Script dan Command RouterOS](#6-script-dan-command-routeros)
7. [Keamanan Sistem](#7-keamanan-sistem)
8. [Kesimpulan dan Rekomendasi](#8-kesimpulan-dan-rekomendasi)

---

## 1. OVERVIEW PROJECT

**Mikhmon V4** adalah aplikasi web-based Hotspot Monitor untuk MikroTik RouterOS yang dikembangkan oleh **Irhabi89**. Aplikasi ini merupakan recode untuk PHP 8.x dan sudah didecode agar bisa diupload ke shared hosting.

### Karakteristik Utama:
| Aspek | Detail |
|-------|--------|
| **Bahasa** | PHP 8.x (Native/Procedural) |
| **Framework** | Tidak menggunakan framework (Native PHP) |
| **Database** | File-based (PHP Arrays), tidak menggunakan SQL |
| **API** | RouterOS API (Port 8728/8729) |
| **Pattern** | Front Controller + Router |
| **Multi-Router** | ✅ Support multiple MikroTik router |
| **Responsive** | ✅ Mobile-friendly design |

### URL Structure:
```
/?{user}/{page}/{action}
```
Contoh: `/?admin/settings`, `/?mikrotik1/dashboard`

---

## 2. STRUKTUR DIREKTORI

```
mikhmon_v4/
├── index.php                    # Entry point utama
├── README.md                    # Dokumentasi singkat
├── robots.txt                   # SEO robots config
│
├── assets/                      # Static assets
│   ├── css/
│   │   ├── mikhmon-ui.{theme}.css   # Main UI styles (5 tema)
│   │   ├── pace.{theme}.css         # Loading indicator
│   │   └── font-awesome/            # Icons
│   ├── js/
│   │   ├── jquery.min.js            # jQuery 3.x
│   │   ├── mikhmon.js               # Main application logic
│   │   ├── func.js                  # Utility functions
│   │   ├── fancyTable.js            # Table sorting/filtering
│   │   ├── highcharts/              # Charting library
│   │   ├── codemirror/              # Code editor
│   │   └── qr/                      # QR code generator
│   ├── img/                     # Logo, favicon, icons
│   └── fonts/                   # Font Awesome
│
├── config/                      # Konfigurasi sistem
│   ├── config.php               # Data router & admin (encoded)
│   ├── connection.php           # Inisialisasi koneksi MikroTik
│   ├── page.php                 # Routing table
│   ├── readcfg.php              # Parser konfigurasi
│   ├── theme.php                # Tema aktif
│   └── settheme.php             # Handler ganti tema
│
├── core/                        # Core framework
│   ├── routeros_api.class.php   # RouterOS API Client
│   ├── route.php                # Router engine
│   ├── page_route.php           # Error handlers & decode functions
│   ├── jsencode.class.php       # XOR encoding untuk JSON
│   ├── generator_functions.php  # Random string generators
│   └── no_cache.php             # Cache control headers
│
├── get/                         # AJAX GET endpoints
│   ├── get_dashboard.php        # System info, traffic, logs
│   ├── get_users.php            # List hotspot users
│   ├── get_hotspot_active.php   # Active sessions
│   └── ... (18 files)
│
├── post/                        # AJAX POST endpoints
│   ├── post_add_user.php        # Add hotspot user
│   ├── post_generate_voucher.php # Generate voucher massal
│   ├── post_expire_monitor.php  # Setup scheduler monitoring
│   └── ... (12 files)
│
├── template/                    # Voucher templates
│   ├── header.{default,small,thermal}.txt
│   ├── row.{default,small,thermal}.txt
│   └── footer.{default,small,thermal}.txt
│
└── view/                        # View/UI files
    ├── header_html.php          # Template header
    ├── footer_html.php          # Template footer
    ├── menu.php                 # Sidebar navigation
    ├── login.php                # Login page
    ├── dashboard.php            # Dashboard utama
    ├── hotspot.php              # Manajemen Hotspot
    ├── admin/
    │   ├── settings.php         # Router settings
    │   └── editor.php           # Template editor
    └── ... (other pages)
```

---

## 3. ARSITEKTUR DAN ALUR KERJA

### 3.1 Alur Request

```
User Request
    ↓
index.php (Entry Point)
    ↓
Parse URI → $m_user, $page, $act
    ↓
Check Session (login?)
    ↓
route() in core/route.php
    ↓
Lookup $s_page array (config/page.php)
    ↓
Include target file (view/ or get/ or post/)
    ↓
If API needed:
    include config/connection.php
    Connect to MikroTik via RouterosAPI
    Execute commands
    ↓
Return HTML/JSON response
```

### 3.2 Key Components

#### A. Entry Point (`index.php`)
- Parse URI: `/?admin/settings` → `$m_user='admin'`, `$page='settings'`
- Session management untuk autentikasi
- Auto-routing ke halaman settings jika belum login
- Memuat semua session dari `config.php`

#### B. Router Engine (`core/route.php`)
```php
function route() {
    // Validasi dan include file berdasarkan URL
    // Support mobile detection
}
```

#### C. Configuration System (`config/`)
- **Format:** PHP arrays dengan encoding
- **Storage:** File-based (`config.php`)
- **Struktur:**
```php
$data['router_name'] = array(
    '1' => 'router_name!|!ip_address',
    '2' => 'router_name@|@api_username', 
    '3' => 'router_name#|#encoded_password',
    '4' => 'router_name%hotspot_name',
    '5' => 'router_name^dns_name',
    '6' => 'router_name&currency',
    // ... delimiter-based encoding
);
```

#### D. RouterosAPI Class (`core/routeros_api.class.php`)
| Property/Method | Deskripsi |
|-----------------|-----------|
| `$port` | Default 8728 (API port) |
| `$ssl` | SSL support (port 8729) |
| `$timeout` | 3 detik |
| `$attempts` | 5 percobaan koneksi |
| `$delay` | 3 detik antar percobaan |
| `connect($ip, $login, $password)` | Koneksi dengan retry mechanism |
| `comm($command, $array)` | Kirim command dengan parameter |
| `disconnect()` | Tutup koneksi |

---

## 4. MENU DAN FITUR FRONTEND

### 4.1 Menu Navigasi

#### **Menu User Biasa:**
| Menu | Icon | Deskripsi |
|------|------|-----------|
| Dashboard | `fa-th-large` | Overview sistem & monitoring |
| Hotspot | `fa-wifi` | Manajemen hotspot users |
| Log | `fa-file-text` | System & hotspot logs |
| Report | `fa-money` | Sales report & income |
| Settings | `fa-gear` | Router settings |

#### **Menu Admin:**
| Menu | Icon | Deskripsi |
|------|------|-----------|
| Settings | `fa-gear` | Admin & router configuration |
| Template Editor | `fa-edit` | Edit voucher templates |
| About | `fa-info-circle` | About Mikhmon |

### 4.2 Detail Fitur per Halaman

#### **🔥 Halaman Hotspot (`hotspot.php`)**

Tab/Sub-menu:

1. **Users** - Manajemen user hotspot
   - Tabel dengan filter & pagination
   - Tambah user manual
   - Generate voucher batch
   - Print voucher (Default/Small)
   - Filter by Profile & Comment
   
2. **User Profile** - Manajemen profil user
   - Add/Edit/Delete profile
   - Setting: Rate Limit, Shared Users, Expired Mode
   - Validity, Price, Selling Price
   - User Lock & Server Lock
   
3. **Active** - User yang sedang online
   - Real-time monitoring
   - Disconnect user
   - Filter by server
   
4. **Hosts** - Hotspot hosts list
   - Status: Authorized, Bypassed, DHCP, Dynamic
   - MAC Address management

#### **📈 Halaman Dashboard (`dashboard.php`)**

Widget/Card:
- **Connection Status** - Session selector & status
- **Hotspot Active** - Jumlah user online
- **Hotspot Users** - Total user terdaftar
- **Income Monitor** - Pendapatan hari ini & bulan ini
- **System Resource** - CPU, Memory, HDD usage (progress bar)
- **System Info** - Uptime, Board Name, Model, Router OS
- **Health** - Voltage & Temperature
- **Traffic Monitor** - Real-time graph (Highcharts)
- **App Log** - Application activity log
- **Hotspot Log** - Authentication log

#### **📊 Halaman Report (`report.php`)**

Fitur:
- Filter by Day, Month, Year
- Search/Filter tabel
- Export CSV & Excel
- Summary statistics (Qty & Total)
- Pagination dengan FancyTable

#### **⚙️ Halaman Settings (`admin/settings.php`)**

Tab:
1. **Admin** - Username & Password admin
2. **Router List** - Daftar router/session
3. **Edit Router** - Konfigurasi router:
   - Session Name, Router IP/Host
   - Username/Password
   - Hotspot Name, DNS Name
   - Currency, Theme
   - Idle Timeout, Live Report
   - Logo upload

#### **🎨 Template Editor (`admin/editor.php`)**

- CodeMirror editor dengan syntax highlighting
- Edit 3 template voucher: Default, Small, Thermal
- Edit 3 bagian per template: Header, Row, Footer
- Live preview (Ctrl+S untuk save)

**Template Variables:**
- `%username%`, `%password%`, `%profile%`
- `%validity%`, `%limitUptime%`, `%limitBytesTotal%`
- `%price%`, `%qrCode%`, `%hotspotName%`
- `%dnsName%`, `%currency%`, `%timeStamp%`
- `%#%` - Nomor urut

### 4.3 Library Frontend

| Library | Kegunaan |
|---------|----------|
| **jQuery** | DOM manipulation, AJAX |
| **Font Awesome** | Icons |
| **Mikhmon UI CSS** | Custom CSS framework (5 tema) |
| **Pace** | Loading progress bar |
| **Notify.js** | Toast notifications |
| **FancyTable** | Table sorting & pagination |
| **Highcharts** | Traffic monitoring charts |
| **CodeMirror** | Template editor syntax highlighting |
| **QRious** | QR Code generator |

**Tema yang Tersedia:** Light, Dark, Blue, Green, Pink

### 4.4 Responsive Design

- **Desktop:** Sidebar navigation, 12-column grid, Card-based layout
- **Mobile:** Top navbar, hamburger menu, Touch-friendly buttons
- **Auto-refresh:** 10 detik untuk dashboard

---

## 5. KOMUNIKASI DENGAN MIKROTIK

### 5.1 Metode Koneksi

Project menggunakan **MikroTik API Protocol** (bukan SSH atau Telnet):

**Default Configuration:**
- **Port:** 8728 (default API port)
- **SSL:** Nonaktif (dapat diaktifkan dengan port 8729)
- **Timeout:** 3 detik
- **Connection Attempts:** 5 kali
- **Delay between attempts:** 3 detik

### 5.2 Authentication Mechanism

Class mendukung **dua metode login**:
1. **Post-v6.43:** Login langsung dengan username/password
2. **Pre-v6.43:** Challenge-response dengan MD5 hash

```php
// Login flow
$this->write('/login', false);
$this->write('=name=' . $login, false);
$this->write('=password=' . $password);
```

### 5.3 Endpoint API MikroTik yang Digunakan

#### **System & Monitoring:**
| Endpoint | Fungsi |
|----------|--------|
| `/system/clock/print` | Get timezone dan system time |
| `/system/resource/print` | Get resource usage (CPU, RAM, disk) |
| `/system/routerboard/print` | Get routerboard info (model) |
| `/system/identity/print` | Get system identity |
| `/system/health/print` | Get health info (temperature, voltage) |
| `/log/print` | Get system logs |
| `/system/script/print` | Get script (untuk reporting) |
| `/system/logging/add` | Add logging configuration |

#### **Hotspot Management:**
| Endpoint | Fungsi |
|----------|--------|
| `/ip/hotspot/user/print` | Get hotspot users |
| `/ip/hotspot/user/add` | Add hotspot user |
| `/ip/hotspot/user/set` | Update hotspot user |
| `/ip/hotspot/user/remove` | Remove hotspot user |
| `/ip/hotspot/user/reset-counters` | Reset user counters |
| `/ip/hotspot/user/profile/print` | Get user profiles |
| `/ip/hotspot/user/profile/add` | Add user profile |
| `/ip/hotspot/user/profile/remove` | Remove user profile |
| `/ip/hotspot/active/print` | Get active sessions |
| `/ip/hotspot/active/remove` | Remove active session |
| `/ip/hotspot/host/print` | Get hotspot hosts |
| `/ip/hotspot/server/print` | Get hotspot servers |

#### **Interface & Traffic:**
| Endpoint | Fungsi |
|----------|--------|
| `/interface/print` | List interfaces |
| `/interface/monitor-traffic` | Monitor traffic TX/RX |

#### **Queue & Network:**
| Endpoint | Fungsi |
|----------|--------|
| `/queue/simple/print` | Get simple queues |
| `/ip/address/print` | Get IP addresses |
| `/ip/firewall/nat/print` | Get NAT rules |
| `/ip/pool/print` | Get address pools |

#### **Scheduler & Script:**
| Endpoint | Fungsi |
|----------|--------|
| `/system/scheduler/add` | Add scheduler |
| `/system/scheduler/remove` | Remove scheduler |
| `/system/scheduler/print` | List schedulers |
| `/system/script/add` | Add script |
| `/system/script/print` | List scripts |

### 5.4 Error Handling Mechanism

#### A. Connection Error Handling:
```php
for ($ATTEMPT = 1; $ATTEMPT <= $this->attempts; $ATTEMPT++) {
    $this->connected = false;
    // ... koneksi
    if ($this->socket) {
        if ($RESPONSE[0] == '!done') {
            $this->connected = true;
            break;
        }
    }
    sleep($this->delay);
}
```

#### B. API Error Handling:
```php
if(!empty($add['!trap'][0]['message'])){
    $message = $add['!trap'][0]['message'];
    $mess = array(
        "message" => "error",
        "data" => array("error" => $message)
    );
    echo json_encode($mess);
}
```

**Jenis Error yang Ditangani:**
- `!trap` - Error dari RouterOS
- `!fatal` - Fatal error dari RouterOS
- Connection timeout
- Invalid username/password

### 5.5 Response Encoding

```php
// jsEncode class untuk encoding response JSON
class jsEncode {
    public function encodeString($str, $decodeKey) {
        for ($i = 0; $i < strlen($str); $i++) {
            $a = $this->_getCharcode($str, $i);
            $b = $a ^ $decodeKey; // XOR dengan key (default 25)
            $result .= $this->_fromCharCode($b);
        }
        return $result;
    }
}
```

---

## 6. SCRIPT DAN COMMAND ROUTEROS

### 6.1 Total Command: 40+ Command

| Kategori | Command | Jumlah |
|----------|---------|--------|
| **Hotspot User** | `/ip/hotspot/user/*` | 5 |
| **Hotspot Profile** | `/ip/hotspot/user/profile/*` | 4 |
| **Hotspot Active** | `/ip/hotspot/active/*` | 2 |
| **Hotspot Host** | `/ip/hotspot/host/*` | 2 |
| **System Scheduler** | `/system/scheduler/*` | 4 |
| **System Script** | `/system/script/*` | 2 |
| **System Monitoring** | `/system/*` | 5 |
| **Logging** | `/system/logging/*`, `/log/*` | 3 |
| **Interface** | `/interface/*` | 2 |
| **Queue** | `/queue/simple/*` | 1 |
| **IP Pool** | `/ip/pool/*` | 1 |
| **Firewall** | `/ip/firewall/nat/*` | 1 |
| **Hotspot Server** | `/ip/hotspot/*` | 1 |

### 6.2 Fitur yang Diimplementasikan

1. **User Management** - Add/Edit/Delete user hotspot
2. **Voucher System** - Generate voucher massal dengan berbagai pola
3. **User Profile** - Profile dengan rate-limit, shared-users, parent-queue
4. **Auto Expiration** - 4 mode: rem, ntf, remc, ntfc
5. **MAC Binding** - User lock ke MAC address dan server lock
6. **Reporting System** - System script untuk menyimpan laporan login
7. **Expire Monitor** - Scheduler berjalan setiap 1 menit
8. **Real-time Monitoring** - Traffic, system resources, active users
9. **Voucher Printing** - Template dengan QR code

### 6.3 Script RouterOS yang Digenerate

#### **On-Login Script** (untuk auto expire):
```routeros
:put (",rem,15000,30d,18000,,Enable,Disable,");
:local mode "X";
{
  :local date [ /system clock get date ];
  :local year [ :pick $date 7 11 ];
  :local month [ :pick $date 0 3 ];
  :local day [ :pick $date 4 6 ];
  
  # Konversi bulan ke angka
  :if ($month = "jan") do={ :set month "01" }
  :if ($month = "feb") do={ :set month "02" }
  # ... dst
  
  :local date "$year/$month/$day";
  :local time [ /system clock get time ];
  
  # Buat scheduler untuk auto-remove
  /sys sch add name="$user" disable=no start-date=$date interval="30d";
  
  # Update comment user dengan expiry
  /ip hotspot user set comment="$exp X" [find where name="$user"];
}
```

#### **Expire Monitor Scheduler**:
```routeros
/system/scheduler/add \
  name="Mikhmon-Expire-Monitor" \
  interval=00:01:00 \
  on-event="...script monitoring..."
```

### 6.4 Mode Expire

| Mode | Deskripsi |
|------|-----------|
| **rem** | Remove - Hapus user saat expire |
| **ntf** | Notification - Hanya notifikasi |
| **remc** | Remove + Cache - Hapus dan simpan ke cache |
| **ntfc** | Notification + Cache - Notifikasi dan simpan ke cache |

### 6.5 Generator Functions

File: `core/generator_functions.php`

| Function | Pola |
|----------|------|
| `randN()` | Numeric (0-9) |
| `randLC()` | Lowercase (a-z) |
| `randUC()` | Uppercase (A-Z) |
| `randULC()` | Uppercase + Lowercase |
| `randNLC()` | Numeric + Lowercase |
| `randNUC()` | Numeric + Uppercase |
| `randNULC()` | Numeric + Uppercase + Lowercase |

---

## 7. KEAMANAN SISTEM

### 7.1 Encoding & Enkripsi

#### **Password Encoding:**
```php
// Fungsi enc_rypt dan dec_rypt di routeros_api.class.php
function enc_rypt($string, $key = 128) {
    $result = '';
    for ($i = 0, $k = strlen($string); $i < $k; $i++) {
        $char = substr($string, $i, 1);
        $keychar = substr($key, ($i % strlen($key)) - 1, 1);
        $char = chr(ord($char) + ord($keychar));
        $result .= $char;
    }
    return base64_encode($result);
}
```

#### **Decode Function:**
```php
function decode($text) {
    // Base64 + XOR + Double Base64
}
```

### 7.2 Session Management

- Menggunakan PHP `$_SESSION`
- Session key: `mikhmon` untuk autentikasi admin
- Idle timeout: 10 menit (default)

### 7.3 File Protection

- `index.php` check untuk mencegah direct access
- `error_reporting(0)` di production

### 7.4 Input Handling

- Minimal sanitization
- Rely on MikroTik API untuk validasi
- JSON encoding dengan XOR untuk response

---

## 8. KESIMPULAN DAN REKOMENDASI

### 8.1 Kelebihan Mikhmon V4

✅ **Ringan dan Sederhana**
- Tidak memerlukan database SQL
- Bisa dijalankan di shared hosting
- File-based configuration

✅ **Fitur Lengkap**
- Multi-router support
- Voucher system dengan QR code
- Auto-expire dengan scheduler
- Real-time monitoring
- Multiple themes

✅ **Mobile Friendly**
- Responsive design
- Touch-friendly interface
- Mobile-optimized navigation

✅ **Robust Connection**
- Retry mechanism (5 attempts)
- Error handling untuk `!trap`
- Session-based caching

### 8.2 Area Perbaikan

⚠️ **Keamanan**
- Password encryption sederhana (XOR + Base64)
- Tidak ada CSRF protection
- Minimal input sanitization
- **Rekomendasi:** Gunakan bcrypt untuk password, tambahkan CSRF token

⚠️ **Performa**
- Tidak ada connection pooling
- Query ke MikroTik synchronous
- **Rekomendasi:** Implementasi async atau queue untuk operasi berat

⚠️ **Skalabilitas**
- File-based storage tidak scalable untuk data besar
- **Rekomendasi:** Gunakan database (MySQL/PostgreSQL) untuk production scale

⚠️ **Error Handling**
- Error handling bisa lebih granular
- Tidak ada retry untuk query panjang
- **Rekomendasi:** Implementasi exponential backoff

### 8.3 Use Case Terbaik

Mikhmon V4 **cocok untuk:**
- 🏢 Small to medium hotspot business
- 🏠 Home/SOHO environment
- 📡 Single or few MikroTik routers
- 💻 Shared hosting deployment

**Tidak direkomendasikan untuk:**
- 🏭 Enterprise dengan ribuan user
- 📊 High-throughput environment
- 🔒 Environment dengan kebutuhan keamanan tinggi

### 8.4 Checklist Deployment

- [ ] RouterOS API enabled (`/ip service enable api`)
- [ ] Port 8728 terbuka (firewall)
- [ ] PHP 8.x terinstall
- [ ] Write permission untuk folder config/
- [ ] SSL/TLS untuk production (HTTPS)
- [ ] Backup config.php secara berkala

---

## 📎 LAMPIRAN

### A. File-file Penting

| File | Fungsi Kunci |
|------|--------------|
| `core/routeros_api.class.php` | Class API untuk komunikasi MikroTik |
| `core/generator_functions.php` | Fungsi random string untuk voucher |
| `post/post_generate_voucher.php` | Generate voucher massal |
| `post/post_add_userprofile.php` | Profile dengan on-login script |
| `post/post_expire_monitor.php` | Setup scheduler monitoring |
| `template/row.default.txt` | Template voucher printing |

### B. Dependensi Eksternal

- MikroTik RouterOS dengan API enabled
- PHP 8.x dengan ekstensi: sockets, json
- Web server: Apache/Nginx

---

**📅 Dokumen ini dibuat:** 6 April 2026  
**🤖 Dibuat dengan:** Multiple AI Agent Analysis  
**📂 Lokasi:** `/home/butterfly-student/Programing/Mikrotik/mikhmon_v4/ANALISIS.md`

---
*End of Analysis*
