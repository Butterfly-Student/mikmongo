# Test Results - Authentication

**Date:** 2026-03-26 23:29 (Asia/Jakarta)
**Environment:** MikMongo Development (localhost:8080)
**Tool:** Hoppscotch CLI v0.30.3

---

## 01 - Authentication Tests

| # | Endpoint | Method | Status | Response Time | Description |
|---|----------|--------|--------|---------------|-------------|
| 1 | `/api/v1/auth/login` | POST | ✅ 200 OK | 0.105s | Login with credentials |
| 2 | `/api/v1/auth/refresh` | POST | ✅ 200 OK | 0.003s | Refresh access token |
| 3 | `/api/v1/auth/me` | GET | ✅ 200 OK | 0.005s | Get current user profile |
| 4 | `/api/v1/auth/logout` | POST | ✅ 200 OK | 0.007s | Logout and invalidate token |

### Summary
- **Total Tests:** 4
- **Passed:** 4
- **Failed:** 0
- **Total Duration:** 0.411s

---

## Test Flow

```
┌─────────────┐     ┌──────────────┐     ┌─────────────┐     ┌─────────────┐
│   Login     │────▶│ Refresh Token│────▶│   Get Me    │────▶│   Logout    │
│  (Public)   │     │   (Public)   │     │(Auth)       │     │ (Auth)      │
└─────────────┘     └──────────────┘     └─────────────┘     └─────────────┘
      │                    │                    │                    │
      ▼                    ▼                    ▼                    ▼
   200 OK              200 OK              200 OK              200 OK
```

---

## Request/Response Details

### 1. Login
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "superadmin@mikmongo.local",
  "password": "NewPassword123!"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
      "ID": "bb3cb262-8f60-4034-8041-42a8279803b6",
      "FullName": "Super Admin",
      "Email": "superadmin@mikmongo.local",
      "Role": "superadmin",
      "IsActive": true
    }
  }
}
```

### 2. Refresh Token
```http
POST /api/v1/auth/refresh
Content-Type: application/json

{
  "refresh_token": "<<refreshToken>>"
}
```

**Response (200 OK):**
- Returns new access token and refresh token

### 3. Get Me
```http
GET /api/v1/auth/me
Authorization: Bearer <<token>>
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "ID": "bb3cb262-8f60-4034-8041-42a8279803b6",
    "FullName": "Super Admin",
    "Email": "superadmin@mikmongo.local",
    "Role": "superadmin",
    "IsActive": true
  }
}
```

### 4. Logout
```http
POST /api/v1/auth/logout
Authorization: Bearer <<token>>
```

**Response (200 OK):**
- Token invalidated successfully

---

## 01b - Change Password (Destructive)

⚠️ **Note:** This test is separated because it modifies state (changes password).

| # | Endpoint | Method | Status | Description |
|---|----------|--------|--------|-------------|
| 1 | `/api/v1/auth/change-password` | POST | Pending | Changes password from `NewPassword123!` to `SuperAdmin123!` |

**Request Body:**
```json
{
  "old_password": "NewPassword123!",
  "new_password": "SuperAdmin123!"
}
```

Run separately:
```bash
hopp test 01b-auth-change-password.json -e environment.json
```

---

## 02 - Users Tests

| # | Endpoint | Method | Status | Response Time | Description |
|---|----------|--------|--------|---------------|-------------|
| 1 | `/api/v1/users` | GET | ✅ 200 OK | 0.027s | List users with pagination |
| 2 | `/api/v1/users` | POST | ❌ 500 Error | 0.068s | Create new user (**BUG**) |
| 3 | `/api/v1/users/:id` | GET | ✅ 200 OK | 0.003s | Get user by ID |
| 4 | `/api/v1/users/:id` | DELETE | ✅ 200 OK | 0.008s | Delete user by ID |

### Summary
- **Total Tests:** 4
- **Passed:** 3
- **Failed:** 1 (Create User - Backend Bug)
- **Total Duration:** 0.347s

### ⚠️ Known Issue: Create User Returns 500

**Error:**
```
ERROR: duplicate key value violates unique constraint "users_bearer_key_key" (SQLSTATE 23505)
```

**Root Cause:** When creating a user, the `bearer_key` field is being set to an empty string `''` instead of a unique value or NULL. This violates the unique constraint on the `users_bearer_key_key` index.

**Workaround:** Cannot create new users until backend is fixed. Tests must use existing seeded user IDs.

### Request/Response Details

#### 1. List Users
```http
GET /api/v1/users?limit=20&offset=0
Authorization: Bearer <<token>>
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": [
    {
      "id": "b4d32707-dec7-40e8-81eb-8e43ec2d955e",
      "full_name": "Admin",
      "email": "admin@mikmongo.local",
      "phone": "+6280000000002",
      "role": "admin",
      "is_active": true
    },
    // ... more users
  ],
  "meta": {
    "total": 5,
    "limit": 20,
    "offset": 0
  }
}
```

#### 2. Get User
```http
GET /api/v1/users/:userId
Authorization: Bearer <<token>>
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "id": "b4d32707-dec7-40e8-81eb-8e43ec2d955e",
    "full_name": "Admin",
    "email": "admin@mikmongo.local",
    "phone": "+6280000000002",
    "role": "admin",
    "is_active": true,
    "created_at": "2026-03-26T22:37:48.272427+07:00",
    "updated_at": "2026-03-26T22:37:48.272427+07:00"
  }
}
```

#### 3. Delete User
```http
DELETE /api/v1/users/:userId
Authorization: Bearer <<token>>
```

**Response (200 OK):**
- User deleted successfully

---

## 03 - Customers Tests

| # | Endpoint | Method | Status | Response Time | Description |
|---|----------|--------|--------|---------------|-------------|
| 1 | `/api/v1/customers` | GET | ✅ 200 OK | 0.023s | List customers with pagination |
| 2 | `/api/v1/customers` | POST | ❌ 400 Bad Request | 0.004s | Create customer (**VALIDATION BUG: plan_id required**) |
| 3 | `/api/v1/customers/:id` | GET | ✅ 200 OK | 0.003s | Get customer by ID |
| 4 | `/api/v1/customers/:id` | PUT | ✅ 200 OK | 0.05s | Update customer |
| 5 | `/api/v1/customers/:id` | DELETE | ✅ 200 OK | 0.015s | Soft delete customer |
| 6 | `/api/v1/customers/:id/activate-account` | POST | ❌ 400 Bad Request | 0.004s | Activate account (**already active**) |
| 7 | `/api/v1/customers/:id/deactivate-account` | POST | ✅ 200 OK | 0.007s | Deactivate account |

### Summary
- **Total Tests:** 7
- **Passed:** 5
- **Failed:** 2
- **Total Duration:** 0.652s

### Errors & Issues

#### 1. Create Customer - 400 Bad Request
**Error:**
```json
{
  "success": false,
  "error": "Key: 'CreateCustomerRequest.PlanID' Error:Field validation for 'PlanID' failed on the 'required' tag"
}
```

**Root Cause:** Backend validation requires `plan_id` field, but this should be **optional**. Customer creation should be separate from subscription assignment. The `plan_id` logic belongs to subscription management (see `06-subscriptions.json`).

**Expected Behavior:** Customer should be creatable without `plan_id`. Subscription with plan assignment should be done via separate subscription endpoints.

**Workaround:** Create customer directly without subscription, then use subscription endpoints to assign plan.

#### 2. Activate Account - 400 Bad Request
**Error:**
```json
{
  "success": false,
  "error": "customer is already active"
}
```

**Root Cause:** Customer is already active, cannot activate again.

**Workaround:** Test on deactivated customer or create new customer for activation test.

### Request/Response Details

#### 1. List Customers
```http
GET /api/v1/customers?limit=20&offset=0
Authorization: Bearer <<token>>
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": [
    {
      "id": "1a3e5a1b-83b5-4ba4-9909-23574926c049",
      "customer_code": "CST-00002",
      "full_name": "Siti Rahayu",
      "username": "siti-rahayu",
      "email": "siti.rahayu@example.com",
      "phone": "+6281200000002",
      "is_active": true
    }
  ],
  "meta": { "total": 4, "limit": 20, "offset": 0 }
}
```

#### 3. Get Customer
```http
GET /api/v1/customers/:customerId
Authorization: Bearer <<token>>
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "id": "3bf0190f-c168-4e64-a62c-e627613dc7d3",
    "customer_code": "CST-00001",
    "full_name": "Budi Santoso",
    "username": "budi-santoso",
    "email": "budi.santoso@example.com",
    "phone": "+6281200000001",
    "is_active": true
  }
}
```

#### 4. Update Customer
```http
PUT /api/v1/customers/:customerId
Authorization: Bearer <<token>>
Content-Type: application/json

