---
phase: 01-auth-api-foundation
plan: 02
subsystem: auth-ui
tags: [react, tanstack-query, react-hook-form, zod, zustand, lucide, sonner]

# Dependency graph
requires:
  - phase: 01-01
    provides: "Auth store, API functions, Axios clients, Zod schemas, TypeScript types"
provides:
  - Admin login form with real API integration and Indonesian text
  - Admin auth hooks (useAdminLogin, useAdminLogout, useAdminChangePassword)
  - MikMongo branded auth layout replacing Shadcn Admin
  - Admin change password page with Zod validation
  - Sign-out dialog with real API logout and Indonesian text
  - NavUser reading real admin user data from Zustand store
affects: [01-03, 01-04]

# Tech tracking
tech-stack:
  added: []
  patterns:
    - "TanStack Query mutations for auth operations (login, logout, change password)"
    - "Auth hooks pattern: mutation wrappers around API functions with toast feedback and navigation"

key-files:
  created:
    - website/src/hooks/use-admin-auth.ts
    - website/src/features/auth/change-password/index.tsx
  modified:
    - website/src/features/auth/auth-layout.tsx
    - website/src/features/auth/sign-in/components/user-auth-form.tsx
    - website/src/features/auth/sign-in/index.tsx
    - website/src/components/sign-out-dialog.tsx
    - website/src/components/layout/nav-user.tsx
    - website/src/components/layout/app-sidebar.tsx

key-decisions:
  - "Login form uses direct API call in onSubmit (not useAdminLogin hook) for simpler form integration"
  - "Logout dialog uses useAdminLogout hook which swallows API errors and always clears auth"
  - "NavUser reads admin user from Zustand store directly, no props needed"
  - "Avatar fallback generates initials from user's full_name dynamically"

patterns-established:
  - "Auth form pattern: react-hook-form + zodResolver + LoginFormSchema + async onSubmit with try/catch"
  - "Toast-based error feedback for API errors, FormMessage for field-level validation"
  - "Loading state: Loader2 spinner + Indonesian text + disabled button during submission"

requirements-completed: [AUTH-01, AUTH-03, AUTH-04, AUTH-05]

# Metrics
duration: 3min
completed: 2026-04-02
---

# Phase 01 Plan 02: Admin Auth UI Summary

**MikMongo branded admin login form with real API integration, change password page, Indonesian text throughout, and TanStack Query auth hooks replacing all mock auth patterns**

## Performance

- **Duration:** 3 min
- **Started:** 2026-04-02T22:43:44Z
- **Completed:** 2026-04-02T22:46:54Z
- **Tasks:** 2
- **Files modified:** 7

## Accomplishments
- Replaced mock sleep-based login with real `adminLogin` API call using react-hook-form + Zod
- Auth layout now shows MikMongo branding (M logo + title) instead of Shadcn Admin
- Login form has full Indonesian text, Remember me checkbox, Forgot password link, loading spinner
- Created TanStack Query auth hooks: useAdminLogin, useAdminLogout, useAdminChangePassword
- Built change password page with three fields and Zod validation (min 8 chars, confirm match)
- Sign-out dialog calls real logout API and shows Indonesian confirmation text
- NavUser reads real admin user from Zustand store, removed all template menu items

## Task Commits

Each task was committed atomically:

1. **Task 1: Rewrite auth layout and admin login form with real API integration** - `5bd0be7` (feat)
2. **Task 2: Build admin change password page and update logout dialog** - `20d8c40` (feat)

_Note: Commits are in the website submodule (git submodule)._

## Files Created/Modified
- `website/src/features/auth/auth-layout.tsx` - Replaced Shadcn Admin with MikMongo M logo + title
- `website/src/features/auth/sign-in/components/user-auth-form.tsx` - Real adminLogin API call, Indonesian text, Remember me, Forgot password, removed mock/social login
- `website/src/features/auth/sign-in/index.tsx` - Admin Portal subtitle, removed redirect param per D-07
- `website/src/hooks/use-admin-auth.ts` - TanStack Query mutations for login/logout/change-password
- `website/src/features/auth/change-password/index.tsx` - Change password page with 3 fields and Zod validation
- `website/src/components/sign-out-dialog.tsx` - Real logout via useAdminLogout, Indonesian text
- `website/src/components/layout/nav-user.tsx` - Reads admin user from Zustand store, simplified menu
- `website/src/components/layout/app-sidebar.tsx` - Removed user prop from NavUser call

## Decisions Made
- Login form uses direct `adminLogin` API call in `onSubmit` rather than `useAdminLogin` hook, keeping form state management simple
- `useAdminLogout` swallows API errors and always clears auth -- ensures user can always log out even if API is unreachable
- NavUser no longer takes props -- reads from Zustand store directly, making it self-contained
- Avatar fallback dynamically generates initials from admin user's `full_name` field

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 1 - Bug] Fixed NavUser consumer after removing user prop**
- **Found during:** Task 2 (Build admin change password page and update logout dialog)
- **Issue:** NavUser component was rewritten to read from Zustand store (no props), but `app-sidebar.tsx` still passed `user={sidebarData.user}` prop
- **Fix:** Removed the `user` prop from `<NavUser />` call in app-sidebar.tsx
- **Files modified:** website/src/components/layout/app-sidebar.tsx
- **Verification:** `npx tsc --noEmit` passes with zero errors
- **Committed in:** `20d8c40` (Task 2 commit)

---

**Total deviations:** 1 auto-fixed (1 bug)
**Impact on plan:** Consumer fix was necessary to maintain correctness. NavUser no longer receives props.

## Issues Encountered
- No issues -- both tasks executed cleanly following Plan 01 interfaces

## User Setup Required

None - no external service configuration required.

## Next Phase Readiness
- Admin auth UI is complete: login, change password, logout flows all working with real API
- use-admin-auth.ts hooks available for route guard integration in Plan 03
- ChangePassword page component ready for route registration in Plan 04
- NavUser shows real admin user data -- ready for authenticated layout rendering

## Self-Check: PASSED

All created/modified files exist. Both commits verified. TypeScript compilation passes with zero errors.

---
*Phase: 01-auth-api-foundation*
*Completed: 2026-04-02*
