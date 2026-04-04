---
phase: 03-customers-routers-subscriptions
plan: 04
subsystem: ui
tags: [react, typescript, tanstack-query, shadcn-ui, mikrotik, router-management]

# Dependency graph
requires:
  - phase: 03-customers-routers-subscriptions
    provides: Router table with create/sync/test-connection actions and Zustand router store

provides:
  - updateRouter, deleteRouter, syncAllRouters API functions in api/router.ts
  - useUpdateRouter, useDeleteRouter, useSyncAllRouters mutation hooks
  - EditRouterDialog component with pre-populated form and optional password update
  - DeleteRouterDialog component with confirmation before DELETE API call
  - Sync All Routers button with loading spinner in router page header
  - Customers and Routers sidebar nav items enabled with real URLs (/customers, /routers)

affects: [phase-04-billing, phase-05-sales, sidebar-navigation, router-management]

# Tech tracking
tech-stack:
  added: []
  patterns:
    - unknown error type narrowing pattern for useMutation onError handlers
    - react-hook-form values prop for pre-populating edit forms from existing entity data
    - Optional password update: omit field from payload if blank rather than sending empty string

key-files:
  created:
    - website/src/features/routers/components/edit-router-dialog.tsx
    - website/src/features/routers/components/delete-router-dialog.tsx
  modified:
    - website/src/api/router.ts
    - website/src/hooks/use-routers.ts
    - website/src/features/routers/data/columns.tsx
    - website/src/features/routers/index.tsx
    - website/src/components/layout/data/sidebar-data.ts

key-decisions:
  - "Edit router omits password from payload when field is blank, preserving existing credentials"
  - "useUpdateRouter/useDeleteRouter/useSyncAllRouters use unknown error type with type assertion narrowing (not any)"
  - "Sidebar nav enabled for Customers (/customers) and Routers (/routers) by removing disabled flag and setting real URLs"

patterns-established:
  - "Edit dialogs use react-hook-form values prop to sync with external router prop, not just defaultValues"
  - "Conditional payload field: include password only when non-empty string provided"

requirements-completed: [RTR-02, RTR-06]

# Metrics
duration: 12min
completed: 2026-04-04
---

# Phase 03 Plan 04: Router CRUD + Sidebar Summary

**Full router CRUD via edit/delete dialogs and sync-all button, plus sidebar nav enabled for Customers and Routers pages**

## Performance

- **Duration:** 12 min
- **Started:** 2026-04-04T12:00:00Z
- **Completed:** 2026-04-04T12:12:00Z
- **Tasks:** 2
- **Files modified:** 7

## Accomplishments

- Added `updateRouter`, `deleteRouter`, `syncAllRouters` API functions and corresponding mutation hooks with proper `unknown` error type handling
- Created `EditRouterDialog` (pre-populates all fields, optional password update) and `DeleteRouterDialog` (confirmation before DELETE)
- Removed `console.log` stub from router page — all CRUD actions now wired to real API hooks
- Added "Sync All" button with spinner to router page header
- Enabled sidebar navigation for Customers (`/customers`) and Routers (`/routers`) by removing `disabled: true` and setting real URLs

## Task Commits

Each task was committed atomically:

1. **Task 1: Add router API functions, hooks, and sidebar navigation fix** - `b701565` (feat)
2. **Task 2: Create edit/delete dialogs and wire into router page and columns** - `dda410d` (feat)

## Files Created/Modified

- `website/src/api/router.ts` - Added updateRouter, deleteRouter, syncAllRouters functions
- `website/src/hooks/use-routers.ts` - Added useUpdateRouter, useDeleteRouter, useSyncAllRouters hooks; fixed useTestRouterConnection to use unknown error type
- `website/src/components/layout/data/sidebar-data.ts` - Enabled Customers and Routers nav items with real URLs
- `website/src/features/routers/components/edit-router-dialog.tsx` - New: edit dialog with pre-populated form
- `website/src/features/routers/components/delete-router-dialog.tsx` - New: delete confirmation dialog
- `website/src/features/routers/data/columns.tsx` - Added onEdit to ColumnActions interface and Edit Router menu item with Pencil icon
- `website/src/features/routers/index.tsx` - Added edit/delete state, Sync All button, EditRouterDialog, DeleteRouterDialog

## Decisions Made

- Edit router dialog omits password from the PUT payload when the field is blank, preserving the router's existing credentials without requiring re-entry
- All new hooks use `unknown` error type with narrowing type assertion rather than `any`, consistent with TypeScript coding style rules
- Sidebar navigation enabled in-place (removed `disabled: true`, set real URL) rather than creating a separate routing file

## Deviations from Plan

None - plan executed exactly as written.

## Issues Encountered

None - TypeScript compilation passed cleanly on first attempt for both tasks.

## User Setup Required

None - no external service configuration required.

## Next Phase Readiness

- Router CRUD is complete: create, edit, delete, sync, sync-all, test-connection all wired
- Customers and Routers pages are reachable from sidebar navigation
- RTR-02 (edit/delete router) and RTR-06 (sync all routers) verification gaps are closed
- Ready for billing phase (invoices, payments) or additional feature phases

---
*Phase: 03-customers-routers-subscriptions*
*Completed: 2026-04-04*