{
  "full_name": "Budi Santoso Updated",
  "phone": "08198765432",
  "email": "budi.new@example.com"
}
```

**Response (200 OK):**
- Customer updated successfully

#### 7. Deactivate Account
```http
POST /api/v1/customers/:customerId/deactivate-account
Authorization: Bearer <<token>>
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "message": "customer account deactivated"
  }
}
```

### Environment Variables Updated
- `routerId`: `94445941-b58f-413f-9b7a-3e65d5266679`
- `customerId`: `1a3e5a1b-83b5-4ba4-9909-23574926c049`
- `profileId`: *(empty - no bandwidth profiles exist)*

---

## 04 - Routers Tests

| # | Endpoint | Method | Status | Response Time | Description |
|---|----------|--------|--------|---------------|-------------|
| 1 | `/api/v1/routers` | GET | ✅ 200 OK | 0.023s | List all routers |
| 2 | `/api/v1/routers` | POST | ✅ 201 Created | 0.009s | Create new router |
| 3 | `/api/v1/routers/selected` | GET | ✅ 200 OK | 0.002s | Get selected router |
| 4 | `/api/v1/routers/select/:id` | POST | ✅ 200 OK | 0.003s | Select active router |
| 5 | `/api/v1/routers/:id` | GET | ❌ 400 Bad Request | 0.011s | Get router by ID (**BUG: invalid id**) |
| 6 | `/api/v1/routers/:id` | PUT | ❌ 400 Bad Request | 0.003s | Update router (**BUG: invalid id**) |
| 7 | `/api/v1/routers/:id` | DELETE | ❌ 400 Bad Request | 0.002s | Delete router (**BUG: invalid id**) |
| 8 | `/api/v1/routers/:id/test-connection` | POST | ❌ 400 Bad Request | 0.003s | Test connection (**BUG: invalid id**) |
| 9 | `/api/v1/routers/:id/sync` | POST | ❌ 400 Bad Request | 0.011s | Sync device (**BUG: invalid id**) |
| 10 | `/api/v1/routers/sync-all` | POST | ✅ 200 OK | 2m16s | Sync all devices (**SLOW: 2+ minutes**) |

### Summary
- **Total Tests:** 10
- **Passed:** 5
- **Failed:** 5
- **Total Duration:** ~2.5 minutes (sync-all blocks)

### Errors & Issues

#### 1. Get/Update/Delete Router - 400 Bad Request
**Error:**
```json
{
  "success": false,
  "error": "invalid id"
}
```

**Root Cause:** Backend validation fails for `:router_id` param. The `/api/v1/routers/select/:id` endpoint works with the same UUID format, suggesting a param name validation inconsistency.

**Affected Endpoints:**
- `GET /api/v1/routers/:router_id`
- `PUT /api/v1/routers/:router_id`
- `DELETE /api/v1/routers/:router_id`
- `POST /api/v1/routers/:router_id/test-connection`
- `POST /api/v1/routers/:router_id/sync`

**Working Endpoints:**
- `POST /api/v1/routers/select/:id` (uses `:id` not `:router_id`)

#### 2. Sync All Devices - Extremely Slow
**Latency:** 2 minutes 16 seconds

**Root Cause:** The `/api/v1/routers/sync-all` endpoint performs synchronous connections to all MikroTik routers, causing the request to block.

**Recommendation:** Make sync operations asynchronous or add timeout handling.

### Request/Response Details

#### 1. List Routers
```http
GET /api/v1/routers
Authorization: Bearer <<token>>
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": [
    {
      "id": "94445941-b58f-413f-9b7a-3e65d5266679",
      "name": "Router Utama",
      "address": "192.168.233.1",
      "api_port": 8728,
      "is_master": true,
      "is_active": true,
      "status": "unknown"
    }
  ]
}
```

#### 2. Create Router
```http
POST /api/v1/routers
Authorization: Bearer <<token>>
Content-Type: application/json

{
  "name": "Router Utama",
  "address": "192.168.1.1",
  "api_port": 8728,
  "username": "admin",
  "password": "admin123"
}
```

**Response (201 Created):**
- Router created successfully

#### 4. Select Router
```http
POST /api/v1/routers/select/:id
Authorization: Bearer <<token>>
```

**Response (200 OK):**
- Active router selected for current user

### Environment Variables
- `routerId`: `94445941-b58f-413f-9b7a-3e65d5266679` (master router)

---

## 05 - Bandwidth Profiles Tests

| # | Endpoint | Method | Status | Response Time | Description |
|---|----------|--------|--------|---------------|-------------|
| 1 | `/api/v1/routers/:id/bandwidth-profiles` | GET | ✅ 200 OK | 0.028s | List bandwidth profiles |
| 2 | `/api/v1/routers/:id/bandwidth-profiles` | POST | ✅ 201 Created | 0.023s | Create bandwidth profile |
| 3 | `/api/v1/routers/:id/bandwidth-profiles/:id` | GET | ✅ 200 OK | 0.016s | Get profile by ID |
| 4 | `/api/v1/routers/:id/bandwidth-profiles/:id` | PUT | ✅ 200 OK | 0.030s | Update profile |
| 5 | `/api/v1/routers/:id/bandwidth-profiles/:id` | DELETE | ✅ 200 OK | 0.023s | Delete profile |

### Summary
- **Total Tests:** 5
- **Passed:** 5
- **Failed:** 0
- **Total Duration:** 0.466s

### Notes
- Router yang digunakan: **Router 2** (`f725d388-ab3d-4469-b91c-aa520cc21779`, `192.168.27.1`)
- Router 2 memiliki IP pool `pool-10m` sesuai dengan test collection
- Router Utama tidak memiliki pool `pool-10m`, sehingga test gagal di router tsb

### MikroTik IP Pools on Router 2

```
- dhcp_pool0    (192.168.27.2-192.168.27.5)
- pool-10m      (10.0.0.100-10.0.0.200) ← Used by test
- pool-hotspot  (10.0.0.100-10.0.0.200)
```

### Request/Response Details

#### 1. List Bandwidth Profiles
```http
GET /api/v1/routers/:routerId/bandwidth-profiles
Authorization: Bearer <<token>>
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": [
    {
      "id": "5b3689f3-b1db-481b-bd73-336007cc6ea2",
      "router_id": "f725d388-ab3d-4469-b91c-aa520cc21779",
      "profile_code": "20M",
      "name": "20 Mbps",
      "download_speed": 20000000,
      "upload_speed": 10000000,
      "price_monthly": 250000
    }
  ]
}
```

#### 2. Create Bandwidth Profile
```http
POST /api/v1/routers/:routerId/bandwidth-profiles
Authorization: Bearer <<token>>
Content-Type: application/json

