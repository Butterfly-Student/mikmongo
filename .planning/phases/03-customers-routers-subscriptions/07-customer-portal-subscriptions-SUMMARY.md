---
phase: 03-customers-routers-subscriptions
plan: "07"
subsystem: customer-portal
tags: [customer-portal, subscriptions, portal-api, tanstack-query]
dependency_graph:
  requires: []
  provides: [customer-portal-subscriptions-page, portal-subscription-api]
  affects: [customer-portal-routes]
tech_stack:
  added: []
  patterns: [feature-based-directory, react-query-hooks, zod-schema-validation, tanstack-router-file-routing]
key_files:
  created:
    - website/src/api/portal/subscription.ts
    - website/src/hooks/use-customer-portal.ts
    - website/src/features/customer-portal/subscriptions.tsx
    - website/src/routes/customer/index.tsx
    - website/src/routes/customer/subscriptions/index.tsx
  modified: []
decisions:
  - Customer portal index route redirects to /customer/subscriptions for direct access convenience
  - Portal subscription API uses separate PortalSubscriptionListResponseSchema (not admin schema) since portal endpoint may not include meta pagination
key_decisions:
  - "Customer portal index redirects to /customer/subscriptions for seamless UX"
  - "Portal subscription list schema defined inline without pagination meta (portal endpoint returns all customer subscriptions)"
metrics:
  duration: "2 minutes"
  completed: "2026-04-04"
  tasks_completed: 1
  tasks_total: 1
  files_created: 5
  files_modified: 0
requirements:
  - SUB-05
---

# Phase 03 Plan 07: Customer Portal Subscriptions Summary

**One-liner:** Customer portal subscriptions page at /customer/subscriptions showing subscription cards with status badges, IP, profile, expiry, via portal-scoped customerClient API.

## What Was Built

A complete customer portal subscriptions feature consisting of:

1. **Portal Subscription API** (`website/src/api/portal/subscription.ts`): Calls `GET /portal/v1/subscriptions` using `customerClient` with Bearer token from Zustand store. Validates response with Zod schema.

2. **Portal Hook** (`website/src/hooks/use-customer-portal.ts`): `usePortalSubscriptions()` hook using TanStack Query with 2-minute stale time and `['portal-subscriptions']` query key.

3. **Customer Portal Subscriptions Page** (`website/src/features/customer-portal/subscriptions.tsx`): Card-based layout with:
   - Loading state: 3 skeleton cards
   - Empty state: WifiOff icon with "No Subscriptions Found" message
   - Subscription cards showing: username, status badge (color-coded), IP address, MikroTik profile, expiry date, activated date

4. **Customer Index Route** (`website/src/routes/customer/index.tsx`): Redirects `/customer` to `/customer/subscriptions` automatically.

5. **Customer Subscriptions Route** (`website/src/routes/customer/subscriptions/index.tsx`): TanStack Router file-based route rendering `CustomerPortalSubscriptions`.

## Verification Results

- TypeScript compilation: PASSED (no errors)
- `listPortalSubscriptions` exported from portal API: VERIFIED
- `usePortalSubscriptions` exported from hook: VERIFIED
- `customerClient` used (not `adminClient`) in portal API: VERIFIED
- Route file exists at `website/src/routes/customer/subscriptions/index.tsx`: VERIFIED
- No `adminClient` usage in portal API directory: VERIFIED

## Deviations from Plan

None — plan executed exactly as written.

## Known Stubs

None. All data is wired to live API via customerClient.

## Self-Check: PASSED

- `website/src/api/portal/subscription.ts`: FOUND
- `website/src/hooks/use-customer-portal.ts`: FOUND
- `website/src/features/customer-portal/subscriptions.tsx`: FOUND
- `website/src/routes/customer/index.tsx`: FOUND
- `website/src/routes/customer/subscriptions/index.tsx`: FOUND
- Commit `012e355`: FOUND
