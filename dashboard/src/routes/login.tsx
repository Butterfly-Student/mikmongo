// src/routes/login.tsx — admin login page (Plan 01-02 adds form logic)
import { createFileRoute } from "@tanstack/react-router"

export const Route = createFileRoute("/login")({
  component: LoginPage,
})

function LoginPage() {
  return (
    <div className="flex min-h-screen items-center justify-center">
      <p className="text-muted-foreground">Admin Login — coming in Plan 01-02</p>
    </div>
  )
}
