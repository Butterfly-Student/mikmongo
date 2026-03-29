// src/routes/agent/_agentAuth/route.tsx
// Pathless layout: wraps all agent-portal routes. Auth guard in Plan 01-02.
import { createFileRoute, Outlet } from "@tanstack/react-router"

export const Route = createFileRoute("/agent/_agentAuth")({
  component: () => <Outlet />,
})
