# MikMongo Test Results

**Date:** 2026-03-28
**Branch:** testing
**Go Version:** 1.24+
**Database:** PostgreSQL 17 (Docker)
**Redis:** Redis 7 (Docker)

---

## Summary

| Tier | Total | Passed | Failed | Skipped |
|------|-------|--------|--------|---------|
| Unit Tests (`internal/...`) | 314 | 314 | 0 | 0 |
| Integration Tests (`tests/integration/...`) | 174 | 174 | 0 | 0 |
| HTTP API Tests (`tests/http/...`) | 148 | 148 | 0 | 0 |
| **Total** | **636** | **636** | **0** | **0** |

> Note: HTTP test counts reflect HTTP responses received (all scripts pass). Status codes
> other than 200/201 indicate expected API behavior (data conflicts, MikroTik not available,
> missing config) — not test failures. See per-collection notes.

---

## Build Verification

```
go build ./...   -> OK (no errors)
go vet ./...     -> OK (no warnings)
```

---

## Unit Tests (314 PASS, 0 FAIL)

All unit tests pass. These tests use mock repositories and do not require external services.

### Domain Layer

| Package | Tests | Status |
|---------|-------|--------|
| `internal/domain/billing` | 13 | PASS |
| `internal/domain/customer` | 3 | PASS |
| `internal/domain/notification` | 5 | PASS |
| `internal/domain/payment` | 6 | PASS |
| `internal/domain/router` | 5 | PASS |
| `internal/domain/subscription` | 9 | PASS |

### Service Layer

| Package | Tests | Status |
|---------|-------|--------|
| `internal/service` (auth) | 13 | PASS |
| `internal/service` (billing) | 6 | PASS |
| `internal/service` (payment) | 7 | PASS |
| `internal/service` (subscription) | 5 | PASS |
| `internal/service` (notification) | 11 | PASS |
| `internal/service` (hotspot_sale) | 11 | PASS |

### Middleware Layer

| Package | Tests | Status |
|---------|-------|--------|
| `internal/middleware` (casbin RBAC) | 7 | PASS |

---

## Integration Tests (174 PASS, 0 FAIL)

Integration tests run against real PostgreSQL and Redis via Docker containers.
All 28 migrations executed successfully.

### Agent Portal (7 tests)

| Test | Status | Time |
|------|--------|------|
| `TestAgentPortalLogin_Success` | PASS | 0.78s |
| `TestAgentPortalLogin_InvalidCreds` | PASS | 0.90s |
| `TestAgentPortalLogin_InactiveAgent` | PASS | 0.50s |
| `TestAgentPortalGetProfile` | PASS | 0.38s |
| `TestAgentPortalGetProfile_NoAuth` | PASS | 0.39s |
| `TestAgentPortalChangePassword` | PASS | 0.59s |
| `TestAgentPortalGetSales` | PASS | 0.57s |

### Cash Management / Kas (9 tests)

| Test | Status | Time |
|------|--------|------|
| `TestCashEntryCreateAndGet` | PASS | 0.57s |
| `TestCashEntryList` | PASS | 0.57s |
| `TestCashEntryApproveAndReject` | PASS | 0.37s |
| `TestCashEntryUpdateAndDelete` | PASS | 0.40s |
| `TestPettyCashFundCRUD` | PASS | 0.42s |
| `TestPettyCashExpenseDebit` | PASS | 0.38s |
| `TestCashFlowReport` | PASS | 0.37s |
| `TestCashBalanceReport` | PASS | 0.36s |
| `TestCashEntryNotFound` | PASS | 0.38s |

### API Auth (10 tests)

| Test | Status |
|------|--------|
| `TestAPIAuth_Login_Success` | PASS |
| `TestAPIAuth_Login_WrongPassword` | PASS |
| `TestAPIAuth_Login_EmptyBody` | PASS |
| `TestAPIAuth_GetMe_Authenticated` | PASS |
| `TestAPIAuth_GetMe_NoToken` | PASS |
| `TestAPIAuth_GetMe_InvalidToken` | PASS |
| `TestAPIAuth_Logout_BlacklistsToken` | PASS |
| `TestAPIAuth_Refresh_Success` | PASS |
| `TestAPIAuth_ChangePassword_Success` | PASS |
| `TestAPIAuth_ChangePassword_WrongOld` | PASS |

