// src/routes/customer/_customerAuth/dashboard.tsx — renders at /customer/dashboard
import { createFileRoute } from "@tanstack/react-router"

export const Route = createFileRoute("/customer/_customerAuth/dashboard")({
  component: CustomerDashboardPage,
})

function CustomerDashboardPage() {
  return (
    <div className="p-6">
      <h1 className="text-2xl font-bold">Customer Portal</h1>
      <p className="text-muted-foreground mt-2">Coming in Phase 5</p>
    </div>
  )
}
