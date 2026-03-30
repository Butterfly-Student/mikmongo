# Requirements: MikMongo Dashboard

**Defined:** 2026-03-30
**Core Value:** Operator ISP dapat mengelola pelanggan, langganan, pembayaran, dan perangkat MikroTik dalam satu antarmuka terpadu dengan akses real-time.

## v1 Requirements

### Foundation & Setup

- [x] **SETUP-01**: Project React + Vite + TypeScript terkonfigurasi di folder `dashboard/`
- [x] **SETUP-02**: TanStack Router v1 (file-based routing) terkonfigurasi dengan 3 route tree terpisah
- [x] **SETUP-03**: TanStack Query v5 dengan QueryClient global terkonfigurasi
- [x] **SETUP-04**: Shadcn/UI + Tailwind CSS v4 terkonfigurasi dengan design tokens
- [x] **SETUP-05**: Dark mode / light mode mengikuti system preference dengan toggle manual
- [ ] **SETUP-06**: API client (Axios) terkonfigurasi dengan interceptor untuk JWT, refresh token, dan error handling
- [ ] **SETUP-07**: Zustand store untuk auth state (token, user, role, portal)
- [x] **SETUP-08**: Shared layout components: AppShell, Sidebar, Topbar, mobile nav
- [x] **SETUP-09**: Mobile-first responsive breakpoints terkonfigurasi (mobile 375px → desktop 1440px)
- [x] **SETUP-10**: Zod schemas untuk validasi form dan response API

### Authentication (All Portals)

- [ ] **AUTH-01**: Admin/superadmin/teknisi dapat login via `POST /api/v1/auth/login` dan mendapatkan JWT
- [ ] **AUTH-02**: Sales agent dapat login via `POST /agent-portal/v1/auth/login`
- [ ] **AUTH-03**: Customer dapat login via `POST /portal/v1/auth/login`
- [ ] **AUTH-04**: Refresh token berjalan otomatis sebelum JWT expired
- [ ] **AUTH-05**: Logout menghapus token dan redirect ke halaman login
- [ ] **AUTH-06**: Protected routes redirect ke login jika tidak authenticated
- [ ] **AUTH-07**: RBAC — superadmin dapat akses semua fitur, admin terbatas, teknisi hanya network management
- [ ] **AUTH-08**: Customer/agent login tidak dapat mengakses admin dashboard

### Admin — Overview & Layout

- [x] **ADMIN-01**: Dashboard overview dengan summary cards: pelanggan aktif, revenue bulan ini, invoice overdue, router online
- [x] **ADMIN-02**: Sidebar navigasi dengan semua menu, collapsible di mobile
- [x] **ADMIN-03**: Topbar dengan info user, notifikasi, theme toggle, logout

### Admin — Customer Management

- [ ] **CUST-01**: Tabel pelanggan dengan pagination server-side, search, filter status
- [ ] **CUST-02**: Form create pelanggan (nama, email, telepon, alamat, sales agent)
- [ ] **CUST-03**: Form edit pelanggan
- [ ] **CUST-04**: Delete pelanggan dengan konfirmasi dialog
- [ ] **CUST-05**: Tombol activate/deactivate account pelanggan
- [ ] **CUST-06**: Detail pelanggan: info lengkap + riwayat langganan + invoice

### Admin — Router Management

- [ ] **ROUT-01**: Tabel router dengan status online/offline
- [ ] **ROUT-02**: Form create router (nama, host, port, username, password, tipe)
- [ ] **ROUT-03**: Form edit router
- [ ] **ROUT-04**: Delete router dengan konfirmasi
- [ ] **ROUT-05**: Tombol "Select Router" untuk set router aktif
- [ ] **ROUT-06**: Tombol "Test Connection" dengan feedback status
- [ ] **ROUT-07**: Tombol "Sync" per router dan "Sync All"
- [ ] **ROUT-08**: Badge router yang sedang terpilih (active router)

### Admin — Bandwidth Profiles

- [ ] **BWP-01**: Tabel bandwidth profiles per router
- [ ] **BWP-02**: Form create/edit bandwidth profile (nama, rate upload, rate download)
- [ ] **BWP-03**: Delete bandwidth profile dengan konfirmasi

