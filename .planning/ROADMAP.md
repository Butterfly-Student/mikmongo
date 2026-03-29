# Roadmap: MikMongo Dashboard

## Overview

Dashboard web frontend untuk MikMongo ISP management system. Dibangun di `dashboard/` menggunakan React 19 + TanStack stack + Shadcn/UI. Menyediakan 3 portal: Admin Dashboard (superadmin/admin/teknisi), Sales Agent Portal, dan Customer Portal. Mengonsumsi API Go Gin yang sudah ada di `/api/v1`, `/agent-portal/v1`, dan `/portal/v1`.

## Phases

**Phase Numbering:**
- Integer phases (1-5): Milestone v1.0 work

- [ ] **Phase 1: Foundation & Auth** - Vite+React scaffold, auth 3 portal, layout, dark/light mode
- [ ] **Phase 2: Admin Network Management** - Router, Bandwidth Profiles, Customer, Subscription CRUD
- [ ] **Phase 3: Billing, Finance & Agents** - Invoice, Payment, Registration, Sales Agent, Cash Management
- [ ] **Phase 4: Reports, Live Monitor & Settings** - Reports charts, WebSocket real-time, MikroTik monitor, Settings
- [ ] **Phase 5: Agent Portal & Customer Portal** - Sales agent portal, customer portal, polish

## Phase Details

### Phase 1: Foundation & Auth
**Goal**: Project scaffold berjalan di `dashboard/`, semua 3 portal dapat login dan mengakses halaman yang terproteksi dengan RBAC.
**Depends on**: Nothing (first phase)
**Requirements**: SETUP-01, SETUP-02, SETUP-03, SETUP-04, SETUP-05, SETUP-06, SETUP-07, SETUP-08, SETUP-09, SETUP-10, AUTH-01, AUTH-02, AUTH-03, AUTH-04, AUTH-05, AUTH-06, AUTH-07, AUTH-08, ADMIN-01, ADMIN-02, ADMIN-03
**Success Criteria** (what must be TRUE):
  1. `npm run dev` berjalan di `dashboard/` tanpa error
  2. Login superadmin/admin/teknisi → redirect ke dashboard admin dengan summary cards
  3. Login agent → redirect ke portal agent halaman kosong
  4. Login customer → redirect ke portal customer halaman kosong
  5. Akses route admin tanpa token → redirect ke `/login`
  6. Teknisi tidak bisa akses halaman user management (403 atau menu disembunyikan)
  7. Dark/light mode toggle berfungsi, mengikuti system preference default
**Plans**: Ready to execute

Plans:
- [x] 01-01: Project scaffold — Vite + React 19 + TypeScript + TanStack Router + Shadcn/UI + Tailwind CSS v4
- [x] 01-02: Auth system — Zustand store, API client, login/logout/refresh untuk 3 portal, RBAC guards
- [ ] 01-03: Shared layout — AppShell, Sidebar, Topbar, mobile nav, dark/light mode, overview page

### Phase 2: Admin Network Management
**Goal**: Admin dapat mengelola router MikroTik, bandwidth profiles, pelanggan, dan subscription melalui tabel dan form.
**Depends on**: Phase 1
**Requirements**: CUST-01, CUST-02, CUST-03, CUST-04, CUST-05, CUST-06, ROUT-01, ROUT-02, ROUT-03, ROUT-04, ROUT-05, ROUT-06, ROUT-07, ROUT-08, BWP-01, BWP-02, BWP-03, SUB-01, SUB-02, SUB-03, SUB-04, SUB-05, SUB-06
**Success Criteria** (what must be TRUE):
  1. Admin dapat CRUD customer, tabel ter-paginate dengan search dan filter
  2. Admin dapat CRUD router, test connection, sync device berhasil
  3. Bandwidth profiles dapat dikelola per router
  4. Subscription dapat dibuat dan status-nya dapat diubah (isolate/restore/suspend/terminate)
  5. Tabel 500+ rows scroll lancar dengan TanStack Virtual
**Plans**: TBD

Plans:
- [ ] 02-01: Customer management — tabel server-side, CRUD forms, activate/deactivate, detail page
- [ ] 02-02: Router management — tabel, CRUD, select/sync/test-connection, bandwidth profiles CRUD
- [ ] 02-03: Subscription management — tabel filter router+status, CRUD, aksi lifecycle, status badges