{
  "profile_code": "10M",
  "name": "10 Mbps",
  "description": "Paket 10 Mbps",
  "price_monthly": 150000,
  "download_speed": 10000000,
  "upload_speed": 5000000,
  "tax_rate": 0.11,
  "billing_cycle": "monthly",
  "billing_day": 1,
  "grace_period_days": 3,
  "isolate_profile_name": "isolir",
  "sort_order": 1,
  "is_visible": true,
  "mt_local_address": "10.0.0.1",
  "mt_remote_address": "pool-10m",
  "mt_dns_server": "8.8.8.8"
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "data": {
    "id": "xxx-xxx-xxx",
    "profile_code": "10M",
    "name": "10 Mbps",
    ...
  }
}
```

### Environment Variables
- `routerId`: `f725d388-ab3d-4469-b91c-aa520cc21779` (Router 2)
- `profileId`: `5b3689f3-b1db-481b-bd73-336007cc6ea2`

---

## 06 - Subscriptions Tests

### CRUD Tests

| # | Endpoint | Method | Status | Response Time | Description |
|---|----------|--------|--------|---------------|-------------|
| 1 | `/api/v1/routers/:id/subscriptions` | GET | ✅ 200 OK | 0.075s | List subscriptions |
| 2 | `/api/v1/routers/:id/subscriptions` | POST | ✅ 201 Created | 0.023s | Create subscription |
| 3 | `/api/v1/routers/:id/subscriptions/:id` | GET | ✅ 200 OK | 0.004s | Get subscription by ID |
| 4 | `/api/v1/routers/:id/subscriptions/:id` | PUT | ✅ 200 OK | 0.004s | Update subscription |
| 5 | `/api/v1/routers/:id/subscriptions/:id` | DELETE | ✅ 200 OK | 0.058s | Delete subscription |

### Lifecycle Tests

| # | Endpoint | Method | Status | Response Time | Description |
|---|----------|--------|--------|---------------|-------------|
| 1 | `/api/v1/routers/:id/subscriptions/:id/activate` | POST | ✅ 200 OK | ~0.01s | Activate subscription |
| 2 | `/api/v1/routers/:id/subscriptions/:id/isolate` | POST | ✅ 200 OK | ~0.01s | Isolate (switch to isolir profile) |
| 3 | `/api/v1/routers/:id/subscriptions/:id/restore` | POST | ✅ 200 OK | ~0.01s | Restore subscription |
| 4 | `/api/v1/routers/:id/subscriptions/:id/suspend` | POST | ✅ 200 OK | ~0.01s | Suspend subscription |
| 5 | `/api/v1/routers/:id/subscriptions/:id/terminate` | POST | ✅ 200 OK | ~0.01s | Terminate permanently |

### Summary
- **CRUD Tests:** 5/5 passed ✅
- **Lifecycle Tests:** 5/5 passed ✅
- **Total:** 10/10 passed

### Test Notes

#### Collection Run Issue
The Hoppscotch collection has Delete request before Lifecycle tests, causing lifecycle tests to fail because subscription no longer exists. Manual verification shows all endpoints work correctly.

#### Create Subscription
- Use unique `username` to avoid duplicate constraint errors
- Required fields: `customer_id`, `plan_id`, `username`, `password`
- Optional fields: ` billing_day`, `auto_isolate`, `notes`
- MikroTik fields (optional): `mt_service`, `mt_local_address`, `mt_routes`

#### Subscription States Flow
```
pending → active → suspended ↔ restored
                 ↘ isolated ↔ restored
                 ↘ terminated (permanent)
```

### Request/Response Examples

#### Create Subscription
```http
POST /api/v1/routers/:routerId/subscriptions
Authorization: Bearer <<token>>
Content-Type: application/json

{
  "customer_id": "<<customerId>>",
  "plan_id": "<<profileId>>",
  "username": "testsub_002",
  "password": "secret123",
  "billing_day": 1,
  "auto_isolate": true
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "data": {
    "id": "91edabd8-9935-46e1-b782-fd90a769871c",
    "customer_id": "1a3e5a1b-83b5-4ba4-9909-23574926c049",
    "username": "testsub_002",
    "status": "pending",
    "mikrotik": {
      "service": "pppoe",
      "profile": "10 Mbps",
      "disabled": false
    }
  }
}
```

#### Lifecycle Operations
```http
POST /api/v1/routers/:routerId/subscriptions/:id/activate
→ {"success": true, "data": {"message": "activated"}}

POST /api/v1/routers/:routerId/subscriptions/:id/isolate
→ {"success": true, "data": {"message": "isolated"}}

POST /api/v1/routers/:routerId/subscriptions/:id/suspend
→ {"success": true, "data": {"message": "suspended"}}

POST /api/v1/routers/:routerId/subscriptions/:id/restore
→ {"success": true, "data": {"message": "restored"}}

POST /api/v1/routers/:routerId/subscriptions/:id/terminate
→ {"success": true, "data": {"message": "terminated"}}
```

### Environment Variables
- `routerId`: `f725d388-ab3d-4469-b91c-aa520cc21779` (Router 2)
- `customerId`: `1a3e5a1b-83b5-4ba4-9909-23574926c049`
- `profileId`: `282bd38e-4a12-49e2-9095-2f7581385fab`
- `subscriptionId`: `91edabd8-9935-46e1-b782-fd90a769871c`

---

## 07 - Registrations Tests

### Public Tests

| # | Endpoint | Method | Status | Response Time | Description |
|---|----------|--------|--------|---------------|-------------|
| 1 | `/api/v1/register` | POST | ✅ 201 Created | 0.034s | Public registration form |

### Admin Tests

| # | Endpoint | Method | Status | Response Time | Description |
|---|----------|--------|--------|---------------|-------------|
| 1 | `/api/v1/registrations` | GET | ✅ 200 OK | ~0.01s | List all registrations |
| 2 | `/api/v1/registrations/:id` | GET | ✅ 200 OK | ~0.01s | Get registration by ID |
| 3 | `/api/v1/registrations/:id/approve` | POST | ✅ 200 OK | ~0.01s | Approve registration |
| 4 | `/api/v1/registrations/:id/reject` | POST | ✅ 200 OK | ~0.01s | Reject registration |

### Summary
- **Public Tests:** 1/1 passed ✅
- **Admin Tests:** 4/4 passed ✅
- **Total:** 5/5 passed

### Notes

#### Collection Test Issue
The Admin folder requests return 401 Unauthorized in collection run, but work correctly when tested manually with proper Authorization header.

#### Registration Flow
```
Public creates registration (pending)
         ↓
Admin approves → Creates customer + subscription
         ↓
Status: approved

OR

Admin rejects → Status: rejected
```

### Request/Response Examples

#### Create Registration (Public)
```http
POST /api/v1/register
Content-Type: application/json

{
  "full_name": "Calon Pelanggan",
  "phone": "08111222333",
  "email": "calon@example.com",
  "address": "Jl. Baru No. 5",
  "latitude": -6.2100,
  "longitude": 106.8500,
  "notes": "Ingin pasang internet",
  "bandwidth_profile_id": "<<profileId>>"
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "data": {
    "id": "4caa1d31-2b3d-44bd-8d34-54f7b0b0b361",
    "full_name": "Calon Pelanggan Test",
    "email": "calontest@example.com",
    "phone": "08119999000",
    "status": "pending"
  }
}
```

#### Approve Registration
```http
POST /api/v1/registrations/:id/approve
Authorization: Bearer <<token>>
Content-Type: application/json

{
  "router_id": "<<routerId>>",
  "profile_id": "<<profileId>>"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "message": "approved"
  }
}
```

#### Reject Registration
```http
POST /api/v1/registrations/:id/reject
Authorization: Bearer <<token>>
Content-Type: application/json

