// Pathless layout for agent portal auth guard
// Renders children at /agent/dashboard, /agent/profile, etc.
import { createFileRoute, redirect, Outlet } from "@tanstack/react-router"

export const Route = createFileRoute("/agent/_agentAuth")({
  beforeLoad: ({ context, location }) => {
    if (!context.agentAuth.isAuthenticated) {
      throw redirect({
        to: "/agent/login",
        search: { redirect: location.href },
      })
    }
  },
  component: () => <Outlet />,
})
