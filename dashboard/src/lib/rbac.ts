// ─── Types ────────────────────────────────────────────────────────────────────
export type AdminRole = "superadmin" | "admin" | "teknisi"

// Resources map to route segments / feature areas
export type Resource =
  | "dashboard"
  | "customers"
  | "routers"
  | "bandwidth-profiles"
  | "subscriptions"
  | "invoices"
  | "payments"
  | "registrations"
  | "agents"
  | "agent-invoices"
  | "cash"
  | "reports"
  | "live-monitor"
  | "settings"
  | "users"         // superadmin only (SYS-02)

export type Action = "read" | "create" | "update" | "delete" | "manage"

// ─── Permission matrix ────────────────────────────────────────────────────────
// Defines the minimum role required for each (resource, action) pair.
// Roles are ordered: superadmin > admin > teknisi
// "manage" means all actions (create + read + update + delete + special actions)

const ROLE_LEVEL: Record<AdminRole, number> = {
  superadmin: 3,
  admin: 2,
  teknisi: 1,
}

// For each resource, record the minimum role level required per action.
// If a resource/action is not listed, default is superadmin-only.
type PermissionMatrix = Partial<
  Record<Resource, Partial<Record<Action, number>>>
>

const PERMISSIONS: PermissionMatrix = {
  dashboard: {
    read: ROLE_LEVEL.teknisi,       // All roles can see the dashboard
  },
  // Network management — teknisi's domain
  routers: {
    read: ROLE_LEVEL.teknisi,
    create: ROLE_LEVEL.admin,
    update: ROLE_LEVEL.teknisi,     // teknisi can update router config
    delete: ROLE_LEVEL.admin,
    manage: ROLE_LEVEL.admin,       // sync, test-connection actions
  },
  "bandwidth-profiles": {
    read: ROLE_LEVEL.teknisi,
    create: ROLE_LEVEL.teknisi,
    update: ROLE_LEVEL.teknisi,
    delete: ROLE_LEVEL.admin,
    manage: ROLE_LEVEL.teknisi,
  },
  subscriptions: {
    read: ROLE_LEVEL.teknisi,
    create: ROLE_LEVEL.admin,
    update: ROLE_LEVEL.teknisi,     // teknisi can isolate/restore/suspend
    delete: ROLE_LEVEL.admin,
    manage: ROLE_LEVEL.teknisi,     // lifecycle actions
  },
  // Customer management — admin and above
  customers: {
    read: ROLE_LEVEL.admin,
    create: ROLE_LEVEL.admin,
    update: ROLE_LEVEL.admin,
    delete: ROLE_LEVEL.admin,
    manage: ROLE_LEVEL.admin,
  },
  // Billing — admin and above
  invoices: {
    read: ROLE_LEVEL.admin,
    create: ROLE_LEVEL.admin,
    update: ROLE_LEVEL.admin,
    delete: ROLE_LEVEL.admin,
    manage: ROLE_LEVEL.admin,
  },
  payments: {
    read: ROLE_LEVEL.admin,
    create: ROLE_LEVEL.admin,
    update: ROLE_LEVEL.admin,
    delete: ROLE_LEVEL.admin,
    manage: ROLE_LEVEL.admin,
  },
  registrations: {
    read: ROLE_LEVEL.admin,
    manage: ROLE_LEVEL.admin,
  },
  agents: {
    read: ROLE_LEVEL.admin,
    create: ROLE_LEVEL.admin,
    update: ROLE_LEVEL.admin,
    delete: ROLE_LEVEL.admin,
    manage: ROLE_LEVEL.admin,
  },
  "agent-invoices": {
    read: ROLE_LEVEL.admin,
    manage: ROLE_LEVEL.admin,
  },
  cash: {
    read: ROLE_LEVEL.admin,
    create: ROLE_LEVEL.admin,
    update: ROLE_LEVEL.admin,
    delete: ROLE_LEVEL.admin,
    manage: ROLE_LEVEL.admin,
  },
  reports: {
    read: ROLE_LEVEL.admin,
    manage: ROLE_LEVEL.admin,
  },
  "live-monitor": {
    read: ROLE_LEVEL.teknisi,
    manage: ROLE_LEVEL.teknisi,
  },
  settings: {
    read: ROLE_LEVEL.admin,
    manage: ROLE_LEVEL.admin,
  },
  // User management — superadmin only (AUTH-07, SYS-02)
  users: {
    read: ROLE_LEVEL.superadmin,
    create: ROLE_LEVEL.superadmin,
    update: ROLE_LEVEL.superadmin,
    delete: ROLE_LEVEL.superadmin,
    manage: ROLE_LEVEL.superadmin,
  },
}

// ─── Public API ───────────────────────────────────────────────────────────────

/**
 * Check if a role has permission to perform an action on a resource.
 *
 * @example
 * hasPermission("teknisi", "users", "read")        // false
 * hasPermission("admin", "users", "read")           // false
 * hasPermission("superadmin", "users", "read")      // true
 * hasPermission("teknisi", "subscriptions", "manage") // true
 * hasPermission("admin", "customers", "read")       // true
 */
export function hasPermission(
  role: AdminRole,
  resource: Resource,
  action: Action
): boolean {
  const resourcePerms = PERMISSIONS[resource]
  if (!resourcePerms) {
    // Resource not in matrix — default deny
    return false
  }
  const requiredLevel = resourcePerms[action]
  if (requiredLevel === undefined) {
    // Action not defined for this resource — default deny
    return false
  }
  return ROLE_LEVEL[role] >= requiredLevel
}

/**
 * Get all resources a role can access (for sidebar menu filtering).
 */
export function getAccessibleResources(role: AdminRole): Resource[] {
  return (Object.keys(PERMISSIONS) as Resource[]).filter((resource) =>
    hasPermission(role, resource, "read")
  )
}