### Phase 3: Billing, Finance & Agents
**Goal**: Admin dapat mengelola siklus billing lengkap: invoice, pembayaran, approval registrasi, komisi sales agent, dan kas.
**Depends on**: Phase 1
**Requirements**: INV-01, INV-02, INV-03, INV-04, INV-05, PAY-01, PAY-02, PAY-03, PAY-04, PAY-05, PAY-06, REG-01, REG-02, REG-03, AGENT-01, AGENT-02, AGENT-03, AGENT-04, AGENT-05, AGENT-06, AGENT-07, AINV-01, AINV-02, AINV-03, AINV-04, AINV-05, CASH-01, CASH-02, CASH-03, CASH-04, CASH-05, CASH-06
**Success Criteria** (what must be TRUE):
  1. Trigger monthly billing → invoice muncul di tabel
  2. Confirm payment → invoice status berubah jadi paid
  3. Approve registrasi → customer + subscription baru terbuat
  4. Upsert harga profile agent → tersimpan dan tampil di tabel
  5. Approve cash entry → saldo berubah di summary kas
  6. Semua angka tampil format IDR (Rp 1.500.000)
**Plans**: TBD

Plans:
- [ ] 03-01: Invoice & Payment management — tabel, detail, CRUD, confirm/reject/refund, gateway
- [ ] 03-02: Registration approvals + Sales Agent management — CRUD, profile prices, agent invoices
- [ ] 03-03: Cash Management — cash entries CRUD+approve/reject, petty cash, summary saldo

### Phase 4: Reports, Live Monitor & Settings
**Goal**: Admin dapat melihat laporan bisnis dengan grafik, memonitor MikroTik secara real-time via WebSocket, dan mengelola pengaturan sistem.
**Depends on**: Phase 2, Phase 3
**Requirements**: RPT-01, RPT-02, RPT-03, RPT-04, RPT-05, RPT-06, LIVE-01, LIVE-02, LIVE-03, LIVE-04, LIVE-05, SYS-01, SYS-02
**Success Criteria** (what must be TRUE):
  1. Summary report dengan date range menampilkan data dari API dengan grafik Recharts
  2. PPP Active Sessions update real-time tanpa refresh halaman
  3. Resource Monitor (CPU/memory/traffic) update grafik setiap detik via WebSocket
  4. Disconnect → indicator WebSocket "disconnected" muncul, reconnect otomatis
  5. Superadmin bisa CRUD user sistem, admin tidak bisa akses menu ini
**Plans**: TBD

Plans:
- [ ] 04-01: Reports — summary, subscriptions, cash flow, cash balance, reconciliation (Recharts)
- [ ] 04-02: MikroTik Live Monitor — WebSocket PPP/hotspot sessions, resource monitor, Mikhmon
- [ ] 04-03: System Settings + User Management (superadmin only)

### Phase 5: Agent Portal & Customer Portal
**Goal**: Sales agent dan customer dapat mengakses portal masing-masing untuk self-service, semua halaman usable di mobile.
**Depends on**: Phase 1
**Requirements**: SPORT-01, SPORT-02, SPORT-03, SPORT-04, SPORT-05, SPORT-06, CPORT-01, CPORT-02, CPORT-03, CPORT-04, CPORT-05, CPORT-06, CPORT-07
**Success Criteria** (what must be TRUE):
  1. Agent login → dashboard agent (bukan admin) dengan summary pendapatan
  2. Agent bisa lihat dan request payment dari agent invoice
  3. Customer login → melihat detail langganan aktif
  4. Customer bisa bayar invoice unpaid via payment gateway
  5. Semua halaman usable di layar 375px (mobile-first verified)
**Plans**: TBD

Plans:
- [ ] 05-01: Sales Agent Portal — login, dashboard, profile, invoices, request payment, hotspot sales
- [ ] 05-02: Customer Portal — login, dashboard, profile, subscriptions, invoices, payments
- [ ] 05-03: Polish — loading skeletons, empty states, error states, responsive audit, a11y

## Progress

**Execution Order:**
Phase 1 → Phase 2 + Phase 3 (paralel) → Phase 4 → Phase 5

| Phase | Plans Complete | Status | Completed |
|-------|----------------|--------|-----------|
| 1. Foundation & Auth | 2/3 | In Progress|  |
| 2. Admin Network Management | 0/3 | Not started | - |
| 3. Billing, Finance & Agents | 0/3 | Not started | - |
| 4. Reports, Live Monitor & Settings | 0/3 | Not started | - |
| 5. Agent Portal & Customer Portal | 0/3 | Not started | - |
