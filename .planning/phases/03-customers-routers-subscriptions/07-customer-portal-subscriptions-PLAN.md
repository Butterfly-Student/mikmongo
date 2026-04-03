---
phase: 03-customers-routers-subscriptions
plan: 07
type: execute
wave: 1
depends_on: []
files_modified:
  - website/src/api/portal/subscription.ts
  - website/src/hooks/use-customer-portal.ts
  - website/src/features/customer-portal/subscriptions.tsx
  - website/src/routes/customer/subscriptions/index.tsx
  - website/src/routes/customer/index.tsx
autonomous: true
gap_closure: true
requirements:
  - SUB-05

must_haves:
  truths:
    - "Customer can view their active subscriptions after logging into the customer portal"
    - "Customer subscription data is fetched from the portal-scoped API endpoint"
  artifacts:
    - path: "website/src/api/portal/subscription.ts"
      provides: "Portal subscription API using customerClient"
      exports: ["listPortalSubscriptions"]
    - path: "website/src/hooks/use-customer-portal.ts"
      provides: "usePortalSubscriptions hook"
      exports: ["usePortalSubscriptions"]
    - path: "website/src/features/customer-portal/subscriptions.tsx"
      provides: "Customer portal subscriptions page component"
    - path: "website/src/routes/customer/index.tsx"
      provides: "Customer portal index route with outlet"
    - path: "website/src/routes/customer/subscriptions/index.tsx"
      provides: "Customer subscriptions route file"
  key_links:
    - from: "website/src/routes/customer/subscriptions/index.tsx"
      to: "website/src/features/customer-portal/subscriptions.tsx"
      via: "import and render Subscriptions component"
      pattern: "CustomerPortalSubscriptions|Subscriptions"
    - from: "website/src/features/customer-portal/subscriptions.tsx"
      to: "website/src/hooks/use-customer-portal.ts"
      via: "usePortalSubscriptions import"
      pattern: "usePortalSubscriptions"
    - from: "website/src/hooks/use-customer-portal.ts"
      to: "website/src/api/portal/subscription.ts"
      via: "listPortalSubscriptions import"
      pattern: "listPortalSubscriptions"
    - from: "website/src/api/portal/subscription.ts"
      to: "website/src/lib/axios/customer-client.ts"
      via: "customerClient import"
      pattern: "customerClient"
---

<objective>
Create a customer portal subscriptions page that shows the logged-in customer's subscriptions.

Purpose: Close verification gap SUB-05 (customer portal shows their active subscriptions). The customer portal route exists at /customer but contains only the auth guard with no child routes.

Output: Customer-facing subscriptions page at /customer/subscriptions showing the customer's subscription list with status, plan details, and expiry dates.
</objective>

<execution_context>
@$HOME/.claude/get-shit-done/workflows/execute-plan.md
@$HOME/.claude/get-shit-done/templates/summary.md
</execution_context>

<context>
@.planning/PROJECT.md
@.planning/ROADMAP.md
@.planning/STATE.md
@.planning/phases/03-customers-routers-subscriptions/03-subscriptions-SUMMARY.md
@website/src/routes/customer/route.tsx
@website/src/lib/axios/customer-client.ts
@website/src/lib/schemas/subscription.ts
@website/src/stores/auth-store.ts
</context>

<interfaces>
<!-- Existing types and contracts the executor needs -->

From website/src/lib/axios/customer-client.ts:
```typescript
export const customerClient = axios.create({
  baseURL: '/portal/v1',
  headers: { 'Content-Type': 'application/json' },
})
// Has Bearer token interceptor reading from useAuthStore.getState().customerToken
// Has 401 interceptor that clears auth and redirects to /customer/login
```

From website/src/lib/schemas/subscription.ts:
```typescript
export const SubscriptionResponseSchema = z.object({
    id: z.string().uuid(),
    customer_id: z.string().uuid(),
    plan_id: z.string().uuid().nullable(),
    router_id: z.string().uuid(),
    username: z.string(),
    static_ip: z.string().nullable(),
    gateway: z.string().nullable(),
    status: z.enum(['pending', 'active', 'suspended', 'isolated', 'expired', 'terminated']),
    activated_at: z.string().nullable(),
    expiry_date: z.string().nullable(),
    billing_day: z.number().int().nullable(),
    auto_isolate: z.boolean().nullable(),
    grace_period_days: z.number().int().nullable(),
    suspend_reason: z.string().nullable(),
    notes: z.string().nullable(),
    created_at: z.string(),
    updated_at: z.string(),
    mikrotik: z.object({
        service: z.string().nullable(),
        profile: z.string().nullable(),
        local_address: z.string().nullable(),
        remote_address: z.string().nullable(),
    }).nullable(),
});
export type SubscriptionResponse = z.infer<typeof SubscriptionResponseSchema>;
```

OpenAPI customer portal subscriptions endpoint:
- GET /portal/v1/subscriptions
- Security: PortalAuth (Bearer token)
- Response: SuccessResponse with data: array of SubscriptionResponse

From website/src/stores/auth-store.ts:
```typescript
customerToken: string | null
customerUser: CustomerUser | null
customerIsAuthenticated: boolean
```

