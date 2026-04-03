---
phase: 03-customers-routers-subscriptions
plan: 06
type: execute
wave: 1
depends_on: []
files_modified:
  - website/src/api/profiles.ts
  - website/src/hooks/use-profiles.ts
  - website/src/features/profiles/components/edit-profile-dialog.tsx
  - website/src/features/profiles/data/columns.tsx
  - website/src/features/profiles/index.tsx
autonomous: true
gap_closure: true
requirements:
  - BW-03

must_haves:
  truths:
    - "Admin can update an existing bandwidth profile's name, speed, price, and other settings"
  artifacts:
    - path: "website/src/api/profiles.ts"
      provides: "updateProfile API function"
      exports: ["updateProfile"]
    - path: "website/src/hooks/use-profiles.ts"
      provides: "useUpdateProfile hook"
      exports: ["useUpdateProfile"]
    - path: "website/src/features/profiles/components/edit-profile-dialog.tsx"
      provides: "Edit profile form dialog"
    - path: "website/src/features/profiles/data/columns.tsx"
      provides: "Edit action in profile dropdown menu"
      contains: "onEdit"
    - path: "website/src/features/profiles/index.tsx"
      provides: "Edit dialog state and rendering"
      contains: "editTarget"
  key_links:
    - from: "website/src/features/profiles/components/edit-profile-dialog.tsx"
      to: "website/src/hooks/use-profiles.ts"
      via: "useUpdateProfile import"
      pattern: "useUpdateProfile"
    - from: "website/src/features/profiles/index.tsx"
      to: "website/src/features/profiles/data/columns.tsx"
      via: "onEdit callback in createColumns"
      pattern: "onEdit"
---

<objective>
Add bandwidth profile update capability with an edit dialog.

Purpose: Close verification gap BW-03 (admin can update bandwidth profiles). Create and delete already work; update is missing.

Output: `updateProfile` API function, `useUpdateProfile` hook, edit profile dialog pre-populated with current values, edit action in profile table dropdown.
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
@website/src/api/profiles.ts
@website/src/hooks/use-profiles.ts
@website/src/features/profiles/index.tsx
@website/src/features/profiles/data/columns.tsx
@website/src/features/profiles/data/schema.ts
@website/src/features/profiles/components/create-profile-dialog.tsx
</context>

<interfaces>
<!-- Existing types and contracts the executor needs -->

From website/src/api/profiles.ts:
```typescript
export async function createProfile(routerId: string, data: CreateProfile): Promise<Profile> {
    const response = await adminClient.post(`/routers/${routerId}/bandwidth-profiles`, data)
    const parsed = SingleProfileResponseSchema.parse(response.data)
    return parsed.data
}
```

From website/src/features/profiles/data/schema.ts:
```typescript
export const createProfileSchema = z.object({
    profile_code: z.string().min(1, { message: "Profile code is required" }),
    name: z.string().min(1, { message: "Name is required" }),
    description: z.string().optional(),
    download_speed: z.coerce.number().int().min(1),
    upload_speed: z.coerce.number().int().min(1),
    price_monthly: z.coerce.number().min(0.01),
    tax_rate: z.coerce.number().optional(),
    billing_cycle: z.enum(["daily", "weekly", "monthly", "yearly"]).default("monthly"),
    billing_day: z.coerce.number().int().optional(),
    grace_period_days: z.coerce.number().int().optional(),
    isolate_profile_name: z.string().optional(),
    sort_order: z.coerce.number().int().optional(),
    is_visible: z.boolean().default(true),
    mt_local_address: z.string().optional(),
    mt_remote_address: z.string().optional(),
    mt_parent_queue: z.string().optional(),
    mt_queue_type: z.string().optional(),
    mt_dns_server: z.string().optional(),
    mt_session_timeout: z.string().optional(),
    mt_idle_timeout: z.string().optional(),
});
export type CreateProfile = z.infer<typeof createProfileSchema>;

export type Profile = z.infer<typeof profileSchema>;
```

From website/src/features/profiles/data/columns.tsx:
```typescript
interface ColumnActions {
  onDelete: (profile: Profile) => void
}
```

