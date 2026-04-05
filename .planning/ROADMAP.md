# Roadmap: MikMongo ISP Management Dashboard

## Overview

Transform the shadcn-admin template into a full-featured ISP management dashboard with custom JWT authentication, real-time MikroTik router management, billing, payment gateways, and self-service portals for customers and agents. The journey starts by replacing Clerk auth and establishing the API layer, then builds out ISP domain features (customers, routers, subscriptions, billing, sales), extends into advanced MikroTik device management with WebSocket-based monitoring, and culminates with customer and agent self-service portals.

## Phases

**Phase Numbering:**
- Integer phases (1, 2, 3): Planned milestone work
- Decimal phases (2.1, 2.2): Urgent insertions (marked with INSERTED)

Decimal phases appear between their surrounding integers in numeric order.

- [x] **Phase 1: Auth & API Foundation** - Replace Clerk with custom JWT auth, establish API client layer, configure three-portal route structure
- [ ] **Phase 2: Layout, Dashboard & Users** - Router selector sidebar, real-time ping header, ISP dashboard widgets, admin user CRUD
- [ ] **Phase 3: Customers, Routers & Subscriptions** - Customer management pipeline, MikroTik router CRUD, bandwidth profiles, subscription lifecycle
- [ ] **Phase 4: Billing & Payments** - Invoice generation, payment management with gateway integration, cash workflow
- [ ] **Phase 5: Sales & Agents** - Sales agent management, agent invoices, hotspot sales tracking, Mikhmon voucher generation
- [ ] **Phase 6: MikroTik Device Management** - PPP profiles/secrets, Hotspot management, Network queues/firewall/pools
- [ ] **Phase 7: Monitoring, Reports & Settings** - Real-time system/interface/log monitoring via WebSocket, Mikhmon config, business reports, system settings
- [ ] **Phase 8: Self-Service Portals** - Customer and Agent portal pages for self-service profile, subscription, invoice, and sales management

## Phase Details

### Phase 1: Auth & API Foundation
**Goal**: Admin, customer, and agent can each log into their respective portal using custom JWT authentication, and all API communication goes through properly configured Axios instances with token refresh
**Depends on**: Nothing (first phase)
**Requirements**: AUTH-01, AUTH-02, AUTH-03, AUTH-04, AUTH-05, AUTH-06, AUTH-07, AUTH-08, AUTH-09, NAV-04
**Success Criteria** (what must be TRUE):
  1. Admin can log in with email/password, see the admin dashboard, and stay logged in across browser refreshes
  2. Admin can change password and log out (token invalidated)
  3. Expired access tokens are silently refreshed using the refresh token without user action
  4. Customer can log into their portal route, and agent can log into their portal route
  5. Unauthenticated users on protected routes are redirected to their portal's login page
**Plans**: 4 plans

Plans:
- [x] 01-01: Auth data layer -- Zustand store with three portal slices, Zod schemas matching OpenAPI, Axios clients with token refresh, API auth functions
- [x] 01-02: Admin auth UI -- login page with Indonesian text, change password page, logout confirmation dialog, MikMongo branding
- [x] 01-03: Customer and agent portal login pages with identifier/username fields matching OpenAPI, auth hooks
- [x] 01-04: Three-portal route guards, route files, hydration gate, Clerk removal, human verification

### Phase 2: Layout, Dashboard & Users
**Goal**: Admin sees a functional ISP dashboard with router selector in sidebar, real-time ping in header, overview widgets, and can manage admin users
**Depends on**: Phase 1
**Requirements**: NAV-01, NAV-02, NAV-03, NAV-05, DASH-01, DASH-02, DASH-03, USER-01, USER-02, USER-03, USER-04
**Success Criteria** (what must be TRUE):
  1. Admin can select an active router from a dropdown in the sidebar, and the dropdown shows online/offline/syncing status badges
  2. Admin sees real-time ping latency (ms) to 8.8.8.8 displayed in the header
  3. Admin dashboard shows widgets for total customers, active subscriptions, monthly revenue, and router health cards
  4. Admin can create, view, and deactivate admin users with role assignment (edit deferred -- no PUT endpoint)
  5. Sidebar navigation groups are organized by ISP management domains (Customers, Billing, MikroTik, etc.)
**Plans**: 5 plans
**UI hint**: yes

Plans:
- [ ] 02-01: Data layer -- schemas, store, API functions, TanStack Query hooks for routers, reports, users
- [ ] 02-02: Sidebar nav restructure -- ISP groups, MikMongo branding, disabled items, app-title
- [ ] 02-03: Router selector + ping display -- Select dropdown, WebSocket hook, header integration
- [ ] 02-04: Dashboard overview -- KPI widgets, router health cards, activity feed with real user data
- [ ] 02-05: Admin user management -- list table, create dialog, delete confirmation

