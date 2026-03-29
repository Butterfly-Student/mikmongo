// Root route with typed context — auth state injected from main.tsx
import { createRootRouteWithContext, Outlet } from "@tanstack/react-router"
// NOTE: TanStackRouterDevtools intentionally omitted here — added in Plan 01-03
// after @tanstack/router-devtools is installed in Plan 01-01
import { ReactQueryDevtools } from "@tanstack/react-query-devtools"
import type { QueryClient } from "@tanstack/react-query"
import type { AdminRole } from "@/lib/rbac"

export interface RouterContext {
  adminAuth: {
    isAuthenticated: boolean
    role: AdminRole | null
    accessToken: string | null
  }
  agentAuth: {
    isAuthenticated: boolean
    accessToken: string | null
  }
  customerAuth: {
    isAuthenticated: boolean
    accessToken: string | null
  }
  queryClient: QueryClient
}

export const Route = createRootRouteWithContext<RouterContext>()({
  component: () => (
    <>
      <Outlet />
      {import.meta.env.DEV && (
        <>
          <ReactQueryDevtools initialIsOpen={false} />
        </>
      )}
    </>
  ),
})
