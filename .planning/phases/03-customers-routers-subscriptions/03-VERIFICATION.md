---
phase: 03-customers-routers-subscriptions
verified: 2026-04-04T13:00:00Z
status: passed
score: 6/6 must-haves verified
re_verification:
  previous_status: gaps_found
  previous_score: 4/6
  gaps_closed:
    - "Admin can add, edit, delete, test connection to, and sync MikroTik routers (RTR-02, RTR-06)"
    - "Admin can update customer details (CUST-03)"
    - "Admin can create and manage bandwidth profiles — update now implemented (BW-03)"
    - "Customer portal displays the customer's active subscriptions (SUB-05)"
    - "Sidebar navigation enables access to Customers and Routers pages"
  gaps_remaining: []
  regressions: []
human_verification:
  - test: "Navigate to Customers, Routers, and Subscriptions from sidebar"
    expected: "All three nav items are clickable and route to the correct pages without errors"
    why_human: "Runtime navigation requires browser interaction"
  - test: "Edit a router and verify changes persist"
    expected: "EditRouterDialog opens pre-populated with current values; saving updates the table"
    why_human: "Pre-population from react-hook-form values prop and optional password omission require visual confirmation"
  - test: "Open customer portal at /customer/subscriptions"
    expected: "Subscription cards render with status badges, IP, profile, expiry; WifiOff empty state when no subscriptions"
    why_human: "Portal uses customerClient with bearer token — requires logged-in customer session"
  - test: "Registration approval with dependent profile select"
    expected: "Router dropdown appears; profiles load only after a router is selected; approval wires profile to customer"
    why_human: "Dependent select behavior requires browser interaction to verify loading state"
---

# Phase 3: Customers, Routers & Subscriptions — Re-Verification Report

**Phase Goal:** Admin can manage the full customer lifecycle (create, activate, registration pipeline) and manage MikroTik routers with bandwidth profiles and subscription plans
**Verified:** 2026-04-04T13:00:00Z
**Status:** passed
**Re-verification:** Yes — after gap closure plans 04, 05, 06, 07

## Goal Achievement

### Observable Truths

| #   | Truth | Status | Evidence |
| --- | ----- | ------ | -------- |
| 1   | Admin can view, create, update, and activate/deactivate customers with pagination and search | VERIFIED | `EditCustomerDialog` + `useUpdateCustomer` wired in `customers/index.tsx`. Create, delete, activate/deactivate all present from previous verification. |
| 2   | Admin can manage the customer registration pipeline (view pending, approve with router/profile assignment, reject with reason) | VERIFIED | No change from initial verification — was VERIFIED. |
| 3   | Admin can add, edit, delete, test connection to, and sync MikroTik routers | VERIFIED | `EditRouterDialog` + `DeleteRouterDialog` + `useSyncAllRouters` wired in `routers/index.tsx`. Console.log stub eliminated. Sync All button present with spinner. |
| 4   | Admin can create and manage bandwidth profiles (rate-limit, burst) per router | VERIFIED | `EditProfileDialog` + `useUpdateProfile` wired in `profiles/index.tsx`. Create and delete confirmed from previous verification. |
| 5   | Admin can create subscriptions, assign profiles to customers on routers, and manage subscription lifecycle (activate, suspend, isolate, restore, terminate) | VERIFIED | No change from initial verification — was VERIFIED. |
| 6   | Customer portal displays the customer's active subscriptions | VERIFIED | `CustomerPortalSubscriptions` feature component exists; `usePortalSubscriptions` hook wired to `listPortalSubscriptions` via `customerClient`; route at `/customer/subscriptions/index.tsx` renders the component; `/customer/` index redirects to `/customer/subscriptions`. |

**Score:** 6/6 truths verified

### Required Artifacts

