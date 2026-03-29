import { StateCreator } from "zustand"
import type { StoreState } from "../types"

export interface CustomerUser {
  id: string
  email: string
  full_name: string
  role: "customer"
}

export interface CustomerAuthSlice {
  customerAccessToken: string | null
  customerRefreshToken: string | null
  customerUser: CustomerUser | null
  customerIsAuthenticated: boolean
  customerSetTokens: (access: string, refresh: string) => void
  customerSetUser: (user: CustomerUser) => void
  customerClearAuth: () => void
}

export const createCustomerAuthSlice: StateCreator<
  StoreState,
  [["zustand/persist", unknown], ["zustand/devtools", never]],
  [],
  CustomerAuthSlice
> = (set) => ({
  customerAccessToken: null,
  customerRefreshToken: null,
  customerUser: null,
  customerIsAuthenticated: false,
  customerSetTokens: (access, refresh) =>
    set({
      customerAccessToken: access,
      customerRefreshToken: refresh,
      customerIsAuthenticated: true,
    }),
  customerSetUser: (user) => set({ customerUser: user }),
  customerClearAuth: () =>
    set({
      customerAccessToken: null,
      customerRefreshToken: null,
      customerUser: null,
      customerIsAuthenticated: false,
    }),
})
