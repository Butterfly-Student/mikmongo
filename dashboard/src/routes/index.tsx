// src/routes/index.tsx
// Redirects root path to /dashboard (admin default landing).
import { createFileRoute, redirect } from "@tanstack/react-router"

export const Route = createFileRoute("/")({
  beforeLoad: () => {
    throw redirect({ to: "/dashboard" })
  },
})
