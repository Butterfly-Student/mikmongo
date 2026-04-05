# Phase 5: Sales & Agents - Context

**Gathered:** 2026-04-05
**Status:** Ready for planning

<domain>
## Phase Boundary

Admin can manage sales agents (create, update, configure profile pricing), manage agent invoices (view, generate, pay, cancel, process), track hotspot sales (view list, create entry, generate Mikhmon vouchers), and manage the Mikhmon integration (profiles, setup script, sales reports, expire monitoring). Agent portal gains a sales history view and invoice management.

</domain>

<decisions>
## Implementation Decisions

### Claude's Discretion

All implementation decisions for Phase 5 are at Claude's discretion. The following established patterns from prior phases apply:

**Locked patterns (from prior phases):**
- TanStack Table + DataTableToolbar + faceted filters for all admin list views
- Side sheet / drawer for detail views (row click → right-side panel)
- Confirmation dialogs for destructive actions; rejection/cancel dialogs include a reason text field
- Feature-based directory structure: `website/src/features/{domain}/`
- Indonesian text (Bahasa Indonesia) for all UI labels, toasts, and messages
- Toast messages for all mutation success/failure
- TanStack Query useQuery + useMutation with queryKey invalidation

**Discretion items (agent decides approach):**
- Agent profile pricing UX — how prices per bandwidth profile are managed within the agent workflow
- Agent invoice generation trigger location and "Process batch" UX
- Hotspot sales — whether to use global or router-scoped list as the primary view
- Mikhmon page organization — tabs vs. separate nav items for vouchers/profiles/reports/expire
- Mikhmon router context selection — how active router is exposed in Mikhmon pages
- Agent portal sales history layout and prominence
- Exact column sets for all list views
- Loading skeletons and empty states
- Form layouts for create/edit dialogs

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### API Contract
- `docs/openapi.docs.yml` — Authoritative API contract. Key sections for Phase 5:
  - Sales Agents: `GET /api/v1/sales-agents`, `POST /api/v1/sales-agents`, `PUT /api/v1/sales-agents/{id}`, `GET /api/v1/sales-agents/{id}/profile-prices`, `PUT /api/v1/sales-agents/{id}/profile-prices/{profile}`
  - Agent Invoices: `GET /api/v1/sales-agents/{id}/invoices`, `POST /api/v1/sales-agents/{id}/invoices/generate`, `GET /api/v1/agent-invoices`, `GET /api/v1/agent-invoices/{id}`, `POST /api/v1/agent-invoices/{id}/pay`, `POST /api/v1/agent-invoices/{id}/cancel`, `POST /api/v1/agent-invoices/process`
  - Hotspot Sales: `GET /api/v1/hotspot-sales`, `POST /api/v1/hotspot-sales`, `GET /api/v1/routers/{router_id}/hotspot-sales`
  - Mikhmon: `POST /api/v1/routers/{router_id}/mikhmon/vouchers/generate`, `GET /api/v1/routers/{router_id}/mikhmon/vouchers`, `GET/POST/PUT/DELETE /api/v1/routers/{router_id}/mikhmon/profiles`, `POST /api/v1/routers/{router_id}/mikhmon/profiles/generate-script`, `GET /api/v1/routers/{router_id}/mikhmon/reports`, `GET /api/v1/routers/{router_id}/mikhmon/reports/summary`, `POST /api/v1/routers/{router_id}/mikhmon/expire/setup`, `GET /api/v1/routers/{router_id}/mikhmon/expire/status`
  - Agent Portal: `GET /agent-portal/v1/invoices`, `GET /agent-portal/v1/invoices/{id}`, `POST /agent-portal/v1/invoices/{id}/request-payment`, `GET /agent-portal/v1/sales`
  - Key schemas: `SalesAgentResponse`, `SalesProfilePriceResponse`, `AgentInvoiceResponse`, `HotspotSaleResponse`, `GenerateAgentInvoiceRequest`

### Prior Phase Context
- `.planning/phases/04-billing-payments/04-CONTEXT.md` — Invoice/payment patterns, gateway URL handling, faceted filter + date range, inline approve vs. dialog reject
- `.planning/phases/03-customers-routers-subscriptions/03-CONTEXT.md` — TanStack Table pattern, destructive-action confirmation dialogs, feature-based structure
- `.planning/phases/02-layout-dashboard-users/02-CONTEXT.md` — Sidebar data structure, router selector context, Indonesian currency formatting

### Requirements
- `.planning/REQUIREMENTS.md` — Phase 5 requirements: AGNT-01 to AGNT-08, HS-01 to HS-06

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- **TanStack Table + DataTableToolbar + FacetedFilter** (`website/src/components/data-table/`): All admin list views
- **Sheet/Drawer** (`website/src/components/ui/sheet.tsx`): Detail drawers for agents and agent invoices
- **ConfirmDialog** (`website/src/components/confirm-dialog.tsx`): Destructive actions
- **Axios adminClient** (`website/src/lib/axios/admin-client.ts`): All admin API calls
- **Axios agentClient** (`website/src/lib/axios/agent-client.ts`): Agent portal API calls
- **Agent Portal existing**: `website/src/features/agent-portal/invoices.tsx` — Phase 4 built this; Phase 5 adds sales alongside it
- **Sidebar Sales section**: Already has "Agents" and "Hotspot Sales" items, both `disabled: true` — Phase 5 enables them and adds Mikhmon sub-items

### Established Patterns
- Feature structure: `features/{domain}/data/schema.ts`, `components/`, `index.tsx`
- Query hooks: TanStack Query with queryKey arrays, invalidation after mutations
- Route files: `website/src/routes/_authenticated/{domain}/index.tsx`
- Currency: "Rp" prefix with Indonesian number format

### Integration Points
- Sidebar: Enable "Agents" (`/agents`) and "Hotspot Sales" (`/hotspot-sales`); add Mikhmon nav items
- Agent portal route tree: Add `/agent/sales` route alongside existing `/agent/invoices`
- Router context: Mikhmon endpoints are `router_id`-scoped — use active router from Zustand store (established in Phase 2)

</code_context>

<specifics>
## Specific Ideas

None — all decisions delegated to Claude's discretion.

</specifics>

<deferred>
## Deferred Ideas

None.

</deferred>

---
*Phase: 05-sales-agents*
*Context gathered: 2026-04-05*
