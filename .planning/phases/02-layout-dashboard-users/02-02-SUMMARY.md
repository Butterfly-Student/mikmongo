---
plan: 02-02
phase: 02-layout-dashboard-users
status: complete
completed: 2026-04-03
---

# Plan 02-02: Sidebar Navigation Restructure

## Summary

Sidebar navigation restructured for ISP management domain. MikMongo branding applied. Disabled items
show "Coming soon" tooltip.

## What Was Built

- `website/src/components/layout/types.ts` — added `disabled?: boolean` to BaseNavItem, removed Team type and teams from SidebarData
- `website/src/components/layout/data/sidebar-data.ts` — full ISP nav groups: Overview (Dashboard), Management (Users, Customers, Routers, Subscriptions), Billing, Sales, MikroTik, Monitor, Reports, Settings. Only Dashboard and Users are active; 15+ items set disabled: true
- `website/src/components/layout/nav-group.tsx` — disabled items render without Link, with opacity-50 cursor-not-allowed, wrapped in Tooltip showing "Coming soon"
- `website/src/components/layout/app-title.tsx` — replaced "Shadcn-Admin / Vite + ShadcnUI" with "MikMongo / ISP Management"
- `website/src/components/layout/app-sidebar.tsx` — removed TeamSwitcher import and usage; added AppTitle; RouterSelector placeholder for Plan 03

## Verification

- TypeScript `npx tsc --noEmit` passes with zero errors
- ISP nav groups present with Overview, Management, Billing, etc.
- Disabled items use opacity-50 and "Coming soon" tooltip
- MikMongo branding in app-title
- No TeamSwitcher in app-sidebar

## Self-Check: PASSED
