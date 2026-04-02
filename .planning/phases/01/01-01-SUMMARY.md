---
phase: 01-auth-api-foundation
plan: 01
subsystem: auth
tags: [jwt, zustand, axios, zod, react, typescript]

# Dependency graph
requires: []
provides:
  - Zod schemas matching OpenAPI auth response shapes
  - TypeScript types derived from auth schemas
  - Zustand auth store with three portal slices and localStorage persistence
  - Admin Axios client with silent token refresh on 401
  - Customer and agent Axios clients with 401 redirect
  - Raw API functions for all auth endpoints (login, refresh, change password, logout, me)
affects: [01-02, 01-03, 01-04]

# Tech tracking
tech-stack:
  added: []
  patterns:
    - "Zustand persist with partialize for token/user persistence"
    - "Axios interceptor with refresh queue for concurrent 401 handling"
    - "Plain axios for unauthenticated endpoints, interceptored clients for authenticated"

key-files:
  created:
    - website/src/lib/schemas/auth.ts
    - website/src/api/types.ts
    - website/src/lib/axios/admin-client.ts
    - website/src/lib/axios/customer-client.ts
    - website/src/lib/axios/agent-client.ts
    - website/src/api/auth.ts
  modified:
    - website/src/stores/auth-store.ts
    - website/vite.config.ts
    - website/src/main.tsx
    - website/src/components/sign-out-dialog.tsx
    - website/src/features/auth/sign-in/components/user-auth-form.tsx

key-decisions:
  - "Zustand persist with partialize excludes actions and isHydrated from localStorage"
  - "onRehydrateStorage callback sets hydration gate for SSR safety"
  - "Admin refresh uses plain axios.post to avoid interceptor loop"
  - "Customer/agent clients redirect to portal login on 401 (no refresh, single token)"
  - "Vite proxy maps /api, /portal, /agent-portal to localhost:8080 for development"

patterns-established:
  - "Portal-specific Axios clients with separate baseURL and interceptor logic"
  - "API response parsing with Zod schemas for runtime type safety"
  - "Vanilla store getters (useAuthStore.getState()) in non-React contexts (interceptors)"

requirements-completed: [AUTH-01, AUTH-02, AUTH-08]

# Metrics
duration: 5min
completed: 2026-04-02
---

# Phase 01 Plan 01: Auth Data Layer Summary

**Zustand auth store with three portal slices (admin/customer/agent), localStorage persistence via partialize, Zod schemas matching OpenAPI exactly, and Axios clients with silent refresh queue for admin and 401 redirect for portals**

## Performance

- **Duration:** 5 min
- **Started:** 2026-04-02T22:35:29Z
- **Completed:** 2026-04-02T22:40:48Z
- **Tasks:** 2
- **Files modified:** 11

## Accomplishments
- Zod schemas matching all OpenAPI auth response shapes (admin login with access_token+refresh_token+user, admin refresh with token+refresh_token, customer portal with token+customer, agent portal with token+agent)
- Zustand store with persist middleware, three portal slices, partialize for selective persistence, and hydration gate
- Admin Axios client with silent token refresh using a failed-queue pattern for concurrent 401 handling
- Customer and agent Axios clients with simple 401 redirect (no refresh, single token per OpenAPI spec)
- Raw API functions for all auth endpoints using plain axios for login and interceptored clients for authenticated calls
- Vite dev proxy configured for all three API prefixes to backend port 8080

## Task Commits

Each task was committed atomically:

1. **Task 1: Create Zod schemas, API types, and Zustand auth store with three slices** - `2b2721b` (feat)
2. **Task 2: Create Axios clients and API auth functions** - `14c20aa` (feat)

_Note: Commits are in the website submodule (git submodule)._