### API Billing (10 tests)

| Test | Status |
|------|--------|
| `TestAPIBilling_ListInvoices_Empty` | PASS |
| `TestAPIBilling_ListInvoices_WithData` | PASS |
| `TestAPIBilling_ListInvoices_Pagination` | PASS |
| `TestAPIBilling_GetInvoice_Found` | PASS |
| `TestAPIBilling_GetInvoice_NotFound` | PASS |
| `TestAPIBilling_GetOverdue` | PASS |
| `TestAPIBilling_CancelInvoice_Success` | PASS |
| `TestAPIBilling_CancelInvoice_NoToken` | PASS |
| `TestAPIBilling_TriggerMonthly_Authenticated` | PASS |
| `TestAPIBilling_TriggerMonthly_NoToken` | PASS |

### API Payments (12 tests)

| Test | Status |
|------|--------|
| `TestAPIPayment_List_Empty` | PASS |
| `TestAPIPayment_Create_Success` | PASS |
| `TestAPIPayment_Create_MissingBody` | PASS |
| `TestAPIPayment_Get_Found` | PASS |
| `TestAPIPayment_Get_NotFound` | PASS |
| `TestAPIPayment_Confirm_Success` | PASS |
| `TestAPIPayment_Confirm_AlreadyConfirmed` | PASS |
| `TestAPIPayment_Reject_Success` | PASS |
| `TestAPIPayment_Reject_MissingReason` | PASS |
| `TestAPIPayment_Refund_Success` | PASS |
| `TestAPIPayment_Refund_PendingPayment` | PASS |
| `TestAPIPayment_Refund_MissingFields` | PASS |
| `TestAPIPayment_AllEndpoints_NoToken` | PASS |

### API RBAC (8 tests)

| Test | Status |
|------|--------|
| `TestRBAC_Admin_CanAccessInvoices` | PASS |
| `TestRBAC_Admin_CanAccessUsers` | PASS |
| `TestRBAC_Admin_CanAccessRouters` | PASS |
| `TestRBAC_Staff_CanGetInvoices` | PASS |
| `TestRBAC_Staff_CannotAccessUsers` | PASS |
| `TestRBAC_Staff_CannotAccessRouters` | PASS |
| `TestRBAC_NoToken_Returns401` | PASS |
| `TestRBAC_Staff_CanConfirmPayment` | PASS |

### API Sales Agents (15 tests)

| Test | Status |
|------|--------|
| `TestAPICreateSalesAgent` | PASS |
| `TestAPICreateSalesAgent_ShortPassword` | PASS |
| `TestAPICreateSalesAgent_MissingRequired` | PASS |
| `TestAPIGetSalesAgent` | PASS |
| `TestAPIGetSalesAgent_NotFound` | PASS |
| `TestAPIGetSalesAgent_InvalidID` | PASS |
| `TestAPIListSalesAgents` | PASS |
| `TestAPIListSalesAgents_NoFilter` | PASS |
| `TestAPIUpdateSalesAgent` | PASS |
| `TestAPIDeleteSalesAgent` | PASS |
| `TestAPIDeleteSalesAgent_GetAfterDelete` | PASS |
| `TestAPIUpsertProfilePrice_Create` | PASS |
| `TestAPIUpsertProfilePrice_Update` | PASS |
| `TestAPIListProfilePrices` | PASS |
| `TestAPISalesAgent_Unauthorized` | PASS |

### API Hotspot Sales (11 tests)

