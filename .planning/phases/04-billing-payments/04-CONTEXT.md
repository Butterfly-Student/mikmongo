# Phase 4: Billing & Payments - Context

**Gathered:** 2026-04-04
**Status:** Ready for planning

<domain>
## Phase Boundary

Admin can generate and view invoices (with overdue filter), manage payments (view, manual confirm/reject/refund, and gateway-initiated via Midtrans/Xendit), and manage cash entries with an approval workflow and petty cash fund. Customer portal shows invoices, invoice details, payment history, and allows gateway-initiated payments. Agent portal shows invoice list with a payment request option.

</domain>

<decisions>
## Implementation Decisions

### Invoice List & Detail (Admin)
- **D-01:** Invoice detail opens in a **side sheet / drawer** — clicking a table row slides open a right-side panel with full invoice details and linked payments. The invoice list table remains visible behind it.
- **D-02:** Overdue filter uses a **faceted filter dropdown** (same pattern as subscription status filter on Phase 3) — not a separate tab or page.

### Payment Actions UX (Admin)
- **D-03:** All payment actions (confirm, reject, refund) use **confirmation dialogs** — confirm shows "Confirm this payment?", reject dialog includes a reason text field, refund dialog shows amount + confirmation. Consistent with the destructive-action pattern from Phase 3 (D-03).
- **D-04:** Payments list uses **faceted filter by payment method and status** (same pattern as invoices) **plus a date range filter**.

### Gateway Payment Initiation
- **D-05:** When admin initiates a gateway payment (`POST /payments/{id}/initiate-gateway`), the returned `payment_url` opens in a **new browser tab** directly.
- **D-06:** When a customer initiates payment via gateway in the Customer Portal (`POST /portal/v1/payments/{id}/pay`), the returned `payment_url` also opens in a **new browser tab**.

### Cash Entry Approval Workflow
- **D-07:** **Inline approve** (single-click button on the row) — no confirmation dialog for approve since it's lower risk. **Reject opens a dialog** with a reason field (consistent with rejection patterns across the app).
- **D-08:** Cash entry list is a **single table with faceted filter** by status (pending/approved/rejected) and type (income/expense).
- **D-09:** Petty cash fund is displayed as a **section within the cash management page** showing current balance and a top-up button — not a separate page.

### Customer Portal Billing
- **D-10:** Customer invoice list uses the **TanStack Table pattern** — consistent with the subscription table already in the customer portal (Phase 3). No card layout.
- **D-11:** Invoice detail in the customer portal opens in a **side sheet** (same pattern as admin — D-01).

### Agent Portal Billing (Agent's Discretion)
- Agent portal invoice list uses a simple table with a "Request Payment" button per row (PAY-07). Implementation details at agent's discretion.

### Claude's Discretion
- Exact columns shown in invoice and payment tables (beyond the obvious: invoice number, customer, amount, status, due date, payment date)
- Loading skeletons and empty states for all list views
- Toast messages for action success/failure
- Exact form layout for cash entry creation dialog
- Specific shadcn components for date range filter (DatePicker vs. Popover with two inputs)

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### API Contract
- `docs/openapi.docs.yml` — Authoritative API contract. Key sections for Phase 4:
  - Invoices: `GET /api/v1/invoices`, `GET /api/v1/invoices/overdue`, `GET /api/v1/invoices/{id}`, `POST /api/v1/invoices/trigger-monthly`
  - Payments: `GET /api/v1/payments`, `POST /api/v1/payments/{id}/confirm`, `POST /api/v1/payments/{id}/reject`, `POST /api/v1/payments/{id}/refund`, `POST /api/v1/payments/{id}/initiate-gateway`
  - Cash: `GET /api/v1/cash-entries`, `POST /api/v1/cash-entries`, `POST /api/v1/cash-entries/{id}/approve`, `POST /api/v1/cash-entries/{id}/reject`
  - Petty Cash: `GET /api/v1/petty-cash`, `POST /api/v1/petty-cash`, `POST /api/v1/petty-cash/{id}/topup`
  - Customer Portal: `GET /portal/v1/invoices`, `GET /portal/v1/invoices/{id}`, `GET /portal/v1/payments`, `GET /portal/v1/payments/{id}`, `POST /portal/v1/payments/{id}/pay`
  - Key schemas: `InvoiceResponse`, `PaymentResponse`, `GatewayPaymentResponse`, `CashEntryResponse`, `PettyCashFundResponse`, `RejectPaymentRequest`, `RefundPaymentRequest`, `CreateCashEntryRequest`

### Prior Phase Context
- `.planning/phases/03-customers-routers-subscriptions/03-CONTEXT.md` — TanStack Table pattern, destructive-action confirmation dialogs, feature-based structure, Zustand store patterns
- `.planning/phases/02-layout-dashboard-users/02-CONTEXT.md` — Established patterns: feature-based structure, TanStack Query hooks, Zustand persist, Indonesian Rp currency formatting, faceted filter pattern

### Requirements
- `.planning/REQUIREMENTS.md` — Phase 4 requirements: INV-01 to INV-05, PAY-01 to PAY-07, CASH-01 to CASH-04

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- **TanStack Table + DataTablePagination + FacetedFilter** (`website/src/components/data-table/`): Already established pattern for all list views — use for invoices, payments, and cash entries
- **Sheet/Drawer component** (`website/src/components/ui/sheet.tsx`): Will be used for invoice detail drawers (admin + customer portal)
- **AlertDialog / ConfirmDialog** (`website/src/components/confirm-dialog.tsx`): Use for payment confirm/reject/refund dialogs and cash entry rejection
- **Axios admin client** (`website/src/lib/axios/admin-client.ts`): All admin billing API calls go through this
- **Feature-based structure**: New billing feature goes in `website/src/features/billing/` (invoices, payments, cash-entries sub-features)

### Established Patterns
- **Feature structure**: `features/{name}/data/schema.ts`, `data/{name}.ts`, `components/`, `index.tsx`
- **Query hooks**: TanStack Query useQuery + useMutation, queryKey invalidation after mutations
- **Currency formatting**: "Rp" prefix with Indonesian number format (e.g., "Rp 15.230.000")
- **Status display**: Badges with color-coded variants (active/inactive/pending) — apply same to invoice/payment status

### Integration Points
- Billing sidebar nav items need enabling (were marked "Coming soon" / disabled in Phase 2)
- Customer portal billing routes connect to the existing customer portal route tree
- Agent portal billing connects to existing agent portal route tree

</code_context>

<specifics>
## Specific Ideas

- Invoice side sheet should show linked payments within the drawer (so admin can see payment status without leaving)
- Petty cash section on the cash management page shows current balance prominently (card or summary row), with a top-up dialog
- Gateway payment initiation: after `window.open(paymentUrl, '_blank')`, show a toast "Payment page opened in new tab" so admin knows action succeeded

</specifics>

<deferred>
## Deferred Ideas

- None — discussion stayed within phase scope.

</deferred>

---
*Phase: 04-billing-payments*
*Context gathered: 2026-04-04*
