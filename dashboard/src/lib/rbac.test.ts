import { describe, it, expect } from "vitest"
import { hasPermission, getAccessibleResources } from "./rbac"

describe("hasPermission", () => {
  // ─── superadmin: full access ────────────────────────────────────────────────
  describe("superadmin", () => {
    it("can read users (superadmin-only resource)", () => {
      expect(hasPermission("superadmin", "users", "read")).toBe(true)
    })
    it("can delete users", () => {
      expect(hasPermission("superadmin", "users", "delete")).toBe(true)
    })
    it("can manage all admin resources", () => {
      expect(hasPermission("superadmin", "customers", "manage")).toBe(true)
      expect(hasPermission("superadmin", "invoices", "manage")).toBe(true)
      expect(hasPermission("superadmin", "settings", "manage")).toBe(true)
    })
    it("can access network management resources", () => {
      expect(hasPermission("superadmin", "routers", "read")).toBe(true)
      expect(hasPermission("superadmin", "subscriptions", "manage")).toBe(true)
    })
  })

  // ─── admin: everything except users ────────────────────────────────────────
  describe("admin", () => {
    it("cannot read users (superadmin-only)", () => {
      expect(hasPermission("admin", "users", "read")).toBe(false)
    })
    it("cannot create users", () => {
      expect(hasPermission("admin", "users", "create")).toBe(false)
    })
    it("can manage customers", () => {
      expect(hasPermission("admin", "customers", "read")).toBe(true)
      expect(hasPermission("admin", "customers", "create")).toBe(true)
      expect(hasPermission("admin", "customers", "delete")).toBe(true)
    })
    it("can manage invoices and payments", () => {
      expect(hasPermission("admin", "invoices", "manage")).toBe(true)
      expect(hasPermission("admin", "payments", "manage")).toBe(true)
    })
    it("can manage agents and cash", () => {
      expect(hasPermission("admin", "agents", "manage")).toBe(true)
      expect(hasPermission("admin", "cash", "manage")).toBe(true)
    })
    it("can read reports", () => {
      expect(hasPermission("admin", "reports", "read")).toBe(true)
    })
    it("can read live-monitor", () => {
      expect(hasPermission("admin", "live-monitor", "read")).toBe(true)
    })
    it("can delete routers", () => {
      expect(hasPermission("admin", "routers", "delete")).toBe(true)
    })
  })

  // ─── teknisi: only network management ──────────────────────────────────────
  describe("teknisi", () => {
    it("cannot read users", () => {
      expect(hasPermission("teknisi", "users", "read")).toBe(false)
    })
    it("cannot manage customers", () => {
      expect(hasPermission("teknisi", "customers", "read")).toBe(false)
      expect(hasPermission("teknisi", "customers", "create")).toBe(false)
    })
    it("cannot manage invoices", () => {
      expect(hasPermission("teknisi", "invoices", "read")).toBe(false)
    })
    it("cannot manage payments", () => {
      expect(hasPermission("teknisi", "payments", "read")).toBe(false)
    })
    it("cannot manage agents", () => {
      expect(hasPermission("teknisi", "agents", "read")).toBe(false)
    })
    it("cannot access reports", () => {
      expect(hasPermission("teknisi", "reports", "read")).toBe(false)
    })
    it("can read and update routers", () => {
      expect(hasPermission("teknisi", "routers", "read")).toBe(true)
      expect(hasPermission("teknisi", "routers", "update")).toBe(true)
    })
    it("cannot delete routers", () => {
      expect(hasPermission("teknisi", "routers", "delete")).toBe(false)
    })
    it("can manage subscriptions (isolate/restore/suspend)", () => {
      expect(hasPermission("teknisi", "subscriptions", "manage")).toBe(true)
    })
    it("can manage bandwidth-profiles", () => {
      expect(hasPermission("teknisi", "bandwidth-profiles", "manage")).toBe(true)
    })
    it("can access live-monitor", () => {
      expect(hasPermission("teknisi", "live-monitor", "read")).toBe(true)
    })
    it("can see dashboard", () => {
      expect(hasPermission("teknisi", "dashboard", "read")).toBe(true)
    })
    it("cannot create customers", () => {
      expect(hasPermission("teknisi", "customers", "create")).toBe(false)
    })
  })

  // ─── Edge cases ─────────────────────────────────────────────────────────────
  describe("edge cases", () => {
    it("returns false for unknown resource", () => {
      expect(
        hasPermission("superadmin", "unknown-resource" as never, "read")
      ).toBe(false)
    })
    it("returns false for unknown action on known resource", () => {
      expect(
        hasPermission("superadmin", "dashboard", "delete" as never)
      ).toBe(false)
    })
  })
})

describe("getAccessibleResources", () => {
  it("superadmin gets all resources including users", () => {
    const resources = getAccessibleResources("superadmin")
    expect(resources).toContain("users")
    expect(resources).toContain("customers")
    expect(resources).toContain("routers")
  })

  it("admin gets resources excluding users", () => {
    const resources = getAccessibleResources("admin")
    expect(resources).not.toContain("users")
    expect(resources).toContain("customers")
    expect(resources).toContain("routers")
  })

  it("teknisi gets only network resources", () => {
    const resources = getAccessibleResources("teknisi")
    expect(resources).not.toContain("users")
    expect(resources).not.toContain("customers")
    expect(resources).not.toContain("invoices")
    expect(resources).toContain("routers")
    expect(resources).toContain("subscriptions")
    expect(resources).toContain("bandwidth-profiles")
    expect(resources).toContain("live-monitor")
    expect(resources).toContain("dashboard")
  })
})
