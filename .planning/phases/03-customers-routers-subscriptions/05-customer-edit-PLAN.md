---
phase: 03-customers-routers-subscriptions
plan: 05
type: execute
wave: 1
depends_on: []
files_modified:
  - website/src/hooks/use-customers.ts
  - website/src/features/customers/components/edit-customer-dialog.tsx
  - website/src/features/customers/index.tsx
  - website/src/features/customers/data/columns.tsx
autonomous: true
gap_closure: true
requirements:
  - CUST-03

must_haves:
  truths:
    - "Admin can update customer details (name, phone, email, address, username, static IP)"
  artifacts:
    - path: "website/src/hooks/use-customers.ts"
      provides: "useUpdateCustomer hook"
      exports: ["useUpdateCustomer"]
    - path: "website/src/features/customers/components/edit-customer-dialog.tsx"
      provides: "Edit customer form dialog pre-populated with current values"
    - path: "website/src/features/customers/data/columns.tsx"
      provides: "Edit action in customer dropdown menu"
      contains: "onEdit"
    - path: "website/src/features/customers/index.tsx"
      provides: "Edit dialog state and rendering"
      contains: "editTarget"
  key_links:
    - from: "website/src/features/customers/data/columns.tsx"
      to: "website/src/features/customers/index.tsx"
      via: "onEdit callback in ColumnActions"
      pattern: "onEdit"
    - from: "website/src/features/customers/components/edit-customer-dialog.tsx"
      to: "website/src/hooks/use-customers.ts"
      via: "useUpdateCustomer import"
      pattern: "useUpdateCustomer"
---

<objective>
Add customer edit capability with a pre-populated form dialog and update hook.

Purpose: Close verification gap CUST-03 (admin can update customer details). The API function `updateCustomer` already exists in customer.ts but no hook or UI dialog was built.

Output: Edit customer dialog with pre-populated fields, useUpdateCustomer hook, edit action in customer table row dropdown.
</objective>

<execution_context>
@$HOME/.claude/get-shit-done/workflows/execute-plan.md
@$HOME/.claude/get-shit-done/templates/summary.md
</execution_context>

<context>
@.planning/PROJECT.md
@.planning/ROADMAP.md
@.planning/STATE.md
@.planning/phases/03-customers-routers-subscriptions/02-customers-SUMMARY.md
@website/src/api/customer.ts
@website/src/hooks/use-customers.ts
@website/src/features/customers/index.tsx
@website/src/features/customers/data/columns.tsx
@website/src/features/customers/data/schema.ts
@website/src/features/customers/components/create-customer-dialog.tsx
@website/src/lib/schemas/customer.ts
</context>

<interfaces>
<!-- Existing types and contracts the executor needs -->

From website/src/api/customer.ts:
```typescript
export async function updateCustomer(
    id: string,
    data: Record<string, unknown>
): Promise<CustomerResponse> {
    const response = await adminClient.put(`/customers/${id}`, data)
    const parsed = CustomerDetailResponseSchema.parse(response.data)
    return parsed.data
}
```

From website/src/features/customers/data/schema.ts:
```typescript
export const createCustomerSchema = z.object({
    full_name: z.string().min(1, 'Full name is required'),
    phone: z.string().min(1, 'Phone is required'),
    email: z.string().email('Invalid email address').optional().or(z.literal('')),
    address: z.string().optional(),
    latitude: z.number().optional(),
    longitude: z.number().optional(),
    plan_id: z.string().uuid().optional(),
    router_id: z.string().uuid().optional(),
    username: z.string().optional(),
    password: z.string().optional(),
    static_ip: z.string().optional(),
});
export type CreateCustomer = z.infer<typeof createCustomerSchema>;
```

From website/src/lib/schemas/customer.ts:
```typescript
export const CustomerResponseSchema = z.object({
    id: z.string().uuid(),
    customer_code: z.string(),
    full_name: z.string(),
    username: z.string().nullable(),
    email: z.string().nullable(),
    phone: z.string(),
    id_card_number: z.string().nullable(),
    address: z.string().nullable(),
    latitude: z.number().nullable(),
    longitude: z.number().nullable(),
    is_active: z.boolean(),
    notes: z.string().nullable(),
    tags: z.array(z.string()),
    created_at: z.string(),
    updated_at: z.string(),
});
export type CustomerResponse = z.infer<typeof CustomerResponseSchema>;
```

From website/src/features/customers/data/columns.tsx:
```typescript
interface ColumnActions {
    onActivate: (customer: CustomerResponse) => void
    onDeactivate: (customer: CustomerResponse) => void
    onDelete: (customer: CustomerResponse) => void
}
```
</interfaces>

<tasks>

<task type="auto">
  <name>Task 1: Add useUpdateCustomer hook and edit customer dialog, wire into columns and page</name>
  <files>website/src/hooks/use-customers.ts, website/src/features/customers/components/edit-customer-dialog.tsx, website/src/features/customers/data/columns.tsx, website/src/features/customers/index.tsx</files>
  <read_first>
    - website/src/hooks/use-customers.ts
    - website/src/api/customer.ts
    - website/src/features/customers/components/create-customer-dialog.tsx
    - website/src/features/customers/data/columns.tsx
    - website/src/features/customers/data/schema.ts
    - website/src/features/customers/index.tsx
    - website/src/lib/schemas/customer.ts
  </read_first>
  <action>
