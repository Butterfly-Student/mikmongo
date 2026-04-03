# Phase 2: Layout, Dashboard & Users - Context

**Gathered:** 2026-04-03
**Status:** Ready for planning

<domain>
## Phase Boundary

Admin sees a functional ISP dashboard with router selector in sidebar, real-time ping in header, overview widgets, and can manage admin users. Covers: sidebar restructure (router selector, ISP nav groups, MikMongo branding), header ping display (WebSocket), dashboard overview page (KPI widgets, router health cards, recent activity feed), and admin user CRUD (list, create, edit, deactivate/delete with data table).

</domain>

<decisions>
## Implementation Decisions

### Router Store & Persistence
- **D-01:** Active router state lives in a **separate** `router-store.ts` (not a slice in auth-store). Cleaner separation — router selection is a different domain than authentication.
- **D-02:** Selected router ID **syncs from API on page load** (`GET /api/v1/routers/selected`). localStorage is a cache/fallback, not the source of truth. This prevents stale state if router was deselected in another session.
- **D-03:** When no router is selected, **dashboard widgets load normally** — they use `/reports/summary` which doesn't depend on a specific router. Ping displays "-- ms" in muted-foreground. Router health section shows "No Routers Configured" empty state.
- **D-04:** When admin **switches routers**, dashboard refreshes **immediately** — router health cards update, ping disconnects old WebSocket and connects to new router, KPI widgets re-fetch.

### Template Cleanup
- **D-05:** Template demo pages (apps, chats, tasks, help-center) are **kept as files** but **removed from sidebar nav**. No file deletion — just hide from navigation.
- **D-06:** Sidebar nav shows **all ISP groups with disabled items** for unimplemented phases. Items appear greyed out with tooltip "Coming soon". Only Dashboard and Users are clickable in Phase 2.

### Claude's Discretion
- User edit functionality — OpenAPI has no PUT endpoint for user updates. Decide at implementation time: skip edit, use workaround, or defer.
- Activity feed data source — if no API endpoint exists, decide between combined-data approach or placeholder message.
- Profile dropdown in header — currently hardcoded to template values. Update to use auth store or defer.
- Exact router store structure (fields, actions, persist config)
- WebSocket reconnection strategy for ping display
- Loading skeleton designs for widgets and table rows

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### API Contract
- `docs/openapi.docs.yml` — Authoritative API contract. Key endpoints for Phase 2:
  - Users: `GET /api/v1/users` (list, limit/offset), `POST /api/v1/users` (create), `GET /api/v1/users/{id}` (detail), `DELETE /api/v1/users/{id}` (delete). Note: No PUT/PATCH for user updates.
  - Routers: `GET /api/v1/routers` (list), `GET /api/v1/routers/selected` (current selection), `POST /api/v1/routers/select/{id}` (set active)
  - Reports: `GET /api/v1/reports/summary` (dashboard KPI data)
  - Ping WebSocket: `GET /api/v1/routers/{router_id}/monitor/ws/ping?address=8.8.8.8`
  - Security scheme: BearerAuth at lines 91-106
  - Schemas: UserResponse, CreateUserRequest, RouterResponse, ReportSummary

### UI Design Contract
- `.planning/phases/02-layout-dashboard-users/02-UI-SPEC.md` — Comprehensive visual and interaction contract. Defines: sidebar router selector behavior, ping display states, KPI widget typography/color, navigation group structure, user management page layout, form fields, validation messages, role display names, toast messages, empty states, error states. All Phase 2 UI must follow this spec.

### Project Constraints
- `.planning/PROJECT.md` — Tech stack constraints (React 19, TypeScript, TanStack, Tailwind, shadcn/ui, Axios, Zustand, Zod). Immutable data patterns. API must match OpenAPI spec exactly.

### Requirements
- `.planning/REQUIREMENTS.md` — Phase 2 requirements: NAV-01, NAV-02, NAV-03, NAV-05, DASH-01, DASH-02, DASH-03, USER-01, USER-02, USER-03, USER-04. Each has specific acceptance criteria.

