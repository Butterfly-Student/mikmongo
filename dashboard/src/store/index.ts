import { create } from "zustand"
import { persist, devtools } from "zustand/middleware"
import { createAdminAuthSlice } from "./slices/adminAuthSlice"
import { createAgentAuthSlice } from "./slices/agentAuthSlice"
import { createCustomerAuthSlice } from "./slices/customerAuthSlice"
import type { StoreState } from "./types"

export type { StoreState }
export type { AdminUser } from "./slices/adminAuthSlice"
export type { AgentUser } from "./slices/agentAuthSlice"
export type { CustomerUser } from "./slices/customerAuthSlice"

export const useStore = create<StoreState>()(
  devtools(
    persist(
      (set, get, api) => ({
        // Hydration flag — gated rendering prevents redirect-on-reload bug
        isHydrated: false,
        setHydrated: (v) => set({ isHydrated: v }),
        ...createAdminAuthSlice(set, get, api),
        ...createAgentAuthSlice(set, get, api),
        ...createCustomerAuthSlice(set, get, api),
      }),
      {
        name: "mikmongo-auth",
        // Only persist tokens and user objects — NOT actions
        partialize: (state) => ({
          adminAccessToken: state.adminAccessToken,
          adminRefreshToken: state.adminRefreshToken,
          adminUser: state.adminUser,
          adminIsAuthenticated: state.adminIsAuthenticated,
          agentAccessToken: state.agentAccessToken,
          agentRefreshToken: state.agentRefreshToken,
          agentUser: state.agentUser,
          agentIsAuthenticated: state.agentIsAuthenticated,
          customerAccessToken: state.customerAccessToken,
          customerRefreshToken: state.customerRefreshToken,
          customerUser: state.customerUser,
          customerIsAuthenticated: state.customerIsAuthenticated,
        }),
        onRehydrateStorage: () => (state) => {
          // Mark hydration complete so route guards don't fire before tokens load
          state?.setHydrated(true)
        },
      }
    ),
    { name: "MikMongo Store", enabled: import.meta.env.DEV }
  )
)

// ─── Vanilla getters — for use OUTSIDE React (Axios interceptors) ───────────
// These call useStore.getState() directly, which is safe outside components.
// Do NOT use hooks (useStore(...)) inside Axios interceptors.

export const getAdminAccessToken = (): string | null =>
  useStore.getState().adminAccessToken

export const getAdminRefreshToken = (): string | null =>
  useStore.getState().adminRefreshToken

export const getAgentAccessToken = (): string | null =>
  useStore.getState().agentAccessToken

export const getAgentRefreshToken = (): string | null =>
  useStore.getState().agentRefreshToken

export const getCustomerAccessToken = (): string | null =>
  useStore.getState().customerAccessToken

export const getCustomerRefreshToken = (): string | null =>
  useStore.getState().customerRefreshToken

export const adminAuthActions = {
  setTokens: (access: string, refresh: string) =>
    useStore.getState().adminSetTokens(access, refresh),
  setUser: (user: Parameters<StoreState["adminSetUser"]>[0]) =>
    useStore.getState().adminSetUser(user),
  clearAuth: () => useStore.getState().adminClearAuth(),
}

export const agentAuthActions = {
  setTokens: (access: string, refresh: string) =>
    useStore.getState().agentSetTokens(access, refresh),
  setUser: (user: Parameters<StoreState["agentSetUser"]>[0]) =>
    useStore.getState().agentSetUser(user),
  clearAuth: () => useStore.getState().agentClearAuth(),
}

export const customerAuthActions = {
  setTokens: (access: string, refresh: string) =>
    useStore.getState().customerSetTokens(access, refresh),
  setUser: (user: Parameters<StoreState["customerSetUser"]>[0]) =>
    useStore.getState().customerSetUser(user),
  clearAuth: () => useStore.getState().customerClearAuth(),
}
