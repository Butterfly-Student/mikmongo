import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import type { AdminUser, CustomerUser, AgentUser } from '@/api/types'

// --- Admin Auth Slice ---
interface AdminAuthSlice {
  adminAccessToken: string | null
  adminRefreshToken: string | null
  adminUser: AdminUser | null
  adminIsAuthenticated: boolean
  adminSetTokens: (accessToken: string, refreshToken: string) => void
  adminSetUser: (user: AdminUser) => void
  adminClearAuth: () => void
}

// --- Customer Auth Slice ---
interface CustomerAuthSlice {
  customerToken: string | null
  customerUser: CustomerUser | null
  customerIsAuthenticated: boolean
  customerSetToken: (token: string) => void
  customerSetUser: (user: CustomerUser) => void
  customerClearAuth: () => void
}

// --- Agent Auth Slice ---
interface AgentAuthSlice {
  agentToken: string | null
  agentUser: AgentUser | null
  agentIsAuthenticated: boolean
  agentSetToken: (token: string) => void
  agentSetUser: (user: AgentUser) => void
  agentClearAuth: () => void
}

// --- Hydration Slice ---
interface HydrationSlice {
  isHydrated: boolean
  setHydrated: () => void
}

type AuthStoreState = AdminAuthSlice & CustomerAuthSlice & AgentAuthSlice & HydrationSlice

export const useAuthStore = create<AuthStoreState>()(
  persist(
    (set) => ({
      // Admin
      adminAccessToken: null,
      adminRefreshToken: null,
      adminUser: null,
      adminIsAuthenticated: false,
      adminSetTokens: (accessToken, refreshToken) =>
        set((state) => ({
          ...state,
          adminAccessToken: accessToken,
          adminRefreshToken: refreshToken,
          adminIsAuthenticated: true,
        })),
      adminSetUser: (user) =>
        set((state) => ({ ...state, adminUser: user })),
      adminClearAuth: () =>
        set((state) => ({
          ...state,
          adminAccessToken: null,
          adminRefreshToken: null,
          adminUser: null,
          adminIsAuthenticated: false,
        })),

      // Customer
      customerToken: null,
      customerUser: null,
      customerIsAuthenticated: false,
      customerSetToken: (token) =>
        set((state) => ({
          ...state,
          customerToken: token,
          customerIsAuthenticated: true,
        })),
      customerSetUser: (user) =>
        set((state) => ({ ...state, customerUser: user })),
      customerClearAuth: () =>
        set((state) => ({
          ...state,
          customerToken: null,
          customerUser: null,
          customerIsAuthenticated: false,
        })),

      // Agent
      agentToken: null,
      agentUser: null,
      agentIsAuthenticated: false,
      agentSetToken: (token) =>
        set((state) => ({
          ...state,
          agentToken: token,
          agentIsAuthenticated: true,
        })),
      agentSetUser: (user) =>
        set((state) => ({ ...state, agentUser: user })),
      agentClearAuth: () =>
        set((state) => ({
          ...state,
          agentToken: null,
          agentUser: null,
          agentIsAuthenticated: false,
        })),

      // Hydration
      isHydrated: false,
      setHydrated: () => set((state) => ({ ...state, isHydrated: true })),
    }),
    {
      name: 'mikmongo-auth',
      partialize: (state) => ({
        adminAccessToken: state.adminAccessToken,
        adminRefreshToken: state.adminRefreshToken,
        adminUser: state.adminUser,
        adminIsAuthenticated: state.adminIsAuthenticated,
        customerToken: state.customerToken,
        customerUser: state.customerUser,
        customerIsAuthenticated: state.customerIsAuthenticated,
        agentToken: state.agentToken,
        agentUser: state.agentUser,
        agentIsAuthenticated: state.agentIsAuthenticated,
      }),
      onRehydrateStorage: () => (state) => {
        state?.setHydrated()
      },
    }
  )
)