TanStack Router file-based routing convention:
- `website/src/routes/customer/route.tsx` exists (auth guard + Outlet)
- New child route: `website/src/routes/customer/subscriptions/index.tsx`
- Index route for customer: `website/src/routes/customer/index.tsx` (layout with Outlet)
</interfaces>

<tasks>

<task type="auto">
  <name>Task 1: Create portal subscription API, hook, and customer subscriptions page</name>
  <files>website/src/api/portal/subscription.ts, website/src/hooks/use-customer-portal.ts, website/src/features/customer-portal/subscriptions.tsx, website/src/routes/customer/index.tsx, website/src/routes/customer/subscriptions/index.tsx</files>
  <read_first>
    - website/src/lib/axios/customer-client.ts
    - website/src/lib/schemas/subscription.ts
    - website/src/routes/customer/route.tsx
    - website/src/stores/auth-store.ts
    - website/src/lib/schemas/auth.ts
  </read_first>
  <action>
1. Create directory `website/src/api/portal/` if it does not exist.

2. Create `website/src/api/portal/subscription.ts`:

   ```typescript
   import { customerClient } from '@/lib/axios/customer-client'
   import { SubscriptionResponseSchema } from '@/lib/schemas/subscription'
   import { z } from 'zod'

   const PortalSubscriptionListResponseSchema = z.object({
       success: z.boolean(),
       data: z.array(SubscriptionResponseSchema),
   })

   export async function listPortalSubscriptions(): Promise<z.infer<typeof PortalSubscriptionListResponseSchema>['data']> {
       const response = await customerClient.get('/subscriptions')
       const parsed = PortalSubscriptionListResponseSchema.parse(response.data)
       return parsed.data
   }
   ```

   Note: Uses `customerClient` (baseURL '/portal/v1') so the actual endpoint is `/portal/v1/subscriptions`.

3. Create `website/src/hooks/use-customer-portal.ts`:

   ```typescript
   import { useQuery } from '@tanstack/react-query'
   import { listPortalSubscriptions } from '@/api/portal/subscription'

   export function usePortalSubscriptions() {
       return useQuery({
           queryKey: ['portal-subscriptions'],
           queryFn: () => listPortalSubscriptions(),
           staleTime: 2 * 60 * 1000,
       })
   }
   ```

4. Create `website/src/routes/customer/index.tsx`:

   This is the customer portal index/layout route. The auth guard is already in `customer/route.tsx` (the parent). This file provides the index page that redirects to subscriptions or shows a welcome.

   ```typescript
   import { createFileRoute, redirect } from '@tanstack/react-router'

   export const Route = createFileRoute('/customer/')({
     beforeLoad: () => {
       throw redirect({ to: '/customer/subscriptions' })
     },
   })
   ```

   This redirects `/customer` to `/customer/subscriptions` automatically.

5. Create directory `website/src/routes/customer/subscriptions/` if it does not exist.

6. Create `website/src/routes/customer/subscriptions/index.tsx`:

   ```typescript
   import { createFileRoute } from '@tanstack/react-router'
   import { CustomerPortalSubscriptions } from '@/features/customer-portal/subscriptions'

   export const Route = createFileRoute('/customer/subscriptions/')({
     component: () => <CustomerPortalSubscriptions />,
   })
   ```

7. Create directory `website/src/features/customer-portal/` if it does not exist.

