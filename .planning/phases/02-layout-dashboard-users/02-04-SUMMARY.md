---
plan: 02-04
phase: 02-layout-dashboard-users
status: complete
completed: 2026-04-03
---

# Plan 02-04: Dashboard Overview Page

## Summary

Admin dashboard completely rebuilt as single-page overview with KPI widgets, router health cards, and recent activity feed - all wired to real API data.

## What Was Built

- `website/src/features/dashboard/components/kpi-card.tsx` — Reusable KPICard with UI-SPEC typography (title: text-sm font-semibold muted, value: text-2xl font-semibold, trend: text-sm muted) and loading skeleton
- `website/src/features/dashboard/components/revenue-chart.tsx` — BarChart with Recharts using /reports/summary total_revenue data, Rp formatting, chart-1 CSS variable colors
- `website/src/features/dashboard/components/router-health-cards.tsx` — Grid of router cards with click-to-select (calls selectRouter API + updates store directly), status badges (online=default, offline=destructive, unknown=secondary), active ring-2 ring-primary, empty state with "No Routers Configured"
- `website/src/features/dashboard/components/recent-activity-feed.tsx` — useUsers(5, 0) combined-data approach showing last 5 registered users with relative timestamps
- `website/src/features/dashboard/index.tsx` — Complete rewrite: single page (no Tabs, TopNav, Search, ConfigDrawer, ProfileDropdown), 4 KPI cards grid, RevenueChart, RouterHealthCards + RecentActivityFeed in 7-col layout

## Verification

- TypeScript `npx tsc --noEmit` passes with zero errors
- No Tabs, TopNav, Search, or ConfigDrawer imports in dashboard/index.tsx
- useReportSummary drives TotalCustomers, ActiveSubscriptions, MonthlyRevenue KPI cards
- useRouters drives ActiveRouters KPI card and RouterHealthCards
- useUsers(5, 0) drives RecentActivityFeed

## Self-Check: PASSED