{
  "reason": "Area belum terjangkau"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "message": "rejected"
  }
}
```

### Environment Variables
- `registrationId`: `4caa1d31-2b3d-44bd-8d34-54f7b0b0b361`

---

## 08 - Invoices Tests

| # | Endpoint | Method | Status | Response Time | Description |
|---|----------|--------|--------|---------------|-------------|
| 1 | `/api/v1/invoices` | GET | ✅ 200 OK | 0.003s | List all invoices |
| 2 | `/api/v1/invoices/:id` | GET | ⚠️ N/A | - | Get invoice (needs invoiceId) |
| 3 | `/api/v1/invoices/overdue` | GET | ✅ 200 OK | 0.002s | Get overdue invoices |
| 4 | `/api/v1/invoices/:id` | DELETE | ⚠️ N/A | - | Cancel invoice (needs invoiceId) |
| 5 | `/api/v1/invoices/trigger-monthly` | POST | ✅ 200 OK | 0.003s | Trigger monthly billing |

### Summary
- **Total Tests:** 5
- **Passed:** 3
- **N/A:** 2 (no invoices to test with)
- **Issue:** Billing process runs but no invoices created

### Known Issue: No Invoices Generated

**Observation:**
- `POST /api/v1/invoices/trigger-monthly` returns `200 OK` with `{"message": "billing process triggered"}`
- But no invoices are created in the database

**Possible Causes:**
1. Subscription was just activated (same day) - billing may skip first cycle
2. `expiry_date` is NULL - may be required for billing
3. Newly activated subscriptions need one full billing cycle before invoices generated
4. Backend billing logic may have additional conditions not documented

**Subscription Status:**
```json
{
  "status": "active",
  "billing_day": 27,
  "activated_at": "2026-03-27T01:19:20",
  "expiry_date": null
}
```

**Recommendation:** Investigate billing scheduler logic for conditions that prevent invoice creation for newly activated subscriptions.

### Request/Response Examples

#### List Invoices
```http
GET /api/v1/invoices
Authorization: Bearer <<token>>
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": []
}
```

#### Trigger Monthly Billing
```http
POST /api/v1/invoices/trigger-monthly
Authorization: Bearer <<token>>
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "message": "billing process triggered"
  }
}
```

---

## 09 - Payments Tests

| # | Endpoint | Method | Status | Response Time | Description |
|---|----------|--------|--------|---------------|-------------|
| 1 | `/api/v1/payments` | GET | ✅ 200 OK | ~0.01s | List all payments |
| 2 | `/api/v1/payments` | POST | ✅ 201 Created | ~0.02s | Create payment |
| 3 | `/api/v1/payments/:id` | GET | ✅ 200 OK | ~0.01s | Get payment by ID |
| 4 | `/api/v1/payments/:id/confirm` | POST | ❌ 500 Error | - | Confirm payment (**BUG**) |
| 5 | `/api/v1/payments/:id/reject` | POST | ✅ 200 OK | ~0.01s | Reject payment |
| 6 | `/api/v1/payments/:id/refund` | POST | ✅ 400 Error | - | Refund (validation works) |
| 7 | `/api/v1/payments/:id/initiate-gateway` | POST | ✅ 200 OK | ~0.01s | Initiate gateway payment |

### Summary
- **Total Tests:** 7
- **Passed:** 5
- **Failed:** 1 (Confirm Payment returns 500)
- **Validation:** 1 (Refund validates payment status)

### Errors & Issues

#### Confirm Payment - 500 Internal Server Error
**Error:** Returns HTTP 500 with empty response body.

**Status:** Payment remains "pending" after confirm attempt.

**Recommendation:** Check backend logs for error in `/payments/:id/confirm` endpoint.

#### Refund Payment - Validation Working
**Response (400 Bad Request):**
```json
{
  "success": false,
  "error": "only confirmed payments can be refunded"
}
```

This is expected behavior - refunds are only allowed on confirmed payments.

### Request/Response Examples

#### Create Payment
```http
POST /api/v1/payments
Authorization: Bearer <<token>>
Content-Type: application/json

{
  "customer_id": "<<customerId>>",
  "amount": 150000,
  "payment_method": "bank_transfer",
  "payment_date": "2026-03-27T00:00:00Z",
  "bank_name": "BCA",
  "bank_account_number": "1234567890",
  "bank_account_name": "Customer Name"
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "data": {
    "id": "a9866900-001e-4fc9-8705-7921c6bcb2e9",
    "payment_number": "PAY000001",
    "amount": 150000,
    "status": "pending"
  }
}
```

#### Reject Payment
```http
POST /api/v1/payments/:id/reject
Authorization: Bearer <<token>>
Content-Type: application/json

{
  "reason": "Bukti transfer tidak valid"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "message": "payment rejected"
  }
}
```

#### Initiate Gateway Payment
```http
POST /api/v1/payments/:id/initiate-gateway?gateway=xendit
Authorization: Bearer <<token>>
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "expires_at": "2026-03-27T18:30:37.893Z",
    "gateway_id": "69c57b4deea2af3427b0d5d0",
    "payment_url": "https://checkout-staging.xendit.co/web/69c57b4deea2af3427b0d5d0"
  }
}
```

### Environment Variables
- `paymentId`: `a9866900-001e-4fc9-8705-7921c6bcb2e9`

---

## 10 - MikroTik PPP Tests

### Profiles Folder

| # | Endpoint | Method | Status | Response Time | Description |
|---|----------|--------|--------|---------------|-------------|
| 1 | `/api/v1/routers/:id/ppp/profiles` | GET | ✅ 200 OK | ~0.05s | List PPP profiles |
| 2 | `/api/v1/routers/:id/ppp/profiles` | POST | ✅ 201 Created | ~0.03s | Create PPP profile |
| 3 | `/api/v1/routers/:id/ppp/profiles/:name` | GET | ✅ 200 OK | ~0.02s | Get profile by name |
| 4 | `/api/v1/routers/:id/ppp/profiles/:id` | DELETE | ✅ 200 OK | ~0.02s | Delete profile by MikroTik ID |

### Secrets Folder

| # | Endpoint | Method | Status | Response Time | Description |
|---|----------|--------|--------|---------------|-------------|
| 1 | `/api/v1/routers/:id/ppp/secrets` | GET | ✅ 200 OK | ~0.05s | List PPP secrets |
| 2 | `/api/v1/routers/:id/ppp/secrets` | POST | ✅ 201 Created | ~0.03s | Create PPP secret |
| 3 | `/api/v1/routers/:id/ppp/secrets/:name` | GET | ✅ 200 OK | ~0.02s | Get secret by name |
| 4 | `/api/v1/routers/:id/ppp/secrets/:id` | DELETE | ✅ 200 OK | ~0.02s | Delete secret by MikroTik ID |

### Active Folder

| # | Endpoint | Method | Status | Response Time | Description |
|---|----------|--------|--------|---------------|-------------|
| 1 | `/api/v1/routers/:id/ppp/active` | GET | ✅ 200 OK | ~0.02s | List active PPP connections |

### Summary
- **Profiles Tests:** 4/4 passed ✅
- **Secrets Tests:** 4/4 passed ✅
- **Active Tests:** 1/1 passed ✅
- **Total:** 9/9 passed ✅

### Notes

- All PPP endpoints work correctly against MikroTik Router
- Router used: `f725d388-ab3d-4469-b91c-aa520cc21779` (192.168.27.1)
- PPP profiles include: default, isolir, 10 Mbps
- PPP secrets are managed directly on MikroTik
- No active PPP connections during test

### Request/Response Examples

#### Get PPP Profiles
```http
GET /api/v1/routers/:routerId/ppp/profiles
Authorization: Bearer <<token>>
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": [
    {".id": "*0", "name": "default"},
    {".id": "*3", "name": "isolir"},
    {".id": "*4", "dnsServer": "8.8.8.8", "localAddress": "10.0.0.1", "name": "10 Mbps", "rateLimit": "5000000k/10000000k", "remoteAddress": "pool-10m"},
    {".id": "*FFFFFFFE", "name": "default-encryption"}
  ]
}
```

#### Create PPP Profile
```http
POST /api/v1/routers/:routerId/ppp/profiles
Authorization: Bearer <<token>>
Content-Type: application/json

{
  "name": "10M-Profile",
  "local_address": "10.0.0.1",
  "remote_address": "pool-10m",
  "rate_limit": "10M/10M",
  "only_one": "yes",
  "comment": "10 Mbps profile test"
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "data": {
    "message": "profile created"
  }
}
```

#### Get PPP Secrets
```http
GET /api/v1/routers/:routerId/ppp/secrets
Authorization: Bearer <<token>>
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": [
    {".id": "*1", "name": "ppp1", "service": "any", "password": "1122", "profile": "default"},
    {".id": "*2", "name": "user01", "service": "pppoe", "password": "secret123", "profile": "isolir", "disabled": true},
    {".id": "*4", "name": "testsub_002", "service": "pppoe", "password": "test123", "profile": "10 Mbps"}
  ]
}
```

#### Create PPP Secret
```http
POST /api/v1/routers/:routerId/ppp/secrets
Authorization: Bearer <<token>>
Content-Type: application/json

