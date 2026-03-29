// src/routes/_admin/dashboard.tsx — redirects to / (overview moved to _admin/index.tsx)
import { createFileRoute, redirect } from "@tanstack/react-router"

export const Route = createFileRoute("/_admin/dashboard")({
  beforeLoad: () => {
    throw redirect({ to: "/" })
  },
  component: () => null,
})