### Phase 3: Customers, Routers & Subscriptions
**Goal**: Admin can manage the full customer lifecycle (create, activate, registration pipeline) and manage MikroTik routers with bandwidth profiles and subscription plans
**Depends on**: Phase 2
**Requirements**: CUST-01, CUST-02, CUST-03, CUST-04, CUST-05, CUST-06, CUST-07, RTR-01, RTR-02, RTR-03, RTR-04, RTR-05, RTR-06, BW-01, BW-02, BW-03, SUB-01, SUB-02, SUB-03, SUB-04, SUB-05
**Success Criteria** (what must be TRUE):
  1. Admin can view, create, update, and activate/deactivate customers with pagination and search
  2. Admin can manage the customer registration pipeline (view pending, approve with router/profile assignment, reject with reason)
  3. Admin can add, edit, delete, test connection to, and sync MikroTik routers
  4. Admin can create and manage bandwidth profiles (rate-limit, burst) per router
  5. Admin can create subscriptions, assign profiles to customers on routers, and manage subscription lifecycle (activate, suspend, isolate, restore, terminate)
  6. Customer portal displays the customer's active subscriptions
**Plans**: 7 plans (3 executed + 4 gap closure)
**UI hint**: yes

Plans:
- [x] 03-01: Router list page with profiles, create router, sync, test connection, bandwidth profile management
- [x] 03-02: Customer list page with CRUD, registration pipeline, approve/reject dialogs
- [x] 03-03: Subscription management with full lifecycle (CRUD, activate, suspend, isolate, restore, terminate)
- [ ] 03-04: Router CRUD completion -- edit/delete router dialogs, sync-all button, sidebar navigation fix (RTR-02, RTR-06)
- [ ] 03-05: Customer edit -- update hook, edit dialog, wire into table (CUST-03)
- [ ] 03-06: Profile update -- update API, hook, edit dialog, wire into profile table (BW-03)
- [ ] 03-07: Customer portal subscriptions -- portal API, hook, subscriptions page, route (SUB-05)

### Phase 4: Billing & Payments
**Goal**: Admin can generate invoices, manage payments (manual confirm/reject/refund and gateway-initiated via Midtrans/Xendit), and manage cash entry workflow
**Depends on**: Phase 3
**Requirements**: INV-01, INV-02, INV-03, INV-04, INV-05, PAY-01, PAY-02, PAY-03, PAY-04, PAY-05, PAY-06, PAY-07, CASH-01, CASH-02, CASH-03, CASH-04
**Success Criteria** (what must be TRUE):
  1. Admin can view invoices with overdue filter, view individual invoice details, and trigger monthly invoice generation
  2. Admin can view payments with filters, confirm/reject manual payments, and refund payments
  3. Admin can initiate gateway payments via Midtrans/Xendit
  4. Admin can view, create, approve, and reject cash entries and manage petty cash fund
  5. Customer portal shows their invoices, invoice details, payment history, and allows initiating gateway payments
  6. Agent portal shows invoice list with payment request option
**Plans**: 5 plans
**UI hint**: yes

Plans:
- [x] 04-01-PLAN.md — Billing data layer: Zod schemas, API functions, TanStack Query hooks for invoices, payments, cash, petty cash, portal billing
- [x] 04-02-PLAN.md — Admin invoice management: list with filters, detail sheet, monthly generation trigger
- [ ] 04-03-PLAN.md — Admin payment management: list with filters, confirm/reject/refund dialogs, gateway initiation
- [ ] 04-04-PLAN.md — Admin cash management: entries table with inline approve/reject, create dialog, petty cash card with top-up
- [ ] 04-05-PLAN.md — Portal billing views: customer invoices + payments, agent invoices with payment request

### Phase 5: Sales & Agents
**Goal**: Admin can manage sales agents (create, profile pricing, invoice generation/payment) and track hotspot sales with Mikhmon voucher generation
**Depends on**: Phase 4
**Requirements**: AGNT-01, AGNT-02, AGNT-03, AGNT-04, AGNT-05, AGNT-06, AGNT-07, AGNT-08, HS-01, HS-02, HS-03, HS-04, HS-05, HS-06
**Success Criteria** (what must be TRUE):
  1. Admin can view, create, and update sales agents with profile pricing configuration
  2. Admin can view agent invoices, generate new ones, and manage payments (pay, cancel, process)
  3. Admin can view and create hotspot sales entries, generate Mikhmon vouchers in batch, and manage Mikhmon profiles
  4. Admin can generate Mikhmon setup script and view Mikhmon sales reports/summary
  5. Agent portal shows sales history and invoice management
**Plans**: TBD
**UI hint**: yes

