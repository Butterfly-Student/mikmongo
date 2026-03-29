import { StateCreator } from "zustand"
import type { StoreState } from "../types"

export interface AdminUser {
  id: string
  email: string
  role: "superadmin" | "admin" | "teknisi"
  full_name: string
}

export interface AdminAuthSlice {
  adminAccessToken: string | null
  adminRefreshToken: string | null
  adminUser: AdminUser | null
  adminIsAuthenticated: boolean
  adminSetTokens: (access: string, refresh: string) => void
  adminSetUser: (user: AdminUser) => void
  adminClearAuth: () => void
}

export const createAdminAuthSlice: StateCreator<
  StoreState,
  [["zustand/persist", unknown], ["zustand/devtools", never]],
  [],
  AdminAuthSlice
> = (set) => ({
  adminAccessToken: null,
  adminRefreshToken: null,
  adminUser: null,
  adminIsAuthenticated: false,
  adminSetTokens: (access, refresh) =>
    set({
      adminAccessToken: access,
      adminRefreshToken: refresh,
      adminIsAuthenticated: true,
    }),
  adminSetUser: (user) => set({ adminUser: user }),
  adminClearAuth: () =>
    set({
      adminAccessToken: null,
      adminRefreshToken: null,
      adminUser: null,
      adminIsAuthenticated: false,
    }),
})