### Admin — Subscription Management

- [ ] **SUB-01**: Tabel subscriptions dengan filter router, status, customer
- [ ] **SUB-02**: Form create subscription (pilih customer, router, bandwidth profile, tanggal mulai)
- [ ] **SUB-03**: Form edit subscription
- [ ] **SUB-04**: Aksi per subscription: activate, isolate, restore, suspend, terminate
- [ ] **SUB-05**: Status badge: active, isolated, suspended, terminated dengan warna berbeda
- [ ] **SUB-06**: Detail subscription dengan riwayat perubahan status

### Admin — Invoice Management

- [ ] **INV-01**: Tabel invoice dengan filter status (paid/unpaid/overdue/cancelled), tanggal, customer
- [ ] **INV-02**: Detail invoice: item, jumlah, customer, due date, status
- [ ] **INV-03**: Cancel invoice dengan konfirmasi
- [ ] **INV-04**: Daftar invoice overdue di halaman terpisah atau tab
- [ ] **INV-05**: Tombol "Trigger Monthly Billing" untuk generate invoice bulanan

### Admin — Payment Management

- [ ] **PAY-01**: Tabel payment dengan filter status (pending/confirmed/rejected/refunded)
- [ ] **PAY-02**: Form create payment manual
- [ ] **PAY-03**: Detail payment dengan bukti transfer (jika ada)
- [ ] **PAY-04**: Aksi: Confirm / Reject / Refund payment
- [ ] **PAY-05**: Initiate payment gateway untuk invoice tertentu
- [ ] **PAY-06**: Status badge payment dengan warna berbeda

### Admin — Registration Approvals

- [ ] **REG-01**: Tabel registrasi pending dengan detail calon pelanggan
- [ ] **REG-02**: Approve registrasi (otomatis buat customer + subscription)
- [ ] **REG-03**: Reject registrasi dengan alasan

### Admin — Sales Agent Management

- [ ] **AGENT-01**: Tabel sales agent dengan filter aktif/nonaktif
- [ ] **AGENT-02**: Form create/edit sales agent
- [ ] **AGENT-03**: Delete sales agent dengan konfirmasi
- [ ] **AGENT-04**: Halaman profile prices per agent — list bandwidth profile + harga jual
- [ ] **AGENT-05**: Form upsert harga jual per bandwidth profile
- [ ] **AGENT-06**: Tabel agent invoice per sales agent
- [ ] **AGENT-07**: Generate agent invoice untuk periode tertentu

### Admin — Agent Invoice Management

- [ ] **AINV-01**: Tabel semua agent invoice dengan filter status
- [ ] **AINV-02**: Detail agent invoice
- [ ] **AINV-03**: Mark agent invoice sebagai paid
- [ ] **AINV-04**: Cancel agent invoice
- [ ] **AINV-05**: Process scheduled agent invoices

### Admin — Cash Management

- [ ] **CASH-01**: Tabel cash entries dengan filter tipe (income/expense), status, tanggal
- [ ] **CASH-02**: Form create/edit cash entry
- [ ] **CASH-03**: Approve/reject cash entry
- [ ] **CASH-04**: Delete cash entry
- [ ] **CASH-05**: Tabel petty cash fund CRUD
- [ ] **CASH-06**: Summary saldo kas di halaman cash management

### Admin — Reports

- [ ] **RPT-01**: Summary report dengan filter date range — total revenue, pelanggan baru, pembayaran
- [ ] **RPT-02**: Subscription report — distribusi status, revenue per bandwidth profile
- [ ] **RPT-03**: Cash flow report — grafik pemasukan vs pengeluaran
- [ ] **RPT-04**: Cash balance report — saldo berjalan
- [ ] **RPT-05**: Reconciliation report — invoice vs pembayaran
- [ ] **RPT-06**: Export laporan (minimal print/PDF view)

### Admin — MikroTik Live Monitor

