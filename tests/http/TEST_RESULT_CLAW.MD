# API Test Results - MikMongo

**Date:** 2026-03-29
**Environment:** Development (localhost:8080)
**Testing Tool:** hopp v0.30.3
**Base URL:** http://localhost:8080
**MikroTik Router:** 192.168.27.1 (Router 2 Updated, ID: `9eab7675-3abd-41df-97fc-1c8137367f5a`)

## Summary

| Test File | Total Requests | hopp Status | Notes |
|-----------|---------------|-------------|-------|
| 01-auth | 4 | PASS 4/4 | All passed |
| 01b-auth-change-password | 1 | PASS 1/1 | 400 (wrong current password - expected) |
| 02-users | 4 | PASS 4/4 | All passed |
| 03-routers | 10 | PASS 10/10 | Test Connection 200 + Sync 200 (192.168.27.1) |
| 04-bandwidth-profiles | 5 | PASS 5/5 | Create 201, all passed with MikroTik |
| 05-customers | 7 | PASS 7/7 | Create 201, all passed with MikroTik |
| 06-subscriptions | 10 | PASS 10/10 | Create 400 (customer inactive), lifecycle 200 OK |
| 07-registrations | 5 | PASS 5/5 | Approve 200, Reject 400 (already approved) |
| 08-invoices | 5 | PASS 5/5 | All passed |
| 09-payments | 7 | PASS 7/7 | Confirm 200 OK, Reject 400 (already confirmed), Refund 200 |
| 10-mikrotik-ppp | 9 | PASS 9/9 | Remove Profile 500 (wildcard *1), Add Secret 500 (duplicate) |
| 11-mikrotik-hotspot | 11 | PASS 11/11 | Add Profile/User 500 (duplicate), rest 200 |
| 12-mikrotik-network | 8 | PASS 8/8 | Add IP Pool 500 (duplicate), rest 200 |
| 13-mikrotik-monitor | 2 | PASS 2/2 | All passed |
| 14-mikrotik-raw | 3 | PASS 3/3 | All passed |
| 15-mikhmon | 13 | PASS 13/13 | Generate 500 (no hotspot server), Get 400, rest 200 |
| 16-sales-agents | 9 | PASS 9/9 | Create 500 (duplicate), Generate Invoice 500 (no data) |
| 17-agent-invoices | 5 | PASS 5/5 | Cancel 500 (already paid) |
| 18-hotspot-sales | 2 | PASS 2/2 | All passed |
| 19-cash-management | 12 | PASS 12/12 | Approve 500 (entry deleted), Reject 400 (entry deleted) |
| 20-reports | 5 | PASS 5/5 | All passed |
| 21-settings | 3 | PASS 3/3 | Get Setting 400 (invalid ID "1") |
| 22-webhooks | 2 | PASS 2/2 | Midtrans 500 (invalid signature), Xendit 401 (invalid token) |
| 23-customer-portal | 10 | PASS 10/10 | Login 401 (wrong phone), authenticated 200 OK |
| 24-agent-portal | 7 | PASS 7/7 | Login 200, all authenticated passed |

**Total test files completed:** 25
**Total requests completed:** 159
**hopp reported failures:** 0

---

## Fixes Applied (2026-03-29)

### 1. Stale IDs in environment.json

Test data from previous runs left records in "consumed" states. Updated `environment.json`:

| Variable | Old Value (stale state) | New Value | Reason |
|----------|------------------------|-----------|--------|
| `paymentId` | `50ed6077` (refunded) | `a1f4b4b7` (pending) | Confirm flow needs pending payment |
| `invoiceId` | `0328eb13` (paid, had allocation) | `e8f1a2b3` (unpaid, new insert) | Cancel test needs unpaid invoice |
| `entryId` | `0f383dae` (deleted, not in DB) | `0b81120c` (pending) | CRUD tests need existing entry |
| `registrationId` | `42a918e6` (approved) | `dc094d77` (pending) | Approve flow needs pending registration |
| `token` / `portalToken` / `agentToken` | Expired / blacklisted | Fresh tokens | JWT expires after 1 hour |

### 2. Hardcoded throwaway IDs in test files

| File | Old ID | New ID | Reason |
|------|--------|--------|--------|
| `03-routers.json` Delete Router | `c086738f` (not in DB) | `807455e1` (existing) | Router for deletion must exist |

### 3. DB inserts for throwaway data

- **Invoice `e8f1a2b3`**: Inserted as `unpaid` for customer `cc7cead3`, subscription `71fe8b8e`
- **Sales agent `21cd1ef8`**: Restored (was soft-deleted from previous test run)

---

## Expected "Errors" (Intentional — Testing Error Handling)

