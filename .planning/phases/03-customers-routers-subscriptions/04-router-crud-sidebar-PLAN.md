---
phase: 03-customers-routers-subscriptions
plan: 04
type: execute
wave: 1
depends_on: []
files_modified:
  - website/src/api/router.ts
  - website/src/hooks/use-routers.ts
  - website/src/features/routers/index.tsx
  - website/src/features/routers/components/edit-router-dialog.tsx
  - website/src/features/routers/components/delete-router-dialog.tsx
  - website/src/features/routers/data/columns.tsx
  - website/src/components/layout/data/sidebar-data.ts
autonomous: true
gap_closure: true
requirements:
  - RTR-02
  - RTR-06

must_haves:
  truths:
    - "Admin can edit a router's name, address, credentials, ports, and other settings"
    - "Admin can delete a router with a confirmation dialog"
    - "Admin can sync all routers simultaneously via a button"
    - "Customers and Routers are reachable from sidebar navigation"
  artifacts:
    - path: "website/src/api/router.ts"
      provides: "updateRouter, deleteRouter, syncAllRouters API functions"
      exports: ["updateRouter", "deleteRouter", "syncAllRouters"]
    - path: "website/src/hooks/use-routers.ts"
      provides: "useUpdateRouter, useDeleteRouter, useSyncAllRouters hooks"
      exports: ["useUpdateRouter", "useDeleteRouter", "useSyncAllRouters"]
    - path: "website/src/features/routers/components/edit-router-dialog.tsx"
      provides: "Edit router form dialog"
    - path: "website/src/features/routers/components/delete-router-dialog.tsx"
      provides: "Delete router confirmation dialog"
    - path: "website/src/features/routers/data/columns.tsx"
      provides: "Edit action in dropdown menu"
      contains: "onEdit"
    - path: "website/src/components/layout/data/sidebar-data.ts"
      provides: "Enabled Customers and Routers sidebar nav items"
      contains: "url: '/customers'"
      contains: "url: '/routers'"
  key_links:
    - from: "website/src/features/routers/index.tsx"
      to: "website/src/hooks/use-routers.ts"
      via: "useUpdateRouter, useDeleteRouter, useSyncAllRouters imports"
      pattern: "useUpdateRouter|useDeleteRouter|useSyncAllRouters"
    - from: "website/src/features/routers/data/columns.tsx"
      to: "website/src/features/routers/components/edit-router-dialog.tsx"
      via: "onEdit callback in ColumnActions"
      pattern: "onEdit"
    - from: "website/src/components/layout/data/sidebar-data.ts"
      to: "/customers"
      via: "url property"
      pattern: "url: '/customers'"
---

<objective>
Complete router CRUD by adding update, delete, and sync-all operations. Also enable sidebar navigation for Customers and Routers pages.

Purpose: Close verification gaps RTR-02 (edit/delete router) and RTR-06 (sync all routers), plus the discovered sidebar navigation gap that makes implemented pages unreachable.

Output: Full router CRUD (create + edit + delete + sync + sync-all + test connection), working sidebar nav for Customers and Routers.
</objective>

<execution_context>
@$HOME/.claude/get-shit-done/workflows/execute-plan.md
@$HOME/.claude/get-shit-done/templates/summary.md
</execution_context>

<context>
@.planning/PROJECT.md
@.planning/ROADMAP.md
@.planning/STATE.md
@.planning/phases/03-customers-routers-subscriptions/01-routers-profiles-SUMMARY.md
@website/src/api/router.ts
@website/src/hooks/use-routers.ts
@website/src/features/routers/index.tsx
@website/src/features/routers/data/columns.tsx
@website/src/features/routers/data/schema.ts
@website/src/features/routers/components/create-router-dialog.tsx
@website/src/lib/schemas/router.ts
@website/src/components/layout/data/sidebar-data.ts
</context>

<interfaces>
<!-- Existing types and contracts the executor needs -->

From website/src/features/routers/data/schema.ts:
```typescript
export const createRouterSchema = z.object({
    name: z.string().min(1, { message: "Name is required" }),
    address: z.string().min(1, { message: "Address is required" }),
    username: z.string().min(1, { message: "Username is required" }),
    password: z.string().min(1, { message: "Password is required" }),
    area: z.string().optional(),
    api_port: z.coerce.number().int().default(8728),
    rest_port: z.coerce.number().int().default(80),
    use_ssl: z.boolean().default(false),
    is_master: z.boolean().default(false),
    notes: z.string().optional(),
});
export type CreateRouter = z.infer<typeof createRouterSchema>;
```

