# MikMongo API - Hoppscotch Test Collections

## Setup

1. Import **environment.json** ke Hoppscotch: Environments > Import
2. Import collection JSON yang dibutuhkan: Collections > Import
3. Set environment variables setelah login (copy token dari response)

## Authentication Flow

1. Jalankan `POST /api/v1/auth/login` di `01-auth.json`
2. Copy `token` dari response ke environment variable `token`
3. Copy `refresh_token` ke `refreshToken`
4. Semua request admin akan otomatis menggunakan `<<token>>` via auth inheritance

## Collections

| File | Deskripsi | Auth |
|------|-----------|------|
| 01-auth.json | Login, refresh, logout, me | Public + Admin |
| 02-users.json | User CRUD | Admin |
| 03-customers.json | Customer CRUD + activate/deactivate | Admin |
| 04-routers.json | Router CRUD + sync + test-connection | Admin |
| 05-bandwidth-profiles.json | Bandwidth profile CRUD (router-scoped) | Admin |
| 06-subscriptions.json | Subscription CRUD + lifecycle | Admin |
| 07-registrations.json | Registration list + approve/reject | Public + Admin |
| 08-invoices.json | Invoice list + overdue + cancel | Admin |
| 09-payments.json | Payment CRUD + confirm/reject/refund | Admin |
| 10-mikrotik-ppp.json | PPP profiles, secrets, active | Admin |
| 11-mikrotik-hotspot.json | Hotspot profiles, users, active | Admin |
| 12-mikrotik-network.json | Queue, Firewall, IP | Admin |
| 13-mikrotik-monitor.json | System resource, interfaces | Admin |
| 14-mikrotik-raw.json | Raw RouterOS command | Admin |
| 15-mikhmon.json | Vouchers, profiles, reports, expire | Admin |
| 16-sales-agents.json | Sales agent CRUD + profile prices | Admin |
| 17-agent-invoices.json | Agent invoice list + pay/cancel | Admin |
| 18-hotspot-sales.json | Hotspot sales list | Admin |
| 19-cash-management.json | Cash entries + petty cash | Admin |
| 20-reports.json | Summary, subscriptions, cash flow | Admin |
| 21-settings.json | System settings CRUD | Admin |
| 22-webhooks.json | Midtrans + Xendit callbacks | Public |
| 23-customer-portal.json | Customer self-service portal | Portal |
| 24-agent-portal.json | Agent self-service portal | Agent |

## WebSocket / Realtime Endpoints

Hoppscotch WebSocket tidak bisa di-export ke collection JSON. Gunakan **Realtime tab** di Hoppscotch untuk testing.

### Cara Testing WebSocket

1. Buka Hoppscotch > tab **Realtime** > pilih **WebSocket**
2. Masukkan URL endpoint (ganti `{routerId}` dengan UUID router)
3. Auth via query parameter: tambahkan `?token=YOUR_JWT_TOKEN` di URL

### Endpoint List

| Endpoint | Params | Deskripsi |
|----------|--------|-----------|
| `ws://localhost:8080/api/v1/routers/{routerId}/ppp/ws/active` | `?token=JWT` | Stream PPP active connections |
| `ws://localhost:8080/api/v1/routers/{routerId}/ppp/ws/inactive` | `?token=JWT` | Stream PPP disconnections |
| `ws://localhost:8080/api/v1/routers/{routerId}/hotspot/ws/active` | `?token=JWT` | Stream hotspot active users |
| `ws://localhost:8080/api/v1/routers/{routerId}/hotspot/ws/inactive` | `?token=JWT` | Stream hotspot disconnections |
| `ws://localhost:8080/api/v1/routers/{routerId}/monitor/ws/system-resource` | `?token=JWT` | Stream CPU/memory/disk |
| `ws://localhost:8080/api/v1/routers/{routerId}/monitor/ws/traffic/{interfaceName}` | `?token=JWT` | Stream interface traffic |
| `ws://localhost:8080/api/v1/routers/{routerId}/monitor/ws/logs` | `?token=JWT&topics=hotspot` | Stream router logs |
| `ws://localhost:8080/api/v1/routers/{routerId}/monitor/ws/ping` | `?token=JWT&address=8.8.8.8` | Stream ping results |
| `ws://localhost:8080/api/v1/routers/{routerId}/raw/ws/listen` | `?token=JWT` | Raw command listener |

### Raw WS Listen

Setelah connect ke `/raw/ws/listen`, kirim message JSON:

```json
{"args": ["/interface/print"]}
```

### Monitor Logs Topics

Query param `topics` untuk filter log:
- `hotspot` - Hotspot events
- `ppp` - PPP events
- `system` - System events
- Custom topic sesuai RouterOS
