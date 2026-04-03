# Discussion Log: Phase 3

**Date:** 2026-04-03
**Phase:** 3 (Customers, Routers & Subscriptions)

## Question 1: Registration Pipeline UX
**Options:**
- A: A standard data table (reusing our existing TanStack Table) with "Approve" / "Reject" buttons on each row.
- B: A visual Kanban/Pipeline board where customers can be moved from "Pending" to "Approved" or "Rejected".
**Selection:** 1.A

## Question 2: Subscription Assignment Flow
**Options:**
- A: Done from the individual Customer profile page (e.g., clicking "Add Subscription" while viewing a customer).
- B: Done from a global "Subscriptions" page using Select dropdowns to pick the customer and the profile.
**Selection:** 2.B

## Question 3: Destructive Subscription Actions
**Options:**
- A: Explicit Confirmation Dialogs for actions like Suspend or Terminate (similar to the Delete User dialog we just built).
- B: Immediate actions with just an undo toast notification.
**Selection:** 3.A