| Artifact | Expected | Status | Details |
| -------- | -------- | ------ | ------- |
| `website/src/features/routers/components/edit-router-dialog.tsx` | Edit router dialog with pre-populated form | VERIFIED | Full form with all fields, password-optional logic, wired to `useUpdateRouter` |
| `website/src/features/routers/components/delete-router-dialog.tsx` | Delete confirmation dialog | VERIFIED | AlertDialog with destructive button, wired to `useDeleteRouter` |
| `website/src/features/routers/index.tsx` | Router page with all CRUD + Sync All | VERIFIED | `editTarget`, `deleteTarget` state; `useSyncAllRouters` with spinner button; all dialogs rendered; no console.log stub |
| `website/src/api/router.ts` | `updateRouter`, `deleteRouter`, `syncAllRouters` API functions | VERIFIED | All three present, calling `adminClient.put`, `.delete`, `.post` respectively |
| `website/src/hooks/use-routers.ts` | `useUpdateRouter`, `useDeleteRouter`, `useSyncAllRouters` hooks | VERIFIED | All three exported with cache invalidation and toast feedback |
| `website/src/features/customers/components/edit-customer-dialog.tsx` | Edit customer dialog | VERIFIED | Pre-populated via `useEffect`; password optional; wired to `useUpdateCustomer` |
| `website/src/hooks/use-customers.ts` | `useUpdateCustomer` hook | VERIFIED | Exported with `invalidateQueries(['customers'])` and toast |
| `website/src/features/customers/index.tsx` | Customer page with edit wiring | VERIFIED | `editTarget` state; `onEdit` callback on columns; `EditCustomerDialog` rendered |
| `website/src/features/profiles/components/edit-profile-dialog.tsx` | Edit profile dialog | VERIFIED | Pre-populated with all profile fields including MikroTik settings; wired to `useUpdateProfile` |
| `website/src/api/profiles.ts` | `updateProfile` API function | VERIFIED | PUT call to `/routers/{routerId}/bandwidth-profiles/{id}` with Zod validation |
| `website/src/hooks/use-profiles.ts` | `useUpdateProfile` hook | VERIFIED | Exported with `['profiles', routerId]` invalidation; all hooks use `unknown` error type |
| `website/src/features/profiles/index.tsx` | Profiles page with edit wiring | VERIFIED | `editTarget` state; `onEdit` in columns factory; `EditProfileDialog` rendered |
| `website/src/api/portal/subscription.ts` | Portal subscription API using `customerClient` | VERIFIED | `listPortalSubscriptions` calls `customerClient.get('/subscriptions')`; validated with Zod |
| `website/src/hooks/use-customer-portal.ts` | `usePortalSubscriptions` hook | VERIFIED | TanStack Query with `['portal-subscriptions']` key, 2-min stale time |
| `website/src/features/customer-portal/subscriptions.tsx` | Customer portal subscriptions page | VERIFIED | Card layout with loading skeleton, empty state (WifiOff icon), subscription cards showing status badge, IP, profile, expiry, activated date |
| `website/src/routes/customer/subscriptions/index.tsx` | `/customer/subscriptions` route | VERIFIED | Renders `CustomerPortalSubscriptions` |
| `website/src/routes/customer/index.tsx` | `/customer` index redirect | VERIFIED | `beforeLoad` throws `redirect({ to: '/customer/subscriptions' })` |
| `website/src/components/layout/data/sidebar-data.ts` | Customers and Routers nav items enabled | VERIFIED | `{ title: 'Customers', url: '/customers' }` and `{ title: 'Routers', url: '/routers' }` — no `disabled` flag |

### Key Link Verification

| From | To | Via | Status | Details |
| ---- | --- | --- | ------ | ------- |
| `EditRouterDialog` | PUT /api/v1/routers/{id} | `useUpdateRouter` -> `updateRouter` | WIRED | Form submit calls `updateRouter({ id, data: payload })` |
| `DeleteRouterDialog` | DELETE /api/v1/routers/{id} | `useDeleteRouter` -> `deleteRouter` | WIRED | Confirm button calls `deleteRouter(router.id)` |
| Sync All button | POST /api/v1/routers/sync-all | `useSyncAllRouters` -> `syncAllRouters` | WIRED | `onClick={() => syncAllRouters()}` in router page header |
| `EditCustomerDialog` | PUT /api/v1/customers/{id} | `useUpdateCustomer` -> `updateCustomer` | WIRED | Submit calls `updateCustomer({ id: customer.id, data: payload })` |
| `EditProfileDialog` | PUT /api/v1/routers/{id}/bandwidth-profiles/{id} | `useUpdateProfile` -> `updateProfile` | WIRED | Submit calls `updateProfile({ routerId, id: profile.id, data })` |
| `CustomerPortalSubscriptions` | GET /portal/v1/subscriptions | `usePortalSubscriptions` -> `listPortalSubscriptions` -> `customerClient` | WIRED | `const { data: subscriptions } = usePortalSubscriptions()` rendered in card grid |
| Sidebar Customers | `/customers` route | `sidebar-data.ts` | WIRED | `url: '/customers'` with no `disabled` flag |
| Sidebar Routers | `/routers` route | `sidebar-data.ts` | WIRED | `url: '/routers'` with no `disabled` flag |
| Sidebar Subscriptions | `/subscriptions` route | `sidebar-data.ts` | WIRED | Confirmed from previous verification |

### Data-Flow Trace (Level 4)

