# Phase 2: Admin Network Management - Research

**Researched:** 2026-03-30
**Domain:** React data tables, CRUD forms, server-side pagination, MikroTik operations
**Confidence:** HIGH

## Summary

Phase 2 builds CRUD management pages for Customers, Routers (with Bandwidth Profiles), and Subscriptions on top of the Phase 1 scaffold. The backend API is fully built -- all endpoints exist in `internal/router/admin.go` with consistent patterns: limit/offset pagination via query params, a `{ success, data, error, meta }` JSON envelope, and UUID-based resource IDs.

The frontend scaffold from Phase 1 provides: TanStack Router file-based routing under `_admin/`, TanStack Query v5 for data fetching, TanStack Form v1 for forms with Zod validation, TanStack Table v8 and TanStack Virtual v3 (both already installed but unused), Shadcn/UI components, and an `adminClient` Axios instance at `/api/v1` with JWT interceptors.

**Primary recommendation:** Build reusable data table and form dialog components in the first plan (02-01), then reuse them across router and subscription pages. Use TanStack Table v8 for column definitions + sorting/filtering state, TanStack Query for server-side data fetching with `keepPreviousData`, and TanStack Virtual v3 only where the success criteria demands it (500+ row scroll). Forms use TanStack Form v1 + Zod (matching the login page pattern from Phase 1).

<phase_requirements>
## Phase Requirements

| ID | Description | Research Support |
|----|-------------|------------------|
| CUST-01 | Tabel pelanggan dengan pagination server-side, search, filter status | API: `GET /api/v1/customers?limit=N&offset=N`. Note: backend currently only supports limit/offset, no search/filter params. Frontend must do client-side search on loaded page data, OR implement search as a local filter on fetched data. |
| CUST-02 | Form create pelanggan (nama, email, telepon, alamat, sales agent) | API: `POST /api/v1/customers` with `CreateCustomerRequest` DTO. Supports optional subscription creation. |
| CUST-03 | Form edit pelanggan | API: `PUT /api/v1/customers/:id` with `UpdateCustomerRequest` (partial update, pointer fields). |
| CUST-04 | Delete pelanggan dengan konfirmasi dialog | API: `DELETE /api/v1/customers/:id`. Use AlertDialog confirmation. |
| CUST-05 | Tombol activate/deactivate account pelanggan | API: `POST /api/v1/customers/:id/activate-account` and `/deactivate-account`. |
| CUST-06 | Detail pelanggan: info lengkap + riwayat langganan + invoice | API: `GET /api/v1/customers/:id`. Subscription/invoice history requires cross-referencing subscription and invoice endpoints. |
| ROUT-01 | Tabel router dengan status online/offline | API: `GET /api/v1/routers?limit=N&offset=N`. Response includes `status` and `is_active` fields. |
| ROUT-02 | Form create router (nama, host, port, username, password, tipe) | API: `POST /api/v1/routers` with `CreateRouterRequest` DTO. |
| ROUT-03 | Form edit router | API: `PUT /api/v1/routers/:router_id` with `UpdateRouterRequest`. |
| ROUT-04 | Delete router dengan konfirmasi | API: `DELETE /api/v1/routers/:router_id`. |
| ROUT-05 | Tombol "Select Router" untuk set router aktif | API: `POST /api/v1/routers/select/:id`. Response returns selected router. |
| ROUT-06 | Tombol "Test Connection" dengan feedback status | API: `POST /api/v1/routers/:router_id/test-connection`. Returns `{ message: "connection successful" }` or error. |
| ROUT-07 | Tombol "Sync" per router dan "Sync All" | API: `POST /api/v1/routers/:router_id/sync` and `POST /api/v1/routers/sync-all`. |
| ROUT-08 | Badge router yang sedang terpilih (active router) | API: `GET /api/v1/routers/selected` returns currently selected router or null. |
| BWP-01 | Tabel bandwidth profiles per router | API: `GET /api/v1/routers/:router_id/bandwidth-profiles?limit=N&offset=N`. Scoped to router. |
| BWP-02 | Form create/edit bandwidth profile | API: `POST` and `PUT /api/v1/routers/:router_id/bandwidth-profiles(/:id)`. Many fields including MikroTik passthrough. |
| BWP-03 | Delete bandwidth profile dengan konfirmasi | API: `DELETE /api/v1/routers/:router_id/bandwidth-profiles/:id`. |
| SUB-01 | Tabel subscriptions dengan filter router, status, customer | API: `GET /api/v1/routers/:router_id/subscriptions?limit=N&offset=N`. Scoped to router. Status filtering must be client-side. |
| SUB-02 | Form create subscription | API: `POST /api/v1/routers/:router_id/subscriptions` with `CreateSubscriptionRequest`. Requires customer_id, plan_id (bandwidth profile), username. |
| SUB-03 | Form edit subscription | API: `PUT /api/v1/routers/:router_id/subscriptions/:id`. |
| SUB-04 | Aksi per subscription: activate, isolate, restore, suspend, terminate | APIs: `POST .../subscriptions/:id/{activate,isolate,restore,suspend,terminate}`. Isolate and suspend accept optional `{ reason }` body. |
| SUB-05 | Status badge: active, isolated, suspended, terminated | Response `status` field. Values: "pending", "active", "isolated", "suspended", "terminated". Use Shadcn Badge with color variants. |
| SUB-06 | Detail subscription dengan riwayat perubahan status | API: `GET .../subscriptions/:id`. Includes `mikrotik` live data. Status history not available from single endpoint -- show current status + key dates (activated_at, expiry_date). |
</phase_requirements>

