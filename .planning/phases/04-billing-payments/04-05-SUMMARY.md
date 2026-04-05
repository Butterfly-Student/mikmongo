---
phase: 04-billing-payments
plan: 05
subsystem: ui
tags: [react, tanstack-table, customer-portal, agent-portal, invoices, payments]

requires:
  - phase: 04-01
    provides: billing schemas, portal invoice/payment API hooks
  - phase: 04-02
    provides: InvoiceDetailSheet component, invoiceStatuses, paymentStatuses

provides:
  - Customer portal invoice list at /customer/invoices with status filter, detail sheet, and gateway pay button
  - Customer portal payment history at /customer/payments (read-only)
  - Agent portal invoice list at /agent/invoices with Ajukan Pembayaran action for unpaid invoices

affects: [customer-portal, agent-portal, phase-05]

tech-stack:
  added: []
  patterns:
    - Inline TanStack Table with status badge filter (button group, no dropdown)
    - Isolated mutation button component (PayButton, RequestPaymentButton) to avoid hook-per-row penalty
    - usePortalInitiatePayment handles window.open + toast internally (hook-level side effect)

key-files:
  created:
    - website/src/features/customer-portal/invoices.tsx
    - website/src/features/customer-portal/payments.tsx
    - website/src/features/agent-portal/invoices.tsx
    - website/src/routes/customer/invoices/index.tsx
    - website/src/routes/customer/payments/index.tsx
    - website/src/routes/agent/invoices/index.tsx
  modified: []

key-decisions:
  - "Isolated PayButton/RequestPaymentButton as sub-components to call useMutation once per row without violating hook rules"
  - "Status filter implemented as simple button group (not faceted filter dropdown) matching customer portal simplicity"

patterns-established:
  - "Agent portal feature directory created at src/features/agent-portal/"
  - "Portal pages export default function component; route files import and wire them"

requirements-completed: [INV-04, INV-05, PAY-05, PAY-06, PAY-07]

duration: 8min
completed: 2026-04-05
---

# Phase 04 Plan 05: Customer & Agent Portal Billing Views Summary

**Customer portal invoice list with gateway pay button and detail sheet, read-only payment history, and agent portal client invoice list with payment request — three portal pages across /customer/invoices, /customer/payments, /agent/invoices**

## Performance

- **Duration:** ~8 min
- **Started:** 2026-04-05T03:00:00Z
- **Completed:** 2026-04-05T03:08:00Z
- **Tasks:** 2
- **Files modified:** 6

## Accomplishments

- Customer invoice page with TanStack Table, status filter (All/Unpaid/Paid/Overdue), invoice detail side sheet, and "Bayar Sekarang" button visible only for unpaid/overdue invoices
- Customer payment history page: read-only table with payment method, amount, date, and status badge
- Agent invoice page with voucher count, period, total, status badge, and "Ajukan Pembayaran" button for unpaid invoices only

## Task Commits

1. **Task 1: Customer portal invoice and payment pages** - `9e5b7d8` (feat)
2. **Task 2: Agent portal invoice page with payment request** - `e6eb392` (feat)

## Files Created/Modified

- `website/src/features/customer-portal/invoices.tsx` - Invoice list page with filter, detail sheet, pay button
- `website/src/features/customer-portal/payments.tsx` - Read-only payment history page
- `website/src/features/agent-portal/invoices.tsx` - Agent client invoice list with payment request
- `website/src/routes/customer/invoices/index.tsx` - Route /customer/invoices
- `website/src/routes/customer/payments/index.tsx` - Route /customer/payments
- `website/src/routes/agent/invoices/index.tsx` - Route /agent/invoices

## Decisions Made

- Isolated `PayButton` and `RequestPaymentButton` as sub-components per row to call `useMutation` once per component instance — this avoids calling hooks in a loop while keeping the mutation accessible per-row.
- Status filter implemented as a button group rather than a faceted dropdown — customer portal is simpler than admin; full DataTableToolbar with combobox facets would be over-engineered here.

## Deviations from Plan

None — plan executed exactly as written.

## Issues Encountered

None.

## User Setup Required

None - no external service configuration required.

## Next Phase Readiness

- All three portal billing pages are wired to real hooks from `use-portal-billing.ts`
- Phase 04 billing-payments is now complete (all 5 plans done)
- Phase 05 can proceed with reports/monitoring work

---
*Phase: 04-billing-payments*
*Completed: 2026-04-05*