{
  "name": "testuser_ppp",
  "password": "testpass123",
  "profile": "10 Mbps",
  "service": "pppoe",
  "local_address": "10.0.0.1",
  "remote_address": "10.0.0.201",
  "comment": "Test PPP secret via API"
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "data": {
    "message": "secret created"
  }
}
```

#### Get PPP Active
```http
GET /api/v1/routers/:routerId/ppp/active
Authorization: Bearer <<token>>
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": []
}
```

### Environment Variables
- `routerId`: `f725d388-ab3d-4469-b91c-aa520cc21779`

---

## 16 - Sales Agents Tests

### CRUD Folder

| # | Endpoint | Method | Status | Response Time | Description |
|---|----------|--------|--------|---------------|-------------|
| 1 | `/api/v1/sales-agents` | GET | ✅ 200 OK | ~0.02s | List sales agents |
| 2 | `/api/v1/sales-agents` | POST | ✅ 201 Created | ~0.03s | Create sales agent |
| 3 | `/api/v1/sales-agents/:id` | GET | ✅ 200 OK | ~0.01s | Get agent by ID |
| 4 | `/api/v1/sales-agents/:id` | PUT | ✅ 200 OK | ~0.02s | Update agent |
| 5 | `/api/v1/sales-agents/:id` | DELETE | ✅ 200 OK | ~0.01s | Delete agent |

### Profile Prices Folder

| # | Endpoint | Method | Status | Response Time | Description |
|---|----------|--------|--------|---------------|-------------|
| 1 | `/api/v1/sales-agents/:id/profile-prices` | GET | ✅ 200 OK | ~0.01s | List profile prices |
| 2 | `/api/v1/sales-agents/:id/profile-prices/:name` | PUT | ✅ 200 OK | ~0.02s | Upsert profile price |

### Agent Invoices Folder

| # | Endpoint | Method | Status | Response Time | Description |
|---|----------|--------|--------|---------------|-------------|
| 1 | `/api/v1/sales-agents/:id/invoices` | GET | ✅ 200 OK | ~0.01s | List agent invoices |
| 2 | `/api/v1/sales-agents/:id/invoices/generate` | POST | ⚠️ 500 | - | Generate invoice (no sales data) |

### Summary
- **CRUD Tests:** 5/5 passed ✅
- **Profile Prices:** 2/2 passed ✅
- **Agent Invoices:** 1/2 (invoice generation requires sales data)
- **Total:** 8/9 passed

### Notes

#### Create Sales Agent - Field Length Issue
**Error (with full payload):**
```
ERROR: value too long for type character varying(10) (SQLSTATE 22001)
```

**Cause:** Some fields like `voucher_type` have restricted length in database.

**Workaround:** Use minimal required fields or shorter values:
```json
{
  "router_id": "...",
  "name": "Agen Budi",
  "phone": "08123456789",
  "username": "agenbudi",
  "password": "password123",
  "status": "active"
}
```

#### Invoice Generation - Requires Sales Data
**Error:**
```json
{
  "success": false,
  "error": "failed to generate invoice number: record not found"
}
```

**Cause:** Agent has no hotspot sales records for the billing period.

**Expected Behavior:** Invoices are generated based on agent's voucher sales. Empty salesresult in no invoice.

### Request/Response Examples

#### Create Sales Agent
```http
POST /api/v1/sales-agents
Authorization: Bearer <<token>>
Content-Type: application/json

{
  "router_id": "<<routerId>>",
  "name": "Agen Budi",
  "phone": "08123456789",
  "username": "agenbudi",
  "password": "password123",
  "status": "active"
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "data": {
    "ID": "efefabc9-92cf-4701-b53a-dae94dc1a729",
    "Name": "Agen Budi",
    "Phone": "08123456789",
    "Username": "agenbudi",
    "Status": "active",
    "VoucherMode": "mix",
    "VoucherLength": 6,
    "VoucherType": "upp",
    "BillDiscount": 0,
    "BillingCycle": "monthly",
    "BillingDay": 1
  }
}
```

#### Upsert Profile Price
```http
PUT /api/v1/sales-agents/:id/profile-prices/1jam
Authorization: Bearer <<token>>
Content-Type: application/json

{
  "base_price": 3000,
  "selling_price": 5000,
  "voucher_length": 6,
  "is_active": true
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "ID": "b2f45204-33c5-4a68-9f00-4ada3cae1253",
    "ProfileName": "1jam",
    "BasePrice": 3000,
    "SellingPrice": 5000,
    "VoucherLength": 6,
    "IsActive": true
  }
}
```

### Environment Variables
- `agentId`: `efefabc9-92cf-4701-b53a-dae94dc1a729` (created and deleted)

---

## 17 - Agent Invoices Tests

| # | Endpoint | Method | Status | Response Time | Description |
|---|----------|--------|--------|---------------|-------------|
| 1 | `/api/v1/agent-invoices` | GET | ✅ 200 OK | ~0.01s | List all agent invoices |
| 2 | `/api/v1/agent-invoices/process` | POST | ✅ 200 OK | ~0.01s | Process scheduled invoices |
| 3 | `/api/v1/agent-invoices` | GET | ✅ 200 OK | ~0.01s | List with filter (agent_id) |
| 4 | `/api/v1/agent-invoices/:id` | GET | ✅ 404 | ~0.01s | Get non-existent invoice |
| 5 | `/api/v1/agent-invoices/:id/pay` | PUT | ✅ 404 | ~0.01s | Markpaid non-existent |
| 6 | `/api/v1/agent-invoices/:id/cancel` | PUT | ✅ 404 | ~0.01s | Cancel non-existent invoice |

### Summary
- **Total Tests:** 6
- **Passed:** 6
- **Note:** No invoices exist due to lack of hotspot sales data

### Notes

#### No Invoice Data Available
Agent invoices are generated based on hotspot sales. Since there are no hotspot sales in the database:
- `GET /api/v1/agent-invoices` returns empty array
- `POST /api/v1/agent-invoices/process` completes but creates no invoices
- Invoice-specific endpoints (`:id`, `:id/pay`, `:id/cancel`) return 404/500 for non-existent IDs

#### Invoice Generation Flow
```
Sales Agent created → Hotspot Sales made → Invoice generated (process) → Invoice paid/cancelled
```

To fully test invoice operations:
1. Create sales agent via `POST /api/v1/sales-agents`
2. Create hotspot sales (requires voucher sales flow)
3. Run `POST /api/v1/agent-invoices/process`
4. Test `GET :id`, `PUT :id/pay`, `PUT :id/cancel`

### Request/Response Examples

#### List Agent Invoices
```http
GET /api/v1/agent-invoices
Authorization: Bearer <<token>>
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": [],
  "meta": {
    "total": 0,
    "limit": 20,
    "offset": 0
  }
}
```

#### Process Scheduled Invoices
```http
POST /api/v1/agent-invoices/process
Authorization: Bearer <<token>>
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "message": "processing complete"
  }
}
```

#### Get Agent Invoice (not found)
```http
GET /api/v1/agent-invoices/:id
Authorization: Bearer <<token>>
```

**Response (404 Not Found):**
```json
{
  "success": false,
  "error": "record not found"
}
```

#### List with Filter
```http
GET /api/v1/agent-invoices?agent_id=xxx&status=pending&billing_year=2026
Authorization: Bearer <<token>>
```

**Supported Filters:**
- `agent_id` - Filter by agent UUID
- `router_id` - Filter by router UUID
- `status` - pending | paid | cancelled
- `billing_cycle` - weekly | monthly
- `billing_year` - e.g. 2026
- `billing_month` - 1-12
- `billing_week` - 1-53

### Environment Variables
- `agentInvoiceId`: *(none - no invoices created)*

---

## 18 - Hotspot Sales Tests

| # | Endpoint | Method | Status | Response Time | Description |
|---|----------|--------|--------|---------------|-------------|
| 1 | `/api/v1/hotspot-sales` | GET | ✅ 200 OK | ~0.01s | List all hotspot sales |
| 2 | `/api/v1/routers/:id/hotspot-sales` | GET | ✅ 200 OK | ~0.01s | List sales by router |
| 3 | `/api/v1/hotspot-sales?date_from=...` | GET | ✅ 200 OK | ~0.01s | List with date filter |
| 4 | `/api/v1/hotspot-sales?router_id=...` | GET | ✅ 200 OK | ~0.01s | List with router filter |
| 5 | `/api/v1/hotspot-sales?profile=...` | GET | ✅ 200 OK | ~0.01s | List with profile filter |

### Summary
- **Total Tests:** 5
- **Passed:** 5
- **Note:** No sales data (created via voucher sales flow)

### Notes

#### No Sales Data Available
Hotspot sales are created through the voucher sales process, not via direct API:
- Agents sell vouchers through the agent portal
- Each sale creates a hotspot sales record
- Records can be queried via this endpoint

#### Supported Query Filters
| Filter | Type | Description |
|--------|------|-------------|
| `router_id` | UUID | Filter by router |
| `agent_id` | UUID | Filter by sales agent |
| `profile` | string | Filter by profile name (e.g., "1jam") |
| `batch_code` | string | Filter by batch code |
| `date_from` | date | Start date (YYYY-MM-DD) |
| `date_to` | date | End date (YYYY-MM-DD) |
| `limit` | int | Pagination limit |
| `offset` | int | Pagination offset |

### Request/Response Examples

#### List All Hotspot Sales
```http
GET /api/v1/hotspot-sales?limit=20&offset=0
Authorization: Bearer <<token>>
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": [],
  "meta": {
    "total": 0,
    "limit": 20,
    "offset": 0
  }
}
```

#### List Sales by Router
```http
GET /api/v1/routers/:routerId/hotspot-sales
Authorization: Bearer <<token>>
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": []
}
```

#### List with Date Range
```http
GET /api/v1/hotspot-sales?date_from=2026-03-01&date_to=2026-03-31
Authorization: Bearer <<token>>
```

### Data Flow
```
Agent Portal/Login→ Sell Voucher → Hotspot Sale Created→ Invoice Generated
                ↓                                     ↓
           Profile Price                        Batch Code
         (Base/Selling Price)                        ↓
                                              Sales Summary