8. Create `website/src/features/customer-portal/subscriptions.tsx`:

   This is the customer-facing subscriptions page. It should show a simple, clean card-based layout of the customer's subscriptions. Use these design elements:
   - Page title "My Subscriptions" with description "View your active and past subscriptions"
   - Loading state: skeleton cards (3 placeholders)
   - Empty state: "No subscriptions found. Contact your provider to set up a subscription."
   - Each subscription displayed as a Card (from @/components/ui/card) with:
     - CardHeader: status badge + username
     - CardContent: details grid showing:
       - Status (colored badge: active=green/default, suspended=orange/secondary, terminated=red/destructive, expired=secondary, isolated=outline, pending=outline)
       - IP Address (static_ip or "Dynamic")
       - Expiry Date (formatted with Intl.DateTimeFormat)
       - Activated At (formatted date)
       - MikroTik Profile (from mikrotik.profile or "-")
     - Use a grid layout: `grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4`

   Implementation:
   ```typescript
   import { usePortalSubscriptions } from '@/hooks/use-customer-portal'
   import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
   import { Badge } from '@/components/ui/badge'
   import { Skeleton } from '@/components/ui/skeleton'
   import { Wifi, WifiOff } from 'lucide-react'
   import type { SubscriptionResponse } from '@/lib/schemas/subscription'

   function formatDate(dateString: string | null): string {
       if (!dateString) return '-'
       return new Intl.DateTimeFormat('id-ID', {
           dateStyle: 'long',
       }).format(new Date(dateString))
   }

   function getStatusBadge(status: SubscriptionResponse['status']) {
       const variants: Record<string, 'default' | 'secondary' | 'destructive' | 'outline'> = {
           active: 'default',
           pending: 'outline',
           suspended: 'secondary',
           isolated: 'outline',
           expired: 'secondary',
           terminated: 'destructive',
       }
       return <Badge variant={variants[status] ?? 'secondary'}>{status.charAt(0).toUpperCase() + status.slice(1)}</Badge>
   }

   export function CustomerPortalSubscriptions() {
       const { data: subscriptions, isLoading } = usePortalSubscriptions()

       if (isLoading) {
           return (
               <div className="space-y-6">
                   <div>
                       <h1 className="text-2xl font-semibold tracking-tight">My Subscriptions</h1>
                       <p className="text-sm text-muted-foreground">View your active and past subscriptions</p>
                   </div>
                   <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                       {[1, 2, 3].map((i) => (
                           <Card key={i}>
                               <CardHeader><Skeleton className="h-6 w-32" /></CardHeader>
                               <CardContent className="space-y-2">
                                   <Skeleton className="h-4 w-full" />
                                   <Skeleton className="h-4 w-3/4" />
                                   <Skeleton className="h-4 w-1/2" />
                               </CardContent>
                           </Card>
                       ))}
                   </div>
               </div>
           )
       }

       return (
           <div className="space-y-6">
               <div>
                   <h1 className="text-2xl font-semibold tracking-tight">My Subscriptions</h1>
                   <p className="text-sm text-muted-foreground">View your active and past subscriptions</p>
               </div>
               {(!subscriptions || subscriptions.length === 0) ? (
                   <div className="flex flex-col items-center justify-center py-12 text-center">
                       <WifiOff className="size-12 text-muted-foreground mb-4" />
                       <h3 className="text-lg font-medium">No Subscriptions Found</h3>
                       <p className="text-sm text-muted-foreground mt-1">
                           Contact your provider to set up a subscription.
                       </p>
                   </div>
               ) : (
                   <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                       {subscriptions.map((sub) => (
                           <Card key={sub.id}>
                               <CardHeader className="flex flex-row items-center justify-between pb-2">
                                   <CardTitle className="text-base font-medium flex items-center gap-2">
                                       <Wifi className="size-4" />
                                       {sub.username}
                                   </CardTitle>
                                   {getStatusBadge(sub.status)}
                               </CardHeader>
                               <CardContent className="space-y-2 text-sm">
                                   <div className="flex justify-between">
                                       <span className="text-muted-foreground">IP Address</span>
                                       <span>{sub.static_ip ?? 'Dynamic'}</span>
                                   </div>
                                   <div className="flex justify-between">
                                       <span className="text-muted-foreground">Profile</span>
                                       <span>{sub.mikrotik?.profile ?? '-'}</span>
                                   </div>
                                   <div className="flex justify-between">
                                       <span className="text-muted-foreground">Expiry</span>
                                       <span>{formatDate(sub.expiry_date)}</span>
                                   </div>
                                   <div className="flex justify-between">
                                       <span className="text-muted-foreground">Activated</span>
                                       <span>{formatDate(sub.activated_at)}</span>
                                   </div>
                               </CardContent>
                           </Card>
                       ))}
                   </div>
               )}
           </div>
       )
   }
   ```
  </action>
  <verify>
    <automated>cd /home/butterfly-student/Programing/Mikrotik/mikmongo-fully/website && npx tsc --noEmit 2>&1 | head -30</automated>
  </verify>
  <done>
    - website/src/api/portal/subscription.ts exists and exports listPortalSubscriptions
    - website/src/hooks/use-customer-portal.ts exists and exports usePortalSubscriptions
    - website/src/features/customer-portal/subscriptions.tsx exists and exports CustomerPortalSubscriptions
    - website/src/routes/customer/index.tsx exists (redirects to /customer/subscriptions)
    - website/src/routes/customer/subscriptions/index.tsx exists (renders CustomerPortalSubscriptions)
    - All files use customerClient (not adminClient) for portal API calls
    - TypeScript compilation passes
  </done>
</task>

</tasks>

<verification>
1. TypeScript compilation passes: `cd website && npx tsc --noEmit`
2. Portal API: `grep "listPortalSubscriptions" website/src/api/portal/subscription.ts`
3. Portal hook: `grep "usePortalSubscriptions" website/src/hooks/use-customer-portal.ts`
4. Route files: `test -f website/src/routes/customer/subscriptions/index.tsx`
5. Customer client usage: `grep "customerClient" website/src/api/portal/subscription.ts`
6. No adminClient usage in portal files: `grep -r "adminClient" website/src/api/portal/` returns nothing
</verification>

<success_criteria>
- Customer portal route /customer/subscriptions renders a subscription list page
- Page fetches from GET /portal/v1/subscriptions using customerClient (Bearer token from Zustand)
- Shows subscription cards with status badges, username, IP, profile, expiry date
- Empty state shown when no subscriptions exist
- Loading skeleton shown during fetch
- /customer redirects to /customer/subscriptions
- TypeScript compiles cleanly
</success_criteria>

<output>
After completion, create `.planning/phases/03-customers-routers-subscriptions/07-customer-portal-subscriptions-SUMMARY.md`
</output>
