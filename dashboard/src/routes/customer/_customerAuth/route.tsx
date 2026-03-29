// src/routes/customer/_customerAuth/route.tsx
// Pathless layout: wraps all customer-portal routes. Auth guard in Plan 01-02.
import { createFileRoute, Outlet } from "@tanstack/react-router"

export const Route = createFileRoute("/customer/_customerAuth")({
  component: () => <Outlet />,
})