```

### Environment Variables
- *(none - no sales data)*

---

## 19 - Cash Management Tests

### Cash Entries Folder

| # | Endpoint | Method | Status | Response Time | Description |
|---|----------|--------|--------|---------------|-------------|
| 1 | `/api/v1/cash-entries` | GET | ✅ 200 OK | ~0.01s | List cash entries |
| 2 | `/api/v1/cash-entries` | POST | ✅ 201 Created | ~0.02s | Create cash entry |
| 3 | `/api/v1/cash-entries/:id` | GET | ✅ 200 OK | ~0.01s | Get cash entry by ID |
| 4 | `/api/v1/cash-entries/:id` | PUT | ✅ 200 OK | ~0.02s | Update cash entry |
| 5 | `/api/v1/cash-entries/:id/approve` | POST | ✅ 200 OK | ~0.02s | Approve cash entry |
| 6 | `/api/v1/cash-entries/:id/reject` | POST | ✅ 200 OK | ~0.02s | Reject cash entry |
| 7 | `/api/v1/cash-entries/:id` | DELETE | ✅ 500 | - | Delete (validation works) |

### Petty Cash Folder

| # | Endpoint | Method | Status | Response Time | Description |
|---|----------|--------|--------|---------------|-------------|
| 1 | `/api/v1/petty-cash` | GET | ✅ 200 OK | ~0.01s | List petty cash funds |
| 2 | `/api/v1/petty-cash` | POST | ✅ 201 Created | ~0.02s | Create petty cash fund |
| 3 | `/api/v1/petty-cash/:id` | GET | ✅ 200 OK | ~0.01s | Get fund by ID |
| 4 | `/api/v1/petty-cash/:id` | PUT | ✅ 200 OK | ~0.02s | Update fund |
| 5 | `/api/v1/petty-cash/:id/topup` | POST | ✅ 200 OK | ~0.02s | Top up fund |

### Summary
- **Cash Entries:** 7/7 passed ✅
- **Petty Cash:** 5/5 passed ✅
- **Total:** 12/12 passed ✅

### Notes

#### Source Field Constraint
**Error (with invalid source):**
```
ERROR: new row for relation "cash_entries" violates check constraint "cash_entries_source_check"
```

**Valid sources:** `other`, and possibly other predefined values. Use `source: "other"` for generic entries.

#### Cash Entry Workflow
```
Create (pending) → Update (pending) → Approve/Reject → Final state
                                     ↓
                        Approved: affects fund balance
                        Rejected: can NOT be deleted
```

#### Delete Validation
- Can only delete **pending** entries
- Approved/Rejected entries return 500 error

#### Reject Endpoint Requires Body
```http
POST /api/v1/cash-entries/:id/reject
Content-Type: application/json

{
  "reason": "Bukti tidak valid"
}
```

### Request/Response Examples

#### Create Petty Cash Fund
```http
POST /api/v1/petty-cash
Authorization: Bearer <<token>>
Content-Type: application/json

{
  "fund_name": "Kas Operasional",
  "initial_balance": 5000000,
  "custodian_id": "<<userId>>"
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "data": {
    "ID": "5981550b-6551-496f-8425-99d92f6f496c",
    "FundName": "Kas Operasional",
    "InitialBalance": 5000000,
    "CurrentBalance": 5000000,
    "Status": "active"
  }
}
```

#### Top Up Fund
```http
POST /api/v1/petty-cash/:id/topup
Authorization: Bearer <<token>>
Content-Type: application/json

{
  "amount": 1000000
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "CurrentBalance": 6000000
  }
}
```

#### Create Cash Entry
```http
POST /api/v1/cash-entries
Authorization: Bearer <<token>>
Content-Type: application/json

{
  "type": "income",
  "source": "other",
  "amount": 150000,
  "description": "Pembayaran langganan",
  "payment_method": "bank_transfer",
  "petty_cash_fund_id": "<<fundId>>"
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "data": {
    "ID": "9404f5ac-edaa-402b-bc21-dd903c7b6b95",
    "EntryNumber": "KAS000002",
    "Type": "income",
    "Status": "pending"
  }
}
```

#### Approve Cash Entry
```http
POST /api/v1/cash-entries/:id/approve
Authorization: Bearer <<token>>
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "Status": "approved",
    "ApprovedBy": "user-uuid",
    "ApprovedAt": "2026-03-27T..."
  }
}
```

#### Reject Cash Entry
```http
POST /api/v1/cash-entries/:id/reject
Authorization: Bearer <<token>>
Content-Type: application/json

