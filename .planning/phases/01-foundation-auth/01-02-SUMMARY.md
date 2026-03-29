---
phase: 01-foundation-auth
plan: "02"
subsystem: ui
tags: [react, zustand, axios, zod, tanstack-query, tanstack-router, rbac, jwt, typescript]

# Dependency graph
requires:
  - 01-01 (Vite + React + TanStack Router scaffold in dashboard/)
provides:
  - Zustand persisted auth store with 3-portal slices and isHydrated gate
  - Three Axios instances (admin/agent/customer) with JWT refresh queue
  - RBAC utility with hasPermission() and getAccessibleResources()
  - Zod schemas for all API auth responses
  - Singleton QueryClient at src/lib/queryClient.ts
  - Auth hooks: useAdminAuth, useAgentAuth, useCustomerAuth + context hooks
  - TanStack Router beforeLoad guards on _admin, _agentAuth, _customerAuth layouts
  - main.tsx with hydration-gated RouterProvider
  - src/api/auth.ts raw login functions for login pages
affects:
  - 01-03-app-shell (uses auth hooks for login forms and AppShell layout)
  - All admin/agent/customer route files (protected by guards)

# Tech tracking
tech-stack:
  added:
    - vitest@4 (dev)
    - "@vitest/coverage-v8" (dev)
  patterns:
    - Zustand slice pattern: AdminAuthSlice + AgentAuthSlice + CustomerAuthSlice combined in one store
    - Zustand persist with partialize (tokens/user only, no actions)
    - isHydrated flag via onRehydrateStorage callback — prevents redirect-on-reload bug
    - Vanilla getters (getAdminAccessToken etc.) for Axios interceptors outside React
    - Per-portal Axios instances with isolated isRefreshing + failedQueue
    - Refresh calls use plain axios (not portal client) to avoid interceptor loop
    - RBAC role hierarchy: superadmin(3) > admin(2) > teknisi(1)
    - hasPermission(role, resource, action) pure function — testable without React
    - RouterContext typed with AdminRole | null for role field
    - beforeLoad guard pattern in TanStack Router pathless layouts

key-files:
  created:
    - dashboard/src/store/types.ts
    - dashboard/src/store/slices/adminAuthSlice.ts
    - dashboard/src/store/slices/agentAuthSlice.ts
    - dashboard/src/store/slices/customerAuthSlice.ts
    - dashboard/src/lib/axios/admin-client.ts
    - dashboard/src/lib/axios/agent-client.ts
    - dashboard/src/lib/axios/customer-client.ts
    - dashboard/src/lib/schemas/auth.ts
    - dashboard/src/lib/queryClient.ts
    - dashboard/src/lib/rbac.ts
    - dashboard/src/lib/rbac.test.ts
    - dashboard/src/hooks/useAuth.ts
    - dashboard/src/hooks/useAdminAuth.ts
    - dashboard/src/hooks/useAgentAuth.ts
    - dashboard/src/hooks/useCustomerAuth.ts
    - dashboard/src/api/auth.ts
    - dashboard/src/vite-env.d.ts
  modified:
    - dashboard/src/store/index.ts (replaced export {} stub with full implementation)
    - dashboard/src/routes/__root.tsx (updated RouterContext with AdminRole + ReactQueryDevtools)
    - dashboard/src/routes/_admin/route.tsx (added beforeLoad auth guard)
    - dashboard/src/routes/agent/_agentAuth/route.tsx (added beforeLoad auth guard)
    - dashboard/src/routes/customer/_customerAuth/route.tsx (added beforeLoad auth guard)
    - dashboard/src/main.tsx (hydration gate + auth context injection + queryClient singleton import)
    - dashboard/package.json (added vitest + @vitest/coverage-v8 dev dependencies)

key-decisions:
  - "isHydrated flag in Zustand store gates RouterProvider rendering — prevents redirect-on-reload when localStorage tokens not yet loaded"
  - "Each Axios portal client has isolated isRefreshing + failedQueue — no cross-portal state sharing"
  - "Refresh calls use plain import axios — not portal client — to prevent interceptor loop"
  - "RBAC hasPermission() is pure function taking (role, resource, action) — no React dependency enables direct unit testing"
  - "RouterContext.adminAuth.role typed as AdminRole | null (not string | null) for type-safe RBAC in route guards"
  - "vite-env.d.ts added (Rule 1 fix) — missing file caused import.meta.env TS2339 error in tsc -b"
  - "vitest installed as dev dependency (Rule 3 fix) — was missing from scaffold but required for RBAC TDD"

metrics:
  duration: "11 minutes"
  completed: "2026-03-30"
  tasks_completed: 5
  tasks_total: 5
  files_created: 17
  files_modified: 6
---

# Phase 1 Plan 02: Auth System Summary

Zustand auth store with 3-portal slices, persisted tokens, JWT refresh queue Axios clients, pure RBAC with 30 unit tests, TanStack Router beforeLoad guards, and hydration-gated main.tsx.

## Tasks Completed

| Task | Name | Commit | Status |
|------|------|--------|--------|
| 1 | Zustand Auth Store — slices + combined store + vanilla getters | 7673f4e | Complete |
| 2 | Three Axios instances with JWT refresh queue | c0afb55 | Complete |
| 3 | Zod schemas, QueryClient, RBAC utility + unit tests | 6657fe2 | Complete |
| 4 | Auth hooks, RouterContext, route guards, main.tsx | 7b51b42 | Complete |
| 5 | Raw auth API functions (`src/api/auth.ts`) | 68d3344 | Complete |

## Verification Results

- TypeScript: `tsc --noEmit` — zero errors
- RBAC tests: 30/30 passed (vitest run)
- Vite build: `npm run build` — succeeded, 205 modules, 361 kB main bundle

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 1 - Bug] Missing `vite-env.d.ts` caused `import.meta.env` TypeScript error**
- **Found during:** Task 4 verification (`npm run build`)
- **Issue:** `tsc -b` reported TS2339 "Property 'env' does not exist on type 'ImportMeta'" on `src/routes/__root.tsx:30` and `src/store/index.ts:47`
- **Fix:** Created `dashboard/src/vite-env.d.ts` with `/// <reference types="vite/client" />`
- **Files modified:** `dashboard/src/vite-env.d.ts` (new)
- **Commit:** 7b51b42

**2. [Rule 3 - Blocking] Vitest not installed — required for Task 3 RBAC test verification**
- **Found during:** Task 3 verification
- **Issue:** `node_modules/.bin/vitest` did not exist — vitest was not in package.json devDependencies
- **Fix:** `npm install --save-dev vitest @vitest/coverage-v8`
- **Files modified:** `dashboard/package.json`, `dashboard/package-lock.json`
- **Commit:** 6657fe2

## Known Stubs

None — all files wire real data. Login page UI (using `useAdminLogin`, `useAgentLogin`, `useCustomerLogin` hooks from this plan) will be implemented in Plan 01-03 (App Shell + Login Forms).

## Self-Check: PASSED

- All 17 created files exist on disk
- All 5 task commits found in git log (7673f4e, c0afb55, 6657fe2, 7b51b42, 68d3344)
- TypeScript: zero errors
- RBAC tests: 30/30 passed
- Vite build: succeeded
