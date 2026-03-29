# MikMongo Dashboard

## What This Is

Dashboard web frontend untuk MikMongo — sistem manajemen ISP berbasis MikroTik. Dibangun di dalam folder `dashboard/` menggunakan React + TanStack stack, menyediakan tiga portal terpisah: Admin Dashboard (superadmin/admin/teknisi), Sales Agent Portal, dan Customer Portal. Dashboard mengonsumsi seluruh REST API yang ada di `/api/v1`, `/agent-portal/v1`, dan `/portal/v1`.

## Core Value

Operator ISP dapat mengelola pelanggan, langganan, pembayaran, dan perangkat MikroTik dalam satu antarmuka terpadu dengan akses real-time.

## Requirements

### Validated

(None yet — ship to validate)

### Active

#### Authentication & Authorization
- [ ] Login dengan JWT + refresh token untuk 3 portal berbeda
- [ ] RBAC di Admin Dashboard: superadmin, admin, teknisi
- [ ] Persistent session (refresh token auto-renew)
- [ ] Protected routes per role/portal

#### Admin Dashboard (3 roles: superadmin, admin, teknisi)
- [ ] **Overview page** — summary stats: pelanggan aktif, pendapatan bulan ini, invoice overdue, status router
- [ ] **Customer management** — CRUD, activate/deactivate, filter/search/pagination
- [ ] **Router management** — CRUD, select active router, test connection, sync device
- [ ] **Bandwidth Profiles** — CRUD scoped per router
- [ ] **Subscription management** — CRUD per router, activate/isolate/restore/suspend/terminate
- [ ] **Invoice management** — list, detail, cancel, trigger monthly billing, overdue list
- [ ] **Payment management** — CRUD, confirm/reject/refund, initiate gateway
- [ ] **Registration approvals** — list pending, approve/reject
- [ ] **Sales Agent management** — CRUD, profile prices, agent invoices
- [ ] **Cash Management** — cash entries CRUD + approve/reject, petty cash fund CRUD
- [ ] **Reports** — summary, subscriptions, cash flow, cash balance, reconciliation
- [ ] **System Settings** — upsert settings
- [ ] **User management** — CRUD users (superadmin only)
- [ ] **MikroTik Live Monitor** — PPP active sessions, hotspot users, resource monitor via WebSocket
- [ ] **MikroTik Mikhmon** — voucher, profile, expire, report

#### Sales Agent Portal (`/agent-portal/v1`)
- [ ] Login khusus agent
- [ ] Profile & ganti password
- [ ] Daftar invoice agent + request payment
- [ ] Daftar hotspot sales

#### Customer Portal (`/portal/v1`)
- [ ] Login khusus customer
- [ ] Profile & ganti password
- [ ] Langganan aktif
- [ ] Invoice history + detail
- [ ] Pembayaran + pay via gateway

#### UX / Technical
- [ ] Mobile-first responsive design
- [ ] Dark mode + Light mode (ikut system preference)
- [ ] TanStack Router untuk routing
- [ ] TanStack Query untuk data fetching + caching
- [ ] TanStack Table untuk tabel data besar
- [ ] TanStack Form untuk semua form
- [ ] TanStack Virtual untuk list panjang
- [ ] Shadcn/UI + Tailwind CSS untuk komponen
- [ ] WebSocket support untuk real-time MikroTik monitoring

### Out of Scope

- Backend changes — dashboard hanya mengonsumsi API yang sudah ada
- Mobile native app — web mobile-first sudah cukup
- Mikhmon chatbot / AI features — tidak diminta
- Multi-language support — bahasa Indonesia/Inggris cukup

## Context

**Backend API:** Go (Gin framework), PostgreSQL + Redis + RabbitMQ, JWT auth, Casbin RBAC. API sudah production-ready dengan endpoint lengkap.

**Response format standar:**
```json
{ "success": true, "data": {...}, "meta": { "total": 100, "limit": 10, "offset": 0 } }
```

**API Base URLs:**
- Admin API: `GET/POST/PUT/DELETE /api/v1/...` (JWT Bearer required)
- Customer Portal: `/portal/v1/...` (Portal JWT)
- Agent Portal: `/agent-portal/v1/...` (Agent JWT)

**Entities utama:** User, Customer, MikrotikRouter, BandwidthProfile, Subscription, Invoice, Payment, Registration, SalesAgent, AgentInvoice, HotspotSale, CashEntry, PettyCashFund, SystemSetting

**WebSocket endpoints tersedia** untuk: PPP sessions, hotspot monitor, resource monitor, mikhmon live

**Tech stack yang dipilih:**
- React 19 + Vite
- TanStack Router v1 (file-based routing)
- TanStack Query v5
- TanStack Table v8
- TanStack Form v1
- TanStack Virtual v3
- Shadcn/UI (Radix UI primitives)
- Tailwind CSS v4
- Zustand (auth state)
- Axios / ky (HTTP client)
- Zod (schema validation)
- Recharts / Tremor (charts)
- Sonner (toast notifications)
- date-fns (date formatting)

## Constraints

- **Tech Stack**: React + TanStack + Shadcn/Tailwind — tidak berubah, sudah ditentukan user
- **Location**: Harus berada di `dashboard/` folder dalam repo yang sama
- **API**: Tidak boleh mengubah backend — frontend-only project
- **Mobile-first**: Semua halaman harus usable di layar 375px ke atas
- **Context7**: Gunakan Context7 untuk referensi dokumentasi library terbaru

## Key Decisions

| Decision | Rationale | Outcome |
|----------|-----------|---------|
| 3 portal dalam 1 project | Semua portal share auth logic dan komponen UI | — Pending |
| TanStack Router (file-based) | Type-safe routing, code splitting otomatis | — Pending |
| Zustand untuk auth state | Lightweight, tidak perlu Redux untuk state sesederhana ini | — Pending |
| Dark/Light mode via system preference | User tidak perlu toggle manual, lebih natural | — Pending |
| Zod untuk form validation | Terintegrasi sempurna dengan TanStack Form | — Pending |

## Evolution

This document evolves at phase transitions and milestone boundaries.

**After each phase transition** (via `/gsd:transition`):
1. Requirements invalidated? → Move to Out of Scope with reason
2. Requirements validated? → Move to Validated with phase reference
3. New requirements emerged? → Add to Active
4. Decisions to log? → Add to Key Decisions
5. "What This Is" still accurate? → Update if drifted

**After each milestone** (via `/gsd:complete-milestone`):
1. Full review of all sections
2. Core Value check — still the right priority?
3. Audit Out of Scope — reasons still valid?
4. Update Context with current state

---
*Last updated: 2026-03-29 after initialization*
