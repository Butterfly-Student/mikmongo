# Requirements: MikMongo ISP Management Dashboard

**Defined:** 2026-04-02
**Core Value:** Admin can manage their entire ISP operation from one dashboard: customers, routers, subscriptions, billing, and monitor MikroTik devices in real-time.

## v1 Requirements

Requirements for initial release. Each maps to roadmap phases.

### Authentication

- [x] **AUTH-01**: Admin can login with email/password and receive JWT tokens (access + refresh)
- [x] **AUTH-02**: Admin can refresh expired access token using refresh token
- [x] **AUTH-03**: Admin can change password (old + new)
- [x] **AUTH-04**: Admin can logout (invalidate token)
- [x] **AUTH-05**: Admin session persists across browser refresh via token storage
- [x] **AUTH-06**: Customer can login to portal with email/password
- [x] **AUTH-07**: Agent can login to agent portal with email/password
- [x] **AUTH-08**: Auth state managed via Zustand store with token persistence
- [x] **AUTH-09**: Protected routes redirect to respective login pages when unauthenticated

### Navigation & Layout

- [ ] **NAV-01**: Sidebar displays router selector dropdown to switch active router
- [ ] **NAV-02**: Sidebar shows router status badges (online/offline/syncing)
- [ ] **NAV-03**: Header displays real-time ping to 8.8.8.8 showing ms latency
- [x] **NAV-04**: Navigation structure supports three portals (admin, customer, agent) with separate route trees
- [ ] **NAV-05**: Sidebar navigation groups reflect ISP management domains (Customers, Billing, MikroTik, etc.)

### Dashboard

- [ ] **DASH-01**: Admin dashboard shows real-time overview widgets (total customers, active subscriptions, monthly revenue)
- [ ] **DASH-02**: Dashboard displays router health status cards
- [ ] **DASH-03**: Dashboard shows recent activity feed (new customers, payments, registrations)

### User Management

- [ ] **USER-01**: Admin can view list of admin users with pagination, search, and filters
- [ ] **USER-02**: Admin can create new admin user with role assignment
- [ ] **USER-03**: Admin can update existing admin user details
- [ ] **USER-04**: Admin can delete/deactivate admin users

### Customer Management

- [x] **CUST-01**: Admin can view customer list with pagination, search, and filters
- [x] **CUST-02**: Admin can create new customer (with optional auto-subscription)
- [x] **CUST-03**: Admin can update customer details
- [x] **CUST-04**: Admin can activate/deactivate customer accounts
- [x] **CUST-05**: Admin can manage customer registrations (pending/approved/rejected pipeline)
- [x] **CUST-06**: Admin can approve registrations with router and profile assignment
- [x] **CUST-07**: Admin can reject registrations with reason

### Router Management

- [ ] **RTR-01**: Admin can view list of MikroTik routers with status
- [x] **RTR-02**: Admin can add/edit/delete router configurations (name, address, credentials, area)
- [ ] **RTR-03**: Admin can select active router for context-dependent operations
- [ ] **RTR-04**: Admin can sync router data from MikroTik device
- [ ] **RTR-05**: Admin can test connection to router
- [x] **RTR-06**: Admin can sync all routers simultaneously

### Bandwidth Profiles

- [ ] **BW-01**: Admin can view bandwidth profiles per router
- [ ] **BW-02**: Admin can create bandwidth profile (name, rate-limit, burst)
- [x] **BW-03**: Admin can update/delete bandwidth profiles

### Subscription Management

- [x] **SUB-01**: Admin can view subscriptions per router with status filters
- [x] **SUB-02**: Admin can create new subscription (assign profile to customer on router)
- [x] **SUB-03**: Admin can update/terminate subscription
- [x] **SUB-04**: Admin can activate, suspend, isolate, and restore subscriptions
- [x] **SUB-05**: Customer portal shows their active subscriptions

### Billing & Invoices

- [x] **INV-01**: Admin can view invoices with overdue filter
- [x] **INV-02**: Admin can trigger monthly invoice generation
- [x] **INV-03**: Admin can view invoice details
- [x] **INV-04**: Customer portal shows their invoices
- [x] **INV-05**: Customer portal can view individual invoice details

### Payments

- [x] **PAY-01**: Admin can view all payments with filters
- [x] **PAY-02**: Admin can confirm/reject manual payments
- [x] **PAY-03**: Admin can refund payments
- [x] **PAY-04**: Admin can initiate gateway payment (Midtrans/Xendit)
- [x] **PAY-05**: Customer portal shows payment history
- [x] **PAY-06**: Customer portal can initiate payment via gateway
- [x] **PAY-07**: Agent portal shows invoice list with payment request option

### Sales Agents

- [ ] **AGNT-01**: Admin can view list of sales agents
- [ ] **AGNT-02**: Admin can create/update sales agents
- [ ] **AGNT-03**: Admin can manage agent profile pricing
- [ ] **AGNT-04**: Admin can view agent invoices
- [ ] **AGNT-05**: Admin can generate agent invoices
- [ ] **AGNT-06**: Admin can manage agent invoice payments (pay/cancel/process)
- [ ] **AGNT-07**: Agent portal shows sales history
- [ ] **AGNT-08**: Agent portal shows invoice management