| Artifact | Data Variable | Source | Produces Real Data | Status |
| -------- | ------------- | ------ | ------------------ | ------ |
| `CustomerPortalSubscriptions` | `subscriptions` | `usePortalSubscriptions` -> `listPortalSubscriptions` -> `customerClient.get('/subscriptions')` -> Zod parse | Yes — live API call via `customerClient` | FLOWING |
| `Routers` page | `data?.routers` | `useRouters` -> `listRouters` -> `adminClient.get('/routers')` | Yes | FLOWING (confirmed prior verification) |
| `Customers` page | `customersData?.customers` | `useCustomers` -> `listCustomers` -> `adminClient.get('/customers')` | Yes | FLOWING (confirmed prior verification) |
| `BandwidthProfiles` | `data?.profiles` | `useProfiles` -> `listProfiles` -> `adminClient.get('/routers/{id}/bandwidth-profiles')` | Yes | FLOWING (confirmed prior verification) |

### Behavioral Spot-Checks

| Behavior | Command | Result | Status |
| -------- | ------- | ------ | ------ |
| TypeScript compilation | `npx tsc --noEmit` | No errors (zero output) | PASS |
| `editTarget` state in routers/index.tsx | code read | `setEditTarget`, `<EditRouterDialog router={editTarget} open={!!editTarget}` present | PASS |
| `deleteTarget` state in routers/index.tsx | code read | `setDeleteTarget`, `<DeleteRouterDialog router={deleteTarget} open={!!deleteTarget}` present | PASS |
| Sync All button | code read | `useSyncAllRouters`, `onClick={() => syncAllRouters()}` with `disabled={isSyncingAll}` | PASS |
| `onEdit` wired to customer columns | code read | `onEdit: (customer) => setEditTarget(customer)` in `createCustomerColumns` call | PASS |
| `onEdit` wired to profile columns | code read | `onEdit: (profile) => setEditTarget(profile)` in `createColumns` call | PASS |
| Customer portal route renders component | file read | `component: () => <CustomerPortalSubscriptions />` in route file | PASS |
| Portal API uses `customerClient` | file read | `import { customerClient } from '@/lib/axios/customer-client'` in `api/portal/subscription.ts` | PASS |
| No console.log in router feature | grep | No matches in `website/src/features/routers/` | PASS |
| Sidebar items not disabled | file read | Customers and Routers entries have no `disabled` property | PASS |

### Requirements Coverage

| Requirement | Source Plan | Description | Status | Evidence |
| ----------- | ----------- | ----------- | ------ | -------- |
| CUST-01 | 02-customers | Admin can view customer list with pagination, search, and filters | SATISFIED | Confirmed from previous verification |
| CUST-02 | 02-customers | Admin can create new customer | SATISFIED | Confirmed from previous verification |
| CUST-03 | 02-customers | Admin can update customer details | SATISFIED | `EditCustomerDialog` + `useUpdateCustomer` now wired (was BLOCKED) |
| CUST-04 | 02-customers | Admin can activate/deactivate customer accounts | SATISFIED | Confirmed from previous verification |
| CUST-05 | 02-customers | Admin can manage customer registrations pipeline | SATISFIED | Confirmed from previous verification |
| CUST-06 | 02-customers | Admin can approve registrations with router and profile assignment | SATISFIED | Confirmed from previous verification |
| CUST-07 | 02-customers | Admin can reject registrations with reason | SATISFIED | Confirmed from previous verification |
| RTR-01 | 01-routers | Admin can view list of MikroTik routers with status | SATISFIED | Confirmed from previous verification |
| RTR-02 | 01-routers | Admin can add/edit/delete router configurations | SATISFIED | `EditRouterDialog` and `DeleteRouterDialog` now implemented (was BLOCKED) |
| RTR-03 | 01-routers | Admin can select active router for context-dependent operations | SATISFIED | Confirmed from previous verification |
| RTR-04 | 01-routers | Admin can sync router data from MikroTik device | SATISFIED | Confirmed from previous verification |
| RTR-05 | 01-routers | Admin can test connection to router | SATISFIED | Confirmed from previous verification |
| RTR-06 | 01-routers | Admin can sync all routers simultaneously | SATISFIED | `useSyncAllRouters` + Sync All button now implemented (was BLOCKED) |
| BW-01 | 01-routers | Admin can view bandwidth profiles per router | SATISFIED | Confirmed from previous verification |
| BW-02 | 01-routers | Admin can create bandwidth profile (name, rate-limit, burst) | SATISFIED | Confirmed from previous verification |
| BW-03 | 01-routers | Admin can update/delete bandwidth profiles | SATISFIED | `EditProfileDialog` + `useUpdateProfile` now implemented (was BLOCKED) |
| SUB-01 | 03-subscriptions | Admin can view subscriptions per router with status filters | SATISFIED | Confirmed from previous verification |
| SUB-02 | 03-subscriptions | Admin can create new subscription | SATISFIED | Confirmed from previous verification |
| SUB-03 | 03-subscriptions | Admin can update/terminate subscription | SATISFIED | Confirmed from previous verification |
| SUB-04 | 03-subscriptions | Admin can activate, suspend, isolate, and restore subscriptions | SATISFIED | Confirmed from previous verification |
| SUB-05 | 03-subscriptions | Customer portal shows their active subscriptions | SATISFIED | `CustomerPortalSubscriptions` + portal API + route now implemented (was BLOCKED) |