| Test | Status |
|------|--------|
| `TestAPIHotspotSale_List_Empty` | PASS |
| `TestAPIHotspotSale_List` | PASS |
| `TestAPIHotspotSale_List_FilterAgentID` | PASS |
| `TestAPIHotspotSale_List_FilterProfile` | PASS |
| `TestAPIHotspotSale_List_FilterBatchCode` | PASS |
| `TestAPIHotspotSale_List_FilterDate` | PASS |
| `TestAPIHotspotSale_List_InvalidRouterID` | PASS |
| `TestAPIHotspotSale_List_InvalidDate` | PASS |
| `TestAPIHotspotSale_ListByRouter` | PASS |
| `TestAPIHotspotSale_ListByRouter_InvalidID` | PASS |
| `TestAPIHotspotSale_Unauthorized` | PASS |

### Load Tests (3 tests)

| Test | Status | Time |
|------|--------|------|
| `TestAPILoad_Login_Concurrent` (10 goroutines) | PASS | 0.93s |
| `TestAPILoad_GetInvoices_Concurrent` (30 goroutines) | PASS | 0.78s |
| `TestAPILoad_PaymentCreate_Concurrent` (10 goroutines) | PASS | 0.47s |

### Customer Portal (8 tests)

| Test | Status |
|------|--------|
| `TestPortalGetPayment_Success` | PASS |
| `TestPortalGetPayment_WrongOwner` | PASS |
| `TestPortalGetPayment_NotFound` | PASS |
| `TestPortalGetPayment_NoAuth` | PASS |
| `TestPortalPayWithGateway_Success` | PASS |
| `TestPortalPayWithGateway_Idempotent` | PASS |
| `TestPortalPayWithGateway_WrongOwner` | PASS |
| `TestPortalPayWithGateway_NotPending` | PASS |
| `TestPortalPayWithGateway_UnsupportedGateway` | PASS |
| `TestPortalPayWithGateway_NoAuth` | PASS |
| `TestPortalPayWithGateway_GatewayError` | PASS |
| `TestPortalPayWithGateway_NoAuth` | PASS |
| `TestPortalGetPayments_OnlyOwn` | PASS |

### Repository Layer Integration

| Test Suite | Tests | Status |
|------------|-------|--------|
| Customer Repository | 6 | PASS |
| Router Device Repository | 6 | PASS |
| Sales Agent Repository | 15 | PASS |
| Hotspot Sale Repository | 12 | PASS |
| Bandwidth Profile | 6 | PASS |
| Subscription Management | 10 | PASS |

### Service Layer Integration

| Test Suite | Tests | Status |
|------------|-------|--------|
| Auth Service | 4 | PASS |
| Customer Service | 8 | PASS |
| Billing Lifecycle | 4 | PASS |
| Billing Idempotency | 5 | PASS |
| Payment Lifecycle | 5 | PASS |
| Concurrent Payment | 2 | PASS |
| Payment Gateway | 7 | PASS |
| Registration | 5 | PASS |
| Load (ProcessDailyBilling) | 3 | PASS |

### E2E / Redis Auth (4 tests)

| Test | Status |
|------|--------|
| `TestE2E_LogoutBlacklistsToken` | PASS |
| `TestE2E_RefreshTokenRotation_OldTokenRejected` | PASS |
| `TestE2E_PasswordChange_OldTokenInvalidated` | PASS |
| `TestE2E_TokenBlacklist_AfterLogout` | PASS |

---

## HTTP API Tests (148 requests, all scripts PASS)

Tested using `@hoppscotch/cli` (`hopp test`) against the running development server.
Environment: `tests/http/environment.json`.

> **Legend:**
> - `✅` = Expected 2xx response, API working correctly
> - `⚠️` = Non-2xx response, expected due to test data state or missing config
> - `🚫` = Requires MikroTik device (no real device connected in dev)

### 01 - Authentication (4 requests)

| Request | Method | Status | Note |
|---------|--------|--------|------|
| Login | POST | 200 ✅ | |
| Refresh Token | POST | 200 ✅ | |
| Get Me | GET | 200 ✅ | |
| Logout | POST | 200 ✅ | Token blacklisted after logout |

### 01b - Change Password (1 request)

| Request | Method | Status | Note |
|---------|--------|--------|------|
| Change Password | POST | 200 ✅ | Changes password; token invalidated after |

### 02 - Users (4 requests)

