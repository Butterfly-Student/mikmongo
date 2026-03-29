// Root route with typed context — auth state injected from main.tsx
import { createRootRouteWithContext, Outlet } from "@tanstack/react-router"
import { ReactQueryDevtools } from "@tanstack/react-query-devtools"
import type { QueryClient } from "@tanstack/react-query"
import { Toaster } from "sonner"
import { ThemeProvider } from "@/components/providers/ThemeProvider"
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

function RootComponent() {
  return (
    <ThemeProvider>
      <Outlet />
      <Toaster richColors position="top-right" />
      {import.meta.env.DEV && <ReactQueryDevtools initialIsOpen={false} />}
    </ThemeProvider>
  )
}

export const Route = createRootRouteWithContext<RouterContext>()({
  component: RootComponent,
})
