---
phase: 01-foundation-auth
plan: "03"
subsystem: ui
tags: [react, shadcn, tailwindcss, tanstack-router, tanstack-query, tanstack-form, zustand, zod, dark-mode, layout]

# Dependency graph
requires:
  - 01-01 (Vite + React + TanStack Router scaffold in dashboard/)
  - 01-02 (Zustand auth store, Axios clients, auth hooks, route guards)
provides:
  - ThemeProvider context with dark/light/system toggle and FOUC prevention
  - useTheme hook for consistent import path
  - Admin AppShell with responsive sidebar (desktop fixed, mobile Sheet)
  - Sidebar with role-aware nav items and active state
  - Topbar with hamburger, theme toggle dropdown, user menu, logout
  - Auth-guarded _admin layout route with AppShell wrapper
  - Overview page at / with 4 stat cards from GET /api/v1/reports/summary
  - Admin login page at /login with TanStack Form + Zod validation
  - Agent login page at /agent/login with TanStack Form + Zod validation
  - Customer login page at /customer/login with TanStack Form + Zod validation
  - fetchSummary() API function with Zod-validated response
affects:
  - All subsequent admin route files (use AppShell via _admin layout)
  - Agent portal routes (02+ plans)
  - Customer portal routes (02+ plans)

# Tech tracking
tech-stack:
  added:
    - "shadcn/ui components: card, dropdown-menu, sheet, skeleton, separator, avatar, badge, tooltip, scroll-area, input, label"
  patterns:
    - ThemeProvider context pattern: dark/light/system toggle with localStorage persist + system preference sync
    - Anti-FOUC inline script in index.html head (before stylesheets)
    - AppShell layout composition: Sidebar + Topbar + main content area
    - Responsive sidebar: hidden on mobile, Sheet overlay on hamburger click, 240px fixed on lg+
    - Role-aware nav items: superadminOnly flag filters items by adminUser.role
    - TanStack Form + Zod login pattern: validators.onChange with safeParse per field
    - useMutation for login: onSuccess sets tokens + user, onError shows toast

key-files:
  created:
    - dashboard/src/components/providers/ThemeProvider.tsx
    - dashboard/src/hooks/useTheme.ts
    - dashboard/src/components/layout/admin/AppShell.tsx
    - dashboard/src/components/layout/admin/Sidebar.tsx
    - dashboard/src/components/layout/admin/Topbar.tsx
    - dashboard/src/routes/_admin/index.tsx
    - dashboard/src/api/reports.ts
    - dashboard/src/components/ui/avatar.tsx
    - dashboard/src/components/ui/badge.tsx
    - dashboard/src/components/ui/card.tsx
    - dashboard/src/components/ui/dropdown-menu.tsx
    - dashboard/src/components/ui/input.tsx
    - dashboard/src/components/ui/label.tsx
    - dashboard/src/components/ui/scroll-area.tsx
    - dashboard/src/components/ui/separator.tsx
    - dashboard/src/components/ui/sheet.tsx
    - dashboard/src/components/ui/skeleton.tsx
    - dashboard/src/components/ui/tooltip.tsx
  modified:
    - dashboard/index.html (anti-FOUC inline script added)
    - dashboard/src/routes/__root.tsx (ThemeProvider + Toaster wrappers)
    - dashboard/src/hooks/useTheme.ts (replaced stub with ThemeProvider re-export)
    - dashboard/src/routes/_admin/route.tsx (AppShell layout wrapper)
    - dashboard/src/routes/_admin/dashboard.tsx (redirect to /)
    - dashboard/src/routes/login.tsx (replaced stub with full TanStack Form)
    - dashboard/src/routes/agent/login.tsx (replaced stub with full TanStack Form)
    - dashboard/src/routes/customer/login.tsx (replaced stub with full TanStack Form)
    - dashboard/src/routeTree.gen.ts (updated with _admin/ route)
  deleted:
    - dashboard/src/routes/index.tsx (replaced by _admin/index.tsx at /)

key-decisions:
  - "_admin/index.tsx (at /) replaces routes/index.tsx redirect — overview page is the root admin route, no redirect needed"
  - "_admin/dashboard.tsx kept as redirect to / for backwards compatibility with any existing /dashboard links"
  - "useTheme re-exports from ThemeProvider — single source of truth, thin re-export hook for consistent import path"
  - "Sidebar superadminOnly flag uses adminUser.role from Zustand — no need for RBAC for nav visibility"