## Standard Stack

### Core (already installed in Phase 1)
| Library | Version | Purpose | Why Standard |
|---------|---------|---------|--------------|
| @tanstack/react-table | ^8.21.3 | Column definitions, sorting, filtering, pagination state | De facto standard for headless tables in React |
| @tanstack/react-virtual | ^3.13.23 | Virtualized row rendering for 500+ row tables | Required by success criteria #5 |
| @tanstack/react-query | ^5.95.2 | Server state management, caching, refetching | Already established in Phase 1 |
| @tanstack/react-form | ^1.28.5 | Form state management | Already established in Phase 1 login pages |
| @tanstack/react-router | ^1.168.8 | File-based routing | Already established in Phase 1 |
| zod | ^3.25.76 | Schema validation for forms and API responses | Already established in Phase 1 |
| shadcn/ui | ^4.1.1 | Component library (need additional: table, dialog, select, alert-dialog, tabs, popover, command) | Already established in Phase 1 |
| axios | ^1.14.0 | HTTP client via adminClient | Already established in Phase 1 |

### Shadcn Components to Add
| Component | Purpose |
|-----------|---------|
| table | Data table UI (thead/tbody/tr/td with proper styling) |
| dialog | Create/edit form modals |
| alert-dialog | Delete confirmation dialogs |
| select | Dropdown selectors (status filter, router select, etc.) |
| tabs | Tab navigation on detail pages |
| popover | Filter popover menus |
| command | Searchable select (customer picker in subscription form) |
| pagination | Page navigation controls |
| textarea | Notes/address fields in forms |
| switch | Toggle fields (is_active, auto_isolate) |

**Installation:**
```bash
cd dashboard && npx shadcn@latest add table dialog alert-dialog select tabs popover command pagination textarea switch
```

### Alternatives Considered
| Instead of | Could Use | Tradeoff |
|------------|-----------|----------|
| TanStack Table | AG Grid | Overkill -- AG Grid is heavier, TanStack Table already installed |
| TanStack Form | React Hook Form | TanStack Form already established in Phase 1 login pattern |
| TanStack Virtual rows | Full pagination only | Success criteria #5 requires smooth 500+ row scroll |

## Architecture Patterns

