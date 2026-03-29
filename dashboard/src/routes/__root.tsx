// src/routes/__root.tsx
// Root route — provides typed router context for all portals.
// Full implementation: Plan 01-02 (auth store wired here).
import { createRootRouteWithContext, Outlet } from "@tanstack/react-router"
import { QueryClient } from "@tanstack/react-query"

export interface RouterContext {
  adminAuth: {
    isAuthenticated: boolean
    role: string | null
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
  component: () => <Outlet />,
})
