---
phase: 04-billing-payments
plan: 02
subsystem: billing/invoices
tags: [invoices, billing, datatable, sheet, sidebar]
dependency_graph:
  requires: [04-01]
  provides: [invoice-page, invoice-route, invoice-sidebar]
  affects: [sidebar-nav, billing-feature]
tech_stack:
  added: []
  patterns:
    - DataTable with DataTableToolbar faceted filter
    - Sheet side panel for detail view
    - ConfirmDialog for destructive/bulk actions
    - Feature module pattern (data/schema, data/columns, components/, index.tsx)
key_files:
  created:
    - website/src/features/billing/invoices/data/schema.ts
    - website/src/features/billing/invoices/data/columns.tsx
    - website/src/features/billing/invoices/components/invoice-table.tsx
    - website/src/features/billing/invoices/components/invoice-detail-sheet.tsx
    - website/src/features/billing/invoices/components/invoice-generation-trigger.tsx
    - website/src/features/billing/invoices/index.tsx
    - website/src/routes/_authenticated/invoices/index.tsx
  modified:
    - website/src/components/layout/data/sidebar-data.ts
decisions:
  - Sidebar Invoices link enabled at /invoices; Payments and Cash remain disabled pending plans 03-04
  - InvoiceDetailSheet uses overflow-y-auto for scrollable long content
  - Invoice status filtering uses DataTableFacetedFilter via columns filterFn for multi-select
metrics:
  duration: 150s
  completed: 2026-04-05
  tasks: 2
  files: 8
---

# Phase 04 Plan 02: Invoice Management Page Summary

Invoice management page with DataTable, faceted status filter, detail side sheet, and monthly billing generation trigger — delivering the core billing visibility for the admin portal.

## Tasks Completed

| Task | Name | Commit | Files |
|------|------|--------|-------|
| 1 | Create invoice feature module | 2b11bee | 6 new files |
| 2 | Create invoice route and enable sidebar nav | c95cc88 | 1 new, 1 modified |

## What Was Built

- **Invoice feature module** at `website/src/features/billing/invoices/` with the full standard structure:
  - `data/schema.ts`: Re-exports `InvoiceResponse`, defines `invoiceStatuses` with Indonesian labels and Tailwind classes, and `overdueOptions`
  - `data/columns.tsx`: TanStack Table `ColumnDef<InvoiceResponse>[]` — invoice number, customer ID, dates with red highlighting for overdue deadlines, total amount formatted as Rp, status badge
  - `components/invoice-table.tsx`: Full-featured DataTable using `DataTableToolbar` with global search and faceted status filter, `DataTablePagination`
  - `components/invoice-detail-sheet.tsx`: Right-side Sheet at `sm:max-w-[540px]` showing billing info, line-item amounts (conditional tax/discount/late fee), status badge, payment history placeholder
  - `components/invoice-generation-trigger.tsx`: "Buat Tagihan Bulanan" button with `ConfirmDialog` calling `useTriggerMonthlyBilling()`
  - `index.tsx`: Page composition with loading skeletons, empty state with FileText icon, and sheet state management

- **Route** at `/_authenticated/invoices/` wired to `InvoicesPage`
- **Sidebar** Invoices item enabled at `/invoices`; Payments and Cash remain `disabled: true`

## Deviations from Plan

None — plan executed exactly as written.

## Known Stubs

- `invoice-detail-sheet.tsx` Section 4 "Riwayat Pembayaran" shows placeholder "Lihat halaman pembayaran untuk detail" — intentional per plan spec; linked payments not embedded in API invoice response. Future plan to enhance.

## Self-Check: PASSED

- `website/src/features/billing/invoices/index.tsx` — FOUND
- `website/src/features/billing/invoices/components/invoice-detail-sheet.tsx` — FOUND
- `website/src/features/billing/invoices/components/invoice-generation-trigger.tsx` — FOUND
- `website/src/routes/_authenticated/invoices/index.tsx` — FOUND
- `website/src/components/layout/data/sidebar-data.ts` updated — FOUND
- Commits `2b11bee`, `c95cc88` — FOUND (verified via git log)