**Coverage Summary:** 21 requirements tracked. 21 SATISFIED. 0 BLOCKED. 0 ORPHANED.

### Anti-Patterns Found

| File | Line | Pattern | Severity | Impact |
| ---- | ---- | ------- | -------- | ------ |
| `website/src/hooks/use-routers.ts` | 39, 54 | `onError: (err: any)` in `useCreateRouter` and `useSyncRouter` | WARNING | Inconsistent with the rest of the file which uses `unknown`; not updated by plan 04 |
| `website/src/api/router.ts` | 32 | `createRouter(data: any)` parameter type | WARNING | Should be typed as `CreateRouter` from the feature schema |
| `website/src/features/routers/components/edit-router-dialog.tsx` | 37 | `zodResolver(createRouterSchema) as any` | WARNING | Type assertion workaround for zodResolver generic mismatch |
| `website/src/features/customers/components/edit-customer-dialog.tsx` | 41 | `zodResolver(createCustomerSchema) as never` | WARNING | Type assertion workaround for zodResolver generic mismatch |
| `website/src/features/profiles/components/edit-profile-dialog.tsx` | 34 | `zodResolver(createProfileSchema) as never` | WARNING | Type assertion workaround for zodResolver generic mismatch |

No BLOCKER anti-patterns. The `zodResolver as any/never` pattern is a known workaround for a type inference gap between `@hookform/resolvers` and `react-hook-form` — it does not affect runtime correctness. The `err: any` in two hooks is a minor inconsistency that does not block goal achievement.

### Human Verification Required

### 1. Full Sidebar Navigation Flow

**Test:** Open the app in a browser logged in as admin. Click Customers, Routers, and Subscriptions from the sidebar.
**Expected:** Each click navigates to the correct page and renders the data table without errors or blank screens.
**Why human:** Runtime TanStack Router navigation requires browser interaction.

### 2. Router Edit Round-Trip

**Test:** Click Edit on a router row. Verify the dialog pre-populates with current values. Change the name and leave password blank. Submit.
**Expected:** Form opens with current values. Name updates in table. No password field is sent in the API request.
**Why human:** `react-hook-form` `values` prop behavior and conditional payload logic require network tab inspection.

### 3. Customer Portal Subscriptions

**Test:** Log in as a customer. Navigate to `/customer/subscriptions`.
**Expected:** Subscription cards appear with status badge, IP address, profile name, expiry date, and activated date. Redirect from `/customer` to `/customer/subscriptions` occurs automatically.
**Why human:** Portal uses `customerClient` with Zustand bearer token — requires an actual customer session and live API responses.

### 4. Bandwidth Profile Edit

**Test:** Select a router, then click Edit on a bandwidth profile.
**Expected:** Edit dialog opens with current profile values pre-populated. Saving a change reflects in the table.
**Why human:** Pre-population via `defaultValues` (not `values` prop) may not re-sync across sequential opens — requires visual confirmation.

### Re-Verification Gap Summary

All five gaps from the initial verification (2026-04-03) have been closed by plans 04, 05, 06, and 07:

1. **RTR-02 / RTR-06 — Router edit, delete, sync-all (plan 04):** `EditRouterDialog`, `DeleteRouterDialog`, `useSyncAllRouters`, and a Sync All button are all present and wired. The `console.log` stub at line 39 is gone. Sidebar Customers and Routers nav items are enabled.

2. **CUST-03 — Customer update (plan 05):** `useUpdateCustomer` is exported from `use-customers.ts`. `EditCustomerDialog` exists and is rendered in `customers/index.tsx` with `onEdit` wired to the columns factory.

3. **BW-03 — Profile update (plan 06):** `updateProfile` API function, `useUpdateProfile` hook, and `EditProfileDialog` are all present and wired in the profiles page. Existing hooks had their `err: any` upgraded to `err: unknown`.

4. **SUB-05 — Customer portal subscriptions (plan 07):** Full implementation: `api/portal/subscription.ts` uses `customerClient`, `use-customer-portal.ts` provides the query hook, `features/customer-portal/subscriptions.tsx` renders the card layout, and two route files establish `/customer/subscriptions` with a redirect from `/customer`.

5. **Sidebar navigation (plan 04):** Customers and Routers nav items in `sidebar-data.ts` now have real URLs and no `disabled` property.

TypeScript compiles with zero errors across all new and modified files.

---

_Verified: 2026-04-04T13:00:00Z_
_Verifier: Claude (gsd-verifier)_