### Recommended Project Structure
```
dashboard/src/
├── api/
│   ├── auth.ts              # (existing)
│   ├── reports.ts           # (existing)
│   ├── customers.ts         # NEW: customer CRUD API functions
│   ├── routers.ts           # NEW: router CRUD + actions API functions
│   ├── bandwidth-profiles.ts # NEW: bandwidth profile CRUD
│   └── subscriptions.ts     # NEW: subscription CRUD + lifecycle actions
├── lib/
│   ├── schemas/
│   │   ├── auth.ts          # (existing)
│   │   ├── customer.ts      # NEW: Zod schemas for customer API
│   │   ├── router.ts        # NEW: Zod schemas for router API
│   │   ├── bandwidth-profile.ts # NEW
│   │   └── subscription.ts  # NEW
│   └── axios/
│       └── admin-client.ts  # (existing, reused)
├── components/
│   ├── ui/                  # (existing Shadcn components)
│   ├── shared/
│   │   ├── DataTable.tsx    # NEW: reusable table with TanStack Table + Virtual
│   │   ├── DataTablePagination.tsx  # NEW: pagination controls
│   │   ├── ConfirmDialog.tsx # NEW: reusable delete confirmation
│   │   └── FormDialog.tsx   # NEW: reusable form modal wrapper
│   └── layout/admin/        # (existing)
├── hooks/
│   ├── useCustomers.ts      # NEW: TanStack Query hooks for customer data
│   ├── useRouters.ts        # NEW: TanStack Query hooks for router data
│   ├── useBandwidthProfiles.ts # NEW
│   └── useSubscriptions.ts  # NEW
└── routes/_admin/
    ├── index.tsx             # (existing overview page)
    ├── customers/
    │   ├── index.tsx         # Customer list page
    │   └── $customerId.tsx   # Customer detail page
    ├── routers/
    │   ├── index.tsx         # Router list page
    │   └── $routerId/
    │       ├── index.tsx     # Router detail + bandwidth profiles
    │       └── bandwidth-profiles.tsx  # (or nested in router detail)
    └── subscriptions/
        └── index.tsx         # Subscription list (requires selected router context)
```

### Pattern 1: API Function with Zod Validation
**What:** Each API module exports typed functions that call `adminClient` and validate responses with Zod
**When to use:** Every API call in this phase
**Example:**
```typescript
// src/api/customers.ts
import { z } from "zod"
import { adminClient } from "@/lib/axios/admin-client"
import { ApiResponseSchema } from "@/lib/schemas/auth"

const CustomerSchema = z.object({
  id: z.string().uuid(),
  customer_code: z.string(),
  full_name: z.string(),
  email: z.string().email().nullable().optional(),
  phone: z.string(),
  address: z.string().nullable().optional(),
  is_active: z.boolean(),
  created_at: z.string(),
  updated_at: z.string(),
})

export type Customer = z.infer<typeof CustomerSchema>

const CustomerListSchema = ApiResponseSchema(z.array(CustomerSchema))

export async function fetchCustomers(params: { limit: number; offset: number }) {
  const { data } = await adminClient.get("/customers", { params })
  const parsed = CustomerListSchema.parse(data)
  return { data: parsed.data, meta: parsed.meta! }
}
```

### Pattern 2: TanStack Query Hook with Pagination
**What:** Custom hook wrapping `useQuery` with pagination state
**When to use:** Every list page
**Example:**
```typescript
// src/hooks/useCustomers.ts
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query"
import { fetchCustomers, createCustomer, deleteCustomer } from "@/api/customers"

export function useCustomerList(limit: number, offset: number) {
  return useQuery({
    queryKey: ["customers", { limit, offset }],
    queryFn: () => fetchCustomers({ limit, offset }),
    placeholderData: keepPreviousData, // smooth page transitions
    staleTime: 30_000,
  })
}

export function useCreateCustomer() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: createCustomer,
    onSuccess: () => qc.invalidateQueries({ queryKey: ["customers"] }),
  })
}
```

### Pattern 3: Reusable DataTable Component
**What:** Generic DataTable component wrapping TanStack Table + TanStack Virtual for optional virtualization
**When to use:** All list pages (customers, routers, bandwidth profiles, subscriptions)
**Key features:**
- Accept `columns` (ColumnDef[]) and `data` as props
- Server-side pagination via `pageCount` + `onPaginationChange`
- Optional virtualizer mode for large datasets (SUB-01 with 500+ rows)
- Column header sorting (client-side within current page)
- Search input that filters current page data

### Pattern 4: Form Dialog with TanStack Form + Zod
**What:** Modal dialog containing a form for create/edit operations
**When to use:** CUST-02, CUST-03, ROUT-02, ROUT-03, BWP-02, SUB-02, SUB-03
**Example:**
```typescript
// Reusable pattern from Phase 1 login.tsx
const form = useForm({
  defaultValues: initialValues,
  onSubmit: async ({ value }) => {
    const result = schema.safeParse(value)
    if (!result.success) return
    await mutation.mutateAsync(result.data)
  },
})
```

### Anti-Patterns to Avoid
- **Fetching all data client-side:** Backend supports limit/offset pagination -- always use it. Never fetch all 500+ customers at once.
- **Hardcoding router_id:** Subscriptions and bandwidth profiles are scoped to `router_id` in the URL. Always read the selected router from state or URL params.
- **Skipping Zod validation on API responses:** Phase 1 established the pattern of validating every API response. Do not bypass this.
- **Inline column definitions:** Define `columns` as a const outside the component to prevent re-renders.
- **Mutating query cache directly:** Use `queryClient.invalidateQueries()` after mutations, not manual cache updates.