- [ ] **LIVE-01**: Halaman PPP Active Sessions — tabel real-time via WebSocket dengan status, IP, uptime
- [ ] **LIVE-02**: Halaman Hotspot Active Users — tabel real-time
- [ ] **LIVE-03**: Resource Monitor — grafik CPU, memory, disk, traffic per interface (real-time)
- [ ] **LIVE-04**: Koneksi WebSocket dengan auto-reconnect dan status indicator (connected/disconnected)
- [ ] **LIVE-05**: MikroTik Mikhmon — voucher management, profile, expire, hotspot report

### Admin — System Settings & User Management

- [ ] **SYS-01**: Halaman system settings — daftar setting + form edit per setting
- [ ] **SYS-02**: User management (superadmin only) — CRUD users sistem

### Sales Agent Portal

- [ ] **SPORT-01**: Login page khusus agent
- [ ] **SPORT-02**: Dashboard agent: summary pendapatan, invoice pending
- [ ] **SPORT-03**: Profile agent + form ganti password
- [ ] **SPORT-04**: Daftar invoice agent dengan detail dan status
- [ ] **SPORT-05**: Request pembayaran commission dari agent invoice
- [ ] **SPORT-06**: Daftar hotspot sales yang berkaitan dengan agent

### Customer Portal

- [ ] **CPORT-01**: Login page khusus customer
- [ ] **CPORT-02**: Dashboard customer: status langganan aktif, tagihan terdekat
- [ ] **CPORT-03**: Profile customer + form ganti password
- [ ] **CPORT-04**: Detail langganan aktif (bandwidth, IP, status)
- [ ] **CPORT-05**: Riwayat invoice + detail per invoice
- [ ] **CPORT-06**: Halaman pembayaran — list invoice unpaid + tombol bayar
- [ ] **CPORT-07**: Initiate payment via gateway (redirect ke payment page)

## v2 Requirements

### Notifications

- **NOTF-01**: Push notification / in-app notification untuk invoice baru, payment confirmed
- **NOTF-02**: Email notification untuk customer dan agent

### Advanced Analytics

- **ANLT-01**: Dashboard analytics lanjutan dengan trend chart historis
- **ANLT-02**: Comparison report antar periode

### Bulk Operations

- **BULK-01**: Bulk activate/suspend subscription
- **BULK-02**: Bulk confirm/reject payments

## Out of Scope

| Feature | Reason |
|---------|--------|
| Backend API changes | Frontend-only project |
| Mobile native app (iOS/Android) | Web mobile-first sudah cukup |
| Multi-language support | Bahasa Indonesia/Inggris sudah cukup |
| Multi-tenant / white-label | Satu instance ISP |
| AI/chatbot features | Tidak diminta |
| RADIUS server integration | Sudah dihandle backend |

## Traceability

| Requirement | Phase | Status |
|-------------|-------|--------|
| SETUP-01 ~ SETUP-10 | Phase 1 | Pending |
| AUTH-01 ~ AUTH-08 | Phase 1 | Pending |
| ADMIN-01 ~ ADMIN-03 | Phase 1 | Pending |
| CUST-01 ~ CUST-06 | Phase 2 | Pending |
| ROUT-01 ~ ROUT-08 | Phase 2 | Pending |
| BWP-01 ~ BWP-03 | Phase 2 | Pending |
| SUB-01 ~ SUB-06 | Phase 2 | Pending |
| INV-01 ~ INV-05 | Phase 3 | Pending |
| PAY-01 ~ PAY-06 | Phase 3 | Pending |
| REG-01 ~ REG-03 | Phase 3 | Pending |
| AGENT-01 ~ AGENT-07 | Phase 3 | Pending |
| AINV-01 ~ AINV-05 | Phase 3 | Pending |
| CASH-01 ~ CASH-06 | Phase 3 | Pending |
| RPT-01 ~ RPT-06 | Phase 4 | Pending |
| LIVE-01 ~ LIVE-05 | Phase 4 | Pending |
| SYS-01 ~ SYS-02 | Phase 4 | Pending |
| SPORT-01 ~ SPORT-06 | Phase 5 | Pending |
| CPORT-01 ~ CPORT-07 | Phase 5 | Pending |

**Coverage:**
- v1 requirements: 95 total
- Mapped to phases: 95
- Unmapped: 0 ✓

---
*Requirements defined: 2026-03-30*
*Last updated: 2026-03-30 after initial definition*
