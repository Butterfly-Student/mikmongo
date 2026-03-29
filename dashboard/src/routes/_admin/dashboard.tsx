// src/routes/_admin/dashboard.tsx — renders at URL: /dashboard
import { createFileRoute } from "@tanstack/react-router"

export const Route = createFileRoute("/_admin/dashboard")({
  component: DashboardPage,
})

function DashboardPage() {
  return (
    <div className="p-6">
      <h1 className="text-2xl font-bold">Dashboard</h1>
      <p className="text-muted-foreground mt-2">Overview — coming in Plan 01-03</p>
    </div>
  )
}
