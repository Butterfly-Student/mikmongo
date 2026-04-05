---
phase: 04-billing-payments
plan: "03"
subsystem: billing/payments
tags: [payments, billing, datatable, dialogs, faceted-filters]
dependency_graph:
  requires: [04-01]
  provides: [payment-management-page]
  affects: [sidebar-nav]
tech_stack:
  added: []
  patterns:
    - DataTable with faceted filters and date-range filter
    - Conditional action menu based on payment status
    - ConfirmDialog reuse for confirm/refund actions
    - Custom AlertDialog for reject with Textarea
    - createPaymentColumns factory accepting action callbacks
key_files:
  created:
    - website/src/features/billing/payments/data/schema.ts
    - website/src/features/billing/payments/data/columns.tsx
    - website/src/features/billing/payments/components/payment-table.tsx
    - website/src/features/billing/payments/components/payment-action-menu.tsx
    - website/src/features/billing/payments/components/confirm-payment-dialog.tsx
    - website/src/features/billing/payments/components/reject-payment-dialog.tsx
    - website/src/features/billing/payments/components/refund-payment-dialog.tsx
    - website/src/features/billing/payments/components/date-range-filter.tsx
    - website/src/features/billing/payments/index.tsx
    - website/src/routes/_authenticated/payments/index.tsx
  modified:
    - website/src/components/layout/data/sidebar-data.ts
decisions:
  - "Payments sidebar link enabled at /payments by removing disabled flag and setting real URL"
  - "createPaymentColumns factory accepts action callbacks so PaymentTable owns dialog open state at page level"
  - "DateRangeFilter filters data client-side before passing to useReactTable, not as a server query param"
  - "RefundPaymentDialog uses hardcoded reason 'Admin refund' per plan spec; amount taken directly from payment.amount"
metrics:
  duration: "150s"
  completed_date: "2026-04-05"
  tasks_completed: 2
  files_changed: 11
---

# Phase 04 Plan 03: Payment Management Page Summary

Payment management page with DataTable, method/status faceted filters, date-range filter, and action dialogs for confirm/reject/refund/gateway payment initiation.

## What Was Built

- **PaymentsPage** at `/payments` — DataTable showing all payments with search, method faceted filter, status faceted filter, and a date-range popover filter
- **Payment action menu** — conditional dropdown: `pending` shows Konfirmasi + Tolak, `confirmed` shows Kembalikan Dana, `gateway` method always shows Buka Halaman Pembayaran
- **ConfirmPaymentDialog** — reuses `ConfirmDialog`, calls `useConfirmPayment().mutate(id)`, closes on success
- **RejectPaymentDialog** — custom AlertDialog with Textarea for required rejection reason, calls `useRejectPayment().mutate({ id, reason })`, submit disabled if reason empty
- **RefundPaymentDialog** — reuses `ConfirmDialog` with `destructive=true`, calls `useRefundPayment().mutate({ id, amount, reason })`
- **DateRangeFilter** — Popover with two `<Input type="date">` fields (Dari / Sampai), filters table client-side
- **Route** registered at `/_authenticated/payments/`
- **Sidebar** Payments link enabled at `/payments`

## Decisions Made

| Decision | Rationale |
|----------|-----------|
| `createPaymentColumns` factory pattern | Action callbacks passed from page level keep dialog state in PaymentsPage, not per-row |
| Client-side date filtering | API usePayments hook doesn't expose date params; filtered via useMemo before table data |
| Hardcoded gateway `xendit` | Plan spec: `initiateGateway({ id, gateway: 'xendit' })` — consistent with Phase 04 decisions |

## Deviations from Plan

None — plan executed exactly as written.

## Self-Check

- [x] All 9 feature files created
- [x] Route file created at `website/src/routes/_authenticated/payments/index.tsx`
- [x] Sidebar updated: `url: '/payments'`, `disabled` removed
- [x] TypeScript passes (`npx tsc --noEmit`)
- [x] All acceptance criteria verified

## Self-Check: PASSED
