---
plan: 02-03
phase: 02-layout-dashboard-users
status: complete
completed: 2026-04-03
---

# Plan 02-03: Router Selector + Ping Display

## Summary

Router selector dropdown added to sidebar header. Real-time ping via WebSocket now shows in header.

## What Was Built

- `website/src/hooks/use-ping.ts` — WebSocket hook with 4 states (connecting/connected/error/unavailable), reconnection (max 5 attempts, 3s delay), proper cleanup on unmount/router change
- `website/src/components/layout/sidebar/router-selector.tsx` — Select dropdown with router list from useRouters hook, status badges (online=default, offline=destructive, unknown=secondary), loading skeleton, API sync on mount via getSelectedRouter()
- `website/src/components/layout/sidebar/ping-display.tsx` — Activity icon + latency text, text-destructive on error, animate-pulse when connecting
- `website/src/components/layout/header.tsx` — PingDisplay added between SidebarTrigger separator and page children
- `website/src/components/layout/app-sidebar.tsx` — RouterSelector added to SidebarHeader below AppTitle

## Verification

- TypeScript `npx tsc --noEmit` passes with zero errors
- RouterSelector uses useRouters + useSelectRouter hooks (from plan 02-01)
- useSelectRouter handles store update + query invalidation (D-04)
- WebSocket URL: `/api/v1/routers/{id}/monitor/ws/ping?address=8.8.8.8&token={token}`
- PingDisplay reads selectedRouterId from useRouterStore

## Self-Check: PASSED
