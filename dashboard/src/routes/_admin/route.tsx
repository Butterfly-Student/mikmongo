// src/routes/_admin/route.tsx
// Pathless layout: wraps all admin routes with auth guard + AppShell.
// Auth guard implementation: Plan 01-02. AppShell: Plan 01-03.
import { createFileRoute, Outlet } from "@tanstack/react-router"

export const Route = createFileRoute("/_admin")({
  component: () => (
    <div>
      <Outlet />
    </div>
  ),
})