From website/src/lib/schemas/router.ts:
```typescript
export const RouterResponseSchema = z.object({
    id: z.string(),
    name: z.string(),
    address: z.string(),
    area: z.string(),
    api_port: z.number().optional(),
    rest_port: z.number().optional(),
    username: z.string(),
    use_ssl: z.boolean(),
    is_master: z.boolean(),
    is_active: z.boolean(),
    status: z.enum(['online', 'offline', 'unknown']),
    last_seen_at: z.string().nullable(),
    notes: z.string().optional(),
    created_at: z.string(),
    updated_at: z.string(),
});
export type RouterResponse = z.infer<typeof RouterResponseSchema>;
```

OpenAPI UpdateRouterRequest fields (all optional): name, address, username, password, area, api_port, rest_port, use_ssl, is_master, notes

From website/src/features/routers/data/columns.tsx:
```typescript
interface ColumnActions {
  onSync: (router: RouterResponse) => void
  onTestConnection: (router: RouterResponse) => void
  onSelectActive: (router: RouterResponse) => void
  onDeleteRouter: (router: RouterResponse) => void
}
```
</interfaces>

<tasks>

<task type="auto">
  <name>Task 1: Add router API functions, hooks, and sidebar navigation fix</name>
  <files>website/src/api/router.ts, website/src/hooks/use-routers.ts, website/src/components/layout/data/sidebar-data.ts</files>
  <read_first>
    - website/src/api/router.ts
    - website/src/hooks/use-routers.ts
    - website/src/components/layout/data/sidebar-data.ts
    - website/src/lib/schemas/router.ts
  </read_first>
  <action>
1. In `website/src/api/router.ts`, add three new exported async functions after the existing `testConnection` function:

   a. `updateRouter(id: string, data: Record<string, unknown>): Promise<RouterResponse>` -- sends PUT to `/routers/${id}` via adminClient, parses response with `SelectedRouterResponseSchema`, returns `parsed.data`. Use `Record<string, unknown>` for the data parameter (same pattern as `createRouter`).

   b. `deleteRouter(id: string): Promise<void>` -- sends DELETE to `/routers/${id}` via adminClient. No parsing needed (returns MessageResponse).

   c. `syncAllRouters(): Promise<void>` -- sends POST to `/routers/sync-all` via adminClient. No parsing needed.

   Add `import type { RouterResponse }` is already imported. No new imports needed since `SelectedRouterResponseSchema` is already imported.

2. In `website/src/hooks/use-routers.ts`, add three new exported hooks after the existing `useTestRouterConnection` hook:

   a. `useUpdateRouter()` -- useMutation, mutationFn receives `{ id: string; data: Record<string, unknown> }`, calls `updateRouter(id, data)`. onSuccess: invalidate `['routers']` query key, toast.success("Router updated successfully"). onError: use `unknown` type (not `any`), narrow with type assertion `(err as { response?: { data?: { error?: string } } })?.response?.data?.error ?? "Failed to update router"`, toast.error(message).

   b. `useDeleteRouter()` -- useMutation, mutationFn receives `id: string`, calls `deleteRouter(id)`. onSuccess: invalidate `['routers']` query key, toast.success("Router deleted successfully"). onError: same `unknown` error pattern as above.

   c. `useSyncAllRouters()` -- useMutation, mutationFn receives no argument, calls `syncAllRouters()`. onSuccess: invalidate `['routers']` query key, toast.success("All routers synced successfully"). onError: same `unknown` error pattern.

   Add imports for `updateRouter`, `deleteRouter`, `syncAllRouters` from `@/api/router` (add to existing import line).

