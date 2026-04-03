---
phase: 03-customers-routers-subscriptions
plan: 03
subsystem: ui
tags: [react, tanstack-table, tanstack-query, zustand, zod, shadcn-ui, alert-dialog]

# Dependency graph
requires:
  - phase: 03-customers-routers-subscriptions
    plan: 01
    provides: "router Zustand store with selected router, router table pattern, API client"
  - phase: 03-customers-routers-subscriptions
    plan: 02
    provides: "customer list hook, bandwidth profiles hook with router-scoped loading, customer schema"
provides:
  - "Subscription feature module (schema, columns, table, dialogs)"
  - "ConfirmActionDialog reusable component for destructive actions"
  - "Subscription API layer and TanStack Query hooks for all CRUD and state transitions"
  - "Subscriptions route at /subscriptions consuming active router from Zustand"
affects: [billing, invoices, payments, customer-portal, agent-portal]

# Tech tracking
tech-stack:
  added: []
  patterns: ["router-scoped data fetching via Zustand store", "destructive action confirmation dialog pattern", "subscription state machine actions (activate/suspend/isolate/restore/terminate)"]

key-files:
  created:
    - "website/src/components/ui/confirm-action-dialog.tsx"
    - "website/src/features/subscriptions/data/schema.ts"
    - "website/src/features/subscriptions/data/columns.tsx"
    - "website/src/features/subscriptions/components/subscription-table.tsx"
    - "website/src/features/subscriptions/components/create-subscription-dialog.tsx"
    - "website/src/features/subscriptions/index.tsx"
    - "website/src/routes/_authenticated/subscriptions/index.tsx"
    - "website/src/lib/schemas/subscription.ts"
    - "website/src/api/subscription.ts"
    - "website/src/hooks/use-subscriptions.ts"
  modified:
    - "website/src/components/layout/data/sidebar-data.ts"

key-decisions:
  - "Subscriptions consume active router from Zustand store rather than router-scoped URL params"
  - "All destructive actions (suspend, isolate, terminate, delete) route through ConfirmActionDialog before API call"
  - "Sidebar updated to enable /subscriptions route (disabled -> active)"

patterns-established:
  - "Router-scoped resource pattern: feature uses router from Zustand store to scope API calls"
  - "Action confirmation pattern: row actions set confirmAction state, single ConfirmActionDialog renders"

requirements-completed: [SUB-01, SUB-02, SUB-03, SUB-04, SUB-05]

# Metrics
duration: 7min
completed: 2026-04-03
---

# Phase 03 Plan 03: Subscriptions Summary

**Subscription management with router-scoped data table, CreateSubscriptionDialog with customer/profile selects, and ConfirmActionDialog for all destructive state transitions**

## Performance

- **Duration:** 7 min
- **Started:** 2026-04-03T14:42:43Z
- **Completed:** 2026-04-03T14:50:05Z
- **Tasks:** 3
- **Files modified:** 11

## Accomplishments
- Reusable ConfirmActionDialog component for explicit destructive action confirmations
- Global Subscriptions view consuming active router from Zustand store with status filtering
- CreateSubscriptionDialog with dependent customer and bandwidth profile selects
- All subscription state transition actions (activate, suspend, isolate, restore, terminate) wired with confirmation
- Sidebar navigation updated to enable Subscriptions route

## Task Commits

Each task was committed atomically:

1. **Task 1: Build Shared Confirmation Dialog** - `1ee740b` (feat)
2. **Task 2: Global Subscriptions View** - `c7a55b5` (feat)
3. **Task 3: Subscription Actions & Assignment** - `1abbae9` (feat)

## Files Created/Modified
- `website/src/components/ui/confirm-action-dialog.tsx` - Reusable AlertDialog wrapper with destructive variant, pending state, custom labels
- `website/src/features/subscriptions/data/schema.ts` - Zod schemas for Subscription and CreateSubscription
- `website/src/features/subscriptions/data/columns.tsx` - TanStack Table column definitions with context-sensitive row actions
- `website/src/features/subscriptions/components/subscription-table.tsx` - Data table with status filter, pagination, empty state
- `website/src/features/subscriptions/components/create-subscription-dialog.tsx` - Form dialog with customer and bandwidth profile selects
- `website/src/features/subscriptions/index.tsx` - Page composition with router-scoped data and action confirmation state machine
- `website/src/routes/_authenticated/subscriptions/index.tsx` - TanStack Router route at /subscriptions
- `website/src/lib/schemas/subscription.ts` - SubscriptionResponse Zod schema for API validation
- `website/src/api/subscription.ts` - API functions for all subscription endpoints (CRUD + state transitions)
- `website/src/hooks/use-subscriptions.ts` - TanStack Query hooks with query invalidation and toast notifications
- `website/src/components/layout/data/sidebar-data.ts` - Subscriptions nav item enabled with /subscriptions URL

## Decisions Made
- Router-scoped data fetching via Zustand store rather than URL params, consistent with profiles pattern
- All destructive actions (suspend, isolate, terminate, delete) route through ConfirmActionDialog with descriptive copy per UI spec
- Subscription status column uses color-coded badges (default for active, outline for suspended/isolated, destructive for terminated)
- Empty state renders router selection prompt when no router is active in Zustand store

## Deviations from Plan

None - plan executed exactly as written.

## Issues Encountered
None

## User Setup Required
None - no external service configuration required.

## Next Phase Readiness
- Subscription feature module complete and ready for billing/invoice integration
- Router-scoped resource pattern established for future MikroTik features (PPP, Hotspot, Network)
- ConfirmActionDialog reusable for any future destructive actions across the dashboard

## Self-Check: PASSED

All 11 files verified present. All 3 commits verified in git history.

---
*Phase: 03-customers-routers-subscriptions*
*Completed: 2026-04-03*
