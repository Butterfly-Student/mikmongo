---
plan: 02-01
phase: 02-layout-dashboard-users
status: complete
completed: 2026-04-03
---

# Plan 02-01: Data Layer — Schemas, Store, API, Hooks

## Summary

Complete data layer established for all Phase 2 UI plans. All 11 files created with Zod-validated
types matching the OpenAPI contract exactly.

## What Was Built

**Zod Schemas:**
- `website/src/lib/schemas/router.ts` — RouterResponseSchema, RouterListResponseSchema, SelectedRouterResponseSchema + MetaSchema
- `website/src/lib/schemas/user.ts` — UserResponseSchema, CreateUserRequestSchema, CreateUserFormSchema (with confirm_password refine), UserListResponseSchema
- `website/src/lib/schemas/report.ts` — ReportSummarySchema, ReportSummaryResponseSchema

**Zustand Store:**
- `website/src/stores/router-store.ts` — selectedRouterId, selectedRouterName, isHydrated; persist + partialize (only data fields) + onRehydrateStorage hydration gate; follows auth-store pattern exactly

**API Functions:**
- `website/src/api/router.ts` — listRouters, getSelectedRouter, selectRouter (all using adminClient + Zod parse)
- `website/src/api/user.ts` — listUsers, getUser, createUser, deleteUser (all using adminClient + Zod parse)
- `website/src/api/report.ts` — getReportSummary (using adminClient + Zod parse)
- `website/src/api/types.ts` — extended with RouterResponse, UserResponse, CreateUserRequest, CreateUserFormValues, ReportSummary, Meta types

**TanStack Query Hooks:**
- `website/src/hooks/use-routers.ts` — useRouters (staleTime 2min), useSelectRouter (updates store + invalidates ['routers'] and ['report-summary'])
- `website/src/hooks/use-report-summary.ts` — useReportSummary (staleTime 5min, queryKey ['report-summary', from, to])
- `website/src/hooks/use-users.ts` — useUsers, useUser, useCreateUser, useDeleteUser (all with query invalidation)

## Verification

- TypeScript `npx tsc --noEmit` passes with zero errors
- Router store has persist + partialize + onRehydrateStorage hydration gate
- useSelectRouter invalidates both ['routers'] and ['report-summary'] on success (per D-04)
- All API functions use adminClient and Zod schema parse

## Key Files Created

- website/src/lib/schemas/router.ts
- website/src/lib/schemas/user.ts
- website/src/lib/schemas/report.ts
- website/src/stores/router-store.ts
- website/src/api/router.ts
- website/src/api/user.ts
- website/src/api/report.ts
- website/src/api/types.ts (modified)
- website/src/hooks/use-routers.ts
- website/src/hooks/use-report-summary.ts
- website/src/hooks/use-users.ts

## Self-Check: PASSED