OpenAPI UpdateBandwidthProfileRequest fields (all optional): profile_code, name, description, download_speed, upload_speed, price_monthly, tax_rate, billing_cycle, billing_day, grace_period_days, isolate_profile_name, sort_order, is_visible, is_active, mt_local_address, mt_remote_address, mt_parent_queue, mt_queue_type, mt_dns_server, mt_session_timeout, mt_idle_timeout

API endpoint: PUT /api/v1/routers/{router_id}/bandwidth-profiles/{id}
</interfaces>

<tasks>

<task type="auto">
  <name>Task 1: Add updateProfile API, useUpdateProfile hook, edit dialog, and wire into UI</name>
  <files>website/src/api/profiles.ts, website/src/hooks/use-profiles.ts, website/src/features/profiles/components/edit-profile-dialog.tsx, website/src/features/profiles/data/columns.tsx, website/src/features/profiles/index.tsx</files>
  <read_first>
    - website/src/api/profiles.ts
    - website/src/hooks/use-profiles.ts
    - website/src/features/profiles/components/create-profile-dialog.tsx
    - website/src/features/profiles/data/columns.tsx
    - website/src/features/profiles/data/schema.ts
    - website/src/features/profiles/index.tsx
  </read_first>
  <action>
1. In `website/src/api/profiles.ts`, add after `deleteProfile`:

   ```typescript
   export async function updateProfile(routerId: string, id: string, data: Partial<CreateProfile>): Promise<Profile> {
       const response = await adminClient.put(`/routers/${routerId}/bandwidth-profiles/${id}`, data)
       const parsed = SingleProfileResponseSchema.parse(response.data)
       return parsed.data
   }
   ```

   Import `Partial` from TypeScript (no import needed, it is a built-in utility type). The `CreateProfile` type is already imported.

2. In `website/src/hooks/use-profiles.ts`, add after `useDeleteProfile`:

   Add `updateProfile` to the existing import from `@/api/profiles`.

   ```typescript
   export function useUpdateProfile() {
       const queryClient = useQueryClient()

       return useMutation({
           mutationFn: ({ routerId, id, data }: { routerId: string; id: string; data: Partial<CreateProfile> }) =>
               updateProfile(routerId, id, data),
           onSuccess: (_, variables) => {
               queryClient.invalidateQueries({ queryKey: ['profiles', variables.routerId] })
               toast.success("Bandwidth profile updated successfully")
           },
           onError: (err: unknown) => {
               const message =
                   (err as { response?: { data?: { error?: string } } })?.response?.data
                       ?.error ?? 'Failed to update profile'
               toast.error(message)
           },
       })
   }
   ```

   Note: Use `unknown` error type (not `any`) matching the project convention.

3. Create `website/src/features/profiles/components/edit-profile-dialog.tsx`:

   Follow the exact same structure as `create-profile-dialog.tsx` but with these differences:
   - Props: `{ profile: Profile | null; routerId: string | null; open: boolean; onOpenChange: (open: boolean) => void }`
   - Import `Profile` from `../data/schema`
   - Import `useUpdateProfile` from `@/hooks/use-profiles`
   - Use `useUpdateProfile` instead of `useCreateProfile`
   - Destructure: `const { mutateAsync: updateProfile, isPending } = useUpdateProfile()`
   - Title: "Edit Bandwidth Profile" with Pencil icon
   - Submit button: "Save Changes" instead of "Create Profile", loading text "Saving..."
   - Default values populated from `profile` prop:
     ```
     defaultValues: {
         profile_code: profile?.profile_code ?? '',
         name: profile?.name ?? '',
         description: profile?.description ?? '',
         download_speed: profile?.download_speed ?? 1,
         upload_speed: profile?.upload_speed ?? 1,
         price_monthly: profile?.price_monthly ?? 0,
         tax_rate: profile?.tax_rate ?? undefined,
         billing_cycle: profile?.billing_cycle ?? 'monthly',
         billing_day: profile?.billing_day ?? undefined,
         grace_period_days: profile?.grace_period_days ?? undefined,
         isolate_profile_name: profile?.isolate_profile_name ?? '',
         sort_order: profile?.sort_order ?? undefined,
         is_visible: profile?.is_visible ?? true,
         mt_local_address: profile?.mikrotik?.local_address ?? '',
         mt_remote_address: profile?.mikrotik?.remote_address ?? '',
         mt_parent_queue: profile?.mikrotik?.parent_queue ?? '',
         mt_queue_type: profile?.mikrotik?.queue_type ?? '',
         mt_dns_server: profile?.mikrotik?.dns_server ?? '',
         mt_session_timeout: profile?.mikrotik?.session_timeout ?? '',
         mt_idle_timeout: profile?.mikrotik?.idle_timeout ?? '',
     },
     ```
   - On submit: call `await updateProfile({ routerId: routerId!, id: profile.id, data })` then close dialog
   - Use same form layout with grid cols as create dialog
   - `zodResolver(createProfileSchema) as never` (same pattern)