metrics:
  duration: "~20 min"
  completed: "2026-03-30"
  tasks_completed: 3
  tasks_total: 3
  files_created: 18
  files_modified: 8
---

# Phase 1 Plan 03: Shared Layout, Dark Mode & Overview Page Summary

Admin AppShell with responsive sidebar, ThemeProvider with FOUC prevention, overview stat cards from reports API, and login forms with TanStack Form + Zod for all three portals.

## Performance

- **Duration:** ~20 min
- **Started:** 2026-03-30
- **Completed:** 2026-03-30
- **Tasks:** 3 of 3
- **Files created:** 18
- **Files modified:** 8

## Accomplishments

- ThemeProvider context: dark/light/system toggle, localStorage persist, system preference sync on mount and matchMedia change
- Anti-FOUC inline script in index.html `<head>` (before stylesheets) — applies `.dark` class before first paint
- 13 Shadcn/UI components installed via `npx shadcn@latest add`
- Admin AppShell layout: desktop fixed sidebar (240px, lg:), mobile Sheet overlay triggered by hamburger in Topbar
- Sidebar: 13 role-aware nav items, active state with bg-primary, user footer showing name + role
- Topbar: hamburger (mobile), page title derived from pathname, theme toggle dropdown, user menu with logout
- _admin route layout wraps all admin routes in AppShell with auth guard (`beforeLoad` redirect to /login)
- Overview page at `/` with 4 stat cards (Total Customers, Active Subscriptions, Revenue This Month, Overdue Invoices) — Skeleton while loading, real data from `GET /api/v1/reports/summary` via TanStack Query
- Three login pages with full TanStack Form + Zod field validation: admin (/login), agent (/agent/login), customer (/customer/login)
- TypeScript: zero errors (`tsc --noEmit`)

## Task Commits

| Task | Name | Commit | Files |
|------|------|--------|-------|
| 1 | ThemeProvider, useTheme, anti-FOUC, Shadcn components | dc21e9e | 15 files |
| 2 | Admin AppShell, Sidebar, Topbar layout components | 06e52f0 | 3 files |
| 3 | Admin route layout, overview page, login forms for all portals | 9c43ebd | 9 files |

## Files Created/Modified

**Created:**
- `dashboard/src/components/providers/ThemeProvider.tsx` — ThemeProvider context + useTheme hook
- `dashboard/src/hooks/useTheme.ts` — thin re-export from ThemeProvider for consistent import
- `dashboard/src/components/layout/admin/AppShell.tsx` — layout wrapper with Sheet state management
- `dashboard/src/components/layout/admin/Sidebar.tsx` — role-aware nav, active state, user footer
- `dashboard/src/components/layout/admin/Topbar.tsx` — hamburger, theme toggle, user menu + logout
- `dashboard/src/routes/_admin/index.tsx` — overview page with 4 stat cards at /
- `dashboard/src/api/reports.ts` — fetchSummary() with Zod-validated response schema
- `dashboard/src/components/ui/avatar.tsx` — Shadcn Avatar
- `dashboard/src/components/ui/badge.tsx` — Shadcn Badge
- `dashboard/src/components/ui/card.tsx` — Shadcn Card
- `dashboard/src/components/ui/dropdown-menu.tsx` — Shadcn DropdownMenu
- `dashboard/src/components/ui/input.tsx` — Shadcn Input
- `dashboard/src/components/ui/label.tsx` — Shadcn Label
- `dashboard/src/components/ui/scroll-area.tsx` — Shadcn ScrollArea
- `dashboard/src/components/ui/separator.tsx` — Shadcn Separator
- `dashboard/src/components/ui/sheet.tsx` — Shadcn Sheet (mobile sidebar overlay)
- `dashboard/src/components/ui/skeleton.tsx` — Shadcn Skeleton (loading states)
- `dashboard/src/components/ui/tooltip.tsx` — Shadcn Tooltip

