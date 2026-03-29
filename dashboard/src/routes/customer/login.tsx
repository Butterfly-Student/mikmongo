// src/routes/customer/login.tsx — customer portal login at /customer/login
import { createFileRoute } from "@tanstack/react-router"

export const Route = createFileRoute("/customer/login")({
  component: CustomerLoginPage,
})

function CustomerLoginPage() {
  return (
    <div className="flex min-h-screen items-center justify-center">
      <p className="text-muted-foreground">Customer Login — coming in Plan 01-02</p>
    </div>
  )
}
