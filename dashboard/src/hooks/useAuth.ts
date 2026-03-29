// Context-facing hooks — return shape matching RouterContext interface
import { useStore } from "@/store"

export function useAdminAuthContext() {
  const isAuthenticated = useStore((s) => s.adminIsAuthenticated)
  const role = useStore((s) => s.adminUser?.role ?? null)
  const accessToken = useStore((s) => s.adminAccessToken)
  return { isAuthenticated, role, accessToken }
}

export function useAgentAuthContext() {
  const isAuthenticated = useStore((s) => s.agentIsAuthenticated)
  const accessToken = useStore((s) => s.agentAccessToken)
  return { isAuthenticated, accessToken }
}

export function useCustomerAuthContext() {
  const isAuthenticated = useStore((s) => s.customerIsAuthenticated)
  const accessToken = useStore((s) => s.customerAccessToken)
  return { isAuthenticated, accessToken }
}