### Hotspot Sales & Vouchers

- [ ] **HS-01**: Admin can view hotspot sales records
- [ ] **HS-02**: Admin can create hotspot sales entry
- [ ] **HS-03**: Admin can generate Mikhmon vouchers (batch)
- [ ] **HS-04**: Admin can view/manage Mikhmon profiles
- [ ] **HS-05**: Admin can generate Mikhmon setup script
- [ ] **HS-06**: Admin can view Mikhmon sales reports and summary

### Cash Management

- [x] **CASH-01**: Admin can view cash entries with approval workflow
- [x] **CASH-02**: Admin can create new cash entry
- [x] **CASH-03**: Admin can approve/reject cash entries
- [x] **CASH-04**: Admin can manage petty cash fund

### MikroTik PPP

- [ ] **PPP-01**: Admin can view PPP profiles per router
- [ ] **PPP-02**: Admin can create/update/delete PPP profiles
- [ ] **PPP-03**: Admin can view PPP secrets (users) per router
- [ ] **PPP-04**: Admin can create/update/delete PPP secrets
- [ ] **PPP-05**: Admin can view active PPP connections
- [ ] **PPP-06**: Admin can view active PPP connections in real-time via WebSocket

### MikroTik Hotspot

- [ ] **HOT-01**: Admin can view Hotspot profiles per router
- [ ] **HOT-02**: Admin can create/update/delete Hotspot profiles
- [ ] **HOT-03**: Admin can view Hotspot users per router
- [ ] **HOT-04**: Admin can create/update/delete Hotspot users
- [ ] **HOT-05**: Admin can view active Hotspot sessions
- [ ] **HOT-06**: Admin can view Hotspot hosts and servers
- [ ] **HOT-07**: Admin can view active Hotspot sessions in real-time via WebSocket

### MikroTik Network

- [ ] **NET-01**: Admin can view simple queues per router
- [ ] **NET-02**: Admin can view firewall filter rules per router
- [ ] **NET-03**: Admin can view firewall NAT rules per router
- [ ] **NET-04**: Admin can view firewall address-list per router
- [ ] **NET-05**: Admin can manage IP pools per router
- [ ] **NET-06**: Admin can view IP addresses per router

### MikroTik Monitor

- [ ] **MON-01**: Admin can view system resource usage (CPU, memory, uptime)
- [ ] **MON-02**: Admin can view network interfaces and traffic
- [ ] **MON-03**: Admin can monitor system resources in real-time via WebSocket
- [ ] **MON-04**: Admin can monitor interface traffic in real-time via WebSocket
- [ ] **MON-05**: Admin can view router logs in real-time via WebSocket
- [ ] **MON-06**: Admin can ping from router in real-time via WebSocket
- [ ] **MON-07**: Admin can execute raw RouterOS commands via WebSocket with live output

### Mikhmon

- [ ] **MKH-01**: Admin can configure Mikhmon expiration monitoring
- [ ] **MKH-02**: Admin can enable/disable expiration monitoring
- [ ] **MKH-03**: Admin can view expiration monitoring status
- [ ] **MKH-04**: Admin can generate expiration monitoring script

### Reports

- [ ] **RPT-01**: Admin can view business reports with charts (Recharts) and data tables
- [ ] **RPT-02**: Reports cover revenue, customers, subscriptions, agent sales

### Settings

- [ ] **SET-01**: Admin can manage system settings

### Portal Features

- [ ] **PORT-01**: Customer can view and update their profile
- [ ] **PORT-02**: Customer can change their password
- [ ] **PORT-03**: Agent can view and update their profile
- [ ] **PORT-04**: Agent can change their password

## v2 Requirements

Deferred to future release. Tracked but not in current roadmap.

### Notifications

- **NOTF-01**: Real-time notification system for admin events
- **NOTF-02**: Email notifications for customer payment confirmations

### Advanced Features

- **ADV-01**: Customer map view (geolocation of customers)
- **ADV-02**: Bulk operations for customers/subscriptions
- **ADV-03**: API documentation page integrated in dashboard

## Out of Scope

| Feature | Reason |
|---------|--------|
| Mobile native apps | Web-first, mobile responsive design sufficient for v1 |
| Clerk authentication | Replaced by custom JWT auth matching API |
| Template demo pages | apps, chats, tasks, help-center replaced by ISP features |
| Multi-language/i18n | Not needed for v1, single language sufficient |
| Automated testing (E2E) | Focus on implementation first, add tests in future milestone |

## Traceability