4. In `website/src/features/profiles/data/columns.tsx`:

   a. Add `onEdit` to `ColumnActions`:
      ```
      onEdit: (profile: Profile) => void
      ```

   b. Add `Pencil` to the lucide-react import.

   c. Add edit menu item before the Delete item:
      ```
      <DropdownMenuItem onClick={() => actions.onEdit(profile)}>
        <Pencil className="mr-2 size-4" /> Edit
      </DropdownMenuItem>
      ```

5. In `website/src/features/profiles/index.tsx`:

   a. Add state: `const [editTarget, setEditTarget] = useState<Profile | null>(null)`

   b. Import `EditProfileDialog` from `./components/edit-profile-dialog`

   c. Import `Profile` from `./data/schema` (if not already imported)

   d. Update the columns call to include onEdit:
      ```
      const columns = createColumns({
          onEdit: (profile) => setEditTarget(profile),
          onDelete: (profile) => {
              if (selectedRouterId) {
                  deleteProfile({ routerId: selectedRouterId, id: profile.id })
              }
          },
      })
      ```

   e. Add EditProfileDialog after CreateProfileDialog:
      ```
      <EditProfileDialog
          profile={editTarget}
          routerId={selectedRouterId}
          open={!!editTarget}
          onOpenChange={(open) => { if (!open) setEditTarget(null) }}
      />
      ```
  </action>
  <verify>
    <automated>cd /home/butterfly-student/Programing/Mikrotik/mikmongo-fully/website && npx tsc --noEmit 2>&1 | head -30</automated>
  </verify>
  <done>
    - `updateProfile` exported from website/src/api/profiles.ts
    - `useUpdateProfile` exported from website/src/hooks/use-profiles.ts
    - website/src/features/profiles/components/edit-profile-dialog.tsx exists and exports EditProfileDialog
    - website/src/features/profiles/data/columns.tsx has `onEdit` in ColumnActions
    - website/src/features/profiles/index.tsx has `editTarget` state and renders EditProfileDialog
    - TypeScript compilation passes
  </done>
</task>

</tasks>

<verification>
1. TypeScript compilation passes: `cd website && npx tsc --noEmit`
2. API function: `grep "export async function updateProfile" website/src/api/profiles.ts`
3. Hook: `grep "export function useUpdateProfile" website/src/hooks/use-profiles.ts`
4. Dialog: `test -f website/src/features/profiles/components/edit-profile-dialog.tsx`
5. Columns wired: `grep "onEdit" website/src/features/profiles/data/columns.tsx`
6. Page wired: `grep "editTarget\|EditProfileDialog" website/src/features/profiles/index.tsx`
</verification>

<success_criteria>
- Admin can click "Edit" in a bandwidth profile row dropdown to open a pre-populated edit dialog
- Edit dialog shows current profile data (name, speeds, price, billing cycle, etc.)
- Saving calls PUT /api/v1/routers/{router_id}/bandwidth-profiles/{id}
- Success toast shown after update, profile list refreshes
- TypeScript compiles cleanly
</success_criteria>

<output>
After completion, create `.planning/phases/03-customers-routers-subscriptions/06-profile-update-SUMMARY.md`
</output>