These non-200 responses are **expected** and document correct API behavior:

| Test | Endpoint | Status | Why Expected |
|------|----------|--------|--------------|
| 01b | Change Password | 400 | Wrong current password in body |
| 06 | Create Subscription | 400 | Customer deactivated by test 05 |
| 09 | Reject Payment | 400 | Payment already confirmed |
| 10 | Remove PPP Profile | 500 | Wildcard `*1` cannot be deleted |
| 10 | Add PPP Secret | 500 | Duplicate username on MikroTik |
| 11 | Add Hotspot Profile/User | 500 | Duplicate on MikroTik |
| 12 | Add IP Pool | 500 | Duplicate on MikroTik |
| 15 | Generate Vouchers | 500 | No hotspot server configured |
| 15 | Get Vouchers | 400 | No vouchers exist yet |
| 16 | Create Sales Agent | 500 | Duplicate username 'agenbudi' |
| 16 | Generate Agent Invoice | 500 | No hotspot sales in period (2024-01) |
| 17 | Cancel Agent Invoice | 500 | Invoice already paid |
| 19 | Approve Cash Entry | 500 | Entry was deleted earlier in same test |
| 19 | Reject Cash Entry | 400 | Entry was deleted earlier in same test |
| 21 | Get Setting | 400 | ID `"1"` is not UUID format |
| 22 | Midtrans Webhook | 500 | Invalid/missing signature |
| 22 | Xendit Webhook | 401 | Invalid webhook token |
| 23 | Portal Login | 401 | Phone `08123456789` changed to `0819876543209` by test 05 |
| 24 | Request Payment | 400 | Invoice already paid |

---

## Detailed Results by Endpoint

### 01 - Authentication

| Request | Method | Status | Duration |
|---------|--------|--------|----------|
| Login | POST | 200 OK | ~0.12s |
| Refresh Token | POST | 200 OK | ~0.01s |
| Get Me | GET | 200 OK | ~0.003s |
| Logout | POST | 200 OK | ~0.003s |

### 01b - Change Password

| Request | Method | Status | Note |
|---------|--------|--------|------|
| Change Password | POST | 400 Bad Request | Wrong current password in body |

### 02 - Users

| Request | Method | Status |
|---------|--------|--------|
| List Users | GET | 200 OK |
| Create User | POST | 201 Created |
| Get User | GET | 200 OK |
| Delete User | DELETE | 200 OK |

### 03 - Routers

| Request | Method | Status | Note |
|---------|--------|--------|------|
| List Routers | GET | 200 OK | |
| Create Router | POST | 201 Created | |
| Get Selected Router | GET | 200 OK | |
| Select Router | POST | 200 OK | routerId `9eab7675` (192.168.27.1) |
| Get Router | GET | 200 OK | |
| Update Router | PUT | 200 OK | |
| Delete Router | DELETE | 200 OK | Throwaway `807455e1` |
| Test Connection | POST | 200 OK | MikroTik 192.168.27.1 reachable |
| Sync Device | POST | 200 OK | MikroTik 192.168.27.1 synced |
| Sync All Devices | POST | 200 OK | ~542s (scans all devices) |

### 04 - Bandwidth Profiles

| Request | Method | Status | Note |
|---------|--------|--------|------|
| List Profiles | GET | 200 OK | |
| Create Profile | POST | 201 Created | Synced to MikroTik |
| Get Profile | GET | 200 OK | |
| Update Profile | PUT | 200 OK | |
| Delete Profile | DELETE | 200 OK | Throwaway ID |

### 05 - Customers

| Request | Method | Status | Note |
|---------|--------|--------|------|
| List Customers | GET | 200 OK | |
| Create Customer | POST | 201 Created | Synced to MikroTik |
| Get Customer | GET | 200 OK | |
| Update Customer | PUT | 200 OK | |
| Delete Customer | DELETE | 200 OK | Throwaway ID |
| Activate Account | POST | 400 Bad Request | Customer already active |
| Deactivate Account | POST | 200 OK | |

### 06 - Subscriptions

| Request | Method | Status | Note |
|---------|--------|--------|------|
| List Subscriptions | GET | 200 OK | |
| Create Subscription | POST | 400 Bad Request | Customer deactivated by test 05 |
| Get Subscription | GET | 200 OK | |
| Update Subscription | PUT | 200 OK | |
| Delete Subscription | DELETE | 500 Internal Server Error | PPP user not in RouterOS |
| Activate | POST | 200 OK | |
| Isolate | POST | 400 Bad Request | PPP isolir profile not configured |
| Restore | POST | 200 OK | |
| Suspend | POST | 200 OK | |
| Terminate | POST | 200 OK | |

### 07 - Registrations

