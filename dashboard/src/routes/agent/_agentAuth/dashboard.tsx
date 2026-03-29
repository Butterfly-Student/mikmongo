// src/routes/agent/_agentAuth/dashboard.tsx — renders at /agent/dashboard
import { createFileRoute } from "@tanstack/react-router"

export const Route = createFileRoute("/agent/_agentAuth/dashboard")({
  component: AgentDashboardPage,
})

function AgentDashboardPage() {
  return (
    <div className="p-6">
      <h1 className="text-2xl font-bold">Agent Portal</h1>
      <p className="text-muted-foreground mt-2">Coming in Phase 5</p>
    </div>
  )
}
