---
phase: 03-customers-routers-subscriptions
plan: 02
subsystem: ui
tags: [react, tanstack-table, zod, tanstack-query, tanstack-router]

# Dependency graph
requires:
  - phase: 01-auth-portals
    provides: admin auth, Axios interceptors, authenticated layout
  - phase: 01-routers-profiles
    provides: router/profile patterns, TanStack Table conventions
provides:
  - Customer CRUD table with activate/deactivate/delete actions
  - Registration pipeline table with approve/reject dialogs
  - Customer and Registration API layer with Zod validation
  - TanStack Query hooks for customer lifecycle management
  - Tabbed Customers page (Customers + Registrations)
affects: [03-subscriptions]

# Tech tracking
tech-stack:
  added: []
  patterns:
    - Tab-based page layout for related entity management
    - Registration pipeline as data table with inline approve/reject
    - Dependent select (profiles load after router selected)

key-files:
  created:
    - website/src/lib/schemas/customer.ts
    - website/src/api/customer.ts
    - website/src/hooks/use-customers.ts
    - website/src/features/customers/data/columns.tsx
    - website/src/features/customers/data/registration-columns.tsx
    - website/src/features/customers/components/customer-table.tsx
    - website/src/features/customers/components/registration-table.tsx
    - website/src/features/customers/components/create-customer-dialog.tsx
    - website/src/features/customers/components/delete-customer-dialog.tsx
    - website/src/features/customers/components/approve-registration-dialog.tsx
    - website/src/features/customers/components/reject-registration-dialog.tsx
    - website/src/features/customers/index.tsx
    - website/src/routes/_authenticated/customers/index.tsx
  modified: []

key-decisions:
  - "Tabbed layout for Customers and Registrations instead of separate pages"
  - "Default registration filter set to 'pending' to surface actionable items"
  - "Approve dialog uses dependent select: profiles load only after router selected"

patterns-established:
  - "Tab-based page: use Tabs/TabsContent for related CRUD entities on same route"
  - "Registration pipeline: approve/reject as dropdown actions on pending rows"
  - "Dependent selects: profile query enabled only when router_id is set"

requirements-completed: [CUST-01, CUST-02, CUST-03, CUST-04, CUST-05, CUST-06, CUST-07]

# Metrics
duration: 6min
completed: 2026-04-03
---

# Phase 03 Plan 02: Customers Summary

**Customer CRUD table with registration pipeline (approve/reject), activate/deactivate actions, and tabbed page layout using TanStack Table**

## Performance

- **Duration:** 6 min
- **Started:** 2026-04-03T14:34:48Z
- **Completed:** 2026-04-03T14:40:42Z
- **Tasks:** 3
- **Files modified:** 13

## Accomplishments
- Full customer data table with search, status filtering, and pagination
- Registration pipeline table with approve/reject actions for pending registrations
- Create Customer dialog with required fields (name, phone) and optional PPP/Hotspot credentials
- Delete Customer dialog with destructive confirmation
- Approve Registration dialog with router selection and dependent profile loading
- Reject Registration dialog with mandatory rejection reason
- TanStack Query hooks for all customer and registration API operations
- Zod schemas at lib level for API response validation

## Task Commits

Each task was committed atomically:

1. **Task 1: Define Customer Schema and Data Layer** - `d7c5fa5` (feat)
2. **Task 2: Build Customer Table and Registration Pipeline UX** - `85bed82` (feat)
3. **Task 3: Build Customer Views and Routing** - `a93685d` (feat)

## Files Created/Modified

- `website/src/lib/schemas/customer.ts` - Zod schemas for CustomerResponse, RegistrationResponse with API list/detail types
- `website/src/api/customer.ts` - API functions for customer CRUD, activate/deactivate, and registration approve/reject
- `website/src/hooks/use-customers.ts` - TanStack Query hooks for customer and registration operations
- `website/src/features/customers/data/columns.tsx` - Customer table column definitions with avatar, status badges, and row actions
- `website/src/features/customers/data/registration-columns.tsx` - Registration table columns with status-based action visibility
- `website/src/features/customers/components/customer-table.tsx` - Customer data table with search, filter, pagination
- `website/src/features/customers/components/registration-table.tsx` - Registration data table with search, filter, pagination
- `website/src/features/customers/components/create-customer-dialog.tsx` - Create customer form with validation
- `website/src/features/customers/components/delete-customer-dialog.tsx` - Delete confirmation dialog
- `website/src/features/customers/components/approve-registration-dialog.tsx` - Approve registration with router/profile selection
- `website/src/features/customers/components/reject-registration-dialog.tsx` - Reject registration with reason input
- `website/src/features/customers/index.tsx` - Main Customers page with tabbed layout
- `website/src/routes/_authenticated/customers/index.tsx` - TanStack Router route file

## Decisions Made
- **Tabbed layout:** Customers and Registrations share a single page with Tabs rather than separate routes, keeping related lifecycle management together
- **Default pending filter:** Registrations tab defaults to showing "pending" status, surfacing the most actionable items for admins
- **Dependent profile select:** In the Approve dialog, bandwidth profiles only load after a router is selected (matching the useProfiles hook's routerId requirement)
- **Type-safe error handling:** Used `unknown` type narrowing instead of `any` in hooks for type safety

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 2 - Missing Critical Functionality] Registration pipeline scope expansion**
- **Found during:** Task 2 (Build Customer Table and Registration Pipeline UX)
- **Issue:** Plan mentioned "Approve" and "Reject" buttons but did not specify the approve dialog requires router_id and optional profile_id per OpenAPI spec, nor the reject dialog requires a reason
- **Fix:** Created full ApproveRegistrationDialog with router dropdown and dependent profile loading, and RejectRegistrationDialog with required reason textarea
- **Files modified:** N/A (new files created)
- **Committed in:** `85bed82` (Task 2 commit)

**2. [Rule 1 - Bug] Fixed approve dialog circular import**
- **Found during:** Task 2 (ApproveRegistrationDialog)
- **Issue:** Initial implementation incorrectly imported non-existent exports from use-customers
- **Fix:** Corrected imports to use useRouters from use-routers and useProfiles from use-profiles
- **Files modified:** approve-registration-dialog.tsx
- **Committed in:** `85bed82` (Task 2 commit)

**3. [Rule 2 - Missing Critical Functionality] Delete Customer dialog**
- **Found during:** Task 3 (Build Customer Views and Routing)
- **Issue:** Plan's delete action was mentioned in columns but no delete dialog was specified; destructive actions require confirmation per context decisions (D-03)
- **Fix:** Added DeleteCustomerDialog with destructive styling and confirmation
- **Files modified:** N/A (new file)
- **Committed in:** `a93685d` (Task 3 commit)

---

**Total deviations:** 3 auto-fixed (2 missing critical, 1 bug)
**Impact on plan:** All auto-fixes necessary for completeness and correctness. The delete dialog follows established D-03 pattern from CONTEXT.md. No scope creep.

## Issues Encountered
- Website directory is a separate git repository (not a submodule), requiring commits within `website/` rather than the project root

## Known Stubs
None - all features are fully wired with real API endpoints.

## Next Phase Readiness
- Customer management complete and ready for subscription linking in plan 03
- Registration pipeline fully functional for admin approval workflow
- All API calls use validated Zod schemas matching OpenAPI spec

## Self-Check: PASSED

All 13 files verified present. All 3 commits verified in git history.

---
*Phase: 03-customers-routers-subscriptions*
*Completed: 2026-04-03*