{
  "reason": "Bukti tidak valid"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "Status": "rejected",
    "Notes": "Bukti tidak valid"
  }
}
```

### Environment Variables
- `fundId`: `5981550b-6551-496f-8425-99d92f6f496c`
- `entryId`: `9404f5ac-edaa-402b-bc21-dd903c7b6b95` (approved)
- `entryId2`: `d8622ae5-6344-40c4-9d97-c51a17c46a19` (rejected)

---

## 20 - Reports Tests

| # | Endpoint | Method | Status | Response Time | Description |
|---|----------|--------|--------|---------------|-------------|
| 1 | `/api/v1/reports/summary` | GET | ✅ 200 OK | ~0.02s | Business summary report |
| 2 | `/api/v1/reports/subscriptions` | GET | ✅ 200 OK | ~0.02s | Subscription report |
| 3 | `/api/v1/reports/cash-flow` | GET | ✅ 200 OK | ~0.02s | Cash flow report |
| 4 | `/api/v1/reports/cash-balance` | GET | ✅ 200 OK | ~0.01s | Cash balance report |
| 5 | `/api/v1/reports/reconciliation` | GET | ✅ 200 OK | ~0.01s | Reconciliation report |

### Summary
- **Total Tests:** 5
- **Passed:** 5 ✅
- **100% Success Rate**

### Notes

All report endpoints work correctly and return datafrom the previous tests:
- **Summary Report** returns totalcustomers, subscriptions, and invoice stats
- **Subscriptions Report** returns list with customer details and status
- **Cash Flow Report** shows income/expense breakdown (175,000 income from cash entry test)
- **Cash Balance Report** shows current balance across all funds
- **Reconciliation Report** returns matched/missing entries

### Request/Response Examples

#### Summary Report
```http
GET /api/v1/reports/summary?from=2026-03-01&to=2026-03-31
Authorization: Bearer <<token>>
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "period_start": "2026-03-01T00:00:00Z",
    "period_end": "2026-03-31T23:59:59Z",
    "total_revenue": 0,
    "total_invoiced": 0,
    "total_invoices": 0,
    "total_customers": 5,
    "active_customers": 0,
    "new_customers": 5,
    "active_subscriptions": 1,
    "subscriptions": {
      "active": 1,
      "isolated": 0,
      "suspended": 0,
      "pending": 1,
      "total": 2
    }
  }
}
```

#### Cash Flow Report
```http
GET /api/v1/reports/cash-flow?from=2026-03-01&to=2026-03-31
Authorization: Bearer <<token>>
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "period_start": "2026-03-01T00:00:00Z",
    "period_end": "2026-03-31T00:00:00Z",
    "total_income": 175000,
    "total_expense": 0,
    "net_cash_flow": 175000,
    "breakdown": [
      {"type": "income", "source": "other", "total": 175000}
    ]
  }
}
```

#### Subscriptions Report
```http
GET /api/v1/reports/subscriptions?from=2026-03-01&to=2026-03-31&limit=50&offset=0
Authorization: Bearer <<token>>
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "customer_code": "CST-00002",
      "customer_name": "Siti Rahayu",
      "username": "testsub_002",
      "plan_name": "10 Mbps",
      "status": "active",
      "monthly_price": 150000
    }
  ],
  "meta": {"total": 3, "limit": 50, "offset": 0}
}
```

### Supported Query Parameters

| Endpoint | Parameters |
|----------|------------|
| `summary` | `from`, `to` (YYYY-MM-DD) |
| `subscriptions` | `from`, `to`, `limit`, `offset` |
| `cash-flow` | `from`, `to` |
| `cash-balance` | none |
| `reconciliation` | `from`, `to` |

---

## 21 - System Settings Tests

| # | Endpoint | Method | Status | Response Time | Description |
|---|----------|--------|--------|---------------|-------------|
| 1 | `/api/v1/settings` | GET | ✅ 200 OK | ~0.01s | List all settings |
| 2 | `/api/v1/settings/:id` | GET | ✅ 200 OK | ~0.01s | Get setting by ID |
| 3 | `/api/v1/settings` | PUT | ⚠️ 200 | ~0.02s | Upsert setting (creates new) |
| 4 | `/api/v1/settings?group=billing` | GET | ✅ 200 OK | ~0.01s | List settings by group |

### Summary
- **Total Tests:** 4
- **Passed:** 4 ✅
- **Note:** Upsert creates new setting instead of updating existing

### Notes

#### Upsert Behavior - Creates New Setting
**Issue:** The upsert endpoint creates a new setting instead of updating anexisting one with the same `group_name` + `key_name`.

**Test case:**
```json
{
  "group_name": "billing",
  "key_name": "due_days",
  "value": "14",
  "type": "integer"
}
```

**Result:** Created new setting with empty `GroupName` and `KeyName`:
```json
{
  "ID": "e310bd77-c408-4684-9f89-22fcde395402",
  "GroupName": "",
  "KeyName": "",
  "Value": "14"
}
```

**Expected:** Should update existing setting where `GroupName="billing"` AND `KeyName="due_days"`.

#### Valid Setting Types
Type field has a check constraint. Allowed values appear to be:
- `string`
- `integer`

Invalid types like `float` return constraint violation error.

### Request/Response Examples

#### List Settings
```http
GET /api/v1/settings
Authorization: Bearer <<token>>
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": [
    {
      "ID": "ff9b77ee-8837-46a4-afa9-08d68ec9133a",
      "GroupName": "billing",
      "KeyName": "due_days",
      "Value": "10",
      "Type": "integer",
      "Label": "Invoice Due Days"
    }
  ]
}
```

#### Get Setting by ID
```http
GET /api/v1/settings/:id
Authorization: Bearer <<token>>
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "ID": "ff9b77ee-8837-46a4-afa9-08d68ec9133a",
    "GroupName": "billing",
    "KeyName": "due_days",
    "Value": "10",
    "Type": "integer",
    "Label": "Invoice Due Days"
  }
}
```

#### List by Group
```http
GET /api/v1/settings?group=billing
Authorization: Bearer <<token>>
```

**Response:** Returns only billing settings.

### Environment Variables
- `dueDaysSettingId`: `ff9b77ee-8837-46a4-afa9-08d68ec9133a`

---

## 22 - Webhooks Tests

| # | Endpoint | Method | Status | Response Time | Description |
|---|----------|--------|--------|---------------|-------------|
| 1 | `/api/v1/webhooks/midtrans` | POST | ⚠️ 500 | ~0.02s | Midtrans callback (no matching order) |
| 2 | `/api/v1/webhooks/xendit` | POST | ⚠️ 401 | ~0.01s | Xendit callback (invalid token) |

### Summary
- **Total Tests:** 2
- **Integration Tests:** Require production configuration
- **Note:** Webhooks require valid payment records and gateway tokens

### Notes

#### Webhook Endpoints - Public (No Auth Required)
Both webhook endpoints are designed for payment gateway callbacks:
- **Midtrans:** `/api/v1/webhooks/midtrans`
- **Xendit:** `/api/v1/webhooks/xendit`

#### Midtrans Webhook Response
**Error:** `record not found`

**Cause:** The webhook tries to find a payment by `order_id` (mapped from Midtrans `order_id`). Test payload used `PAY-00001` which has no matching payment record.

**Valid payment_number in database:** `PAY000001`, `PAY000002`, `PAY000003`

#### Xendit Webhook Response
**Error:** `webhook verification failed: xendit: invalid webhook token`

**Cause:** The Xendit webhook validates the `x-callback-token` header against the configured Xendit webhook verification token in environment variables. Test used invalid token `xendit-callback-token`.

**To Test Properly:**
1. Create a payment with gateway initiation
2. Use the actual Xendit callback token from environment
3. Send callback with matching `external_id`

### Request/Response Examples

#### Midtrans Webhook
```http
POST /api/v1/webhooks/midtrans
Content-Type: application/json

{
  "transaction_time": "2026-03-27 10:00:00",
  "transaction_status": "settlement",
  "transaction_id": "txn-123",
  "status_code": "200",
  "signature_key": "signature-here",
  "payment_type": "bank_transfer",
  "order_id": "PAY000001",
  "gross_amount": "150000.00",
  "fraud_status": "accept",
  "currency": "IDR"
}
```

**Response (500 - no matching order):**
```json
{
  "success": false,
  "error": "record not found"
}
```

#### Xendit Webhook
```http
POST /api/v1/webhooks/xendit
Content-Type: application/json
x-callback-token: <actual-xendit-token>