Plans:
- [ ] 05-01: Build sales agent management (list, create, update, profile pricing)
- [ ] 05-02: Implement agent invoice management (view, generate, payment workflow)
- [ ] 05-03: Build hotspot sales tracking and Mikhmon voucher generation
- [ ] 05-04: Implement Mikhmon profile management, setup script generation, and sales reports
- [ ] 05-05: Add sales history and invoice management views to Agent Portal

### Phase 6: MikroTik Device Management
**Goal**: Admin can manage PPP profiles/secrets, Hotspot users/sessions, and Network configuration (queues, firewall, IP pools) across all routers
**Depends on**: Phase 3
**Requirements**: PPP-01, PPP-02, PPP-03, PPP-04, PPP-05, PPP-06, HOT-01, HOT-02, HOT-03, HOT-04, HOT-05, HOT-06, HOT-07, NET-01, NET-02, NET-03, NET-04, NET-05, NET-06
**Success Criteria** (what must be TRUE):
  1. Admin can view, create, update, and delete PPP profiles and secrets (users) per router
  2. Admin can view active PPP connections with real-time updates via WebSocket
  3. Admin can view, create, update, and delete Hotspot profiles and users per router
  4. Admin can view active Hotspot sessions, hosts, and servers with real-time session updates via WebSocket
  5. Admin can view simple queues, firewall filter/NAT/address-list rules, IP pools, and IP addresses per router
**Plans**: TBD
**UI hint**: yes

Plans:
- [ ] 06-01: Build PPP management (profiles, secrets, active connections) with WebSocket for real-time
- [ ] 06-02: Build Hotspot management (profiles, users, active sessions, hosts, servers) with WebSocket for real-time
- [ ] 06-03: Build Network management (simple queues, firewall filter/NAT/address-list, IP pools, IP addresses)

### Phase 7: Monitoring, Reports & Settings
**Goal**: Admin can monitor router health in real-time via WebSocket (resources, traffic, logs, ping), manage Mikhmon expiration monitoring, view business reports with charts, and configure system settings
**Depends on**: Phase 6
**Requirements**: MON-01, MON-02, MON-03, MON-04, MON-05, MON-06, MON-07, MKH-01, MKH-02, MKH-03, MKH-04, RPT-01, RPT-02, SET-01
**Success Criteria** (what must be TRUE):
  1. Admin can view system resource usage (CPU, memory, uptime) and network interface traffic per router
  2. Admin can monitor system resources and interface traffic in real-time via WebSocket
  3. Admin can view router logs in real-time and ping from router in real-time via WebSocket
  4. Admin can execute raw RouterOS commands via WebSocket and see live output
  5. Admin can configure Mikhmon expiration monitoring (enable/disable, view status, generate script)
  6. Admin can view business reports with charts (Recharts) covering revenue, customers, subscriptions, and agent sales
  7. Admin can manage system settings
**Plans**: TBD
**UI hint**: yes

Plans:
- [ ] 07-01: Build system resource and interface traffic monitoring pages
- [ ] 07-02: Implement real-time WebSocket monitoring (resources, traffic, logs, ping)
- [ ] 07-03: Build raw RouterOS command execution interface with WebSocket live output
- [ ] 07-04: Implement Mikhmon expiration monitoring configuration and reports
- [ ] 07-05: Build business reports with Recharts and data tables
- [ ] 07-06: Build system settings management page

### Phase 8: Self-Service Portals
**Goal**: Customers can view and update their profile and change password; Agents can view and update their profile and change password
**Depends on**: Phase 5, Phase 7
**Requirements**: PORT-01, PORT-02, PORT-03, PORT-04
**Success Criteria** (what must be TRUE):
  1. Customer can view and update their profile details from the Customer Portal
  2. Customer can change their password from the Customer Portal
  3. Agent can view and update their profile details from the Agent Portal
  4. Agent can change their password from the Agent Portal
**Plans**: TBD
**UI hint**: yes

Plans:
- [ ] 08-01: Build Customer Portal profile and password change pages
- [ ] 08-02: Build Agent Portal profile and password change pages

## Progress

**Execution Order:**
Phases execute in numeric order: 1 -> 2 -> 3 -> 4 -> 5 -> 6 -> 7 -> 8
Note: Phase 6 depends on Phase 3 and can theoretically run in parallel with Phases 4-5, but executes sequentially for simplicity.

| Phase | Plans Complete | Status | Completed |
|-------|----------------|--------|-----------|
| 1. Auth & API Foundation | 4/4 | Complete | - |
| 2. Layout, Dashboard & Users | 0/5 | Planning complete | - |
| 3. Customers, Routers & Subscriptions | 6/7 | In Progress|  |
| 4. Billing & Payments | 2/5 | In Progress|  |
| 5. Sales & Agents | 0/5 | Not started | - |
| 6. MikroTik Device Management | 0/3 | Not started | - |
| 7. Monitoring, Reports & Settings | 0/6 | Not started | - |
| 8. Self-Service Portals | 0/2 | Not started | - |
