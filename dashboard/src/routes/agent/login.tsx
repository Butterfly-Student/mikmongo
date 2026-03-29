// src/routes/agent/login.tsx — agent portal login at /agent/login
import { createFileRoute } from "@tanstack/react-router"

export const Route = createFileRoute("/agent/login")({
  component: AgentLoginPage,
})

function AgentLoginPage() {
  return (
    <div className="flex min-h-screen items-center justify-center">
      <p className="text-muted-foreground">Agent Login — coming in Plan 01-02</p>
    </div>
  )
}