3. In `website/src/components/layout/data/sidebar-data.ts`, make two changes:

   a. Change the Customers item (currently lines 51-55) from:
      ```
      {
        title: 'Customers',
        url: '#',
        icon: UsersRound,
        disabled: true,
      },
      ```
      to:
      ```
      {
        title: 'Customers',
        url: '/customers',
        icon: UsersRound,
      },
      ```

   b. Change the Routers item (currently lines 57-61) from:
      ```
      {
        title: 'Routers',
        url: '#',
        icon: Server,
        disabled: true,
      },
      ```
      to:
      ```
      {
        title: 'Routers',
        url: '/routers',
        icon: Server,
      },
      ```

   Remove the `disabled: true` property from both items.
  </action>
  <verify>
    <automated>cd /home/butterfly-student/Programing/Mikrotik/mikmongo-fully/website && npx tsc --noEmit 2>&1 | head -30</automated>
  </verify>
  <done>
    - `updateRouter`, `deleteRouter`, `syncAllRouters` exported from website/src/api/router.ts
    - `useUpdateRouter`, `useDeleteRouter`, `useSyncAllRouters` exported from website/src/hooks/use-routers.ts
    - All three hooks use `unknown` error type (not `any`)
    - sidebar-data.ts has `url: '/customers'` without `disabled` property
    - sidebar-data.ts has `url: '/routers'` without `disabled` property
    - TypeScript compilation passes
  </done>
</task>

<task type="auto">
  <name>Task 2: Create edit and delete router dialogs, wire into router page and columns</name>
  <files>website/src/features/routers/components/edit-router-dialog.tsx, website/src/features/routers/components/delete-router-dialog.tsx, website/src/features/routers/data/columns.tsx, website/src/features/routers/index.tsx</files>
  <read_first>
    - website/src/features/routers/components/create-router-dialog.tsx
    - website/src/features/customers/components/delete-customer-dialog.tsx
    - website/src/features/routers/data/columns.tsx
    - website/src/features/routers/index.tsx
    - website/src/features/routers/data/schema.ts
    - website/src/hooks/use-routers.ts
    - website/src/lib/schemas/router.ts
  </read_first>
  <action>
1. Create `website/src/features/routers/components/edit-router-dialog.tsx`:

   This is a dialog component for editing an existing router. Follow the exact same structure and styling as `create-router-dialog.tsx` but with these differences:
   - Props: `{ router: RouterResponse | null; open: boolean; onOpenChange: (open: boolean) => void }`
   - Import `RouterResponse` from `@/lib/schemas/router`
   - Import `useUpdateRouter` from `@/hooks/use-routers`
   - Use `useUpdateRouter` instead of `useCreateRouter`
   - Title: "Edit Router" instead of "Add New Router"
   - Submit button text: "Save Changes" instead of "Add Router", loading text: "Saving..."
   - Default values populated from `router` prop: name, address, username, area, api_port, rest_port, use_ssl, is_master, notes
   - password field should be empty by default (not pre-filled from router data) -- the password is optional for updates
   - On submit: call `mutateAsync({ id: router.id, data })` where data is the form values (omit empty password)
   - Use `createRouterSchema` from `../data/schema` for validation (same schema works since the fields are the same)
   - Add `Pencil` icon from lucide-react to the DialogTitle

2. Create `website/src/features/routers/components/delete-router-dialog.tsx`:

   Follow the exact same pattern as `website/src/features/customers/components/delete-customer-dialog.tsx`:
   - Props: `{ router: RouterResponse | null; open: boolean; onOpenChange: (open: boolean) => void }`
   - Import `RouterResponse` from `@/lib/schemas/router`
   - Import `useDeleteRouter` from `@/hooks/use-routers`
   - Title: "Delete Router" with Trash2 icon
   - Description: "Are you sure you want to delete router `{router?.name}`? This action cannot be undone. All associated bandwidth profiles and subscriptions will be affected."
   - Cancel button + destructive "Delete Router" button with Loader2 spinner when pending
   - handleDelete: `await deleteRouter(router.id)` then close dialog

3. In `website/src/features/routers/data/columns.tsx`, add edit action:

   a. Add `onEdit` to the `ColumnActions` interface:
      ```
      onEdit: (router: RouterResponse) => void
      ```

   b. Add an edit menu item BEFORE the Delete menu item (after the second DropdownMenuSeparator):
      ```
      <DropdownMenuItem onClick={() => actions.onEdit(router)}>
        <Pencil className="mr-2 size-4" /> Edit Router
      </DropdownMenuItem>
      ```

   c. Add `Pencil` to the lucide-react import.