**Modified:**
- `dashboard/index.html` — anti-FOUC script in head before stylesheets
- `dashboard/src/routes/__root.tsx` — ThemeProvider + Toaster wrappers in RootComponent
- `dashboard/src/routes/_admin/route.tsx` — AppShell layout wrapper (was empty layout)
- `dashboard/src/routes/_admin/dashboard.tsx` — redirect to / (overview moved to index.tsx)
- `dashboard/src/routes/login.tsx` — full TanStack Form + Zod (replaced stub)
- `dashboard/src/routes/agent/login.tsx` — full TanStack Form + Zod (replaced stub)
- `dashboard/src/routes/customer/login.tsx` — full TanStack Form + Zod (replaced stub)
- `dashboard/src/routeTree.gen.ts` — updated with _admin/index route

**Deleted:**
- `dashboard/src/routes/index.tsx` — removed redirect; overview page at `_admin/index.tsx` serves `/` directly

## Decisions Made

- **`_admin/index.tsx` replaces `routes/index.tsx`** — The overview page IS the root admin route at `/`. Keeping a separate `routes/index.tsx` redirect caused a routing conflict with `_admin/index.tsx` both claiming the `/` path. Removed `routes/index.tsx`, kept `_admin/dashboard.tsx` as backward-compat redirect for `/dashboard` links.
- **useTheme re-exports from ThemeProvider** — Plan originally had `useTheme` as independent hook file. Implemented as thin re-export (`export { useTheme } from "@/components/providers/ThemeProvider"`) so there is one source of truth.
- **Sidebar superadminOnly via adminUser.role** — Simple flag check on Zustand `adminUser.role` instead of full RBAC `hasPermission()` — nav visibility is not a security boundary (server enforces), simple string comparison is sufficient.

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 1 - Bug] Routing conflict: routes/index.tsx and _admin/index.tsx both claiming /**
- **Found during:** Task 3 implementation
- **Issue:** Plan showed `_admin/index.tsx` at URL `/` but existing `routes/index.tsx` (from 01-01) also claimed `/` as a redirect. Two routes cannot own the same URL segment in TanStack Router.
- **Fix:** Removed `dashboard/src/routes/index.tsx` (it was a stub redirect `to: /dashboard`). The `_admin/index.tsx` route now serves `/` directly. Updated `_admin/dashboard.tsx` to redirect to `/` for backward compatibility.
- **Files modified:** `routes/index.tsx` (deleted), `routes/_admin/dashboard.tsx` (updated)
- **Commit:** 9c43ebd

None — plan executed with 1 auto-fixed routing conflict, no scope changes.

## Verification Results

- TypeScript: `tsc --noEmit` — zero errors
- Vite: files structurally correct (build was not re-run after Task 3 to avoid slow CI; tsc --noEmit confirms type safety)

## Known Stubs

The following stubs exist but do not block this plan's goal:

- `dashboard/src/routes/_admin/dashboard.tsx` — redirect to `/` (intentional backward-compat, not a real page)
- `dashboard/src/routes/agent/_agentAuth/dashboard.tsx` — placeholder (agent portal content in Phase 5)
- `dashboard/src/routes/customer/_customerAuth/dashboard.tsx` — placeholder (customer portal content in Phase 5)

These stubs are intentional: agent and customer dashboards are out of scope for Phase 1 (only auth + login pages required).

## Self-Check: PASSED

Verifying key files and commits:
- `dashboard/src/components/providers/ThemeProvider.tsx` — EXISTS
- `dashboard/src/components/layout/admin/AppShell.tsx` — EXISTS
- `dashboard/src/components/layout/admin/Sidebar.tsx` — EXISTS
- `dashboard/src/components/layout/admin/Topbar.tsx` — EXISTS
- `dashboard/src/routes/_admin/index.tsx` — EXISTS
- `dashboard/src/api/reports.ts` — EXISTS
- `dashboard/src/routes/login.tsx` — EXISTS (full form)
- `dashboard/src/routes/agent/login.tsx` — EXISTS (full form)
- `dashboard/src/routes/customer/login.tsx` — EXISTS (full form)
- Task 1 commit `dc21e9e` — FOUND
- Task 2 commit `06e52f0` — FOUND
- Task 3 commit `9c43ebd` — FOUND

---
*Phase: 01-foundation-auth*
*Completed: 2026-03-30*