## Don't Hand-Roll

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| Data table rendering | Custom table with manual DOM | TanStack Table v8 + Shadcn Table | Sorting, filtering, pagination state management is complex |
| Virtual scrolling | Custom IntersectionObserver rows | TanStack Virtual v3 | Window calculations, overscan, dynamic height |
| Form state | useState per field | TanStack Form v1 | Validation, touched state, submit handling |
| Confirmation dialog | Custom modal + state | Shadcn AlertDialog | Accessible, handles focus trap |
| Searchable select | Custom dropdown + filter | Shadcn Command (cmdk) | Keyboard navigation, search, accessibility |
| Pagination controls | Custom prev/next buttons | Shadcn Pagination + TanStack Table state | Consistent UI, page size selector |

**Key insight:** All data table, form, and dialog primitives are already installed. The work is composing them into page-specific implementations, not building infrastructure.

## API Endpoint Reference

### Customers (`/api/v1/customers`)
| Method | Path | Query Params | Body | Response |
|--------|------|-------------|------|----------|
| GET | `/customers` | `limit`, `offset` | - | `{ success, data: Customer[], meta: { total, limit, offset } }` |
| POST | `/customers` | - | `CreateCustomerRequest` | `{ success, data: { customer, subscription? } }` |
| GET | `/customers/:id` | - | - | `{ success, data: Customer }` |
| PUT | `/customers/:id` | - | `UpdateCustomerRequest` (partial) | `{ success, data: Customer }` |
| DELETE | `/customers/:id` | - | - | `{ success, data: { message } }` |
| POST | `/customers/:id/activate-account` | - | - | `{ success, data: { message } }` |
| POST | `/customers/:id/deactivate-account` | - | - | `{ success, data: { message } }` |

### Routers (`/api/v1/routers`)
| Method | Path | Query Params | Body | Response |
|--------|------|-------------|------|----------|
| GET | `/routers` | `limit`, `offset` | - | `{ success, data: Router[] }` (no meta -- see note) |
| POST | `/routers` | - | `CreateRouterRequest` | `{ success, data: Router }` |
| GET | `/routers/selected` | - | - | `{ success, data: Router \| null }` |
| POST | `/routers/select/:id` | - | - | `{ success, data: Router }` |
| POST | `/routers/sync-all` | - | - | `{ success, data: { message } }` |
| GET | `/routers/:router_id` | - | - | `{ success, data: Router }` |
| PUT | `/routers/:router_id` | - | `UpdateRouterRequest` | `{ success, data: Router }` |
| DELETE | `/routers/:router_id` | - | - | `{ success, data: { message } }` |
| POST | `/routers/:router_id/sync` | - | - | `{ success, data: { message } }` |
| POST | `/routers/:router_id/test-connection` | - | - | `{ success, data: { message } }` |

**Note:** Router List handler uses `response.OK` (no meta), not `response.WithMeta`. Pagination total count not returned. Frontend should handle gracefully.

### Bandwidth Profiles (`/api/v1/routers/:router_id/bandwidth-profiles`)
| Method | Path | Body | Response |
|--------|------|------|----------|
| GET | `...bandwidth-profiles?limit=N&offset=N` | - | `{ success, data: BandwidthProfile[], meta }` |
| POST | `...bandwidth-profiles` | `CreateBandwidthProfileRequest` | `{ success, data: BandwidthProfile }` |
| GET | `...bandwidth-profiles/:id` | - | `{ success, data: BandwidthProfile }` (with mikrotik live data) |
| PUT | `...bandwidth-profiles/:id` | `UpdateBandwidthProfileRequest` | `{ success, data: BandwidthProfile }` |
| DELETE | `...bandwidth-profiles/:id` | - | `{ success, data: { message } }` |

