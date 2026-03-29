// CRITICAL: Zustand persist is async. Gate RouterProvider behind isHydrated
// to prevent beforeLoad guards from firing before tokens load from localStorage.
import React from "react"
import ReactDOM from "react-dom/client"
import { RouterProvider, createRouter } from "@tanstack/react-router"
import { QueryClientProvider } from "@tanstack/react-query"
import { Toaster } from "sonner"
import { routeTree } from "./routeTree.gen"
import { queryClient } from "./lib/queryClient"
import { useAdminAuthContext, useAgentAuthContext, useCustomerAuthContext } from "./hooks/useAuth"
import { useStore } from "./store"
// Global styles — MUST be imported here so Tailwind v4 + Shadcn CSS variables load
import "./styles/globals.css"

const router = createRouter({
  routeTree,
  context: {
    // Placeholder values — real values injected via RouterProvider context prop below
    adminAuth: { isAuthenticated: false, role: null, accessToken: null },
    agentAuth: { isAuthenticated: false, accessToken: null },
    customerAuth: { isAuthenticated: false, accessToken: null },
    queryClient,
  },
})

declare module "@tanstack/react-router" {
  interface Register {
    router: typeof router
  }
}

function App() {
  const adminAuth = useAdminAuthContext()
  const agentAuth = useAgentAuthContext()
  const customerAuth = useCustomerAuthContext()
  const isHydrated = useStore((s) => s.isHydrated)

  // Wait for Zustand to rehydrate from localStorage before rendering routes.
  // Without this gate, beforeLoad guards see null tokens and redirect to login.
  if (!isHydrated) {
    return (
      <div className="flex h-screen items-center justify-center">
        <div className="h-8 w-8 animate-spin rounded-full border-4 border-primary border-t-transparent" />
      </div>
    )
  }

  return (
    <RouterProvider
      router={router}
      context={{ adminAuth, agentAuth, customerAuth, queryClient }}
    />
  )
}

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <QueryClientProvider client={queryClient}>
      <App />
      <Toaster richColors position="top-right" />
    </QueryClientProvider>
  </React.StrictMode>
)
