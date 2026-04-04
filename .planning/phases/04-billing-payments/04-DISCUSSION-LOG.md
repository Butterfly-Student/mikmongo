# Phase 4: Billing & Payments - Discussion Log

> **Audit trail only.** Do not use as input to planning, research, or execution agents.
> Decisions are captured in CONTEXT.md — this log preserves the alternatives considered.

**Date:** 2026-04-04
**Phase:** 04-billing-payments
**Areas discussed:** Invoice List & Detail, Payment Confirmation UX, Gateway Payment Initiation, Cash Entry Approval Workflow, Customer Portal Billing UI

---

## Invoice List & Detail

| Option | Description | Selected |
|--------|-------------|----------|
| Separate detail page | Navigate to `/invoices/{id}` — full page | |
| Side sheet / drawer | Right-side panel slides open; list stays visible | ✓ |
| Agent's discretion | Agent picks best pattern | |

**Overdue filter:**

| Option | Description | Selected |
|--------|-------------|----------|
| Toggle tab/chip ("All \| Overdue") | Top of table tab switch | |
| Faceted filter dropdown | Same pattern as subscription status filter | ✓ |
| Agent's discretion | | |

**User's choice:** Side sheet for detail, faceted filter dropdown for overdue.
**Notes:** Consistent with fast review UX — list stays visible while inspecting invoice details.

---

## Payment Confirmation UX

| Option | Description | Selected |
|--------|-------------|----------|
| Confirmation dialog for all | Confirm, reject (with reason), refund all open dialogs | ✓ |
| Inline confirm, dialog for reject/refund | Inline green checkmark for confirm, dialogs for destructive | |
| Agent's discretion | | |

**Payments list filter:**

| Option | Description | Selected |
|--------|-------------|----------|
| Faceted filter by method & status + date range | Both filter types combined | ✓ |
| Simple date range only | | |
| Agent's discretion | | |

**User's choice:** All three actions (confirm/reject/refund) via confirmation dialogs. Payments list has faceted filter by method/status plus date range.
**Notes:** Deliberate safety — even confirm uses a dialog. Reject dialog includes reason field.

---

## Gateway Payment Initiation

**Admin side:**

| Option | Description | Selected |
|--------|-------------|----------|
| Open in new tab | `payment_url` opens in new browser tab | ✓ |
| Modal with copy link + Open button | Dialog shows URL to share with customer, Open button | |
| Agent's discretion | | |

**Customer Portal side:**

| Option | Description | Selected |
|--------|-------------|----------|
| Redirect in same tab | Full redirect to gateway | |
| Open in new tab | New browser tab | ✓ |
| Agent's discretion | | |

**User's choice:** New tab for both admin-initiated and customer-initiated gateway payments.

---

## Cash Entry Approval Workflow

**Approve/reject mechanism:**

| Option | Description | Selected |
|--------|-------------|----------|
| Confirmation dialog for both | Approve and reject both open dialogs | |
| Inline approve, dialog for reject | Approve is single-click; reject opens reason dialog | ✓ |
| Agent's discretion | | |

**List layout:**

| Option | Description | Selected |
|--------|-------------|----------|
| Single table with faceted filter (status + type) | Combined filter approach | ✓ |
| Tabbed view (Pending \| All Entries) | Separate tabs | |
| Agent's discretion | | |

**Petty cash fund:**

| Option | Description | Selected |
|--------|-------------|----------|
| Section on cash management page (balance + top-up) | Sub-section on same page | ✓ |
| Agent's discretion | | |

**User's choice:** Inline approve (fast), reject via dialog with reason. Single table with faceted status+type filter. Petty cash as a section on the same page.
**Notes:** Approve is low-risk enough for inline click; reject requires reason, so dialog is appropriate.

---

## Customer Portal Billing UI

**Invoice display:**

| Option | Description | Selected |
|--------|-------------|----------|
| Card-based layout | Invoice cards with status chip + "Pay Now" button | |
| Simple TanStack Table | Consistent with portal subscription table (Phase 3) | ✓ |
| Agent's discretion | | |

**Invoice detail in portal:**

| Option | Description | Selected |
|--------|-------------|----------|
| Separate detail page `/customer/invoices/{id}` | Full page navigation | |
| Side sheet (same as admin) | Right-side panel | ✓ |
| Agent's discretion | | |

**User's choice:** Simple table (not cards) for invoice list; side sheet for detail — keeps portal UX consistent.

---

## Claude's Discretion

- Exact table columns beyond core fields (invoice number, amount, status, due date)
- Loading skeletons and empty states for all list views
- Toast messages for action feedback
- Cash entry creation dialog layout
- Date range filter component choice (DatePicker vs. Popover with two inputs)
- Agent portal invoice table details (PAY-07)

## Deferred Ideas

None — discussion stayed within phase scope.