4. In `website/src/features/routers/index.tsx`:

   a. Add state: `const [editTarget, setEditTarget] = useState<RouterResponse | null>(null)`
   b. Add state: `const [deleteTarget, setDeleteTarget] = useState<RouterResponse | null>(null)`
   c. Import `useUpdateRouter`, `useDeleteRouter`, `useSyncAllRouters` from `@/hooks/use-routers`
   d. Import `useSyncAllRouters` and destructure: `const { mutate: syncAllRouters, isPending: isSyncingAll } = useSyncAllRouters()`
   e. Import `EditRouterDialog` from `./components/edit-router-dialog`
   f. Import `DeleteRouterDialog` from `./components/delete-router-dialog`
   g. Import `RouterResponse` from `@/lib/schemas/router`
   h. Import `Button` from `@/components/ui/button` and `RefreshCw` from lucide-react

   i. Update the columns call to include onEdit and pass the real delete handler:
      ```
      const columns = createColumns({
        onSync: (router) => syncRouter(router.id),
        onTestConnection: (router) => testConnection(router.id),
        onSelectActive: (router) => selectRouter(router.id),
        onEdit: (router) => setEditTarget(router),
        onDeleteRouter: (router) => setDeleteTarget(router),
      })
      ```
      This replaces the current `onDeleteRouter: (router) => console.log(...)` stub.

   j. Add a "Sync All Routers" button in the header area, next to the existing "Add Router" functionality. Add it after the `<p>` description and before `<RouterTable>`:
      ```
      <div className="flex items-center justify-between">
        <p className='text-sm text-muted-foreground'>Manage connected MikroTik routers</p>
        <Button
          variant="outline"
          size="sm"
          onClick={() => syncAllRouters()}
          disabled={isSyncingAll}
        >
          <RefreshCw className={`mr-2 size-4 ${isSyncingAll ? 'animate-spin' : ''}`} />
          {isSyncingAll ? 'Syncing...' : 'Sync All'}
        </Button>
      </div>
      ```
      Remove the original standalone `<p>` tag for the description.

   k. Add dialog components after `<CreateRouterDialog>`:
      ```
      <EditRouterDialog
        router={editTarget}
        open={!!editTarget}
        onOpenChange={(open) => { if (!open) setEditTarget(null) }}
      />
      <DeleteRouterDialog
        router={deleteTarget}
        open={!!deleteTarget}
        onOpenChange={(open) => { if (!open) setDeleteTarget(null) }}
      />
      ```
  </action>
  <verify>
    <automated>cd /home/butterfly-student/Programing/Mikrotik/mikmongo-fully/website && npx tsc --noEmit 2>&1 | head -30</automated>
  </verify>
  <done>
    - website/src/features/routers/components/edit-router-dialog.tsx exists and exports EditRouterDialog
    - website/src/features/routers/components/delete-router-dialog.tsx exists and exports DeleteRouterDialog
    - website/src/features/routers/data/columns.tsx has `onEdit` in ColumnActions interface
    - website/src/features/routers/index.tsx imports useUpdateRouter, useDeleteRouter, useSyncAllRouters
    - website/src/features/routers/index.tsx has NO console.log in onDeleteRouter handler
    - website/src/features/routers/index.tsx renders EditRouterDialog and DeleteRouterDialog
    - website/src/features/routers/index.tsx has "Sync All" button
    - TypeScript compilation passes
  </done>
</task>

</tasks>

<verification>
1. TypeScript compilation passes: `cd website && npx tsc --noEmit`
2. No console.log stubs in routers/index.tsx: `grep -n "console.log" website/src/features/routers/index.tsx` returns nothing
3. Sidebar navigation: `grep -A3 "title: 'Customers'" website/src/components/layout/data/sidebar-data.ts` shows `url: '/customers'` without `disabled`
4. Sidebar navigation: `grep -A3 "title: 'Routers'" website/src/components/layout/data/sidebar-data.ts` shows `url: '/routers'` without `disabled`
5. API functions exist: `grep "export async function updateRouter\|export async function deleteRouter\|export async function syncAllRouters" website/src/api/router.ts`
6. Hooks exist: `grep "export function useUpdateRouter\|export function useDeleteRouter\|export function useSyncAllRouters" website/src/hooks/use-routers.ts`
</verification>

<success_criteria>
- Admin can edit a router (edit dialog pre-populates fields, saves via PUT /api/v1/routers/{id})
- Admin can delete a router (confirmation dialog calls DELETE /api/v1/routers/{id})
- Admin can sync all routers (button calls POST /api/v1/routers/sync-all)
- No console.log stubs remain in router page
- Customers and Routers sidebar items navigate to /customers and /routers respectively
- TypeScript compiles cleanly
</success_criteria>

<output>
After completion, create `.planning/phases/03-customers-routers-subscriptions/04-router-crud-sidebar-SUMMARY.md`
</output>
