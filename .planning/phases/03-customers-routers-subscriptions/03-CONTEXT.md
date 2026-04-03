# Phase 3: Customers, Routers & Subscriptions - Context

**Gathered:** 2026-04-03
**Status:** Ready for planning

<domain>
## Phase Boundary

Admin can manage the full customer lifecycle (create, activate, registration pipeline), manage MikroTik routers (sync, test connection), construct bandwidth profiles (plans), and manage subscriptions linking profiles to customers on specific routers. Follows existing CRUD table patterns.

</domain>

<decisions>
## Implementation Decisions

### Registration Pipeline UX
- **D-01:** Built using a data table with standard actions, leveraging the existing TanStack Table and pagination patterns instead of a Kanban board.

### Subscription Assignment Flow
- **D-02:** Assignments happen on a global "Subscriptions" page, utilizing select dropdowns to select the Customer and Profile, rather than deep linking strictly from the Customer profile page.

### Destructive Subscription Actions
- **D-03:** Actions such as Suspend, Isolate, or Terminate will use explicit Confirmation Dialogs (similar to Delete User) to prevent accidental billing and service interruptions.

### Router Sync Behavior
- **D-04:** [Unspecified - the agent's Discretion] The interface should accommodate a balance of responsiveness and safety for sync actions. Default to blocking loading states for connection tests and non-interruptable syncs if not stated otherwise.

### Claude's Discretion
- Background loading states and table empty states.
- Re-use of Zod schemas and query invalidation triggers.
- Specific form layouts for Customer CRUD.
</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### API Contract
- `docs/openapi.docs.yml` — Authoritative API contract for Router, Customer, Subscription, and Bandwidth Profile endpoints.

### UI Design Contract
- Prior phase specs (`.planning/phases/02-layout-dashboard-users/02-UI-SPEC.md` and `01-UI-SPEC.md`) establish generic Table, Select, Profile, and Toast conventions to follow tightly.

### Requirements
- `.planning/REQUIREMENTS.md` — Phase 3 requirements: CUST-01 to CUST-07, RTR-01 to RTR-06, BW-01 to BW-03, SUB-01 to SUB-05.
</canonical_refs>