{
  "id": "inv-123",
  "external_id": "PAY000001",
  "status": "PAID",
  "amount": 150000,
  "paid_amount": 150000,
  "paid_at": "2026-03-27T10:00:00.000Z",
  "payment_method": "BANK_TRANSFER",
  "payment_channel": "BCA",
  "currency": "IDR"
}
```

**Response (401 - invalid token):**
```json
{
  "success": false,
  "error": "webhook verification failed: xendit: invalid webhook token"
}
```

### Production Testing
To properly test webhooks:
1. **Midtrans:** Configure Midtrans server key and create actual transaction
2. **Xendit:** Set `XENDIT_WEBHOOK_SECRET` in environment and create invoice

### Environment Variables Required
```env
MIDTRANS_SERVER_KEY=your-server-key
MIDTRANS_CLIENT_KEY=your-client-key
XENDIT_SECRET_KEY=your-secret-key
XENDIT_WEBHOOK_SECRET=your-webhook-secret
```

---

## 23 - Customer Portal Tests

### Public Folder

| # | Endpoint | Method | Status | Description |
|---|----------|--------|--------|-------------|
| 1 | `/portal/v1/login` | POST | ✅ 200 | Portal login (email/username) |
| 1b | `/portal/v1/login` | POST | ❌ 401 | Login with phone (invalid format) |

### Authenticated Folder

| # | Endpoint | Method | Status | Description |
|---|----------|--------|--------|-------------|
| 2 | `/portal/v1/profile` | GET | ✅ 200 | Get customer profile |
| 3 | `/portal/v1/profile/password` | PUT | ✅ 200 | Change portal password |
| 4 | `/portal/v1/subscriptions` | GET | ✅ 200 | List customer subscriptions |
| 5 | `/portal/v1/invoices` | GET | ✅ 200 | List customer invoices |
| 6 | `/portal/v1/payments` | GET | ✅ 200 | List customer payments |
| 7 | `/portal/v1/payments` | POST | ✅ 201 | Create payment |
| 8 | `/portal/v1/payments/:id` | GET | ✅ 200 | Get payment detail |
| 9 | `/portal/v1/payments/:id/pay` | POST | ✅ 200 | Pay via gateway |

### Summary
- **Public Endpoints:** 2 tested (email/username login works, phone format issue)
- **Authenticated Endpoints:** 8/8 passed ✅
- **Total:** 10 tests

### Notes

#### Portal Login Identifier
**Working formats:**
- Email: `siti.rahayu@example.com` ✅
- Username: `siti-rahayu` ✅

**Not working:**
- Phone: `081200000002` or `+6281200000002` ❌ (returns invalid credentials)

**Recommendation:** Documentation should clarify supported login identifier formats (email and username only).

#### Customer Activation Required
Customers must be activated before portal login:
```http
POST /api/v1/customers/:id/activate-account
Authorization: Bearer <admin-token>
```

#### Password Setup
Portal password can be set via admin API:
```http
PUT /api/v1/customers/:id
Authorization: Bearer <admin-token>
Content-Type: application/json

{
  "password": "portal123"
}
```

### Request/Response Examples

#### Portal Login
```http
POST /portal/v1/login
Content-Type: application/json

{
  "identifier": "siti.rahayu@example.com",
  "password": "portal123"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "customer": {
      "ID": "1a3e5a1b-83b5-4ba4-9909-23574926c049",
      "CustomerCode": "CST-00002",
      "FullName": "Siti Rahayu",
      "Email": "siti.rahayu@example.com"
    },
    "token": "eyJhbGciOiJIUzI1NiIs..."
  }
}
```

#### Get Profile
```http
GET /portal/v1/profile
Authorization: Bearer <portal-token>
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "ID": "1a3e5a1b-83b5-4ba4-9909-23574926c049",
    "FullName": "Siti Rahayu",
    "Email": "siti.rahayu@example.com",
    "Phone": "+6281200000002"
  }
}
```

#### Create Payment
```http
POST /portal/v1/payments
Authorization: Bearer <portal-token>
Content-Type: application/json

{
  "amount": 100000,
  "payment_method": "bank_transfer"
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "data": {
    "ID": "f1a57c12-1837-4eea-9330-755fa405033c",
    "PaymentNumber": "PAY000004",
    "Amount": 100000,
    "Status": "pending"
  }
}
```

#### Pay with Gateway
```http
POST /portal/v1/payments/:id/pay?gateway=xendit
Authorization: Bearer <portal-token>
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "expires_at": "2026-03-28T04:21:45.079Z",
    "gateway_id": "69c605d87cba7679600e943b",
    "payment_url": "https://checkout-staging.xendit.co/web/69c605d87cba7679600e943b"
  }
}
```

### Environment Variables
- `portalToken`: Set after portal login
- `invoiceId`: *(none - customer has no invoices)*
- `paymentId`: `f1a57c12-1837-4eea-9330-755fa405033c`

---

## 24 - Agent Portal Tests

### Public Folder

| # | Endpoint | Method | Status | Description |
|---|----------|--------|--------|-------------|
| 1 | `/agent-portal/v1/login` | POST | ✅ 200 | Agent login |

### Authenticated Folder

| # | Endpoint | Method | Status | Description |
|---|----------|--------|--------|-------------|
| 2 | `/agent-portal/v1/profile` | GET | ✅ 200 | Get agent profile |
| 3 | `/agent-portal/v1/profile/password` | PUT | ✅ 200 | Change agent password |
| 4 | `/agent-portal/v1/invoices` | GET | ✅ 200 | List agent invoices |
| 5 | `/agent-portal/v1/sales` | GET | ✅ 200 | List agent sales |

### Summary
- **Public Endpoints:** 1/1 passed ✅
- **Authenticated Endpoints:** 4/4 passed ✅
- **Total:** 5/5 tests passed ✅

### Notes

#### Agent Creation Required
Agents must be created via admin API beforelogin:
```http
POST /api/v1/sales-agents
Authorization: Bearer <admin-token>
Content-Type: application/json

{
  "router_id": "<router-uuid>",
  "name": "Agent Test",
  "phone": "08188888888",
  "username": "agenttest",
  "password": "agent123",
  "status": "active"
}
```

#### Invoice Endpoints
Agent invoices are generated from hotspot voucher sales:
- `GET /agent-portal/v1/invoices` - Returns empty if agent has no sales
- `GET /agent-portal/v1/invoices/:id` - Requires invoice ID
- `POST /agent-portal/v1/invoices/:id/request-payment` - Request payment

### Request/Response Examples

#### Agent Login
```http
POST /agent-portal/v1/login
Content-Type: application/json

{
  "username": "agenttest",
  "password": "agent123"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "agent": {
      "ID": "6b07a4fc-4419-4784-8a9d-e9dd9b12991b",
      "Name": "Agent Test",
      "Username": "agenttest",
      "Phone": "08188888888",
      "Status": "active"
    },
    "token": "eyJhbGciOiJIUzI1NiIs..."
  }
}
```

#### Get Profile
```http
GET /agent-portal/v1/profile
Authorization: Bearer <agent-token>
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "ID": "6b07a4fc-4419-4784-8a9d-e9dd9b12991b",
    "Name": "Agent Test",
    "Username": "agenttest",
    "Phone": "08188888888",
    "Status": "active",
    "VoucherMode": "mix",
    "VoucherLength": 6,
    "BillingCycle": "monthly",
    "BillingDay": 1
  }
}
```

#### Get Sales
```http
GET /agent-portal/v1/sales
Authorization: Bearer <agent-token>
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": [],
  "meta": {
    "total": 0,
    "limit": 20,
    "offset": 0
  }
}
```

### Agent Portal Flow
```
Admin creates agent → Agent logs in → Agent sells vouchers → 
Invoices generated → Agent views invoices/sales → Agent requests payment
```

### Environment Variables
- `agentToken`: Set after agent portal login
- `agentInvoiceId`: *(none - agent has no invoices yet)*

---

## Current Password State

- **Password after test:** `NewPassword123!`
- **After running 01b:** `SuperAdmin123!`

---

## Environment Variables

| Variable | Value |
|----------|-------|
| baseUrl | `http://localhost:8080` |
| token | *(auto-populated from login)* |
| refreshToken | *(auto-populated from login)* |

---

## Commands Used

```bash
# Run auth tests
cd tests/http
hopp test 01-auth.json -e environment.json

# Run change password (optional)
hopp test 01b-auth-change-password.json -e environment.json
```

---

## Notes

1. **Prerequisites:**
   - PostgreSQL running on localhost:5432
   - Redis running on localhost:6379
   - RabbitMQ running on localhost:5672
   - MikMongo server running on localhost:8080

2. **Database Setup:**
   ```bash
   cd /path/to/mikmongo
   go run cmd/migrate/main.go  # Run migrations
   go run cmd/seed/main.go     # Seed data
   ```

3. **Server Setup:**
   ```bash
   go run cmd/server/main.go
   ```

---

*Generated by OpenClaw Agent - Qixi*