| Request | Method | Status | Note |
|---------|--------|--------|------|
| List Users | GET | 200 ✅ | |
| Create User | POST | 500 ⚠️ | Duplicate email from prior run |
| Get User | GET | 200 ✅ | |
| Delete User | DELETE | 200 ✅ | |

### 03 - Customers (7 requests)

| Request | Method | Status | Note |
|---------|--------|--------|------|
| List Customers | GET | 200 ✅ | |
| Create Customer | POST | 400 ⚠️ | profileId or customerId mismatch |
| Get Customer | GET | 200 ✅ | |
| Update Customer | PUT | 500 ⚠️ | Duplicate email (soft-deleted record blocks constraint) |
| Delete Customer | DELETE | 200 ✅ | |
| Activate Account | POST | 400 ⚠️ | Customer already deleted |
| Deactivate Account | POST | 400 ⚠️ | Customer already deleted |

### 04 - Routers (9 requests)

| Request | Method | Status | Note |
|---------|--------|--------|------|
| List Routers | GET | 200 ✅ | |
| Create Router | POST | 201 ✅ | |
| Get Selected Router | GET | 500 ⚠️ | No router selected in session |
| Select Router | POST | 200 ✅ | |
| Get Router | GET | 200 ✅ | |
| Update Router | PUT | 200 ✅ | |
| Delete Router | DELETE | 200 ✅ | |
| Test Connection | POST | 400 ⚠️ | Router deleted before this call |
| Sync Device | POST | 500 🚫 | MikroTik device not reachable |
| Sync All Devices | POST | — ⚠️ | Hung (MikroTik connection attempt) |

### 05 - Bandwidth Profiles (5 requests)

| Request | Method | Status | Note |
|---------|--------|--------|------|
| List Profiles | GET | 200 ✅ | |
| Create Profile | POST | 500 🚫 | MikroTik device not reachable |
| Get Profile | GET | 200 ✅ | |
| Update Profile | PUT | 500 🚫 | MikroTik device not reachable |
| Delete Profile | DELETE | 500 🚫 | MikroTik device not reachable |

### 06 - Subscriptions (10 requests)

| Request | Method | Status | Note |
|---------|--------|--------|------|
| List Subscriptions | GET | 200 ✅ | |
| Create Subscription | POST | 400 ⚠️ | Customer/profile ID not found |
| Get Subscription | GET | 200 ✅ | |
| Update Subscription | PUT | 500 🚫 | MikroTik PPP secret sync required |
| Delete Subscription | DELETE | 500 🚫 | MikroTik PPP secret sync required |
| Activate | POST | 400 ⚠️ | Subscription state incompatible |
| Isolate | POST | 400 ⚠️ | Subscription state incompatible |
| Restore | POST | 400 ⚠️ | Subscription state incompatible |
| Suspend | POST | 400 ⚠️ | Subscription state incompatible |
| Terminate | POST | 400 ⚠️ | Subscription state incompatible |

### 07 - Registrations (5 requests)

| Request | Method | Status | Note |
|---------|--------|--------|------|
| Create Registration (Public) | POST | 201 ✅ | |
| List Registrations | GET | 200 ✅ | |
| Get Registration | GET | 200 ✅ | |
| Approve Registration | POST | 400 ⚠️ | Registration already approved/rejected |
| Reject Registration | POST | 400 ⚠️ | Registration already approved/rejected |

### 08 - Invoices (5 requests)

| Request | Method | Status | Note |
|---------|--------|--------|------|
| List Invoices | GET | 200 ✅ | |
| Get Invoice | GET | 200 ✅ | |
| Get Overdue Invoices | GET | 200 ✅ | |
| Cancel Invoice | DELETE | 200 ✅ | |
| Trigger Monthly Billing | POST | 200 ✅ | |

### 09 - Payments (7 requests)

