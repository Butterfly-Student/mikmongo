# Roadmap: MikMongo Dashboard

**Created:** 2026-03-30
**Granularity:** Coarse (5 phases)
**Execution:** Parallel where possible

## Phase 1 — Foundation & Auth
**Goal:** Project scaffold berjalan di `dashboard/`, semua 3 portal dapat login dan mengakses halaman yang terproteksi.

**Deliverables:**
- Vite + React 19 + TypeScript project di `dashboard/`
- TanStack Router v1 file-based routing dengan 3 route trees: `/_admin`, `/_agent`, `/_customer`
- Auth flow: login, JWT storage, auto-refresh, logout untuk ketiga portal
- Zustand auth store dengan RBAC (superadmin/admin/teknisi)
- Axios API client dengan interceptors (auth header, refresh token, error handling)
- Shared AppShell layout: Sidebar, Topbar, mobile nav (hamburger)
- Dark/light mode dengan system preference + toggle
- Shadcn/UI + Tailwind CSS v4 design tokens (warna, spacing, typography)
- Loading states, error boundaries, toast notifications (Sonner)
- Admin dashboard overview page (summary cards — data dari API `/api/v1/reports/summary`)

**Requirements covered:** SETUP-01~10, AUTH-01~08, ADMIN-01~03

**Success criteria:**
- `npm run dev` berjalan di `dashboard/`
- Login sebagai superadmin/admin/teknisi → redirect ke dashboard admin
- Login sebagai agent → redirect ke portal agent
- Login sebagai customer → redirect ke portal customer
- Akses route admin tanpa token → redirect ke `/login`
- Teknisi tidak bisa akses halaman user management
- Dark/light mode toggle berfungsi

---

## Phase 2 — Admin Network Management
**Goal:** Admin dapat mengelola router MikroTik, bandwidth profiles, customers, dan subscriptions.

**Deliverables:**
- **Customer management:** tabel dengan server-side pagination/search/filter, CRUD forms, activate/deactivate, detail page
- **Router management:** tabel, CRUD, select active router, test connection, sync (TanStack Table + TanStack Form)
- **Bandwidth Profiles:** CRUD scoped per router
- **Subscription management:** tabel dengan filter router+status, CRUD, aksi activate/isolate/restore/suspend/terminate, status badges
- Shared reusable table components dengan TanStack Table v8
- Reusable form components dengan TanStack Form v1 + Zod validation
- Confirmation dialog component (delete, destructive actions)
- TanStack Virtual untuk tabel dengan ribuan baris

**Requirements covered:** CUST-01~06, ROUT-01~08, BWP-01~03, SUB-01~06

**Success criteria:**
- Buat customer baru, muncul di tabel ✓
- Edit customer, perubahan tersimpan ✓
- Buat router, test connection, sync berhasil ✓
- Buat subscription, ubah status (isolate → restore), status badge berubah ✓
- Tabel 500+ rows scroll lancar dengan TanStack Virtual ✓
- Semua form validasi dengan Zod (field required, format email/phone) ✓

---

## Phase 3 — Billing, Finance & Agents
**Goal:** Admin dapat mengelola invoice, pembayaran, registrasi pelanggan baru, komisi sales agent, dan kas.

**Deliverables:**
- **Invoice management:** tabel filter status/tanggal, detail, cancel, overdue list, trigger monthly billing
- **Payment management:** tabel, create manual, detail, confirm/reject/refund, initiate gateway
- **Registration approvals:** tabel pending, approve/reject dengan form
- **Sales Agent management:** CRUD, profile prices per bandwidth profile, agent invoices, generate invoice
- **Agent Invoice management:** tabel, detail, mark paid, cancel, process scheduled
- **Cash Management:** cash entries CRUD + approve/reject, petty cash CRUD, summary saldo
- Currency formatter IDR, date formatter WIB/WITA/WIT dengan date-fns + id locale
- Status badges untuk semua entity (invoice, payment, subscription, cash)

**Requirements covered:** INV-01~05, PAY-01~06, REG-01~03, AGENT-01~07, AINV-01~05, CASH-01~06

**Success criteria:**
- Generate monthly billing → invoice muncul di tabel ✓
- Confirm payment → invoice status berubah jadi paid ✓
- Approve registrasi → customer + subscription baru terbuat ✓
- Upsert harga profile agent → tersimpan ✓
- Approve cash entry → saldo berubah di summary ✓
- Semua angka tampil dalam format IDR (Rp 1.500.000) ✓

---

## Phase 4 — Reports, Live Monitor & Settings
**Goal:** Admin dapat melihat laporan bisnis, monitoring MikroTik real-time via WebSocket, dan mengelola pengaturan sistem.

