---
phase: 04-billing-payments
plan: 01
subsystem: api
tags: [zod, tanstack-query, axios, billing, payments, invoices, typescript]

requires:
  - phase: 03-customers-routers-subscriptions
    provides: "admin-client, customer-client, agent-client, ApiResponseSchema, MetaSchema patterns"

provides:
  - "Zod schemas for all billing entities: Invoice, Payment, CashEntry, PettyCashFund, GatewayPayment, AgentInvoice"
  - "Admin API functions for invoices, payments, cash entries, petty cash"
  - "Portal API functions for customer invoices/payments and agent invoices"
  - "TanStack Query hooks for all billing operations with toast feedback"

affects: [04-02-PLAN.md, 04-03-PLAN.md, 04-04-PLAN.md, 04-05-PLAN.md]

tech-stack:
  added: []
  patterns:
    - "Billing schemas follow ApiResponseSchema wrapper + MetaSchema for paginated lists"
    - "Portal APIs use separate axios clients: customerClient, agentClient, adminClient"
    - "useMutation hooks use unknown error type with type assertion narrowing"

key-files:
  created:
    - website/src/lib/schemas/billing.ts
    - website/src/api/invoice.ts
    - website/src/api/payment.ts
    - website/src/api/cash.ts
    - website/src/api/portal/invoice.ts
    - website/src/api/portal/payment.ts
    - website/src/api/portal/agent-invoice.ts
    - website/src/hooks/use-invoices.ts
    - website/src/hooks/use-payments.ts
    - website/src/hooks/use-cash.ts
    - website/src/hooks/use-portal-billing.ts
  modified: []

key-decisions:
  - "CashEntryListResponseSchema uses optional meta since cash-entries endpoint may not return pagination meta"
  - "listOverdueInvoices uses dedicated ApiResponseSchema(z.array(invoiceResponseSchema)) since endpoint returns array not paginated list"
  - "useInitiateGatewayPayment and usePortalInitiatePayment both call window.open(data.payment_url, '_blank') on success before showing toast"

patterns-established:
  - "Billing list schemas: ApiResponseSchema(z.array(schema)).extend({ meta: MetaSchema }) for paginated"
  - "Portal APIs co-locate schema definitions inline when only used by that module"

requirements-completed: [INV-01, INV-02, INV-03, PAY-01, PAY-02, PAY-03, PAY-04, CASH-01, CASH-02, CASH-03, CASH-04, INV-04, INV-05, PAY-05, PAY-06, PAY-07]

duration: 4min
completed: 2026-04-05
---

# Phase 04 Plan 01: Billing Data Layer Summary

**Zod schemas and TanStack Query hooks for all billing entities — invoices, payments, cash entries, petty cash, and agent invoices — across admin and portal APIs**

## Performance

- **Duration:** 4 min
- **Started:** 2026-04-05T02:30:25Z
- **Completed:** 2026-04-05T02:34:00Z
- **Tasks:** 2
- **Files modified:** 11

## Accomplishments

- Created `billing.ts` with 9 Zod schemas covering all OpenAPI billing response types and 11 list/detail response wrapper schemas
- Created 7 API files covering 20+ endpoints across admin (invoices, payments, cash, petty cash) and portals (customer, agent)
- Created 4 hook files with 20 TanStack Query hooks (reads + mutations) with Indonesian toast messages and unknown error narrowing

## Task Commits

1. **Task 1: Create billing Zod schemas and API functions** - `a43344b` (feat)
2. **Task 2: Create TanStack Query hooks for all billing operations** - `690eb01` (feat)

## Files Created/Modified

- `website/src/lib/schemas/billing.ts` - All billing Zod schemas: invoiceResponseSchema, paymentResponseSchema, cashEntryResponseSchema, pettyCashFundResponseSchema, agentInvoiceResponseSchema, gatewayPaymentResponseSchema, form schemas, list/detail response schemas
- `website/src/api/invoice.ts` - Admin invoice API: listInvoices, listOverdueInvoices, getInvoice, triggerMonthlyBilling
- `website/src/api/payment.ts` - Admin payment API: listPayments, confirmPayment, rejectPayment, refundPayment, initiateGatewayPayment
- `website/src/api/cash.ts` - Admin cash API: listCashEntries, createCashEntry, approveCashEntry, rejectCashEntry, listPettyCashFunds, createPettyCashFund, topUpPettyCashFund
- `website/src/api/portal/invoice.ts` - Customer portal invoice API using customerClient
- `website/src/api/portal/payment.ts` - Customer portal payment API using customerClient
- `website/src/api/portal/agent-invoice.ts` - Agent portal invoice API using agentClient
- `website/src/hooks/use-invoices.ts` - useInvoices, useOverdueInvoices, useInvoice, useTriggerMonthlyBilling
- `website/src/hooks/use-payments.ts` - usePayments, useConfirmPayment, useRejectPayment, useRefundPayment, useInitiateGatewayPayment
- `website/src/hooks/use-cash.ts` - useCashEntries, useCreateCashEntry, useApproveCashEntry, useRejectCashEntry, usePettyCashFunds, useCreatePettyCashFund, useTopUpPettyCashFund
- `website/src/hooks/use-portal-billing.ts` - usePortalInvoices, usePortalInvoice, usePortalPayments, usePortalInitiatePayment, useAgentPortalInvoices, useAgentRequestPayment

## Decisions Made

- `CashEntryListResponseSchema` uses `meta: MetaSchema.optional()` since cash-entries endpoint may not return pagination
- `listOverdueInvoices` uses `ApiResponseSchema(z.array(invoiceResponseSchema))` since overdue endpoint returns array not paginated
- Gateway payment hooks open `payment_url` in new tab via `window.open` before showing toast

## Deviations from Plan

None - plan executed exactly as written.

## Issues Encountered

None.

## User Setup Required

None - no external service configuration required.

## Next Phase Readiness

- All billing data layer ready for Phase 04-02 through 04-05 UI implementation
- Schemas, API functions, and hooks importable via `@/lib/schemas/billing`, `@/api/invoice`, `@/api/payment`, `@/api/cash`, `@/hooks/use-invoices`, etc.

---
*Phase: 04-billing-payments*
*Completed: 2026-04-05*