1. In `website/src/hooks/use-customers.ts`, add `updateCustomer` to the import from `@/api/customer`. Then add the following hook after `useDeleteCustomer`:

   ```typescript
   export function useUpdateCustomer() {
       const queryClient = useQueryClient()

       return useMutation({
           mutationFn: ({ id, data }: { id: string; data: Record<string, unknown> }) => updateCustomer(id, data),
           onSuccess: () => {
               queryClient.invalidateQueries({ queryKey: ['customers'] })
               toast.success('Customer updated successfully')
           },
           onError: (err: unknown) => {
               const message =
                   (err as { response?: { data?: { error?: string } } })?.response?.data
                       ?.error ?? 'Failed to update customer'
               toast.error(message)
           },
       })
   }
   ```

2. Create `website/src/features/customers/components/edit-customer-dialog.tsx`:

   Follow the exact structure of `create-customer-dialog.tsx` with these differences:
   - Props: `{ customer: CustomerResponse | null; open: boolean; onOpenChange: (open: boolean) => void }`
   - Import `CustomerResponse` from `@/lib/schemas/customer`
   - Import `useUpdateCustomer` from `@/hooks/use-customers`
   - Use `useUpdateCustomer` instead of `useCreateCustomer`
   - Destructure: `const { mutateAsync: updateCustomer, isPending } = useUpdateCustomer()`
   - Title: "Edit Customer" instead of "Add New Customer"
   - DialogDescription: "Update customer details. Fields marked with * are required."
   - Submit button: loading text "Saving..." instead of loading spinner, normal text "Save Changes" instead of "Create Customer"
   - Default values populated from `customer` prop:
     ```
     defaultValues: {
         full_name: customer?.full_name ?? '',
         phone: customer?.phone ?? '',
         email: customer?.email ?? '',
         address: customer?.address ?? '',
         username: customer?.username ?? '',
         password: '',
         static_ip: '',
     },
     ```
     Note: password defaults to empty string (not pre-filled). Only send password in payload if user explicitly fills it.
   - On submit: build payload like create-customer-dialog does, but also add `id: customer.id`:
     ```
     const payload: Record<string, unknown> = {
         full_name: data.full_name,
         phone: data.phone,
     }
     if (data.email) payload.email = data.email
     if (data.address) payload.address = data.address
     if (data.username) payload.username = data.username
     if (data.password) payload.password = data.password
     if (data.static_ip) payload.static_ip = data.static_ip

     await updateCustomer({ id: customer.id, data: payload })
     ```
   - Use `createCustomerSchema` from `../data/schema` for validation (same schema works)
   - Add `Pencil` icon from lucide-react to the DialogTitle (same pattern as delete-customer-dialog uses Trash2)
   - `zodResolver(createCustomerSchema) as never` (same pattern as create dialog)

3. In `website/src/features/customers/data/columns.tsx`:

   a. Add `onEdit` to the `ColumnActions` interface:
      ```
      onEdit: (customer: CustomerResponse) => void
      ```

   b. Add `Pencil` to the lucide-react import (alongside existing MoreHorizontal, User, Trash2, Power, PowerOff).

   c. Add an edit menu item in the dropdown, AFTER the Activate/Deactivate toggle and BEFORE the DropdownMenuSeparator + Delete:
      ```
      <DropdownMenuItem onClick={() => actions.onEdit(customer)}>
          <Pencil className='mr-2 size-4' />
          Edit
      </DropdownMenuItem>
      ```

4. In `website/src/features/customers/index.tsx`:

   a. Add state: `const [editTarget, setEditTarget] = useState<CustomerResponse | null>(null)`

   b. Import `EditCustomerDialog` from `./components/edit-customer-dialog`

   c. Update `createCustomerColumns` call to include `onEdit`:
      ```
      const customerColumns = createCustomerColumns({
          onActivate: (customer) => activateCustomer(customer.id),
          onDeactivate: (customer) => deactivateCustomer(customer.id),
          onDelete: (customer) => setDeleteTarget(customer),
          onEdit: (customer) => setEditTarget(customer),
      })
      ```

   d. Add `EditCustomerDialog` rendering after `CreateCustomerDialog`:
      ```
      <EditCustomerDialog
          customer={editTarget}
          open={!!editTarget}
          onOpenChange={(open) => {
              if (!open) setEditTarget(null)
          }}
      />
      ```
  </action>
  <verify>
    <automated>cd /home/butterfly-student/Programing/Mikrotik/mikmongo-fully/website && npx tsc --noEmit 2>&1 | head -30</automated>
  </verify>
  <done>
    - `useUpdateCustomer` exported from website/src/hooks/use-customers.ts
    - website/src/features/customers/components/edit-customer-dialog.tsx exists and exports EditCustomerDialog
    - website/src/features/customers/data/columns.tsx has `onEdit` in ColumnActions interface
    - website/src/features/customers/index.tsx imports EditCustomerDialog and renders it
    - website/src/features/customers/index.tsx has `editTarget` state and passes `onEdit` to createCustomerColumns
    - TypeScript compilation passes
  </done>
</task>

</tasks>

<verification>
1. TypeScript compilation passes: `cd website && npx tsc --noEmit`
2. Hook exists: `grep "export function useUpdateCustomer" website/src/hooks/use-customers.ts`
3. Dialog exists: `test -f website/src/features/customers/components/edit-customer-dialog.tsx`
4. Columns wired: `grep "onEdit" website/src/features/customers/data/columns.tsx`
5. Page wired: `grep "editTarget\|EditCustomerDialog\|onEdit" website/src/features/customers/index.tsx`
</verification>

<success_criteria>
- Admin can click "Edit" in a customer row dropdown to open a pre-populated edit dialog
- Edit dialog shows current customer data (name, phone, email, address, username)
- Saving calls PUT /api/v1/customers/{id} with updated data
- Success toast shown after update, customer list refreshes
- TypeScript compiles cleanly
</success_criteria>

<output>
After completion, create `.planning/phases/03-customers-routers-subscriptions/05-customer-edit-SUMMARY.md`
</output>
