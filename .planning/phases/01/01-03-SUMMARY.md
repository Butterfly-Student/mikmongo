---
phase: 01-auth-api-foundation
plan: 03
subsystem: auth
tags: [react, typescript, tanstack-query, zustand, zod, shadcn-ui]

# Dependency graph
requires:
  - phase: 01-01
    provides: "Zustand auth store, Zod schemas, API auth functions, Axios clients"
  - phase: 01-02
    provides: "AuthLayout with MikMongo branding, admin login pattern, Indonesian text conventions"
provides:
  - Customer portal login page with identifier-based auth matching OpenAPI PortalLoginRequest
  - Agent portal login page with username-based auth matching OpenAPI AgentPortalLoginRequest
  - Customer auth hooks (useCustomerLogin, useCustomerLogout, useCustomerUser)
  - Agent auth hooks (useAgentLogin, useAgentLogout, useAgentUser)
affects: [01-04]

# Tech tracking
tech-stack:
  added: []
  patterns:
    - "Portal login pages follow same Card + AuthLayout pattern as admin"
    - "Login pages use direct useState + async onSubmit (not TanStack Query mutations) for form components"
    - "Separate TanStack Query mutation hooks provided for programmatic use"

key-files:
  created:
    - website/src/features/auth/customer-login/index.tsx
    - website/src/features/auth/agent-login/index.tsx
    - website/src/hooks/use-customer-auth.ts
    - website/src/hooks/use-agent-auth.ts
  modified: []

key-decisions:
  - "Customer login uses identifier field (not email) matching OpenAPI PortalLoginRequest schema"
  - "Agent login uses username field matching OpenAPI AgentPortalLoginRequest schema"
  - "Both portals store single token (not access+refresh) per their OpenAPI response shapes"
  - "Forgot password link shows toast info directing user to contact admin (no backend endpoint)"

patterns-established:
  - "Three-portal login consistency: same AuthLayout wrapper, Card structure, form pattern, Indonesian text"
  - "Portal-specific login hooks follow same naming convention: use{Portal}Login/Logout/User"

requirements-completed: [AUTH-06, AUTH-07]

# Metrics
duration: 4min
completed: 2026-04-02
---

# Phase 01 Plan 03: Customer & Agent Login Pages Summary

**Customer login with identifier field (email/phone/username) and agent login with username field, both using single-token auth, Indonesian text, MikMongo branding, and TanStack Query auth hooks**

## Performance

- **Duration:** 4 min
- **Started:** 2026-04-02T22:43:12Z
- **Completed:** 2026-04-02T22:47:52Z
- **Tasks:** 2
- **Files modified:** 4

## Accomplishments
- Customer portal login page with `identifier` field matching OpenAPI `PortalLoginRequest` (accepts email, phone, or username)
- Agent portal login page with `username` field matching OpenAPI `AgentPortalLoginRequest`
- Both pages follow established admin login pattern: AuthLayout, Card wrapper, Zod form validation, PasswordInput, Remember me checkbox, Forgot password link
- Customer and agent auth hooks (login/logout/user) using TanStack Query mutations for programmatic use

## Task Commits

Each task was committed atomically:

1. **Task 1: Customer portal login page and auth hooks** - `942b38e` (feat)
2. **Task 2: Agent portal login page and auth hooks** - `55f674c` (feat)

_Note: Commits are in the website submodule (git submodule)._

## Files Created/Modified
- `website/src/features/auth/customer-login/index.tsx` - Customer portal login with identifier field, Indonesian text, MikMongo branding
- `website/src/features/auth/agent-login/index.tsx` - Agent portal login with username field, Indonesian text, MikMongo branding
- `website/src/hooks/use-customer-auth.ts` - useCustomerLogin, useCustomerLogout, useCustomerUser hooks
- `website/src/hooks/use-agent-auth.ts` - useAgentLogin, useAgentLogout, useAgentUser hooks

## Decisions Made
- Customer login uses `identifier` field (not `email`) -- matches OpenAPI `PortalLoginRequest` which accepts email, phone, or username in a single field
- Agent login uses `username` field -- matches OpenAPI `AgentPortalLoginRequest` which uses username-only authentication
- Both portal login pages use direct `useState` + `async onSubmit` pattern in the component (matching admin pattern from Plan 02), while also providing TanStack Query mutation hooks for programmatic use
- Forgot password link shows `toast.info('Hubungi admin untuk reset password')` since there is no backend forgot-password endpoint

## Deviations from Plan

None - plan executed exactly as written.

## Issues Encountered
- Worktree branch did not include the `website` directory (it was on an older commit). Resolved by merging the `fully` branch into the worktree to get the submodule reference.
- The `website` directory is tracked as a git submodule (mode 160000) in the main repo but has no `.gitmodules` file. Commits go to the website's own git history.

## Known Stubs
- "Remember me" checkbox is present but non-functional (no backend support for persistent sessions)
- "Forgot password?" link shows a toast directing to admin rather than navigating to a dedicated page (no forgot-password endpoint)

## User Setup Required

None - no external service configuration required.

## Next Phase Readiness
- All three portal login pages are complete (admin, customer, agent)
- Auth hooks available for all three portals (login, logout, user access)
- Route guards (Plan 01-04) can import CustomerLogin and AgentLogin components for their respective login routes
- No blockers for Plan 01-04

## Self-Check: PASSED

All created files exist. Both commits verified. TypeScript compilation passes with zero errors.

---
*Phase: 01-auth-api-foundation*
*Completed: 2026-04-02*