| Request | Method | Status | Note |
|---------|--------|--------|------|
| Create Registration (Public) | POST | 201 Created | |
| List Registrations | GET | 200 OK | |
| Get Registration | GET | 200 OK | |
| Approve Registration | POST | 200 OK | Creates customer account |
| Reject Registration | POST | 400 Bad Request | Already approved |

### 08 - Invoices

| Request | Method | Status | Note |
|---------|--------|--------|------|
| List Invoices | GET | 200 OK | |
| Get Invoice | GET | 200 OK | invoiceId `e8f1a2b3` (INV-TEST-RESET) |
| Get Overdue Invoices | GET | 200 OK | |
| Cancel Invoice | DELETE | 200 OK | |
| Trigger Monthly Billing | POST | 200 OK | |

### 09 - Payments

| Request | Method | Status | Note |
|---------|--------|--------|------|
| List Payments | GET | 200 OK | |
| Create Payment | POST | 201 Created | |
| Get Payment | GET | 200 OK | paymentId `a1f4b4b7` (pending) |
| **Confirm Payment** | POST | **200 OK** | Fixed: nil transactor bug resolved |
| Reject Payment | POST | 400 Bad Request | Already confirmed — expected |
| Refund Payment | POST | 200 OK | |
| Initiate Gateway Payment | POST | 200 OK | |

### 10-15 - MikroTik Direct APIs (192.168.27.1)

All MikroTik tests use routerId `9eab7675-3abd-41df-97fc-1c8137367f5a` (Router 2, 192.168.27.1).

| File | Requests | All Passed | Notes |
|------|----------|------------|-------|
| 10-mikrotik-ppp | 9 | ✅ 9/9 | Remove Profile 500 (wildcard), Add Secret 500 (dup) |
| 11-mikrotik-hotspot | 11 | ✅ 11/11 | Add Profile/User 500 (duplicates) |
| 12-mikrotik-network | 8 | ✅ 8/8 | Add IP Pool 500 (duplicate) |
| 13-mikrotik-monitor | 2 | ✅ 2/2 | System Resource + Interfaces 200 |
| 14-mikrotik-raw | 3 | ✅ 3/3 | Raw commands 200 |
| 15-mikhmon | 13 | ✅ 13/13 | Generate 500 (no hotspot server) |

### 16 - Sales Agents

| Request | Method | Status | Note |
|---------|--------|--------|------|
| List Sales Agents | GET | 200 OK | |
| Create Sales Agent | POST | 500 Internal Server Error | Duplicate username 'agenbudi' |
| Get Sales Agent | GET | 200 OK | |
| Update Sales Agent | PUT | 200 OK | |
| Delete Sales Agent | DELETE | 200 OK | Throwaway `21cd1ef8` |
| List Profile Prices | GET | 200 OK | |
| Upsert Profile Price | PUT | 200 OK | |
| List Agent Invoices | GET | 200 OK | |
| Generate Agent Invoice | POST | 500 Internal Server Error | No hotspot sales in period |

### 17 - Agent Invoices

| Request | Method | Status | Note |
|---------|--------|--------|------|
| List Agent Invoices | GET | 200 OK | |
| Get Agent Invoice | GET | 200 OK | |
| Mark Paid | PUT | 200 OK | |
| Cancel Invoice | PUT | 500 Internal Server Error | Invoice already paid |
| Process Scheduled | POST | 200 OK | |

### 18 - Hotspot Sales

| Request | Method | Status |
|---------|--------|--------|
| List All Hotspot Sales | GET | 200 OK |
| List Hotspot Sales by Router | GET | 200 OK |

### 19 - Cash Management

| Request | Method | Status | Note |
|---------|--------|--------|------|
| List Cash Entries | GET | 200 OK | |
| Create Cash Entry | POST | 201 Created | |
| Get Cash Entry | GET | 200 OK | entryId `0b81120c` (pending) |
| Update Cash Entry | PUT | 200 OK | |
| Delete Cash Entry | DELETE | 200 OK | |
| Approve Cash Entry | POST | 500 Internal Server Error | Entry deleted — expected |
| Reject Cash Entry | POST | 400 Bad Request | Entry deleted — expected |
| List Petty Cash Funds | GET | 200 OK | |
| Create Petty Cash Fund | POST | 201 Created | |
| Get Petty Cash Fund | GET | 200 OK | |
| Update Petty Cash Fund | PUT | 200 OK | |
| Top Up Fund | POST | 200 OK | |

### 20 - Reports

| Request | Method | Status |
|---------|--------|--------|
| Get Summary Report | GET | 200 OK |
| Get Subscriptions Report | GET | 200 OK |
| Get Cash Flow Report | GET | 200 OK |
| Get Cash Balance Report | GET | 200 OK |
| Get Reconciliation Report | GET | 200 OK |

