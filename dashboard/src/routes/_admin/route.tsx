// Pathless layout route — wraps ALL admin routes at /, /customers, etc.
// The underscore prefix means this does NOT add a path segment.
// Auth guard runs in beforeLoad before any child route renders.
import { createFileRoute, redirect, Outlet } from "@tanstack/react-router"
import { AppShell } from "@/components/layout/admin/AppShell"

export const Route = createFileRoute("/_admin")({
  beforeLoad: ({ context, location }) => {
    if (!context.adminAuth.isAuthenticated) {
      throw redirect({
        to: "/login",
        search: {
          // Preserve the URL the user was trying to access for post-login redirect
          redirect: location.href,
        },
      })
    }
  },
  component: AdminLayout,
})

function AdminLayout() {
  return (
    <AppShell>
      <Outlet />
    </AppShell>
  )
}
