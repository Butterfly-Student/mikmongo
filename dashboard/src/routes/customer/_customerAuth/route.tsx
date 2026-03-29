// Pathless layout for customer portal auth guard
// Renders children at /customer/dashboard, /customer/invoices, etc.
import { createFileRoute, redirect, Outlet } from "@tanstack/react-router"

export const Route = createFileRoute("/customer/_customerAuth")({
  beforeLoad: ({ context, location }) => {
    if (!context.customerAuth.isAuthenticated) {
      throw redirect({
        to: "/customer/login",
        search: { redirect: location.href },
      })
    }
  },
  component: () => <Outlet />,
})