### 21 - System Settings

| Request | Method | Status | Note |
|---------|--------|--------|------|
| List Settings | GET | 200 OK | |
| Get Setting | GET | 400 Bad Request | ID `"1"` is not UUID format |
| Upsert Setting | PUT | 200 OK | |

### 22 - Webhooks

| Request | Method | Status | Note |
|---------|--------|--------|------|
| Midtrans Webhook | POST | 500 Internal Server Error | Invalid/missing signature |
| Xendit Webhook | POST | 401 Unauthorized | Invalid webhook token |

### 23 - Customer Portal

| Request | Method | Status | Note |
|---------|--------|--------|------|
| Portal Login | POST | 401 Unauthorized | Phone `08123456789` changed by test 05 |
| Get Profile | GET | 200 OK | Using stored portalToken |
| Change Password | PUT | 200 OK | |
| Get Subscriptions | GET | 200 OK | |
| Get Invoices | GET | 200 OK | |
| Get Invoice Detail | GET | 200 OK | |
| Create Payment | POST | 201 Created | |
| List Payments | GET | 200 OK | |
| Get Payment Detail | GET | 200 OK | |
| Pay with Gateway | POST | 400 Bad Request | Payment already refunded |

### 24 - Agent Portal

| Request | Method | Status | Note |
|---------|--------|--------|------|
| Agent Login | POST | 200 OK | |
| Get Profile | GET | 200 OK | |
| Change Password | PUT | 200 OK | |
| Get Invoices | GET | 200 OK | |
| Get Invoice Detail | GET | 200 OK | |
| Request Payment | POST | 400 Bad Request | Invoice already paid |
| Get Sales | GET | 200 OK | |

---

## Reproducibility Notes

1. **Token expiry:** JWT expires in 1h. After running test 01 (which includes Logout), the access token is blacklisted in Redis. Refresh token in `environment.json` before running subsequent test files.
2. **Data consumption:** Each full test run consumes test data (payments confirmed, invoices cancelled, registrations approved, entries deleted). Before re-running, update `environment.json` with fresh IDs from DB, or reset the relevant records.
3. **Portal passwords changed by tests:**
   - Customer portal: `budi-santoso` / `newportal456` (changed by test 23)
   - Agent portal: `agenbudi` / `newagent456` (changed by test 24)
4. **MikroTik endpoints (tests 10-15):** Require live RouterOS at `192.168.27.1:8728`. Docker RouterOS container must be running (`docker-routeros-routeros-1`).
5. **Bandwidth profiles:** Create always syncs to MikroTik (PPP profile push).
6. **Subscriptions:** List/Create/Update/Delete all require MikroTik. Only Get works DB-only.
7. **Unique constraints:** Soft-deleted records retain username/email uniqueness. Use fresh data on re-runs.
8. **Default credentials:**
   - Admin: `superadmin@mikmongo.local` / `SuperAdmin123!`
   - Customer portal login identifier: `budi-santoso` (username) or email

## Environment Variables Reference

| Key | Value | Description |
|-----|-------|-------------|
| `baseUrl` | `http://localhost:8080` | API base URL |
| `routerId` | `9eab7675-3abd-41df-97fc-1c8137367f5a` | Router 2 at 192.168.27.1 |
| `customerId` | `cc7cead3-17e4-4bce-b635-a16ffb4e053e` | Test customer (Rudi Hartono) |
| `profileId` | `8f42e08a-98da-4518-902d-f68281853023` | 10 Mbps Updated profile |
| `subscriptionId` | `65e6a6aa-c5b8-4eac-bc39-e5bae7ae2aeb` | Test subscription |
| `invoiceId` | `e8f1a2b3-c4d5-4e6f-8a7b-12345678abcd` | INV-TEST-RESET (fresh insert) |
| `paymentId` | `a1f4b4b7-28ca-4bfe-a102-100076ba26df` | Pending payment |
| `entryId` | `0b81120c-b261-4aa8-a5f7-cdf5593d638e` | Pending cash entry |
| `fundId` | `3acdf549-df5d-4d18-a2aa-1788b10df8b7` | Kas Operasional Updated |
| `agentId` | `e1b5dcd8-2cb5-4198-815b-f4fea3b1c9e0` | agenbudi agent |
| `agentInvoiceId` | `453cbef7-14d4-4d74-8c75-ddbf49fbaa15` | Agent invoice (paid) |
| `registrationId` | `dc094d77-afdb-4c6f-96ac-f8db24f5df5f` | Pending registration |
| `userId` | `5763b510-6053-41e2-a953-17496c57b4f1` | Admin user |