### Subscriptions (`/api/v1/routers/:router_id/subscriptions`)
| Method | Path | Body | Response |
|--------|------|------|----------|
| GET | `...subscriptions?limit=N&offset=N` | - | `{ success, data: Subscription[], meta }` |
| POST | `...subscriptions` | `CreateSubscriptionRequest` | `{ success, data: Subscription }` |
| GET | `...subscriptions/:id` | - | `{ success, data: Subscription }` (with mikrotik live data) |
| PUT | `...subscriptions/:id` | `UpdateSubscriptionRequest` | `{ success, data: Subscription }` |
| DELETE | `...subscriptions/:id` | - | `{ success, data: { message } }` |
| POST | `...subscriptions/:id/activate` | - | `{ success, data: { message } }` |
| POST | `...subscriptions/:id/isolate` | `{ reason? }` | `{ success, data: { message } }` |
| POST | `...subscriptions/:id/restore` | - | `{ success, data: { message } }` |
| POST | `...subscriptions/:id/suspend` | `{ reason? }` | `{ success, data: { message } }` |
| POST | `...subscriptions/:id/terminate` | - | `{ success, data: { message } }` |

## Common Pitfalls

### Pitfall 1: Router-Scoped Resources
**What goes wrong:** Subscriptions and bandwidth profiles are nested under `/routers/:router_id/...`. Forgetting the router_id scope causes 404 errors.
**Why it happens:** Developers treat subscriptions as top-level resources.
**How to avoid:** Always require a selected router before showing subscription/bandwidth profile pages. Store selected router in Zustand or URL state.
**Warning signs:** 404 errors when listing subscriptions without router context.

### Pitfall 2: Router List Endpoint Has No Meta
**What goes wrong:** `GET /api/v1/routers` uses `response.OK` not `response.WithMeta`, so no total count is returned. Frontend pagination breaks if it expects meta.
**Why it happens:** Inconsistency in backend handler implementation.
**How to avoid:** For the router list, either (a) assume no pagination needed (typically <50 routers) and fetch all, or (b) handle missing meta gracefully by checking `data.meta?.total`.
**Warning signs:** `meta` is `undefined` in router list response.

### Pitfall 3: TanStack Virtual + Table Integration
**What goes wrong:** TanStack Virtual requires a fixed container height and a ref to the scroll container. Using it inside a TanStack Table with server-side pagination adds complexity.
**Why it happens:** Virtual and Table are separate libraries that must be manually wired.
**How to avoid:** Only use Virtual for the subscription table where 500+ rows are expected. For customers and routers (typically <100), standard pagination is sufficient. When using Virtual, set `overscan: 5` and fixed `estimateSize: () => 48`.
**Warning signs:** Blank rows, flickering scroll, incorrect row heights.

### Pitfall 4: Customer Create Has Optional Subscription
**What goes wrong:** `POST /api/v1/customers` can optionally create a subscription if `plan_id` is provided. Form must handle both modes.
**Why it happens:** Backend combines customer + subscription creation in one endpoint.
**How to avoid:** Make the subscription fields conditional in the form -- only show when "Create with subscription" toggle is enabled.
**Warning signs:** Required field errors when submitting without subscription data.

### Pitfall 5: Subscription Status Transitions
**What goes wrong:** Not all status transitions are valid. Calling activate on a terminated subscription will error.
**Why it happens:** Backend enforces state machine logic.
**How to avoid:** Show only valid action buttons based on current status:
  - `pending` -> activate
  - `active` -> isolate, suspend, terminate
  - `isolated` -> restore, terminate
  - `suspended` -> restore, terminate
  - `terminated` -> (no actions)
**Warning signs:** 400 errors from lifecycle action endpoints.

### Pitfall 6: Search/Filter Not Supported Server-Side for Customers
**What goes wrong:** `GET /api/v1/customers` only accepts `limit` and `offset`. No `search` or `status` query params.
**Why it happens:** Backend handler has simple List implementation.
**How to avoid:** Implement client-side filtering on the current page of data using TanStack Table's built-in column filtering. For full-text search across all pages, fetch a larger page size or accept the limitation.
**Warning signs:** Search bar that doesn't filter results.

## Code Examples

### Example 1: Zod Schema for API Response Envelope (Reuse from Phase 1)
```typescript
// Already exists at src/lib/schemas/auth.ts
export const ApiResponseSchema = <T extends z.ZodTypeAny>(dataSchema: T) =>
  z.object({
    success: z.boolean(),
    data: dataSchema,
    error: z.string().optional().nullable(),
    meta: z.object({
      total: z.number(),
      limit: z.number(),
      offset: z.number(),
    }).optional().nullable(),
  })
```

