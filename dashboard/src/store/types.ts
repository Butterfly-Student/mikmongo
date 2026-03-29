import type { AdminAuthSlice } from "./slices/adminAuthSlice"
import type { AgentAuthSlice } from "./slices/agentAuthSlice"
import type { CustomerAuthSlice } from "./slices/customerAuthSlice"

export type StoreState = AdminAuthSlice & AgentAuthSlice & CustomerAuthSlice & {
  isHydrated: boolean
  setHydrated: (v: boolean) => void
}