## Files Created/Modified
- `website/src/lib/schemas/auth.ts` - Zod schemas for all auth request/response types matching OpenAPI exactly
- `website/src/api/types.ts` - TypeScript types derived from Zod schemas
- `website/src/stores/auth-store.ts` - Rewrote cookie-based mock auth with Zustand persist, three portal slices, hydration gate
- `website/src/lib/axios/admin-client.ts` - Admin Axios client with Bearer token and silent 401 refresh with retry queue
- `website/src/lib/axios/customer-client.ts` - Customer Axios client with Bearer token and 401 redirect
- `website/src/lib/axios/agent-client.ts` - Agent Axios client with Bearer token and 401 redirect
- `website/src/api/auth.ts` - Raw API functions for login, refresh, change password, logout, me for all portals
- `website/vite.config.ts` - Added dev proxy for /api, /portal, /agent-portal to localhost:8080
- `website/src/main.tsx` - Updated 401 handler to use new adminClearAuth instead of old auth.reset
- `website/src/components/sign-out-dialog.tsx` - Updated to use adminClearAuth instead of old auth.reset
- `website/src/features/auth/sign-in/components/user-auth-form.tsx` - Updated to use adminSetTokens instead of old auth.setUser/auth.setAccessToken

## Decisions Made
- Used `partialize` to exclude all actions and `isHydrated` from localStorage persistence -- only tokens, users, and isAuthenticated flags are persisted
- Used `onRehydrateStorage` callback to set `isHydrated` flag after rehydration -- prevents hydration mismatch in SSR/SSG scenarios
- Admin refresh response uses `response.data.data.token` (not `access_token`) -- matches OpenAPI RefreshTokenResponse schema which returns `token`, not `access_token`
- Admin client uses plain `axios.post` for refresh to avoid infinite interceptor loop
- Customer and agent clients redirect to portal-specific login pages on 401 (no token refresh) -- portals use single token, not access/refresh pair
- Login API functions use plain `axios` (not interceptored clients) since there is no token yet at login time

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 3 - Blocking] Fixed broken consumers after auth store rewrite**
- **Found during:** Task 1 (Create Zod schemas, API types, and Zustand auth store)
- **Issue:** Rewriting auth-store.ts changed its exported shape from `{ auth: { user, setUser, accessToken, ... } }` to flat slices `{ adminAccessToken, adminSetTokens, ... }`. Three existing files imported the old shape and would cause TypeScript errors: `main.tsx` (used `useAuthStore.getState().auth.reset()`), `sign-out-dialog.tsx` (used `const { auth } = useAuthStore()`), `user-auth-form.tsx` (used `auth.setUser()` and `auth.setAccessToken()`)
- **Fix:** Updated all three consumers to use the new flat slice API: `adminClearAuth()` in main.tsx and sign-out-dialog.tsx, `adminSetTokens()` in user-auth-form.tsx
- **Files modified:** website/src/main.tsx, website/src/components/sign-out-dialog.tsx, website/src/features/auth/sign-in/components/user-auth-form.tsx
- **Verification:** `npx tsc --noEmit` passes with zero errors
- **Committed in:** `2b2721b` (Task 1 commit)

---

**Total deviations:** 1 auto-fixed (1 blocking)
**Impact on plan:** Consumer fix was necessary to maintain TypeScript compilation. The plan stated these would be updated in later plans, but TypeScript errors are blocking and were handled immediately per deviation rules.

## Issues Encountered
- Zod 4 compatibility: The project uses Zod v4.3.6 but the plan's code examples used Zod 3 syntax. Verified all used APIs (z.enum, z.literal, z.string().email(), z.string().nullable(), z.object().refine()) work identically in Zod 4 -- no changes needed.
- Website is a git submodule: Discovered the website directory is tracked as a git submodule (mode 160000) in the main repo. All commits go to the website submodule's own git history.

## User Setup Required

None - no external service configuration required.

## Next Phase Readiness
- Auth data layer is complete and ready for UI integration in plans 01-02 through 01-04
- Admin login page (01-02) can import `adminLogin` from `@/api/auth` and use `adminSetTokens`/`adminSetUser` from auth store
- Customer/agent portal login pages can use their respective `customerLogin`/`agentLogin` functions
- Change password form can use `adminChangePassword` from auth API
- All Axios clients are ready for authenticated API calls in subsequent phases

## Self-Check: PASSED

All created files exist. Both commits verified. TypeScript compilation passes with zero errors.

---
*Phase: 01-auth-api-foundation*
*Completed: 2026-04-02*