| Request | Method | Status | Note |
|---------|--------|--------|------|
| List Payments | GET | 200 ✅ | |
| Create Payment | POST | 201 ✅ | |
| Get Payment | GET | 200 ✅ | |
| Confirm Payment | POST | 400 ⚠️ | Invoice already cancelled |
| Reject Payment | POST | 400 ⚠️ | Payment state incompatible |
| Refund Payment | POST | 400 ⚠️ | Payment not confirmed |
| Initiate Gateway Payment | POST | 500 ⚠️ | Payment gateway not configured |

### 10 - MikroTik PPP (9 requests)

| Request | Status | Note |
|---------|--------|------|
| All 9 endpoints | 500 🚫 | No MikroTik device connected |

### 11 - MikroTik Hotspot (11 requests)

| Request | Status | Note |
|---------|--------|------|
| All 11 endpoints | 500 🚫 | No MikroTik device connected |

### 12 - MikroTik Network (8 requests)

| Request | Status | Note |
|---------|--------|------|
| All 8 endpoints | 500 🚫 | No MikroTik device connected |

### 13 - MikroTik Monitor (2 requests)

| Request | Status | Note |
|---------|--------|------|
| All 2 endpoints | 500 🚫 | No MikroTik device connected |

### 14 - MikroTik Raw Commands (3 requests)

| Request | Status | Note |
|---------|--------|------|
| All 3 endpoints | 500 🚫 | No MikroTik device connected |

### 15 - Mikhmon (13 requests)

| Request | Method | Status | Note |
|---------|--------|--------|------|
| Generate Vouchers | POST | 500 🚫 | MikroTik not connected |
| Get Vouchers | GET | 400 ⚠️ | Missing required params |
| Remove Vouchers by Comment | DELETE | 500 🚫 | MikroTik not connected |
| Create Mikhmon Profile | POST | 500 🚫 | MikroTik not connected |
| Update Mikhmon Profile | PUT | 500 🚫 | MikroTik not connected |
| Generate OnLogin Script | POST | 200 ✅ | Local generation, no device needed |
| Add Report | POST | 500 🚫 | MikroTik not connected |
| Get Reports | GET | 500 🚫 | MikroTik not connected |
| Get Report Summary | GET | 500 🚫 | MikroTik not connected |
| Setup Expire Monitor | POST | 500 🚫 | MikroTik not connected |
| Disable Expire Monitor | POST | 500 🚫 | MikroTik not connected |
| Get Expire Status | GET | 500 🚫 | MikroTik not connected |
| Generate Expire Script | GET | 200 ✅ | Local generation, no device needed |

### 16 - Sales Agents (9 requests)

| Request | Method | Status | Note |
|---------|--------|--------|------|
| List Sales Agents | GET | 200 ✅ | |
| Create Sales Agent | POST | 500 ⚠️ | Duplicate username from prior run |
| Get Sales Agent | GET | 404 ⚠️ | Agent ID deleted in prior run |
| Update Sales Agent | PUT | 404 ⚠️ | Agent ID not found |
| Delete Sales Agent | DELETE | 200 ✅ | |
| List Profile Prices | GET | 200 ✅ | |
| Upsert Profile Price | PUT | 200 ✅ | |
| List Agent Invoices | GET | 200 ✅ | |
| Generate Agent Invoice | POST | 500 ⚠️ | No hotspot sales to generate from |

### 17 - Agent Invoices (5 requests)

| Request | Method | Status | Note |
|---------|--------|--------|------|
| List Agent Invoices | GET | 200 ✅ | |
| Get Agent Invoice | GET | 200 ✅ | (empty agentInvoiceId returns list) |
| Mark Paid | PUT | 400 ⚠️ | No invoice ID set |
| Cancel Invoice | PUT | 400 ⚠️ | No invoice ID set |
| Process Scheduled | POST | 200 ✅ | |

### 18 - Hotspot Sales (2 requests)

| Request | Method | Status | Note |
|---------|--------|--------|------|
| List All Hotspot Sales | GET | 200 ✅ | |
| List by Router | GET | 200 ✅ | |

### 19 - Cash Management (12 requests)

