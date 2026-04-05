# MikMongo ISP Management Dashboard

## What This Is

A full-featured ISP management dashboard built on the shadcn-admin template, implementing the MikMongo API (80+ endpoints). It serves three portals — Admin, Customer self-service, and Agent self-service — with complete MikroTik router management, real-time monitoring via WebSockets, billing, payments (Midtrans/Xendit), hotspot voucher sales, and business reports.

## Core Value

Admin can manage their entire ISP operation from one dashboard: customers, routers, subscriptions, billing, and monitor MikroTik devices in real-time.

## Requirements

### Validated

- [x] Custom JWT authentication replacing Clerk (login, refresh, change password, logout, /me) — Validated in Phase 01: auth-foundation
- [x] Admin dashboard with real-time widgets (customers, revenue, subscriptions, router health) — Validated in Phase 02: admin-dashboard
- [x] Three-portal architecture (Admin, Customer Portal, Agent Portal) with separate routes and auth — Validated in Phase 01: auth-foundation
- [x] Router selector in sidebar with status badges (online/offline/syncing) — Validated in Phase 02: admin-dashboard
- [x] Real-time ping display (8.8.8.8) in header showing ms latency — Validated in Phase 02: admin-dashboard
- [x] Customer management (CRUD, activate/deactivate, registration pipeline) — Validated in Phase 03: customers-routers-subscriptions
- [x] MikroTik router management (CRUD, select active, sync, test connection) — Validated in Phase 03: customers-routers-subscriptions
- [x] Bandwidth profile management per router (plans) — Validated in Phase 03: customers-routers-subscriptions
- [x] Subscription management (CRUD, activate, suspend, isolate, restore, terminate) — Validated in Phase 03: customers-routers-subscriptions
- [x] Invoice management with monthly trigger — Validated in Phase 04: billing-payments
- [x] Payment management with gateway integration (Midtrans/Xendit) — Validated in Phase 04: billing-payments
- [x] Cash entry and petty cash fund management — Validated in Phase 04: billing-payments
- [x] Customer Portal billing (invoices with pay button, payment history) — Validated in Phase 04: billing-payments
- [x] Agent Portal billing (invoice list with payment request) — Validated in Phase 04: billing-payments

### Active

- [ ] Sales agent management with profile pricing
- [ ] Agent invoice management with payment workflow
- [ ] Hotspot sales tracking and voucher generation (Mikhmon)
- [ ] Business reports with charts (Recharts) and data tables
- [ ] MikroTik PPP management (profiles, secrets, active connections, WebSocket)
- [ ] MikroTik Hotspot management (profiles, users, active, hosts, servers, WebSocket)
- [ ] MikroTik Network management (queues, firewall filter/NAT/address-list, IP pools, addresses)
- [ ] MikroTik real-time monitoring (system resource, interfaces, traffic WebSocket, logs WebSocket, ping WebSocket)
- [ ] Raw RouterOS command execution with WebSocket
- [ ] Mikhmon integration (voucher generation, profiles, reports, expiration monitoring)
- [ ] WebSocket connections for real-time data (PPP active, Hotspot active, monitor, logs, ping, raw commands)
- [ ] Customer Portal (login, profile, subscriptions, invoices, payments)
- [ ] Agent Portal (login, profile, invoices, sales tracking)

### Out of Scope

- Mobile native apps (web-first, mobile responsive)
- Clerk authentication (replaced by custom JWT)
- Template demo pages (apps, chats, tasks, help-center — replaced by ISP features)
- Multi-language/i18n support

## Context

- **Backend API**: Go-based MikMongo ISP API with 80+ REST endpoints and WebSocket channels
- **OpenAPI spec**: `docs/openapi.docs.yml` — authoritative API contract
- **Frontend template**: shadcn-admin v2.2.1 — feature-based architecture with TanStack Router, TanStack Query, TanStack Table, Zustand, React Hook Form + Zod, shadcn/ui, Tailwind CSS
- **Existing patterns**: Feature modules use `data/schema.ts` (Zod schemas), `data/[entity].ts` (mock data), `components/` (table, dialogs, forms), `index.tsx` (page composition)
- **Current auth**: Template uses Clerk; needs replacement with custom JWT auth
- **API response format**: Standard envelope `{ success, data, error, meta: { total, limit, offset } }`

## Constraints

- **Tech stack**: Must use existing template stack (React 19, TypeScript, TanStack, Tailwind, shadcn/ui, Axios, Zustand, Zod)
- **Template structure**: Must follow the feature-based directory convention; no structural rewrites
- **API contract**: Must match `docs/openapi.docs.yml` exactly — schemas, endpoints, auth schemes
- **Immutability**: Immutable data patterns throughout (no mutation of existing objects)
- **Real-time**: WebSocket endpoints must provide live updates without page refresh

## Key Decisions

| Decision | Rationale | Outcome |
|----------|-----------|---------|
| Three separate portals | API defines distinct auth schemes (BearerAuth, PortalAuth, AgentPortalAuth) | — Pending |
| Replace Clerk with custom JWT | API uses its own JWT auth, Clerk doesn't fit the ISP context | — Pending |
| Axios + React Query | Template already uses both; consistent with existing patterns | — Pending |
| Full MikroTik features | User needs complete router management including PPP/Hotspot/Mikhmon | — Pending |
| Real-time WebSocket | API provides WebSocket endpoints for monitoring — must use them | — Pending |
| Router selector in sidebar | User explicitly requested this for quick router switching | — Pending |
| Ping in header | User wants live 8.8.8.8 ping showing ms in the header | — Pending |
| Recharts for reports | Template already has recharts dependency | — Pending |

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
*Last updated: 2026-04-02 after initialization*