### Prior Phase Context
- `.planning/phases/01/01-CONTEXT.md` — Phase 1 decisions: Indonesian UI text, Zustand persist pattern, Axios client architecture, feature-based structure, auth store design.

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- **Sidebar primitives** (`website/src/components/ui/sidebar.tsx`): Full shadcn sidebar with SidebarProvider, SidebarHeader, SidebarContent, SidebarFooter, SidebarGroup, SidebarMenu, SidebarMenuButton, SidebarMenuBadge, etc. 729 lines, very feature-rich.
- **Sidebar components** (`website/src/components/layout/`): `app-sidebar.tsx` (main sidebar), `team-switcher.tsx` (to be replaced by RouterSelector), `nav-group.tsx` (nav group renderer with 3 item types), `nav-user.tsx` (sidebar footer user widget), `app-title.tsx` (MikMongo branding), `sidebar-data.ts` (nav data to restructure), `types.ts` (SidebarData, NavGroup, NavItem types).
- **Header** (`website/src/components/layout/header.tsx`): Sticky h-16 header with SidebarTrigger, children slot. Ping display goes in the right side.
- **Dashboard** (`website/src/features/dashboard/`): Existing page with hardcoded data, Recharts bar chart, tabs (Overview/Analytics/Reports). Must be rebuilt per UI spec.
- **Users feature** (`website/src/features/users/`): Full CRUD with TanStack Table, react-hook-form + zod, AlertDialog confirmations, DataTablePagination, faceted filters, bulk actions. Template uses faker data — must be rewired to API.
- **Shared components**: `confirm-dialog.tsx`, `select-dropdown.tsx`, `password-input.tsx`, `data-table/` (DataTablePagination, DataTableColumnHeader, DataTableToolbar, FacetedFilter).
- **Auth store** (`website/src/stores/auth-store.ts`): Zustand with persist, partialize, hydration gate pattern. Router store should follow same patterns.
- **Axios clients** (`website/src/lib/axios/`): `admin-client.ts` with Bearer token + silent refresh. All Phase 2 API calls use this.
- **API functions** (`website/src/api/auth.ts`): Pattern for Zod-validated API calls. New API functions should follow same pattern.
- **Zod schemas** (`website/src/lib/schemas/auth.ts`): AdminUserSchema already defined. New schemas needed for Router, CreateUser, ReportSummary.

### Established Patterns
- **Feature-based structure**: `features/{name}/data/schema.ts`, `data/{name}.ts`, `components/`, `index.tsx`
- **Data fetching**: Currently only useMutation (no useQuery hooks). Phase 2 introduces first useQuery hooks for lists.
- **State management**: Zustand with persist + partialize + hydration gate
- **Form handling**: react-hook-form + zod validation + FormMessage for field errors + toast for API errors
- **Table pattern**: TanStack Table with URL-synced state, faceted filters, column headers, pagination
- **Immutability**: New objects via set(), never mutate existing state
- **Error handling**: Sonner toasts for user-facing errors, per-field validation for forms

### Integration Points
- Router selector must integrate with Zustand router store and call `POST /api/v1/routers/select/{id}`
- Ping display needs WebSocket connection managed by a custom hook, reading active router from store
- Dashboard widgets fetch from `/reports/summary` via TanStack Query useQuery
- Router health cards fetch from `/api/v1/routers` — same data as sidebar selector
- Users table fetches from `/api/v1/users` — replacing faker data
- Profile dropdown should eventually use auth store (currently hardcoded)

### Known Issues
- No PUT/PATCH endpoint for user updates in OpenAPI — edit user may need workaround or deferral
- No activity feed API endpoint — UI spec allows placeholder or combined-data approach
- `website/src/api/types.ts` has limited types — needs expansion for Router, ReportSummary schemas
- Template user fields (firstName/lastName/username) don't match OpenAPI (full_name) — must realign

</code_context>

<specifics>
## Specific Ideas

- Router selector replaces TeamSwitcher in SidebarHeader — uses shadcn Select (not DropdownMenu) for better keyboard navigation
- Ping display is small and non-intrusive: `text-sm text-muted-foreground` with Activity icon — not a badge, not a card
- Dashboard removes the Overview/Analytics/Reports tabs — single page with all widgets visible
- Sidebar shows all ISP navigation groups even if most items are disabled/"Coming soon"
- Currency formatting uses "Rp" prefix with Indonesian number formatting (e.g., "Rp 15.230.000")
- Role display names mapped from API values: superadmin→"Super Admin", cs→"Customer Service", readonly→"Read Only", etc.

</specifics>

<deferred>
## Deferred Ideas

- User edit functionality — deferred pending API support (no PUT endpoint exists)
- Template demo page files (apps, chats, tasks, help-center) — kept but hidden, can delete in a future cleanup phase
- Profile dropdown update to use auth store — cosmetic, not blocking

</deferred>

---
*Phase: 02-layout-dashboard-users*
*Context gathered: 2026-04-03*