**Deliverables:**
- **Reports:** summary report, subscription report, cash flow + balance, reconciliation — dengan Recharts untuk grafik, date range picker
- **MikroTik Live Monitor:**
  - WebSocket client dengan auto-reconnect (reconnecting-websocket library)
  - PPP Active Sessions table (real-time update)
  - Hotspot Active Users table (real-time update)
  - Resource Monitor: grafik CPU/memory/traffic per interface (Recharts, update setiap 2-3 detik)
  - Connection status indicator di UI
- **Mikhmon:** voucher management, hotspot profile, expire management, hotspot report
- **System Settings:** list settings, form edit per setting (key-value)
- **User Management** (superadmin only): CRUD users sistem
- Export/print view untuk laporan

**Requirements covered:** RPT-01~06, LIVE-01~05, SYS-01~02

**Success criteria:**
- Summary report dengan date range menampilkan data dari API ✓
- Cash flow chart render dengan data bulanan ✓
- WebSocket PPP sessions update real-time tanpa full refresh ✓
- Resource monitor CPU chart update setiap detik ✓
- Disconnect WiFi → indicator WebSocket "disconnected" muncul, reconnect otomatis ✓
- Superadmin bisa buat user baru, admin tidak bisa akses menu ini ✓

---

## Phase 5 — Agent Portal & Customer Portal
**Goal:** Sales agent dan customer dapat mengakses portal masing-masing untuk self-service.

**Deliverables:**
- **Sales Agent Portal:**
  - Login page agent (`/agent/login`)
  - Dashboard: summary pendapatan, invoice pending, recent sales
  - Profile + ganti password form
  - Tabel agent invoices + detail
  - Request payment dari agent invoice
  - Daftar hotspot sales
- **Customer Portal:**
  - Login page customer (`/customer/login`)
  - Dashboard: status langganan, tagihan terdekat
  - Profile + ganti password form
  - Detail langganan aktif (bandwidth, IP, status, expire)
  - Riwayat invoice + detail
  - Halaman pembayaran: list unpaid invoice + tombol bayar via gateway
  - Redirect ke payment gateway
- Polishing: loading skeletons di semua halaman, empty states, error states
- Responsive audit: test di 375px, 768px, 1024px, 1440px
- Accessibility: keyboard navigation, focus states, aria labels

**Requirements covered:** SPORT-01~06, CPORT-01~07

**Success criteria:**
- Agent login → dashboard agent (bukan admin dashboard) ✓
- Agent lihat invoices, klik "Request Payment" → status berubah ✓
- Customer login → melihat langganan aktif ✓
- Customer lihat invoice unpaid, klik bayar → redirect gateway ✓
- Semua halaman usable di layar 375px (mobile) ✓
- Lighthouse mobile score ≥ 80 ✓

---

## Phase Dependencies

```
Phase 1 (Foundation)
    ↓
Phase 2 (Network Mgmt) ──┐
    ↓                     │ dapat paralel setelah Phase 1
Phase 3 (Billing)  ──────┘
    ↓
Phase 4 (Reports + Live)
    ↓
Phase 5 (Portals)
```

> Phase 2 dan 3 **dapat dikerjakan paralel** karena tidak saling bergantung.
> Phase 4 membutuhkan Phase 2+3 selesai (laporan butuh data subscription + payment).
> Phase 5 dapat mulai setelah Phase 1 selesai (portals hanya butuh auth + API client).

---

## Technology Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Routing | TanStack Router v1 (file-based) | Type-safe, code-splitting otomatis, route guards via `beforeLoad` |
| Data fetching | TanStack Query v5 | Caching, background refetch, optimistic updates, devtools |
| Tables | TanStack Table v8 + TanStack Virtual v3 | Server-side pagination, virtual scroll untuk ribuan baris |
| Forms | TanStack Form v1 + Zod | Type-safe form state, schema validation |
| UI Components | Shadcn/UI | Headless, accessible, customizable, tidak opinionated |
| Styling | Tailwind CSS v4 | Utility-first, mobile-first, CSS variables via @theme |
| State | Zustand v5 | Lightweight, tidak over-engineered untuk auth state |
| Charts | Recharts v2 | React-native, composable, cukup untuk ISP charts |
| WebSocket | reconnecting-websocket | Auto-reconnect, drop-in replacement for native WebSocket |
| Date | date-fns v4 + id locale | Lightweight, timezone-aware, Indonesian locale |
| Toast | Sonner | Modern, accessible, minimal config |
| HTTP | Axios | Interceptors untuk auth + refresh token |

---
*Roadmap created: 2026-03-30*
*Run `/gsd:plan-phase 1` to start Phase 1 execution.*