### Example 2: DataTable with TanStack Table v8
```typescript
import {
  useReactTable,
  getCoreRowModel,
  flexRender,
  type ColumnDef,
  type PaginationState,
} from "@tanstack/react-table"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"

interface DataTableProps<TData> {
  columns: ColumnDef<TData, unknown>[]
  data: TData[]
  pageCount: number
  pagination: PaginationState
  onPaginationChange: (updater: PaginationState | ((old: PaginationState) => PaginationState)) => void
}

export function DataTable<TData>({ columns, data, pageCount, pagination, onPaginationChange }: DataTableProps<TData>) {
  const table = useReactTable({
    data,
    columns,
    pageCount,
    state: { pagination },
    onPaginationChange,
    getCoreRowModel: getCoreRowModel(),
    manualPagination: true,
  })
  // ... render table header/body/footer with flexRender
}
```

### Example 3: Subscription Status Badge Colors
```typescript
const STATUS_VARIANTS: Record<string, string> = {
  pending: "bg-yellow-100 text-yellow-800 dark:bg-yellow-900/30 dark:text-yellow-400",
  active: "bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400",
  isolated: "bg-orange-100 text-orange-800 dark:bg-orange-900/30 dark:text-orange-400",
  suspended: "bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-400",
  terminated: "bg-gray-100 text-gray-800 dark:bg-gray-900/30 dark:text-gray-400",
}
```

### Example 4: TanStack Virtual Integration for Large Tables
```typescript
import { useVirtualizer } from "@tanstack/react-virtual"

// Inside component with a ref to scroll container
const parentRef = useRef<HTMLDivElement>(null)
const virtualizer = useVirtualizer({
  count: rows.length,
  getScrollElement: () => parentRef.current,
  estimateSize: () => 48,
  overscan: 5,
})
```

## State of the Art

| Old Approach | Current Approach | When Changed | Impact |
|--------------|------------------|--------------|--------|
| React Hook Form + Yup | TanStack Form v1 + Zod | Phase 1 decision | Already established -- use TanStack Form |
| Client-side only tables | TanStack Table v8 + server pagination | Standard since TanStack Table v8 | Use manual pagination mode |
| CSS table styling | Shadcn Table component (Radix) | shadcn v4 | Use `npx shadcn add table` |

## Open Questions

1. **Customer search across all pages**
   - What we know: Backend `GET /customers` only supports `limit`/`offset`, no `search` param
   - What's unclear: Whether to implement a debounced search that re-fetches or client-side filter only
   - Recommendation: Client-side filter on current page data. If full search is critical, increase default page size to 100 and add a note for future backend enhancement.

2. **Router selection persistence**
   - What we know: `POST /routers/select/:id` and `GET /routers/selected` exist server-side per user
   - What's unclear: Whether to also cache in Zustand for instant access
   - Recommendation: Store selected router in Zustand (fetched on app load), update on select action. This avoids an API call every time subscription page loads.

3. **Subscription list without router selection**
   - What we know: Subscriptions are scoped to `/routers/:router_id/subscriptions`
   - What's unclear: What to show on `/subscriptions` if no router is selected
   - Recommendation: Show a prompt to select a router first, or automatically use the selected router from Zustand/API. The sidebar already links to `/subscriptions` -- this page must handle the "no router selected" state.

## Sources

### Primary (HIGH confidence)
- Backend source code: `internal/router/admin.go` -- all API routes verified
- Backend source code: `internal/handler/customer_handler.go`, `router_handler.go`, `subscription_handler.go`, `bandwidth_profile_handler.go` -- all DTO shapes verified
- Backend source code: `internal/dto/*.go` -- all request/response structures verified
- Backend source code: `pkg/response/response.go` -- API envelope structure verified
- Frontend source code: `dashboard/package.json` -- all library versions verified (installed)
- Frontend source code: Phase 1 SUMMARY files -- established patterns verified

### Secondary (MEDIUM confidence)
- TanStack Table v8 patterns: based on training data + installed version ^8.21.3
- TanStack Virtual v3 patterns: based on training data + installed version ^3.13.23
- TanStack Form v1 patterns: verified from existing login.tsx implementation

## Metadata

**Confidence breakdown:**
- Standard stack: HIGH - all libraries already installed, versions verified from package.json
- Architecture: HIGH - patterns established in Phase 1, backend API fully inspected
- Pitfalls: HIGH - derived from direct code inspection of backend handlers and response patterns
- API endpoints: HIGH - verified from `internal/router/admin.go` source code

**Research date:** 2026-03-30
**Valid until:** 2026-04-30 (stable -- no external dependencies, all based on inspected source code)