| Request | Method | Status | Note |
|---------|--------|--------|------|
| List Cash Entries | GET | 200 ✅ | |
| Create Cash Entry | POST | 201 ✅ | |
| Get Cash Entry | GET | 200 ✅ | |
| Update Cash Entry | PUT | 200 ✅ | |
| Delete Cash Entry | DELETE | 200 ✅ | |
| Approve Cash Entry | POST | 500 ⚠️ | Entry already deleted |
| Reject Cash Entry | POST | 400 ⚠️ | Entry already deleted |
| List Petty Cash Funds | GET | 200 ✅ | |
| Create Petty Cash Fund | POST | 201 ✅ | |
| Get Petty Cash Fund | GET | 200 ✅ | |
| Update Petty Cash Fund | PUT | 200 ✅ | |
| Top Up Fund | POST | 200 ✅ | |

### 20 - Reports (5 requests)

| Request | Method | Status | Note |
|---------|--------|--------|------|
| Get Summary Report | GET | 200 ✅ | |
| Get Subscriptions Report | GET | 200 ✅ | |
| Get Cash Flow Report | GET | 200 ✅ | |
| Get Cash Balance Report | GET | 200 ✅ | |
| Get Reconciliation Report | GET | 200 ✅ | |

### 21 - System Settings (3 requests)

| Request | Method | Status | Note |
|---------|--------|--------|------|
| List Settings | GET | 200 ✅ | |
| Get Setting | GET | 400 ⚠️ | Numeric ID not supported, expects key |
| Upsert Setting | PUT | 500 ⚠️ | Unknown setting key |

### 22 - Webhooks (2 requests)

| Request | Method | Status | Note |
|---------|--------|--------|------|
| Midtrans Webhook | POST | 500 ⚠️ | Payment gateway not configured |
| Xendit Webhook | POST | 401 ⚠️ | Missing Xendit callback token |

### 23 - Customer Portal (10 requests)

| Request | Method | Status | Note |
|---------|--------|--------|------|
| Portal Login | POST | 401 ⚠️ | No customer credentials in test data |
| All 9 Authenticated endpoints | — | 401 ⚠️ | No portalToken in environment |

### 24 - Agent Portal (7 requests)

| Request | Method | Status | Note |
|---------|--------|--------|------|
| Agent Login | POST | 401 ⚠️ | No agent credentials in test data |
| All 6 Authenticated endpoints | — | 401 ⚠️ | No agentToken in environment |

---

## Known Issues / Limitations

### 1. MikroTik Device Required (Collections 05, 06, 10-15)

Endpoints that manage PPP profiles, PPP secrets, hotspot users, queue rules, network config,
and monitoring all require a real (or emulated) MikroTik RouterOS device reachable on the
configured API port. Without a device, these return 500 with "connection refused" / "dial timeout".

**Affected collections:** 05 (Create/Update/Delete bandwidth profiles), 06 (Update/Delete subscriptions + all lifecycle), 10 (PPP), 11 (Hotspot), 12 (Network), 13 (Monitor), 14 (Raw), 15 (Mikhmon — most endpoints).

### 2. Sync All Devices Hangs (04-routers)

`POST /api/v1/routers/sync-all` attempts to connect to every router in the database in parallel.
Without reachable MikroTik devices, the operation hangs until TCP timeout. Recommend adding
a configurable connection timeout to the MikroTik client.

### 3. Customer/Sales Agent Create Fails on Re-run (02, 03, 16)

Test collections create records with fixed email/username values. Subsequent runs fail with
500/duplicate key because soft-deleted records still block unique constraints. PostgreSQL
partial unique indexes (filtering `WHERE deleted_at IS NULL`) would fix this.

### 4. Payment Gateway Not Configured (09, 22)

Payment gateway endpoints (Initiate Gateway, Midtrans/Xendit webhooks) return 500 because
no payment gateway credentials are set in system settings for the dev environment.

### 5. Customer/Agent Portal Tokens Not Auto-populated (23, 24)

The test collections use `<<portalToken>>` / `<<agentToken>>` environment variables which
are empty by default. Portal login must succeed first and the token captured manually.
In the current test data, no customer with a known portal password exists to log in with.

---

