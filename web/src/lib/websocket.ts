import { useAuthStore } from '@/stores/auth-store'

/**
 * Creates a WebSocket connection to the given path with a Bearer token.
 * Token is passed as a query param because WS headers aren't supported by browsers.
 */
export function createWebSocket(path: string, token: string): WebSocket {
  const base = window.location.origin.replace(/^http/, 'ws')
  const url = `${base}${path}?token=${encodeURIComponent(token)}`
  return new WebSocket(url)
}

/** Creates an admin WebSocket using the token from the auth store. */
export function createAdminWs(path: string): WebSocket {
  const token = useAuthStore.getState().adminAccessToken ?? ''
  return createWebSocket(path, token)
}

/** Creates a customer portal WebSocket using the portal token from the auth store. */
export function createCustomerWs(path: string): WebSocket {
  const token = useAuthStore.getState().customerToken ?? ''
  return createWebSocket(path, token)
}

/** Creates an agent portal WebSocket using the agent token from the auth store. */
export function createAgentWs(path: string): WebSocket {
  const token = useAuthStore.getState().agentToken ?? ''
  return createWebSocket(path, token)
}

// ── Convenience helpers for each real-time endpoint ──────────────────

export const ws = {
  pppActive: (routerId: string) =>
    createAdminWs(`/api/v1/routers/${routerId}/ppp/ws/active`),

  pppInactive: (routerId: string) =>
    createAdminWs(`/api/v1/routers/${routerId}/ppp/ws/inactive`),

  hotspotActive: (routerId: string) =>
    createAdminWs(`/api/v1/routers/${routerId}/hotspot/ws/active`),

  hotspotInactive: (routerId: string) =>
    createAdminWs(`/api/v1/routers/${routerId}/hotspot/ws/inactive`),

  systemResource: (routerId: string) =>
    createAdminWs(`/api/v1/routers/${routerId}/monitor/ws/system-resource`),

  traffic: (routerId: string, iface: string) =>
    createAdminWs(`/api/v1/routers/${routerId}/monitor/ws/traffic/${iface}`),

  logs: (routerId: string) =>
    createAdminWs(`/api/v1/routers/${routerId}/monitor/ws/logs`),

  ping: (routerId: string) =>
    createAdminWs(`/api/v1/routers/${routerId}/monitor/ws/ping`),

  rawListen: (routerId: string) =>
    createAdminWs(`/api/v1/routers/${routerId}/raw/ws/listen`),
}
