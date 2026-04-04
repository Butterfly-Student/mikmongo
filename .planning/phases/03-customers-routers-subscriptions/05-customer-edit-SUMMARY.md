---
phase: 03-customers-routers-subscriptions
plan: 05
subsystem: customers
tags: [customers, edit, dialog, mutation, tanstack-query]
dependency_graph:
  requires: [03-02-customers]
  provides: [customer-edit-capability]
  affects: [customers-page]
tech_stack:
  added: []
  patterns: [useMutation with invalidateQueries, pre-populated form dialog, useEffect reset on prop change]
key_files:
  created:
    - website/src/features/customers/components/edit-customer-dialog.tsx
  modified:
    - website/src/hooks/use-customers.ts
    - website/src/features/customers/data/columns.tsx
    - website/src/features/customers/index.tsx
decisions:
  - useEffect to reset form when customer prop changes, ensuring fresh pre-population on each open
  - Password field defaults to empty string and is only sent in payload if user fills it
  - Edit action positioned between activate/deactivate and delete in dropdown (logical UX grouping)
metrics:
  duration: 2m28s
  completed: 2026-04-04T11:58:27Z
  tasks_completed: 1
  files_changed: 4
---

# Phase 03 Plan 05: Customer Edit Summary

## One-liner

Customer edit dialog with pre-populated form fields and `useUpdateCustomer` mutation hook calling PUT /customers/{id}.

## What Was Built

Closed verification gap CUST-03 (admin can update customer details). The `updateCustomer` API function already existed but had no hook or UI surface.

### Components

**`useUpdateCustomer` hook** (`website/src/hooks/use-customers.ts`)
- Calls `updateCustomer(id, data)` via `useMutation`
- Invalidates `['customers']` query on success
- Shows success/error toasts

**`EditCustomerDialog`** (`website/src/features/customers/components/edit-customer-dialog.tsx`)
- Props: `{ customer: CustomerResponse | null; open: boolean; onOpenChange }`
- Pre-populates full_name, phone, email, address, username from `customer` prop
- `useEffect` resets form when customer prop changes
- Password field defaults empty — only included in payload if user fills it
- Submit calls `updateCustomer({ id: customer.id, data: payload })`
- Pencil icon in dialog title

**Columns (`website/src/features/customers/data/columns.tsx`)**
- Added `onEdit: (customer: CustomerResponse) => void` to `ColumnActions` interface
- Added `Pencil` to lucide-react imports
- Edit dropdown item appears between activate/deactivate and the separator+delete

**Customers page (`website/src/features/customers/index.tsx`)**
- Added `editTarget` state (`CustomerResponse | null`)
- Imported `EditCustomerDialog`
- Passed `onEdit: (customer) => setEditTarget(customer)` to `createCustomerColumns`
- Renders `<EditCustomerDialog customer={editTarget} open={!!editTarget} onOpenChange=...>`

## Deviations from Plan

None — plan executed exactly as written.

## Known Stubs

None — all wired to real API.

## Self-Check

- [x] `website/src/features/customers/components/edit-customer-dialog.tsx` exists
- [x] `useUpdateCustomer` exported from `use-customers.ts`
- [x] `onEdit` in `ColumnActions` interface in `columns.tsx`
- [x] `editTarget`, `EditCustomerDialog`, `onEdit` present in `index.tsx`
- [x] TypeScript compilation passes cleanly
- [x] Task commit `35469da` exists in website repo

## Self-Check: PASSED