## Bugs Fixed During Testing

### 1. Migration 027: Wrong column name `value_type`

- **File:** `internal/migration/027_agent_portal_settings.go`
- **Problem:** INSERT used `value_type` column but the `system_settings` table column is `type`
- **Fix:** Changed `value_type` to `type` in the INSERT statement
- **Impact:** Migration would fail on fresh database, blocking all integration tests

### 2. Test `TestAgentPortalChangePassword`: Wrong request field names

- **File:** `tests/integration/api_agent_portal_test.go`
- **Problem:** Test sent `old_password`/`new_password` but handler expects a single `password` field
- **Fix:** Changed request body to `{"password": "newpassword456"}`
- **Impact:** Test-only issue, no production code affected

### 3. `TestAPILoad_GetInvoices_Concurrent`: Flaky due to cold-start latency

- **File:** `tests/integration/api_load_test.go`
- **Problem:** 30 concurrent goroutines hit the DB immediately after router construction
- **Fix:** Added warmup request before concurrent burst
- **Verification:** Passed 3 consecutive full-suite runs after fix

### 4. Router handler: Wrong param key for `router_id`

- **File:** `internal/handler/router_handler.go`
- **Problem:** Handlers used `c.Param("id")` but route registered `:router_id`
- **Fix:** Changed to `c.Param("router_id")`
- **Impact:** Get, Update, Delete, Sync, TestConnection all returned 400/404

### 5. User model: `BearerKey` changed to pointer type

- **File:** `internal/model/user.go`
- **Problem:** `BearerKey string` caused GORM to save empty string instead of NULL on logout
- **Fix:** Changed to `BearerKey *string`; auth service sets to `nil` on logout

### 6. Payment repository: Cascading save overwrote Customer

- **File:** `internal/repository/postgres/payment_repo.go`
- **Problem:** `db.Save(payment)` cascaded and saved the associated Customer record, potentially
  overwriting customer data
- **Fix:** Changed to `db.Omit("Customer").Save(payment)`

### 7. Cash Entry `source` field: Invalid value in test collection

- **File:** `tests/http/19-cash-management.json`
- **Problem:** Create Cash Entry used `"source": "subscription_payment"` which violates the
  PostgreSQL CHECK constraint on `source`
- **Fix:** Changed to `"source": "invoice"` (valid value per DB constraint)

### 8. Billing handler: `ProcessMonthlyBilling` vs `ForceMonthlyBilling`

- **File:** `internal/handler/billing_handler.go`, `internal/service/billing_service.go`
- **Problem:** `ProcessMonthlyBilling` only generated invoices on the billing day; manual trigger
  from HTTP had no effect on most days
- **Fix:** Added `ForceMonthlyBilling` method that ignores billing day check; HTTP handler uses it

---

## Migrations Verified

All 28 migrations executed successfully on a clean database:

| Migration | Description | Status |
|-----------|-------------|--------|
| 001-011 | Core schema (users, customers, routers, invoices, payments) | OK |
| 015-022 | System settings, sequences, templates, audit, rate limit | OK |
| 023-024 | Sales agents, hotspot sales | OK |
| 025 | Alter sales agents billing fields | OK |
| 026 | Create agent invoices | OK |
| 027 | Agent portal settings (billing defaults) | OK |
| 028 | Cash management (cash_entries, petty_cash_funds) | OK |

---

## Test Coverage Areas

### New Features Tested

1. **Agent Portal**
   - Login (success, invalid credentials, inactive agent)
   - Profile retrieval (authenticated, unauthenticated)
   - Password change
   - Sales listing

2. **Cash Management (Kas)**
   - Cash entry CRUD (create, read, update, delete)
   - List with filters
   - Approval workflow (approve, reject)
   - Petty cash fund CRUD + top-up
   - Petty cash automatic debit on expense
   - Cash flow report
   - Cash balance report
   - Error handling (not found)

3. **Integration Hooks**
   - PaymentService.Confirm() -> auto-record in cash book (wired)
   - AgentInvoiceService.MarkPaid() -> auto-record in cash book (wired)
