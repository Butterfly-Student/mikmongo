import { StateCreator } from "zustand"
import type { StoreState } from "../types"

export interface AgentUser {
  id: string
  email: string
  full_name: string
  role: "agent"
}

export interface AgentAuthSlice {
  agentAccessToken: string | null
  agentRefreshToken: string | null
  agentUser: AgentUser | null
  agentIsAuthenticated: boolean
  agentSetTokens: (access: string, refresh: string) => void
  agentSetUser: (user: AgentUser) => void
  agentClearAuth: () => void
}

export const createAgentAuthSlice: StateCreator<
  StoreState,
  [["zustand/persist", unknown], ["zustand/devtools", never]],
  [],
  AgentAuthSlice
> = (set) => ({
  agentAccessToken: null,
  agentRefreshToken: null,
  agentUser: null,
  agentIsAuthenticated: false,
  agentSetTokens: (access, refresh) =>
    set({
      agentAccessToken: access,
      agentRefreshToken: refresh,
      agentIsAuthenticated: true,
    }),
  agentSetUser: (user) => set({ agentUser: user }),
  agentClearAuth: () =>
    set({
      agentAccessToken: null,
      agentRefreshToken: null,
      agentUser: null,
      agentIsAuthenticated: false,
    }),
})