| Requirement | Phase | Status |
|-------------|-------|--------|
| AUTH-01 | Phase 1 | Complete |
| AUTH-02 | Phase 1 | Complete |
| AUTH-03 | Phase 1 | Complete |
| AUTH-04 | Phase 1 | Complete |
| AUTH-05 | Phase 1 | Complete |
| AUTH-06 | Phase 1 | Complete |
| AUTH-07 | Phase 1 | Complete |
| AUTH-08 | Phase 1 | Complete |
| AUTH-09 | Phase 1 | Complete |
| NAV-01 | Phase 2 | Pending |
| NAV-02 | Phase 2 | Pending |
| NAV-03 | Phase 2 | Pending |
| NAV-04 | Phase 1 | Complete |
| NAV-05 | Phase 2 | Pending |
| DASH-01 | Phase 2 | Pending |
| DASH-02 | Phase 2 | Pending |
| DASH-03 | Phase 2 | Pending |
| USER-01 | Phase 2 | Pending |
| USER-02 | Phase 2 | Pending |
| USER-03 | Phase 2 | Pending |
| USER-04 | Phase 2 | Pending |
| CUST-01 | Phase 3 | Complete |
| CUST-02 | Phase 3 | Complete |
| CUST-03 | Phase 3 | Complete |
| CUST-04 | Phase 3 | Complete |
| CUST-05 | Phase 3 | Complete |
| CUST-06 | Phase 3 | Complete |
| CUST-07 | Phase 3 | Complete |
| RTR-01 | Phase 3 | Pending |
| RTR-02 | Phase 3 | Complete |
| RTR-03 | Phase 3 | Pending |
| RTR-04 | Phase 3 | Pending |
| RTR-05 | Phase 3 | Pending |
| RTR-06 | Phase 3 | Complete |
| BW-01 | Phase 3 | Pending |
| BW-02 | Phase 3 | Pending |
| BW-03 | Phase 3 | Complete |
| SUB-01 | Phase 3 | Complete |
| SUB-02 | Phase 3 | Complete |
| SUB-03 | Phase 3 | Complete |
| SUB-04 | Phase 3 | Complete |
| SUB-05 | Phase 3 | Complete |
| INV-01 | Phase 4 | Complete |
| INV-02 | Phase 4 | Complete |
| INV-03 | Phase 4 | Complete |
| INV-04 | Phase 4 | Complete |
| INV-05 | Phase 4 | Complete |
| PAY-01 | Phase 4 | Complete |
| PAY-02 | Phase 4 | Complete |
| PAY-03 | Phase 4 | Complete |
| PAY-04 | Phase 4 | Complete |
| PAY-05 | Phase 4 | Complete |
| PAY-06 | Phase 4 | Complete |
| PAY-07 | Phase 4 | Complete |
| CASH-01 | Phase 4 | Complete |
| CASH-02 | Phase 4 | Complete |
| CASH-03 | Phase 4 | Complete |
| CASH-04 | Phase 4 | Complete |
| AGNT-01 | Phase 5 | Pending |
| AGNT-02 | Phase 5 | Pending |
| AGNT-03 | Phase 5 | Pending |
| AGNT-04 | Phase 5 | Pending |
| AGNT-05 | Phase 5 | Pending |
| AGNT-06 | Phase 5 | Pending |
| AGNT-07 | Phase 5 | Pending |
| AGNT-08 | Phase 5 | Pending |
| HS-01 | Phase 5 | Pending |
| HS-02 | Phase 5 | Pending |
| HS-03 | Phase 5 | Pending |
| HS-04 | Phase 5 | Pending |
| HS-05 | Phase 5 | Pending |
| HS-06 | Phase 5 | Pending |
| PPP-01 | Phase 6 | Pending |
| PPP-02 | Phase 6 | Pending |
| PPP-03 | Phase 6 | Pending |
| PPP-04 | Phase 6 | Pending |
| PPP-05 | Phase 6 | Pending |
| PPP-06 | Phase 6 | Pending |
| HOT-01 | Phase 6 | Pending |
| HOT-02 | Phase 6 | Pending |
| HOT-03 | Phase 6 | Pending |
| HOT-04 | Phase 6 | Pending |
| HOT-05 | Phase 6 | Pending |
| HOT-06 | Phase 6 | Pending |
| HOT-07 | Phase 6 | Pending |
| NET-01 | Phase 6 | Pending |
| NET-02 | Phase 6 | Pending |
| NET-03 | Phase 6 | Pending |
| NET-04 | Phase 6 | Pending |
| NET-05 | Phase 6 | Pending |
| NET-06 | Phase 6 | Pending |
| MON-01 | Phase 7 | Pending |
| MON-02 | Phase 7 | Pending |
| MON-03 | Phase 7 | Pending |
| MON-04 | Phase 7 | Pending |
| MON-05 | Phase 7 | Pending |
| MON-06 | Phase 7 | Pending |
| MON-07 | Phase 7 | Pending |
| MKH-01 | Phase 7 | Pending |
| MKH-02 | Phase 7 | Pending |
| MKH-03 | Phase 7 | Pending |
| MKH-04 | Phase 7 | Pending |
| RPT-01 | Phase 7 | Pending |
| RPT-02 | Phase 7 | Pending |
| SET-01 | Phase 7 | Pending |
| PORT-01 | Phase 8 | Pending |
| PORT-02 | Phase 8 | Pending |
| PORT-03 | Phase 8 | Pending |
| PORT-04 | Phase 8 | Pending |

**Coverage:**
- v1 requirements: 109 total
- Mapped to phases: 109
- Unmapped: 0

---
*Requirements defined: 2026-04-02*
*Last updated: 2026-04-02 after roadmap creation*
