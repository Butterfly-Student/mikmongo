---
phase: 03-customers-routers-subscriptions
plan: 06
subsystem: ui
tags: [react, tanstack-query, zod, react-hook-form, shadcn-ui, bandwidth-profiles]

# Dependency graph
requires:
  - phase: 03-customers-routers-subscriptions
    provides: profiles API (createProfile, deleteProfile), useProfiles, useDeleteProfile, profile table with columns
provides:
  - updateProfile API function (PUT /routers/{router_id}/bandwidth-profiles/{id})
  - useUpdateProfile React Query mutation hook
  - EditProfileDialog pre-populated with current profile values
  - onEdit action in profile dropdown menu
  - Profile edit workflow wired in profiles/index.tsx
affects: [gap-closure, BW-03]

# Tech tracking
tech-stack:
  added: []
  patterns:
    - Edit dialog pattern mirrors create dialog with defaultValues from existing entity
    - unknown error type narrowing for onError callbacks (no any)

key-files:
  created:
    - website/src/features/profiles/components/edit-profile-dialog.tsx
  modified:
    - website/src/api/profiles.ts
    - website/src/hooks/use-profiles.ts
    - website/src/features/profiles/data/columns.tsx
    - website/src/features/profiles/index.tsx

key-decisions:
  - "Edit dialog resets on close via handleOpenChange to prevent stale form state"
  - "editTarget state in parent controls dialog open/close without separate boolean flag"
  - "useUpdateProfile invalidates profiles queryKey scoped by routerId matching existing patterns"

patterns-established:
  - "Edit dialogs pre-populate from entity prop using ?? fallback to empty/default values"
  - "Error type in onError callbacks uses unknown with explicit cast, never any"

requirements-completed: [BW-03]

# Metrics
duration: 8min
completed: 2026-04-04
---

# Phase 03 Plan 06: Profile Update Summary

**Bandwidth profile edit dialog with updateProfile API, useUpdateProfile hook, and onEdit dropdown action closing gap BW-03**

## Performance

- **Duration:** 8 min
- **Started:** 2026-04-04T12:00:00Z
- **Completed:** 2026-04-04T12:08:00Z
- **Tasks:** 1
- **Files modified:** 5

## Accomplishments
- Added `updateProfile` API function calling PUT /routers/{id}/bandwidth-profiles/{id}
- Added `useUpdateProfile` React Query mutation hook with success toast and cache invalidation
- Created `EditProfileDialog` pre-populated with all current profile fields including MikroTik settings
- Wired Edit action into profile table dropdown menu with Pencil icon
- Connected edit dialog to `editTarget` state in profiles page

## Task Commits

Each task was committed atomically:

1. **Task 1: Add updateProfile API, useUpdateProfile hook, edit dialog, and wire into UI** - `a6c1fc1` (feat)

**Plan metadata:** pending docs commit

## Files Created/Modified
- `website/src/api/profiles.ts` - Added `updateProfile` function (PUT endpoint)
- `website/src/hooks/use-profiles.ts` - Added `useUpdateProfile` hook; fixed any->unknown error types in existing hooks
- `website/src/features/profiles/components/edit-profile-dialog.tsx` - New edit dialog pre-populated from profile entity
- `website/src/features/profiles/data/columns.tsx` - Added `onEdit` to ColumnActions, Pencil icon, Edit menu item
- `website/src/features/profiles/index.tsx` - Added `editTarget` state, `EditProfileDialog` render, `onEdit` callback

## Decisions Made
- Edit dialog resets form on close to prevent stale form data carrying over to next edit
- `editTarget` state doubles as open flag (`!!editTarget`) avoiding redundant boolean state
- Error handling uses `unknown` type with explicit narrowing cast, consistent with TypeScript strict rules

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 1 - Bug] Fixed any->unknown error types in useCreateProfile and useDeleteProfile**
- **Found during:** Task 1 (adding useUpdateProfile)
- **Issue:** Existing hooks used `err: any` which violates TypeScript strict rules and project coding style
- **Fix:** Changed to `err: unknown` with explicit narrowing cast matching the new useUpdateProfile pattern
- **Files modified:** website/src/hooks/use-profiles.ts
- **Verification:** TypeScript compilation passes with no errors
- **Committed in:** a6c1fc1 (Task 1 commit)

---

**Total deviations:** 1 auto-fixed (Rule 1 - type safety fix)
**Impact on plan:** Necessary for consistent code quality. No scope creep.

## Issues Encountered
None - TypeScript compiled cleanly after all changes.

## User Setup Required
None - no external service configuration required.

## Next Phase Readiness
- BW-03 gap closure complete: admin can now create, read, update, and delete bandwidth profiles
- Profile CRUD is fully operational
- Ready for next gap closure plan or phase transition

---
*Phase: 03-customers-routers-subscriptions*
*Completed: 2026-04-04*
